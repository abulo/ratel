package xgrpc

import (
	"fmt"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Deployment string `json:"deployment"`
	// Network network type, tcp4 by default
	Network string `json:"network" toml:"network"`
	// EnableTrace  Trace Interceptor, false by default
	EnableTrace bool
	// EnableMetric disable Metric Interceptor, false by default
	EnableMetric bool
	// SlowQueryThresholdInMill, request will be colored if cost over this threshold value
	SlowQueryThresholdInMill int64
	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string
	// EnableTLS
	EnableTLS bool
	// CaFile
	CaFile string
	// CertFile
	CertFile string
	// PrivateFile
	PrivateFile        string
	Labels             map[string]string `json:"labels"`
	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// New ...
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

// WithServerOption inject server option to grpc server User should not inject interceptor option, which is recommend by WithStreamInterceptor  and WithUnaryInterceptor
func (config *Config) WithServerOption(options ...grpc.ServerOption) *Config {
	if config.serverOptions == nil {
		config.serverOptions = make([]grpc.ServerOption, 0)
	}
	config.serverOptions = append(config.serverOptions, options...)
	return config
}

// WithStreamInterceptor inject stream interceptors to server option
func (config *Config) WithStreamInterceptor(intes ...grpc.StreamServerInterceptor) *Config {
	if config.streamInterceptors == nil {
		config.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
	}
	config.streamInterceptors = append(config.streamInterceptors, intes...)
	return config
}

// WithUnaryInterceptor inject unary interceptors to server option
func (config *Config) WithUnaryInterceptor(intes ...grpc.UnaryServerInterceptor) *Config {
	if config.unaryInterceptors == nil {
		config.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
	}

	config.unaryInterceptors = append(config.unaryInterceptors, intes...)
	return config
}

// MustBuild ...
func (config *Config) MustBuild() *Server {
	server, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("build grpc server: %v", err)
	}
	return server
}

// Build ...
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

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
