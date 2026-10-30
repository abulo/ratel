package rabbitmq

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/streadway/amqp"
)

const defaultRetryInterval = time.Second * 3
const defaultRetryTimes = 10

// Dial 创建 Connection 并立即完成首次拨号连接；retryable 为 nil 时不启用断线重连。
// 拨号失败时返回错误。
func Dial(url string, retryable Retryable) (*Connection, error) {
	conn := NewConnection(url, retryable)
	err := conn.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Adder 是无符号整数累加器，每次调用返回累加后的新值。
type Adder func() uint64

// SAdder 是字符串形式的累加器，每次调用返回一个唯一标识字符串。
type SAdder func() string

// NewDefaultSAdder 返回从 0 开始、每次递增 1 的字符串累加器。
func NewDefaultSAdder() SAdder {
	return SAdderGenerator(1)
}

// AdderGenerator 生成一个从 0 开始、每次调用按 delta 递增的 Adder。
func AdderGenerator(delta uint64) Adder {
	var i uint64 = 0
	return func() uint64 { return atomic.AddUint64(&i, delta) }
}

// SAdderGenerator 生成一个字符串累加器，每次调用将内部计数器按 delta 递增并转为十六进制字符串返回。
func SAdderGenerator(delta uint64) SAdder {
	adder := AdderGenerator(delta)
	return func() string {
		return strconv.FormatUint(adder(), 16)
	}
}

func getNonNilArgs(args *amqp.Table) *amqp.Table {
	if args == nil {
		return &amqp.Table{}
	}
	return args
}

func getNonNilRetryable(retryable Retryable) Retryable {
	if retryable != nil {
		return retryable
	}
	return emptyRetryable
}

func getNonNilMessageFactory(factory MessageFactory) MessageFactory {
	if factory != nil {
		return factory
	}
	return MessagePlainPersistent
}

// Retryable 定义重试策略的契约：执行可重试操作、查询是否已放弃重试以及主动放弃重试。
// 实现须保证重试过程可被终止，避免出现无限重试。
type Retryable interface {
	// retry 反复执行 retryOperation；回调返回 brk 为 true 表示终止循环，否则按策略继续尝试，
	// 直到次数耗尽或策略主动放弃。
	retry(retryOperation func() (brk bool))
	// hasGaveUp 返回是否已经放弃重试（达到重试结束条件）。
	hasGaveUp() bool
	// GiveUp 主动放弃重试，用于在应当停止时立即终止后续重试，避免多余尝试或无限重试。
	GiveUp()
}

var emptyRetryable emptyRetry

type emptyRetry struct{}

// retry 仅执行一次操作，不做任何重试。
func (emptyRetry) retry(fn func() bool) { fn() }
func (emptyRetry) hasGaveUp() bool      { return true }
func (emptyRetry) GiveUp()              {}

// TimesRetry 按次数或无限循环执行重试，使用读写锁保护放弃标记，可并发查询 hasGaveUp。
type TimesRetry struct {
	// RetryTimes 为最大重试次数；Always 为 true 时该字段不生效。
	RetryTimes int
	// Interval 为两次重试之间的等待间隔。
	Interval time.Duration
	// Always 指定是否一直重试，不受 RetryTimes 限制。
	Always bool
	// gaveUp 标记是否已放弃重试。
	gaveUp bool
	sync.RWMutex
}

// NewTimesRetry 创建按次数（或无限）重试的策略；always 为 true 时忽略 retryTimes。
func NewTimesRetry(always bool, interval time.Duration, retryTimes int) *TimesRetry {
	return &TimesRetry{Always: always, Interval: interval, RetryTimes: retryTimes}
}

// DefaultTimesRetry 返回默认重试策略：Always 持续重试，间隔 3 秒；因 Always 为 true，
// RetryTimes 不生效。
func DefaultTimesRetry() *TimesRetry {
	return &TimesRetry{Always: true, Interval: defaultRetryInterval, RetryTimes: defaultRetryTimes}
}

// retry 按 RetryTimes 或 Always 规则反复执行操作，成功或达到次数上限后退出并标记放弃。
func (r *TimesRetry) retry(retryOperation func() (brk bool)) {
	var brk bool
	// 已放弃、非 Always 模式且次数耗尽时退出循环。
	retryTimes := r.RetryTimes
	for !r.hasGaveUp() || r.Always || retryTimes > 0 {
		if !r.Always {
			retryTimes--
		}
		brk = retryOperation()
		if brk {
			return
		}
		time.Sleep(r.Interval)
	}
	if !r.Always && retryTimes == 0 {
		r.GiveUp()
	}
}

// hasGaveUp 返回是否已放弃重试。
func (r *TimesRetry) hasGaveUp() bool {
	r.RLock()
	defer r.RUnlock()
	return r.gaveUp
}

// GiveUp 标记后续不再重试。
func (r *TimesRetry) GiveUp() {
	r.Lock()
	r.gaveUp = true
	r.Unlock()
}

// TimesRetryBuilder 以链式方式构造 TimesRetry。
type TimesRetryBuilder struct {
	timesRetry *TimesRetry
}

// NewTimesRetryBuilder 创建基于默认策略的 TimesRetryBuilder。
func NewTimesRetryBuilder() *TimesRetryBuilder {
	return &TimesRetryBuilder{DefaultTimesRetry()}
}

// SetRetryTimes 设置最大重试次数，返回构造器自身以便链式调用。
func (bld *TimesRetryBuilder) SetRetryTimes(retryTimes int) *TimesRetryBuilder {
	bld.timesRetry.RetryTimes = retryTimes
	return bld
}

// SetInterval 设置重试间隔，返回构造器自身以便链式调用。
func (bld *TimesRetryBuilder) SetInterval(interval time.Duration) *TimesRetryBuilder {
	bld.timesRetry.Interval = interval
	return bld
}

// SetAlways 设置是否持续重试，返回构造器自身以便链式调用。
func (bld *TimesRetryBuilder) SetAlways(always bool) *TimesRetryBuilder {
	bld.timesRetry.Always = always
	return bld
}

// Builder 完成策略构造，返回最终的 TimesRetry。
func (bld *TimesRetryBuilder) Builder() *TimesRetry {
	return bld.timesRetry
}

// CtxRetry 受上下文控制的重试策略：上下文结束后不再重试，使用读写锁保护放弃标记。
type CtxRetry struct {
	// Ctx 用于取消重试，上下文被取消或超时后 hasGaveUp 返回 true。
	Ctx context.Context
	// Interval 为两次重试之间的等待间隔。
	Interval time.Duration
	// gaveUp 标记是否已放弃重试。
	gaveUp bool
	sync.RWMutex
}

// NewCtxRetry 创建受指定上下文控制、固定间隔重试的策略。
func NewCtxRetry(ctx context.Context, interval time.Duration) *CtxRetry {
	return &CtxRetry{Ctx: ctx, Interval: interval}
}

// DefaultCtxRetry 创建受上下文控制、默认间隔 3 秒重试的策略。
func DefaultCtxRetry(ctx context.Context) *CtxRetry {
	return &CtxRetry{Ctx: ctx, Interval: defaultRetryInterval}
}

// retry 反复执行操作，直到操作成功、主动放弃或上下文结束。
func (r *CtxRetry) retry(retryOperation func() (brk bool)) {
	var brk bool
	for !r.hasGaveUp() {
		brk = retryOperation()
		if brk {
			return
		}
		time.Sleep(r.Interval)
	}
	logger.Logger.Debug("Gave up retrying or CtxRetry context done!")
}

// hasGaveUp 返回是否已放弃重试；上下文结束时同样返回 true。
func (r *CtxRetry) hasGaveUp() bool {
	var gaveUp bool
	r.RLock()
	gaveUp = r.gaveUp
	r.RUnlock()
	return gaveUp || r.Ctx.Err() != nil
}

// GiveUp 标记后续不再重试。
func (r *CtxRetry) GiveUp() {
	r.Lock()
	defer r.Unlock()
	r.gaveUp = true
}

var (
	// MessagePlainTransient 构造 text/plain、非持久化的消息。
	MessagePlainTransient MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         body,
		}
	}
	// MessagePlainPersistent 构造 text/plain、持久化的消息。
	MessagePlainPersistent MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		}
	}
	// MessageJSONTransient 构造 text/json、非持久化的消息。
	MessageJSONTransient MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/json",
			DeliveryMode: amqp.Transient,
			Body:         body,
		}
	}
	// MessageJSONPersistent 构造 text/json、持久化的消息。
	MessageJSONPersistent MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		}
	}
)

// MessageFactory 将消息体转换为 amqp.Publishing，决定消息的 ContentType 与持久化级别。
// 包内已提供 Plain/JSON、持久化/非持久化四个常用工厂；没有匹配需求时由调用方自行实现。
type MessageFactory func(body []byte) amqp.Publishing

// ReceiveListener 定义消费者监听器契约：负责消息消费、结束清理，以及在主动停止消费时从连接上
// 注销自身。消费者主动退出后必须在 Remove 中移除对应操作，否则断线重连会再次触发该消费监听。
type ReceiveListener interface {
	// Consumer 处理单条消息，返回 true 表示终止消费；详细语义见 ConsumerFunc。
	//
	// 消费者主动退出后应在 Finish/Remove 流程中注销当前监听器，避免重连后重复接收。
	Consumer(*amqp.Delivery) (brk bool)
	// Finish 在消费正常结束或发生错误时被调用，用于清理资源；无清理需求时可提供空实现。
	// 消费无错误时参数 err 为 nil。
	Finish(err error)
	// Remove 在主动停止消费时被调用，从连接上移除当前监听操作；断线重连后该监听器不应再被触发。
	Remove(key string, ch *Channel)
}

// AbsReceiveListener 为 ReceiveListener 提供基于方法字段的抽象实现，方便按消息注入消费与清理逻辑。
//
// ConsumerMethod 为 nil 时 Consumer 会 panic；FinishMethod 为 nil 时 Finish 不执行任何操作。
type AbsReceiveListener struct {
	// ConsumerMethod 为消费回调，为 nil 时 Consumer 会 panic。
	ConsumerMethod ConsumerFunc
	// FinishMethod 为结束清理回调，为 nil 时 Finish 直接返回。
	FinishMethod func(err error)
}

// Consumer 调用注入的 ConsumerMethod 处理消息；ConsumerMethod 未赋值时直接 panic。
func (lis *AbsReceiveListener) Consumer(delivery *amqp.Delivery) (brk bool) {
	if lis.ConsumerMethod == nil {
		panic("AbsReceiveListener.ConsumerMethod must not be nil")
	}
	return lis.ConsumerMethod(delivery)
}

// Finish 调用注入的 FinishMethod 执行清理；未注入时不做任何操作。
func (lis *AbsReceiveListener) Finish(err error) {
	if lis.FinishMethod == nil {
		return
	}
	lis.FinishMethod(err)
}

// Remove 从所属连接上移除当前监听操作。
func (lis *AbsReceiveListener) Remove(key string, ch *Channel) {
	ch.RemoveOperation(key)
}
