package proxy

import (
	"github.com/abulo/ratel/v3/store/hbase"
)

type ProxyHbase struct {
	*hbase.Client
}

//NewProxyHbase 缓存
func NewProxyHbase() *ProxyHbase {
	return &ProxyHbase{}
}

//Store 设置写库
func (proxy *ProxyHbase) Store(client *hbase.Client) {
	proxy.Client = client
}

//StoreEs 设置组
func (proxypool *ProxyPool) StoreHbase(group string, proxy *ProxyHbase) {
	proxypool.m.Store(group, proxy)
}

//LoadEs 获取分组
func (proxypool *ProxyPool) LoadHbase(group string) *hbase.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyHbase).Client
	}
	return nil
}
