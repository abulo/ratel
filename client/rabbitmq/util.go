package rabbitmq

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v2/core/logger"
	"github.com/streadway/amqp"
)

const defaultRetryInterval = time.Second * 3
const defaultRetryTimes = 10

func init() {
	rand.Seed(time.Now().Unix())
}

// Dial 如果 retryable 为 nil，则表示不启用断线重连
func Dial(url string, retryable Retryable) (*Connection, error) {
	conn := NewConnection(url, retryable)
	err := conn.Dial()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Adder 累加器。每次执行累加一定数额，返回一个 uint64。
type Adder func() uint64

// SAdder 返回一个 uint64 的字符串
type SAdder func() string

// NewDefaultSAdder 获取一个从 0 开始累加，每次加 1 的累加器。返回一个 uint64 的字符串
func NewDefaultSAdder() SAdder {
	return SAdderGenerator(1)
}

// AdderGenerator 累加器生成器。生成的累加器从 0 开始累加，delta 表示需要累加的数字
func AdderGenerator(delta uint64) Adder {
	var i uint64 = 0
	return func() uint64 { return atomic.AddUint64(&i, delta) }
}

// SAdderGenerator ..
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

// Retryable ..
type Retryable interface {
	// retry 会尝试重试 retryOperation 中的操作。retryOperation 返回 brk 表示终止循环；
	// 否则继续尝试，直到尝试次数结束。
	retry(retryOperation func() (brk bool))
	// 是否已放弃重试（即，达到了重试结束的标志）
	hasGaveUp() bool
	// 放弃重试。在应该放弃重试的时候主动放弃重试，防止多余的重试或无限重试。
	GiveUp()
}

var emptyRetryable emptyRetry

type emptyRetry struct{}

func (emptyRetry) retry(fn func() bool) { fn() } // 只执行一次 listener
func (emptyRetry) hasGaveUp() bool      { return true }
func (emptyRetry) GiveUp()              {}

// TimesRetry ..
type TimesRetry struct {
	RetryTimes int           // 重试次数。如果 Always 为 true，此选项不可用。
	Interval   time.Duration // 间隔时间，指定断线后间隔多久再尝试重试。
	Always     bool          // 是否一直重试
	gaveUp     bool          // 是否已放弃重试
	sync.RWMutex
	//fastRetry  bool            //是否启用使用快速重试。只在重试方法是 Always 时可用。表示断线后是否先尝试快速重试，然后再尝试间隔指定时间发起连接
}

// NewTimesRetry 创建根据次数结束重试的配置
func NewTimesRetry(always bool, interval time.Duration, retryTimes int) *TimesRetry {
	return &TimesRetry{Always: always, Interval: interval, RetryTimes: retryTimes}
}

// DefaultTimesRetry 创建一个默认的重试配置：总是重试，且间隔三秒
func DefaultTimesRetry() *TimesRetry {
	return &TimesRetry{Always: true, Interval: defaultRetryInterval, RetryTimes: defaultRetryTimes}
}

// retry 见 Retryable.retry()
func (r *TimesRetry) retry(retryOperation func() (brk bool)) {
	var brk bool
	// 超出指定连接次数或连接成功，则退出循环
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

// hasGaveUp 见 Retryable.hasGaveUp()
func (r *TimesRetry) hasGaveUp() bool {
	r.RLock()
	defer r.RUnlock()
	return r.gaveUp
}

// GiveUp ..
func (r *TimesRetry) GiveUp() {
	r.Lock()
	r.gaveUp = true
	r.Unlock()
}

// TimesRetryBuilder ..
type TimesRetryBuilder struct {
	timesRetry *TimesRetry
}

// NewTimesRetryBuilder ..
func NewTimesRetryBuilder() *TimesRetryBuilder {
	return &TimesRetryBuilder{DefaultTimesRetry()}
}

// SetRetryTimes ...
func (bld *TimesRetryBuilder) SetRetryTimes(retryTimes int) *TimesRetryBuilder {
	bld.timesRetry.RetryTimes = retryTimes
	return bld
}

// SetInterval ...
func (bld *TimesRetryBuilder) SetInterval(interval time.Duration) *TimesRetryBuilder {
	bld.timesRetry.Interval = interval
	return bld
}

// SetAlways ...
func (bld *TimesRetryBuilder) SetAlways(always bool) *TimesRetryBuilder {
	bld.timesRetry.Always = always
	return bld
}

// Builder ...
func (bld *TimesRetryBuilder) Builder() *TimesRetry {
	return bld.timesRetry
}

// CtxRetry ...
type CtxRetry struct {
	Ctx      context.Context
	Interval time.Duration // 间隔时间，指定断线后间隔多久再尝试重试。
	gaveUp   bool          // 是否已放弃重试
	sync.RWMutex
}

// NewCtxRetry ...
func NewCtxRetry(ctx context.Context, interval time.Duration) *CtxRetry {
	return &CtxRetry{Ctx: ctx, Interval: interval}
}

// DefaultCtxRetry ...
func DefaultCtxRetry(ctx context.Context) *CtxRetry {
	return &CtxRetry{Ctx: ctx, Interval: defaultRetryInterval}
}

// retry 见 Retryable.retry()
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

// hasGaveUp 见 Retryable.hasGaveUp()
func (r *CtxRetry) hasGaveUp() bool {
	var gaveUp bool
	r.RLock()
	gaveUp = r.gaveUp
	r.RUnlock()
	return gaveUp || r.Ctx.Err() != nil
}

// GiveUp ..
func (r *CtxRetry) GiveUp() {
	r.Lock()
	defer r.Unlock()
	r.gaveUp = true
}

var (
	// MessagePlainTransient 无格式、非持久化消息工厂方法
	MessagePlainTransient MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         body,
		}
	}
	// MessagePlainPersistent 无格式、持久化消息工厂方法
	MessagePlainPersistent MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		}
	}
	// MessageJSONTransient JSON、非持久化消息工厂方法
	MessageJSONTransient MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/json",
			DeliveryMode: amqp.Transient,
			Body:         body,
		}
	}
	// MessageJSONPersistent JSON、持久化消息工厂方法
	MessageJSONPersistent MessageFactory = func(body []byte) amqp.Publishing {
		return amqp.Publishing{
			ContentType:  "text/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		}
	}
)

// MessageFactory 消息工厂方法。默认提供了如： MessagePlainTransient, MessagePlainPersistent, MessageJSONPersistent 等
// 在内的工厂方法。
// 如果没有需要的工厂方法，则需要调用者自己提供对应的工厂方法。
type MessageFactory func(body []byte) amqp.Publishing

// ReceiveListener ...
type ReceiveListener interface {
	// Consumer 用于实现消费操作。详见 ConsumerFunc。
	//
	// 如果消费者主动终止了退出，应该在 Finish 中主动移除当前 ReceiveListener。
	// 否则下次断线重连会再次触发该消息的接收操作。
	Consumer(*amqp.Delivery) (brk bool)
	// Finish 处理接收完成需要执行的操作，比如用于清理或关闭某些内容。
	// 如果没有相关操作需要执行，可以提供空实现。
	// Finish 可能由于主动停止接收消息或因为产生错误被调用。
	// 如果消费时没有错误，则参数 err 为 nil。
	Finish(err error)
	// 当主动停止消费时，应当实现该方法，主动移除当前 ReceiveListener。
	// 否则，一旦断线重连，该 ReceiveListener 会继续消费。
	Remove(key string, ch *Channel)
}

// AbsReceiveListener 的抽象实现。
// 如果 ConsumerMethod 为 nil 或不赋值，将 panic;
// 如果 FinishMethod 为 nil 或不赋值，则默认不做任何操作。
type AbsReceiveListener struct {
	ConsumerMethod ConsumerFunc
	FinishMethod   func(err error)
}

// Consumer ..
func (lis *AbsReceiveListener) Consumer(delivery *amqp.Delivery) (brk bool) {
	if lis.ConsumerMethod == nil {
		panic("AbsReceiveListener.ConsumerMethod must not be nil")
	}
	return lis.ConsumerMethod(delivery)
}

// Finish TODO
func (lis *AbsReceiveListener) Finish(err error) {
	if lis.FinishMethod == nil {
		return
	}
	lis.FinishMethod(err)
}

// Remove TODO
func (lis *AbsReceiveListener) Remove(key string, ch *Channel) {
	ch.RemoveOperation(key)
}
