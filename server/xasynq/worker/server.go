package worker

import (
	"context"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	"github.com/hibiken/asynq"
)

// Server 封装 Asynq 任务工作者，提供任务消费的生命周期管理与注册信息查询。
//
// 零值不可直接使用，须由 Config.Build 创建；Server 实现 server.Server 接口。
type Server struct {
	// Worker 为底层 asynq.Server，负责从 Redis 队列拉取并分发任务。
	Worker *asynq.Server
	config *Config
}

func newServer(config *Config) *Server {
	return &Server{
		Worker: asynq.NewServerFromRedisClient(config.Redis, config.Option),
		config: config,
	}
}

// Serve 启动任务工作者并阻塞运行，按 Mux 中注册的处理器消费队列任务，返回退出时的错误。
func (s *Server) Serve() error {
	return s.Worker.Run(s.config.Mux)
}

// Stop 立即停止任务工作者，不再等待在途任务完成。
func (s *Server) Stop() error {
	s.Worker.Stop()
	return nil
}

// GracefulStop 等待在途任务处理完成后优雅停止任务工作者。
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Worker.Shutdown()
	return nil
}

// Info 返回任务工作者的注册信息，协议为 asynq；注册地址优先取配置中的 ServiceAddress，
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

// Health 固定返回 true，任务工作者不做独立探活。
func (s *Server) Health() bool {
	return true
}
