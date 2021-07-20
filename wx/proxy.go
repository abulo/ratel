package wx

import "sync"

//Proxy 代理
type Proxy struct {
	namespace map[string]*Client
	mu        sync.RWMutex
}

//NewProxy 代理池
func NewProxy() *Proxy {
	return &Proxy{
		namespace: make(map[string]*Client),
	}
}

//SetNameSpace 设置组
func (proxy *Proxy) SetNameSpace(group string, client *Client) {
	proxy.mu.Lock()
	proxy.namespace[group] = client
	proxy.mu.Unlock()
}

//NameSpace 获取分组
func (proxy *Proxy) NameSpace(group string) *Client {
	proxy.mu.RLock()
	res := proxy.namespace[group]
	proxy.mu.RUnlock()
	return res
}
