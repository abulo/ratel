package redis

import (
	"fmt"
	"io"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/resource"
	"github.com/abulo/ratel/v3/util"
	"github.com/redis/go-redis/v9"
)

// // Client  配置
// type Client struct {
// 	ClientType          string            // 模式(normal => 单节点,cluster =>  集群,failover => 哨兵,ring => 分片)
// 	Hosts         []string          // 集群 哨兵 需要填写
// 	Password      string            // 密码
// 	Database      int               // 数据库
// 	PoolSize      int               // 连接池大小
// 	KeyPrefix     string            // 前缀标识
// 	DisableMetric bool              // 关闭指标采集
// 	DisableTrace  bool              // 关闭链路追踪
// 	Addr          string            // 单节点客户端
// 	Addrs         map[string]string // 分片客户端  shardName => host:port
// 	MasterName    string            // 哨兵
// }

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

// WithClientType customizes the given Redis with given Type.
func WithClientType(ClientType string) Option {
	return func(r *Client) {
		r.ClientType = ClientType
	}
}

// WithHosts customizes the given Redis with given Hosts.
func WithHosts(Hosts []string) Option {
	return func(r *Client) {
		r.Hosts = Hosts
	}
}

// WithPassword customizes the given Redis with given Password.
func WithPassword(Password string) Option {
	return func(r *Client) {
		r.Password = Password
	}
}

// WithDatabase customizes the given Redis with given Database.
func WithDatabase(Database int) Option {
	return func(r *Client) {
		r.Database = Database
	}
}

// WithPoolSize customizes the given Redis with given PoolSize.
func WithPoolSize(PoolSize int) Option {
	return func(r *Client) {
		r.PoolSize = PoolSize
	}
}

// WithKeyPrefix customizes the given Redis with given KeyPrefix.
func WithKeyPrefix(KeyPrefix string) Option {
	return func(r *Client) {
		r.KeyPrefix = KeyPrefix + "%s"
	}
}

// WithDisableMetric customizes the given Redis with given DisableMetric.
func WithDisableMetric(DisableMetric bool) Option {
	return func(r *Client) {
		r.DisableMetric = DisableMetric
	}
}

// WithDisableTrace customizes the given Redis with given DisableTrace.
func WithDisableTrace(DisableTrace bool) Option {
	return func(r *Client) {
		r.DisableTrace = DisableTrace
	}
}

// WithAddr customizes the given Redis with given Addr.
func WithAddr(Addr string) Option {
	return func(r *Client) {
		r.Addr = Addr
	}
}

// WithAddrs customizes the given Redis with given Addrs.
func WithAddrs(Addrs map[string]string) Option {
	return func(r *Client) {
		r.Addrs = Addrs
	}
}

// WithMasterName customizes the given Redis with given MasterName.
func WithMasterName(MasterName string) Option {
	return func(r *Client) {
		r.MasterName = MasterName
	}
}

// RedisNode interface represents a redis node.
type RedisNode interface {
	redis.UniversalClient
	redis.Cmdable
	redis.BitMapCmdable
	// FTList(ctx context.Context) *redis.StringSliceCmd
}

// getRedis new redis client
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

// getClient new redis client
func getClient(r *Client) (RedisNode, error) {
	driverName := r.Addr
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetClientConfig()
		store := redis.NewClient(opt)
		if !r.DisableTrace || !r.DisableMetric {
			store.AddHook(OpenTraceHook{
				DisableMetric: r.DisableMetric,
				DisableTrace:  r.DisableTrace,
				DB:            r.Database,
				Addr:          driverName,
			})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// var clusterManager = resource.NewResourceManager()

// getCluster new redis  cluster client
func getCluster(r *Client) (RedisNode, error) {
	driverName := util.Implode(";", r.Hosts)
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetClusterClientConfig()
		store := redis.NewClusterClient(opt)
		if !r.DisableTrace || !r.DisableMetric {
			store.AddHook(OpenTraceHook{
				DisableMetric: r.DisableMetric,
				DisableTrace:  r.DisableTrace,
				DB:            r.Database,
				Addr:          driverName,
			})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// var failoverManager = resource.NewResourceManager()

// getFailover new redis  failover client
func getFailover(r *Client) (RedisNode, error) {
	driverName := r.MasterName + "://" + util.Implode(";", r.Hosts)
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetFailoverClientConfig()
		store := redis.NewFailoverClient(opt)
		if !r.DisableTrace || !r.DisableMetric {
			store.AddHook(OpenTraceHook{
				DisableMetric: r.DisableMetric,
				DisableTrace:  r.DisableTrace,
				DB:            r.Database,
				Addr:          driverName,
			})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Client), nil
}

// var ringManager = resource.NewResourceManager()

// getRing new redis  ring client
func getRing(r *Client) (RedisNode, error) {
	var driverName string
	for k, v := range r.Addrs {
		driverName += k + ":" + v + ";"
	}
	val, err := clientManager.GetResource(driverName, func() (io.Closer, error) {
		opt := r.GetRingClientConfig()
		store := redis.NewRing(opt)
		if !r.DisableTrace || !r.DisableMetric {
			store.AddHook(OpenTraceHook{
				DisableMetric: r.DisableMetric,
				DisableTrace:  r.DisableTrace,
				DB:            r.Database,
				Addr:          driverName,
			})
		}
		return store, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(*redis.Ring), nil
}
