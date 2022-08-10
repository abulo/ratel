package proxy

import "github.com/abulo/ratel/v3/stores/elasticsearch"

// ProxyElasticSearch ...
type ProxyElasticSearch struct {
	*elasticsearch.Client
}

// NewProxyElasticSearch 缓存
func NewProxyElasticSearch() *ProxyElasticSearch {
	return &ProxyElasticSearch{}
}

// Store 设置写库
func (proxy *ProxyElasticSearch) Store(client *elasticsearch.Client) {
	proxy.Client = client
}

// StoreElasticSearch StoreEs 设置组
func (proxypool *ProxyPool) StoreElasticSearch(group string, proxy *ProxyElasticSearch) {
	proxypool.m.Store(group, proxy)
}

// LoadElasticSearch LoadEs 获取分组
func (proxypool *ProxyPool) LoadElasticSearch(group string) *elasticsearch.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyElasticSearch).Client
	}
	return nil
}
