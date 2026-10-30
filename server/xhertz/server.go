package xhertz

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	hserver "github.com/cloudwego/hertz/pkg/app/server"
)

// Server 在 Hertz 实例之上封装 HTTP 服务端的生命周期管理与注册信息查询。
//
// 零值不可直接使用，须由 Config.Build 创建；Server 实现 server.Server 接口。
type Server struct {
	*hserver.Hertz
	config *Config
}

func newServer(config *Config) *Server {
	return &Server{
		Hertz: hserver.New(
			hserver.WithHostPorts(config.Address()),
			hserver.WithDisablePrintRoute(config.Mode == ReleaseMode),
			hserver.WithMaxRequestBodySize(20*1204*1204),
			hserver.WithMaxKeepBodySize(20*1204*1204),
			hserver.WithReadBufferSize(20*1204*1204),
		),
		config: config,
	}
}

// Serve 启动 Hertz 服务并阻塞运行，返回底层服务退出时的错误。
func (s *Server) Serve() error {
	return s.Run()
}

// Stop 立即关闭 Hertz 服务，不再等待在途请求完成。
func (s *Server) Stop() error {
	return s.Close()
}

// GracefulStop 在 ctx 取消前等待在途请求处理完成后关闭服务；ctx 超时后由底层服务强制退出。
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

// Info 返回服务注册信息，协议为 http；注册地址优先取配置中的 ServiceAddress，否则使用监听地址。
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

// Health 固定返回 true，Hertz 服务端不做独立端口探活。
func (s *Server) Health() bool {
	return true
}
