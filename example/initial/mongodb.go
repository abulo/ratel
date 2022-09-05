package initial

import (
	"time"

	"github.com/abulo/ratel/v3/stores/mongodb"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/spf13/cast"
)

// InitMongoDB load mongodb && returns an mongodb instance.
func (initial *Initial) InitMongoDB() *Initial {
	configs := initial.Config.Get("mongodb")
	list := configs.(map[string]interface{})
	links := make(map[string]*mongodb.MongoDB)
	for node, nodeConfig := range list {
		opt := &mongodb.Config{}
		res := nodeConfig.(map[string]interface{})
		if URI := cast.ToString(res["URI"]); URI != "" {
			opt.URI = URI
		}
		if MaxConnIdleTime := cast.ToInt64(res["MaxConnIdleTime"]); MaxConnIdleTime > 0 {
			opt.MaxConnIdleTime = cast.ToDuration(MaxConnIdleTime) * time.Minute
		}
		if MaxPoolSize := cast.ToInt64(res["MaxPoolSize"]); MaxPoolSize > 0 {
			opt.MaxPoolSize = cast.ToUint64(MaxPoolSize)
		}
		if MinPoolSize := cast.ToInt64(res["MinPoolSize"]); MinPoolSize > 0 {
			opt.MinPoolSize = cast.ToUint64(MinPoolSize)
		}

		opt.DisableMetric = cast.ToBool(res["DisableMetric"])
		opt.DisableTrace = cast.ToBool(res["DisableTrace"])
		conn := mongodb.NewClient(opt)
		links["mongodb."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxymongodb")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewMongoDB()
		if node := cast.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreMongoDB(Name, proxyPool)
		}
	}
	return initial
}
