package xgrpc

import (
	"context"
	"strings"
	"sync"
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
	// PERF-6-11 [P1] 装配期一次性捕获拦截器列表与固定嵌套顺序；运行期不再经 build 工厂逐层包装。
	// info 属每 RPC 变化参数，仍由最外层签名透传；拦截器嵌套方向与参数传递语义与原实现保持一致：
	// interceptors[0] 为最外层（先进入），最后一个拦截器最贴近终端 handler。
	n := len(interceptors)
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		chain := handler
		for i := n - 1; i >= 0; i-- {
			current := interceptors[i]
			downstream := chain
			chain = func(srv any, stream grpc.ServerStream) error {
				return current(srv, stream, info, downstream)
			}
		}
		return chain(srv, stream)
	}
}

// UnaryInterceptorChain 将多个一元服务端拦截器按传入顺序串联为单个拦截器；
// 未传入拦截器时返回直接调用处理器的拦截器。
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// PERF-6-11 [P1] 装配期一次性捕获拦截器列表与固定嵌套顺序；运行期不再经 build 工厂逐层包装。
	// info 属每 RPC 变化参数，仍由最外层签名透传；拦截器嵌套方向与参数传递语义与原实现保持一致。
	n := len(interceptors)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		chain := handler
		for i := n - 1; i >= 0; i-- {
			current := interceptors[i]
			downstream := chain
			chain = func(ctx context.Context, req any) (any, error) {
				return current(ctx, req, info, downstream)
			}
		}
		return chain(ctx, req)
	}
}

// prometheusUnaryServerInterceptor 记录一元 RPC 的耗时直方图与请求计数器，错误码从返回错误中提取。
func prometheusUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	// PERF-6-12 [P1] AID 每条 RPC 只从 metadata 提取一次，直方图与计数器复用同一结果。
	aid := extractAID(ctx)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, aid)
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, aid, code.GetMessage())
	return resp, err
}

// prometheusStreamServerInterceptor 记录流式 RPC 的耗时直方图与请求计数器，错误码从返回错误中提取。
func prometheusStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	// PERF-6-12 [P1] AID 每条 RPC 只从流上下文 metadata 提取一次，直方图与计数器复用同一结果。
	aid := extractAID(ss.Context())
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, aid)
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, aid, code.GetMessage())
	return err
}

// NewTraceUnaryServerInterceptor 返回一元 RPC 服务端链路追踪拦截器：创建服务端 Span、
// 填充方法、主机与对端属性，请求结束时记录错误状态并结束 Span。
func NewTraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// PERF-6-13 [P1] tracer 与主机名在装配期获取一次并闭包复用；恒定属性在工厂层构造一次。
	// （hostname 系统调用与重复分配的实际收益量级需实测验证）
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	staticAttrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, err error) {
		var remote string
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		// PERF-6-13 [P1] 服务端只读取 metadata 做上下文提取，无需 md.Copy() 深拷贝；
		// 方法名与方法属性按 FullMethod 缓存复用，避免每条 RPC 重复拆分字符串。
		operation, mAttrs := parseFullMethodCached(info.FullMethod)
		attrs := make([]attribute.KeyValue, 0, len(staticAttrs)+len(mAttrs)+2)
		attrs = append(attrs, staticAttrs...)
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
	// PERF-6-13 [P1] tracer 上提到工厂层创建一次并闭包复用；主机名与恒定属性装配期获取/构造。
	// （hostname 系统调用与重复分配的实际收益量级需实测验证）
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	staticAttrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
		semconv.HostName(hostname.Hostname()),
	}
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var remote string
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.MD{}
		}
		// PERF-6-13 [P1] 服务端只读 metadata，去掉每 RPC 的 md.Copy() 深拷贝；方法属性按 FullMethod 缓存。
		operation, mAttrs := parseFullMethodCached(info.FullMethod)
		attrs := make([]attribute.KeyValue, 0, len(staticAttrs)+len(mAttrs)+3)
		attrs = append(attrs, staticAttrs...)
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

// methodAttrEntry 缓存按 FullMethod 解析出的操作名与方法属性；FullMethod 在服务注册后恒定不变。
type methodAttrEntry struct {
	operation string
	attrs     []attribute.KeyValue
}

// fullMethodCache 以 FullMethod 为键缓存 ParseFullMethod 结果，避免每条 RPC 重复拆分字符串并分配属性切片。
var fullMethodCache sync.Map // map[string]*methodAttrEntry

func parseFullMethodCached(fullMethod string) (string, []attribute.KeyValue) {
	if v, ok := fullMethodCache.Load(fullMethod); ok {
		e := v.(*methodAttrEntry)
		return e.operation, e.attrs
	}
	operation, attrs := globalTrace.ParseFullMethod(fullMethod)
	e := &methodAttrEntry{operation: operation, attrs: attrs}
	fullMethodCache.Store(fullMethod, e)
	return operation, attrs
}
