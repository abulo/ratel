package proxy

import "github.com/abulo/ratel/v3/stores/redis"

// ProxyRedis ...
type ProxyRedis struct {
	*redis.Client
}

// NewProxyRedis 缓存
func NewProxyRedis() *ProxyRedis {
	return &ProxyRedis{}
}

// Store 设置写库
func (proxy *ProxyRedis) Store(client *redis.Client) {
	proxy.Client = client
}

// StoreRedis StoreCache 设置组
func (proxypool *ProxyPool) StoreRedis(group string, proxy *ProxyRedis) {
	proxypool.m.Store(group, proxy)
}

// LoadRedis LoadCache 获取分组
func (proxypool *ProxyPool) LoadRedis(group string) *redis.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyRedis).Client
	}
	return nil
}
