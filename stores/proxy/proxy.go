package proxy

import "sync"

//ProxyPool 代理池
type ProxyPool struct {
	m sync.Map
}

//NewProxyPool 代理池
func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		m: sync.Map{},
	}
}
