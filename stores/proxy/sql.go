package proxy

import (
	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/abulo/ratel/v3/util"
)

// SQL Proxy 代理
type SQL struct {
	write []sql.SqlConn
	read  []sql.SqlConn
}

// NewSQL 代理池
func NewSQL() *SQL {
	return &SQL{}
}

// SetWrite 设置写库
func (proxy *SQL) SetWrite(query sql.SqlConn) {
	proxy.write = append(proxy.write, query)
}

// SetRead 设置读库
func (proxy *SQL) SetRead(query sql.SqlConn) {
	proxy.read = append(proxy.read, query)
}

// Write 获取写库
func (proxy *SQL) Write() sql.SqlConn {
	len := len(proxy.write)
	write := util.Rand(0, len-1)
	return proxy.write[write]
}

// Read 获取读库
func (proxy *SQL) Read() sql.SqlConn {
	len := len(proxy.read)
	if len < 1 {
		return proxy.Write()
	}
	read := util.Rand(0, len-1)
	return proxy.read[read]
}

// StoreSQL 设置组
func (proxyPool *Proxy) StoreSQL(group string, proxy *SQL) {
	proxyPool.m.Store(group, proxy)
}

// LoadSQL 获取分组
func (proxyPool *Proxy) LoadSQL(group string) *SQL {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*SQL)
	}
	return nil
}
