package proxy

import "sync"

// Proxy 代理池
type Proxy struct {
	m sync.Map
}

// NewProxy 代理池
func NewProxy() *Proxy {
	return &Proxy{
		m: sync.Map{},
	}
}
