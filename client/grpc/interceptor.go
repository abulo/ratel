package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	errSlowCommand = errors.New("grpc unary slow command")
)

// PERF-5-4 [P2]：复核结论——该拦截器仅在 EnableDebug 显式开启时装配，默认不进入运行热路径；
// 保持现有调试输出语义，不做级别门控或响应体裁剪，避免改变本地排查行为。
//
// debugUnaryClientInterceptor 构造调试用单向拦截器，在控制台打印目标地址、远端地址、时间戳、
// 错误或响应内容，供本地排查使用。
func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("2006-01-02 15:04:05.000000"), "Error: "+err.Error())
		} else {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("2006-01-02 15:04:05.000000"), "Recv: "+util.JSONString(reply))
		}
		return err
	}
}

// aidUnaryClientInterceptor 构造应用 ID 透传拦截器，将当前应用标识写入出站 metadata 的 aid 键，
// 供服务端识别调用方。
func aidUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		// PERF-5-3 [P2]：存在既有 metadata 时原地 Append aid，避免 metadata.Join 合并两个切片的分配；
		// 无既有 metadata 时直接构造只含 aid 的新 MD，语义与原 Pairs+Join 保持一致。（预期收益需实测验证）
		if ok {
			md.Append("aid", env.AppID())
		} else {
			md = metadata.Pairs("aid", env.AppID())
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// timeoutUnaryClientInterceptor 构造单向调用超时与慢调用统计拦截器：上下文未携带截止时间时按
// timeout 补加超时；调用耗时超过 slowThreshold 时按慢命令记录错误日志。
func timeoutUnaryClientInterceptor(timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		// 上下文未显式设置截止时间时，使用配置的调用超时兜底。
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		du := time.Since(now)
		remoteIP := "unknown"
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			remoteIP = remote.Addr.String()
		}
		if slowThreshold > time.Duration(0) && du > slowThreshold {
			logger.Logger.WithFields(logrus.Fields{
				"err":      errSlowCommand,
				"method":   method,
				"target":   cc.Target(),
				"ip":       remoteIP,
				"duration": du,
			}).Error("slow")
		}
		return err
	}
}

// traceUnaryClientInterceptor 构造单向调用链路追踪拦截器，启动客户端 Span、透传并回写 metadata、
// 记录方法名与错误状态。
func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	tracer := globalTrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemGRPC,
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		// PERF-5-3 [P2]：FromOutgoingContext 返回的已是当次新建的 MD（值切片已拷贝），
		// 去掉冗余的 md.Copy()；carrier.Set 只改写这份当次 MD，不会回写 ctx 中原有的 MD。（预期收益需实测验证）
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

// loggerUnaryClientInterceptor 构造访问日志拦截器：系统级错误记 Error，业务错误记 Warn；
// 调用成功且级别为 info 时记录访问日志，日志携带方法名、耗时、请求与响应内容。
func loggerUnaryClientInterceptor(name string, accessInterceptorLevel string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		spbStatus := ecode.ExtractCodes(err)
		// PERF-5-1 [P1]：Error/Warn/Info 三个分支重复构造同一组日志字段，抽为局部闭包复用，
		// 仅在实际命中分支时才构造字段；req/reply 仍按原逻辑 JSON 序列化输出，日志内容、级别与
		// 触发条件保持不变（按语义红线不做级别门控跳过序列化、不做 payload 截断）。JSONString 内部
		// []byte→string 二次拷贝已由模块 7 在 util 层修复。
		accessFields := func() logrus.Fields {
			return logrus.Fields{
				"code":   spbStatus.Code,
				"msg":    spbStatus.Message,
				"name":   name,
				"method": method,
				"time":   time.Since(beg),
				"req":    json.RawMessage(util.JSONString(req)),
				"reply":  json.RawMessage(util.JSONString(reply)),
			}
		}
		if err != nil {
			// 系统级别错误使用 Error 级别记录。
			if spbStatus.Code < ecode.EcodeNum {
				logger.Logger.WithFields(accessFields()).Error("access_unary")
			} else {
				// 业务错误降级为 Warn，避免淹没系统故障日志。
				logger.Logger.WithFields(accessFields()).Warn("access")
			}
			return err
		}
		if accessInterceptorLevel == "info" {
			logger.Logger.WithFields(accessFields()).Info("access_unary")
		}
		return nil
	}
}

// metricUnaryClientInterceptor 构造指标采集拦截器，按方法与目标统计调用次数与耗时；
// 系统错误与业务错误分别标注错误信息或 biz error。
func metricUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 解析错误码以区分系统错误与业务错误。
		spbStatus := ecode.ExtractCodes(err)
		// PERF-5-2 [P2]：同一请求目标地址恒定，cc.Target() 只取一次供 Inc/Observe 复用，
		// 消除四个标签参数中的重复目标字符串构造。
		target := cc.Target()
		if spbStatus.Code < ecode.EcodeNum {
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, target, spbStatus.GetMessage())
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, target)
		} else {
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, target, "biz error")
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, target)
		}
		return err
	}
}
