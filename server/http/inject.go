package http

import (
	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/trace"
)

func spanStartName(ctx *gin.Context) string {
	return ctx.Request.Method + " " + ctx.FullPath()
}

//TraceServerInterceptor 链路跟踪
func TraceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		span, ctx := trace.StartSpanFromContext(
			c.Request.Context(),
			spanStartName(c),
			trace.TagComponent("http"),
			trace.TagSpanKind("server"),
			trace.HeaderExtractor(c.Request.Header),
			trace.CustomTag("http.url", c.Request.URL.Path),
			trace.CustomTag("http.method", c.Request.Method),
			trace.CustomTag("peer.ipv4", c.ClientIP()),
		)

		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
		c.Next()

	}
}
