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

func GenerateService(moduleParam base.ModuleParam, fullServiceDir, tableName string) {
	tpl := template.Must(template.New("service").Funcs(template.FuncMap{
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
	if exists := base.Config.Exists("template.Service"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.Service"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
	outString := `
package {{.Pkg}}

import (
	"{{.ModName}}/code"
	"{{.ModName}}/module/{{.PkgPath}}"
	"context"

	globalLogger "github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/server/xgrpc"
	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	{{- if .Page}}
	"{{.ModName}}/service/pagination"
	{{- end}}
)
// {{.Table.TableName}} {{.Table.TableComment}}


// Srv{{CamelStr .Table.TableName}}ServiceServer {{.Table.TableComment}}
type Srv{{CamelStr .Table.TableName}}ServiceServer struct {
	Unimplemented{{CamelStr .Table.TableName}}ServiceServer
	Server *xgrpc.Server
}

{{- range .Method}}
{{- if eq .Type "Create"}}
// {{.Name}} 创建数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{.Name}}(ctx context.Context,request *{{.Name}}Request) (*{{.Name}}Response,error){
	req := {{CamelStr .Table.TableName}}Dao(request.GetData())
	data, err := {{.Pkg}}.{{.Name}}(ctx, *req)
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": req,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: data,
	}, nil
}
{{- else if eq .Type "Update"}}
// {{.Name}} 更新数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{.Name}}(ctx context.Context,request *{{.Name}}Request) (*{{.Name}}Response,error){
	{{Helper .Primary.AlisaColumnName}} := request.Get{{CamelStr .Primary.AlisaColumnName}}()
	if {{Helper .Primary.AlisaColumnName}} < 1 {
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.ParamInvalid), code.StatusText(code.ParamInvalid))
	}
	req := {{CamelStr .Table.TableName}}Dao(request.GetData())
	_, err := {{.Pkg}}.{{.Name}}(ctx, {{Helper .Primary.AlisaColumnName}}, *req)
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": req,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, nil
}
{{- else if eq .Type "Drop"}}
// {{.Name}} 清理数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	{{Helper .Primary.AlisaColumnName}} := request.Get{{CamelStr .Primary.AlisaColumnName}}()
	if {{Helper .Primary.AlisaColumnName}} < 1 {
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.ParamInvalid), code.StatusText(code.ParamInvalid))
	}
	_, err := {{.Pkg}}.{{.Name}}(ctx, {{Helper .Primary.AlisaColumnName}})
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": {{Helper .Primary.AlisaColumnName}},
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, nil
}
{{- else if eq .Type "Delete"}}
// {{.Name}} 删除数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	{{Helper .Primary.AlisaColumnName}} := request.Get{{CamelStr .Primary.AlisaColumnName}}()
	if {{Helper .Primary.AlisaColumnName}} < 1 {
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.ParamInvalid), code.StatusText(code.ParamInvalid))
	}
	_, err := {{.Pkg}}.{{.Name}}(ctx, {{Helper .Primary.AlisaColumnName}})
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": {{Helper .Primary.AlisaColumnName}},
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, nil
}
{{- else if eq .Type "Recover"}}
// {{.Name}} 恢复数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	{{Helper .Primary.AlisaColumnName}} := request.Get{{CamelStr .Primary.AlisaColumnName}}()
	if {{Helper .Primary.AlisaColumnName}} < 1 {
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.ParamInvalid), code.StatusText(code.ParamInvalid))
	}
	_, err := {{.Pkg}}.{{.Name}}(ctx, {{Helper .Primary.AlisaColumnName}})
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": {{Helper .Primary.AlisaColumnName}},
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
	}, nil
}
{{- else if eq .Type "Show"}}
// {{.Name}} 查询单条数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	{{Helper .Primary.AlisaColumnName}} := request.Get{{CamelStr .Primary.AlisaColumnName}}()
	if {{Helper .Primary.AlisaColumnName}} < 1 {
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.ParamInvalid), code.StatusText(code.ParamInvalid))
	}
	res,err := {{.Pkg}}.{{.Name}}(ctx ,{{Helper .Primary.AlisaColumnName}})
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": {{Helper .Primary.AlisaColumnName}},
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{}, status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: {{CamelStr .Table.TableName}}Proto(res),
	}, nil
}
{{- else if eq .Type "Item"}}
// {{.Name}} 查询单条数据
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer) {{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	// 数据库查询条件
	condition := make(map[string]any)
	// 构造查询条件
	{{ModuleProtoConvertMap .Condition "request"}}
	if util.Empty(condition) {
		err := errors.New("condition is empty")
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": condition,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	res,err := {{.Pkg}}.{{.Name}}(ctx ,condition)
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": condition,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data:  {{CamelStr .Table.TableName}}Proto(res),
	}, nil
}
{{- else if eq .Type "List"}}
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{.Name}}(ctx context.Context,request *{{.Name}}Request)(*{{.Name}}Response,error){
	// 数据库查询条件
	condition := make(map[string]any)
	// 构造查询条件
	{{ModuleProtoConvertMap .Condition "request"}}
	{{- if .Page}}
	paginationRequest := request.GetPagination()
	if paginationRequest != nil {
		// 当前页面
		pageNum := paginationRequest.GetPageNum()
		// 每页多少数据
		pageSize := paginationRequest.GetPageSize()
		if pageNum < 1 {
			pageNum = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}
		// 分页数据
		offset := pageSize * (pageNum - 1)
		pagination := &sql.Pagination{
			Offset: &offset,
			Limit:  &pageSize,
		}
		condition["pagination"] = pagination
	}
	{{- end}}
	// 获取数据集合
	list, err := {{.Pkg}}.{{.Name}}(ctx, condition)
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": condition,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		return &{{.Name}}Response{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	var res []*{{CamelStr .Table.TableName}}Object
	for _, item := range list {
		res = append(res, {{CamelStr .Table.TableName}}Proto(item))
	}
	return &{{.Name}}Response{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: res,
	}, nil
}
{{- if .Page}}
// {{.Name}}Total 获取总数
func (srv Srv{{CamelStr .Table.TableName}}ServiceServer){{.Name}}Total(ctx context.Context,request *{{.Name}}TotalRequest)(*{{CamelStr .Table.TableName}}TotalResponse,error){
	// 数据库查询条件
	condition := make(map[string]any)
	// 构造查询条件
	{{ModuleProtoConvertMap .Condition "request"}}
	// 获取数据集合
	total, err := {{.Pkg}}.{{.Name}}Total(ctx, condition)
	if sql.ResultAccept(err) != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": condition,
			"err": err,
		}).Error("Sql:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}Total")
		return &{{CamelStr .Table.TableName}}TotalResponse{},  status.Error(code.ConvertToGrpc(code.SqlError), err.Error())
	}
	return &{{CamelStr .Table.TableName}}TotalResponse{
		Code: code.Success,
		Msg:  code.StatusText(code.Success),
		Data: total,
	}, nil
}
{{- end}}
{{- end}}
{{- end}}
`
	return outString
}
