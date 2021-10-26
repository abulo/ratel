package monitor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/abulo/ratel/v2/logger"
	"github.com/abulo/ratel/v2/server"
)

// Config ...
type Config struct {
	Host    string
	Port    int
	Network string
	Enable  bool
	Name    string
}

// Server ...
type Server struct {
	*http.Server
	listener net.Listener
	*Config
	serverInfo *server.ServiceInfo
}

var (
	// DefaultServeMux ...
	DefaultServeMux = http.NewServeMux()
)

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

// WithName ...
func (config *Config) WithName(name string) *Config {
	config.Name = name
	return config
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {

	var listener, err = net.Listen(config.Network, config.Address())
	if err != nil {
		logger.Logger.Panic(err)
	}
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(config.Address()),
		server.WithName(config.Name),
	)

	return &Server{
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener:   listener,
		Config:     config,
		serverInfo: &info,
	}
}

//Serve ..
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

//Stop ..
func (s *Server) Stop() error {
	return s.Server.Close()
}

//GracefulStop ..
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown(context.Background())
}

//Info ..
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.listener.Addr().String()),
		server.WithName(s.Name),
	)
	return &info
}

// HandleFunc ...
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	// todo: 增加安全管控
	DefaultServeMux.HandleFunc(pattern, handler)
}
