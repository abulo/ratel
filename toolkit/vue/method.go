package vue

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
)

// ApiList 返回前端各操作对应的请求路径，key 取值 add、update、item、delete、drop、
// recover、list；路径前缀取自 base.Url（由 GenerateMethod 通过 SetUrl 设置）。
// 未识别的 key 返回空字符串。
func ApiList(key string) string {
	apiUrl := cast.ToString(base.Url)
	list := make(map[string]string)
	list["add"] = base.PathChar(apiUrl)
	list["update"] = base.PathChar(apiUrl, "/${id}/update")
	list["item"] = base.PathChar(apiUrl, "/${id}/item")
	list["delete"] = base.PathChar(apiUrl, "/${id}/delete")
	list["drop"] = base.PathChar(apiUrl, "/${id}/drop")
	list["recover"] = base.PathChar(apiUrl, "/${id}/recover")
	list["list"] = base.PathChar(apiUrl)
	if val, ok := list[key]; ok {
		return val
	}
	return ""
}

// GenerateMethod 按 MethodTemplate 渲染前端 API 请求方法模块，写入 fullMethodDir 下以
// 表名命名的 .ts 文件；同名文件已存在时先删除再创建。apiUrl 作为接口前缀写入 base.Url
// 供 ApiList 使用。
//
// 文件创建或模板解析失败时打印错误并返回，不向调用方返回错误。
func GenerateMethod(moduleParam base.ModuleParam, apiUrl, fullMethodDir, tableName string) {
	base.SetUrl(apiUrl)
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
		"ApiList":               ApiList,
	}).Parse(MethodTemplate()))

	outMethodFile := path.Join(fullMethodDir, base.Helper(tableName)+".ts")
	if util.FileExists(outMethodFile) {
		util.Delete(outMethodFile)
	}

	file, err := os.OpenFile(outMethodFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}

	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outMethodFile))
}

// MethodTemplate 返回前端 API 方法模块模板字符串。
//
// 配置项 template.VueMethod 指定了自定义模板文件且文件存在时，读取该文件内容并返回；
// 否则返回包内置模板，模板占位符由 GenerateMethod 渲染填充。
func MethodTemplate() string {
	if exists := base.Config.Exists("template.VueMethod"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.VueMethod"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
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
  return http.post(PORT + {{ApiList "add"}}, params);
};
{{- else if eq .Type "Update"}}
// {{.Table.TableComment}}更新数据
export const update{{CamelStr .Table.TableName}}Api = (id: number, params: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  return http.put(PORT + {{ApiList "update"}}, params);
};
{{- else if eq .Type "Show"}}
// {{.Table.TableComment}}查询单条数据
export const get{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.get<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>(PORT + {{ApiList "item"}});
};
{{- else if eq .Type "Delete"}}
// {{.Table.TableComment}}删除数据
export const delete{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.delete(PORT + {{ApiList "delete"}});
};
{{- else if eq .Type "Drop"}}
// {{.Table.TableComment}}清理数据
export const drop{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.delete(PORT + {{ApiList "drop"}});
};
{{- else if eq .Type "Recover"}}
// {{.Table.TableComment}}恢复数据
export const recover{{CamelStr .Table.TableName}}Api = (id: number) => {
  return http.put(PORT + {{ApiList "recover"}});
};
{{- else if eq .Type "List"}}
// {{.Table.TableComment}}列表数据
{{- if .Page}}
export const get{{CamelStr .Table.TableName}}ListApi = (params?: {{CamelStr .Table.TableName}}.Req{{CamelStr .Table.TableName}}List) => {
  return http.get<ResPage<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>>(PORT + {{ApiList "list"}}, params);
};
{{- else}}
export const get{{CamelStr .Table.TableName}}ListApi = (params?: {{CamelStr .Table.TableName}}.Req{{CamelStr .Table.TableName}}List) => {
  return http.get<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item[]>(PORT + {{ApiList "list"}}, params);
};
{{- end}}
{{- end}}
{{- end}}`
	return outString
}
