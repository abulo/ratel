package trace

import (
	"github.com/abulo/ratel/v2/gin"
)

func ginExtractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

func SpanHttpStartName(ctx *gin.Context) string {
	return "http " + ctx.Request.Method + " " + ctx.FullPath()
}

func SpanRpcClientStartName(method string) string {
	return "grpcClient " + method
}

func SpanRpcServerStartName(method string) string {
	return "grpcServer " + method
}
