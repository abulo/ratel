package hook

import (
	"context"

	"github.com/abulo/ratel/registry"
	"github.com/abulo/ratel/server"
	"golang.org/x/sync/errgroup"
)

type hookRegistry struct {
	registries []registry.Registry
}

// ListServices ...
func (c hookRegistry) ListServices(ctx context.Context, name string, scheme string) ([]*server.ServiceInfo, error) {
	var eg errgroup.Group
	var services = make([]*server.ServiceInfo, 0)
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			infos, err := registry.ListServices(ctx, name, scheme)
			if err != nil {
				return err
			}
			services = append(services, infos...)
			return nil
		})
	}
	err := eg.Wait()
	return services, err
}

// WatchServices ...
func (c hookRegistry) WatchServices(ctx context.Context, s string, s2 string) (chan registry.Endpoints, error) {
	panic("compound registry doesn't support watch services")
}

// RegisterService ...
func (c hookRegistry) RegisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.RegisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// UnregisterService ...
func (c hookRegistry) UnregisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.UnregisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// Close ...
func (c hookRegistry) Close() error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.Close()
		})
	}
	return eg.Wait()
}

// New ...
func New(registries ...registry.Registry) registry.Registry {
	return hookRegistry{
		registries: registries,
	}
}
