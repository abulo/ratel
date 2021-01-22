package worker

import (
	"context"

	"github.com/abulo/ratel/server"
	"github.com/abulo/ratel/util"
)

// Config ...
type Config struct {
	Worker *Worker
	Name   string
	Listen []string
}

// Server ...
type Server struct {
	*Config
	serverInfo *server.ServiceInfo
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	info := server.ApplyOptions(
		server.WithScheme("worker"),
		server.WithAddress(util.Implode(";", config.Listen)),
		server.WithName(config.Name),
	)
	return &Server{
		Config:     config,
		serverInfo: &info,
	}
}

//Info ..
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("worker"),
		server.WithAddress(util.Implode(";", s.Listen)),
		server.WithName(s.Name),
	)
	return &info
}

//Serve ..
func (s *Server) Serve() error {
	err := s.Worker.WorkerReady()
	if err != nil {
		s.Worker.WorkerClose()
	} else {
		s.Worker.WorkerDo()
	}
	return err
}

//Stop ..
func (s *Server) Stop() error {
	return s.Worker.WorkerClose()
}

//GracefulStop ..
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Stop()
}
