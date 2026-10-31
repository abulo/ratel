package resource

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// letterIdxBits 为表示一个字母索引所需的位数。
	letterIdxBits  = 6
	idLen          = 8
	defaultRandLen = 8
	// letterIdxMask 为取出低位 letterIdxBits 位的全 1 掩码。
	letterIdxMask = 1<<letterIdxBits - 1
	// letterIdxMax 为一次 Int63 的 63 位随机数可容纳的索引个数。
	letterIdxMax      = 63 / letterIdxBits
	numHistoryReasons = 5
	timeFormat        = "15:04:05"
	windowInit        = time.Second * 10
	buckets           = 40
	k                 = 1.5
	protection        = 5
)

// ErrServiceUnavailable 是熔断器处于拒绝状态时返回的错误，表示当前请求应快速失败。
var ErrServiceUnavailable = errors.New("circuit breaker is open")

var initTime = time.Now().AddDate(-1, -1, -1)

type (
	// Acceptable 判定一次调用的返回错误是否视为成功：返回 true 计入成功统计，
	// 返回 false 计入失败统计。
	Acceptable func(err error) bool

	// Breaker 是熔断器抽象，在错误率过高或系统过载时拒绝请求，避免压垮下游服务。
	//
	// 调用方既可以通过 Allow 手动管理请求生命周期，也可以使用 Do 系列方法直接包装
	// 业务函数；业务函数内部发生 panic 时，熔断器按失败处理并将 panic 原样重新抛出。
	Breaker interface {
		// Name 返回熔断器的名称。
		Name() string

		// Allow 在允许放行时返回一个 Promise，调用成功后须调用 Accept，失败后须调用
		// Reject；熔断器拒绝放行时返回 ErrServiceUnavailable。
		Allow() (Promise, error)

		// Do 在熔断器放行时执行 req，被拒绝时立即返回错误；req 内的 panic 会被记为
		// 失败后重新抛出。
		Do(req func() error) error

		// DoWithAcceptable 行为同 Do，acceptable 用于自定义错误是否视为成功。
		DoWithAcceptable(req func() error, acceptable Acceptable) error

		// DoWithFallback 在熔断器放行时执行 req，被拒绝时执行 fallback 并返回其错误。
		DoWithFallback(req func() error, fallback func(err error) error) error

		// DoWithFallbackAcceptable 组合 fallback 降级与 acceptable 成功判定。
		DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}

	circuitBreaker struct {
		name string
		throttle
	}

	throttle interface {
		allow() (Promise, error)
		doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}

	// Promise 表示一次已放行请求的结果回调，由 Breaker.Allow 返回，调用方须且仅调用
	// Accept 与 Reject 之一。
	Promise interface {
		// Accept 通知熔断器本次调用成功。
		Accept()
		// Reject 通知熔断器本次调用失败，reason 记录失败原因，供最近历史错误排查使用。
		Reject(reason string)
	}

	// Option 以函数式选项模式定制 NewBreaker 创建的熔断器。
	Option func(breaker *circuitBreaker)

	lockedSource struct {
		source rand.Source
		lock   sync.Mutex
	}

	loggedThrottle struct {
		name string
		internalThrottle
		errWin *errorWindow
	}
	internalThrottle interface {
		allow() (internalPromise, error)
		doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}

	internalPromise interface {
		Accept()
		Reject()
	}

	errorWindow struct {
		reasons [numHistoryReasons]string
		index   int
		count   int
		lock    sync.Mutex
	}
	promiseWithReason struct {
		promise internalPromise
		errWin  *errorWindow
	}

	// RollingWindowOption 以函数式选项模式定制 NewRollingWindow 创建的滚动窗口。
	RollingWindowOption func(rollingWindow *RollingWindow)

	// Bucket 为滚动窗口中的一个时间桶，累计一段时间内事件的数值与次数。
	//
	// PERF-4-1 [P1]：桶内 Sum/Count 改为原子存储，同桶累加无需排他写锁；
	// 导出字段 Sum/Count 为 Reduce 回调读取用的快照，由 Reduce 在锁内同步。
	Bucket struct {
		sumBits atomic.Uint64
		count   atomic.Int64

		// Sum 为桶内所有事件数值之和（快照字段，供 Reduce 回调读取）。
		Sum float64
		// Count 为桶内事件次数（快照字段，供 Reduce 回调读取）。
		Count int64
	}

	// RollingWindow 是按时间片滚动的统计窗口，由固定数量的桶组成，周期性覆盖已过期的桶。
	//
	// 窗口方法为并发安全；零值不可直接使用，须通过 NewRollingWindow 创建。
	RollingWindow struct {
		lock          sync.RWMutex
		size          int
		win           *window
		interval      time.Duration
		offset        int
		ignoreCurrent bool
		lastTime      time.Duration

		// PERF-4-2 [P1]：窗口累计 Sum/Count，熔断判定直接读取，O(1) 替代全桶遍历；
		// 与各桶原子字段保持一致，跨桶复位时同步扣减（40 桶遍历占比预期收益需实测验证）。
		cumSumBits atomic.Uint64
		cumCount   atomic.Int64
	}
	window struct {
		buckets []*Bucket
		size    int
	}

	// Proba 在受保护的随机数源上按给定概率进行随机判定，方法为并发安全。
	//
	// PERF-4-3 [P2]：改用无锁 xorshift64 状态（atomic CAS），替代互斥锁保护的
	// 非并发安全 *rand.Rand（预期收益需实测验证）。
	Proba struct {
		state atomic.Uint64
	}

	googleBreaker struct {
		k     float64
		stat  *RollingWindow
		proba *Proba
	}

	googlePromise struct {
		b *googleBreaker
	}
)

// NewBreaker 创建一个熔断器；未通过 WithName 指定名称时自动生成随机名称。
//
// 返回的熔断器可并发使用，opts 用于覆盖默认配置。
func NewBreaker(opts ...Option) Breaker {
	var b circuitBreaker
	for _, opt := range opts {
		opt(&b)
	}
	if len(b.name) == 0 {
		b.name = Rand()
	}
	b.throttle = newLoggedThrottle(b.name, newGoogleBreaker())

	return &b
}

func (cb *circuitBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *circuitBreaker) Do(req func() error) error {
	return cb.throttle.doReq(req, nil, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *circuitBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.throttle.doReq(req, fallback, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

// WithName 设置熔断器名称，便于日志与熔断记录中区分不同下游。
func WithName(name string) Option {
	return func(b *circuitBreaker) {
		b.name = name
	}
}

func defaultAcceptable(err error) bool {
	return err == nil
}

func (lt loggedThrottle) allow() (Promise, error) {
	promise, err := lt.internalThrottle.allow()
	return promiseWithReason{
		promise: promise,
		errWin:  lt.errWin,
	}, lt.logError(err)
}

func (lt loggedThrottle) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return lt.logError(lt.internalThrottle.doReq(req, fallback, func(err error) bool {
		accept := acceptable(err)
		if !accept && err != nil {
			lt.errWin.add(err.Error())
		}
		return accept
	}))
}

// logError 为熔断器错误日志的预留挂接点，当前原样返回错误，不做额外处理。
func (lt loggedThrottle) logError(err error) error {
	return err
}

// Pid 返回当前进程的操作系统进程 ID。
func Pid() int {
	return os.Getpid()
}

// ProcessName 返回当前进程名，即可执行文件的命令名。
func ProcessName() string {
	return filepath.Base(os.Args[0])
}

// Rand 返回长度为默认 8 的随机字母数字字符串。
func Rand() string {
	return Randn(defaultRandLen)
}

// Randn 返回长度为 n 的随机字母数字字符串；n 为 0 时返回空字符串。
func Randn(n int) string {
	b := make([]byte, n)
	// 一次 Int63 产生 63 个随机位，可填充 letterIdxMax 个字符。
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var src = newLockedSource(time.Now().UnixNano())

func newLockedSource(seed int64) *lockedSource {
	return &lockedSource{
		source: rand.NewSource(seed),
	}
}

func (ls *lockedSource) Int63() int64 {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	return ls.source.Int63()
}

func newLoggedThrottle(name string, t internalThrottle) loggedThrottle {
	return loggedThrottle{
		name:             name,
		internalThrottle: t,
		errWin:           new(errorWindow),
	}
}

func (ew *errorWindow) add(reason string) {
	ew.lock.Lock()
	ew.reasons[ew.index] = fmt.Sprintf("%s %s", time.Now().Format(timeFormat), reason)
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = MinInt(ew.count+1, numHistoryReasons)
	ew.lock.Unlock()
}

func (ew *errorWindow) String() string {
	var reasons []string

	ew.lock.Lock()
	// 按最近发生的顺序反向输出历史原因。
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		reasons = append(reasons, ew.reasons[(i+numHistoryReasons)%numHistoryReasons])
	}
	ew.lock.Unlock()

	return strings.Join(reasons, "\n")
}

func (p promiseWithReason) Accept() {
	p.promise.Accept()
}

func (p promiseWithReason) Reject(reason string) {
	p.errWin.add(reason)
	p.promise.Reject()
}

// MinInt 返回 a、b 中较小的整数。
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// NewProba 创建一个以当前时间纳秒为种子的 Proba。
func NewProba() *Proba {
	p := &Proba{}
	// PERF-4-3 [P2]：种子避免落到 xorshift 的全零退化状态。
	seed := uint64(time.Now().UnixNano())
	if seed == 0 {
		seed = 1
	}
	p.state.Store(seed)
	return p
}

// TrueOnProba 以 proba 指定的阈值概率返回 true。
func (p *Proba) TrueOnProba(proba float64) (truth bool) {
	// PERF-4-3 [P2]：CAS 循环实现无锁 xorshift64，去掉互斥锁（预期收益需实测验证）。
	return p.nextFloat64() < proba
}

// nextFloat64 通过无锁 xorshift64 生成 [0,1) 区间的均匀随机数。
func (p *Proba) nextFloat64() float64 {
	for {
		old := p.state.Load()
		x := old
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
		if x == 0 {
			x = 1
		}
		if p.state.CompareAndSwap(old, x) {
			// 取高 53 位映射到 [0,1)，与 rand.Float64 精度一致。
			return float64(x>>11) / (1 << 53)
		}
	}
}
func newGoogleBreaker() *googleBreaker {
	bucketDuration := time.Duration(int64(windowInit) / int64(buckets))
	st := NewRollingWindow(buckets, bucketDuration)
	return &googleBreaker{
		stat:  st,
		k:     k,
		proba: NewProba(),
	}
}

func (b *googleBreaker) accept() error {
	accepts, total := b.history()
	weightedAccepts := b.k * float64(accepts)
	// 拒绝比例公式参考 Google SRE《Handling Overload》章节：
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	dropRatio := math.Max(0, (float64(total-protection)-weightedAccepts)/float64(total+1))
	if dropRatio <= 0 {
		return nil
	}

	if b.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}

	return nil
}

func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		return nil, err
	}

	return googlePromise{
		b: b,
	}, nil
}

func (b *googleBreaker) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		if fallback != nil {
			return fallback(err)
		}

		return err
	}

	defer func() {
		if e := recover(); e != nil {
			b.markFailure()
			panic(e)
		}
	}()

	err := req()
	if acceptable(err) {
		b.markSuccess()
	} else {
		b.markFailure()
	}

	return err
}

func (b *googleBreaker) markSuccess() {
	b.stat.Add(1)
}

func (b *googleBreaker) markFailure() {
	b.stat.Add(0)
}

func (b *googleBreaker) history() (accepts, total int64) {
	// PERF-4-2 [P1]：默认路径直接读取窗口累计计数，O(1) 替代 Reduce 全桶遍历；
	// ignoreCurrent 场景语义特殊（跳过当前桶），保留原有遍历实现（预期收益需实测验证）。
	if !b.stat.ignoreCurrent {
		total = b.stat.cumCount.Load()
		accepts = int64(math.Float64frombits(b.stat.cumSumBits.Load()))
		return
	}
	b.stat.Reduce(func(b *Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})

	return
}

// NewRollingWindow 创建包含 size 个桶、桶时间间隔为 interval 的滚动窗口。
//
// size 必须大于 0，否则直接 panic；opts 用于覆盖窗口默认行为。
func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Add 将数值 v 计入当前时间桶。
//
// PERF-4-1 [P1]：同桶累加走原子操作、不加排他写锁；仅跨桶轮转时在短临界区内
// 复位过期桶；累计计数同步原子维护（锁竞争占比预期收益需实测验证）。
func (rw *RollingWindow) Add(v float64) {
	rw.lock.RLock()
	rotated := rw.span() > 0
	rw.lock.RUnlock()

	if rotated {
		rw.lock.Lock()
		// 双检：等待写锁期间可能已被其他协程完成轮转。
		if rw.span() > 0 {
			rw.updateOffset(Now())
		}
		rw.lock.Unlock()
	}

	rw.lock.RLock()
	idx := rw.offset % rw.size
	rw.lock.RUnlock()
	b := rw.win.buckets[idx]

	b.add(v)
	rw.cumCount.Add(1)
	if v != 0 {
		addFloatBits(&rw.cumSumBits, v)
	}
}

// Reduce 遍历全部历史时间桶并对每个桶执行 fn；配置忽略当前桶时，尚未结束的当前桶
// 不计入统计。
func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// 当前桶数据尚未写完，配置 ignoreCurrent 时跳过，避免被部分数据污染统计。
	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		rw.win.reduce(offset, diff, fn)
	}
}

func (rw *RollingWindow) span() int {
	return rw.spanAt(Now())
}

// spanAt 以给定当前时刻计算距基准时间已经越过的桶偏移；PERF-4-4 [P2]：供 Add 复用
// 同一时刻，避免临界区内重复取系统时钟。
func (rw *RollingWindow) spanAt(now time.Duration) int {
	offset := int((now - rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}

	return rw.size
}

// updateOffset 复位已过期的桶并推进偏移；调用方必须持有写锁。
//
// PERF-4-4 [P2]：now 由调用方传入，span 判定与时间对齐复用同一时刻。
func (rw *RollingWindow) updateOffset(now time.Duration) {
	span := rw.spanAt(now)
	if span <= 0 {
		return
	}

	offset := rw.offset
	for i := range span {
		b := rw.win.buckets[(offset+i+1)%rw.size]
		// PERF-4-2 [P1]：复位旧桶前从累计值扣除其残留，保证累计与可见桶之和一致。
		if c := b.count.Load(); c != 0 {
			rw.cumCount.Add(-c)
		}
		if bits := b.sumBits.Load(); bits != 0 {
			addFloatBits(&rw.cumSumBits, -math.Float64frombits(bits))
		}
		b.reset()
	}

	rw.offset = (offset + span) % rw.size
	// 将基准时间对齐到桶间隔边界，保证下一轮窗口统计从整段间隔开始。
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}

func newWindow(size int) *window {
	buckets := make([]*Bucket, size)
	for i := range size {
		buckets[i] = new(Bucket)
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := range count {
		b := w.buckets[(start+i)%w.size]
		// Reduce 回调读取前，把原子字段同步到导出快照字段。
		b.snapshot()
		fn(b)
	}
}

// add 以原子方式累加一次事件：次数恒加 1，数值 v 累加到 Sum。
func (b *Bucket) add(v float64) {
	b.count.Add(1)
	if v != 0 {
		addFloatBits(&b.sumBits, v)
	}
}

func (b *Bucket) reset() {
	b.count.Store(0)
	b.sumBits.Store(0)
}

// snapshot 在 Reduce 锁内调用，把原子内部字段复制到导出快照字段供回调读取。
func (b *Bucket) snapshot() {
	b.Sum = math.Float64frombits(b.sumBits.Load())
	b.Count = b.count.Load()
}

// addFloatBits 对以 math.Float64bits 存储的浮点累计值做无锁原子加法。
func addFloatBits(acc *atomic.Uint64, v float64) {
	for {
		old := acc.Load()
		next := math.Float64bits(math.Float64frombits(old) + v)
		if acc.CompareAndSwap(old, next) {
			return
		}
	}
}

func (p googlePromise) Accept() {
	p.b.markSuccess()
}

func (p googlePromise) Reject() {
	p.b.markFailure()
}

// Now 返回自程序启动基线时刻起经过的时长，供滚动窗口计算时间偏移使用。
func Now() time.Duration {
	return time.Since(initTime)
}

// Since 返回自 d（同为 Now 基线时钟下的时长）起经过的时长。
func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}
