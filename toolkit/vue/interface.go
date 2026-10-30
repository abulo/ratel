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

// GenerateInterface 按 InterfaceTemplate 渲染 TypeScript 接口定义文件，写入
// fullInterfaceDir 下以表名命名的 .ts 文件；同名文件已存在时先删除再创建。
//
// 文件创建或模板解析失败时打印错误并返回，不向调用方返回错误。
func GenerateInterface(moduleParam base.ModuleParam, fullInterfaceDir, tableName string) {
	tpl := template.Must(template.New("interface").Funcs(template.FuncMap{
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
		"TypeScriptCondition":   base.TypeScriptCondition,
		"TypeScript":            base.TypeScript,
	}).Parse(InterfaceTemplate()))

	outInterfaceFile := path.Join(fullInterfaceDir, base.Helper(tableName)+".ts")
	if util.FileExists(outInterfaceFile) {
		util.Delete(outInterfaceFile)
	}

	file, err := os.OpenFile(outInterfaceFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}

	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outInterfaceFile))
}

// InterfaceTemplate 返回 TypeScript 接口定义模板字符串。
//
// 配置项 template.VueInterface 指定了自定义模板文件且文件存在时，读取该文件内容并返回；
// 否则返回包内置模板，模板占位符由 GenerateInterface 渲染填充。
func InterfaceTemplate() string {
	if exists := base.Config.Exists("template.VueInterface"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.VueInterface"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
// {{.Table.TableName}} {{.Table.TableComment}}
{{- if .Page}}
import { ReqPage } from "./index";
{{- end}}
export namespace {{CamelStr .Table.TableName}} {
	{{- range .Method}}
	{{- if eq .Type "List"}}
	{{- if .Page}}
	export interface Req{{CamelStr .Table.TableName}}List extends ReqPage {
	{{- else}}
	export interface Req{{CamelStr .Table.TableName}}List {
	{{- end}}
		{{TypeScriptCondition .Condition}}
	}
	export interface Res{{CamelStr .Table.TableName}}Item {
		{{TypeScript .TableColumn}}
	}
	{{- end}}
	{{- end}}
}`
	return outString
}
