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

// NewRedisClient 根据给定的 Option 构造 Redis 客户端配置 Client。
//
// 该方法只完成配置装配与模式校验，并不立即建立底层连接；实际连接在首次执行命令时
// 按 ClientType 懒加载并复用。未指定 ClientType 时按单节点处理，PoolSize 为 0 时
// 默认设为 64；分片模式缺少 Addrs、哨兵模式缺少 MasterName、集群模式缺少 Hosts 时
// 返回对应错误。
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

// WithClientType 设置客户端连接模式，取值 normal、cluster、failover、ring，
// 对应 ClientNormal、ClientCluster、ClientFailover、ClientRing 常量。
func WithClientType(ClientType string) Option {
	return func(r *Client) {
		r.ClientType = ClientType
	}
}

// WithHosts 设置节点地址列表，格式 host:port；集群与哨兵模式必填。
func WithHosts(Hosts []string) Option {
	return func(r *Client) {
		r.Hosts = Hosts
	}
}

// WithPassword 设置连接认证密码。
func WithPassword(Password string) Option {
	return func(r *Client) {
		r.Password = Password
	}
}

// WithDatabase 设置连接后选中的数据库编号。
func WithDatabase(Database int) Option {
	return func(r *Client) {
		r.Database = Database
	}
}

// WithPoolSize 设置连接池最大连接数；为 0 时由 NewRedisClient 默认为 64。
func WithPoolSize(PoolSize int) Option {
	return func(r *Client) {
		r.PoolSize = PoolSize
	}
}

// WithEnableMetric 开关是否采集 Redis 操作指标。
func WithEnableMetric(EnableMetric bool) Option {
	return func(r *Client) {
		r.EnableMetric = EnableMetric
	}
}

// WithEnableTrace 开关是否为命令挂载链路追踪钩子。
func WithEnableTrace(EnableTrace bool) Option {
	return func(r *Client) {
		r.EnableTrace = EnableTrace
	}
}

// WithAddr 设置单节点模式下的 host:port 地址。
func WithAddr(Addr string) Option {
	return func(r *Client) {
		r.Addr = Addr
	}
}

// WithName 设置客户端节点名称，参与底层连接实例的去重键。
func WithName(Name string) Option {
	return func(r *Client) {
		r.Name = Name
	}
}

// WithAddrs 设置分片模式下分片名称到 host:port 地址的映射。
func WithAddrs(Addrs map[string]string) Option {
	return func(r *Client) {
		r.Addrs = Addrs
	}
}

// WithMasterName 设置哨兵模式下的主节点名称，哨兵模式必填。
func WithMasterName(MasterName string) Option {
	return func(r *Client) {
		r.MasterName = MasterName
	}
}

// WithDialTimeout 设置建立新连接的超时时间。
func WithDialTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.DialTimeout = t
	}
}

// WithReadTimeout 设置套接字读取超时时间，超时后命令失败而非阻塞。
func WithReadTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.ReadTimeout = t
	}
}

// WithWriteTimeout 设置套接字写入超时时间，超时后命令失败而非阻塞。
func WithWriteTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.WriteTimeout = t
	}
}

// WithPoolTimeout 设置连接池耗尽时等待连接的超时时间。
func WithPoolTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.PoolTimeout = t
	}
}

// WithIdleTimeout 设置关闭空闲连接的等待时长。
func WithIdleTimeout(t time.Duration) Option {
	return func(r *Client) {
		r.IdleTimeout = t
	}
}

// WithIdleCheckFrequency 设置空闲连接检查频率，传入负值可关闭检查。
func WithIdleCheckFrequency(t time.Duration) Option {
	return func(r *Client) {
		r.IdleCheckFrequency = t
	}
}

// RedisNode 是各连接模式底层实例的统一命令接口，组合 go-redis 的 UniversalClient
// 与 Cmdable，使上层方法无需区分单节点、集群、哨兵与分片实例类型。
type RedisNode interface {
	redis.UniversalClient
	redis.Cmdable
}

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
