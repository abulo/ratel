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
		"Convert":         base.Convert,
		"SymbolChar":      base.SymbolChar,
		"Char":            base.Char,
		"Helper":          base.Helper,
		"CamelStr":        base.CamelStr,
		"Add":             base.Add,
		"ProtoConvertDao": base.ProtoConvertDao,
		"DaoConvertProto": base.DaoConvertProto,
		"ProtoConvertMap": base.ProtoConvertMap,
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

	"github.com/pkg/errors"
	"github.com/abulo/ratel/v3/server/xgrpc"
	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/util"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	{{ProtoConvertDao .TableColumn "res" "request"}}
	_, err := {{.Pkg}}.{{CamelStr .Table.TableName}}ItemCreate(ctx, res)
	if err != nil {
		return &{{CamelStr .Table.TableName}}ItemCreateResponse{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	return &{{CamelStr .Table.TableName}}ItemCreateResponse{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, err
}

// {{CamelStr .Table.TableName}}ItemUpdate 更新数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,request *{{CamelStr .Table.TableName}}ItemUpdateRequest) (*{{CamelStr .Table.TableName}}ItemUpdateResponse,error){
	{{.Primary.ColumnName}} := request.Get{{CamelStr .Primary.ColumnName}}()
	if {{.Primary.ColumnName}} < 1 {
		return &{{CamelStr .Table.TableName}}ItemUpdateResponse{
			Code: code.ParamInvalid,
			Msg:  code.StatusText(code.ParamInvalid),
		}, errors.New(code.StatusText(code.ParamInvalid))
	}
	var res dao.{{CamelStr .Table.TableName}}
	{{ProtoConvertDao .TableColumn "res" "request"}}
	_, err := {{.Pkg}}.{{CamelStr .Table.TableName}}ItemUpdate(ctx, {{.Primary.ColumnName}}, res)
	if err != nil {
		return &{{CamelStr .Table.TableName}}ItemUpdateResponse{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	return &{{CamelStr .Table.TableName}}ItemUpdateResponse{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, err
}

// {{CamelStr .Table.TableName}}Item 获取数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{CamelStr .Table.TableName}}Item(ctx context.Context,request *{{CamelStr .Table.TableName}}ItemRequest)(*{{CamelStr .Table.TableName}}ItemResponse,error){
	{{.Primary.ColumnName}} := request.Get{{CamelStr .Primary.ColumnName}}()
	if {{.Primary.ColumnName}} < 1 {
		return &{{CamelStr .Table.TableName}}ItemResponse{
			Code: code.ParamInvalid,
			Msg:  code.StatusText(code.ParamInvalid),
		}, errors.New(code.StatusText(code.ParamInvalid))
	}
	res,err := {{.Pkg}}.{{CamelStr .Table.TableName}}Item(ctx ,{{.Primary.ColumnName}})
	if err != nil {
		return &{{CamelStr .Table.TableName}}ItemResponse{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
			Data: &{{CamelStr .Table.TableName}}Object{},
		}, err
	}
	return &{{CamelStr .Table.TableName}}ItemResponse{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data:  &{{CamelStr .Table.TableName}}Object{
			{{DaoConvertProto .TableColumn "res"}}
		},
	}, err
}

// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,request *{{CamelStr .Table.TableName}}ItemDeleteRequest)(*{{CamelStr .Table.TableName}}ItemDeleteResponse,error){
	{{.Primary.ColumnName}} := request.Get{{CamelStr .Primary.ColumnName}}()
	if {{.Primary.ColumnName}} < 1 {
		return &{{CamelStr .Table.TableName}}ItemDeleteResponse{
			Code: code.ParamInvalid,
			Msg:  code.StatusText(code.ParamInvalid),
		}, errors.New(code.StatusText(code.ParamInvalid))
	}

	_, err := {{.Pkg}}.{{CamelStr .Table.TableName}}ItemDelete(ctx, {{.Primary.ColumnName}})
	if err != nil {
		return &{{CamelStr .Table.TableName}}ItemDeleteResponse{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	return &{{CamelStr .Table.TableName}}ItemDeleteResponse{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, err
}

{{- range .Method}}
{{- if eq .Type "List"}}
{{- if .Default}}
// {{CamelStr .Table.TableName}}{{CamelStr .Name}} 列表数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx context.Context,request *{{CamelStr .Table.TableName}}{{CamelStr .Name}}Request)(*{{CamelStr .Table.TableName}}{{CamelStr .Name}}Response,error){
	// 数据库查询条件
	condition := make(map[string]interface{})
	// 当前页面
	pageNumber := request.GetPageNumber()
	// 每页多少数据
	resultPerPage := request.GetResultPerPage()
	if pageNumber < 1 {
		pageNumber = 1
	}
	if resultPerPage < 1 {
		resultPerPage = 10
	}
	// 分页数据
	offset := resultPerPage * (pageNumber - 1)
	condition["pageOffset"] = offset
	condition["pageSize"] = resultPerPage
	// 构造查询条件
	{{ProtoConvertMap .Condition "request"}}
	// 获取数据量
	total, err := {{.Pkg}}.{{CamelStr .Table.TableName}}{{CamelStr .Name}}Total(ctx, condition)
	if err != nil {
		return &{{CamelStr .Table.TableName}}{{CamelStr .Name}}Response{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	// 获取数据集合
	list, err := {{.Pkg}}.{{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx, condition)
	if err != nil {
		return &{{CamelStr .Table.TableName}}{{CamelStr .Name}}Response{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	var res []*{{CamelStr .Table.TableName}}Object
	for _, item := range list {
		res = append(res, &{{CamelStr .Table.TableName}}Object{
			{{DaoConvertProto .TableColumn "item"}}
		})
	}
	return &{{CamelStr .Table.TableName}}{{CamelStr .Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: &{{CamelStr .Table.TableName}}ListObject{
			Total: total,
			List:  res,
		},
	}, nil
}
{{- else}}
// {{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}} 列表数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}(ctx context.Context,request *{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Request)(*{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response,error){
	// 数据库查询条件
	condition := make(map[string]interface{})
	// 当前页面
	pageNumber := request.GetPageNumber()
	// 每页多少数据
	resultPerPage := request.GetResultPerPage()
	if pageNumber < 1 {
		pageNumber = 1
	}
	if resultPerPage < 1 {
		resultPerPage = 10
	}
	// 分页数据
	offset := resultPerPage * (pageNumber - 1)
	condition["pageOffset"] = offset
	condition["pageSize"] = resultPerPage
	// 构造查询条件
	{{ProtoConvertMap .Condition "request"}}
	// 获取数据量
	total, err := {{.Pkg}}.{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Total(ctx, condition)
	if err != nil {
		return &{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	// 获取数据集合
	list, err := {{.Pkg}}.{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}(ctx, condition)
	if err != nil {
		return &{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
		}, err
	}
	var res []*{{CamelStr .Table.TableName}}Object
	for _, item := range list {
		res = append(res, &{{CamelStr .Table.TableName}}Object{
			{{DaoConvertProto .TableColumn "item"}}
		})
	}
	return &{{CamelStr .Table.TableName}}ListBy{{CamelStr .Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: &{{CamelStr .Table.TableName}}ListObject{
			Total: total,
			List:  res,
		},
	}, nil
}
{{- end}}
{{- else}}
// {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}} 单列数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}(ctx context.Context,request *{{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Request)(*{{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response,error){
	// 数据库查询条件
	condition := make(map[string]interface{})
	// 构造查询条件
	{{ProtoConvertMap .Condition "request"}}
	res,err := {{.Pkg}}.{{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}(ctx ,condition)
	if err != nil {
		return &{{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response{
			Code: code.Fail,
			Msg:  code.StatusText(code.Fail),
			Data: &{{CamelStr .Table.TableName}}Object{},
		}, err
	}
	return &{{CamelStr .Table.TableName}}ItemBy{{CamelStr .Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data:  &{{CamelStr .Table.TableName}}Object{
			{{DaoConvertProto .TableColumn "res"}}
		},
	}, err
}
{{- end}}
{{- end}}
`
	return outString
}
