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
}

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		logger.Logger.Panic("new gin server err:", err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(config.Mode)
	return &Server{
		Engine:   gin.New(),
		config:   config,
		listener: listener,
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
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
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

// Health 通过尝试接受一条探测连接判断监听端口是否可用，连接可接受时返回 true。
func (s *Server) Health() bool {
	if s.listener == nil {
		return false
	}

	conn, err := s.listener.Accept()
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}
