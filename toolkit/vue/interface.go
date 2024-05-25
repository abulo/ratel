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

func GenerateInterface(moduleParam base.ModuleParam, fullInterfaceDir, tableName string) {
	// 模板变量
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

	// 文件夹路径
	outInterfaceFile := path.Join(fullInterfaceDir, base.Helper(tableName)+".ts")
	if util.FileExists(outInterfaceFile) {
		util.Delete(outInterfaceFile)
	}

	file, err := os.OpenFile(outInterfaceFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outInterfaceFile))
}

func InterfaceTemplate() string {
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
