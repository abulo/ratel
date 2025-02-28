package xgin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ModName ..
const ModName = "server.gin"

// Config HTTP config
type Config struct {
	Host                     string
	Port                     int
	Deployment               string
	Mode                     string
	EnableMetric             bool
	EnableTrace              bool
	EnableSlowQuery          bool
	ServiceAddress           string // ServiceAddress service address in registry info, default to 'Host:Port'
	SlowQueryThresholdInMill int64
}

// New ...
func New() *Config {
	return &Config{
		Mode:                     gin.ReleaseMode,
		SlowQueryThresholdInMill: 500, // 500ms
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
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// WithEnableSlowQuery ...
func (config *Config) WithEnableSlowQuery(disableSlowQuery bool) *Config {
	config.EnableSlowQuery = disableSlowQuery
	return config
}

// WithEnableMetric  ...
func (config *Config) WithEnableMetric(disableMetric bool) *Config {
	config.EnableMetric = disableMetric
	return config
}

// WithEnableTrace ...
func (config *Config) WithEnableTrace(disableTrace bool) *Config {
	config.EnableTrace = disableTrace
	return config
}

// WithServiceAddress ...
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// WithSlowQueryThresholdInMilli WithPort ...
func (config *Config) WithSlowQueryThresholdInMilli(milli int64) *Config {
	config.SlowQueryThresholdInMill = milli
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	server := newServer(config)
	if config.Mode == gin.DebugMode {
		server.Use(gin.Logger())
	}
	server.Use(gin.Recovery())
	if config.EnableSlowQuery {
		//慢日志查询
		server.Use(recoverMiddleware(config.SlowQueryThresholdInMill))
	}
	if config.EnableMetric {
		server.Use(metricServerInterceptor())
	}
	if config.EnableTrace {
		server.Use(traceServerInterceptor())
	}
	return server
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
