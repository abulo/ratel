package trace

import (
	"context"

	"github.com/abulo/ratel/v3/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

var (
	// String ...
	String = log.String
)

// SetGlobalTracer ...
func SetGlobalTracer(tracer opentracing.Tracer) {
	logger.Logger.Info("set global tracer")
	opentracing.SetGlobalTracer(tracer)
}

// Start ...
func StartSpanFromContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

// SpanFromContext ...
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}
