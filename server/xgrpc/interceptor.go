package xgrpc

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/hostname"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// StreamInterceptorChain 将多个流式服务端拦截器按传入顺序串联为单个拦截器；
// 未传入拦截器时返回直接调用处理器的拦截器。
func StreamInterceptorChain(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	build := func(c grpc.StreamServerInterceptor, n grpc.StreamHandler, info *grpc.StreamServerInfo) grpc.StreamHandler {
		return func(srv any, stream grpc.ServerStream) error {
			return c(srv, stream, info, n)
		}
	}
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		chain := handler
		for _, interceptor := range slices.Backward(interceptors) {
			chain = build(interceptor, chain, info)
		}
		return chain(srv, stream)
	}
}

// UnaryInterceptorChain 将多个一元服务端拦截器按传入顺序串联为单个拦截器；
// 未传入拦截器时返回直接调用处理器的拦截器。
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	build := func(c grpc.UnaryServerInterceptor, n grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
		return func(ctx context.Context, req any) (any, error) {
			return c(ctx, req, info, n)
		}
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		chain := handler
		for _, interceptor := range slices.Backward(interceptors) {
			chain = build(interceptor, chain, info)
		}
		return chain(ctx, req)
	}
}

// prometheusUnaryServerInterceptor 记录一元 RPC 的耗时直方图与请求计数器，错误码从返回错误中提取。
func prometheusUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx), code.GetMessage())
	return resp, err
}

// prometheusStreamServerInterceptor 记录流式 RPC 的耗时直方图与请求计数器，错误码从返回错误中提取。
func prometheusStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), code.GetMessage())
	return err
}

// NewTraceUnaryServerInterceptor 返回一元 RPC 服务端链路追踪拦截器：创建服务端 Span、
// 填充方法、主机与对端属性，请求结束时记录错误状态并结束 Span。
func NewTraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, err error) {
		var remote string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		operation, mAttrs := globalTrace.ParseFullMethod(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.HostName(hostname.Hostname()),
		}
		attrs = append(attrs, mAttrs...)
		if p, ok := peer.FromContext(ctx); ok {
			remote = p.Addr.String()
		}
		attrs = append(attrs, globalTrace.PeerAttr(remote)...)
		ctx, span := tracer.Start(ctx, operation, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer func() {
			if err != nil {
				span.RecordError(err)
				s, ok := status.FromError(err)
				if ok {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(s.Code())))
				} else {
					span.SetStatus(codes.Error, err.Error())
				}
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return handler(ctx, req)
	}
}

// contextedServerStream 在包装底层 ServerStream 的同时携带替换后的上下文，供追踪拦截器向下传递。
type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context 返回携带追踪信息的上下文，覆盖底层流的原始上下文。
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

// NewTraceStreamServerInterceptor 返回流式 RPC 服务端链路追踪拦截器：创建服务端 Span、
// 填充方法、主机与对端属性，并将携带追踪上下文的流传给下游处理器。
func NewTraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		tracer := globalTrace.NewTracer(trace.SpanKindServer)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.HostName(hostname.Hostname()),
		}
		var remote string
		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		operation, mAttrs := globalTrace.ParseFullMethod(info.FullMethod)
		attrs = append(attrs, mAttrs...)
		if p, ok := peer.FromContext(ss.Context()); ok {
			remote = p.Addr.String()
		}
		attrs = append(attrs, globalTrace.PeerAttr(remote)...)
		attrs = append(attrs,
			semconv.HostNameKey.String(remote),
		)
		ctx, span := tracer.Start(ss.Context(), operation, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer span.End()

		return handler(srv, contextedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}
