package trace

import (
	"net/http"
	"time"

	"github.com/abulo/ratel/v2/gin"
	"github.com/abulo/ratel/v2/metric"
	"github.com/abulo/ratel/v2/util"
)

//HTTPMetricServerInterceptor 监控程序跟踪
func HTTPMetricServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := util.Now()
		peer := c.ClientIP()
		if aid := ginExtractAID(c); aid != "" {
			peer += "?aid=" + aid
		}
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, peer)
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, peer, http.StatusText(c.Writer.Status()))
		c.Next()
		return
	}
}

//HTTPTraceServerInterceptor 链路跟踪
func HTTPTraceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := StartSpanFromContext(
			c.Request.Context(),
			SpanHttpStartName(c),
			TagComponent("http"),
			TagSpanKind("server"),
			HeaderExtractor(c.Request.Header),
			CustomTag("http.url", c.Request.URL.Path),
			CustomTag("http.method", c.Request.Method),
			CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
		c.Next()
	}
}
