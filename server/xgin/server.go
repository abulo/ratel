package xgin

import (
	"context"
	"net"
	"net/http"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/server"
	"github.com/gin-gonic/gin"
)

// Server 在 gin.Engine 之上封装 HTTP 服务端的监听、生命周期管理与注册信息查询。
//
// 零值不可直接使用，须由 Config.Build 创建；Server 实现 server.Server 接口，
// 监听地址在创建时由系统端口分配后回填到配置中。
type Server struct {
	*gin.Engine
	// Server 为底层标准库 http.Server，由 Serve 方法初始化。
	Server   *http.Server
	config   *Config
	listener net.Listener
	// info 为构造期一次性缓存的注册信息，避免每次 Info() 重复重建 ServiceInfo。
	info *server.ServiceInfo
}

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		logger.Logger.Panic("new gin server err:", err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(config.Mode)
	// PERF-6-10 [P2] 注册信息在绑定监听地址后构造一次并缓存，Info() 直接返回。
	serviceAddr := listener.Addr().String()
	if config.ServiceAddress != "" {
		serviceAddr = config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	return &Server{
		Engine:   gin.New(),
		config:   config,
		listener: listener,
		info:     &info,
	}
}

// Upgrade 将 WebSocket 处理器注册为 GET 路由，收到请求后完成协议升级并交给 ws.Handler 处理。
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve 启动 HTTP 服务并阻塞监听；服务被正常关闭（http.ErrServerClosed）时返回 nil，其余错误原样返回。
func (s *Server) Serve() error {
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s,
	}
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		logger.Logger.Info("close gin:", s.config.Address())
		return nil
	}

	return err
}

// Stop 立即关闭底层 HTTP 服务，不再等待在途请求完成。
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop 在 ctx 取消前等待在途请求处理完成后关闭服务；ctx 超时后由底层服务强制退出。
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info 返回服务注册信息，协议为 http；注册地址优先取配置中的 ServiceAddress，否则使用实际监听地址。
//
// PERF-6-10 [P2] 直接返回构造期缓存的注册信息，不再重复重建。
func (s *Server) Info() *server.ServiceInfo {
	return s.info
}

// Health 判断服务监听 listener 是否就绪；不抢占共享 listener 上的真实连接。
//
// PERF-6-7 [P1] 原实现直接 Accept 共享 listener 并关闭连接，会与 http.Server 的接受循环竞争、
// 吞掉真实客户端请求；改为仅依据 listener 是否存在判断，保持签名与 bool 返回语义不变。
func (s *Server) Health() bool {
	return s.listener != nil
}
