package manager

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// ModName 为 Asynq 定时任务管理器在框架内的模块标识，用于注册与监控归类。
const ModName = "server.asynq.manager"

// Config 保存 Asynq 定时任务管理器的注册信息与 PeriodicTaskManager 运行选项配置。
//
// 零值不可直接使用，须通过 New 结合各 With 方法构造；Config 为配置载体，
// 应在 Build 之前一次性设置完成，不保证并发安全。
type Config struct {
	// Host 为注册信息中登记的监听主机地址。
	Host                     string
	// Port 为注册信息中登记的监听端口。
	Port                     int
	// Mode 为应用运行模式。
	Mode                     string
	// Deployment 为部署组，用于服务注册时的流量隔离。
	Deployment               string
	// EnableMetric 为 true 时启用指标采集（预留开关）。
	EnableMetric             bool
	// EnableTrace 为 true 时启用链路追踪（预留开关）。
	EnableTrace              bool
	// EnableSlowQuery 为 true 时启用慢任务记录（预留开关）。
	EnableSlowQuery          bool
	// ServiceAddress 为注册中心中登记的服务地址；为空时使用 Host:Port。
	ServiceAddress           string
	// SlowQueryThresholdInMill 为慢任务阈值（单位毫秒），当前作为预留配置项未被消费，默认 500。
	SlowQueryThresholdInMill int64
	// Option 为 asynq.PeriodicTaskManager 的运行选项，包含 Redis 连接、任务规格与日志等配置。
	Option                   asynq.PeriodicTaskManagerOpts
}

// New 返回带有默认慢任务阈值（500 毫秒）的 Config；Option 须在 Build 前补齐。
func New() *Config {
	return &Config{
		SlowQueryThresholdInMill: 500,
	}
}

// WithHost 设置注册信息中的监听主机地址，返回 Config 以便链式调用。
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort 设置注册信息中的监听端口，返回 Config 以便链式调用。
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// WithDeployment 设置部署组，返回 Config 以便链式调用。
func (config *Config) WithDeployment(deployment string) *Config {
	config.Deployment = deployment
	return config
}

// WithEnableSlowQuery 设置是否启用慢任务记录（预留开关），返回 Config 以便链式调用。
func (config *Config) WithEnableSlowQuery(disableSlowQuery bool) *Config {
	config.EnableSlowQuery = disableSlowQuery
	return config
}

// WithEnableMetric 设置是否启用指标采集（预留开关），返回 Config 以便链式调用。
func (config *Config) WithEnableMetric(disableMetric bool) *Config {
	config.EnableMetric = disableMetric
	return config
}

// WithEnableTrace 设置是否启用链路追踪（预留开关），返回 Config 以便链式调用。
func (config *Config) WithEnableTrace(disableTrace bool) *Config {
	config.EnableTrace = disableTrace
	return config
}

// WithServiceAddress 设置注册中心中登记的服务地址，返回 Config 以便链式调用。
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// WithSlowQueryThresholdInMilli 设置慢任务阈值（预留配置项，单位毫秒），返回 Config 以便链式调用。
func (config *Config) WithSlowQueryThresholdInMilli(milli int64) *Config {
	config.SlowQueryThresholdInMill = milli
	return config
}

// WithMode 设置应用运行模式，返回 Config 以便链式调用。
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// Build 根据当前配置创建定时任务管理器 Server；Option 须已配置完成，构造失败时直接 panic。
func (config *Config) Build() *Server {
	serverInstance := newServer(config)
	return serverInstance
}

// Address 返回 Host:Port 形式的服务注册地址。
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
