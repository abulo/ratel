package proxy

import "github.com/abulo/ratel/v2/store/elasticsearch"

type ProxyElasticSearch struct {
	*elasticsearch.Client
}

//NewProxyElasticSearch 缓存
func NewProxyElasticSearch() *ProxyElasticSearch {
	return &ProxyElasticSearch{}
}

//Store 设置写库
func (proxy *ProxyElasticSearch) Store(client *elasticsearch.Client) {
	proxy.Client = client
}

//StoreEs 设置组
func (proxypool *ProxyPool) StoreElasticSearch(group string, proxy *ProxyElasticSearch) {
	proxypool.m.Store(group, proxy)
}

//LoadEs 获取分组
func (proxypool *ProxyPool) LoadElasticSearch(group string) *ProxyElasticSearch {
	if f, ok := proxypool.m.Load(group); ok {
		return f.(*ProxyElasticSearch)
	}
	return nil
}
