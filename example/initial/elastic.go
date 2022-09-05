package initial

import (
	"github.com/abulo/ratel/v3/stores/elasticsearch"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/spf13/cast"
)

// InitElasticSearch load elasticsearch && returns an elasticsearch instance.
func (initial *Initial) InitElasticSearch() *Initial {
	configs := initial.Config.Get("elasticsearch")
	list := configs.(map[string]interface{})
	links := make(map[string]*elasticsearch.Client)
	for node, nodeConfig := range list {
		opts := &elasticsearch.Config{}
		res := nodeConfig.(map[string]interface{})
		opts.URL = cast.ToStringSlice(res["URL"])
		opts.Username = cast.ToString(res["Username"])
		opts.Password = cast.ToString(res["Password"])
		opts.DisableMetric = cast.ToBool(res["DisableMetric"])
		opts.DisableTrace = cast.ToBool(res["DisableTrace"])
		conn := elasticsearch.NewClient(opts)
		links["elasticsearch."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxyelasticsearch")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewElasticSearch()
		if node := cast.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreElasticSearch(Name, proxyPool)
		}
	}
	return initial
}
