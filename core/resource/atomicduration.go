package resource

import (
	"sync/atomic"
	"time"
)

// AtomicDuration 是以原子操作读写的 time.Duration，适用于多 goroutine 共享的
// 超时阈值等配置。零值表示时长 0，指针接收者方法为并发安全。
type AtomicDuration int64

// NewAtomicDuration 创建一个时长为 0 的 AtomicDuration。
func NewAtomicDuration() *AtomicDuration {
	return new(AtomicDuration)
}

// ForAtomicDuration 以指定初始值 val 创建一个 AtomicDuration。
func ForAtomicDuration(val time.Duration) *AtomicDuration {
	d := NewAtomicDuration()
	d.Set(val)
	return d
}

// CompareAndSwap 在当前值等于 old 时将其原子替换为 val，返回是否交换成功。
func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(old), int64(val))
}

// Load 原子读取当前时长。
func (d *AtomicDuration) Load() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

// Set 将当前时长原子写入为 val。
func (d *AtomicDuration) Set(val time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(val))
}
