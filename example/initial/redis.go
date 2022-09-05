package initial

import (
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/abulo/ratel/v3/stores/redis"
	"github.com/spf13/cast"
)

// InitRedis load redis && returns an redis instance.
func (initial *Initial) InitRedis() *Initial {
	configs := initial.Config.Get("redis")
	list := configs.(map[string]interface{})
	links := make(map[string]*redis.Client)
	for node, nodeConfig := range list {
		opt := &redis.Config{}
		res := nodeConfig.(map[string]interface{})
		if KeyPrefix := cast.ToString(res["KeyPrefix"]); KeyPrefix != "" {
			opt.KeyPrefix = KeyPrefix
		}
		if Password := cast.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Database := cast.ToInt(res["Database"]); Database > 0 {
			opt.Database = cast.ToInt(Database)
		}
		if PoolSize := cast.ToInt(res["PoolSize"]); PoolSize > 0 {
			opt.PoolSize = cast.ToInt(PoolSize)
		}
		opt.Type = cast.ToBool(res["Type"])
		if Hosts := cast.ToStringSlice(res["Hosts"]); len(Hosts) > 0 {
			opt.Hosts = Hosts
		}
		opt.DisableMetric = cast.ToBool(res["DisableMetric"])
		opt.DisableTrace = cast.ToBool(res["DisableTrace"])
		conn := redis.New(opt)
		links["redis."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxyredis")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewRedis()
		if node := cast.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreRedis(Name, proxyPool)
		}
	}
	return initial
}
