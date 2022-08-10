package proxy

import (
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
)

// ProxySQL Proxy 代理
type ProxySQL struct {
	write []*query.QueryDb
	read  []*query.QueryDb
}

// NewProxySQL 代理池
func NewProxySQL() *ProxySQL {
	return &ProxySQL{}
}

// SetWrite 设置写库
func (proxy *ProxySQL) SetWrite(query *query.QueryDb) {
	proxy.write = append(proxy.write, query)
}

// SetRead 设置读库
func (proxy *ProxySQL) SetRead(query *query.QueryDb) {
	proxy.read = append(proxy.read, query)
}

// Write 获取写库
func (proxy *ProxySQL) Write() *query.QueryDb {
	len := len(proxy.write)
	write := util.Rand(0, len-1)
	return proxy.write[write]
}

// Read 获取读库
func (proxy *ProxySQL) Read() *query.QueryDb {
	len := len(proxy.read)
	if len < 1 {
		return proxy.Write()
	}
	read := util.Rand(0, len-1)
	return proxy.read[read]
}

// StoreSQL 设置组
func (proxypool *ProxyPool) StoreSQL(group string, proxy *ProxySQL) {
	proxypool.m.Store(group, proxy)
}

// LoadSQL 获取分组
func (proxypool *ProxyPool) LoadSQL(group string) *ProxySQL {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxySQL)
	}
	return nil
}
