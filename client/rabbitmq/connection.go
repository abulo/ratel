package rabbitmq

import (
	"net"
	"os"
	"sync"
	"syscall"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// 协议方案名，用于拼接 AMQP 连接 URL。
const (
	// AMQP 为明文 AMQP 协议方案。
	AMQP Schema = "amqp"
	// AMQPS 为 TLS 加密的 AMQP 协议方案。
	AMQPS Schema = "amqps"
)

// Schema 表示 AMQP 连接 URL 使用的协议方案名。
type Schema string

// Operation 表示断线重连后需要重新执行的通道操作，key 为注册时分配的唯一标识。
type Operation func(key string, ch *Channel)

// Operations 以唯一标识索引一组待重连后重新执行的通道操作。
type Operations map[string]Operation

// Connection 封装一条 AMQP 连接，负责断线检测、按配置重试重连，以及重连后重新执行已注册的操作。
//
// 须通过 NewConnection 创建并调用 Dial 完成首次连接，Dial 仅允许调用一次；内部锁保护连接与
// 操作表的并发访问。
type Connection struct {
	// c 为底层 AMQP 连接。
	c *amqp.Connection
	// cMut 保护 c 的并发读写。
	cMut sync.RWMutex
	url  string
	// retryable 为重连重试策略；nil 会被替换为不重试的空实现。
	retryable  Retryable
	operations Operations
	// oMut 保护 operations 的并发读写。
	oMut sync.Mutex
	// genOptKeyFunc 为已注册操作生成唯一标识。
	genOptKeyFunc func() string
	sync.Once

	// PERF-5-13 [P1]：普通发布、确认模式发布、拉取消费各持一条互斥保护的复用通道；
	// amqp 单通道非并发安全，互斥锁保证同一时刻只有一个使用者操作该通道，三条槽位互不干扰。
	pubCh     *Channel
	pubMu     sync.Mutex
	confirmCh *Channel
	confirmMu sync.Mutex
	getCh     *Channel
	getMu     sync.Mutex
}

// NewConnection 创建一条尚未连接的 Connection；retryable 为 nil 时使用不做任何重试的空实现。
// 创建后必须调用 Dial 才会真正连接服务器。
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

// Dial 建立与服务器的首次连接；借助 sync.Once 保证仅执行一次，重连过程遵循配置的重试策略。
func (c *Connection) Dial() error {
	var err error
	c.Do(func() {
		err = c.reDial()
	})
	return err
}

// dial 直接发起一次 AMQP 拨号，不携带断线重连逻辑。
func (c *Connection) dial() error {
	var err error
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.setConn(conn)
	return nil
}

// reDial 按重试策略完成拨号，并在连接成功后注册关闭事件监听器，用于后续断线检测。
func (c *Connection) reDial() error {
	err := c.doReDial()
	if err != nil {
		return err
	}

	c.cMut.RLock()
	defer c.cMut.RUnlock()
	if c.c == nil {
		return amqp.ErrClosed
	}

	monitor := c.c.NotifyClose(make(chan *amqp.Error))
	go c.reconnectListener(monitor)

	return nil
}

// doReDial 依据重试策略反复拨号，连接成功或遇到非网络类错误时停止重试。
func (c *Connection) doReDial() error {
	retryable := c.retryable
	var err error
	retryable.retry(func() (brk bool) {
		err = c.dial()
		// 连接成功退出循环；非网络连接错误同样退出循环，交由上层判断。
		if err == nil || !isConnectedErr(err) {
			return true
		}
		logger.Logger.Debug("try to re-dial...")
		return false
	})
	return err
}

// reconnect 根据错误类型判断是否需要重连；err 为 nil 直接返回 true。非可恢复的连接错误
// 不会重试，返回 false。
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

		// 重连前先关闭旧连接，避免残留的底层连接。
		_ = c.Close()

		err = c.reDial()
		if err == nil {
			logger.Logger.Debug("reconnected!")
			break
		}
	}
	return true
}

// Channel 在当前底层连接上新建一条通道并封装为 Channel；底层连接由后台断线监听器维护，
// 连接尚未建立或已关闭时调用会失败。极少数情况下通道创建前连接再次断开，也会因此返回网络错误。
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

// trackChannel 返回指定槽位中的复用通道；槽位为空时在当前连接上懒创建一条通道，并登记关闭监听，
// 通道被连接或服务器关闭时自动清空槽位，下次调用重新创建。
//
// PERF-5-13 [P1]：以"单通道 + 互斥串行"方式复用通道，避免每次发布/拉取都承担 channel.open 与
// channel.close 两次 broker 往返；调用期间由互斥锁串行化，保证单通道不被并发使用。
func (c *Connection) trackChannel(mu *sync.Mutex, slot **Channel) (*Channel, error) {
	mu.Lock()
	defer mu.Unlock()
	if *slot != nil {
		return *slot, nil
	}
	ch, err := c.Channel()
	if err != nil {
		return nil, err
	}
	notify := ch.NotifyClose(make(chan *amqp.Error, 1))
	go func() {
		if _, ok := <-notify; ok {
			mu.Lock()
			if *slot == ch {
				*slot = nil
			}
			mu.Unlock()
		}
	}()
	*slot = ch
	return ch, nil
}

// dropChannelOnErr 在遇到可恢复的连接/通道错误时关闭并清空对应复用槽位，下次调用重建通道；
// 非可恢复错误保持槽位不变，交由通道自身的关闭监听处理。
func (c *Connection) dropChannelOnErr(mu *sync.Mutex, slot **Channel, ch *Channel, err error) {
	if !isConnectedErr(err) {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if *slot == ch {
		_ = ch.Close()
		*slot = nil
	}
}

// RetryChannel 在通道创建失败后按传入策略反复重试；重试耗尽仍失败时返回最后一次错误。
// retryable 为 nil 时使用不重试的空实现。
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

// RegisterAndExec 注册一个断线重连后需重新执行的操作，并立即在新通道上执行一次。
//
// 操作执行完后本方法会主动关闭通道，实现内无需手动关闭；操作在后台 goroutine 中执行，方法可能
// 先于操作完成返回。opt 为 nil 时会 panic。
func (c *Connection) RegisterAndExec(opt Operation) {
	if opt == nil {
		panic("Operation must not be nil")
	}
	var key = c.addOperation(opt)
	c.execOperation(key, opt)
}

// execOperation 创建通道失败时立即返回错误；成功则在后台 goroutine 中执行操作并关闭通道。
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

// addOperation 登记一条断线重连后需要重新执行的操作，并返回其唯一标识。
//
// 一般只在此注册“接收消息”类操作：通常无需在每次重连后重新发送已发出的消息。
func (c *Connection) addOperation(opt Operation) (key string) {
	c.oMut.Lock()
	defer c.oMut.Unlock()
	// 标识由自增计数器生成（十六进制字符串），取值空间足以支撑长期运行。
	key = c.genOptKeyFunc()
	c.operations[key] = opt
	return key
}

// exec 遍历执行全部已注册操作，通常用于断线重连成功后恢复消费；每个操作在独立 goroutine 中运行，
// 通道创建失败时记录告警并中止后续操作。
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

// RemoveOperation 异步移除指定标识的已注册操作，之后断线重连不再重新执行该操作。
func (c *Connection) RemoveOperation(key string) {
	go func() {
		c.oMut.Lock()
		defer c.oMut.Unlock()
		delete(c.operations, key)
	}()
}

// Close 关闭底层 AMQP 连接并置空；重复调用安全返回 nil。
func (c *Connection) Close() error {
	// 加锁防止底层连接被并发关闭。
	c.cMut.Lock()
	defer c.cMut.Unlock()

	if c.c != nil {
		defer func() { c.c = nil }()
		return c.c.Close()
	}
	return nil
}

// IsOpen 返回底层连接当前是否存在。
func (c *Connection) IsOpen() bool {
	c.cMut.RLock()
	defer c.cMut.RUnlock()
	return c.c != nil
}

// CanRetry 返回当前重试策略是否尚未放弃重试。
func (c *Connection) CanRetry() bool {
	return !c.retryable.hasGaveUp()
}

// Consumer 返回绑定到当前连接的消费者入口。
func (c *Connection) Consumer() *Consumer {
	return &Consumer{c}
}

// Producer 返回绑定到当前连接的生产者入口。
func (c *Connection) Producer() *Producer {
	return &Producer{c}
}

// QueueBuilder 返回绑定到当前连接的队列构造器。
func (c *Connection) QueueBuilder() *QueueBuilder {
	return NewQueueBuilder(c)
}

// reconnectListener 监听底层连接的关闭事件：收到可恢复错误时尝试重连，重连成功后重新执行
// 全部已注册操作；监听到通道关闭时直接退出。
func (c *Connection) reconnectListener(monitor chan *amqp.Error) {
	err, ok := <-monitor
	if !ok {
		return
	}
	if c.reconnect(err) {
		c.exec()
	}
}

// isAmqpConnectedErr 判断错误是否为 AMQP 层的可恢复连接错误（连接被关闭、EOF 或被强制断开）。
func isAmqpConnectedErr(err error) bool {
	var amqpErr *amqp.Error
	return errors.As(err, &amqpErr) &&
		(errors.Is(amqpErr, amqp.ErrClosed) ||
			amqpErr.Reason == "EOF" ||
			amqpErr.Code == amqp.ConnectionForced)
}

// isTransportNetError 判断错误是否为传输层网络错误（超时、DNS 解析失败、连接被拒绝或连接超时）。
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

// isConnectedErr 判断错误是否属于可触发重连的连接错误（传输层网络错误或 AMQP 层连接错误）。
func isConnectedErr(err error) bool {
	return isTransportNetError(err) || isAmqpConnectedErr(err)
}
