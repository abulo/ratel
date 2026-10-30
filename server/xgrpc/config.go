package xgrpc

import (
	"fmt"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"google.golang.org/grpc"
)

// Config 保存 gRPC 服务端的监听地址、网络类型、TLS、拦截器开关等配置。
//
// 零值不可直接使用，须通过 New 结合各 With 方法构造；Config 为配置载体，
// 应在 Build 之前一次性设置完成，不保证并发安全。
type Config struct {
	// Name 为服务名称。
	Name       string `json:"name"`
	// Host 为监听主机地址。
	Host       string `json:"host"`
	// Port 为监听端口；Build 时由实际绑定端口回填。
	Port       int    `json:"port"`
	// Deployment 为部署组，用于服务注册时的流量隔离。
	Deployment string `json:"deployment"`
	// Network 为监听网络类型，默认 tcp4。
	Network string `json:"network" toml:"network"`
	// EnableTrace 为 true 时追加服务端链路追踪拦截器，默认 false。
	EnableTrace bool
	// EnableMetric 为 true 时追加服务指标采集拦截器，默认 false。
	EnableMetric bool
	// SlowQueryThresholdInMill 为慢请求阈值（单位毫秒），当前作为预留配置项，未参与拦截器判定，默认 500。
	SlowQueryThresholdInMill int64
	// ServiceAddress 为注册中心中登记的服务地址；为空时使用 Host:Port。
	ServiceAddress string
	// EnableTLS 为 true 时启用双向 TLS 认证，须同时提供证书、私钥与 CA 文件。
	EnableTLS bool
	// CaFile 为受信任 CA 证书文件路径，仅 EnableTLS 为 true 时使用。
	CaFile string
	// CertFile 为服务端证书文件路径，仅 EnableTLS 为 true 时使用。
	CertFile string
	// PrivateFile 为服务端私钥文件路径，仅 EnableTLS 为 true 时使用。
	PrivateFile        string
	// Labels 为随服务一起注册的标签。
	Labels             map[string]string `json:"labels"`
	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// New 返回带有默认值的 Config：网络为 tcp4，追踪、指标、TLS 默认关闭，慢查询阈值为 500 毫秒。
func New() *Config {
	return &Config{
		Network:                  "tcp4",
		Deployment:               constant.DefaultDeployment,
		EnableMetric:             false,
		EnableTrace:              false,
		EnableTLS:                false,
		SlowQueryThresholdInMill: 500,
		serverOptions:            []grpc.ServerOption{},
		streamInterceptors:       []grpc.StreamServerInterceptor{},
		unaryInterceptors:        []grpc.UnaryServerInterceptor{},
	}
}

// WithServerOption 向 gRPC 服务端追加 ServerOption；拦截器应分别通过 WithStreamInterceptor、
// WithUnaryInterceptor 注入，不要在此处传入拦截器选项。
func (config *Config) WithServerOption(options ...grpc.ServerOption) *Config {
	if config.serverOptions == nil {
		config.serverOptions = make([]grpc.ServerOption, 0)
	}
	config.serverOptions = append(config.serverOptions, options...)
	return config
}

// WithStreamInterceptor 向服务端追加流式 RPC 拦截器，返回 Config 以便链式调用。
func (config *Config) WithStreamInterceptor(intes ...grpc.StreamServerInterceptor) *Config {
	if config.streamInterceptors == nil {
		config.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
	}
	config.streamInterceptors = append(config.streamInterceptors, intes...)
	return config
}

// WithUnaryInterceptor 向服务端追加一元 RPC 拦截器，返回 Config 以便链式调用。
func (config *Config) WithUnaryInterceptor(intes ...grpc.UnaryServerInterceptor) *Config {
	if config.unaryInterceptors == nil {
		config.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
	}

	config.unaryInterceptors = append(config.unaryInterceptors, intes...)
	return config
}

// MustBuild 调用 Build 创建 Server；构建失败时直接 panic，适用于启动阶段配置已确定的场景。
func (config *Config) MustBuild() *Server {
	server, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("build grpc server: %v", err)
	}
	return server
}

// Build 根据当前配置组装拦截器链并创建 Server；TLS 证书读取、监听绑定失败时返回错误。
func (config *Config) Build() (*Server, error) {
	if config.EnableTrace {
		config.unaryInterceptors = append(config.unaryInterceptors, NewTraceUnaryServerInterceptor())
		config.streamInterceptors = append(config.streamInterceptors, NewTraceStreamServerInterceptor())
	}

	if config.EnableMetric {
		config.unaryInterceptors = append(config.unaryInterceptors, prometheusUnaryServerInterceptor)
		config.streamInterceptors = append(config.streamInterceptors, prometheusStreamServerInterceptor)
	}

	return newServer(config)
}

// Address 返回 Host:Port 形式的监听地址。
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
