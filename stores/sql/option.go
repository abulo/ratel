package sql

import "time"

// WithUsername 设置账号
func WithUsername(username string) Option {
	return func(r *Client) {
		r.Username = username
	}
}

// WithPassword 设置密码
func WithPassword(password string) Option {
	return func(r *Client) {
		r.Password = password
	}
}

// WithHost 设置host
func WithHost(host string) Option {
	return func(r *Client) {
		r.Host = host
	}
}

// WithPort 设置端口
func WithPort(port string) Option {
	return func(r *Client) {
		r.Port = port
	}
}

// WithCharset 设置字符编码
func WithCharset(charset string) Option {
	return func(r *Client) {
		r.Charset = charset
	}
}

// WithDatabase 设置默认连接数据库
func WithDatabase(database string) Option {
	return func(r *Client) {
		r.Database = database
	}
}

// WithTimeZone 设置数据库时区
func WithTimeZone(timeZone string) Option {
	return func(r *Client) {
		r.TimeZone = timeZone
	}
}

// WithMaxOpenConns 设置连接池最多同时打开的连接数
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(r *Client) {
		r.MaxOpenConns = maxOpenConns
	}
}

// WithMaxIdleConns 设置连接池里最大空闲连接数。必须要比maxOpenConns小
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(r *Client) {
		r.MaxIdleConns = maxIdleConns
	}
}

// WithMaxLifetime 设置连接池里面的连接最大存活时长
func WithMaxLifetime(maxLifetime time.Duration) Option {
	return func(r *Client) {
		r.MaxLifetime = maxLifetime
	}
}

// WithMaxIdleTime 设置连接池里面的连接最大空闲时长
func WithMaxIdleTime(maxIdleTime time.Duration) Option {
	return func(r *Client) {
		r.MaxIdleTime = maxIdleTime
	}
}

// WithDriverName 设置驱动名称
func WithDriverName(driverName string) Option {
	return func(r *Client) {
		r.DriverName = driverName
	}
}

// WithDisableMetric 关闭指标采集
func WithDisableMetric(disableMetric bool) Option {
	return func(r *Client) {
		r.DisableMetric = disableMetric
	}
}

// WithDisableTrace 关闭链路追踪
func WithDisableTrace(disableTrace bool) Option {
	return func(r *Client) {
		r.DisableTrace = disableTrace
	}
}

// WithAddr 设置ip:port
func WithAddr(addr []string) Option {
	return func(r *Client) {
		r.Addr = addr
	}
}

// WithDialTimeout 设置200ms
func WithDialTimeout(dialTimeout string) Option {
	return func(r *Client) {
		r.DialTimeout = dialTimeout
	}
}

// WithOpenStrategy 设置random/in_order (default random)
func WithOpenStrategy(openStrategy string) Option {
	return func(r *Client) {
		r.OpenStrategy = openStrategy
	}
}

// WithCompress 设置enable lz4 compression
func WithCompress(compress bool) Option {
	return func(r *Client) {
		r.Compress = compress
	}
}

// WithMaxExecutionTime 设置执行超时时间
func WithMaxExecutionTime(maxExecutionTime string) Option {
	return func(r *Client) {
		r.MaxExecutionTime = maxExecutionTime
	}
}

// WithDisableDebug 关闭 debug模式
func WithDisableDebug(disableDebug bool) Option {
	return func(r *Client) {
		r.DisableDebug = disableDebug
	}
}

func WithDisablePrepare(disablePrepare bool) Option {
	return func(r *Client) {
		r.DisablePrepare = disablePrepare
	}
}
