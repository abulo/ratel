package xgrpc

import (
	"context"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// StreamInterceptorChain returns stream interceptors chain.
func StreamInterceptorChain(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	build := func(c grpc.StreamServerInterceptor, n grpc.StreamHandler, info *grpc.StreamServerInfo) grpc.StreamHandler {
		return func(srv any, stream grpc.ServerStream) error {
			return c(srv, stream, info, n)
		}
	}
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}
		return chain(srv, stream)
	}
}

// UnaryInterceptorChain returns interceptors chain.
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	build := func(c grpc.UnaryServerInterceptor, n grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
		return func(ctx context.Context, req any) (any, error) {
			return c(ctx, req, info, n)
		}
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}
		return chain(ctx, req)
	}
}

func prometheusUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx), code.GetMessage())
	return resp, err
}

func prometheusStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), code.GetMessage())
	return err
}

func traceUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	span, ctx := trace.StartSpanFromContext(
		ctx,
		info.FullMethod,
		trace.FromIncomingContext(ctx),
		trace.TagComponent("gRPC"),
		trace.TagSpanKind("server.unary"),
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
		span.LogFields(trace.String("event", "error"), trace.String("message", err.Error()))
	}
	return resp, err
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	span, ctx := trace.StartSpanFromContext(
		ss.Context(),
		info.FullMethod,
		trace.FromIncomingContext(ss.Context()),
		trace.TagComponent("gRPC"),
		trace.TagSpanKind("server.stream"),
		trace.CustomTag("isServerStream", info.IsServerStream),
	)
	defer span.Finish()

	return handler(srv, contextedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}
