package vue

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
)

func GenerateMethod(moduleParam base.ModuleParam, apiUrl, fullMethodDir, tableName string) {
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
	}).Parse(MethodTemplate(apiUrl)))

	// 文件夹路径
	outMethodFile := path.Join(fullMethodDir, base.Helper(tableName)+".ts")
	if util.FileExists(outMethodFile) {
		util.Delete(outMethodFile)
	}

	file, err := os.OpenFile(outMethodFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outMethodFile))
}

func MethodTemplate(apiUrl string) string {
	addApiUrl := base.SymbolChar() + apiUrl + base.SymbolChar()
	updateApiUrl := base.SymbolChar() + apiUrl + "/${id}/update" + base.SymbolChar()
	showApiUrl := base.SymbolChar() + apiUrl + "/${id}/item" + base.SymbolChar()
	deleteApiUrl := base.SymbolChar() + apiUrl + "/${id}/delete" + base.SymbolChar()
	recoverApiUrl := base.SymbolChar() + apiUrl + "/${id}/recover" + base.SymbolChar()
	listApiUrl := base.SymbolChar() + apiUrl + base.SymbolChar()
	outString := `
// {{.Table.TableName}} {{.Table.TableComment}}
{{- if .Page}}
import { ResPage } from "@/api/interface/index";
{{- end}}
import { PORT } from "@/api/config/servicePort";
import http from "@/api";
import { {{CamelStr .Table.TableName}} } from "@/api/interface/{{Helper .Table.TableName}}";
{{- range .Method}}
{{- if eq .Type "Create"}}
// {{.Table.TableComment}}创建数据
export const add{{CamelStr .Table.TableName}}Api = (params: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  return http.post(PORT + ` + addApiUrl + `, params);
};
{{- else if eq .Type "Update"}}
// {{.Table.TableComment}}更新数据
export const update{{CamelStr .Table.TableName}}Api = (id: number, params: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  return http.put(PORT + ` + updateApiUrl + `, params);
};
{{- else if eq .Type "Show"}}
// {{.Table.TableComment}}查询单条数据
export const get{{CamelStr .Table.TableName}}ItemApi = (id: number) => {
  return http.get<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>(PORT + ` + showApiUrl + `);
};
{{- else if eq .Type "Delete"}}
// {{.Table.TableComment}}删除数据
export const delete{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.delete(PORT + ` + deleteApiUrl + `);
};
{{- else if eq .Type "Recover"}}
// {{.Table.TableComment}}恢复数据
export const recover{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.put(PORT + ` + recoverApiUrl + `);
};
{{- else if eq .Type "List"}}
// {{.Table.TableComment}}列表数据
{{- if .Page}}
export const get{{CamelStr .Table.TableName}}ListApi = (params?: {{CamelStr .Table.TableName}}.Req{{CamelStr .Table.TableName}}List) => {
  return http.get<ResPage<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>>(PORT + ` + listApiUrl + `, params);
};
{{- else}}
export const get{{CamelStr .Table.TableName}}ListApi = (params?: {{CamelStr .Table.TableName}}.Req{{CamelStr .Table.TableName}}List) => {
  return http.get<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item[]>(PORT + ` + listApiUrl + `, params);
};
{{- end}}
{{- end}}
{{- end}}`
	return outString
}
