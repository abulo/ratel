package trace

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	jaegerHeader        = "x-dy-traceid"
	separator           = ":"
	traceID128bitsWidth = 128 / 4
	spanIDWidth         = 64 / 4

	idPaddingChar = "0"

	flagsDebug      = 0x02
	flagsSampled    = 0x01
	flagsNotSampled = 0x00

	// PERF-3-16 [P2]：三个 flags 组合取值恒定，预存十六进制字符串常量，去掉每次 Inject 的 fmt.Sprintf("%x")。
	// 注意 "%x" 不补零，取值必须与原输出逐字一致：0x03→"3"、0x01→"1"、0x00→"0"。
	flagDebugSampledStr = "3"
	flagSampledStr      = "1"
	flagNotSampledStr   = "0"

	deprecatedParentSpanID = "0"
)

var (
	empty = trace.SpanContext{}

	errMalformedTraceContextVal = errors.New("header value of uber-trace-id should contain four different part separated by : ")
	errInvalidTraceIDLength     = errors.New("invalid trace id length, must be either 16 or 32")
	errMalformedTraceID         = errors.New("cannot decode trace id from header, should be a string of hex, lowercase trace id can't be all zero")
	errInvalidSpanIDLength      = errors.New("invalid span id length, must be 16")
	errMalformedSpanID          = errors.New("cannot decode span id from header, should be a string of hex, lowercase span id can't be all zero")
	errMalformedFlag            = errors.New("cannot decode flag")
)

// Jaeger 实现 OpenTelemetry 文本传播器，按 Jaeger uber-trace-id 格式在载体上序列化与反序列化 SpanContext。
//
// 头格式为 trace-id:span-id:parent-span-id:flags；父 SpanID 已废弃，Inject 时固定填充占位零值。
type Jaeger struct{}

var _ propagation.TextMapPropagator = &Jaeger{}

// Inject 将当前上下文的 SpanContext 按 Jaeger 格式写入载体；TraceID 或 SpanID 无效时不写入任何头。
func (jaeger Jaeger) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanFromContext(ctx).SpanContext()
	// PERF-3-16 [P2]：切片按最大 4 段预分配容量，避免多次 append 触发扩容；flags 直接使用预存常量字符串。
	headers := make([]string, 0, 4)
	if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
		return
	}
	headers = append(headers, sc.TraceID().String(), sc.SpanID().String(), deprecatedParentSpanID)
	if debugFromContext(ctx) {
		headers = append(headers, flagDebugSampledStr)
	} else if sc.IsSampled() {
		headers = append(headers, flagSampledStr)
	} else {
		headers = append(headers, flagNotSampledStr)
	}

	carrier.Set(jaegerHeader, strings.Join(headers, separator))
}

// Extract 从载体读取 Jaeger 头并还原远程 SpanContext；头缺失或格式非法时返回原上下文。
func (jaeger Jaeger) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	if h := carrier.Get(jaegerHeader); h != "" {
		ctx, sc, err := extract(ctx, h)
		if err == nil && sc.IsValid() {
			return trace.ContextWithRemoteSpanContext(ctx, sc)
		}
	}

	return ctx
}

func extract(ctx context.Context, headerVal string) (context.Context, trace.SpanContext, error) {
	var (
		scc = trace.SpanContextConfig{}
		err error
	)

	parts := strings.Split(headerVal, separator)
	if len(parts) != 4 {
		return ctx, empty, errMalformedTraceContextVal
	}

	if parts[0] != "" {
		id := parts[0]
		if len(id) > traceID128bitsWidth {
			return ctx, empty, errInvalidTraceIDLength
		}
		// 长度不足 32 时在左侧补零，对齐 128 位 TraceID。
		if len(id) < traceID128bitsWidth {
			padCharCount := traceID128bitsWidth - len(id)
			id = strings.Repeat(idPaddingChar, padCharCount) + id
		}
		scc.TraceID, err = trace.TraceIDFromHex(id)
		if err != nil {
			return ctx, empty, errMalformedTraceID
		}
	}

	if parts[1] != "" {
		id := parts[1]
		if len(id) > spanIDWidth {
			return ctx, empty, errInvalidSpanIDLength
		}
		// 长度不足 16 时在左侧补零，对齐 64 位 SpanID。
		if len(id) < spanIDWidth {
			padCharCount := spanIDWidth - len(id)
			id = strings.Repeat(idPaddingChar, padCharCount) + id
		}
		scc.SpanID, err = trace.SpanIDFromHex(id)
		if err != nil {
			return ctx, empty, errMalformedSpanID
		}
	}

	// 第三段 parent-span-id 已废弃，直接跳过。

	if parts[3] != "" {
		flagStr := parts[3]
		flag, err := strconv.ParseInt(flagStr, 16, 64)
		if err != nil {
			return ctx, empty, errMalformedFlag
		}
		if flag&flagsSampled == flagsSampled {
			// 采样位置位时再判断 debug 位是否同时置位。
			if flag&flagsDebug == flagsDebug {
				scc.TraceFlags |= trace.FlagsSampled
				ctx = withDebug(ctx, true)
			} else {
				scc.TraceFlags |= trace.FlagsSampled
			}
		}
		// 其余标志位（含 firehose）在 W3C Trace Context 中没有对应位，忽略不处理。
	}
	return ctx, trace.NewSpanContext(scc), nil
}

// Fields 返回 Inject 会写入的 Jaeger 头键名。
func (jaeger Jaeger) Fields() []string {
	return []string{jaegerHeader}
}

type jaegerKeyType int

const (
	debugKey jaegerKeyType = iota
)

// withDebug 返回携带 debug 标记的上下文副本。
func withDebug(parent context.Context, debug bool) context.Context {
	return context.WithValue(parent, debugKey, debug)
}

// debugFromContext 读取上下文中保存的 debug 标记；上下文为空或未写入该标记时返回 false。
func debugFromContext(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	if debug, ok := ctx.Value(debugKey).(bool); ok {
		return debug
	}
	return false
}
