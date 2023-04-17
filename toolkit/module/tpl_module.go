package module

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
)

func GenerateModule(moduleParam base.ModuleParam, fullModuleDir, tableName string) {
	// 模板变量
	tpl := template.Must(template.New("module").Funcs(template.FuncMap{
		"Convert":               base.Convert,
		"SymbolChar":            base.SymbolChar,
		"Char":                  base.Char,
		"Helper":                base.Helper,
		"CamelStr":              base.CamelStr,
		"Add":                   base.Add,
		"ModuleProtoConvertDao": base.ModuleProtoConvertDao,
		"ModuleDaoConvertProto": base.ModuleDaoConvertProto,
		"ModuleProtoConvertMap": base.ModuleProtoConvertMap,
		"ApiToProto":            base.ApiToProto,
	}).Parse(ModuleTemplate()))
	// 文件夹路径
	outModuleFile := path.Join(fullModuleDir, tableName+".go")
	if util.FileExists(outModuleFile) {
		util.Delete(outModuleFile)
	}
	file, err := os.OpenFile(outModuleFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	_ = os.Chdir(fullModuleDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullModuleDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outModuleFile))
}

// ModuleTemplate 模板
func ModuleTemplate() string {
	outString := `
package {{.Pkg}}

import (
	"context"
	"{{.ModName}}/dao"
	"{{.ModName}}/initial"

	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/abulo/ratel/v3/util"
	{{- if .Page}}
	"github.com/spf13/cast"
	{{- end}}
)
// {{.Table.TableName}} {{.Table.TableComment}}


// {{CamelStr .Table.TableName}}ItemCreate 创建数据
func {{CamelStr .Table.TableName}}ItemCreate(ctx context.Context,data dao.{{CamelStr .Table.TableName}})(res int64, err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Insert(data)
	if err != nil {
		return
	}
	res, err = db.Insert(ctx ,query,args...)
	return
}

// {{CamelStr .Table.TableName}}ItemUpdate 更新数据
func {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}},data dao.{{CamelStr .Table.TableName}})(res int64, err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Update(data)
	if err != nil {
		return
	}
	res, err = db.Update(ctx ,query,args...)
	return
}

// {{CamelStr .Table.TableName}}Item 获取数据
func {{CamelStr .Table.TableName}}Item(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res dao.{{CamelStr .Table.TableName}},err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	// var res dao.{{CamelStr .Table.TableName}}
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Row()
	if err != nil {
		return
	}
	err = db.QueryRow(ctx ,query,args...).ToStruct(&res)
	return
}

// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func {{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Delete()
	if err != nil {
		return
	}
	res, err = db.Delete(ctx ,query,args...)
	return
}
{{- range .Method}}
{{- if eq .Type "List"}}
{{- if .Default}}

// {{CamelStr .Table.TableName}}{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx context.Context,condition map[string]any)(res []dao.{{CamelStr .Table.TableName}}, err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	// var res []dao.{{CamelStr .Table.TableName}}
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	{{- if .Page}}
	if !util.Empty(condition["offset"]) {
		builder.Offset(cast.ToInt64(condition["offset"]))
	}
	if !util.Empty(condition["limit"]) {
		builder.Limit(cast.ToInt64(condition["limit"]))
	}
	{{- end}}
	builder.OrderBy("{{Char .Primary.ColumnName}}", sql.DESC)
	query,args,err := builder.Rows()
	if err != nil {
		return
	}
	err = db.QueryRows(ctx ,query,args...).ToStruct(&res)
	return
}

{{- if .Page}}
// {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}Total(ctx context.Context,condition map[string]any)(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	query,args,err := builder.Count()
	if err != nil {
		return
	}
	res,err = db.Count(ctx ,query,args...)
	return
}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}(ctx context.Context,condition map[string]any)(res []dao.{{CamelStr .Table.TableName}},err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	// var res []dao.{{CamelStr .Table.TableName}}
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	{{- if .Page}}
	if !util.Empty(condition["offset"]) {
		builder.Offset(cast.ToInt64(condition["offset"]))
	}
	if !util.Empty(condition["limit"]) {
		builder.Limit(cast.ToInt64(condition["limit"]))
	}
	{{- end}}
	query,args,err := builder.OrderBy("{{Char .Primary.ColumnName}}", sql.DESC).Rows()
	if err != nil {
		return
	}
	err = db.QueryRows(ctx ,query,args...).ToStruct(&res)
	return
}

{{- if .Page}}
// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total 列表数据总量
func {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total(ctx context.Context,condition map[string]any)(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	query,args,err := builder.Count()
	if err != nil {
		return
	}
	res,err = db.Count(ctx ,query,args...)
	return
}
{{- end}}
{{- end}}
{{- else}}

// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}} 单列数据
func {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}(ctx context.Context,condition map[string]any)(res dao.{{CamelStr .Table.TableName}},err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	// var res dao.{{CamelStr .Table.TableName}}
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	query,args,err := builder.Row()
	if err != nil {
		return
	}
	err = db.QueryRow(ctx ,query,args...).ToStruct(&res)
	return
}
{{- end}}
{{- end}}
`
	return outString
}
