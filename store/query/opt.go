package query

import "time"

type Opt struct {
	MaxLifetime  time.Duration //连接池里面的连接最大存活时长
	MaxIdleTime  time.Duration //连接池里面的连接最大空闲时长
	MaxOpenConns int           //连接池最多同时打开的连接数
	MaxIdleConns int           //连接池里最大空闲连接数。必须要比maxOpenConns小
}
