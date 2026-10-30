package manager

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	"github.com/hibiken/asynq"
)

// Server 封装 Asynq 定时任务管理器，负责按调度计划周期性入队任务，并提供生命周期管理与注册信息查询。
//
// 零值不可直接使用，须由 Config.Build 创建；Server 实现 server.Server 接口。
type Server struct {
	// Manager 为底层 asynq.PeriodicTaskManager，负责定时任务的调度与入队。
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

// Serve 启动定时任务管理器并阻塞运行，按已注册的定时计划周期性创建任务。
func (s *Server) Serve() error {
	return s.Manager.Run()
}

// Stop 立即关闭定时任务管理器，不再等待在途调度完成。
func (s *Server) Stop() error {
	s.Manager.Shutdown()
	return nil
}

// GracefulStop 优雅关闭定时任务管理器，等待当前调度与在途任务收尾后退出。
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Manager.Shutdown()
	return nil
}

// Info 返回定时任务管理器的注册信息，协议为 asynq；注册地址优先取配置中的 ServiceAddress，
// 否则使用配置的 Host:Port。
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

// Health 固定返回 true，定时任务管理器不做独立探活。
func (s *Server) Health() bool {
	return true
}
