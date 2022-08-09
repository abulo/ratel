package initial

import (
	"time"

	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/toml"
	"github.com/abulo/ratel/v3/stores/clickhouse"
	"github.com/abulo/ratel/v3/stores/elasticsearch"
	"github.com/abulo/ratel/v3/stores/mongodb"
	"github.com/abulo/ratel/v3/stores/mysql"
	"github.com/abulo/ratel/v3/stores/proxy"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/stores/redis"
	"github.com/abulo/ratel/v3/stores/session"
	"github.com/abulo/ratel/v3/trace"
	"github.com/abulo/ratel/v3/trace/jaeger"
	"github.com/abulo/ratel/v3/util"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
)

type Initial struct {
	Path       string           // 应用程序执行路径
	Config     *config.Config   // 配置文件
	Store      *proxy.ProxyPool // 数据库链接
	Session    *session.Session // 回话保存实例
	LaunchTime time.Time        //时间设置
}

// 系统
var Core *Initial

// Default returns an Initial instance.
func Default() *Initial {
	engine := New()
	return engine
}

func New() *Initial {
	Core = &Initial{
		Store: proxy.NewProxyPool(),
	}
	Core.InitPath(util.GetAppRootPath())
	Core.InitLaunchTime(util.Now())
	return Core
}

// InitLaunchTime set binary start time
func (initial *Initial) InitLaunchTime(launchTime time.Time) *Initial {
	initial.LaunchTime = launchTime
	return initial
}

// InitPath binary file path
func (initial *Initial) InitPath(path string) *Initial {
	initial.Path = path
	return initial
}

func (initial *Initial) GetEnvironment(dir, key string) string {
	envConfig := config.NewWithOptions("go-ratel-evn", config.Readonly, config.EnableCache)
	driver := toml.Driver
	envConfig.AddDriver(driver)
	envConfig.LoadDir(dir, driver.Name())
	return envConfig.String(key)
}

// InitConfig set app config toml files
func (initial *Initial) InitConfig(dirs ...string) *Initial {
	Config := config.NewWithOptions("go-ratel", config.Readonly, config.EnableCache)
	driver := toml.Driver
	Config.AddDriver(driver)
	for _, dir := range dirs {
		isDir, err := util.IsDir(dir)
		if err != nil {
			panic(dir)
		}
		if !isDir {
			panic(dir + "not a directory")
		}
		//load file
		Config.LoadDir(dir, driver.Name())
	}
	Config.OnConfigChange(func(e fsnotify.Event) {
		Config.LoadFiles(e.Name)
	})
	Config.WatchConfig(driver.Name())
	initial.Config = Config
	return initial
}

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
		proxyPool := proxy.NewProxyMongoDB()
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

// InitSession load session && returns an session instance.
func (initial *Initial) InitSession(name string) *Initial {
	initial.Session = &session.Session{
		Name:   name,
		Driver: initial.Store.LoadRedis("redis"),
		TTL:    initial.Config.Int64("cookie.expires", 300),
	}
	return initial
}

// InitMysql load mysql && returns an mysql instance.
func (initial *Initial) InitMysql() *Initial {
	configs := initial.Config.Get("mysql")
	list := configs.(map[string]interface{})
	links := make(map[string]*query.QueryDb)
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
		proxyPool := proxy.NewProxySQL()
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

// InitClickHouse load clickhouse && returns an clickhouse instance.
func (initial *Initial) InitClickHouse() *Initial {
	configs := initial.Config.Get("clickhouse")
	list := configs.(map[string]interface{})

	links := make(map[string]*query.QueryDb)
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
		proxyPool := proxy.NewProxySQL()
		if node := cast.ToString(val["Node"]); node != "" {
			proxyPool.SetWrite(links[node])
		}
		if Name := cast.ToString(val["Name"]); Name != "" {
			initial.Store.StoreSQL(Name, proxyPool)
		}
	}
	return initial
}

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
		proxyPool := proxy.NewProxyElasticSearch()
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
		proxyPool := proxy.NewProxyRedis()
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

func (initial *Initial) InitTrace() {
	opt := jaeger.NewJaeger()
	conf := initial.Config.Get("trace")
	res := conf.(map[string]interface{})
	opt.EnableRPCMetrics = cast.ToBool(res["EnableRPCMetrics"])
	opt.LocalAgentHostPort = cast.ToString(res["LocalAgentHostPort"])
	opt.LogSpans = cast.ToBool(res["LogSpans"])
	opt.Param = cast.ToFloat64(res["Param"])
	opt.PanicOnError = cast.ToBool(res["PanicOnError"])
	client := opt.Build().Build()
	trace.SetGlobalTracer(client)
}
