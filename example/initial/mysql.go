package initial

import (
	"time"

	"github.com/abulo/ratel/v3/stores/mysql"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/spf13/cast"
)

// InitMysql load mysql && returns an mysql instance.
func (initial *Initial) InitMysql() *Initial {
	configs := initial.Config.Get("mysql")
	list := configs.(map[string]interface{})
	links := make(map[string]*query.Query)
	for node, nodeConfig := range list {
		opt := &mysql.Config{}
		res := nodeConfig.(map[string]interface{})
		if Username := cast.ToString(res["Username"]); Username != "" {
			opt.Username = Username
		}
		if Password := cast.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Host := cast.ToString(res["Host"]); Host != "" {
			opt.Host = Host
		}
		if Port := cast.ToString(res["Port"]); Port != "" {
			opt.Port = Port
		}
		if Charset := cast.ToString(res["Charset"]); Charset != "" {
			opt.Charset = Charset
		}
		if Database := cast.ToString(res["Database"]); Database != "" {
			opt.Database = Database
		}

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

		opt.DriverName = "mysql"
		opt.DisableMetric = cast.ToBool(res["DisableMetric"])
		opt.DisableTrace = cast.ToBool(res["DisableTrace"])
		conn := mysql.NewClient(opt)
		links["mysql."+node] = conn
	}

	proxyConfigs := initial.Config.Get("proxymysql")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewSQL()
		if Master := cast.ToStringSlice(val["Master"]); len(Master) > 0 {
			for _, v := range Master {
				proxyPool.SetWrite(links[v])
			}
		}
		if Slave := cast.ToStringSlice(val["Slave"]); len(Slave) > 0 {
			for _, v := range Slave {
				proxyPool.SetRead(links[v])
			}
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreSQL(Name, proxyPool)
		}
	}
	return initial
}
