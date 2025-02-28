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

// SetGlobalTracer ...
func SetGlobalTracer(tp trace.TracerProvider) {
	logger.Logger.Info("set global tracer")

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, Jaeger{})

	// be compatible with opentracing
	bridge, wrapperTracerProvider := otelOpentracing.NewTracerPair(tp.Tracer(""))
	bridge.SetTextMapPropagator(propagator)
	opentracing.SetGlobalTracer(bridge)

	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(wrapperTracerProvider)
}

// Tracer is otel span tracer
type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
}

// NewTracer create tracer instance
func NewTracer(kind trace.SpanKind) *Tracer {
	return &Tracer{tracer: otel.Tracer("ratel"), kind: kind}
}

// Start start tracing span
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
