package redis

import (
	"crypto/tls"
	"time"

	"github.com/abulo/ratel/v3/core/resource"
	"github.com/redis/go-redis/v9"
)

type Option func(r *Client)

// Options options to initiate your client
type Client struct {
	ClientType string // 模式(normal => 单节点,cluster =>  集群,failover => 哨兵,ring => 分片)
	// Host address with port number
	// For normal client will only used the first value
	// A seed list of host:port addresses of sentinel nodes.
	Hosts []string // 集群 哨兵 需要填写
	// host:port address.
	Addr string // 单节点客户端
	// Map of name => host:port addresses of ring shards.
	Addrs map[string]string // 分片客户端  shardName => host:port
	// The master name.
	MasterName string
	// The network type, either tcp or unix.
	// Default is tcp.
	// Only for normal client
	Network string
	// Database to be selected after connecting to the server.
	Database int // 数据库
	// Automatically adds a prefix to all keys
	KeyPrefix string // 前缀标识
	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 16.
	// In normal client this is the MaxRetries option
	MaxRedirects int
	// Enables read queries for a connection to a Redis Cluster slave node.
	IsReadOnly bool
	// Enables routing read-only queries to the closest master or slave node.
	// If set will change this client to read-only mode
	RouteByLatency bool
	// Following options are copied from Options struct.
	Password string // 密码
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 3 seconds.
	WriteTimeout time.Duration
	// PoolSize applies per cluster node and not for the whole cluster.
	// Maximum number of socket connections.
	// Default is 10 connections.
	PoolSize int // 连接池大小
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is to not close idle connections.
	IdleTimeout time.Duration
	// Frequency of idle checks.
	// Default is 1 minute.
	// When minus value is set, then idle check is disabled.
	IdleCheckFrequency time.Duration
	// TLS Config to use. When set TLS will be negotiated.
	// Only for normal client
	TLSConfig     *tls.Config
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	brk           resource.Breaker
}

// GetRingClientConfig 获取分片配置
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
	}
	return opts
}

// GetFailoverClient 获取哨兵配置
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
	}
	return opts
}

// GetClusterClientConfig 获取集群配置
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
	}
	if o.MaxRedirects > 0 {
		opts.MaxRedirects = o.MaxRedirects
	}
	return opts
}

// GetClientConfig 获取单节点配置
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
	}
	return opts
}

// ClientType type to define a redis client connector

// ClientNormal ...
const (
	// ClientNormal for standard instance client
	ClientNormal = "normal"
	// ClientCluster for official redis cluster
	ClientCluster = "cluster"
	// FailoverClient for official redis failover
	ClientFailover = "failover"
	// RingClient for official redis ring
	ClientRing = "ring"
)

// RWType Client Reader and Writer
type RWType string

// OnlyRead ...
const (
	// OnlyRead serves as a search suffix for configuration parameters
	OnlyRead RWType = "READER"
	// OnlyWrite serves as a search suffix for configuration parameters
	OnlyWrite RWType = "WRITER"
	// ReadAndWrite serves as a search suffix for configuration parameters
	ReadAndWrite RWType = ""
)

// IsReadOnly will return Is it read-only
func (rw *RWType) IsReadOnly() bool {
	return *rw == OnlyRead
}

// FmtSuffix get fmtstring of  key+ "_" + suffix
func (rw *RWType) FmtSuffix(key string) string {
	if *rw == ReadAndWrite {
		return key
	}
	return key + "_" + string(*rw)
}
