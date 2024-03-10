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
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/proto"
)

// {{.Table.TableName}} {{.Table.TableComment}}

{{- range .Method}}
{{- if eq .Type "Create"}}
// {{.Name}} 创建数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext) {
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	request := &{{.Pkg}}.{{.Name}}Request{}
	// 数据绑定
	var reqInfo dao.{{CamelStr .Table.TableName}}
	if err := newCtx.BindAndValidate(&reqInfo); err != nil {
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ParamInvalid,
			"msg":  code.StatusText(code.ParamInvalid),
		})
		return
	}
	request.Data = {{.Pkg}}.{{CamelStr .Table.TableName}}Proto(reqInfo)
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}
{{- else if eq .Type "Update"}}
// {{.Name}} 更新数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext) {
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{.Name}}Request{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 数据绑定
	var reqInfo dao.{{CamelStr .Table.TableName}}
	if err := newCtx.BindAndValidate(&reqInfo); err != nil {
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ParamInvalid,
			"msg":  code.StatusText(code.ParamInvalid),
		})
		return
	}
	request.Data = {{.Pkg}}.{{CamelStr .Table.TableName}}Proto(reqInfo)
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}
{{- else if eq .Type "Show"}}
// {{.Name}} 查询单条数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext){
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{.Name}}Request{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
		"data": {{.Pkg}}.{{CamelStr .Table.TableName}}Dao(res.GetData()),
	})
}
{{- else if eq .Type "Delete"}}
// {{.Name}} 删除数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext){
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{.Name}}Request{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}
{{- else if eq .Type "Recover"}}
// {{.Name}} 恢复数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext){
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	{{Helper .Primary.AlisaColumnName }} := cast.ToInt64(newCtx.Param("{{Helper .Primary.AlisaColumnName }}"))
	request := &{{.Pkg}}.{{.Name}}Request{}
	request.{{CamelStr .Primary.AlisaColumnName }} = {{Helper .Primary.AlisaColumnName }}
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
	})
}
{{- else if eq .Type "Item"}}
// {{.Name}} 查询单条数据
func {{.Name}}(ctx context.Context, newCtx *app.RequestContext) {
	//判断这个服务能不能链接
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	request := &{{.Pkg}}.{{.Name}}Request{}
	// 构造查询条件
	{{ApiToProto .Condition "request" "newCtx.GetQuery" false}}
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
		"data": {{.Pkg}}.{{CamelStr .Table.TableName}}Dao(res.GetData()),
	})
}
{{- else if eq .Type "List"}}
// {{.Name}} 列表数据
func {{.Name}}(ctx context.Context,newCtx *app.RequestContext){
	grpcClient, err := initial.Core.Client.LoadGrpc("grpc").Singleton()
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Grpc:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.RPCError,
			"msg":  code.StatusText(code.RPCError),
		})
		return
	}
	//链接服务
	client := {{.Pkg}}.New{{CamelStr .Table.TableName}}ServiceClient(grpcClient)
	// 构造查询条件
	request := &{{.Pkg}}.{{.Name}}Request{}
	{{- if .Page}}
	requestTotal := &{{.Pkg}}.{{.Name}}TotalRequest{}
	{{ApiToProto .Condition "request" "newCtx.GetQuery" true}}
	{{- else}}
	{{ApiToProto .Condition "request" "newCtx.GetQuery" false}}
	{{- end}}

	{{- if .Page}}
	// 执行服务,获取数据量
	resTotal, err := client.{{.Name}}Total(ctx, requestTotal)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	var total int64
	request.PageNum = proto.Int64(cast.ToInt64(newCtx.Query("pageNum")))
	request.PageSize = proto.Int64(cast.ToInt64(newCtx.Query("pageSize")))
	if resTotal.GetCode() == code.Success {
		total = resTotal.GetData()
	}
	{{- end}}
	// 执行服务
	res, err := client.{{.Name}}(ctx, request)
	if err != nil {
		globalLogger.Logger.WithFields(logrus.Fields{
			"req": request,
			"err": err,
		}).Error("GrpcCall:{{.Table.TableComment}}:{{.Table.TableName}}:{{.Name}}")
		fromError := status.Convert(err)
		newCtx.JSON(consts.StatusOK, utils.H{
			"code": code.ConvertToHttp(fromError.Code()),
			"msg":  code.StatusText(code.ConvertToHttp(fromError.Code())),
		})
		return
	}
	var list []*dao.{{CamelStr .Table.TableName}}
	if res.GetCode() == code.Success {
		rpcList := res.GetData()
		for _, item := range rpcList {
			list = append(list, {{.Pkg}}.{{CamelStr .Table.TableName}}Dao(item))
		}
	}
	newCtx.JSON(consts.StatusOK, utils.H{
		"code": res.GetCode(),
		"msg":  res.GetMsg(),
		{{- if .Page}}
		"data": utils.H{
			"total": total,
			"list":  list,
			"pageNum": request.PageNum,
			"pageSize": request.PageSize,
		},
		{{- else}}
		"data": list,
		{{- end}}
	})
}
{{- end}}
{{- end}}
`
	return outString
}
