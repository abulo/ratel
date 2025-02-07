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

var (
	CmdNew = &cobra.Command{
		Use:   "dao",
		Short: "数据访问对象",
		Long:  "创建数据访问对象: toolkit dao",
		Run:   Run,
	}
)

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	// 创建文件夹
	dir := path.Join(cast.ToString(base.Path), "dao")
	_ = os.MkdirAll(dir, os.ModePerm)

	// 初始化上下文
	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	//获取表信息
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
	//格式化代码
	_ = os.Chdir(dir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(dir, "*.go"))
	cmdImport.CombinedOutput()
}

// GenerateDao
func GenerateDao(table base.Table, column []base.Column) {
	filePath := path.Join(cast.ToString(base.Path), "dao", table.TableName+".go")
	//存在文件,需要先将文件删除掉
	if util.FileExists(filePath) {
		util.Delete(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	//go文件生成地址
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{
		"CamelStr":   base.CamelStr,
		"Helper":     base.Helper,
		"SymbolChar": base.SymbolChar,
		"Pointer":    base.Pointer,
	}).Parse(DaoTemplate()))

	//定义结构体接收数据
	data := base.DaoParam{
		Table:       table,
		TableColumn: column,
	}
	//渲染输出
	err = tpl.Execute(file, data)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(filePath))
}

// DaoTemplate 模板
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

// package dao

// import "github.com/abulo/ratel/v3/stores/null"

// // {{CamelStr .Table.TableName}} {{.Table.TableComment}} {{.Table.TableName}}
// type {{CamelStr .Table.TableName}} struct {
// 	{{- range .TableColumn }}
// 	{{- if eq .IsNullable "YES" }}
// 	{{CamelStr .ColumnName}}	{{.DataTypeMap.Empty}}	{{SymbolChar}}db:"{{.ColumnName}}" json:"{{Helper .ColumnName}}" form:"{{Helper .ColumnName}}" uri:"{{Helper .ColumnName}}" xml:"{{Helper .ColumnName}}" proto:"{{Helper .ColumnName}}"{{SymbolChar}}  //{{.DataType}} {{.ColumnComment}}
// 	{{- else }}
// 	{{CamelStr .ColumnName}}	{{.DataTypeMap.Default}}	{{SymbolChar}}db:"{{.ColumnName}}" json:"{{Helper .ColumnName}}" form:"{{Helper .ColumnName}}" uri:"{{Helper .ColumnName}}" xml:"{{Helper .ColumnName}}" proto:"{{Helper .ColumnName}}"{{SymbolChar}}  //{{.DataType}} {{.ColumnComment}}
// 	{{- end}}
// 	{{- end}}
// }
