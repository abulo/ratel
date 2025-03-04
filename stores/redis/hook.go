package redis

import (
	"context"
	"net"
	"time"

	"github.com/abulo/ratel/v3/core/call"
	"github.com/abulo/ratel/v3/core/hostname"
	"github.com/abulo/ratel/v3/core/metric"
	globalTrace "github.com/abulo/ratel/v3/core/trace"
	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// OpenTraceHook ...
type OpenTraceHook struct {
	redis.Hook
}

// DialHook 返回redis连接hook hook: 原始redis连接hook
func (op OpenTraceHook) DialHook(hook redis.DialHook) redis.DialHook {

	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("redis"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		ctx, span := tracer.Start(ctx, "redis.dial", globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer span.End()
		conn, err := hook(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// ProcessHook 返回redis命令处理hook hook: 原始redis命令处理hook
func (op OpenTraceHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("redis"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, cmd redis.Cmder) error {
		fn, file, line := call.Caller(11)
		attrs = append(attrs,
			semconv.CodeFunction(fn),
			semconv.CodeFilepath(file),
			semconv.CodeLineNumber(line),
		)
		cmdString := rediscmd.CmdString(cmd)
		attrs = append(attrs, semconv.DBStatement(cmdString))
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		ctx, span := tracer.Start(ctx, "redis."+cmd.FullName(), globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer span.End()
		if err := hook(ctx, cmd); err != nil {
			return err
		}
		return nil
	}
}

// ProcessPipelineHook 返回redis管道命令处理hook hook: 原始redis管道命令处理hook
func (op OpenTraceHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("redis"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, cmds []redis.Cmder) error {
		fn, file, line := call.Caller(11)
		attrs = append(attrs,
			semconv.CodeFunction(fn),
			semconv.CodeFilepath(file),
			semconv.CodeLineNumber(line),
			attribute.Int("db.redis.num_cmd", len(cmds)),
		)
		summary, cmdsString := rediscmd.CmdsString(cmds)
		attrs = append(attrs, semconv.DBStatement(cmdsString))
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		ctx, span := tracer.Start(ctx, "redis.pipeline "+summary, globalTrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer span.End()
		if err := hook(ctx, cmds); err != nil {
			return err
		}
		return nil
	}
}

type MetricsHook struct {
	redis.Hook
}

func (op MetricsHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := hook(ctx, network, addr)
		dur := time.Since(start)
		metric.LibHandleHistogram.WithLabelValues("redis", network, addr).Observe(milliseconds(dur))
		return conn, err
	}
}

func (op MetricsHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := hook(ctx, cmd)
		dur := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("redis", cmd.FullName(), "ERR").Inc()
		} else {
			metric.LibHandleCounter.WithLabelValues("redis", cmd.FullName(), "OK").Inc()
		}
		metric.LibHandleHistogram.WithLabelValues("redis", cmd.FullName()).Observe(milliseconds(dur))
		return err
	}
}

func (op MetricsHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := hook(ctx, cmds)
		dur := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues("redis", "pipeline", "ERR").Inc()
		} else {
			metric.LibHandleCounter.WithLabelValues("redis", "pipeline", "OK").Inc()
		}
		metric.LibHandleHistogram.WithLabelValues("redis", "pipeline").Observe(milliseconds(dur))
		return err
	}
}

func milliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}
