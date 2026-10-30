package proxy

import (
	"github.com/abulo/ratel/v3/stores/redis"
)

// Redis 包装 redis.Client，作为代理池中 Redis 后端的载体。
//
// 零值的 Client 字段为 nil，须先调用 Store 注入底层客户端后再使用。
type Redis struct {
	*redis.Client
}

// NewRedis 构造一个空的 Redis 代理实例，待 Store 注入底层客户端。
func NewRedis() *Redis {
	return &Redis{}
}

// Store 向代理注入底层 Redis 客户端；传入 nil 会覆盖已有客户端，调用方须保证非空。
func (proxy *Redis) Store(client *redis.Client) {
	proxy.Client = client
}

// StoreRedis 将 Redis 代理按分组名注册到代理池，同名分组会被覆盖。
func (proxyPool *Proxy) StoreRedis(group string, proxy *Redis) {
	proxyPool.m.Store(group, proxy)
}

// LoadRedis 按分组名取出已注册的 Redis 客户端；分组不存在时返回 nil。
func (proxyPool *Proxy) LoadRedis(group string) *redis.Client {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*Redis).Client
	}
	return nil
}
