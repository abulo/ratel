package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Queue struct {
	c           *Connection
	declareOpts *QueueDeclareOpts
	bindOpts    *QueueBindOpts
	retryable   Retryable
	*amqp.Queue
}

// 根据 queueName 声明队列，并绑定 queueName, key 到指定的 exchange。
// Queue 可能会因为网络原因创建失败，不提供一定创建成功保证。
func (q *Queue) DeclareAndBind(queueName, key, exchange string) error {
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

type QueueBuilder struct {
	que *Queue
}

func NewQueueBuilder(c *Connection) *QueueBuilder {
	return &QueueBuilder{&Queue{c: c}}
}

func (bld *QueueBuilder) SetQueueDeclareOpts(builderFn func(builder *QueueDeclareOptsBuilder) *QueueDeclareOpts) *QueueBuilder {
	bld.que.declareOpts = builderFn(NewQueueDeclareOptsBuilder())
	return bld
}

func (bld *QueueBuilder) SetQueueBindOpts(builderFn func(builder *QueueBindOptsBuilder) *QueueBindOpts) *QueueBuilder {
	bld.que.bindOpts = builderFn(NewQueueBindOptsBuilder())
	return bld
}

func (bld *QueueBuilder) SetRetryable(retryable Retryable) *QueueBuilder {
	bld.que.retryable = retryable
	return bld
}

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

type QueueDeclareOpts struct {
	durable, autoDelete, exclusive, noWait bool
	args                                   *amqp.Table
}

func DefaultQueueDeclareOpts() *QueueDeclareOpts {
	return &QueueDeclareOpts{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

type QueueDeclareOptsBuilder struct {
	opts *QueueDeclareOpts
}

func NewQueueDeclareOptsBuilder() *QueueDeclareOptsBuilder {
	return &QueueDeclareOptsBuilder{DefaultQueueDeclareOpts()}
}

func (bld *QueueDeclareOptsBuilder) SetDurable(b bool) *QueueDeclareOptsBuilder {
	bld.opts.durable = b
	return bld
}

func (bld *QueueDeclareOptsBuilder) SetAutoDelete(b bool) *QueueDeclareOptsBuilder {
	bld.opts.autoDelete = b
	return bld
}

func (bld *QueueDeclareOptsBuilder) SetExclusive(b bool) *QueueDeclareOptsBuilder {
	bld.opts.exclusive = b
	return bld
}

func (bld *QueueDeclareOptsBuilder) SetNowait(b bool) *QueueDeclareOptsBuilder {
	bld.opts.noWait = b
	return bld
}

func (bld *QueueDeclareOptsBuilder) SetArgs(args *amqp.Table) *QueueDeclareOptsBuilder {
	bld.opts.args = args
	return bld
}

func (bld *QueueDeclareOptsBuilder) Build() *QueueDeclareOpts {
	return bld.opts
}

type QueueBindOpts struct {
	noWait bool
	args   *amqp.Table
}

func DefaultBindOpts() *QueueBindOpts {
	return &QueueBindOpts{
		noWait: false,
		args:   nil,
	}
}

type QueueBindOptsBuilder struct {
	opts *QueueBindOpts
}

func NewQueueBindOptsBuilder() *QueueBindOptsBuilder {
	return &QueueBindOptsBuilder{DefaultBindOpts()}
}

func (bld *QueueBindOptsBuilder) SetNoWait(b bool) *QueueBindOptsBuilder {
	bld.opts.noWait = b
	return bld
}

func (bld *QueueBindOptsBuilder) SetArgs(args *amqp.Table) *QueueBindOptsBuilder {
	bld.opts.args = args
	return bld
}

func (bld *QueueBindOptsBuilder) Build() *QueueBindOpts {
	return bld.opts
}
