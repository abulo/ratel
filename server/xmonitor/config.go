package xmonitor

import (
	"fmt"
)

// ModName 为监控服务端在框架内的模块标识，用于注册与监控归类。
const ModName = "server.monitor"

// Config 保存监控 HTTP 服务端的监听地址与注册信息配置。
//
// 零值不可直接使用，须通过 New 结合各 With 方法构造；Config 为配置载体，
// 应在 Build 之前一次性设置完成，不保证并发安全。
type Config struct {
	// Host 为监听主机地址。
	Host           string
	// Port 为监听端口。
	Port           int
	// Deployment 为部署组，用于服务注册时的流量隔离。
	Deployment     string
	// Mode 为应用运行模式。
	Mode           string
	// ServiceAddress 为注册中心中登记的服务地址；为空时使用 Host:Port。
	ServiceAddress string
}

// New 返回空的 Config，各字段须通过 With 方法设置。
func New() *Config {
	return &Config{}
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

// WithMode 设置应用运行模式，返回 Config 以便链式调用。
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// WithServiceAddress 设置注册中心中登记的服务地址，返回 Config 以便链式调用。
func (config *Config) WithServiceAddress(serviceAddress string) *Config {
	config.ServiceAddress = serviceAddress
	return config
}

// Build 根据当前配置创建监控 Server。
func (config *Config) Build() *Server {
	return newServer(config)
}

// Address 返回 Host:Port 形式的监听地址。
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
