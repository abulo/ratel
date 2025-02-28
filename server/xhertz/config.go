package xhertz

import (
	"fmt"
)

// ModName ..
const ModName = "server.hertz"

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

// Config HTTP config
type Config struct {
	Host                     string
	Port                     int
	Mode                     string
	Deployment               string
	EnableMetric             bool
	EnableTrace              bool
	EnableSlowQuery          bool
	ServiceAddress           string
	SlowQueryThresholdInMill int64
}

// New ...
func New() *Config {
	return &Config{
		SlowQueryThresholdInMill: 500, // 500ms
		Mode:                     DebugMode,
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

// WithMode ...
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	serverInstance := newServer(config)
	if config.EnableSlowQuery {
		//慢日志查询
		serverInstance.Use(recoverMiddleware(config.SlowQueryThresholdInMill))
	}
	if config.EnableMetric {
		serverInstance.Use(metricServerInterceptor())
	}
	if config.EnableTrace {
		serverInstance.Use(traceServerInterceptor())
	}
	return serverInstance
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
