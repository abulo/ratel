package grpc

import (
	"context"
	"time"

	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/sirupsen/logrus"
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
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Panic("dial grpc server,", ecode.ErrKindRequestErr)

		} else {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("dial grpc server,", ecode.ErrKindRequestErr)
		}
	}
	logger.Logger.Info("start grpc client")
	return cc
}
