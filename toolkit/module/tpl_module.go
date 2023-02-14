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
func {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}},data dao.{{CamelStr .Table.TableName}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Update(data)
}

// {{CamelStr .Table.TableName}}Item 获取数据
func {{CamelStr .Table.TableName}}Item(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(dao.{{CamelStr .Table.TableName}},error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{CamelStr .Table.TableName}}
	err := db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Row().ToStruct(&res)
	return res, err
}

// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func {{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(int64,error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Delete()
}
{{- range .Method}}
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
