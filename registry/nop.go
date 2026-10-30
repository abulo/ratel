package registry

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/server"
	"github.com/sirupsen/logrus"
)

// Local 是面向本地开发与调试的空实现注册器：注册与注销仅记录日志，
// 不访问任何外部注册中心；查询与监听能力尚未实现。
type Local struct{}

// ListServices 尚未实现，调用时直接 panic("implement me")。
func (n Local) ListServices(ctx context.Context, s string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices 尚未实现，调用时直接 panic("implement me")。
func (n Local) WatchServices(ctx context.Context, s string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService 仅记录本地注册日志，不向外部注册中心写入，始终返回 nil。
func (n Local) RegisterService(ctx context.Context, si *server.ServiceInfo) error {

	logger.Logger.WithFields(logrus.Fields{
		"action": "registry",
		"name":   si.Name,
		"label":  si.Label(),
	}).Info("register service locally")

	return nil
}

// UnregisterService 仅记录本地注销日志，不向外部注册中心删除数据，始终返回 nil。
func (n Local) UnregisterService(ctx context.Context, si *server.ServiceInfo) error {
	logger.Logger.WithFields(logrus.Fields{
		"action": "registry",
		"name":   si.Name,
		"label":  si.Label(),
	}).Info("unregister service locally")
	return nil
}

// Close 不持有任何外部资源，始终返回 nil。
func (n Local) Close() error { return nil }

// Kind 返回本地实现的类型标识 "local"。
func (n Local) Kind() string { return "local" }
