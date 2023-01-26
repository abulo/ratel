package client

import (
	"google.golang.org/grpc"
)

// Grpc ...
type Grpc struct {
	*grpc.ClientConn
}

// NewGrpc 缓存
func NewGrpc() *Grpc {
	return &Grpc{}
}

// Store 设置写库
func (proxy *Grpc) Store(client *grpc.ClientConn) {
	proxy.ClientConn = client
}

// StoreGrpc Store 设置组
func (proxypool *Proxy) StoreGrpc(group string, proxy *Grpc) {
	proxypool.m.Store(group, proxy)
}

// LoadGrpc Load 获取分组
func (proxypool *Proxy) LoadGrpc(group string) *grpc.ClientConn {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*Grpc).ClientConn
	}
	return nil
}
