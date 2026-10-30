package trace

import (
	"context"
	"time"

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

// Config 保存链路追踪导出器的配置项，由 NewConfig 结合 Option 初始化。
type Config struct {
	// Name 为上报到追踪系统的服务名称。
	Name     string
	// Endpoint 为 OTLP 采集器端点，格式为 host:port。
	Endpoint string
	// Sampler 为采样比例，取值范围 0 到 1；采样决策基于父链路结果进行父优先采样。
	Sampler  float64
	// Protocol 为导出协议，取值 http 或 grpc；NewConfig 缺省为 http。
	Protocol string
}

// Option 是 Config 的配置选项函数，以函数式选项模式在 NewConfig 中修改配置。
type Option func(r *Config)

// WithName 设置上报到追踪系统的服务名称。
func WithName(name string) Option {
	return func(r *Config) {
		r.Name = name
	}
}

// WithEndpoint 设置 OTLP 采集器端点，格式为 host:port。
func WithEndpoint(endpoint string) Option {
	return func(r *Config) {
		r.Endpoint = endpoint
	}
}

// WithSampler 设置采样比例，取值范围 0 到 1。
func WithSampler(sampler float64) Option {
	return func(r *Config) {
		r.Sampler = sampler
	}
}

// WithProtocol 设置导出协议，取值 http 或 grpc；不传时 NewConfig 默认使用 http。
func WithProtocol(protocol string) Option {
	return func(r *Config) {
		r.Protocol = protocol
	}
}

// NewConfig 以给定选项构造配置，Protocol 未显式设置时缺省为 http。
func NewConfig(opts ...Option) *Config {
	c := &Config{
		Protocol: "http",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Build 根据配置创建 TracerProvider，并将其设置为 OpenTelemetry 全局默认提供者。
//
// 导出器或资源构造失败时通过日志器 Panic 终止进程；成功时返回提供者，调用方也可直接使用全局 Tracer。
func (config *Config) Build() trace.TracerProvider {

	var exp *otlptrace.Exporter
	var err error
	if config.Protocol == "http" {
		exp, err = otlptracehttp.New(context.TODO(),
			otlptracehttp.WithEndpoint(config.Endpoint),
			// 内网采集端点未启用 TLS，使用明文 HTTP 传输。
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithRetry(
				otlptracehttp.RetryConfig{
					Enabled:         true,
					InitialInterval: 1 * time.Second,
					MaxInterval:     10 * time.Second,
					MaxElapsedTime:  30 * time.Second,
				},
			),
		)
	} else {
		exp, err = otlptracegrpc.New(context.TODO(),
			otlptracegrpc.WithEndpoint(config.Endpoint),
			// 内网采集端点未启用 TLS，使用明文 gRPC 传输。
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithReconnectionPeriod(5*time.Second),
			otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: 1 * time.Second,
				MaxInterval:     10 * time.Second,
				MaxElapsedTime:  30 * time.Second,
			}),
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
