package redis

import (
	"fmt"
	"io"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/resource"
	"github.com/abulo/ratel/v3/util"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 创建新的Redis客户端 opts: 配置选项
func NewRedisClient(opts ...Option) (*Client, error) {
	c := &Client{}
	c.brk = resource.NewBreaker()
	for _, o := range opts {
		o(c)
	}
	if c.ClientType == "" {
		c.ClientType = ClientNormal
	}
	if c.PoolSize == 0 {
		c.PoolSize = 64
	}
	if c.Database == 0 {
		c.Database = 0
	}
	if c.ClientType == ClientRing {
		if len(c.Addrs) == 0 {
			return nil, fmt.Errorf("redis ring client addrs is empty")
		}
	}
	if c.ClientType == ClientFailover {
		if c.MasterName == "" {
			return nil, fmt.Errorf("redis failover client master name is empty")
		}
	}
	if c.ClientType == ClientCluster {
		if len(c.Hosts) == 0 {
			return nil, fmt.Errorf("redis cluster client hosts is empty")
		}
	}
	return c, nil
}

// WithClientType 设置Redis客户端类型 ClientType: 客户端类型(normal,cluster,failover,ring)
func WithClientType(ClientType string) Option {
	return func(r *Client) {
		r.ClientType = ClientType
	}
}

// WithHosts 设置Redis主机地址 Hosts: 主机地址列表
func WithHosts(Hosts []string) Option {
	return func(r *Client) {
		r.Hosts = Hosts
	}
}

// WithPassword 设置Redis密码 Password: 认证密码
func WithPassword(Password string) Option {
	return func(r *Client) {
		r.Password = Password
	}
}

// WithDatabase 设置Redis数据库 Database: 数据库编号
func WithDatabase(Database int) Option {
	return func(r *Client) {
		r.Database = Database
	}
}

// WithPoolSize 设置连接池大小 PoolSize: 连接池大小
func WithPoolSize(PoolSize int) Option {
	return func(r *Client) {
		r.PoolSize = PoolSize
	}
}

// WithEnableMetric 禁用指标采集 EnableMetric: 是否禁用指标采集
func WithEnableMetric(EnableMetric bool) Option {
	return func(r *Client) {
		r.EnableMetric = EnableMetric
	}
}

// WithEnableTrace 禁用链路追踪 EnableTrace: 是否禁用链路追踪
func WithEnableTrace(EnableTrace bool) Option {
	return func(r *Client) {
		r.EnableTrace = EnableTrace
	}
}

// WithAddr 设置单节点地址 Addr: 单节点地址
func WithAddr(Addr string) Option {
	return func(r *Client) {
		r.Addr = Addr
	}
}

func WithName(Name string) Option {
	return func(r *Client) {
		r.Name = Name
	}
}

// WithAddrs 设置分片地址 Addrs: 分片地址映射
func WithAddrs(Addrs map[string]string) Option {
	return func(r *Client) {
		r.Addrs = Addrs
	}
}

// WithMasterName 设置哨兵主节点名称 MasterName: 主节点名称
func WithMasterName(MasterName string) Option {
	return func(r *Client) {
		r.MasterName = MasterName
	}
}
func WithDialTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.DialTimeout = t
	}
}
func WithReadTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.WriteTimeout = t
	}
}

func WithPoolTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.PoolTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.IdleTimeout = t
	}
}

func WithIdleCheckFrequency(t time.Duration) Option {
	return func(r *Client) {
		r.IdleCheckFrequency = t
	}
}

// RedisNode interface represents a redis node.
type RedisNode interface {
	redis.UniversalClient
	redis.Cmdable
}

// getRedis 获取Redis客户端实例 r: 客户端配置
func getRedis(r *Client) (RedisNode, error) {
	switch r.ClientType {
	case ClientNormal:
		return getClient(r)
	case ClientCluster:
		return getCluster(r)
	case ClientFailover:
		return getFailover(r)
	case ClientRing:
		return getRing(r)
	default:
		err := fmt.Errorf("redis type '%s' is not supported", r.ClientType)
		logger.Logger.Panic(err)
		return nil, err
	}
}

var clientManager = resource.NewResourceManager()

// getClient 获取单节点Redis客户端 r: 客户端配置
func getClient(r *Client) (RedisNode, error) {
	driverName := r.Addr + "@" + r.Name
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetClientConfig()
		store := redis.NewClient(opt)
		if r.EnableTrace {
			store.AddHook(OpenTraceHook{})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// getCluster 获取集群Redis客户端 r: 客户端配置
func getCluster(r *Client) (RedisNode, error) {
	driverName := util.Implode(";", r.Hosts) + "@" + r.Name
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetClusterClientConfig()
		store := redis.NewClusterClient(opt)
		if r.EnableTrace {
			store.AddHook(OpenTraceHook{})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// getFailover 获取哨兵Redis客户端 r: 客户端配置
func getFailover(r *Client) (RedisNode, error) {
	driverName := r.MasterName + "://" + util.Implode(";", r.Hosts) + "@" + r.Name
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetFailoverClientConfig()
		store := redis.NewFailoverClient(opt)
		if r.EnableTrace {
			store.AddHook(OpenTraceHook{})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// getRing 获取分片Redis客户端 r: 客户端配置
func getRing(r *Client) (RedisNode, error) {
	var driverName string
	for k, v := range r.Addrs {
		driverName += k + ":" + v + ";"
	}
	driverName = driverName + "@" + r.Name
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetRingClientConfig()
		store := redis.NewRing(opt)
		if r.EnableTrace {
			store.AddHook(OpenTraceHook{})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Ring), nil
}
