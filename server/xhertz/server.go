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
	// info 为构造期一次性缓存的注册信息，避免每次 Info() 重复重建 ServiceInfo。
	info *server.ServiceInfo
}

func newServer(config *Config) *Server {
	// PERF-6-16 [P1] 修正 1204 系 1024 笔误：原三个上限均为 28,992,320 字节（≈27.65MiB），
	// 连接级缓冲随并发放大内存占用。取值依据：
	//   - MaxRequestBodySize：单次请求体上限，按业务常见 JSON/表单报文取 10MiB 量级；
	//   - MaxKeepBodySize：读取后保留给 handler 访问的内存体，无需与整体 body 上限同量级，取 1MiB；
	//   - ReadBufferSize：连接级读取缓冲（头部 + 小包），回落到 32KiB 常规量级。
	// （各缓冲在 Hertz 内部分配时机与实际内存峰值需压测确认，预期收益需实测验证）
	serviceAddr := config.Address()
	if config.ServiceAddress != "" {
		serviceAddr = config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	return &Server{
		Hertz: hserver.New(
			hserver.WithHostPorts(config.Address()),
			hserver.WithDisablePrintRoute(config.Mode == ReleaseMode),
			hserver.WithMaxRequestBodySize(10*1024*1024),
			hserver.WithMaxKeepBodySize(1*1024*1024),
			hserver.WithReadBufferSize(32*1024),
		),
		config: config,
		info:   &info,
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
//
// PERF-6-20 [P2] 直接返回构造期缓存的注册信息，不再重复重建。
func (s *Server) Info() *server.ServiceInfo {
	return s.info
}

// Health 固定返回 true，Hertz 服务端不做独立端口探活。
func (s *Server) Health() bool {
	return true
}
