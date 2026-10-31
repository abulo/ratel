package grpc

import (
	"fmt"

	"github.com/abulo/ratel/v3/client/grpc/resolver"
	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// newGRPCClient 组装拨号选项并创建 gRPC ClientConn；首次拨号失败时按 panic 级别记录日志，
// 随后重试一次，再次失败返回错误。
func newGRPCClient(config *Config) (*grpc.ClientConn, error) {
	dialOptions := getDialOptions(config)

	conn, err := grpc.NewClient(config.Address, dialOptions...)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("dial grpc server failed, connect without block,", ecode.ErrKindRequestErr)
		conn, err = grpc.NewClient(config.Address, dialOptions...)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Panic("connect without block failed", ecode.ErrKindRequestErr)
			return nil, err
		}
	}
	logger.Logger.Info("start grpc client")
	return conn, nil
}

// getDialOptions 汇总配置携带的拨号选项，追加 keepalive 参数、insecure 凭证、etcd 解析器与
// 默认服务配置（负载均衡策略）。调用方须保证 config.Etcd 非 nil。
func getDialOptions(config *Config) []grpc.DialOption {
	dialOptions := config.dialOptions

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	// PERF-5-5 [P2]：insecure 凭证在本单点装配，New() 预置的重复凭证已移除，避免同一连接上出现
	// 两条重复的 WithTransportCredentials 选项；拨号连接语义保持不变。
	dialOptions = append(dialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolver.NewEtcdBuilder(config.Etcd.GetNode(), config.Etcd.MustSingleton())),
		grpc.WithDisableServiceConfig(),
	)

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.BalancerName)
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(svcCfg))

	return dialOptions
}
