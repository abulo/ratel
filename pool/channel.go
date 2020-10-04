package pool

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// PoolConfig 连接池相关配置
type PoolConfig struct {
	//是否初始化连接池
	// IsInit bool
	//最大链接数
	MaxOpenConns int
	//最大闲置链接
	MaxIdleConns int
	//生成连接的方法
	New func() (interface{}, error)
	//关闭连接的方法
	Close func(interface{}) error
	//检查连接是否有效的方法
	Ping func(interface{}) error
	//连接最大空闲时间，超过该事件则将失效
	ConnTimeout time.Duration
	//配置参数
	Config interface{}
}

// channelPool 存放连接信息
type channelPool struct {
	mutex              sync.Mutex
	conns              chan *Conn
	new                func() (interface{}, error)
	close              func(interface{}) error
	ping               func(interface{}) error
	connTimeout        time.Duration
	maxOpenConns       int
	overstepConnNumber int
	config             interface{}
}

type Conn struct {
	conn interface{}
	t    time.Time
}

// NewChannelPool 初始化连接
func NewChannelPool(poolConfig *PoolConfig) (Pool, error) {
	if poolConfig.MaxIdleConns < 0 || poolConfig.MaxOpenConns <= 0 || poolConfig.MaxIdleConns > poolConfig.MaxOpenConns {
		return nil, errors.New("invalid capacity settings")
	}
	if poolConfig.New == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}

	c := &channelPool{
		conns:       make(chan *Conn, poolConfig.MaxIdleConns),
		new:         poolConfig.New,
		close:       poolConfig.Close,
		connTimeout: poolConfig.ConnTimeout,
		config:      poolConfig.Config,
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}
	// if poolConfig.IsInit {
	// 	for i := 0; i < poolConfig.MaxOpenConns; i++ {
	// 		conn, err := c.new()
	// 		if err != nil {
	// 			c.CloseAll()
	// 			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
	// 		}
	// 		c.conns <- &Conn{conn: conn, t: time.Now()}
	// 	}
	// }
	return c, nil
}

// getConns 获取所有连接
func (c *channelPool) getConns() chan *Conn {
	c.mutex.Lock()
	conns := c.conns
	c.mutex.Unlock()
	return conns
}

// Get 从pool中取一个连接
func (c *channelPool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.connTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.Close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					fmt.Println("conn is not able to be connected: ", err)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			if c.maxOpenConns > (c.Len() + c.overstepConnNumber) {
				return <-conns, nil
			}
			c.mutex.Lock()
			if c.new == nil {
				c.mutex.Unlock()
				continue
			}
			conn, err := c.new()
			c.mutex.Unlock()
			if err != nil {
				return nil, err
			}
			c.overstepConnNumber++
			return conn, nil
		}
	}
}

// Put 将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mutex.Lock()

	if c.conns == nil {
		c.mutex.Unlock()
		return c.Close(conn)
	}

	select {
	case c.conns <- &Conn{conn: conn, t: time.Now()}:
		c.mutex.Unlock()
		return nil
	default:
		c.mutex.Unlock()
		c.overstepConnNumber--
		//连接池已满，直接关闭该连接
		return c.Close(conn)
	}
}

// Close 关闭单条连接
func (c *channelPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.close == nil {
		return nil
	}
	return c.close(conn)
}

// Ping 检查单条连接是否有效
func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.ping(conn)
}

// Release 释放连接池中所有连接
func (c *channelPool) CloseAll() {
	c.mutex.Lock()
	conns := c.conns
	c.conns = nil
	c.new = nil
	c.ping = nil
	closeFun := c.close
	c.close = nil
	c.mutex.Unlock()

	if conns == nil {
		return
	}
	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

// Len 连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}
