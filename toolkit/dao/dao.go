package dao

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"text/template"
	"time"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// PERF-8-1 [P1] DAO 模板对所有表完全相同：模板字符串读取、FuncMap 注册与 template.Parse
// 在整个命令生命周期内只做一次并缓存复用，避免批量生成时每张表重复完成模板词法/语法编译。
var (
	daoTplOnce sync.Once
	daoTpl     *template.Template
)

// loadDaoTemplate 懒加载并缓存预解析完成的 DAO 模板；解析失败保持原有 template.Must 直接 panic 语义。
func loadDaoTemplate() *template.Template {
	daoTplOnce.Do(func() {
		daoTpl = template.Must(template.New("name").Funcs(template.FuncMap{
			"CamelStr":   base.CamelStr,
			"Helper":     base.Helper,
			"SymbolChar": base.SymbolChar,
			"Pointer":    base.Pointer,
		}).Parse(DaoTemplate()))
	})
	return daoTpl
}

// CmdNew 是 toolkit dao 子命令，按数据库表清单生成 GORM 模型结构体代码。
var (
	CmdNew = &cobra.Command{
		Use:   "dao",
		Short: "数据访问对象",
		Long:  "创建数据访问对象: toolkit dao",
		Run:   Run,
	}
)

// Run 执行 dao 子命令：初始化运行环境后读取配置数据库的表清单，命令行首个参数可指定
// 单表或以逗号分隔的多表进行过滤；逐表读取字段元数据并调用 GenerateDao 生成模型文件，
// 最后在 dao 目录执行 go fmt 与 goimports 整理格式。
//
// 环境初始化或表清单读取失败时打印错误并返回；单表字段读取失败时跳过该表，不中断其余表。
func Run(cmd *cobra.Command, args []string) {
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	dir := path.Join(cast.ToString(base.Path), "dao")
	_ = os.MkdirAll(dir, os.ModePerm)

	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	tableList, err := base.TableList(ctx, base.Config.String("db.Database"))
	if err != nil {
		fmt.Println("数据库信息获取:", color.RedString(err.Error()))
		return
	}
	tables := make([]string, 0)
	if len(args) > 0 {
		tableName := args[0]
		if tableName != "" {
			tables = util.Explode(",", tableName)
		}
	}

	// PERF-8-7 [P2] 命令行指定的表名预建集合，循环内 O(1) 命中判断，替代反射型 InArray 线性查找。
	tablesSet := make(map[string]struct{}, len(tables))
	for _, name := range tables {
		tablesSet[name] = struct{}{}
	}

	// PERF-8-1 [P1] 命令启动后即完成一次模板读取与预解析，避免逐表重复解析。
	loadDaoTemplate()

	// PERF-8-18 [P2] 跨表 TableColumn 查询受限并发（预期收益需实测验证）：各表查询互不依赖，
	// 4 路信号量并发取字段；单表失败保持原有跳过、不中断其余表的语义；结果按下标回填，
	// 后续渲染与写文件仍按原表顺序串行，生成产物顺序与内容保持不变。
	filtered := make([]base.Table, 0, len(tableList))
	for _, table := range tableList {
		if !util.Empty(tables) {
			if _, ok := tablesSet[table.TableName]; !ok {
				continue
			}
		}
		filtered = append(filtered, table)
	}

	type daoColumns struct {
		columns []base.Column
		ok      bool
	}
	results := make([]daoColumns, len(filtered))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, table := range filtered {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, tbl base.Table) {
			defer wg.Done()
			defer func() { <-sem }()
			column, qerr := base.TableColumn(ctx, base.Config.String("db.Database"), tbl.TableName)
			if qerr != nil {
				return
			}
			results[idx] = daoColumns{columns: column, ok: true}
		}(i, table)
	}
	wg.Wait()

	for i, table := range filtered {
		if !results[i].ok {
			continue
		}
		GenerateDao(table, results[i].columns)
	}
	_ = os.Chdir(dir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(dir, "*.go"))
	cmdImport.CombinedOutput()
}

// GenerateDao 按 DaoTemplate 渲染一张表的 GORM 模型文件，写入项目 dao 目录下以表名
// 命名的 Go 文件；同名文件已存在时先删除再创建。
//
// 文件创建或模板解析失败时打印错误并返回，不向调用方返回错误。
func GenerateDao(table base.Table, column []base.Column) {
	filePath := path.Join(cast.ToString(base.Path), "dao", table.TableName+".go")
	if util.FileExists(filePath) {
		util.Delete(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	// PERF-8-1 [P1] 复用命令级预解析缓存，不再逐表重建 FuncMap 与重新 Parse。
	tpl := loadDaoTemplate()

	data := base.DaoParam{
		Table:       table,
		TableColumn: column,
	}
	err = tpl.Execute(file, data)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(filePath))
}

// DaoTemplate 返回 DAO 模型结构体代码模板字符串。
//
// 配置项 template.Dao 指定了自定义模板文件且文件存在时，读取该文件内容并返回；
// 否则返回包内置模板，模板占位符由 GenerateDao 渲染填充。
func DaoTemplate() string {
	if exists := base.Config.Exists("template.Dao"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.Dao"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
package dao

import "github.com/abulo/ratel/v3/stores/null"

// {{CamelStr .Table.TableName}} {{.Table.TableComment}} {{.Table.TableName}}
type {{CamelStr .Table.TableName}} struct {
	{{- range .TableColumn }}
	{{- if eq .IsNullable "YES" }}
	{{CamelStr .ColumnName}}	{{Pointer .DataTypeMap.Empty}}{{.DataTypeMap.Empty}}	{{SymbolChar}}gorm:"column:{{.ColumnName}}" json:"{{Helper .ColumnName}}"{{SymbolChar}}  //{{.DataType}} {{.ColumnComment}}
	{{- else }}
	{{- if eq .ColumnKey "PRI" }}
	{{CamelStr .ColumnName}}	*{{.DataTypeMap.Default}}	{{SymbolChar}}gorm:"column:{{.ColumnName}}" json:"{{Helper .ColumnName}}"{{SymbolChar}}  //{{.DataType}} {{.ColumnComment}},PRI
	{{- else }}
	{{CamelStr .ColumnName}}	{{Pointer .DataTypeMap.Default}}{{.DataTypeMap.Default}}	{{SymbolChar}}gorm:"column:{{.ColumnName}}" json:"{{Helper .ColumnName}}"{{SymbolChar}}  //{{.DataType}} {{.ColumnComment}}
	{{- end}}
	{{- end}}
	{{- end}}
}

func ({{CamelStr .Table.TableName}}) TableName() string {
	return "{{.Table.TableName}}"
}
`
	return outString
}
