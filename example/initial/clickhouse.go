package initial

import (
	"time"

	"github.com/abulo/ratel/v3/stores/clickhouse"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/spf13/cast"
)

// InitClickHouse load clickhouse && returns an clickhouse instance.
func (initial *Initial) InitClickHouse() *Initial {
	configs := initial.Config.Get("clickhouse")
	list := configs.(map[string]interface{})

	links := make(map[string]*query.Query)
	for node, nodeConfig := range list {
		opt := &clickhouse.Config{}
		res := nodeConfig.(map[string]interface{})
		if Username := cast.ToString(res["Username"]); Username != "" {
			opt.Username = Username
		}
		if Password := cast.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Addr := cast.ToStringSlice(res["Addr"]); len(Addr) > 0 {
			opt.Addr = Addr
		}
		if Database := cast.ToString(res["Database"]); Database != "" {
			opt.Database = Database
		}
		if Local := cast.ToString(res["Local"]); Local != "" {
			opt.Local = Local
		}
		if DialTimeout := cast.ToString(res["DialTimeout"]); DialTimeout != "" {
			opt.DialTimeout = DialTimeout
		}
		if OpenStrategy := cast.ToString(res["OpenStrategy"]); OpenStrategy != "" {
			opt.OpenStrategy = OpenStrategy
		}
		if Compress := cast.ToBool(res["Compress"]); Compress {
			opt.Compress = true
		} else {
			opt.Compress = false
		}
		if MaxExecutionTime := cast.ToString(res["MaxExecutionTime"]); MaxExecutionTime != "" {
			opt.MaxExecutionTime = MaxExecutionTime
		}

		opt.DisableDebug = cast.ToBool(res["DisableDebug"])
		// # MaxOpenConns 连接池最多同时打开的连接数
		// MaxOpenConns = 128
		// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
		// MaxIdleConns = 32
		// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
		// MaxLifetime = 10
		// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
		// MaxIdleTime = 5

		if MaxLifetime := cast.ToInt(res["MaxLifetime"]); MaxLifetime > 0 {
			opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
		}
		if MaxIdleTime := cast.ToInt(res["MaxIdleTime"]); MaxIdleTime > 0 {
			opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
		}
		if MaxIdleConns := cast.ToInt(res["MaxIdleConns"]); MaxIdleConns > 0 {
			opt.MaxIdleConns = cast.ToInt(MaxIdleConns)
		}
		if MaxOpenConns := cast.ToInt(res["MaxOpenConns"]); MaxOpenConns > 0 {
			opt.MaxOpenConns = cast.ToInt(MaxOpenConns)
		}
		opt.DriverName = "clickhouse"
		opt.DisableDebug = cast.ToBool(res["DisableDebug"])
		opt.DisableMetric = cast.ToBool(res["DisableMetric"])
		opt.DisableTrace = cast.ToBool(res["DisableTrace"])
		conn := clickhouse.NewClient(opt)
		links["clickhouse."+node] = conn
	}

	proxyConfigs := initial.Config.Get("proxyclickhouse")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewSQL()
		if node := cast.ToString(val["Node"]); node != "" {
			proxyPool.SetWrite(links[node])
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreSQL(Name, proxyPool)
		}
	}
	return initial
}
