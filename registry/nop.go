package registry

import (
	"context"

	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/server"
	"github.com/sirupsen/logrus"
)

// Local Nop registry, used for local development/debugging
type Local struct{}

// ListServices ...
func (n Local) ListServices(ctx context.Context, s string, s2 string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Local) WatchServices(ctx context.Context, s string, s2 string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Local) RegisterService(ctx context.Context, si *server.ServiceInfo) error {

	logger.Logger.WithFields(logrus.Fields{
		"action": "registry",
		"name":   si.Name,
		"label":  si.Label(),
	}).Info("register service locally")

	return nil
}

// UnregisterService ...
func (n Local) UnregisterService(ctx context.Context, si *server.ServiceInfo) error {
	logger.Logger.WithFields(logrus.Fields{
		"action": "registry",
		"name":   si.Name,
		"label":  si.Label(),
	}).Info("unregister service locally")
	return nil
}

// Close ...
func (n Local) Close() error { return nil }

// Kind Close ...
func (n Local) Kind() string { return "local" }
