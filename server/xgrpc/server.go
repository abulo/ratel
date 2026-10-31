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
	// info 为构造期一次性缓存的注册信息，避免每次 Info() 重复重建 ServiceInfo。
	info *server.ServiceInfo
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
		MaxConnectionIdle: 15 * time.Second,
		// PERF-6-15 [P2] 原 30 秒 MaxConnectionAge 会强制所有长连接频繁重建（含 TLS 握手），
		// 放宽到 10 分钟量级，配合在途 RPC 宽限期自然老化连接（预期收益需实测验证）。
		MaxConnectionAge: 10 * time.Minute,
		// 给在途 RPC 5 秒宽限期后再强制关闭连接。
		MaxConnectionAgeGrace: 5 * time.Second,
		// 连接空闲 5 秒后向客户端发送 Ping 确认存活。
		Time: 5 * time.Second,
		// 等待 Ping 应答 1 秒，超时即判定连接已失效。
		Timeout: 1 * time.Second,
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
	// PERF-6-3 [P2] 注册信息在绑定监听地址后构造一次并缓存，Info() 直接返回，避免重复重建。
	serviceAddress := listener.Addr().String()
	if config.ServiceAddress != "" {
		serviceAddress = config.ServiceAddress
	}
	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(serviceAddress),
		server.WithKind(constant.ServiceProvider),
	)
	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
		info:     &info,
	}, nil
}

// Health 判断服务监听 listener 是否就绪；不抢占共享 listener 上的真实连接。
//
// PERF-6-14 [P1] 原实现直接 Accept 共享 listener 并关闭连接，会与 grpc.Server 的接受循环竞争、
// 吞掉真实客户端连接；改为仅依据 listener 是否存在判断，保持签名与 bool 返回语义不变。
func (s *Server) Health() bool {
	return s.listener != nil
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
//
// PERF-6-3 [P2] 直接返回构造期缓存的注册信息，不再重复重建。
func (s *Server) Info() *server.ServiceInfo {
	return s.info
}
