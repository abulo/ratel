package module

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// CmdNew represents the new command.
	CmdNew = &cobra.Command{
		Use:   "module",
		Short: "数据模型层",
		Long:  "创建数据库模型层: toolkit module dir table_name",
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
	dir := path.Join(base.Path, "module")
	_ = os.MkdirAll(dir, os.ModePerm)

	//创建数据
	moduleDir := ""
	tableName := ""
	if len(args) == 0 {
		if err := survey.AskOne(&survey.Input{
			Message: "模型路径",
			Help:    "文件夹路径",
		}, &moduleDir); err != nil || moduleDir == "" {
			return
		}
		if err := survey.AskOne(&survey.Input{
			Message: "表名称",
			Help:    "数据库中某个表名称",
		}, &tableName); err != nil || tableName == "" {
			return
		}
	} else {
		moduleDir = args[0]
		tableName = args[1]
	}
	if tableName == "" || moduleDir == "" {
		fmt.Println("初始化:", color.RedString("模型层名称 & 表名称 必须填写"))
		return
	}
	// 文件夹的路径
	fullDir := path.Join(base.Path, "module", moduleDir)
	_ = os.MkdirAll(fullDir, os.ModePerm)

	// 初始化上下文
	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// 表结构信息
	tableColumn, err := base.TableColumn(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表结构信息:", color.RedString(err.Error()))
		return
	}
	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}
	// 表信息
	tableItem, err := base.TableItem(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表信息:", color.RedString(err.Error()))
		return
	}
	// 表索引
	tableIndex, err := base.TableIndex(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表索引:", color.RedString(err.Error()))
		return
	}
	// 表主键
	tablePrimary, err := base.TablePrimary(ctx, base.Config.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表主键:", color.RedString(err.Error()))
		return
	}
	var methodList []base.Method

	//获取的索引信息没有
	if err != nil {
		method := base.Method{
			Table:       tableItem,
			TableColumn: tableColumn,
			Type:        "List",
			Name:        "List",
			Default:     true,
			Condition:   nil,
			Primary:     tablePrimary,
		}
		methodList = append(methodList, method)
	} else {
		//存储条件信息
		field := make([]string, 0)
		//有索引信息
		for _, v := range tableIndex {
			//查询条件
			condition := make([]base.Column, 0)
			//数据库索引
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
				//构造查询条件
				condition = append(condition, tableColumnMap[fieldValue])
				if !util.InArray(fieldValue, field) {
					field = append(field, fieldValue)
				}
			}
			// 数据库中的索引名称
			indexName := v.IndexName
			// 拆分字符串,得到索引类型和索引名称
			indexNameSlice := util.Explode(":", indexName)
			if len(indexNameSlice) < 2 {
				continue
			}
			// 自定义函数名称和索引信息
			customIndexType := util.UCWords(indexNameSlice[0])
			customIndexName := util.UCWords(indexNameSlice[1])
			method := base.Method{
				Table:       tableItem,
				TableColumn: tableColumn,
				Type:        customIndexType,
				Name:        customIndexName,
				Default:     false,
				Condition:   condition,
				Primary:     tablePrimary,
			}
			//添加到集合中
			methodList = append(methodList, method)
		}
		condition := make([]base.Column, 0)
		for _, fieldValue := range field {
			//构造查询条件
			condition = append(condition, tableColumnMap[fieldValue])
		}
		method := base.Method{
			Table:       tableItem,
			TableColumn: tableColumn,
			Type:        "List",
			Name:        "List",
			Default:     true,
			Condition:   condition,
			Primary:     tablePrimary,
		}
		methodList = append(methodList, method)
	}

	//获取 go.mod
	mod, err := base.ModulePath(path.Join(base.Path, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		mod = "test"
	}
	// 数字长度
	strLen := strings.LastIndex(moduleDir, "/")
	// 数据模型
	moduleParam := base.ModuleParam{
		Pkg:         moduleDir[strLen+1:],
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      methodList,
		ModName:     mod,
	}
	// 模板变量
	tpl := template.Must(template.New("name").Funcs(template.FuncMap{
		"Convert":    base.Convert,
		"SymbolChar": base.SymbolChar,
		"Char":       base.Char,
		"Helper":     base.Helper,
		"CamelStr":   base.CamelStr,
	}).Parse(ModuleTemplate()))
	// 文件夹路径
	outFile := path.Join(fullDir, tableName+".go")
	if util.FileExists(outFile) {
		util.Delete(outFile)
	}
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	//渲染输出
	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	_ = os.Chdir(fullDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outFile))
}

// ModuleTemplate 模板
func ModuleTemplate() string {
	outString := `
package {{.Pkg}}

import (
	"context"
	"{{.ModName}}/dao"
	"{{.ModName}}/initial"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)
// {{.Table.TableName}} {{.Table.TableComment}}


// {{CamelStr .Table.TableName}}ItemCreate 创建数据
func {{CamelStr .Table.TableName}}ItemCreate(ctx context.Context,data dao.{{CamelStr .Table.TableName}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Insert(data)
}

// {{CamelStr .Table.TableName}}ItemUpdate 更新数据
func {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}},data dao.{{CamelStr .Table.TableName}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Update(data)
}

// {{CamelStr .Table.TableName}}Item 获取数据
func {{CamelStr .Table.TableName}}Item(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}})(dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{CamelStr .Table.TableName}}
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Row().ToStruct(&res)
}

// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func {{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,{{.Primary.ColumnName}} {{.Primary.DataTypeMap.Default}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{.Primary.ColumnName}}).Delete()
}
{{- range .Method}}
//多条数据
{{- if eq .Type "List"}}
{{- if .Default}}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})([]dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Char .Primary.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total(ctx context.Context,condition map[string]interface{})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	return builder.Count()
}
{{- else}}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})([]dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	if !util.Empty(condition["pageOffset"]) {
		builder.Offset(cast.ToInt64(condition["pageOffset"]))
	}
	if !util.Empty(condition["pageSize"]) {
		builder.Limit(cast.ToInt64(condition["pageSize"]))
	}
	err := builder.OrderBy("{{Char .Primary.ColumnName}}", query.DESC).Rows().ToStruct(&res)
	return res, err
}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total(ctx context.Context,condition map[string]interface{})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	return builder.Count()
}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}} 单列数据
func {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}(ctx context.Context,condition map[string]interface{})(dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{CamelStr .Table.TableName}}
	builder := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	err := builder.Row().ToStruct(&res)
	return res, err
}
{{- end}}
{{- end}}
`
	return outString
}
