package redis

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// OpenTraceHook ...
type OpenTraceHook struct {
	redis.Hook
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	DB            int
	Addr          string
}

// DialHook 返回redis连接hook hook: 原始redis连接hook
func (op OpenTraceHook) DialHook(hook redis.DialHook) redis.DialHook {
	return hook
}

// ProcessHook 返回redis命令处理hook hook: 原始redis命令处理hook
func (op OpenTraceHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		ctx, err := op.BeforeProcess(ctx, cmd)
		if err != nil {
			return err
		}
		hook(ctx, cmd)
		return op.AfterProcess(ctx, cmd)
	}
}

// ProcessPipelineHook 返回redis管道命令处理hook hook: 原始redis管道命令处理hook
func (op OpenTraceHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmd []redis.Cmder) error {
		ctx, err := op.BeforeProcessPipeline(ctx, cmd)
		if err != nil {
			return err
		}
		hook(ctx, cmd)
		return op.AfterProcessPipeline(ctx, cmd)
	}
}

// CmdStart ...
type CmdStart string

// RequestCmdStart ...
const RequestCmdStart = CmdStart("start")

// Caller 获取调用者信息 skip: 调用栈跳过的层数
func Caller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"path": file + ":" + cast.ToString(lineNo),
		"func": name,
	}
}

// BeforeProcess 在执行redis命令前处理 ctx: 上下文 cmd: redis命令
func (op OpenTraceHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	b := make([]byte, 32)
	b = appendCmd(b, cmd)
	ctx = getCtx(ctx)
	if !op.DisableTrace {
		call := Caller(11)
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.LogFields(log.Object("call", call))
			span.LogFields(log.String("cmd", cast.ToString(b)))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	if !op.DisableMetric {
		start := time.Now()
		ctx = context.WithValue(ctx, RequestCmdStart, start)
	}

	return ctx, nil
}

// AfterProcess 在执行redis命令后处理 ctx: 上下文 cmd: redis命令
func (op OpenTraceHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	ctx = getCtx(ctx)
	if !op.DisableTrace {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			defer span.Finish()
		}
	}
	if !op.DisableMetric {
		start := ctx.Value(RequestCmdStart)
		var cost time.Duration
		if start == nil {
			cost = time.Since(time.Now())
		} else {
			cost = time.Since(start.(time.Time))
		}
		if cmd.Err() != nil {
			metric.LibHandleCounter.WithLabelValues("redis", cast.ToString(op.DB), op.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("redis", cast.ToString(op.DB), op.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("redis", cast.ToString(op.DB), op.Addr).Observe(cost.Seconds())
	}
	return nil
}

// BeforeProcessPipeline 在执行redis管道命令前处理 ctx: 上下文 cmds: redis命令数组
func (op OpenTraceHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx = getCtx(ctx)

	const numCmdLimit = 100
	const numNameLimit = 10

	seen := make(map[string]struct{}, len(cmds))
	unqNames := make([]string, 0, len(cmds))

	b := make([]byte, 0, 32*len(cmds))

	for i, cmd := range cmds {
		if i > numCmdLimit {
			break
		}
		if i > 0 {
			b = append(b, '\n')
		}
		b = appendCmd(b, cmd)

		if len(unqNames) >= numNameLimit {
			continue
		}
		name := cmd.FullName()
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			unqNames = append(unqNames, name)
		}
	}
	if !op.DisableTrace {
		call := Caller(11)
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.LogFields(log.Object("call", call))
			span.LogFields(log.String("cmds", cast.ToString(b)))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	if !op.DisableMetric {
		start := time.Now()
		ctx = context.WithValue(ctx, RequestCmdStart, start)
	}

	return ctx, nil
}

// AfterProcessPipeline 在执行redis管道命令后处理 ctx: 上下文 cmds: redis命令数组
func (op OpenTraceHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	ctx = getCtx(ctx)
	if !op.DisableTrace {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			defer span.Finish()
		}
	}
	if !op.DisableMetric {
		start := ctx.Value(RequestCmdStart)
		cost := time.Since(start.(time.Time))
		if cmds != nil {
			metric.LibHandleCounter.WithLabelValues("redis", cast.ToString(op.DB), op.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("redis", cast.ToString(op.DB), op.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("redis", cast.ToString(op.DB), op.Addr).Observe(cost.Seconds())
	}
	return nil
}

// appendCmd 将redis命令追加到字节数组 b: 目标字节数组 cmd: redis命令
func appendCmd(b []byte, cmd redis.Cmder) []byte {
	const lenLimit = 64

	for i, arg := range cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}

		start := len(b)
		b = AppendArg(b, arg)
		if len(b)-start > lenLimit {
			b = append(b[:start+lenLimit], "..."...)
		}
	}

	if err := cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	}

	return b
}

// AppendArg 将参数追加到字节数组 b: 目标字节数组 v: 要追加的参数
func AppendArg(b []byte, v any) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		return appendUTF8String(b, v)
	case []byte:
		return appendUTF8String(b, cast.ToString(v))
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}

// appendUTF8String 将UTF8字符串追加到字节数组 b: 目标字节数组 s: 要追加的字符串
func appendUTF8String(b []byte, s string) []byte {
	for _, r := range s {
		b = appendRune(b, r)
	}
	return b
}

// appendRune 将rune追加到字节数组 b: 目标字节数组 r: 要追加的rune
func appendRune(b []byte, r rune) []byte {
	if r < utf8.RuneSelf {
		switch c := byte(r); c {
		case '\n':
			return append(b, "\\n"...)
		case '\r':
			return append(b, "\\r"...)
		default:
			return append(b, c)
		}
	}

	l := len(b)
	b = append(b, make([]byte, utf8.UTFMax)...)
	n := utf8.EncodeRune(b[l:l+utf8.UTFMax], r)
	b = b[:l+n]

	return b
}
