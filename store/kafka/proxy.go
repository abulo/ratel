package kafka

import "sync"

//Proxy 代理
type Proxy struct {
	namespace map[string]*Clientkafka
	mu        sync.RWMutex
}

//NewProxy 代理池
func NewProxy() *Proxy {
	return &Proxy{
		namespace: make(map[string]*Clientkafka),
	}
}

//SetNameSpace 设置组
func (proxy *Proxy) SetNameSpace(group string, client *Clientkafka) {
	proxy.mu.Lock()
	proxy.namespace[group] = client
	proxy.mu.Unlock()
}

//NameSpace 获取分组
func (proxy *Proxy) NameSpace(group string) *Clientkafka {
	proxy.mu.RLock()
	res := proxy.namespace[group]
	proxy.mu.RUnlock()
	return res
}
