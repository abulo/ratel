package trace

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Name     string
	Endpoint string
	Sampler  float64
	Protocol string
}

// Option 选项
type Option func(r *Config)

// WithName 设置名称
func WithName(name string) Option {
	return func(r *Config) {
		r.Name = name
	}
}

// WithEndpoint 设置端点
func WithEndpoint(endpoint string) Option {
	return func(r *Config) {
		r.Endpoint = endpoint
	}
}

// WithSampler 设置采样率
func WithSampler(sampler float64) Option {
	return func(r *Config) {
		r.Sampler = sampler
	}
}

func WithProtocol(protocol string) Option {
	return func(r *Config) {
		r.Protocol = protocol
	}
}

func NewConfig(opts ...Option) *Config {
	c := &Config{
		Protocol: "http",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (config *Config) Build() trace.TracerProvider {

	var exp *otlptrace.Exporter
	var err error
	if config.Protocol == "http" {
		exp, err = otlptracehttp.New(context.TODO(),
			otlptracehttp.WithEndpoint(config.Endpoint),
			otlptracehttp.WithInsecure(), // 如果是 http 而不是 https
		)
	} else {
		exp, err = otlptracegrpc.New(context.TODO(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
			otlptracegrpc.WithInsecure(), // 如果未启用 TLS
		)
	}
	if err != nil {
		logger.Logger.Panic("new jaeger", err)
		return nil
	}

	resource, err := resource.New(context.TODO(),
		resource.WithHost(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.Name),
		),
	)
	if err != nil {
		logger.Logger.Panic("new resource", err)
		return nil
	}
	tp := sdk.NewTracerProvider(
		sdk.WithSampler(sdk.ParentBased(sdk.TraceIDRatioBased(config.Sampler))),
		sdk.WithBatcher(exp),
		sdk.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	return tp
}
