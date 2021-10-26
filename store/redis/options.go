package redis

import (
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
)

// Options options to initiate your client
type Options struct {
	Type ClientType
	// Host address with port number
	// For normal client will only used the first value
	Hosts []string

	// The network type, either tcp or unix.
	// Default is tcp.
	// Only for normal client
	Network string

	// Database to be selected after connecting to the server.
	Database int
	// Automatically adds a prefix to all keys
	KeyPrefix string

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 16.
	// In normal client this is the MaxRetries option
	MaxRedirects int

	// Enables read queries for a connection to a Redis Cluster slave node.
	ReadOnly bool

	// Enables routing read-only queries to the closest master or slave node.
	// If set will change this client to read-only mode
	RouteByLatency bool

	// Following options are copied from Options struct.
	Password string

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
	PoolSize int
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
	TLSConfig *tls.Config

	Trace bool
}

// GetClusterConfig translates current configuration into a *redis.ClusterOptions
func (o Options) GetClusterConfig() *redis.ClusterOptions {
	opts := &redis.ClusterOptions{
		Addrs:              o.Hosts,
		ReadOnly:           o.ReadOnly,
		RouteByLatency:     o.RouteByLatency,
		Password:           o.Password,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
	}
	if o.MaxRedirects > 0 {
		opts.MaxRedirects = o.MaxRedirects
	}
	return opts
}

// GetNormalConfig translates current configuration into a *redis.Options struct
func (o Options) GetNormalConfig() *redis.Options {
	opts := &redis.Options{
		Addr:               o.Hosts[0],
		Password:           o.Password,
		DB:                 o.Database,
		MaxRetries:         o.MaxRedirects,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
	return opts
}

// ClientType type to define a redis client connector
type ClientType string

const (
	// ClientNormal for standard instance client
	ClientNormal ClientType = "normal"
	// ClientCluster for official redis cluster
	ClientCluster ClientType = "cluster"
)

// Client Reader and Writer
type RWType string

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

	if *rw == OnlyRead {
		return true
	} else {
		return false
	}
}

// FmtSuffix get fmtstring of  key+ "_" + suffix
func (rw *RWType) FmtSuffix(key string) string {
	if *rw == ReadAndWrite {
		return key
	}
	return key + "_" + string(*rw)
}
