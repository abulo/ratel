package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/server"
)

// Config HTTP config
type Config struct {
	Host   string
	Port   int
	Mode   string
	Name   string
	Health string
}

// Server ...
type Server struct {
	*gin.Engine
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

// WithMode ...
func (config *Config) WithMode(mode string) *Config {
	config.Mode = mode
	return config
}

// WithName ...
func (config *Config) WithName(name string) *Config {
	config.Name = name
	return config
}

// WithHealth ...
func (config *Config) WithHealth(health string) *Config {
	config.Health = health
	return config
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {

	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		logger.Panic("new gin server err", err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	gin.SetMode(config.Mode)

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(config.Address()),
		server.WithName(config.Name),
		server.WithHealth(config.Health),
	)

	return &Server{
		Engine:     gin.New(),
		config:     config,
		listener:   listener,
		serverInfo: &info,
	}
}

//ServerInterceptor ...
func (s *Server) ServerInterceptor(fn gin.HandlerFunc) *Server {
	s.Use(fn)
	return s
}

//Upgrade protocol to WebSocket
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.GET(ws.Pattern, ws.Name, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		logger.Info("close gin", s.config.Address())
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

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.config.Address()),
		server.WithName(s.config.Name),
		server.WithHealth(s.config.Health),
	)
	return &info
}
