package rpc

import (
	"context"
	"fmt"
	"net"

	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/server"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Host    string
	Port    int
	Network string // Network network type, tcp4 by default

	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor

	Name string
}

// Server ...
type Server struct {
	*grpc.Server
	listener   net.Listener
	config     *Config
	serverInfo *server.ServiceInfo
}

// Build ...
func (config *Config) Build() *Server {

	newServer := grpc.NewServer(config.serverOptions...)

	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		logger.Logger.Info("new grpc server err", config.Address(), err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(config.Address()),
		server.WithName(config.Name),
	)

	return &Server{
		Server:     newServer,
		listener:   listener,
		config:     config,
		serverInfo: &info,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	return err
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(s.config.Address()),
		server.WithName(s.config.Name),
	)
	return &info
}

// WithServerOption inject server option to grpc server
// User should not inject interceptor option, which is recommend by WithStreamInterceptor
// and WithUnaryInterceptor
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

func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
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

// WithNetwork ...
func (config *Config) WithNetwork(network string) *Config {
	config.Network = network
	return config
}

// WithName ...
func (config *Config) WithName(name string) *Config {
	config.Name = name
	return config
}
