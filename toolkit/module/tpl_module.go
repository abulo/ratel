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
	"github.com/spf13/cast"
)

// GenerateModule 使用 ModuleTemplate 渲染并写出数据模型层 Go 文件到 fullModuleDir，
// 文件已存在时先删除再生成，随后执行 go fmt 与 goimports 整理格式；模板渲染失败时打印错误并返回。
func GenerateModule(moduleParam base.ModuleParam, fullModuleDir, tableName string) {
	// PERF-8-20 [P2] 公共模板函数复用 base 包级定义，避免每次调用重复构建同一 FuncMap 字面量。
	funcs := base.CommonTemplateFuncs()
	tpl := template.Must(template.New("module").Funcs(funcs).Parse(ModuleTemplate()))
	outModuleFile := path.Join(fullModuleDir, tableName+".go")
	if util.FileExists(outModuleFile) {
		util.Delete(outModuleFile)
	}
	file, err := os.OpenFile(outModuleFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
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

// ModuleTemplate 返回数据模型层代码模板；配置中 template.Module 指定了自定义模板文件时优先读取该文件内容。
func ModuleTemplate() string {
	if exists := base.Config.Exists("template.Module"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.Module"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
package {{.Pkg}}

import (
	"context"
	"{{.ModName}}/dao"
	"{{.ModName}}/initial"

	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
	{{- if .Page}}
	"github.com/spf13/cast"
	{{- end}}
)
// {{.Table.TableName}} {{.Table.TableComment}}

{{- range .Method}}
{{- if eq .Type "Create"}}
// {{.Name}} 创建数据
func {{.Name}}(ctx context.Context,data dao.{{CamelStr .Table.TableName}})(res int64, err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Insert(data)
	if err != nil {
		return
	}
	res, err = db.Insert(ctx ,query,args...)
	return
}
{{- else if eq .Type "Update"}}
// {{.Name}} 更新数据
func {{.Name}}(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}},data dao.{{CamelStr .Table.TableName}})(res int64, err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Update(data)
	if err != nil {
		return
	}
	res, err = db.Update(ctx ,query,args...)
	return
}
{{- else if eq .Type "Show"}}
// {{.Name}} 查询单条数据
func {{.Name}}(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res dao.{{CamelStr .Table.TableName}},err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Row()
	if err != nil {
		return
	}
	err = db.QueryRow(ctx ,query,args...).ToStruct(&res)
	return
}
{{- else if eq .Type "Drop"}}
// {{.Name}} 清理数据
func {{.Name}}(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Delete()
	if err != nil {
		return
	}
	res, err = db.Delete(ctx ,query,args...)
	return
}
{{- else if eq .Type "Delete"}}
// {{.Name}} 删除数据
func {{.Name}}(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	{{- if .SoftDelete}}
	data := make(map[string]any)
	data["deleted"] = 1
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Update(data)
	if err != nil {
		return
	}
	res, err = db.Update(ctx ,query,args...)
	{{- else}}
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Delete()
	if err != nil {
		return
	}
	res, err = db.Delete(ctx ,query,args...)
	{{- end}}
	return
}
{{- else if eq .Type "Recover"}}
{{- if .SoftDelete}}
// {{.Name}} 恢复数据
func {{.Name}}(ctx context.Context,{{Helper .Primary.AlisaColumnName}} {{.Primary.DataTypeMap.Default}})(res int64,err error){
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder := sql.NewBuilder()
	data := make(map[string]any)
	data["deleted"] = 0
	query,args,err := builder.Table("{{Char .Table.TableName}}").Where("{{Char .Primary.ColumnName}}",{{Helper .Primary.AlisaColumnName}}).Update(data)
	if err != nil {
		return
	}
	res, err = db.Update(ctx ,query,args...)
	return
}
{{- end}}
{{- else if eq .Type "Item"}}
// {{.Name}} 查询单条数据
func {{.Name}}(ctx context.Context,condition map[string]any)(res dao.{{CamelStr .Table.TableName}},err error){
	if util.Empty(condition) {
		err = errors.New("condition is empty")
		return
	}
	db := initial.Core.Store.LoadSQL("mysql").Read()
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
{{- else if eq .Type "List"}}
// {{.Name}} 查询列表数据
func {{.Name}}(ctx context.Context,condition map[string]any)(res []dao.{{CamelStr .Table.TableName}}, err error){
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder := sql.NewBuilder()
	builder.Table("{{Char .Table.TableName}}")
	{{Convert .Condition}}
	{{- if .Page}}
	if val, ok := condition["pagination"]; ok {
		pagination := val.(*sql.Pagination)
		if pagination != nil {
			builder.Offset(pagination.GetOffset())
			builder.Limit(pagination.GetLimit())
		}
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
// {{.Name}}Total 查询列表数据总量
func {{.Name}}Total(ctx context.Context,condition map[string]any)(res int64,err error){
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
{{- end}}
`
	return outString
}
