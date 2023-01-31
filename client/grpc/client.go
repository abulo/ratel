package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/client/grpc/resolver"
	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCClient(config *Config) (*grpc.ClientConn, error) {
	var ctx = context.Background()
	dialOptions := getDialOptions(config)
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}
		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	conn, err := grpc.DialContext(ctx, config.Address, dialOptions...)
	// conn, err := grpc.DialContext(ctx, config.Address, append(dialOptions, grpc.WithBlock())...)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("dial grpc server failed, connect without block,", ecode.ErrKindRequestErr)
		conn, err = grpc.DialContext(context.Background(), config.Address, dialOptions...)

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

func getDialOptions(config *Config) []grpc.DialOption {
	dialOptions := config.dialOptions

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions,
		grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolver.NewEtcdBuilder(config.Etcd.GetNode(), config.Etcd.MustSingleton())),
		grpc.WithDisableServiceConfig(),
	)

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.BalancerName)
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(svcCfg))

	return dialOptions
}
