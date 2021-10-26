package proxy

import "github.com/abulo/ratel/v2/store/redis"

type ProxyRedis struct {
	*redis.Client
}

//NewProxyRedis 缓存
func NewProxyRedis() *ProxyRedis {
	return &ProxyRedis{}
}

//Store 设置写库
func (proxy *ProxyRedis) Store(client *redis.Client) {
	proxy.Client = client
}

//StoreCache 设置组
func (proxypool *ProxyPool) StoreRedis(group string, proxy *ProxyRedis) {
	proxypool.m.Store(group, proxy)
}

//LoadCache 获取分组
func (proxypool *ProxyPool) LoadRedis(group string) *ProxyRedis {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyRedis)
	}
	return nil
}
