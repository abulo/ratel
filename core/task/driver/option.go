package driver

import (
	"time"
)

const (
	// OptionTypeTimeout 为超时配置选项的类型标识。
	OptionTypeTimeout = 0x600
)

// Option 是驱动配置选项接口，实现方通过 Type 返回自身选项种类。
type Option interface {
	Type() int
}

// TimeoutOption 设置驱动的心跳与过期判定超时时长。
type TimeoutOption struct{ timeout time.Duration }

// Type 返回选项类型标识 OptionTypeTimeout。
func (to TimeoutOption) Type() int { return OptionTypeTimeout }

// NewTimeoutOption 创建指定超时时长的 TimeoutOption。
func NewTimeoutOption(timeout time.Duration) TimeoutOption { return TimeoutOption{timeout: timeout} }
