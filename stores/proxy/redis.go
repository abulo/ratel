package proxy

import "github.com/abulo/ratel/stores/redis"

// Redis ...
type Redis struct {
	*redis.Client
}

// NewRedis 缓存
func NewRedis() *Redis {
	return &Redis{}
}

// Store 设置写库
func (proxy *Redis) Store(client *redis.Client) {
	proxy.Client = client
}

// StoreRedis StoreCache 设置组
func (proxypool *Proxy) StoreRedis(group string, proxy *Redis) {
	proxypool.m.Store(group, proxy)
}

// LoadRedis LoadCache 获取分组
func (proxypool *Proxy) LoadRedis(group string) *redis.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*Redis).Client
	}
	return nil
}
