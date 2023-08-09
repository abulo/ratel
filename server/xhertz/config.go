package xhertz

import (
	"fmt"
)

// ModName ..
const ModName = "server.hertz"

// Config HTTP config
type Config struct {
	Host                      string
	Port                      int
	Deployment                string
	DisableMetric             bool
	DisableTrace              bool
	DisableSlowQuery          bool
	ServiceAddress            string
	SlowQueryThresholdInMilli int64
}

// New ...
func New() *Config {
	return &Config{
		SlowQueryThresholdInMilli: 500, // 500ms
	}
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

// WithDisableSlowQuery ...
func (config *Config) WithDisableSlowQuery(disableSlowQuery bool) *Config {
	config.DisableSlowQuery = disableSlowQuery
	return config
}

// WithDisableMetric  ...
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

// WithSlowQueryThresholdInMilli WithPort ...
func (config *Config) WithSlowQueryThresholdInMilli(milli int64) *Config {
	config.SlowQueryThresholdInMilli = milli
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	serverInstance := newServer(config)
	if !config.DisableSlowQuery {
		//慢日志查询
		serverInstance.Use(recoverMiddleware(config.SlowQueryThresholdInMilli))
	}
	if !config.DisableMetric {
		serverInstance.Use(metricServerInterceptor())
	}
	if !config.DisableTrace {
		serverInstance.Use(traceServerInterceptor())
	}
	return serverInstance
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
