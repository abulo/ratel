package proxy

import (
	"github.com/abulo/ratel/v2/store/kylin"
)

type ProxyKylin struct {
	*kylin.Client
}

//NewProxyKylin 缓存
func NewProxyKylin() *ProxyKylin {
	return &ProxyKylin{}
}

//Store 设置写库
func (proxy *ProxyKylin) Store(client *kylin.Client) {
	proxy.Client = client
}

//StoreEs 设置组
func (proxypool *ProxyPool) StoreKylin(group string, proxy *ProxyKylin) {
	proxypool.m.Store(group, proxy)
}

//LoadEs 获取分组
func (proxypool *ProxyPool) LoadKylin(group string) *kylin.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyKylin).Client
	}
	return nil
}
