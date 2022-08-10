package proxy

import (
	"google.golang.org/grpc"
)

type ProxyGrpc struct {
	*grpc.ClientConn
}

// NewProxyGrpc 缓存
func NewProxyGrpc() *ProxyGrpc {
	return &ProxyGrpc{}
}

// Store 设置写库
func (proxy *ProxyGrpc) Store(client *grpc.ClientConn) {
	proxy.ClientConn = client
}

// Store 设置组
func (proxypool *ProxyPool) StoreGrpc(group string, proxy *ProxyGrpc) {
	proxypool.m.Store(group, proxy)
}

// Load 获取分组
func (proxypool *ProxyPool) LoadGrpc(group string) *grpc.ClientConn {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyGrpc).ClientConn
	}
	return nil
}
