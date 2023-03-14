package sql

import "time"

type pool struct {
	MaxLifetime    time.Duration // 连接池里面的连接最大存活时长
	MaxIdleTime    time.Duration // 连接池里面的连接最大空闲时长
	MaxOpenConns   int           // 连接池最多同时打开的连接数
	MaxIdleConns   int           // 连接池里最大空闲连接数。必须要比maxOpenConns小
	DisableMetric  bool          // 关闭指标采集
	DisableTrace   bool          // 关闭链路追踪
	DisablePrepare bool          // 关闭预处理
	DbName         string        // 数据库名称
	Addr           string        // 数据库地址
	DriverName     string        // 驱动名称
}
