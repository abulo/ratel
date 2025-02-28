package trace

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Name     string
	Endpoint string
	Sampler  float64
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

func NewConfig(opts ...Option) *Config {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (config *Config) Build() trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Endpoint)))
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
