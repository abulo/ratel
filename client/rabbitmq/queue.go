package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Queue 绑定所属连接与队列声明、绑定选项，提供队列声明并绑定到交换器的能力。
//
// 由 QueueBuilder.Build 构造，零值缺少连接与选项，不可直接使用。
type Queue struct {
	c           *Connection
	declareOpts *QueueDeclareOpts
	bindOpts    *QueueBindOpts
	retryable   Retryable
	*amqp.Queue
}

// DeclareAndBind 声明指定队列，并将其按 key 绑定到 exchange；网络原因可能导致声明失败，
// 不保证一定创建成功。
func (q *Queue) DeclareAndBind(queueName, key, exchange string) error {
	// PERF-5-13 [P1]：复核结论——队列声明/绑定属于启动期低频操作，按报告建议维持"每次新建通道"现状，
	// 不纳入生产者/消费者的常驻通道复用，避免声明状态与发布/消费通道相互影响。
	ch, err := q.c.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	queue, err := ch.QueueDeclare(
		queueName,
		q.declareOpts.durable,
		q.declareOpts.autoDelete,
		q.declareOpts.exclusive,
		q.declareOpts.noWait,
		*getNonNilArgs(q.declareOpts.args),
	)
	if err != nil {
		return err
	}
	q.Queue = &queue
	err = ch.QueueBind(queueName, key, exchange, q.declareOpts.noWait, *getNonNilArgs(q.declareOpts.args))
	if err != nil {
		return err
	}
	return nil
}

// RetryDeclareAndBind 按队列自带的重试策略反复执行声明与绑定，重试耗尽仍失败时返回最后一次错误。
func (q *Queue) RetryDeclareAndBind(queueName, key, exchange string) error {
	var err error
	retryable := q.retryable
	retryable.retry(func() (brk bool) {
		err = q.DeclareAndBind(queueName, key, exchange)
		if err == nil || !q.c.CanRetry() {
			return true
		}
		return false
	})
	return err
}

// QueueBuilder 以链式方式为 Connection 装配队列声明、绑定选项与重试策略，最终构造 Queue。
type QueueBuilder struct {
	que *Queue
}

// NewQueueBuilder 创建绑定到指定连接的 QueueBuilder。
func NewQueueBuilder(c *Connection) *QueueBuilder {
	return &QueueBuilder{&Queue{c: c}}
}

// SetQueueDeclareOpts 设置队列声明选项，返回构造器自身以便链式调用。
func (bld *QueueBuilder) SetQueueDeclareOpts(builderFn func(builder *QueueDeclareOptsBuilder) *QueueDeclareOpts) *QueueBuilder {
	bld.que.declareOpts = builderFn(NewQueueDeclareOptsBuilder())
	return bld
}

// SetQueueBindOpts 设置队列绑定选项，返回构造器自身以便链式调用。
func (bld *QueueBuilder) SetQueueBindOpts(builderFn func(builder *QueueBindOptsBuilder) *QueueBindOpts) *QueueBuilder {
	bld.que.bindOpts = builderFn(NewQueueBindOptsBuilder())
	return bld
}

// SetRetryable 设置声明与绑定失败时的重试策略，返回构造器自身以便链式调用。
func (bld *QueueBuilder) SetRetryable(retryable Retryable) *QueueBuilder {
	bld.que.retryable = retryable
	return bld
}

// Build 完成队列构造；未显式设置的选项分别回落到默认声明选项、默认绑定选项与默认持续重试策略。
func (bld *QueueBuilder) Build() *Queue {
	que := bld.que
	if que.declareOpts == nil {
		que.declareOpts = DefaultQueueDeclareOpts()
	}
	if que.bindOpts == nil {
		que.bindOpts = DefaultBindOpts()
	}
	if que.retryable == nil {
		que.retryable = DefaultTimesRetry()
	}
	return que
}

// QueueDeclareOpts 封装队列声明的参数，由 QueueDeclareOptsBuilder 构造。
type QueueDeclareOpts struct {
	durable, autoDelete, exclusive, noWait bool
	args                                   *amqp.Table
}

// DefaultQueueDeclareOpts 返回默认声明选项：队列持久化，不自动删除、不独占、不等服务器确认。
func DefaultQueueDeclareOpts() *QueueDeclareOpts {
	return &QueueDeclareOpts{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

// QueueDeclareOptsBuilder 以链式方式构造 QueueDeclareOpts。
type QueueDeclareOptsBuilder struct {
	opts *QueueDeclareOpts
}

// NewQueueDeclareOptsBuilder 创建基于默认选项的 QueueDeclareOptsBuilder。
func NewQueueDeclareOptsBuilder() *QueueDeclareOptsBuilder {
	return &QueueDeclareOptsBuilder{DefaultQueueDeclareOpts()}
}

// SetDurable 设置队列是否持久化（服务器重启后是否保留），返回构造器自身以便链式调用。
func (bld *QueueDeclareOptsBuilder) SetDurable(b bool) *QueueDeclareOptsBuilder {
	bld.opts.durable = b
	return bld
}

// SetAutoDelete 设置连接断开后是否自动删除队列，返回构造器自身以便链式调用。
func (bld *QueueDeclareOptsBuilder) SetAutoDelete(b bool) *QueueDeclareOptsBuilder {
	bld.opts.autoDelete = b
	return bld
}

// SetExclusive 设置队列是否独占（仅当前连接可访问），返回构造器自身以便链式调用。
func (bld *QueueDeclareOptsBuilder) SetExclusive(b bool) *QueueDeclareOptsBuilder {
	bld.opts.exclusive = b
	return bld
}

// SetNoWait 设置是否不等服务器确认即返回，返回构造器自身以便链式调用。
func (bld *QueueDeclareOptsBuilder) SetNoWait(b bool) *QueueDeclareOptsBuilder {
	bld.opts.noWait = b
	return bld
}

// SetArgs 设置传给队列声明的自定义参数表，返回构造器自身以便链式调用。
func (bld *QueueDeclareOptsBuilder) SetArgs(args *amqp.Table) *QueueDeclareOptsBuilder {
	bld.opts.args = args
	return bld
}

// Build 完成声明选项构造，返回最终的 QueueDeclareOpts。
func (bld *QueueDeclareOptsBuilder) Build() *QueueDeclareOpts {
	return bld.opts
}

// QueueBindOpts 封装队列绑定到交换器的参数，由 QueueBindOptsBuilder 构造。
type QueueBindOpts struct {
	noWait bool
	args   *amqp.Table
}

// DefaultBindOpts 返回默认绑定选项：等待服务器确认，无自定义参数。
func DefaultBindOpts() *QueueBindOpts {
	return &QueueBindOpts{
		noWait: false,
		args:   nil,
	}
}

// QueueBindOptsBuilder 以链式方式构造 QueueBindOpts。
type QueueBindOptsBuilder struct {
	opts *QueueBindOpts
}

// NewQueueBindOptsBuilder 创建基于默认选项的 QueueBindOptsBuilder。
func NewQueueBindOptsBuilder() *QueueBindOptsBuilder {
	return &QueueBindOptsBuilder{DefaultBindOpts()}
}

// SetNoWait 设置是否不等服务器确认即返回，返回构造器自身以便链式调用。
func (bld *QueueBindOptsBuilder) SetNoWait(b bool) *QueueBindOptsBuilder {
	bld.opts.noWait = b
	return bld
}

// SetArgs 设置传给绑定命令的自定义参数表，返回构造器自身以便链式调用。
func (bld *QueueBindOptsBuilder) SetArgs(args *amqp.Table) *QueueBindOptsBuilder {
	bld.opts.args = args
	return bld
}

// Build 完成绑定选项构造，返回最终的 QueueBindOpts。
func (bld *QueueBindOptsBuilder) Build() *QueueBindOpts {
	return bld.opts
}
