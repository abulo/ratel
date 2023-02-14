package client

import (
	"github.com/abulo/ratel/client/grpc"
	"github.com/abulo/ratel/registry/etcdv3"
)

// Grpc ...
type GrpcConfig struct {
	*grpc.Config
}

func NewGrpcConfig() *GrpcConfig {
	return &GrpcConfig{}
}

func (proxy *GrpcConfig) Store(client *grpc.Config) {
	proxy.Config = client
}

func (proxypool *Proxy) StoreGrpc(group string, proxy *GrpcConfig) {
	proxypool.m.Store(group, proxy)
}

func (proxypool *Proxy) LoadGrpc(group string) *grpc.Config {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*GrpcConfig).Config
	}
	return nil
}

type EtcdConfig struct {
	*etcdv3.Config
}

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

func (proxy *EtcdConfig) Store(client *etcdv3.Config) {
	proxy.Config = client
}

func (proxypool *Proxy) StoreEtcd(group string, proxy *EtcdConfig) {
	proxypool.m.Store(group, proxy)
}

func (proxypool *Proxy) LoadEtcd(group string) *etcdv3.Config {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*EtcdConfig).Config
	}
	return nil
}
