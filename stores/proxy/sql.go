package proxy

import (
	"github.com/abulo/ratel/v3/util"
	"gorm.io/gorm"
)

// SQL Proxy 代理
type SQL struct {
	write []*gorm.DB
	read  []*gorm.DB
}

// NewSQL 代理池
func NewSQL() *SQL {
	return &SQL{}
}

// SetWrite 设置写库
func (proxy *SQL) SetWrite(db *gorm.DB) {
	proxy.write = append(proxy.write, db)
}

// SetRead 设置读库
func (proxy *SQL) SetRead(db *gorm.DB) {
	proxy.read = append(proxy.read, db)
}

// Write 获取写库
func (proxy *SQL) Write() *gorm.DB {
	len := len(proxy.write)
	write := util.Rand(0, len-1)
	return proxy.write[write]
}

// Read 获取读库
func (proxy *SQL) Read() *gorm.DB {
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
