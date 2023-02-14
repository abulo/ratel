package xmonitor

import (
	"context"
	"net"
	"net/http"

	"github.com/abulo/ratel/core/constant"
	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/server"
)

// Server ...
type Server struct {
	*http.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) *Server {
	var listener, err = net.Listen("tcp4", config.Address())
	if err != nil {
		logger.Logger.Panic("start error:", err)
	}

	return &Server{
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener: listener,
		Config:   config,
	}
}

// Serve ..
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

// Stop ..
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop ..
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Health
// TODO(roamerlv):
func (s *Server) Health() bool {
	return true
}

// Info ..
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddr = s.Config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceMonitor),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
