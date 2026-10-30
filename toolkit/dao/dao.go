package dao

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

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
	for _, table := range tableList {
		if !util.Empty(tables) {
			if !util.InArray(table.TableName, tables) {
				continue
			}
		}
		column, err := base.TableColumn(ctx, base.Config.String("db.Database"), table.TableName)
		if err != nil {
			continue
		}
		GenerateDao(table, column)
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
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{
		"CamelStr":   base.CamelStr,
		"Helper":     base.Helper,
		"SymbolChar": base.SymbolChar,
		"Pointer":    base.Pointer,
	}).Parse(DaoTemplate()))

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
