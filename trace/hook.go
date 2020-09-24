package trace

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/abulo/ratel/ecode"
	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/hbase"
	"github.com/abulo/ratel/metric"
	"github.com/abulo/ratel/util"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tsuna/gohbase/hrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ginExtractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

func grpcExtractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

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

//GRPCUnaryServerInterceptor 监控程序跟踪
func GRPCUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, grpcExtractAID(ctx))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, grpcExtractAID(ctx), code.GetMessage())
	return resp, err
}

//GRPCStreamServerInterceptor 监控程序跟踪
func GRPCStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, grpcExtractAID(ss.Context()))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, grpcExtractAID(ss.Context()), code.GetMessage())
	return err
}

func spanStartName(ctx *gin.Context) string {
	return ctx.Request.Method + " " + ctx.FullPath()
}

//HTTPTraceServerInterceptor 链路跟踪
func HTTPTraceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := StartSpanFromContext(
			c.Request.Context(),
			spanStartName(c),
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

// RPCTraceUnaryServerInterceptor ...
func RPCTraceUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := StartSpanFromContext(
		ctx,
		info.FullMethod,
		FromIncomingContext(ctx),
		TagComponent("gRPC"),
		TagSpanKind("server.unary"),
	)

	defer span.Finish()

	resp, err := handler(ctx, req)

	if err != nil {
		code := codes.Unknown
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		}
		span.SetTag("code", code)
		ext.Error.Set(span, true)
		span.LogFields(String("event", "error"), String("message", err.Error()))
	}
	return resp, err
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

//RPCTraceStreamServerInterceptor ...
func RPCTraceStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	span, ctx := StartSpanFromContext(
		ss.Context(),
		info.FullMethod,
		FromIncomingContext(ss.Context()),
		TagComponent("gRPC"),
		TagSpanKind("server.stream"),
		CustomTag("isServerStream", info.IsServerStream),
	)
	defer span.Finish()

	return handler(srv, contextedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}

//HBaseTrace ...
func HBaseTrace(component, instance string) hbase.HookFunc {
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		if customName == "" {
			customName = call.Name()
		}
		statement := string(call.Table()) + " " + string(call.Key())
		span, ctx := StartSpanFromContext(
			ctx,
			"Hbase:"+customName,
			TagComponent("Hbase"),
			TagSpanKind("client"),
			CustomTag("statement", statement),
		)
		return func(err error) {
			if err == io.EOF {
				err = nil
			}
			defer span.Finish()
		}
	}
}
