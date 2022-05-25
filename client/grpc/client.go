package grpc

import (
	"context"
	"time"

	"github.com/abulo/ratel/v3/ecode"
	"github.com/abulo/ratel/v3/logger"
	"google.golang.org/grpc"
)

func newGRPCClient(config *Config) *grpc.ClientConn {
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}
	grpcServiceConfig := `{"loadBalancingPolicy":"` + config.BalancerName + `"}`
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(grpcServiceConfig))
	cc, err := grpc.DialContext(ctx, config.Address, dialOptions...)
	if err != nil {
		if config.OnDialError == "panic" {
			logger.Logger.Panic("dial grpc server", ecode.ErrKindRequestErr, err)
		} else {
			logger.Logger.Error("dial grpc server", ecode.ErrKindRequestErr, err)
		}
	}
	logger.Logger.Info("start grpc client")
	return cc
}
