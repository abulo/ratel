package rabbitmq

import (
	"fmt"

	"github.com/abulo/ratel/v2/core/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ReceiveOpts 消息接收选项。
// 如果 autoAck 设为 false，RabbitMQ 会等待消费者显式地回复确认信号后才从内存(或者磁盘)中移去
// 消息。因此调用者应该主动调用 (*amqp.Delivery).Ack 确认消费，防止消息在内存(或者磁盘)中积累。
// 如果 autoAck 设为 true，消息会被服务器默认为已消费，可能会导致消费者无法处理数据时造成数据丢失。
// 如果 RabbitMQ 一直没有收到消费者的确认信号，并且消费此消息的消费者己经 断开连接，则 RabbitMQ
// 会安排该消息重新进入队列，等待投递给下一个消费者，当然也有可能还是原来的那个消费者。
// consumerTag 用于唯一识别一个消费者，如果不填可自动生成。
// 其他参数如果没有特别需求，默认不填即可。
type ReceiveOpts struct {
	autoAck, exclusive, noLocal, noWait bool
	args                                *amqp.Table
	consumerTag                         string
}

// DefaultReceiveOpts 将 ReceiveOpts.autoAck 默认设置为 true
func DefaultReceiveOpts() *ReceiveOpts {
	return &ReceiveOpts{autoAck: true}
}

type ReceiveOptsBuilder struct {
	opts *ReceiveOpts
}

func NewReceiveOptsBuilder() *ReceiveOptsBuilder {
	return &ReceiveOptsBuilder{DefaultReceiveOpts()}
}

// 作用详见 ReceiveOpts
func (bld *ReceiveOptsBuilder) SetAutoAck(b bool) *ReceiveOptsBuilder {
	bld.opts.autoAck = b
	return bld
}

func (bld *ReceiveOptsBuilder) SetExclusive(b bool) *ReceiveOptsBuilder {
	bld.opts.exclusive = b
	return bld
}

func (bld *ReceiveOptsBuilder) SetNoLocal(b bool) *ReceiveOptsBuilder {
	bld.opts.noLocal = b
	return bld
}

func (bld *ReceiveOptsBuilder) SetNoWait(b bool) *ReceiveOptsBuilder {
	bld.opts.noWait = b
	return bld
}

func (bld *ReceiveOptsBuilder) SetArgs(args *amqp.Table) *ReceiveOptsBuilder {
	bld.opts.args = args
	return bld
}

func (bld *ReceiveOptsBuilder) SetConsumerTag(tag string) *ReceiveOptsBuilder {
	bld.opts.consumerTag = tag
	return bld
}

func (bld *ReceiveOptsBuilder) Build() *ReceiveOpts {
	return bld.opts
}

// SendOpts 消息发送选项。
// mandatory 设为 true 时，交换器无法根据自身的类型和路由键找到一个符合条件的队列，
// 那么 RabbitMQ 会调用 Basic.Return 命令将消息返回给生产者。当该选项设置为 false 时，
// 出现上述情形，则消息直接被丢弃。
// immediate 设为 true 时，如果交换器在将消息路由到队列时发现队列上并不存在任何消费者，
// 那么这条消息将不会存入队列中。当与路由键匹配的所有队列都没有消费者时，该消息会通过 Basic.Return
// 返回至生产者。
// RabbitMQ 3.0版本开始去掉了对 immediate 参数的支持，对此 RabbitMQ 官方解释是: immediate
// 参数会影响镜像队列的性能，增加了代码复杂性，建议采用 TTL 和 DLX 的方法替代。
// messageFactory 如果未设置该选项，则默认使用 MessagePlainTransient 生产消息。
// retryable 如果不设置该选项，表示不启用消息重发功能。
type SendOpts struct {
	mandatory      bool
	immediate      bool
	messageFactory MessageFactory
	retryable      Retryable
}

// DefaultSendOpts 默认消息发送选项：消息无格式，非持久化，启用默认重试配置(DefaultTimesRetry)
func DefaultSendOpts() *SendOpts {
	return &SendOpts{messageFactory: MessagePlainTransient, retryable: DefaultTimesRetry()}
}

// SendOptsBuilder ...
type SendOptsBuilder struct {
	opts *SendOpts
}

// NewSendOptsBuilder creates a new SendOptsBuilder
func NewSendOptsBuilder() *SendOptsBuilder {
	return &SendOptsBuilder{DefaultSendOpts()}
}

// SetMandatory  设为 true 时，交换器无法根据自身的类型和路由键找到一个符合条件的队列，
// 那么 RabbitMQ 会调用 Basic.Return 命令将消息返回给生产者。当该选项设置为 false 时，
// 出现上述情形，则消息直接被丢弃。
func (bld *SendOptsBuilder) SetMandatory(b bool) *SendOptsBuilder {
	bld.opts.mandatory = b
	return bld
}

// SetImmediate 设为 true 时，如果交换器在将消息路由到队列时发现队列上并不存在任何消费者，
// 那么这条消息将不会存入队列中。当与路由键匹配的所有队列都没有消费者时，该消息会通过 Basic.Return
// 返回至生产者。
func (bld *SendOptsBuilder) SetImmediate(b bool) *SendOptsBuilder {
	bld.opts.immediate = b
	return bld
}

// SetMessageFactory 设置消息工厂方法。默认为 Plain Transient （无格式，非持久化）形式。
func (bld *SendOptsBuilder) SetMessageFactory(factory MessageFactory) *SendOptsBuilder {
	bld.opts.messageFactory = factory
	return bld
}

// SetRetryable 设置重试配置
func (bld *SendOptsBuilder) SetRetryable(retryable Retryable) *SendOptsBuilder {
	bld.opts.retryable = retryable
	return bld
}

// Build the sendOptsBuilder
func (bld *SendOptsBuilder) Build() *SendOpts {
	return bld.opts
}

// 如果没有消息则该方法阻塞等待；否则本方法会被持续调用，
// 直到主动停止消费（即本方法返回 true）。
// 返回值 brk 表示是否 break，即在循环消费过程中是否需要终止消费。
type ConsumerFunc func(*amqp.Delivery) (brk bool)

// Channel represents a channel
type Channel struct {
	*amqp.Channel
	conn       *Connection            // 用于断线重连
	confirming bool                   // producer
	confirms   chan amqp.Confirmation // producer
}

func newChannel(ch *amqp.Channel, conn *Connection) *Channel {
	return &Channel{
		Channel: ch,
		conn:    conn,
	}
}

// RemoveOperation removes an operation from the channel
func (c *Channel) RemoveOperation(key string) {
	c.conn.RemoveOperation(key)
}

// ReceiveOpts 持续接收消息并消费，除非 `<-chan amqp.Delivery` 关闭或 ConsumerFunc 主动放弃接收。
// 参数 opts 表示接收选项。opts 如果为 nil，将使用 DefaultReceiveOpts() 作为默认配置。
// 如果将参数 opts 的 autoAck 属性设为 false，则应该在 ReceiveListener.Consumer()
// 函数中调用 (*amqp.Delivery).Ack 手动确认消费；如果设为 true，已发送的消息会被服务器认为已被消费，
// 可能因网络状况不好等原因导致消息未被接收，从而造成数据丢失。 autoAck 为 false 时不提供自动确认，
// 是因为消费者有可能会需要拒绝确认，或在消费出现错误时不进行确认。
// 参数 consumer 用于处理接收操作。参数 consumer 一定不能为 nil，否则将 panic。
// 返回值：当 `<-chan amqp.Delivery` 关闭或 ConsumerFunc 主动放弃接收，返回 nil；其他情况则返回 error
func (c *Channel) ReceiveOpts(queue string, consumer ConsumerFunc, opts *ReceiveOpts) error {
	var err error
	if consumer == nil {
		panic("ConsumerFunc can't be nil")
	}
	if opts == nil {
		opts = DefaultReceiveOpts()
	}

	deliveries, err := c.Consume(
		queue,
		opts.consumerTag,
		opts.autoAck,
		opts.exclusive,
		opts.noLocal,
		opts.noWait,
		*getNonNilArgs(opts.args),
	)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		if consumer(&delivery) {
			break
		}
	}
	return nil
}

// Receive returns error if delivery
func (c *Channel) Receive(queue string, consumer ConsumerFunc) error {
	return c.ReceiveOpts(queue, consumer, nil)
}

// Send 如果使用了 Channel.SetXxx 设置了配置，将使用设定的配置发送消息，否则使用默认配置
func (c *Channel) Send(exchange string, routingKey string, body []byte) error {
	return c.SendOpts(exchange, routingKey, body, nil)
}

// SendOpts 发送消息。此方法不支持并发操作，如果需要并发发送，请先创建新的 Channel，再执行此方法。
// 参数 body 即需要发送的消息。
// 参数 opts 即发送消息需要配置的选项。如果 opts 为 nil，则表示使用默认配置。可以通过配置 SendOpts.retryable
// 启用消息重发的能力。请注意，由于消息重发使用的是同步的方式处理 ack，因此启用消息重发会极大降低 QPS。
func (c *Channel) SendOpts(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	if opts == nil {
		opts = DefaultSendOpts()
	}
	if opts.retryable == nil {
		return c.sendOpts(exchange, routingKey, body, opts)
	}
	return c.reSendSyncOpts(exchange, routingKey, body, opts)
}

// sendOpts 发送消息，但不确保送达。参数 opts 一定不能为 nil。
func (c *Channel) sendOpts(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	opts.messageFactory = getNonNilMessageFactory(opts.messageFactory)
	return c.Publish(exchange, routingKey, opts.mandatory, opts.immediate, opts.messageFactory(body))
}

// reSendSyncOpts 按照 Retryable 的配置内容确保发送消息是否到达。
// 该方法会在发送后等待确认消息，由于消息的发送和确认是同步的，所以在消息确认之前，不会继续发送下一个消息。
// 如果不想后续的消息被阻塞，请使用不同的 Channel 或 Connection 发送。
func (c *Channel) reSendSyncOpts(exchange string, routingKey string, body []byte, opts *SendOpts) (err error) {
	err = c.enableConfirm()
	if err != nil && !isConnectedErr(err) {
		return err
	}
	var retryable = opts.retryable
	var confirm *amqp.Confirmation
	retryable.retry(func() (brk bool) {
		confirm, err = c.sendAndWaitConfirmation(exchange, routingKey, body, opts)
		if confirm.Ack || !c.conn.CanRetry() {
			return true
		}
		c.resetChannelIfNeeded(err)
		return false
	})
	if !confirm.Ack {
		if err != nil {
			err = fmt.Errorf("send failed, cause nack: %w", err)
		} else {
			err = errors.New("send failed, cause nack")
		}
	}
	return
}

// sendAndWaitConfirmation 发送消息并等待确认信息。需要配合 enableConfirm 一起使用。
func (c *Channel) sendAndWaitConfirmation(exchange string, routingKey string, body []byte, opts *SendOpts) (*amqp.Confirmation, error) {
	err := c.sendOpts(exchange, routingKey, body, opts)
	confirm := <-c.confirms
	return &confirm, err
}

// resetChannelIfNeeded 如果必要（发生网络错误），则重置 Channel.Channel
func (c *Channel) resetChannelIfNeeded(err error) bool {
	if err == nil || !isConnectedErr(err) {
		// 如果 err==nil 也会从此退出
		return false
	}
	var e error
	var ch *amqp.Channel
	var conn = c.conn

	if !conn.IsOpen() {
		return false
	}

	if ch, e = conn.channel(); e != nil {
		logger.Logger.Debug(e)
		return false
	}

	c.resetChannel(ch)
	if e = c.enableConfirm(); e != nil {
		logger.Logger.Debug(e)
		return false
	}
	return true
}

// enableConfirm 启用 Confirm Mode。该方法会创建一个监听通道，用于监听发送是否成功
// 详见 (*amqp.Channel).enableConfirm
func (c *Channel) enableConfirm() error {
	if c.confirming {
		return nil
	}
	defer func() { c.confirming = true }()
	c.confirms = c.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := c.Channel.Confirm(false); err != nil {
		return err
	}
	return nil
}

// resetChannel 重置 Channel。关闭旧的 amqp.Channel，赋予新的 amqp.Channel，重置 Confirm Mode。
func (c *Channel) resetChannel(ch *amqp.Channel) {
	err := c.Channel.Close()
	if err != nil {
		logger.Logger.Debug(err)
	}
	c.Channel = ch
	// 重置 Confirm Mode
	c.confirming = false
	c.confirms = nil
}
