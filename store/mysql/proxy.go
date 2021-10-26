package mysql

import (
	"sync"

	"github.com/abulo/ratel/v2/store/query"
	"github.com/abulo/ratel/v2/util"
)

type (
	//Proxy 代理
	Proxy struct {
		write []*query.QueryDb
		read  []*query.QueryDb
	}
	//ProxyPool 代理池
	ProxyPool struct {
		namespace map[string]*Proxy
		mu        sync.RWMutex
	}
)

//NewProxyPool 代理池
func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		namespace: make(map[string]*Proxy),
	}
}

//NewProxy 代理池
func NewProxy() *Proxy {
	return &Proxy{
		write: make([]*query.QueryDb, 0),
		read:  make([]*query.QueryDb, 0),
	}
}

//SetNameSpace 设置组
func (proxypool *ProxyPool) SetNameSpace(group string, proxy *Proxy) {
	proxypool.mu.Lock()
	proxypool.namespace[group] = proxy
	proxypool.mu.Unlock()
}

//NameSpace 获取分组
func (proxypool *ProxyPool) NameSpace(group string) *Proxy {
	proxypool.mu.RLock()
	res := proxypool.namespace[group]
	proxypool.mu.RUnlock()
	return res
}

//SetWrite 设置写库
func (proxy *Proxy) SetWrite(query *query.QueryDb) {
	proxy.write = append(proxy.write, query)
}

//SetRead 设置读库
func (proxy *Proxy) SetRead(query *query.QueryDb) {
	proxy.read = append(proxy.read, query)
}

//Write 获取写库
func (proxy *Proxy) Write() *query.QueryDb {
	len := len(proxy.write)
	write := util.Rand(0, len-1)
	return proxy.write[write]
}

//Read 获取读库
func (proxy *Proxy) Read() *query.QueryDb {
	len := len(proxy.read)
	if len < 1 {
		return proxy.Write()
	}
	read := util.Rand(0, len-1)
	return proxy.read[read]
}
