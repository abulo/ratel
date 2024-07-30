package driver

import (
	"time"
)

const (
	OptionTypeTimeout = 0x600
)

type Option interface {
	Type() int
}

type TimeoutOption struct{ timeout time.Duration }

func (to TimeoutOption) Type() int                         { return OptionTypeTimeout }
func NewTimeoutOption(timeout time.Duration) TimeoutOption { return TimeoutOption{timeout: timeout} }
