package consul

import "github.com/abulo/ratel/registry"

func NewRegistry(opts ...registry.Option) registry.Registry {
	return registry.NewRegistry(opts...)
}
