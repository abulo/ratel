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
)

type Initial struct {
	Path       string           // 应用程序执行路径
	Config     *config.Config   // 配置文件
	Store      *proxy.ProxyPool // 数据库链接
	Session    *session.Session // 回话保存实例
	LaunchTime time.Time        //时间设置
}

//系统
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
		if URI := util.ToString(res["URI"]); URI != "" {
			opt.URI = URI
		}
		if MaxConnIdleTime := util.ToInt64(res["MaxConnIdleTime"]); MaxConnIdleTime > 0 {
			opt.MaxConnIdleTime = util.ToDuration(MaxConnIdleTime) * time.Minute
		}
		if MaxPoolSize := util.ToInt64(res["MaxPoolSize"]); MaxPoolSize > 0 {
			opt.MaxPoolSize = util.ToUint64(MaxPoolSize)
		}
		if MinPoolSize := util.ToInt64(res["MinPoolSize"]); MinPoolSize > 0 {
			opt.MinPoolSize = util.ToUint64(MinPoolSize)
		}

		opt.DisableMetric = util.ToBool(res["DisableMetric"])
		opt.DisableTrace = util.ToBool(res["DisableTrace"])
		conn := mongodb.NewClient(opt)
		links["mongodb."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxymongodb")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewProxyMongoDB()
		if node := util.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := util.ToString(val["Name"]); Name != "" {
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
		if Username := util.ToString(res["Username"]); Username != "" {
			opt.Username = Username
		}
		if Password := util.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Host := util.ToString(res["Host"]); Host != "" {
			opt.Host = Host
		}
		if Port := util.ToString(res["Port"]); Port != "" {
			opt.Port = Port
		}
		if Charset := util.ToString(res["Charset"]); Charset != "" {
			opt.Charset = Charset
		}
		if Database := util.ToString(res["Database"]); Database != "" {
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

		if MaxLifetime := util.ToInt(res["MaxLifetime"]); MaxLifetime > 0 {
			opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
		}
		if MaxIdleTime := util.ToInt(res["MaxIdleTime"]); MaxIdleTime > 0 {
			opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
		}
		if MaxIdleConns := util.ToInt(res["MaxIdleConns"]); MaxIdleConns > 0 {
			opt.MaxIdleConns = util.ToInt(MaxIdleConns)
		}
		if MaxOpenConns := util.ToInt(res["MaxOpenConns"]); MaxOpenConns > 0 {
			opt.MaxOpenConns = util.ToInt(MaxOpenConns)
		}

		opt.DriverName = "mysql"
		opt.DisableMetric = util.ToBool(res["DisableMetric"])
		opt.DisableTrace = util.ToBool(res["DisableTrace"])
		conn := mysql.NewClient(opt)
		links["mysql."+node] = conn
	}

	proxyConfigs := initial.Config.Get("proxymysql")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewProxySQL()
		if Master := util.ToStringSlice(val["Master"]); len(Master) > 0 {
			for _, v := range Master {
				proxyPool.SetWrite(links[v])
			}
		}
		if Slave := util.ToStringSlice(val["Slave"]); len(Slave) > 0 {
			for _, v := range Slave {
				proxyPool.SetRead(links[v])
			}
		}
		if Name := util.ToString(val["Name"]); Name != "" {
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
		if Username := util.ToString(res["Username"]); Username != "" {
			opt.Username = Username
		}
		if Password := util.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Addr := util.ToStringSlice(res["Addr"]); len(Addr) > 0 {
			opt.Addr = Addr
		}
		if Database := util.ToString(res["Database"]); Database != "" {
			opt.Database = Database
		}
		if DialTimeout := util.ToString(res["DialTimeout"]); DialTimeout != "" {
			opt.DialTimeout = DialTimeout
		}
		if OpenStrategy := util.ToString(res["OpenStrategy"]); OpenStrategy != "" {
			opt.OpenStrategy = OpenStrategy
		}
		if Compress := util.ToBool(res["Compress"]); Compress {
			opt.Compress = true
		} else {
			opt.Compress = false
		}
		if MaxExecutionTime := util.ToString(res["MaxExecutionTime"]); MaxExecutionTime != "" {
			opt.MaxExecutionTime = MaxExecutionTime
		}

		opt.DisableDebug = util.ToBool(res["DisableDebug"])
		// # MaxOpenConns 连接池最多同时打开的连接数
		// MaxOpenConns = 128
		// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
		// MaxIdleConns = 32
		// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
		// MaxLifetime = 10
		// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
		// MaxIdleTime = 5

		if MaxLifetime := util.ToInt(res["MaxLifetime"]); MaxLifetime > 0 {
			opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
		}
		if MaxIdleTime := util.ToInt(res["MaxIdleTime"]); MaxIdleTime > 0 {
			opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
		}
		if MaxIdleConns := util.ToInt(res["MaxIdleConns"]); MaxIdleConns > 0 {
			opt.MaxIdleConns = util.ToInt(MaxIdleConns)
		}
		if MaxOpenConns := util.ToInt(res["MaxOpenConns"]); MaxOpenConns > 0 {
			opt.MaxOpenConns = util.ToInt(MaxOpenConns)
		}
		opt.DriverName = "clickhouse"
		opt.DisableDebug = util.ToBool(res["DisableDebug"])
		opt.DisableMetric = util.ToBool(res["DisableMetric"])
		opt.DisableTrace = util.ToBool(res["DisableTrace"])
		conn := clickhouse.NewClient(opt)
		links["clickhouse."+node] = conn
	}

	proxyConfigs := initial.Config.Get("proxyclickhouse")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewProxySQL()
		if node := util.ToString(val["Node"]); node != "" {
			proxyPool.SetWrite(links[node])
		}
		if Name := util.ToString(val["Name"]); Name != "" {
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
		opts.URL = util.ToStringSlice(res["URL"])
		opts.Username = util.ToString(res["Username"])
		opts.Password = util.ToString(res["Password"])
		opts.DisableMetric = util.ToBool(res["DisableMetric"])
		opts.DisableTrace = util.ToBool(res["DisableTrace"])
		conn := elasticsearch.NewClient(opts)
		links["elasticsearch."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxyelasticsearch")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewProxyElasticSearch()
		if node := util.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := util.ToString(val["Name"]); Name != "" {
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
		if KeyPrefix := util.ToString(res["KeyPrefix"]); KeyPrefix != "" {
			opt.KeyPrefix = KeyPrefix
		}
		if Password := util.ToString(res["Password"]); Password != "" {
			opt.Password = Password
		}
		if Database := util.ToInt(res["Database"]); Database > 0 {
			opt.Database = util.ToInt(Database)
		}
		if PoolSize := util.ToInt(res["PoolSize"]); PoolSize > 0 {
			opt.PoolSize = util.ToInt(PoolSize)
		}
		opt.Type = util.ToBool(res["Type"])
		if Hosts := util.ToStringSlice(res["Hosts"]); len(Hosts) > 0 {
			opt.Hosts = Hosts
		}
		opt.DisableMetric = util.ToBool(res["DisableMetric"])
		opt.DisableTrace = util.ToBool(res["DisableTrace"])
		conn := redis.New(opt)
		links["redis."+node] = conn
	}
	proxyConfigs := initial.Config.Get("proxyredis")
	proxyRes := proxyConfigs.([]map[string]interface{})
	for _, val := range proxyRes {
		proxyPool := proxy.NewProxyRedis()
		if node := util.ToStringSlice(val["Node"]); len(node) > 0 {
			for _, v := range node {
				proxyPool.Store(links[v])
			}
		}
		if Name := util.ToString(val["Name"]); Name != "" {
			initial.Store.StoreRedis(Name, proxyPool)
		}
	}
	return initial
}

func (initial *Initial) InitTrace() {
	opt := jaeger.NewJaeger()
	conf := initial.Config.Get("trace")
	res := conf.(map[string]interface{})
	opt.EnableRPCMetrics = util.ToBool(res["EnableRPCMetrics"])
	opt.LocalAgentHostPort = util.ToString(res["LocalAgentHostPort"])
	opt.LogSpans = util.ToBool(res["LogSpans"])
	opt.Param = util.ToFloat64(res["Param"])
	opt.PanicOnError = util.ToBool(res["PanicOnError"])
	client := opt.Build().Build()
	trace.SetGlobalTracer(client)
}
