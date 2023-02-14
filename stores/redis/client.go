package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/util"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// Config 配置
type Config struct {
	Type          bool     //是否集群
	Hosts         []string //IP
	Password      string   //密码
	Database      int      //数据库
	PoolSize      int      //连接池大小
	KeyPrefix     string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
}

// New 新连接
func New(config *Config) *Client {
	opts := Options{}
	if config.Type {
		opts.Type = ClientCluster
	} else {
		opts.Type = ClientNormal
	}
	opts.Hosts = config.Hosts
	opts.KeyPrefix = config.KeyPrefix

	if config.PoolSize > 0 {
		opts.PoolSize = config.PoolSize
	} else {
		opts.PoolSize = 64
	}
	if config.Database > 0 {
		opts.Database = config.Database
	} else {
		opts.Database = 0
	}
	if config.Password != "" {
		opts.Password = config.Password
	}
	if !config.DisableMetric {
		opts.DisableMetric = config.DisableMetric
	}
	if !config.DisableTrace {
		opts.DisableTrace = config.DisableTrace
	}
	client := NewClient(opts)
	ctx := context.TODO()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Logger.Panic(err.Error())
	}
	return client
}

// RedisNil means nil reply, .e.g. when key does not exist.
const RedisNil = redis.Nil

// Client a struct representing the redis client
type Client struct {
	opts          Options
	client        *redis.Client
	clusterClient *redis.ClusterClient
	fmtString     string
	clientType    ClientType
}

// NewClient 新客户端
func NewClient(opts Options) *Client {
	r := &Client{opts: opts}
	switch opts.Type {
	// 群集客户端
	case ClientCluster:
		// NewClusterClient 返回一个 Redis 集群客户端
		tc := redis.NewClusterClient(opts.GetClusterConfig())
		ctx := context.TODO()
		if !opts.DisableTrace || !opts.DisableMetric {
			_ = tc.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
				shard.AddHook(OpenTraceHook{
					DisableMetric: opts.DisableMetric,
					DisableTrace:  opts.DisableTrace,
					DB:            opts.Database,
					Addr:          util.Implode(";", opts.Hosts),
				})
				return nil
			})
		}
		r.clusterClient = tc
		r.clientType = ClientCluster
	// 标准客户端也是默认值
	case ClientNormal:
		fallthrough
	default:
		// NewClient 根据 Options 指定的 Redis Server 返回一个客户端。
		tc := redis.NewClient(opts.GetNormalConfig())
		if !opts.DisableTrace || !opts.DisableMetric {
			tc.AddHook(OpenTraceHook{
				DisableMetric: opts.DisableMetric,
				DisableTrace:  opts.DisableTrace,
				DB:            opts.Database,
				Addr:          util.Implode(";", opts.Hosts),
			})
		}
		r.client = tc
		r.clientType = ClientNormal
	}
	r.fmtString = opts.KeyPrefix + "%s"

	return r
}

// Prefix 返回前缀+键
func (r *Client) Prefix(key string) string {
	return fmt.Sprintf(r.fmtString, key)
}

// k 格式化并返回带前缀的密钥
func (r *Client) k(key string) string {
	return fmt.Sprintf(r.fmtString, key)
}

// ks 使用前缀格式化并返回一组键
func (r *Client) ks(key ...string) []string {
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = r.k(k)
	}
	return keys
}

// GetClient 返回客户端
func (r *Client) GetClient() (res redis.Cmdable) {
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient
	case ClientNormal:
		res = r.client
	}
	return res
}

// MGetByPipeline gets multiple values from keys,Pipeline is used when
// redis is a cluster,This means higher IO performance
// params: keys ...string
// return: []string, error
func (r *Client) MGetByPipeline(ctx context.Context, keys ...string) ([]string, error) {

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	var res []string

	if r.clientType == ClientCluster {
		start := time.Now()
		pipeLineLen := 100
		pipeCount := len(keys)/pipeLineLen + 1
		pipes := make([]redis.Pipeliner, pipeCount)
		for i := 0; i < pipeCount; i++ {
			pipes[i] = r.client.Pipeline()
		}
		for i, k := range keys {
			p := pipes[i%pipeCount]
			p.Get(ctx, r.k(k))
		}
		logger.Logger.Debug("process cost: ", time.Since(start))
		start = time.Now()
		var wg sync.WaitGroup
		var lock sync.Mutex
		errors := make(chan error, pipeCount)
		for _, p := range pipes {
			p := p
			wg.Add(1)
			go func() {
				defer wg.Done()
				cmders, err := p.Exec(ctx)
				if err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}
				lock.Lock()
				defer lock.Unlock()
				for _, cmder := range cmders {
					result, _ := cmder.(*redis.StringCmd).Result()
					res = append(res, result)
				}
			}()
		}
		wg.Wait()
		logger.Logger.Debug("exec cost: ", time.Since(start))
		if len(errors) > 0 {
			return nil, <-errors
		}

		return res, nil
	}

	vals, err := r.client.MGet(ctx, keys...).Result()

	if redis.Nil != err && nil != err {
		return nil, err
	}

	for _, item := range vals {
		res = append(res, fmt.Sprintf("%s", item))
	}

	return res, err
}

// ErrNotImplemented not implemented error
var ErrNotImplemented = errors.New("Not implemented")
