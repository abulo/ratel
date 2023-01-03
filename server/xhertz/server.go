package xhertz

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
)

// Server ...
type Server struct {
	*hserver.Hertz
	config *Config
}

func newServer(config *Config) *Server {
	return &Server{
		Hertz: hserver.New(
			hserver.WithHostPorts(config.Address()),
		),
		config: config,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	s.Spin()
	return nil
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.config.Address()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}

// Health ...
func (s *Server) Health() bool {
	return true
}
