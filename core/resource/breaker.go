package resource

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	Bucket struct {
		// Sum 为桶内所有事件数值之和。
		Sum float64
		// Count 为桶内事件次数。
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
	}
	window struct {
		buckets []*Bucket
		size    int
	}

	// Proba 在受保护的随机数源上按给定概率进行随机判定，方法为并发安全。
	Proba struct {
		// rand.Rand 本身不是并发安全的，因此通过锁串行化随机数访问。
		r    *rand.Rand
		lock sync.Mutex
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
	return &Proba{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TrueOnProba 以 proba 指定的阈值概率返回 true。
func (p *Proba) TrueOnProba(proba float64) (truth bool) {
	p.lock.Lock()
	truth = p.r.Float64() < proba
	p.lock.Unlock()
	return
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
func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
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
	offset := int(Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}

	return rw.size
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	offset := rw.offset
	for i := range span {
		rw.win.resetBucket((offset + i + 1) % rw.size)
	}

	rw.offset = (offset + span) % rw.size
	now := Now()
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

func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := range count {
		fn(w.buckets[(start+i)%w.size])
	}
}

func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
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
