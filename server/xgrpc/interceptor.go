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

// StreamInterceptorChain returns stream interceptors chain.
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

// UnaryInterceptorChain returns interceptors chain.
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

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

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
