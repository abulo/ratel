package proxy

import "github.com/abulo/ratel/v2/stores/elasticsearch"

// ElasticSearch ...
type ElasticSearch struct {
	*elasticsearch.Client
}

// NewElasticSearch 缓存
func NewElasticSearch() *ElasticSearch {
	return &ElasticSearch{}
}

// Store 设置写库
func (proxy *ElasticSearch) Store(client *elasticsearch.Client) {
	proxy.Client = client
}

// StoreElasticSearch StoreEs 设置组
func (proxypool *Proxy) StoreElasticSearch(group string, proxy *ElasticSearch) {
	proxypool.m.Store(group, proxy)
}

// LoadElasticSearch LoadEs 获取分组
func (proxypool *Proxy) LoadElasticSearch(group string) *elasticsearch.Client {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ElasticSearch).Client
	}
	return nil
}
