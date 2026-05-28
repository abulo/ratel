package manager

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	"github.com/hibiken/asynq"
)

// Server ...
type Server struct {
	Manager *asynq.PeriodicTaskManager
	config  *Config
}

func newServer(config *Config) *Server {
	m, err := asynq.NewPeriodicTaskManager(config.Option)
	if err != nil {
		panic(err)
	}
	return &Server{
		Manager: m,
		config:  config,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	return s.Manager.Run()
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	s.Manager.Shutdown()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Manager.Shutdown()
	return nil
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.config.Address()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("asynq"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}

// Health ...
func (s *Server) Health() bool {
	return true
}
