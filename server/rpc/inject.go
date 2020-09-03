package rpc

import (
	"context"

	"github.com/abulo/ratel/trace"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TraceUnaryServerInterceptor ...
func TraceUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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

//TraceStreamServerInterceptor ...
func TraceStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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
