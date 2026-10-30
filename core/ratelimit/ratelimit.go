package ratelimit

import (
	"math"
	"strconv"
	"sync"
	"time"
)

// Bucket 是按预定速率填充令牌的令牌桶；由 NewBucket 等构造函数创建，零值不可直接使用。
//
// Take、Wait 等取用操作可被多个 goroutine 并发调用；桶初始为满桶状态。
type Bucket struct {
	clock Clock

	startTime time.Time

	capacity int64

	quantum int64

	fillInterval time.Duration

	mu sync.Mutex

	availableTokens int64

	latestTick int64
}

// NewBucket 创建一个令牌桶：每经过 fillInterval 补充 1 个令牌，令牌总量上限为 capacity；
// 两个参数都必须为正数，桶创建时即为满桶。
func NewBucket(fillInterval time.Duration, capacity int64) *Bucket {
	return NewBucketWithClock(fillInterval, capacity, nil)
}

// NewBucketWithClock 与 NewBucket 行为一致，但允许注入可替换的 Clock 时间源；
// clock 为 nil 时使用系统时钟。
func NewBucketWithClock(fillInterval time.Duration, capacity int64, clock Clock) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, 1, clock)
}

// rateMargin 为允许的实际填充速率与指定速率之间的最大相对偏差，取 1%。
const rateMargin = 0.01

// NewBucketWithRate 创建一个按每秒 rate 个令牌的速率填充、容量上限为 capacity 的令牌桶。
//
// 受时钟分辨率限制，高速率下实际填充速率与指定速率的偏差不超过 1%；rate 与 capacity 必须为正数。
func NewBucketWithRate(rate float64, capacity int64) *Bucket {
	return NewBucketWithRateAndClock(rate, capacity, nil)
}

// NewBucketWithRateAndClock 与 NewBucketWithRate 行为一致，但允许注入可替换的 Clock 时间源；
// clock 为 nil 时使用系统时钟。
func NewBucketWithRateAndClock(rate float64, capacity int64, clock Clock) *Bucket {
	// 循环复用同一个桶，避免反复分配。
	tb := NewBucketWithQuantumAndClock(1, capacity, 1, clock)
	for quantum := int64(1); quantum < 1<<50; quantum = nextQuantum(quantum) {
		fillInterval := time.Duration(1e9 * float64(quantum) / rate)
		if fillInterval <= 0 {
			continue
		}
		tb.fillInterval = fillInterval
		tb.quantum = quantum
		if diff := math.Abs(tb.Rate() - rate); diff/rate <= rateMargin {
			return tb
		}
	}
	panic("cannot find suitable quantum for " + strconv.FormatFloat(rate, 'g', -1, 64))
}

// nextQuantum 返回 q 之后下一个待尝试的量子值；量子按指数缓慢增长，以便在较小数值区间获得良好拟合。
func nextQuantum(q int64) int64 {
	q1 := q * 11 / 10
	if q1 == q {
		q1++
	}
	return q1
}

// NewBucketWithQuantum 与 NewBucket 类似，但每个 fillInterval 间隔补充 quantum 个令牌；
// fillInterval、capacity 与 quantum 都必须为正数。
func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, quantum, nil)
}

// NewBucketWithQuantumAndClock 是全参数构造函数：每经过 fillInterval 补充 quantum 个令牌，
// 容量上限为 capacity；clock 为 nil 时使用系统时钟。
//
// fillInterval、capacity 与 quantum 任一非正数时直接 panic。
func NewBucketWithQuantumAndClock(fillInterval time.Duration, capacity, quantum int64, clock Clock) *Bucket {
	if clock == nil {
		clock = realClock{}
	}
	if fillInterval <= 0 {
		panic("token bucket fill interval is not > 0")
	}
	if capacity <= 0 {
		panic("token bucket capacity is not > 0")
	}
	if quantum <= 0 {
		panic("token bucket quantum is not > 0")
	}
	return &Bucket{
		clock:           clock,
		startTime:       clock.Now(),
		latestTick:      0,
		fillInterval:    fillInterval,
		capacity:        capacity,
		quantum:         quantum,
		availableTokens: capacity,
	}
}

// Wait 从桶中取走 count 个令牌，令牌不足时阻塞睡眠到令牌可用为止。
func (tb *Bucket) Wait(count int64) {
	if d := tb.Take(count); d > 0 {
		tb.clock.Sleep(d)
	}
}

// WaitMaxDuration 与 Wait 类似，但仅在等待时长不超过 maxWait 时才取走令牌；
// 若需要等待超过 maxWait，则不扣减令牌并立即返回 false，否则返回 true。
func (tb *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	d, ok := tb.TakeMaxDuration(count, maxWait)
	if d > 0 {
		tb.clock.Sleep(d)
	}
	return ok
}

const infinityDuration time.Duration = 0x7fffffffffffffff

// Take 非阻塞地从桶中取走 count 个令牌，返回调用方应等待到令牌实际可用的时长。
//
// 取用一旦提交即不可撤销：方法返回后无法再把令牌归还桶中；令牌充足时返回 0。
func (tb *Bucket) Take(count int64) time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	d, _ := tb.take(tb.clock.Now(), count, infinityDuration)
	return d
}

// TakeMaxDuration 与 Take 类似，但仅在令牌等待时长不超过 maxWait 时才扣减令牌；
// 等待超过 maxWait 时不做任何扣减并返回 false，否则返回应等待时长与 true。
func (tb *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.take(tb.clock.Now(), count, maxWait)
}

// TakeAvailable 立即取走桶中至多 count 个现有令牌，不阻塞；
// 返回实际取走的令牌数，桶中没有可用令牌时返回 0。
func (tb *Bucket) TakeAvailable(count int64) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.takeAvailable(tb.clock.Now(), count)
}

func (tb *Bucket) takeAvailable(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}
	tb.adjustavailableTokens(tb.currentTick(now))
	if tb.availableTokens <= 0 {
		return 0
	}
	if count > tb.availableTokens {
		count = tb.availableTokens
	}
	tb.availableTokens -= count
	return count
}

// Available 返回当前可用令牌数；存在等待取用者时该值可能为负数。
//
// 返回值大于 0 并不保证后续取用一定成功，期间令牌数可能已被其他并发调用改变；
// 本方法主要用于指标上报与调试。
func (tb *Bucket) Available() int64 {
	return tb.available(tb.clock.Now())
}

func (tb *Bucket) available(now time.Time) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.adjustavailableTokens(tb.currentTick(now))
	return tb.availableTokens
}

// Capacity 返回桶创建时设定的容量上限。
func (tb *Bucket) Capacity() int64 {
	return tb.capacity
}

// Rate 返回桶的填充速率，单位为令牌每秒。
func (tb *Bucket) Rate() float64 {
	return 1e9 * float64(tb.quantum) / float64(tb.fillInterval)
}

func (tb *Bucket) take(now time.Time, count int64, maxWait time.Duration) (time.Duration, bool) {
	if count <= 0 {
		return 0, true
	}

	tick := tb.currentTick(now)
	tb.adjustavailableTokens(tick)
	avail := tb.availableTokens - count
	if avail >= 0 {
		tb.availableTokens = avail
		return 0, true
	}
	// 将缺口令牌向上取整到 quantum 的整数倍；endTick 记录请求的全部令牌
	// 补齐所需到达的刻度，令牌在下一个对应刻度才会补齐。
	endTick := tick + (-avail+tb.quantum-1)/tb.quantum
	endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
	waitTime := endTime.Sub(now)
	if waitTime > maxWait {
		return 0, false
	}
	tb.availableTokens = avail
	return waitTime, true
}

func (tb *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(tb.startTime) / tb.fillInterval)
}

func (tb *Bucket) adjustavailableTokens(tick int64) {
	lastTick := tb.latestTick
	tb.latestTick = tick
	if tb.availableTokens >= tb.capacity {
		return
	}
	tb.availableTokens += (tick - lastTick) * tb.quantum
	if tb.availableTokens > tb.capacity {
		tb.availableTokens = tb.capacity
	}
}

// Clock 抽象时间流逝过程，允许在测试中伪造时间推进；实现需提供获取当前时间与睡眠指定时长的能力。
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// realClock 基于标准库 time 实现 Clock。
type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
