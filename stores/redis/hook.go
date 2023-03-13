package redis

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"

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

func (op OpenTraceHook) DialHook(hook redis.DialHook) redis.DialHook {
	return hook
}

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

// BeforeProcess ...
func (op OpenTraceHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	b := make([]byte, 32)
	b = appendCmd(b, cmd)

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	// pc, file, lineNo, _ := runtime.Caller(10)
	// name := runtime.FuncForPC(pc).Name()
	// Path := file + ":" + cast.ToString(lineNo)
	// Func := name
	// fmt.Println(Path, Func, "ddd")

	if !op.DisableTrace {
		pc, file, lineNo, _ := runtime.Caller(5)
		name := runtime.FuncForPC(pc).Name()
		Path := file + ":" + cast.ToString(lineNo)
		Func := name
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("call.func", Func)
			span.SetTag("call.path", Path)
			span.LogFields(log.String("cmd", String(b)))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	if !op.DisableMetric {
		start := time.Now()
		ctx = context.WithValue(ctx, RequestCmdStart, start)
	}

	return ctx, nil
}

// AfterProcess ...
func (op OpenTraceHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !op.DisableTrace {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			defer span.Finish()
		}
	}

	if !op.DisableMetric {
		start := ctx.Value(RequestCmdStart)
		cost := time.Since(start.(time.Time))
		if cmd.Err() != nil {
			metric.LibHandleCounter.WithLabelValues("redis", cast.ToString(op.DB), op.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc("redis", cast.ToString(op.DB), op.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues("redis", cast.ToString(op.DB), op.Addr).Observe(cost.Seconds())
	}

	return nil
}

// BeforeProcessPipeline ...
func (op OpenTraceHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

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
		pc, file, lineNo, _ := runtime.Caller(5)
		name := runtime.FuncForPC(pc).Name()
		Path := file + ":" + cast.ToString(lineNo)
		Func := name
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			span.SetTag("call.func", Func)
			span.SetTag("call.path", Path)
			span.LogFields(log.String("cmds", String(b)))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	if !op.DisableMetric {
		start := time.Now()
		ctx = context.WithValue(ctx, RequestCmdStart, start)
	}

	return ctx, nil
}

// AfterProcessPipeline ...
func (op OpenTraceHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !op.DisableTrace {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			defer span.Finish()
		}
	}
	if !op.DisableMetric {
		start := ctx.Value(RequestCmdStart)
		cost := time.Since(start.(time.Time))
		// if cmds != nil {
		// metric.LibHandleCounter.WithLabelValues("redis", util.ToString(op.DB), op.Addr, "ERR").Inc()
		// } else {
		metric.LibHandleCounter.Inc("redis", cast.ToString(op.DB), op.Addr, "OK")
		// }
		metric.LibHandleHistogram.WithLabelValues("redis", cast.ToString(op.DB), op.Addr).Observe(cost.Seconds())
	}
	return nil
}

// String ...
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

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

// AppendArg ...
func AppendArg(b []byte, v any) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		return appendUTF8String(b, v)
	case []byte:
		return appendUTF8String(b, String(v))
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

func appendUTF8String(b []byte, s string) []byte {
	for _, r := range s {
		b = appendRune(b, r)
	}
	return b
}

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
