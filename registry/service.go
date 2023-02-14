package registry

import "github.com/abulo/ratel/v2/server"

type service struct {
	server.Server
}

// Service ...
func Service(srv server.Server) server.Server {
	return &service{Server: srv}
}

// Serve ...
func (s *service) Serve() error {

	return s.Server.Serve()
}
