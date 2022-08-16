package xgin

import (
	"fmt"

	"github.com/abulo/ratel/v3/gin"
)

// ModName ..
const ModName = "server.gin"

// Config HTTP config
type Config struct {
	Host                      string
	Port                      int
	Deployment                string
	Mode                      string
	DisableMetric             bool
	DisableTrace              bool
	DisableSlowQuery          bool
	ServiceAddress            string // ServiceAddress service address in registry info, default to 'Host:Port'
	SlowQueryThresholdInMilli int64
}

// New ...
func New() *Config {
	return &Config{
		Mode:                      gin.ReleaseMode,
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

// WithMode ...
func (config *Config) With(mode string) *Config {
	config.Mode = mode
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
	server := newServer(config)
	server.Use(gin.Recovery())
	if !config.DisableSlowQuery {
		//慢日志查询
		server.Use(recoverMiddleware(config.SlowQueryThresholdInMilli))
	}
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
