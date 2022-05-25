package gin

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

//ModName ..
const ModName = "server.gin"

// Config HTTP config
type Config struct {
	Host                      string
	Port                      int
	Deployment                string
	Mode                      string
	DisableMetric             bool
	DisableTrace              bool
	ServiceAddress            string // ServiceAddress service address in registry info, default to 'Host:Port'
	SlowQueryThresholdInMilli int64
	logger                    *logrus.Logger
}

func New() *Config {
	return &Config{}
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// WithDeployment ...
func (config *Config) WithDeployment(deployment string) *Config {
	config.Deployment = deployment
	return config
}

// WithMode ...
func (config *Config) With(mode string) *Config {
	config.Mode = mode
	return config
}

// WithDisableTrace ...
func (config *Config) WithDisableMetric(disableMetric bool) *Config {
	config.DisableMetric = disableMetric
	return config
}

// WithDisableTrace ...
func (config *Config) WithDisableTrace(disableTrace bool) *Config {
	config.DisableTrace = disableTrace
	return config
}

// WithServiceAddress ...
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// WithPort ...
func (config *Config) WithSlowQueryThresholdInMilli(milli int64) *Config {
	config.SlowQueryThresholdInMilli = milli
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *logrus.Logger) *Config {
	config.logger = logger
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	server := newServer(config)
	//慢日志查询
	server.Use(recoverMiddleware(config.logger, config.SlowQueryThresholdInMilli))
	if !config.DisableMetric {
		server.Use(metricServerInterceptor())
	}
	if !config.DisableTrace {
		server.Use(traceServerInterceptor())
	}
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
