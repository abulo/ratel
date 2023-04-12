package api

// HertzTemplate 模板
func HertzTemplate() string {
	outString := `
package {{.Pkg}}

import (
	"context"

	"{{.ModName}}/code"
	"{{.ModName}}/dao"
	"{{.ModName}}/initial"
	"{{.ModName}}/service/{{.PkgPath}}"

	globalLogger "github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
	"github.com/sirupsen/logrus"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// {{.Table.TableName}} {{.Table.TableComment}}

// {{CamelStr .Table.TableName}}Dao 数据转换
func {{CamelStr .Table.TableName}}Dao(item *{{.Pkg}}.{{CamelStr .Table.TableName}}Object) dao.{{CamelStr .Table.TableName}} {
	daoItem := dao.{{CamelStr .Table.TableName}}{}
	{{ModuleProtoConvertDao .TableColumn "daoItem" "item"}}
	return daoItem
}

// {{CamelStr .Table.TableName}}Proto 数据绑定
func {{CamelStr .Table.TableName}}Proto(item dao.{{CamelStr .Table.TableName}}) *{{.Pkg}}.{{CamelStr .Table.TableName}}Object {
	return  &{{.Pkg}}.{{CamelStr .Table.TableName}}Object{
		{{ModuleDaoConvertProto .TableColumn "item"}}
	}
}

// {{CamelStr .Table.TableName}}ItemCreate 创建数据
func {{CamelStr .Table.TableName}}ItemCreate(ctx context.Context,newCtx *app.RequestContext) {

	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemCreate")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  err.Error(),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	request := &{{.Pkg}}.{{CamelStr .Table.TableName}}ItemCreateRequest{}
	// 数据绑定
	var reqInfo dao.{{CamelStr .Table.TableName}}
	if err := newCtx.BindAndValidate(&reqInfo); err != nil {
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ParamInvalid,
			"msg":  err.Error(),
		})
		return
	}
	request.Data = {{CamelStr .Table.TableName}}Proto(reqInfo)
	// 执行服务
	res, err := client.{{CamelStr .Table.TableName}}ItemCreate(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemCreate")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RemoteServiceError,
			"msg":  err.Error(),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}

// {{CamelStr .Table.TableName}}ItemUpdate 更新数据
func {{CamelStr .Table.TableName}}ItemUpdate(ctx context.Context,newCtx *app.RequestContext) {
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemUpdate")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  err.Error(),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{CamelStr .Table.TableName}}ItemUpdateRequest{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}

	// 数据绑定
	var reqInfo dao.{{CamelStr .Table.TableName}}
	if err := newCtx.BindAndValidate(&reqInfo); err != nil {
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ParamInvalid,
			"msg":  err.Error(),
		})
		return
	}
	request.Data = {{CamelStr .Table.TableName}}Proto(reqInfo)

	// 执行服务
	res, err := client.{{CamelStr .Table.TableName}}ItemUpdate(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemUpdate")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RemoteServiceError,
			"msg":  err.Error(),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}

// {{CamelStr .Table.TableName}}Item 获取数据
func {{CamelStr .Table.TableName}}Item(ctx context.Context,newCtx *app.RequestContext){
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}Item")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  err.Error(),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{CamelStr .Table.TableName}}ItemRequest{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 执行服务
	res, err := client.{{CamelStr .Table.TableName}}Item(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}Item")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RemoteServiceError,
			"msg":  err.Error(),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
		"data": {{CamelStr .Table.TableName}}Dao(res.GetData()),
	})
}


// {{CamelStr .Table.TableName}}ItemDelete 删除数据
func {{CamelStr .Table.TableName}}ItemDelete(ctx context.Context,newCtx *app.RequestContext){
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemDelete")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  err.Error(),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{CamelStr .Table.TableName}}ItemDeleteRequest{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 执行服务
	res, err := client.{{CamelStr .Table.TableName}}ItemDelete(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}ItemDelete")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RemoteServiceError,
			"msg":  err.Error(),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}

{{- range .Method}}
{{- if eq .Type "List"}}
{{- if .Default}}
// {{CamelStr .Table.TableName}}{{CamelStr .Name}} 列表数据
func {{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx context.Context,newCtx *app.RequestContext){
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}{{CamelStr .Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  err.Error(),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	request := &{{.Pkg}}.{{CamelStr .Table.TableName}}{{CamelStr .Name}}Request{}
	// 构造查询条件
	{{ApiToProto .Condition "request" "newCtx.Query"}}
	request.PageNum = cast.ToInt64(newCtx.Query("pageNum"))
	request.PageSize = cast.ToInt64(newCtx.Query("pageSize"))
	// 执行服务
	res, err := client.{{CamelStr .Table.TableName}}{{CamelStr .Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{CamelStr .Table.TableName}}{{CamelStr .Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RemoteServiceError,
			"msg":  err.Error(),
		})
		return
	}

	var total int64
	var list []dao.{{CamelStr .Table.TableName}}

	if res.GetCode() == code.Success {
		total = res.GetData().GetTotal()
		rpcList := res.GetData().GetList()
		for _, item := range rpcList {
			list = append(list, {{CamelStr .Table.TableName}}Dao(item))
		}
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
		"data": utils.H{
			"total": total,
			"list":  list,
			"pageNum": request.PageNum,
			"pageSize": request.PageSize,
		},
	})
}
{{- end}}
{{- end}}
{{- end}}
`
	return outString
}
