package rabbitmq

import (
	"net"
	"os"
	"sync"
	"syscall"

	"github.com/abulo/ratel/core/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	AMQP  Schema = "amqp"
	AMQPS Schema = "amqps"
)

type Schema string
type Operation func(key string, ch *Channel)
type Operations map[string]Operation

// Connection amqp 连接。 Connection 创建后不会直接连接服务器，而是要调用 Dial 后才会执行连接服务器操作
type Connection struct {
	c             *amqp.Connection // 用于真正发起一个 amqp 连接
	cMut          sync.RWMutex     // 用于读写 c 时加锁
	url           string
	retryable     Retryable // 重试配置
	operations    Operations
	oMut          sync.Mutex    // 用于读写 operations 时加锁
	genOptKeyFunc func() string // 用于生成 operations 的 key，每次调用都会生成新的 key
	sync.Once                   // 用于保证 Dial 只被调用一次
}

// NewConnection retryable 如果为 nil，则使用 emptyRetryable 替换。emptyRetryable 不会尝试重试操作。
func NewConnection(url string, retryable Retryable) *Connection {
	return &Connection{
		url:           url,
		retryable:     getNonNilRetryable(retryable),
		operations:    make(Operations, 0),
		genOptKeyFunc: NewDefaultSAdder(),
	}
}

func (c *Connection) setConn(conn *amqp.Connection) {
	c.cMut.Lock()
	defer c.cMut.Unlock()
	c.c = conn
}

// Dial 连接服务器。仅允许被调用一次。
func (c *Connection) Dial() error {
	var err error
	c.Do(func() {
		err = c.reDial()
	})
	return err
}

// dial 连接服务器，但不提供断线重连
func (c *Connection) dial() error {
	var err error
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.setConn(conn)
	return nil
}

// 尝试重连，如果重连成功，执行监听操作
func (c *Connection) reDial() error {
	err := c.doReDial()
	if err != nil {
		return err
	}

	// 注册重连监听器
	c.cMut.RLock()
	defer c.cMut.RUnlock()
	if c.c == nil {
		return amqp.ErrClosed
	}

	// 执行监听操作
	monitor := c.c.NotifyClose(make(chan *amqp.Error))
	go c.reconnectListener(monitor)

	return nil
}

// 根据配置尝试重连
func (c *Connection) doReDial() error {
	retryable := c.retryable
	var err error
	retryable.retry(func() (brk bool) {
		err = c.dial()
		// 连接成功，退出循环；
		// 不是网络连接错误，退出循环。
		if err == nil || !isConnectedErr(err) {
			return true
		}
		logger.Logger.Debug("try to re-dial...")
		return false
	})
	return err
}

// reconnect 根据错误判断是否需要重连，并返回成功与否。如果 err 为 nil，将直接返回 true。
func (c *Connection) reconnect(err error) bool {
	if err == nil {
		return true
	}
	for {
		logger.Logger.Debug("try to reconnect...")
		if !isAmqpConnectedErr(err) {
			logger.Logger.Info("reconnect failed: ", err)
			return false
		}

		// 先关闭以前的 conn
		_ = c.Close()

		err = c.reDial()
		if err == nil {
			logger.Logger.Debug("reconnected!")
			break
		}
	}
	return true
}

// Channel 可用于发送、接收消息。
// 函数会先判断是否已连接，否则将尝试重连（使用您之前设置的 Retryable 配置）。
// 在获得连接的情况下，会立刻创建 Channel。但可能会存在极少数情况下，因为网络不稳定等因素，
// Channel 创建之前，连接又断开，则会因为网络原因产生错误。
func (c *Connection) Channel() (*Channel, error) {
	var ch, err = c.channel()
	if err != nil {
		return nil, err
	}
	return newChannel(ch, c), nil
}

func (c *Connection) channel() (*amqp.Channel, error) {
	c.cMut.RLock()
	defer c.cMut.RUnlock()
	return c.c.Channel()
}

// RetryChannel 将在 Channel 创建失败后，尝试重试。如果仍然失败，将返回造成失败的原因。
func (c *Connection) RetryChannel(retryable Retryable) (ch *Channel, err error) {
	retryable = getNonNilRetryable(retryable)
	ch, err = c.Channel()
	if err != nil {
		retryable.retry(func() (brk bool) {
			ch, err = c.Channel()
			return err == nil
		})
	}
	return ch, err
}

// RegisterAndExec 注册并执行 Operation 。
// 每次重连后，重连监听器会自动使用 exec 执行一遍所有的 Operation 函数。
// 注意：
//   - 函数会在 Operation 执行完后主动关闭 Channel，因此我们无需在 Operation 中手动关闭 Channel。
//   - 由于使用了 go routine，该方法可能会在 Operation 操作执行完毕前返回。
func (c *Connection) RegisterAndExec(opt Operation) {
	if opt == nil {
		panic("Operation must not be nil")
	}
	var key = c.addOperation(opt)
	c.execOperation(key, opt)
}

// 如果 Channel 创建出错，立即返回错误；否则使用 go routine 执行 Operation。
func (c *Connection) execOperation(key string, opt Operation) error {
	channel, err := c.Channel()
	if err != nil {
		return err
	}
	go func() {
		defer channel.Close()
		opt(key, channel)
	}()
	return nil
}

// 添加你想通过 Channel 执行的断线重连操作。
//
// 一般建议添加用于"接收消息"的操作，因为我们通常不会需要每次断线重连后重发消息。
func (c *Connection) addOperation(opt Operation) (key string) {
	c.oMut.Lock()
	defer c.oMut.Unlock()
	// genOptKeyFunc 能生成从 0 开始到 2^64-1 个数，一纳秒一个 key，可以用 585 年。
	key = c.genOptKeyFunc()
	c.operations[key] = opt
	return key
}

// exec 将会使用 go routine 逐个执行 Operation
func (c *Connection) exec() {
	var opts Operations
	c.oMut.Lock()
	defer c.oMut.Unlock()
	opts = c.operations

	for key, opt := range opts {
		fn := opt
		err := c.execOperation(key, fn)
		if err != nil {
			logger.Logger.Warn(err)
			return
		}
	}
}

func (c *Connection) RemoveOperation(key string) {
	go func() {
		c.oMut.Lock()
		defer c.oMut.Unlock()
		delete(c.operations, key)
	}()
}

// Close 关闭 Connection。
func (c *Connection) Close() error {
	// 防止 c.Connection 并发关闭
	c.cMut.Lock()
	defer c.cMut.Unlock()

	if c.c != nil {
		defer func() { c.c = nil }()
		return c.c.Close()
	}
	return nil
}

// IsOpen returns true
func (c *Connection) IsOpen() bool {
	c.cMut.RLock()
	defer c.cMut.RUnlock()
	return c.c != nil
}

// CanRetry returns true
func (c *Connection) CanRetry() bool {
	return !c.retryable.hasGaveUp()
}

// Consumer returns the consumer for the connection
func (c *Connection) Consumer() *Consumer {
	return &Consumer{c}
}

// Producer returns the producer for the connection
func (c *Connection) Producer() *Producer {
	return &Producer{c}
}

// QueueBuilder returns the queue builder for the connection
func (c *Connection) QueueBuilder() *QueueBuilder {
	return NewQueueBuilder(c)
}

// reconnectListener 重连监听器。会不断监听是否断线。如果断线了，则尝试重连。
// 如果重连成功，则执行注册的 Operation 。
// 理论上只要通过 dial 部分的连接，后面不大可能存在连接不上的问题。
// 因此此处使用了无限循环重试的方式，除非 dial 失败达到指定次数。
func (c *Connection) reconnectListener(monitor chan *amqp.Error) {
	err, ok := <-monitor
	if !ok {
		return
	}
	if c.reconnect(err) {
		c.exec()
	}
}

func isAmqpConnectedErr(err error) bool {
	var amqpErr *amqp.Error
	return errors.As(err, &amqpErr) &&
		(errors.Is(amqpErr, amqp.ErrClosed) ||
			amqpErr.Reason == "EOF" ||
			amqpErr.Code == amqp.ConnectionForced)
}

// 是否为传输层网络错误
func isTransportNetError(err error) bool {
	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}

	if netErr.Timeout() {
		logger.Logger.Debug(err)
		return true
	}

	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return false
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		logger.Logger.Debug("net.DNSError:%+v", t)
		return true
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNREFUSED:
				logger.Logger.Debug("connect refused")
				return true
			case syscall.ETIMEDOUT:
				logger.Logger.Debug("timeout")
				return true
			}
		}
	}

	return false
}

// 判断是否为连接错误
func isConnectedErr(err error) bool {
	return isTransportNetError(err) || isAmqpConnectedErr(err)
}
