package proxy

import "github.com/abulo/ratel/v3/stores/mongodb"

type ProxyMongoDB struct {
	*mongodb.MongoDB
}

// NewProxyMongoDB 缓存
func NewProxyMongoDB() *ProxyMongoDB {
	return &ProxyMongoDB{}
}

// Store 设置写库
func (proxy *ProxyMongoDB) Store(client *mongodb.MongoDB) {
	proxy.MongoDB = client
}

// StoreNoSQL 设置组
func (proxypool *ProxyPool) StoreMongoDB(group string, proxy *ProxyMongoDB) {
	proxypool.m.Store(group, proxy)
}

// LoadNoSQL 获取分组
func (proxypool *ProxyPool) LoadMongoDB(group string) *mongodb.MongoDB {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyMongoDB).MongoDB
	}
	return nil
}
