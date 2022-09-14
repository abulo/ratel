package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Queue represents a queue of messages
type Queue struct {
	c           *Connection
	declareOpts *QueueDeclareOpts
	bindOpts    *QueueBindOpts
	retryable   Retryable
	*amqp.Queue
}

// DeclareAndBind 根据 queueName 声明队列，并绑定 queueName, key 到指定的 exchange。
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

// RetryDeclareAndBind will retry the specified	 messages	 Channels
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

// QueueBuilder ...
type QueueBuilder struct {
	que *Queue
}

// NewQueueBuilder creates a new QueueBuilder instance
func NewQueueBuilder(c *Connection) *QueueBuilder {
	return &QueueBuilder{&Queue{c: c}}
}

// SetQueueDeclareOpts 	sets the default value for the Queue declaration
func (bld *QueueBuilder) SetQueueDeclareOpts(builderFn func(builder *QueueDeclareOptsBuilder) *QueueDeclareOpts) *QueueBuilder {
	bld.que.declareOpts = builderFn(NewQueueDeclareOptsBuilder())
	return bld
}

// SetQueueBindOpts sets the binding options for the queue builde
func (bld *QueueBuilder) SetQueueBindOpts(builderFn func(builder *QueueBindOptsBuilder) *QueueBindOpts) *QueueBuilder {
	bld.que.bindOpts = builderFn(NewQueueBindOptsBuilder())
	return bld
}

// SetRetryable sets the retryable flag to true for
func (bld *QueueBuilder) SetRetryable(retryable Retryable) *QueueBuilder {
	bld.que.retryable = retryable
	return bld
}

// Build the queue
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

// QueueDeclareOpts returns the list of queued operations
type QueueDeclareOpts struct {
	durable, autoDelete, exclusive, noWait bool
	args                                   *amqp.Table
}

// DefaultQueueDeclareOpts  holds the default queue declaration
func DefaultQueueDeclareOpts() *QueueDeclareOpts {
	return &QueueDeclareOpts{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	}
}

// QueueDeclareOptsBuilder is the builder for QueueDeclareOpts
type QueueDeclareOptsBuilder struct {
	opts *QueueDeclareOpts
}

// NewQueueDeclareOptsBuilder creates a new QueueDeclareOptsBuilder
func NewQueueDeclareOptsBuilder() *QueueDeclareOptsBuilder {
	return &QueueDeclareOptsBuilder{DefaultQueueDeclareOpts()}
}

// SetDurable sets whether the queue declaration
func (bld *QueueDeclareOptsBuilder) SetDurable(b bool) *QueueDeclareOptsBuilder {
	bld.opts.durable = b
	return bld
}

// SetAutoDelete sets the auto delete option
func (bld *QueueDeclareOptsBuilder) SetAutoDelete(b bool) *QueueDeclareOptsBuilder {
	bld.opts.autoDelete = b
	return bld
}

// SetExclusive sets the exclusive
func (bld *QueueDeclareOptsBuilder) SetExclusive(b bool) *QueueDeclareOptsBuilder {
	bld.opts.exclusive = b
	return bld
}

// SetNoWait sets the number of seconds to wait before returning
func (bld *QueueDeclareOptsBuilder) SetNoWait(b bool) *QueueDeclareOptsBuilder {
	bld.opts.noWait = b
	return bld
}

// SetArgs sets the arguments to be passed to the command
func (bld *QueueDeclareOptsBuilder) SetArgs(args *amqp.Table) *QueueDeclareOptsBuilder {
	bld.opts.args = args
	return bld
}

// Build returns the queue declaration options
func (bld *QueueDeclareOptsBuilder) Build() *QueueDeclareOpts {
	return bld.opts
}

// QueueBindOpts is the queue binding options
type QueueBindOpts struct {
	noWait bool
	args   *amqp.Table
}

// DefaultBindOpts returns the default queue binding options
func DefaultBindOpts() *QueueBindOpts {
	return &QueueBindOpts{
		noWait: false,
		args:   nil,
	}
}

// QueueBindOptsBuilder is the builder for QueueBindOptsBuilder
type QueueBindOptsBuilder struct {
	opts *QueueBindOpts
}

// NewQueueBindOptsBuilder returns a new QueueBindOptsBuilder
func NewQueueBindOptsBuilder() *QueueBindOptsBuilder {
	return &QueueBindOptsBuilder{DefaultBindOpts()}
}

// SetNoWait sets the default value
func (bld *QueueBindOptsBuilder) SetNoWait(b bool) *QueueBindOptsBuilder {
	bld.opts.noWait = b
	return bld
}

// SetArgs sets the arguments
func (bld *QueueBindOptsBuilder) SetArgs(args *amqp.Table) *QueueBindOptsBuilder {
	bld.opts.args = args
	return bld
}

// Build sets the parameters
func (bld *QueueBindOptsBuilder) Build() *QueueBindOpts {
	return bld.opts
}
