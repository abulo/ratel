package jaeger

import (
	"os"
	"time"

	"github.com/abulo/ratel/v3/env"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/stack"
	"github.com/abulo/ratel/v3/util"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

type JConfig struct {

	// 	[jaeger]
	// EnableRPCMetrics= true
	// [jaeger.Reporter]
	// LocalAgentHostPort = "127.0.0.1:6831"
	// LogSpans = true
	// [jaeger.Sampler]
	// Param = 0.0001

	EnableRPCMetrics   bool
	LocalAgentHostPort string
	LogSpans           bool
	Param              float64
	PanicOnError       bool
}

func NewJaeger() *JConfig {
	return &JConfig{}
}

func (jConfig *JConfig) Build() *Config {
	agentAddr := "127.0.0.1:6831"
	headerName := "x-trace-id"
	if addr := os.Getenv("JAEGER_AGENT_ADDR"); addr != "" {
		agentAddr = addr
	}
	if util.Empty(jConfig.LocalAgentHostPort) {
		agentAddr = jConfig.LocalAgentHostPort
	}
	return &Config{
		ServiceName: env.Name(),
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: jConfig.Param,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:            jConfig.LogSpans,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddr,
		},
		EnableRPCMetrics: true,
		Headers: &jaeger.HeadersConfig{
			TraceBaggageHeaderPrefix: "ctx-",
			TraceContextHeaderName:   headerName,
		},
		tags: []opentracing.Tag{
			{Key: "hostname", Value: env.HostName()},
		},
		PanicOnError: jConfig.EnableRPCMetrics,
	}
}

// Config ...
type Config struct {
	ServiceName      string
	Sampler          *jaegerConfig.SamplerConfig
	Reporter         *jaegerConfig.ReporterConfig
	Headers          *jaeger.HeadersConfig
	EnableRPCMetrics bool
	tags             []opentracing.Tag
	options          []jaegerConfig.Option
	PanicOnError     bool
}

// WithTag ...
func (config *Config) WithTag(tags ...opentracing.Tag) *Config {
	if config.tags == nil {
		config.tags = make([]opentracing.Tag, 0)
	}
	config.tags = append(config.tags, tags...)
	return config
}

// WithOption ...
func (config *Config) WithOption(options ...jaegerConfig.Option) *Config {
	if config.options == nil {
		config.options = make([]jaegerConfig.Option, 0)
	}
	config.options = append(config.options, options...)
	return config
}

// Build ...
func (config *Config) Build(options ...jaegerConfig.Option) opentracing.Tracer {
	var configuration = jaegerConfig.Configuration{
		ServiceName: config.ServiceName,
		Sampler:     config.Sampler,
		Reporter:    config.Reporter,
		RPCMetrics:  config.EnableRPCMetrics,
		Headers:     config.Headers,
		Tags:        config.tags,
	}

	logger.Logger.WithFields(logrus.Fields{
		"configuration": configuration,
	}).Info("trace")

	tracer, closer, err := configuration.NewTracer(config.options...)
	if err != nil {
		if config.PanicOnError {
			logger.Logger.WithFields(logrus.Fields{
				"jaeger": err,
			}).Info("new jaeger")
		} else {
			logger.Logger.WithFields(logrus.Fields{
				"jaeger": err,
			}).Info("new jaeger")
		}
	}
	stack.Register(closer.Close)
	return tracer
}
