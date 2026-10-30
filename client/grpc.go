package client

import (
	"github.com/abulo/ratel/v3/client/grpc"
	"github.com/abulo/ratel/v3/registry/etcdv3"
)

// GrpcConfig 包装 grpc.Config，供 Proxy 按分组保存与加载 gRPC 客户端配置。
type GrpcConfig struct {
	*grpc.Config
}

// NewGrpcConfig 创建一个空的 GrpcConfig，内嵌 Config 为 nil，须调用 Store 注入实际配置。
func NewGrpcConfig() *GrpcConfig {
	return &GrpcConfig{}
}

// Store 将给定的 gRPC 配置存入当前包装对象，覆盖原有配置。
func (proxy *GrpcConfig) Store(client *grpc.Config) {
	proxy.Config = client
}

// StoreGrpc 将指定分组的 gRPC 配置存入代理池，覆盖同分组已有的配置。
func (proxypool *Proxy) StoreGrpc(group string, proxy *GrpcConfig) {
	proxypool.m.Store(group, proxy)
}

// LoadGrpc 按分组取出 gRPC 配置；分组不存在时返回 nil。
func (proxypool *Proxy) LoadGrpc(group string) *grpc.Config {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*GrpcConfig).Config
	}
	return nil
}

// EtcdConfig 包装 etcdv3.Config，供 Proxy 按分组保存与加载 etcd 客户端配置。
type EtcdConfig struct {
	*etcdv3.Config
}

// NewEtcdConfig 创建一个空的 EtcdConfig，内嵌 Config 为 nil，须调用 Store 注入实际配置。
func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

// Store 将给定的 etcd 配置存入当前包装对象，覆盖原有配置。
func (proxy *EtcdConfig) Store(client *etcdv3.Config) {
	proxy.Config = client
}

// StoreEtcd 将指定分组的 etcd 配置存入代理池，覆盖同分组已有的配置。
func (proxypool *Proxy) StoreEtcd(group string, proxy *EtcdConfig) {
	proxypool.m.Store(group, proxy)
}

// LoadEtcd 按分组取出 etcd 配置；分组不存在时返回 nil。
func (proxypool *Proxy) LoadEtcd(group string) *etcdv3.Config {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*EtcdConfig).Config
	}
	return nil
}
