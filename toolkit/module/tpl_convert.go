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

func GenerateConvert(moduleParam base.ModuleParam, fullConvertDir, tableName string) {
	// 模板变量
	tpl := template.Must(template.New("convert").Funcs(template.FuncMap{
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
	}).Parse(ConvertTemplate()))
	// 文件夹路径
	outConvertFile := path.Join(fullConvertDir, tableName+".go")
	if util.FileExists(outConvertFile) {
		util.Delete(outConvertFile)
	}
	file, err := os.OpenFile(outConvertFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	_ = os.Chdir(fullConvertDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullConvertDir, "*.go"))
	cmdImport.CombinedOutput()

	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outConvertFile))
}

func ConvertTemplate() string {
	if exists := base.Config.Exists("template.Convert"); exists {
		filePath := path.Join(base.Path, base.Config.String("template.Convert"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
package {{.Pkg}}

import (
	"{{.ModName}}/dao"

	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)
// {{.Table.TableName}} {{.Table.TableComment}}

// {{CamelStr .Table.TableName}}Dao 数据转换
func {{CamelStr .Table.TableName}}Dao(item *{{CamelStr .Table.TableName}}Object) *dao.{{CamelStr .Table.TableName}} {
	daoItem := &dao.{{CamelStr .Table.TableName}}{}
	{{ModuleProtoConvertDao .TableColumn "daoItem" "item"}}
	return daoItem
}

// {{CamelStr .Table.TableName}}Proto 数据绑定
func {{CamelStr .Table.TableName}}Proto(item dao.{{CamelStr .Table.TableName}}) *{{CamelStr .Table.TableName}}Object {
	res := &{{CamelStr .Table.TableName}}Object{}
	{{ModuleDaoConvertProto .TableColumn "res" "item"}}
	return res
}

`
	return outString
}
