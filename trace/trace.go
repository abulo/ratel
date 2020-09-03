package trace

import (
	"context"
	"io"
	"net"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc/metadata"
)

var (
	// String ...
	String = log.String
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

// Start ...
func StartSpanFromContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

// SpanFromContext ...
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// CustomTag ...
func CustomTag(key string, val interface{}) opentracing.Tag {
	return opentracing.Tag{
		Key:   key,
		Value: val,
	}
}

// TagComponent ...
func TagComponent(component string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "component",
		Value: component,
	}
}

// TagSpanKind ...
func TagSpanKind(kind string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "span.kind",
		Value: kind,
	}
}

// TagSpanURL ...
func TagSpanURL(url string) opentracing.Tag {
	return opentracing.Tag{
		Key:   "span.url",
		Value: url,
	}
}

// FromIncomingContext ...
func FromIncomingContext(ctx context.Context) opentracing.StartSpanOption {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		return NullStartSpanOption{}
	}
	return ext.RPCServerOption(sc)
}

// HeaderExtractor ...
func HeaderExtractor(hdr map[string][]string) opentracing.StartSpanOption {
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: hdr})
	if err != nil {
		return NullStartSpanOption{}
	}
	return opentracing.ChildOf(sc)
}

type hdrRequestKey struct{}

// HeaderInjector ...
func HeaderInjector(ctx context.Context, hdr map[string][]string) context.Context {
	span := opentracing.SpanFromContext(ctx)
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, MetadataReaderWriter{MD: hdr})
	if err != nil {
		span.LogFields(log.String("event", "inject failed"), log.Error(err))
		return ctx
	}
	return context.WithValue(ctx, hdrRequestKey{}, hdr)
}

// MetadataExtractor ...
func MetadataExtractor(md map[string][]string) opentracing.StartSpanOption {
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		return NullStartSpanOption{}
	}
	return opentracing.ChildOf(sc)
}

// MetadataInjector ...
func MetadataInjector(ctx context.Context, md metadata.MD) context.Context {
	span := opentracing.SpanFromContext(ctx)
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, MetadataReaderWriter{MD: md})
	if err != nil {
		span.LogFields(log.String("event", "inject failed"), log.Error(err))
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// NullStartSpanOption ...
type NullStartSpanOption struct{}

// Apply ...
func (sso NullStartSpanOption) Apply(options *opentracing.StartSpanOptions) {}

// MetadataReaderWriter ...
type MetadataReaderWriter struct {
	MD map[string][]string
}

// Set ...
func (w MetadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

// ForeachKey ...
func (w MetadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
