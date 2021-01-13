package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/server"
)

// Config HTTP config
type Config struct {
	Host    string
	Port    int
	Name    string
	Network string
	Proxy   *Proxy
}

// Server ...
type Server struct {
	Server     *http.Server
	config     *Config
	listener   net.Listener
	serverInfo *server.ServiceInfo
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

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {

	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		logger.Logger.Panic("new proxy server err", err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	info := server.ApplyOptions(
		server.WithScheme("proxy"),
		server.WithAddress(config.Address()),
		server.WithName(config.Name),
	)

	return &Server{
		config:     config,
		listener:   listener,
		serverInfo: &info,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s.config.Proxy,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		logger.Logger.Info("close proxy", s.config.Address())
		return nil
	}
	return err
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// WithName ...
func (config *Config) WithName(name string) *Config {
	config.Name = name
	return config
}

// WithProxy ...
func (config *Config) WithProxy(proxy *Proxy) *Config {
	config.Proxy = proxy
	return config
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("proxy"),
		server.WithAddress(s.config.Address()),
		server.WithName(s.config.Name),
	)
	return &info
}
