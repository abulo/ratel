package task

import (
	"time"

	"github.com/abulo/ratel/v3/core/task/cron"
)

type Option func(*Task)

// WithLocation is wrap cron.Cron with location
func WithLocation(loc *time.Location) Option {
	return func(c *Task) { c.opts = append(c.opts, cron.WithLocation(loc)) }
}

// WithSeconds is wrap cron.Cron with seconds
func WithSeconds() Option {
	return func(c *Task) { c.opts = append(c.opts, cron.WithSeconds()) }
}

// WithParser is wrap cron.Cron with schedules
func WithParser(p cron.ScheduleParser) Option {
	return func(c *Task) { c.opts = append(c.opts, cron.WithParser(p)) }
}

// WithChain is wrap cron.Cron with chains
func WithChain(wrappers ...cron.JobWrapper) Option {
	return func(c *Task) { c.opts = append(c.opts, cron.WithChain(wrappers...)) }
}

// WithLogger is wrap cron.Cron with logger
func WithLogger(logger cron.Logger) Option {
	return func(c *Task) { c.opts = append(c.opts, cron.WithLogger(logger)) }
}

// WithNodeUpdateInterval set node update interval
func WithNodeUpdateInterval(dur time.Duration) Option {
	return func(c *Task) {
		c.updateInterval = dur
	}
}

// WithLazyPick set lazy pick option
func WithLazyPick(lazy bool) Option {
	return func(c *Task) {
		c.lazyPick = lazy
	}
}
