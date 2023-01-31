package resolver

import (
	"context"
	"strings"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Register ...
// func Register(name string, reg registry.Registry) {
// 	resolver.Register(&baseBuilder{
// 		name: name,
// 		reg:  reg,
// 	})
// }

// NewEtcdBuilder returns a new etcdv3 resolver builder.
func NewEtcdBuilder(name string, registry registry.Registry) resolver.Builder {
	return &baseBuilder{
		name: name,
		reg:  registry,
	}
}

type baseBuilder struct {
	name string
	reg  registry.Registry
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	reg := b.reg
	if !strings.HasSuffix(target.Endpoint, "/") {
		target.Endpoint += "/"
	}
	endpoints, err := reg.WatchServices(context.Background(), target.Endpoint)
	if err != nil {
		logger.Logger.Error("watch services failed:", err)
		return nil, err
	}
	var stop = make(chan struct{})
	goroutine.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.
						New(constant.KeyRouteConfig, endpoint.RouteConfigs).             // 路由配置
						WithValue(constant.KeyProviderConfig, endpoint.ProviderConfigs). // 服务提供方元信息
						WithValue(constant.KeyConsumerConfig, endpoint.ConsumerConfigs), // 服务消费方配置信息,
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = target.Endpoint
					address.Attributes = attributes.New(constant.KeyServiceInfo, node)
					state.Addresses = append(state.Addresses, address)
				}
				_ = cc.UpdateState(state)
			case <-stop:
				return
			}
		}
	})

	return &baseResolver{
		stop: stop,
	}, nil
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	stop chan struct{}
}

// ResolveNow ...
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {}

// Close ...
func (b *baseResolver) Close() { b.stop <- struct{}{} }
