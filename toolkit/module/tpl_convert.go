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

// GenerateConvert 使用 ConvertTemplate 渲染并写出 DAO 与 Proto 对象互转的 Go 文件到 fullConvertDir，
// 文件已存在时先删除再生成，随后执行 go fmt 与 goimports 整理格式；模板渲染失败时打印错误并返回。
func GenerateConvert(moduleParam base.ModuleParam, fullConvertDir, tableName string) {
	// PERF-8-20 [P2] 公共模板函数复用 base 包级定义，避免每次调用重复构建同一 FuncMap 字面量。
	funcs := base.CommonTemplateFuncs()
	tpl := template.Must(template.New("convert").Funcs(funcs).Parse(ConvertTemplate()))
	outConvertFile := path.Join(fullConvertDir, tableName+".go")
	if util.FileExists(outConvertFile) {
		util.Delete(outConvertFile)
	}
	file, err := os.OpenFile(outConvertFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
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

// ConvertTemplate 返回 DAO 与 Proto 对象互转代码模板；配置中 template.Convert 指定了自定义模板文件时优先读取该文件内容。
func ConvertTemplate() string {
	if exists := base.Config.Exists("template.Convert"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.Convert"))
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
