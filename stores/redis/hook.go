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

// OpenTraceHook 是 Redis 链路追踪钩子，包装 go-redis 的连接、单命令与管道命令处理流程，
// 为每段流程生成 OpenTelemetry span。零值可用，直接通过 AddHook 挂载即可。
type OpenTraceHook struct {
	redis.Hook
}

// DialHook 包装底层连接钩子，在建立连接前后创建 redis.dial span；连接失败时错误原样返回。
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

// maxDBStatementRunes 为 span 中命令语句属性的最大字符长度，超长截断以避免与命令体等长的分配。
const maxDBStatementRunes = 512

// truncateDBStatement 按字符（rune）截断命令语句，避免多字节字符被从中间切断；长度未超限时原样返回。
// PERF-1-3 [P1] 大 key / 多参数命令不再把整段参数写入 span 属性（预期收益需实测验证）。
func truncateDBStatement(statement string) string {
	// 字节数不超限时直接返回，避免为每条命令额外分配 rune 切片。
	if len(statement) <= maxDBStatementRunes {
		return statement
	}
	runes := []rune(statement)
	if len(runes) <= maxDBStatementRunes {
		return statement
	}
	return string(runes[:maxDBStatementRunes])
}

// ProcessHook 包装单命令处理钩子，创建以命令名命名的 span，记录调用位置与命令语句；
// 命令执行错误原样透传给上层。
func (op OpenTraceHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	// 基础属性在钩子注册期固定一次，命令级属性在闭包内单独构建，避免跨命令累积。
	baseAttrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("redis"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, cmd redis.Cmder) error {
		fn, file, line := call.Caller(11)
		// PERF-1-3 [P1] 每条命令使用独立切片并复制基础属性，不再复用闭包外层切片的底层数组，
		// 修复属性跨命令累积导致的残留与单调内存增长。
		attrs := make([]attribute.KeyValue, 0, len(baseAttrs)+4)
		attrs = append(attrs, baseAttrs...)
		attrs = append(attrs,
			semconv.CodeFunction(fn),
			semconv.CodeFilepath(file),
			semconv.CodeLineNumber(line),
		)
		attrs = append(attrs, semconv.DBStatement(truncateDBStatement(rediscmd.CmdString(cmd))))
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

// ProcessPipelineHook 包装管道命令处理钩子，创建 redis.pipeline span，记录命令数量与汇总语句；
// 管道执行错误原样透传给上层。
func (op OpenTraceHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	tracer := globalTrace.NewTracer(trace.SpanKindServer)
	// 基础属性在钩子注册期固定一次，命令级属性在闭包内单独构建，避免跨命令累积。
	baseAttrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("redis"),
		semconv.HostName(hostname.Hostname()),
	}
	return func(ctx context.Context, cmds []redis.Cmder) error {
		fn, file, line := call.Caller(11)
		// PERF-1-3 [P1] 每条管道命令使用独立切片并复制基础属性，不再复用闭包外层切片的底层数组。
		attrs := make([]attribute.KeyValue, 0, len(baseAttrs)+5)
		attrs = append(attrs, baseAttrs...)
		attrs = append(attrs,
			semconv.CodeFunction(fn),
			semconv.CodeFilepath(file),
			semconv.CodeLineNumber(line),
			attribute.Int("db.redis.num_cmd", len(cmds)),
		)
		summary, cmdsString := rediscmd.CmdsString(cmds)
		attrs = append(attrs, semconv.DBStatement(truncateDBStatement(cmdsString)))
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

// MetricsHook 是 Redis 指标采集钩子，记录连接、单命令与管道命令的执行耗时，以及成功/失败计数。
// 零值可用，直接通过 AddHook 挂载即可。
type MetricsHook struct {
	redis.Hook
}

// DialHook 包装连接钩子，测量建连耗时并写入连接耗时直方图。
func (op MetricsHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := hook(ctx, network, addr)
		dur := time.Since(start)
		metric.LibHandleHistogram.WithLabelValues("redis", network, addr).Observe(milliseconds(dur))
		return conn, err
	}
}

// ProcessHook 包装单命令钩子，按命令成功或失败累加计数，并记录命令执行耗时。
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

// ProcessPipelineHook 包装管道钩子，按管道成功或失败累加计数，并记录管道执行耗时。
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
