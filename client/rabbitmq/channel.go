package rabbitmq

import (
	"fmt"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ReceiveOpts 是消费消息的选项集合，通过 ReceiveOptsBuilder 构造。
//
// autoAck 为 false 时，RabbitMQ 要等消费者显式确认后才从内存或磁盘移除消息，调用方须主动
// 调用 (*amqp.Delivery).Ack 确认，防止消息在服务端堆积；autoAck 为 true 时消息由服务器自动
// 视为已消费，消费者处理失败时可能造成消息丢失。若一直未收到确认且消费者已断开，RabbitMQ 会将
// 消息重新入队投递给其他消费者。consumerTag 用于唯一标识消费者，留空由服务器自动生成；其余参数
// 无特殊需求时保持默认即可。
type ReceiveOpts struct {
	autoAck, exclusive, noLocal, noWait bool
	args                                *amqp.Table
	consumerTag                         string
}

// DefaultReceiveOpts 返回默认接收选项，autoAck 默认为 true。
func DefaultReceiveOpts() *ReceiveOpts {
	return &ReceiveOpts{autoAck: true}
}

// ReceiveOptsBuilder 以链式方式构造 ReceiveOpts。
type ReceiveOptsBuilder struct {
	opts *ReceiveOpts
}

// NewReceiveOptsBuilder 创建基于默认选项的 ReceiveOptsBuilder。
func NewReceiveOptsBuilder() *ReceiveOptsBuilder {
	return &ReceiveOptsBuilder{DefaultReceiveOpts()}
}

// SetAutoAck 设置是否由服务器自动确认消息，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetAutoAck(b bool) *ReceiveOptsBuilder {
	bld.opts.autoAck = b
	return bld
}

// SetExclusive 设置是否声明独占消费者，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetExclusive(b bool) *ReceiveOptsBuilder {
	bld.opts.exclusive = b
	return bld
}

// SetNoLocal 设置是否不接收本连接发布的消息，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetNoLocal(b bool) *ReceiveOptsBuilder {
	bld.opts.noLocal = b
	return bld
}

// SetNoWait 设置是否不等服务器确认即返回，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetNoWait(b bool) *ReceiveOptsBuilder {
	bld.opts.noWait = b
	return bld
}

// SetArgs 设置传给消费声明的自定义参数表，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetArgs(args *amqp.Table) *ReceiveOptsBuilder {
	bld.opts.args = args
	return bld
}

// SetConsumerTag 设置消费者唯一标识，返回构造器自身以便链式调用。
func (bld *ReceiveOptsBuilder) SetConsumerTag(tag string) *ReceiveOptsBuilder {
	bld.opts.consumerTag = tag
	return bld
}

// Build 完成选项构造，返回最终的 ReceiveOpts。
func (bld *ReceiveOptsBuilder) Build() *ReceiveOpts {
	return bld.opts
}

// SendOpts 是发送消息的选项集合，通过 SendOptsBuilder 构造。
//
// mandatory 为 true 时，交换器无法按类型和路由键匹配到队列，RabbitMQ 会通过 Basic.Return 将
// 消息退回生产者；为 false 时这类消息会被直接丢弃。immediate 为 true 时，路由到的队列上没有
// 消费者则消息不入队并通过 Basic.Return 退回生产者；RabbitMQ 3.0 起已移除 immediate 支持，
// 官方建议改用 TTL 与 DLX 实现等价语义。messageFactory 在发送路径中未设置时回退为
// MessagePlainPersistent；retryable 未设置时不启用消息重发。
type SendOpts struct {
	mandatory      bool
	immediate      bool
	messageFactory MessageFactory
	retryable      Retryable
}

// DefaultSendOpts 返回默认发送选项：text/plain、非持久化，并启用默认持续重试策略。
func DefaultSendOpts() *SendOpts {
	return &SendOpts{messageFactory: MessagePlainTransient, retryable: DefaultTimesRetry()}
}

// SendOptsBuilder 以链式方式构造 SendOpts。
type SendOptsBuilder struct {
	opts *SendOpts
}

// NewSendOptsBuilder 创建基于默认选项的 SendOptsBuilder。
func NewSendOptsBuilder() *SendOptsBuilder {
	return &SendOptsBuilder{DefaultSendOpts()}
}

// SetMandatory 设置无法路由时是否将消息退回生产者，返回构造器自身以便链式调用。
func (bld *SendOptsBuilder) SetMandatory(b bool) *SendOptsBuilder {
	bld.opts.mandatory = b
	return bld
}

// SetImmediate 设置无消费者时是否将消息退回生产者；注意 RabbitMQ 3.0 起不再支持该参数。
func (bld *SendOptsBuilder) SetImmediate(b bool) *SendOptsBuilder {
	bld.opts.immediate = b
	return bld
}

// SetMessageFactory 设置消息工厂，返回构造器自身以便链式调用；未设置时发送路径回退为持久化工厂。
func (bld *SendOptsBuilder) SetMessageFactory(factory MessageFactory) *SendOptsBuilder {
	bld.opts.messageFactory = factory
	return bld
}

// SetRetryable 设置发布确认重试策略，返回构造器自身以便链式调用；nil 表示不启用重发。
func (bld *SendOptsBuilder) SetRetryable(retryable Retryable) *SendOptsBuilder {
	bld.opts.retryable = retryable
	return bld
}

// Build 完成选项构造，返回最终的 SendOpts。
func (bld *SendOptsBuilder) Build() *SendOpts {
	return bld.opts
}

// ConsumerFunc 是消息消费回调：delivery 为当前投递的消息，返回 true 表示终止消费并退出循环，
// 返回 false 继续接收下一条；队列投递通道关闭时消费循环自然结束。
type ConsumerFunc func(*amqp.Delivery) (brk bool)

// Channel 在 amqp.Channel 基础上绑定所属 Connection，提供带选项的消息收发与发布确认支持。
//
// 零值不可直接使用，由 Connection.Channel 创建；单条 Channel 不支持并发发送，需要并发发送时
// 应先创建新的 Channel。
type Channel struct {
	*amqp.Channel
	// conn 为所属连接，用于断线后重建通道。
	conn *Connection
	// confirming 标记当前通道是否已进入 Confirm 模式。
	confirming bool
	// confirms 为发布确认结果的监听通道。
	confirms chan amqp.Confirmation
}

func newChannel(ch *amqp.Channel, conn *Connection) *Channel {
	return &Channel{
		Channel: ch,
		conn:    conn,
	}
}

// RemoveOperation 转发到所属连接，移除指定标识的已注册操作。
func (c *Channel) RemoveOperation(key string) {
	c.conn.RemoveOperation(key)
}

// ReceiveOpts 在指定队列上持续接收消息并交给 consumer 消费，直到投递通道关闭或 consumer 返回 true
// 主动放弃。
//
// opts 为 nil 时使用 DefaultReceiveOpts；consumer 为 nil 会 panic。autoAck 为 false 时须在回调中
// 调用 (*amqp.Delivery).Ack 手动确认，为 true 时消息由服务器自动确认，网络异常时可能丢失。正常
// 结束返回 nil，订阅或消费过程中的错误直接返回。
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

// Receive 使用默认接收选项在指定队列上持续消费消息，等价于 ReceiveOpts(queue, consumer, nil)。
func (c *Channel) Receive(queue string, consumer ConsumerFunc) error {
	return c.ReceiveOpts(queue, consumer, nil)
}

// Send 使用默认发送选项向指定交换器与路由键发布消息。
func (c *Channel) Send(exchange string, routingKey string, body []byte) error {
	return c.SendOpts(exchange, routingKey, body, nil)
}

// SendOpts 向指定交换器与路由键发布消息；opts 为 nil 时使用默认配置。
//
// 该方法不支持并发操作，并发发送须先创建新的 Channel。opts.retryable 非空时进入同步发布确认
// 流程，消息确认前不会继续发送下一条，会显著降低 QPS；不希望阻塞后续消息时应使用不同的 Channel
// 或 Connection 发送。
func (c *Channel) SendOpts(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	if opts == nil {
		opts = DefaultSendOpts()
	}
	if opts.retryable == nil {
		return c.sendOpts(exchange, routingKey, body, opts)
	}
	return c.reSendSyncOpts(exchange, routingKey, body, opts)
}

// sendOpts 直接发布消息，不保证送达；调用方须保证 opts 非 nil。
func (c *Channel) sendOpts(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	opts.messageFactory = getNonNilMessageFactory(opts.messageFactory)
	return c.Publish(exchange, routingKey, opts.mandatory, opts.immediate, opts.messageFactory(body))
}

// reSendSyncOpts 按重试策略确保消息送达：发送后同步等待发布确认，确认前阻塞后续消息；收到 nack
// 或连接错误时按策略重试，重试耗尽仍未确认则返回 nack 错误。
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

// sendAndWaitConfirmation 发布消息并阻塞等待一条发布确认结果，须在 enableConfirm 之后调用。
func (c *Channel) sendAndWaitConfirmation(exchange string, routingKey string, body []byte, opts *SendOpts) (*amqp.Confirmation, error) {
	err := c.sendOpts(exchange, routingKey, body, opts)
	confirm := <-c.confirms
	return &confirm, err
}

// resetChannelIfNeeded 在发生可恢复连接错误时重建底层通道并重新进入 Confirm 模式；错误不可恢复
// 或连接已关闭时不做处理，返回 false。
func (c *Channel) resetChannelIfNeeded(err error) bool {
	if err == nil || !isConnectedErr(err) {
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

// enableConfirm 将通道切换到 Confirm 模式并创建确认监听通道；重复调用直接返回 nil。
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

// resetChannel 关闭旧的底层通道并替换为新通道，同时重置 Confirm 模式状态与确认通道。
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
