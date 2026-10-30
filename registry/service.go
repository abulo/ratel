package registry

import "github.com/abulo/ratel/v3/server"

type service struct {
	server.Server
}

// Service 用给定的服务实例包装出可接入注册中心的服务对象，内部委托原服务实例完成启动。
func Service(srv server.Server) server.Server {
	return &service{Server: srv}
}

// Serve 委托内部服务实例执行监听与服务启动，并将其错误原样返回。
func (s *service) Serve() error {

	return s.Server.Serve()
}
