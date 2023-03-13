package sql

import (
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/pkg/errors"
)

const (
	driverMysql      = "mysql"
	driverClickhouse = "clickhouse"
	driverPostgres   = "postgres"
)

type (
	Option func(r *Client)
	Client struct {
		Username         string        // 账号 root
		Password         string        // 密码
		Host             string        // host localhost   mysql
		Port             string        // 端口 3306
		Charset          string        // 字符编码 utf8mb4
		Database         string        // 默认连接数据库
		TimeZone         string        // 数据库时区
		MaxOpenConns     int           // 连接池最多同时打开的连接数
		MaxIdleConns     int           // 连接池里最大空闲连接数。必须要比maxOpenConns小
		MaxLifetime      time.Duration // 连接池里面的连接最大存活时长
		MaxIdleTime      time.Duration // 连接池里面的连接最大空闲时长
		DriverName       string        // 驱动名称
		DisableMetric    bool          // 关闭指标采集
		DisableTrace     bool          // 关闭链路追踪
		DisablePrepare   bool          // 关闭预处理
		Addr             []string      // ip:port  clickhouse
		DialTimeout      string        // 200ms    clickhouse
		OpenStrategy     string        // random/in_order (default random)  clickhouse
		Compress         bool          // enable lz4 compression   clickhouse
		MaxExecutionTime string        // 执行超时时间 clickhouse
		DisableDebug     bool          // 关闭 debug模式  clickhouse
	}
)

func NewSqlClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// driverUrl
func (c *Client) dns() string {
	switch c.DriverName {
	case driverMysql:
		return c.mysqlDns()
	case driverClickhouse:
		return c.clickhouseDns()
	case driverPostgres:
		return c.postgresDns()
	default:
		logger.Logger.Panic(errors.New("driverName not support"))
		return ""
	}
}
