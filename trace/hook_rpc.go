package trace

import (
	"context"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
