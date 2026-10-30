package proxy

import "sync"

// Proxy 是按分组管理各类存储代理实例的代理池，实例的注册与获取支持并发调用；
// 分组不存在时对应的 Load 方法返回 nil。
type Proxy struct {
	m sync.Map
}

// NewProxy 构造一个空的 Proxy 代理池，随后通过 StoreRedis、StoreMongoDB 等方法按分组注册实例。
func NewProxy() *Proxy {
	return &Proxy{
		m: sync.Map{},
	}
}
