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

func GenerateService(moduleParam base.ModuleParam, fullServiceDir, tableName string) {
	tpl := template.Must(template.New("service").Funcs(template.FuncMap{
		"Convert":    base.Convert,
		"SymbolChar": base.SymbolChar,
		"Char":       base.Char,
		"Helper":     base.Helper,
		"CamelStr":   base.CamelStr,
		"Add":        base.Add,
		"ConvertDao": base.ConvertDao,
	}).Parse(ServiceTemplate()))

	// 文件夹路径
	outServiceFile := path.Join(fullServiceDir, tableName+"_service.go")
	if util.FileExists(outServiceFile) {
		util.Delete(outServiceFile)
	}
	file, err := os.OpenFile(outServiceFile, os.O_CREATE|os.O_WRONLY, 0755)
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
	_ = os.Chdir(fullServiceDir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("代码格式化错误:", color.RedString(err.Error()))
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(fullServiceDir, "*.go"))
	cmdImport.CombinedOutput()
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outServiceFile))

}
func ServiceTemplate() string {
	outString := `
package {{.Pkg}}

import (
	"{{.ModName}}/code"
	"{{.ModName}}/dao"
	"{{.ModName}}/module/{{.Pkg}}"
	"context"

	"github.com/abulo/ratel/v3/server/xgrpc"
	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/util"
)
// {{.Table.TableName}} {{.Table.TableComment}}


// Srv{{CamelStr .Table.TableName}}ServiceServer {{.Table.TableComment}}
type Srv{{CamelStr .Table.TableName}}ServiceServer struct {
	Unimplemented{{CamelStr .Table.TableName}}ServiceServer
	Server *xgrpc.Server
}

// {{CamelStr .Table.TableName}}ItemCreate 创建数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{CamelStr .Table.TableName}}ItemCreate(ctx context.Context,request *{{CamelStr .Table.TableName}}ItemCreateRequest) (*{{CamelStr .Table.TableName}}ItemCreateResponse,error){
	var res dao.{{CamelStr .Table.TableName}}
	{{ConvertDao .TableColumn "res" "request"}}
}

`
	return outString
}
