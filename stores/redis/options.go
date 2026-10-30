package redis

import (
	"crypto/tls"
	"time"

	"github.com/abulo/ratel/v3/core/resource"
	"github.com/redis/go-redis/v9"
)

// Option 是 Client 的配置选项函数，以函数式选项模式在 NewRedisClient 中初始化客户端。
type Option func(r *Client)

// Client 保存 Redis 客户端的全部配置项，支持单节点、集群、哨兵与分片四种模式。
//
// 零值不可直接使用，须通过 NewRedisClient 结合 Option 初始化并完成模式校验；
// Client 为配置载体，应在初始化阶段一次性写入，不保证并发安全。
type Client struct {
	// ClientType 指定客户端模式：normal 为单节点，cluster 为集群，
	// failover 为哨兵，ring 为分片；为空时默认按单节点处理。
	ClientType string
	// Hosts 为节点地址列表，格式 host:port；集群与哨兵模式必填，单节点模式仅取第一个值。
	Hosts []string
	// Addr 为单节点模式下的 host:port 地址。
	Addr string
	// Addrs 为分片模式下分片名称到 host:port 地址的映射。
	Addrs map[string]string
	// MasterName 为哨兵模式下的主节点名称，必填。
	MasterName string
	// Network 为网络类型，取值 tcp 或 unix，仅单节点模式生效，默认 tcp。
	Network string
	// Database 为连接后选中的数据库编号。
	Database int
	// MaxRedirects 在单节点模式下作为最大重试次数，在集群模式下作为最大重定向次数。
	MaxRedirects int
	// IsReadOnly 允许将读请求路由到集群从节点。
	IsReadOnly bool
	// RouteByLatency 将只读请求路由到延迟最低的主节点或从节点，开启后客户端进入只读模式。
	RouteByLatency bool
	// Password 为连接密码。
	Password string
	// DialTimeout 为建立新连接的超时时间，默认 5 秒。
	DialTimeout time.Duration
	// ReadTimeout 为套接字读取超时时间，超时后命令失败而非阻塞，默认 3 秒。
	ReadTimeout time.Duration
	// WriteTimeout 为套接字写入超时时间，超时后命令失败而非阻塞，默认 3 秒。
	WriteTimeout time.Duration
	// PoolSize 为连接池中最大的套接字连接数，集群模式下按节点计；为 0 时由 NewRedisClient 设为 64。
	PoolSize int
	// PoolTimeout 为连接池耗尽时的等待超时时间，默认 ReadTimeout 加 1 秒。
	PoolTimeout time.Duration
	// IdleTimeout 为关闭空闲连接的等待时长，应小于服务端超时时间，默认不关闭空闲连接。
	IdleTimeout time.Duration
	// IdleCheckFrequency 为空闲连接检查频率，默认 1 分钟；设置为负值可关闭检查。
	IdleCheckFrequency time.Duration
	// TLSConfig 为 TLS 配置，设置后启用 TLS 加密，仅单节点模式生效。
	TLSConfig *tls.Config
	// EnableMetric 是否开启指标采集。
	EnableMetric bool
	// EnableTrace 是否开启链路追踪。
	EnableTrace bool
	// brk 为熔断器。
	brk resource.Breaker
	// Name 为节点名称。
	Name string
}

// GetRingClientConfig 将 Client 配置转换为分片（ring）模式的 redis.RingOptions。
func (o *Client) GetRingClientConfig() *redis.RingOptions {
	opts := &redis.RingOptions{
		Addrs:           o.Addrs,
		Password:        o.Password,
		DialTimeout:     o.DialTimeout,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.IdleTimeout,
		ConnMaxLifetime: o.IdleCheckFrequency,
		DisableIdentity: true,
		Protocol:        2,
	}
	return opts
}

// GetFailoverClientConfig 将 Client 配置转换为哨兵（failover）模式的 redis.FailoverOptions。
func (o *Client) GetFailoverClientConfig() *redis.FailoverOptions {
	opts := &redis.FailoverOptions{
		MasterName:      o.MasterName,
		SentinelAddrs:   o.Hosts,
		RouteByLatency:  o.RouteByLatency,
		Password:        o.Password,
		DialTimeout:     o.DialTimeout,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.IdleTimeout,
		ConnMaxLifetime: o.IdleCheckFrequency,
		DisableIdentity: true,
		Protocol:        2,
	}
	return opts
}

// GetClusterClientConfig 将 Client 配置转换为集群（cluster）模式的 redis.ClusterOptions。
func (o *Client) GetClusterClientConfig() *redis.ClusterOptions {
	opts := &redis.ClusterOptions{
		Addrs:           o.Hosts,
		ReadOnly:        o.IsReadOnly,
		RouteByLatency:  o.RouteByLatency,
		Password:        o.Password,
		DialTimeout:     o.DialTimeout,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.IdleTimeout,
		ConnMaxLifetime: o.IdleCheckFrequency,
		DisableIdentity: true,
		Protocol:        2,
	}
	if o.MaxRedirects > 0 {
		opts.MaxRedirects = o.MaxRedirects
	}
	return opts
}

// GetClientConfig 将 Client 配置转换为单节点（normal）模式的 redis.Options。
func (o *Client) GetClientConfig() *redis.Options {
	opts := &redis.Options{
		Addr:            o.Addr,
		Password:        o.Password,
		DB:              o.Database,
		MaxRetries:      o.MaxRedirects,
		DialTimeout:     o.DialTimeout,
		ReadTimeout:     o.ReadTimeout,
		WriteTimeout:    o.WriteTimeout,
		PoolSize:        o.PoolSize,
		PoolTimeout:     o.PoolTimeout,
		ConnMaxIdleTime: o.IdleTimeout,
		ConnMaxLifetime: o.IdleCheckFrequency,
		TLSConfig:       o.TLSConfig,
		DisableIdentity: true,
		Protocol:        2,
	}
	return opts
}

// ClientType 类型常量，定义 Redis 客户端的连接模式。
const (
	// ClientNormal 为标准单节点实例客户端。
	ClientNormal = "normal"
	// ClientCluster 为官方 Redis 集群客户端。
	ClientCluster = "cluster"
	// ClientFailover 为官方 Redis 哨兵（failover）客户端。
	ClientFailover = "failover"
	// ClientRing 为官方 Redis 分片（ring）客户端。
	ClientRing = "ring"
)

// RWType 表示客户端读写类型，用于拼接配置参数与键的后缀。
type RWType string

// 读写类型常量，作为配置参数搜索后缀使用。
const (
	// OnlyRead 仅读模式，后缀为 "READER"。
	OnlyRead RWType = "READER"
	// OnlyWrite 仅写模式，后缀为 "WRITER"。
	OnlyWrite RWType = "WRITER"
	// ReadAndWrite 读写模式，不追加后缀。
	ReadAndWrite RWType = ""
)

// IsReadOnly 返回当前读写类型是否为仅读模式。
func (rw *RWType) IsReadOnly() bool {
	return *rw == OnlyRead
}

// FmtSuffix 返回 key 拼接读写类型后缀的结果；读写模式不拼接，直接返回原 key。
func (rw *RWType) FmtSuffix(key string) string {
	if *rw == ReadAndWrite {
		return key
	}
	return key + "_" + string(*rw)
}
