package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abulo/ratel/v2"
	"github.com/abulo/ratel/v2/ecode"
	"github.com/abulo/ratel/v2/metric"
	"github.com/abulo/ratel/v2/trace"
	"github.com/abulo/ratel/v2/util"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	errSlowCommand = errors.New("grpc unary slow command")
)

// metric统计
func metricUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 收敛err错误，将err过滤后，可以知道err是否为系统错误码
		spbStatus := ecode.ExtractCodes(err)
		// 只记录系统级别错误
		if spbStatus.Code < ecode.EcodeNum {
			// 只记录系统级别的详细错误码
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, cc.Target(), spbStatus.GetMessage())
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, cc.Target())
		} else {
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, cc.Target(), "biz error")
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, cc.Target())
		}
		return err
	}
}

func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("04:05.000"), "Send: "+method+" | "+util.JsonString(req))
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("04:05.000"), "Erro: "+err.Error())
		} else {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("04:05.000"), "Recv: "+util.JsonString(reply))
		}

		return err
	}
}

func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		remoteIP := "unknown"
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			remoteIP = remote.Addr.String()
		}

		span, ctx := trace.StartSpanFromContext(
			ctx,
			trace.SpanRpcClientStartName(method),
			trace.TagSpanKind("client"),
			trace.TagComponent("grpc"),
			trace.CustomTag("peer.ipv4", remoteIP),
		)

		span.LogFields(log.Object("req", req))
		span.LogFields(log.Object("reply", reply))
		defer span.Finish()

		err := invoker(trace.MetadataInjector(ctx, md), method, req, reply, cc, opts...)
		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("response_code", code)
			ext.Error.Set(span, true)

			span.LogFields(trace.String("event", "error"), trace.String("message", err.Error()))
		}
		return err
	}
}

func aidUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		clientAidMD := metadata.Pairs("aid", ratel.AppID())
		if ok {
			md = metadata.Join(md, clientAidMD)
		} else {
			md = clientAidMD
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// timeoutUnaryClientInterceptor gRPC客户端超时拦截器
func timeoutUnaryClientInterceptor(_logger *logrus.Logger, timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		// 若无自定义超时设置，默认设置超时
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
			_logger.WithFields(logrus.Fields{
				"err":    errSlowCommand,
				"method": method,
				"target": cc.Target(),
				"du":     du,
				"ip":     remoteIP,
			}).Info("slow")
		}
		return err
	}
}

// loggerUnaryClientInterceptor gRPC客户端日志中间件
func loggerUnaryClientInterceptor(_logger *logrus.Logger, name string, accessInterceptorLevel string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		spbStatus := ecode.ExtractCodes(err)
		if err != nil {
			// 只记录系统级别错误
			if spbStatus.Code < ecode.EcodeNum {
				// 只记录系统级别错误
				_logger.WithFields(logrus.Fields{
					"code":    spbStatus.Code,
					"message": spbStatus.Message,
					"name":    name,
					"method":  method,
					"beg":     time.Since(beg),
					"req":     req,
					"reply":   reply,
				}).Info("access")

			} else {
				// 业务报错只做warning
				_logger.WithFields(logrus.Fields{
					"code":    spbStatus.Code,
					"message": spbStatus.Message,
					"name":    name,
					"method":  method,
					"beg":     time.Since(beg),
					"req":     req,
					"reply":   reply,
				}).Info("access")
			}
			return err
		} else {
			if accessInterceptorLevel == "info" {
				_logger.WithFields(logrus.Fields{
					"code":   spbStatus.Code,
					"name":   name,
					"method": method,
					"beg":    time.Since(beg),
					"req":    req,
					"reply":  reply,
				}).Info("info")
			}
		}
		return nil
	}
}
