package etcdv3

import (
	"context"

	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// traceUnaryClientInterceptor 构造单向调用的客户端链路追踪拦截器，负责启动客户端 Span、
// 透传并回写 metadata、记录方法名与错误状态。
func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	tracer := globalTrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		ctx, span := tracer.Start(ctx, method, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		ctx = metadata.NewOutgoingContext(ctx, md)
		span.SetAttributes(
			semconv.RPCMethodKey.String(method),
		)
		err = invoker(ctx, method, req, reply, cc, opts...)
		span.SetStatus(codes.Ok, "ok")
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		return err
	}

}

// traceStreamClientInterceptor 构造流式调用的客户端链路追踪拦截器，生命周期覆盖整个流的创建过程。
func traceStreamClientInterceptor() grpc.StreamClientInterceptor {

	tracer := globalTrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
	}

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		ctx, span := tracer.Start(ctx, method, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		ctx = metadata.NewOutgoingContext(ctx, md)
		span.SetAttributes(
			semconv.RPCMethodKey.String(method),
		)
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		span.SetStatus(codes.Ok, "ok")
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
		return clientStream, nil
	}
}
