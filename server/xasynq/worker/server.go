package worker

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	"github.com/hibiken/asynq"
)

// Server ...
type Server struct {
	Worker *asynq.Server
	config *Config
}

func newServer(config *Config) *Server {
	return &Server{
		Worker: asynq.NewServerFromRedisClient(config.Redis, config.Option),
		config: config,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	return s.Worker.Run(s.config.Mux)
	// return
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	s.Worker.Stop()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Worker.Shutdown()
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
