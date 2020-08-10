package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type OpenTelemetryHook struct{}

var _ redis.Hook = OpenTelemetryHook{}

func (OpenTelemetryHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	b := make([]byte, 32)
	b = appendCmd(b, cmd)

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "redis")
			span.SetTag("redis.cmd", String(b))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return ctx, nil
}

func (OpenTelemetryHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			span.Finish()
		}
	}
	return nil
}

func (OpenTelemetryHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if trace {

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

		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan("redis", opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, "redis")
			span.SetTag("redis.cmds", String(b))
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	return ctx, nil
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (OpenTelemetryHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if trace {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			span.Finish()
		}
	}
	return nil
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

func AppendArg(b []byte, v interface{}) []byte {
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
