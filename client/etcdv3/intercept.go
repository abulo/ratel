package etcdv3

import (
	"context"

	"github.com/abulo/ratel/v3/core/trace"
	"google.golang.org/grpc"
)

func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		span, ctx := trace.StartSpanFromContext(
			ctx,
			method,
			trace.FromIncomingContext(ctx),
			trace.TagComponent("gRPC"),
			trace.TagSpanKind("server.unary"),
		)
		defer span.Finish()

		err = invoker(ctx, method, req, reply, cc, opts...)
		// span.SetTag("code", ecode.OK.Code)

		// if err != nil {
		// 	span.SetTag("code", ecode.ExtractCodes())
		// }
		return err
	}

}

func traceStreamClientInterceptor() grpc.StreamClientInterceptor {

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		span, ctx := trace.StartSpanFromContext(
			ctx,
			method,
			trace.FromIncomingContext(ctx),
			trace.TagComponent("gRPC"),
			trace.TagSpanKind("server.unary"),
		)
		defer span.Finish()

		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		// span.SetTag("code", codes.Ok)

		// if err != nil {
		// 	span.SetTag("code", codes.Error)
		// }

		return clientStream, err
	}
}
