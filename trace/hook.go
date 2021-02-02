package trace

import (
	"github.com/abulo/ratel/gin"
)

func ginExtractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

func spanStartName(ctx *gin.Context) string {
	return ctx.Request.Method + " " + ctx.FullPath()
}
