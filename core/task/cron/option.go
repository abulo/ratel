package cron

import (
	"time"
)

// Option 表示对 Cron 默认行为的一项修改，通过 New 传入。
type Option func(*Cron)

// WithLocation 覆盖 Cron 实例解析调度表时使用的时区，默认为机器本地时区。
func WithLocation(loc *time.Location) Option {
	return func(c *Cron) {
		c.location = loc
	}
}

// WithSeconds 将调度表解析器切换为首位包含秒字段的格式。
func WithSeconds() Option {
	return WithParser(NewParser(
		Second | Minute | Hour | Dom | Month | Dow | Descriptor,
	))
}

// WithParser 覆盖用于解析调度表字符串的解析器。
func WithParser(p ScheduleParser) Option {
	return func(c *Cron) {
		c.parser = p
	}
}

// WithChain 指定应用到本 Cron 所有任务上的包装器链。
func WithChain(wrappers ...JobWrapper) Option {
	return func(c *Cron) {
		c.chain = NewChain(wrappers...)
	}
}
