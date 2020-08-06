package trace

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// TraceConfig 实列
type TraceConfig struct {
	ServiceName         string
	SamplingServerURL   string
	SamplingParam       float64
	SamplingType        string
	BufferFlushInterval time.Duration
	LogSpans            bool
	QueueSize           int
	PropagationFormat   string
}

//jagerLogger 日志
type jagerLogger struct {
}

//Error 错误日志
func (s *jagerLogger) Error(msg string) {
	logrus.Error(msg)
}

//Infof 日志记录
func (s *jagerLogger) Infof(msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}

// Errorf - 日志记录
func (l *jagerLogger) Errorf(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...)
}

//InitConfig 初始化 config
func InitConfig(url string) TraceConfig {
	_, _, err := net.SplitHostPort(url)
	if url == "" || err != nil {
		return initNilConfig()
	}
	return initJaegerConfig(url)
}

// initNilConfig 返回一个空的配置
func initNilConfig() TraceConfig {
	return TraceConfig{}
}

// initJaegerConfig 返回一个jaeger的配置
func initJaegerConfig(serverURL string) TraceConfig {
	return TraceConfig{
		ServiceName:         "jaeger",
		SamplingServerURL:   serverURL,
		SamplingParam:       1.0,
		SamplingType:        "const",
		BufferFlushInterval: time.Second,
		LogSpans:            false,
		QueueSize:           1000,
	}
}

// Tracing tracing的基类
type Tracing struct {
	config TraceConfig
	tracer opentracing.Tracer
	closer io.Closer
}

type noopCloser struct{}

func (n noopCloser) Close() error { return nil }

// New 按照配置信息初始化tracing
func New(cfg TraceConfig) *Tracing {
	return &Tracing{
		config: cfg,
	}
}

// Setup 根据tracing的属性选择合适的opentracing客户端
func (t *Tracing) Setup() (err error) {
	t.tracer, t.closer, err = t.buildJaeger(t.config)
	// 设置全局的tracer
	opentracing.SetGlobalTracer(t.tracer)
	return
}

// Close 关闭tracer
func (t *Tracing) Close() {
	if t.closer != nil {
		t.closer.Close()
	}
}

// buildJaeger 创建jaeger system
func (t *Tracing) buildJaeger(cfg TraceConfig) (opentracing.Tracer, io.Closer, error) {
	svrName := cfg.ServiceName
	conf := config.Configuration{
		ServiceName: svrName,
		Sampler: &config.SamplerConfig{
			Param: cfg.SamplingParam,
			Type:  cfg.SamplingType,
		},
		Reporter: &config.ReporterConfig{
			QueueSize:           cfg.QueueSize,
			LocalAgentHostPort:  cfg.SamplingServerURL,
			BufferFlushInterval: cfg.BufferFlushInterval,
			LogSpans:            cfg.LogSpans,
		},
	}

	tracerMetrics := jaeger.NewMetrics(metrics.NullFactory, nil)
	sampler, err := conf.Sampler.NewSampler(svrName, tracerMetrics)
	tracerLogger := &jagerLogger{}
	if err != nil {
		return nil, nil, err
	}

	reporter, err := conf.Reporter.NewReporter(svrName, tracerMetrics, tracerLogger)
	if err != nil {
		return nil, nil, err
	}

	var (
		tracer opentracing.Tracer
		closer io.Closer
	)

	tracer, closer = jaeger.NewTracer(svrName, sampler, reporter,
		jaeger.TracerOptions.Metrics(tracerMetrics),
		jaeger.TracerOptions.Logger(tracerLogger),
	)

	return tracer, closer, nil
}

// ToContext 将span写入request的context中
func ToContext(r *http.Request, span opentracing.Span) *http.Request {
	return r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
}

// FromContext 从context中取出span
func FromContext(ctx context.Context, name string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, name)
	return span
}
