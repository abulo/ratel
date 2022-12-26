package clickhouse

import (
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abulo/ratel/v2/core/logger"
	"github.com/abulo/ratel/v2/stores/query"
	"github.com/abulo/ratel/v2/util"
)

// Config 数据库配置
type Config struct {
	Username         string   //账号 root
	Password         string   //密码
	Addr             []string //ip:port
	Database         string   //连接数据库
	Local            string   //数据库时区
	DialTimeout      string   //200ms
	OpenStrategy     string   //random/in_order (default random)
	Compress         bool     //enable lz4 compression
	MaxExecutionTime string
	MaxOpenConns     int           //连接池最多同时打开的连接数
	MaxIdleConns     int           //连接池里最大空闲连接数。必须要比maxOpenConns小
	MaxLifetime      time.Duration //连接池里面的连接最大存活时长
	MaxIdleTime      time.Duration //连接池里面的连接最大空闲时长
	DriverName       string
	DisableDebug     bool // 关闭 debug模式
	DisableMetric    bool // 关闭指标采集
	DisableTrace     bool // 关闭链路追踪
}

// URI 构造数据库连接
func (config *Config) URI() string {
	//clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60

	link := "clickhouse://"
	if !util.Empty(config.Username) {
		link = link + config.Username
	}
	link = link + ":"
	if !util.Empty(config.Password) {
		link = link + config.Password
	}
	link = link + "@"
	link = link + util.Implode(",", config.Addr)
	link = link + "/"
	if !util.Empty(config.Database) {
		link = link + config.Database
	}
	param := make([]string, 0)
	if !util.Empty(config.DialTimeout) {
		param = append(param, "dial_timeout="+config.DialTimeout)
	}
	if !util.Empty(config.MaxExecutionTime) {
		param = append(param, "max_execution_time="+config.MaxExecutionTime)
	}
	if !util.Empty(config.OpenStrategy) {
		param = append(param, "connection_open_strategy="+config.OpenStrategy)
	}
	if config.Compress {
		param = append(param, "compress=true")
	} else {
		param = append(param, "compress=false")
	}
	if !config.DisableDebug {
		param = append(param, "debug=true")
	} else {
		param = append(param, "debug=false")
	}
	// param = append(param, "loc="+config.Local)
	// param = append(param, "parseTime=true")
	return link + "?" + util.Implode("&", param)
}

// NewClient New 新连接
func NewClient(config *Config) *query.Query {
	opt := &query.Opt{
		MaxOpenConns: config.MaxOpenConns,
		MaxIdleConns: config.MaxIdleConns,
		MaxLifetime:  config.MaxLifetime,
		MaxIdleTime:  config.MaxIdleTime,
	}

	db, err := query.NewSQLConn(config.DriverName, config.URI(), opt)
	if err != nil {
		logger.Logger.Panic(err)
	}
	return &query.Query{DB: db, DriverName: config.DriverName, DisableMetric: config.DisableMetric, DisableTrace: config.DisableTrace, Prepare: false, DBName: config.Database, Addr: util.Implode(";", config.Addr)}
}
