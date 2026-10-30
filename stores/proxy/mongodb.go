package proxy

import "github.com/abulo/ratel/v3/stores/mongodb"

// MongoDB 包装 mongodb.MongoDB，作为代理池中 MongoDB 后端的载体。
//
// 零值的底层实例字段为 nil，须先调用 Store 注入底层客户端后再使用。
type MongoDB struct {
	*mongodb.MongoDB
}

// NewMongoDB 构造一个空的 MongoDB 代理实例，待 Store 注入底层客户端。
func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

// Store 向代理注入底层 MongoDB 客户端；传入 nil 会覆盖已有客户端，调用方须保证非空。
func (proxy *MongoDB) Store(client *mongodb.MongoDB) {
	proxy.MongoDB = client
}

// StoreMongoDB 将 MongoDB 代理按分组名注册到代理池，同名分组会被覆盖。
func (proxyPool *Proxy) StoreMongoDB(group string, proxy *MongoDB) {
	proxyPool.m.Store(group, proxy)
}

// LoadMongoDB 按分组名取出已注册的 MongoDB 客户端；分组不存在时返回 nil。
func (proxyPool *Proxy) LoadMongoDB(group string) *mongodb.MongoDB {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*MongoDB).MongoDB
	}
	return nil
}
