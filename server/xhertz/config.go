package xhertz

import (
	"fmt"
)

// ModName 为 Hertz 服务端在框架内的模块标识，用于注册与监控归类。
const ModName = "server.hertz"

// Hertz 运行模式常量。
const (
	// DebugMode 为调试模式，输出详细调试信息。
	DebugMode = "debug"
	// ReleaseMode 为发布模式，关闭路由打印等调试输出。
	ReleaseMode = "release"
	// TestMode 为测试模式。
	TestMode = "test"
)

// Config 保存 Hertz HTTP 服务端的监听地址、运行模式与中间件开关等配置。
//
// 零值不可直接使用，须通过 New 结合各 With 方法构造；Config 为配置载体，
// 应在 Build 之前一次性设置完成，不保证并发安全。
type Config struct {
	// Host 为监听主机地址。
	Host                     string
	// Port 为监听端口。
	Port                     int
	// Mode 为 Hertz 运行模式，取值 DebugMode、ReleaseMode、TestMode，默认 DebugMode。
	Mode                     string
	// Deployment 为部署组，用于服务注册时的流量隔离。
	Deployment               string
	// EnableMetric 为 true 时启用服务指标采集中间件。
	EnableMetric             bool
	// EnableTrace 为 true 时启用服务端链路追踪中间件。
	EnableTrace              bool
	// EnableSlowQuery 为 true 时启用慢请求记录与 panic 恢复中间件。
	EnableSlowQuery          bool
	// ServiceAddress 为注册中心中登记的服务地址；为空时使用 Host:Port。
	ServiceAddress           string
	// SlowQueryThresholdInMill 为慢请求阈值，单位毫秒，超过该耗时的请求会被标记为慢查询，默认 500。
	SlowQueryThresholdInMill int64
}

// New 返回带有默认值的 Config：运行模式为 DebugMode，慢查询阈值为 500 毫秒。
func New() *Config {
	return &Config{
		SlowQueryThresholdInMill: 500,
		Mode:                     DebugMode,
	}
}

// WithHost 设置监听主机地址，返回 Config 以便链式调用。
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort 设置监听端口，返回 Config 以便链式调用。
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// WithDeployment 设置部署组，返回 Config 以便链式调用。
func (config *Config) WithDeployment(deployment string) *Config {
	config.Deployment = deployment
	return config
}

// WithEnableSlowQuery 设置是否启用慢请求记录与 panic 恢复中间件，返回 Config 以便链式调用。
func (config *Config) WithEnableSlowQuery(disableSlowQuery bool) *Config {
	config.EnableSlowQuery = disableSlowQuery
	return config
}

// WithEnableMetric 设置是否启用服务指标采集中间件，返回 Config 以便链式调用。
func (config *Config) WithEnableMetric(disableMetric bool) *Config {
	config.EnableMetric = disableMetric
	return config
}

// WithEnableTrace 设置是否启用服务端链路追踪中间件，返回 Config 以便链式调用。
func (config *Config) WithEnableTrace(disableTrace bool) *Config {
	config.EnableTrace = disableTrace
	return config
}

// WithServiceAddress 设置注册中心中登记的服务地址，返回 Config 以便链式调用。
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// WithSlowQueryThresholdInMilli 设置慢请求阈值（单位毫秒），返回 Config 以便链式调用。
func (config *Config) WithSlowQueryThresholdInMilli(milli int64) *Config {
	config.SlowQueryThresholdInMill = milli
	return config
}

// WithMode 设置 Hertz 运行模式，返回 Config 以便链式调用。
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// Build 根据当前配置创建 Server，并按开关挂载慢查询、指标与链路追踪中间件。
func (config *Config) Build() *Server {
	serverInstance := newServer(config)
	if config.EnableSlowQuery {
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

// Address 返回 Host:Port 形式的监听地址。
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
