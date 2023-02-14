package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abulo/ratel/core/ecode"
	"github.com/abulo/ratel/core/env"
	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/core/metric"
	"github.com/abulo/ratel/core/trace"
	"github.com/abulo/ratel/util"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
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

func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("04:05.000"), "Erro: "+err.Error())
		} else {
			fmt.Printf("%-50s[%s] => %s\n", prefix, time.Now().Format("04:05.000"), "Recv: "+util.JSONString(reply))
		}
		return err
	}
}

func aidUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		clientAidMD := metadata.Pairs("aid", env.AppID())
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
func timeoutUnaryClientInterceptor(timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
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

func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		span, ctx := trace.StartSpanFromContext(
			ctx,
			method,
			trace.TagSpanKind("client"),
			trace.TagComponent("grpc"),
		)
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

// loggerUnaryClientInterceptor gRPC客户端日志中间件
func loggerUnaryClientInterceptor(name string, accessInterceptorLevel string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		spbStatus := ecode.ExtractCodes(err)
		if err != nil {
			// 只记录系统级别错误
			if spbStatus.Code < ecode.EcodeNum {
				// 只记录系统级别错误
				logger.Logger.WithFields(logrus.Fields{
					"code":   spbStatus.Code,
					"msg":    spbStatus.Message,
					"name":   name,
					"method": method,
					"time":   time.Since(beg),
					"req":    json.RawMessage(util.JSONString(req)),
					"reply":  json.RawMessage(util.JSONString(reply)),
				}).Error("access_unary")

			} else {
				// 业务报错只做warning
				logger.Logger.WithFields(logrus.Fields{
					"code":   spbStatus.Code,
					"msg":    spbStatus.Message,
					"name":   name,
					"method": method,
					"time":   time.Since(beg),
					"req":    json.RawMessage(util.JSONString(req)),
					"reply":  json.RawMessage(util.JSONString(reply)),
				}).Warn("access")
			}
			return err
		}
		if accessInterceptorLevel == "info" {
			logger.Logger.WithFields(logrus.Fields{
				"code":   spbStatus.Code,
				"msg":    spbStatus.Message,
				"name":   name,
				"method": method,
				"time":   time.Since(beg),
				"req":    json.RawMessage(util.JSONString(req)),
				"reply":  json.RawMessage(util.JSONString(reply)),
			}).Info("access_unary")
		}
		return nil
	}
}

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
