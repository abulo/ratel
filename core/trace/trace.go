package trace

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	otelOpentracing "go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// SetGlobalTracer 将给定提供者设为全局追踪器，并同时适配 OpenTelemetry 与 OpenTracing 两套体系。
//
// 传播格式组合 W3C TraceContext 与 Jaeger；重复调用会覆盖已注册的全局设置。
func SetGlobalTracer(tp trace.TracerProvider) {
	logger.Logger.Info("set global tracer")

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, Jaeger{})

	// 桥接 OpenTracing 旧体系，保持基于 opentracing 的已有调用方兼容。
	bridge, wrapperTracerProvider := otelOpentracing.NewTracerPair(tp.Tracer(""))
	bridge.SetTextMapPropagator(propagator)
	opentracing.SetGlobalTracer(bridge)

	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(wrapperTracerProvider)
}

// Tracer 包装 OpenTelemetry Tracer 并固定 Span 类型，简化跨上下文的 Span 启动与传播。
type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
}

// NewTracer 创建指定 SpanKind 的 Tracer；内部从全局 TracerProvider 获取名为 "ratel" 的 Tracer。
func NewTracer(kind trace.SpanKind) *Tracer {
	return &Tracer{tracer: otel.Tracer("ratel"), kind: kind}
}

// Start 启动新 Span 并返回携带 Span 的上下文。
//
// 服务端/消费者角色先从 carrier 提取上游上下文，客户端/生产者角色再向下游注入；
// carrier 为 nil 时不做跨上下文传播处理。
func (t *Tracer) Start(ctx context.Context, operation string, carrier propagation.TextMapCarrier, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if (t.kind == trace.SpanKindServer || t.kind == trace.SpanKindConsumer) && carrier != nil {
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	}
	opts = append(opts, trace.WithSpanKind(t.kind))

	ctx, span := t.tracer.Start(ctx, operation, opts...)

	if (t.kind == trace.SpanKindClient || t.kind == trace.SpanKindProducer) && carrier != nil {
		otel.GetTextMapPropagator().Inject(ctx, carrier)
	}
	return ctx, span
}
