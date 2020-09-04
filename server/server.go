package server

import (
	"context"
	"fmt"
)

type Option func(c *ServiceInfo)

// ServiceInfo represents service info
type ServiceInfo struct {
	Name    string `json:"name"`
	Scheme  string `json:"scheme"`
	Address string `json:"address"`
}

// Label ...
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

// Server ...
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func WithName(name string) Option {
	return func(c *ServiceInfo) {
		c.Name = name
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{}
	return si
}
