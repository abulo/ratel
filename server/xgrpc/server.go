package xgrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"time"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server 在 grpc.Server 之上封装 gRPC 服务端的监听、生命周期管理与注册信息查询。
//
// 零值不可直接使用，须由 Config.Build 或 MustBuild 创建；Server 实现 server.Server 接口，
// 监听地址在创建时由系统端口分配后回填到配置中。
type Server struct {
	*grpc.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) (*Server, error) {
	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{},
		config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{},
		config.unaryInterceptors...,
	)

	if config.EnableTLS {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.PrivateFile)
		if err != nil {
			return nil, errors.Wrap(err, "tls.LoadX509KeyPair failed")
		}

		certPool := x509.NewCertPool()
		rootBuf, err := os.ReadFile(config.CaFile)
		if err != nil {
			return nil, errors.Wrap(err, "os.ReadFile failed")
		}
		if !certPool.AppendCertsFromPEM(rootBuf) {
			return nil, errors.New("certPool.AppendCertsFromPEM failed")
		}

		tlsConf := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		}

		config.serverOptions = append(config.serverOptions,
			grpc.Creds(credentials.NewTLS(tlsConf)),
		)
	}

	config.serverOptions = append(config.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	kasp := keepalive.ServerParameters{
		// 客户端空闲 15 秒后发送 GOAWAY。
		MaxConnectionIdle:     15 * time.Second,
		// 连接存活超过 30 秒后发送 GOAWAY。
		MaxConnectionAge:      30 * time.Second,
		// 给在途 RPC 5 秒宽限期后再强制关闭连接。
		MaxConnectionAgeGrace: 5 * time.Second,
		// 连接空闲 5 秒后向客户端发送 Ping 确认存活。
		Time:                  5 * time.Second,
		// 等待 Ping 应答 1 秒，超时即判定连接已失效。
		Timeout:               1 * time.Second,
	}

	kaep := keepalive.EnforcementPolicy{
		// 客户端 Ping 间隔短于 5 秒时断开连接。
		MinTime:             5 * time.Second,
		// 允许没有活跃流时发送 Ping。
		PermitWithoutStream: true,
	}
	config.serverOptions = append(config.serverOptions,
		grpc.KeepaliveParams(kasp),
		grpc.KeepaliveEnforcementPolicy(kaep),
	)

	newServer := grpc.NewServer(config.serverOptions...)
	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		return nil, errors.Wrap(err, "net.Listen failed")
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	reflection.Register(newServer)
	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
	}, nil
}

// Health 通过尝试接受一条探测连接判断监听端口是否可用，连接可接受时返回 true。
func (s *Server) Health() bool {
	conn, err := s.listener.Accept()
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}

// Serve 启动 gRPC 服务并阻塞监听，返回底层服务端退出时的错误。
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	return err
}

// Stop 立即停止 gRPC 服务，不再等待在途 RPC 完成。
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop 等待在途 RPC 处理完成后优雅停止 gRPC 服务。
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

// Info 返回服务注册信息，协议为 grpc；注册地址优先取配置中的 ServiceAddress，否则使用实际监听地址。
func (s *Server) Info() *server.ServiceInfo {
	serviceAddress := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddress = s.Config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(serviceAddress),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}
