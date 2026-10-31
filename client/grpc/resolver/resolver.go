package resolver

import (
	"context"
	"strings"
	"sync"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// NewEtcdBuilder 创建一个基于指定注册器的 gRPC 解析器 Builder；name 为该解析器注册使用的
// Scheme，registry 负责服务节点的监听与变更推送。
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

// Build 启动对目标服务的监听：后台协程接收节点变更，将节点地址、服务信息以及路由、提供方、
// 消费方配置组装为 resolver.State 上报给 ClientConn；首次监听失败时返回错误。
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	reg := b.reg

	serviceName := target.Endpoint()
	if !strings.HasSuffix(serviceName, "/") {
		serviceName += "/"
	}

	endpoints, err := reg.WatchServices(context.Background(), serviceName)
	if err != nil {
		logger.Logger.Error("watch services failed:", err)
		return nil, err
	}
	var stop = make(chan struct{})
	goroutine.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				// 状态属性依次携带路由配置、服务提供方元信息与服务消费方配置信息。
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.
						New(constant.KeyRouteConfig, endpoint.RouteConfigs).
						WithValue(constant.KeyProviderConfig, endpoint.ProviderConfigs).
						WithValue(constant.KeyConsumerConfig, endpoint.ConsumerConfigs),
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = serviceName
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

// Scheme 返回该 Builder 注册使用的解析器 Scheme 名称。
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	stop chan struct{}
	once sync.Once
}

// ResolveNow 为空实现：节点变更由后台监听通道主动推送，不支持按需触发重新解析。
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {}

// Close 停止后台监听协程。
func (b *baseResolver) Close() {
	// PERF-5-12 [P2]：以 close(stop) 替代向无缓冲通道发送信号，重复 Close 不再永久阻塞；
	// sync.Once 保证关闭幂等，后台协程经 `case <-stop` 感知通道关闭后退出。
	b.once.Do(func() { close(b.stop) })
}
