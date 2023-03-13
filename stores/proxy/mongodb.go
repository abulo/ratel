package proxy

import "github.com/abulo/ratel/v3/stores/mongodb"

// MongoDB ...
type MongoDB struct {
	*mongodb.MongoDB
}

// NewMongoDB 缓存
func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

// Store 设置写库
func (proxy *MongoDB) Store(client *mongodb.MongoDB) {
	proxy.MongoDB = client
}

// StoreMongoDB StoreNoSQL 设置组
func (proxyPool *Proxy) StoreMongoDB(group string, proxy *MongoDB) {
	proxyPool.m.Store(group, proxy)
}

// LoadMongoDB LoadNoSQL 获取分组
func (proxyPool *Proxy) LoadMongoDB(group string) *mongodb.MongoDB {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*MongoDB).MongoDB
	}
	return nil
}
