package proxy

import (
	"github.com/abulo/ratel/v2/stores/query"
	"github.com/abulo/ratel/v2/util"
)

// SQL Proxy 代理
type SQL struct {
	write []*query.Query
	read  []*query.Query
}

// NewSQL 代理池
func NewSQL() *SQL {
	return &SQL{}
}

// SetWrite 设置写库
func (proxy *SQL) SetWrite(query *query.Query) {
	proxy.write = append(proxy.write, query)
}

// SetRead 设置读库
func (proxy *SQL) SetRead(query *query.Query) {
	proxy.read = append(proxy.read, query)
}

// Write 获取写库
func (proxy *SQL) Write() *query.Query {
	len := len(proxy.write)
	write := util.Rand(0, len-1)
	return proxy.write[write]
}

// Read 获取读库
func (proxy *SQL) Read() *query.Query {
	len := len(proxy.read)
	if len < 1 {
		return proxy.Write()
	}
	read := util.Rand(0, len-1)
	return proxy.read[read]
}

// StoreSQL 设置组
func (proxypool *Proxy) StoreSQL(group string, proxy *SQL) {
	proxypool.m.Store(group, proxy)
}

// LoadSQL 获取分组
func (proxypool *Proxy) LoadSQL(group string) *SQL {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*SQL)
	}
	return nil
}
