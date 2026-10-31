package xmonitor

import (
	"context"
	"net"
	"net/http"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/server"
)

// Server 封装独立的监控 HTTP 服务端，对外暴露 pprof、路由清单与构建信息等治理接口。
//
// 零值不可直接使用，须由 Config.Build 创建；Server 实现 server.Server 接口，
// 路由统一注册到 DefaultServeMux。
type Server struct {
	*http.Server
	listener net.Listener
	*Config
	// info 为构造期一次性缓存的注册信息，避免每次 Info() 重复重建 ServiceInfo。
	info *server.ServiceInfo
}

func newServer(config *Config) *Server {
	var listener, err = net.Listen("tcp4", config.Address())
	if err != nil {
		logger.Logger.Panic("start error:", err)
	}

	// PERF-6-22 [P2] 注册信息在绑定监听地址后构造一次并缓存，Info() 直接返回，避免重复重建。
	serviceAddr := listener.Addr().String()
	if config.ServiceAddress != "" {
		serviceAddr = config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceMonitor),
	)

	return &Server{
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener: listener,
		Config:   config,
		info:     &info,
	}
}

// Serve 启动监控 HTTP 服务并阻塞监听；服务被正常关闭（http.ErrServerClosed）时返回 nil，其余错误原样返回。
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

// Stop 立即关闭监控 HTTP 服务，不再等待在途请求完成。
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop 在 ctx 取消前等待在途请求处理完成后关闭服务；ctx 超时后由底层服务强制退出。
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Health 当前固定返回 true。
//
// TODO(roamerlv): 待补充基于实际进程状态的健康检查逻辑。
func (s *Server) Health() bool {
	return true
}

// Info 返回监控服务的注册信息，服务类型为 ServiceMonitor；注册地址优先取配置中的
// ServiceAddress，否则使用实际监听地址。
//
// PERF-6-22 [P2] 直接返回构造期缓存的注册信息，不再重复重建。
func (s *Server) Info() *server.ServiceInfo {
	return s.info
}
