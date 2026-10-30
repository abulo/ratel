package proxy

import "github.com/abulo/ratel/v3/stores/elasticsearch"

// ElasticSearch 包装 elasticsearch.Client，作为代理池中 Elasticsearch 后端的载体。
//
// 零值的底层客户端字段为 nil，须先调用 Store 注入底层客户端后再使用。
type ElasticSearch struct {
	*elasticsearch.Client
}

// NewElasticSearch 构造一个空的 Elasticsearch 代理实例，待 Store 注入底层客户端。
func NewElasticSearch() *ElasticSearch {
	return &ElasticSearch{}
}

// Store 向代理注入底层 Elasticsearch 客户端；传入 nil 会覆盖已有客户端，调用方须保证非空。
func (proxy *ElasticSearch) Store(client *elasticsearch.Client) {
	proxy.Client = client
}

// StoreElasticSearch 将 Elasticsearch 代理按分组名注册到代理池，同名分组会被覆盖。
func (proxyPool *Proxy) StoreElasticSearch(group string, proxy *ElasticSearch) {
	proxyPool.m.Store(group, proxy)
}

// LoadElasticSearch 按分组名取出已注册的 Elasticsearch 客户端；分组不存在时返回 nil。
func (proxyPool *Proxy) LoadElasticSearch(group string) *elasticsearch.Client {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*ElasticSearch).Client
	}
	return nil
}
