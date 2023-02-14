package grpc

import (
	"time"

	"github.com/abulo/ratel/core/constant"
	"github.com/abulo/ratel/core/logger"
	"github.com/abulo/ratel/core/singleton"
	"github.com/abulo/ratel/registry/etcdv3"
	"github.com/abulo/ratel/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	Name                      string
	BalancerName              string
	Address                   string
	Block                     bool
	DialTimeout               time.Duration
	ReadTimeout               time.Duration
	Direct                    bool
	KeepAlive                 *keepalive.ClientParameters
	dialOptions               []grpc.DialOption
	SlowThreshold             time.Duration
	Debug                     bool
	DisableTraceInterceptor   bool
	DisableAidInterceptor     bool
	DisableTimeoutInterceptor bool
	DisableMetricInterceptor  bool
	DisableAccessInterceptor  bool
	AccessInterceptorLevel    string
	Etcd                      *etcdv3.Config
}

// New ...
func New() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 3,
		ReadTimeout:            util.Duration("1s"),
		SlowThreshold:          util.Duration("600ms"),
		AccessInterceptorLevel: "info",
		Block:                  true,
	}
}

func (config *Config) SetEtcd(option *etcdv3.Config) *Config {
	config.Etcd = option
	return config
}

// SetName ...
func (config *Config) SetName(Name string) *Config {
	config.Name = Name
	return config
}

// SetBalancerName ...
func (config *Config) SetBalancerName(BalancerName string) *Config {
	config.BalancerName = BalancerName
	return config
}

// SetAddress ...
func (config *Config) SetAddress(Address string) *Config {
	config.Address = Address
	return config
}

// SetBlock ...
func (config *Config) SetBlock(Block bool) *Config {
	config.Block = Block
	return config
}

// SetDialTimeout ...
func (config *Config) SetDialTimeout(DialTimeout time.Duration) *Config {
	config.DialTimeout = DialTimeout
	return config
}

// SetReadTimeout ...
func (config *Config) SetReadTimeout(ReadTimeout time.Duration) *Config {
	config.ReadTimeout = ReadTimeout
	return config
}

// SetDirect ...
func (config *Config) SetDirect(Direct bool) *Config {
	config.Direct = Direct
	return config
}

// SetSlowThreshold ...
func (config *Config) SetSlowThreshold(SlowThreshold time.Duration) *Config {
	config.SlowThreshold = SlowThreshold
	return config
}

// SetDebug ...
func (config *Config) SetDebug(Debug bool) *Config {
	config.Debug = Debug
	return config
}

// SetDisableTraceInterceptor ...
func (config *Config) SetDisableTraceInterceptor(DisableTraceInterceptor bool) *Config {
	config.DisableTraceInterceptor = DisableTraceInterceptor
	return config
}

// SetDisableAidInterceptor ...
func (config *Config) SetDisableAidInterceptor(DisableAidInterceptor bool) *Config {
	config.DisableAidInterceptor = DisableAidInterceptor
	return config
}

// SetDisableTimeoutInterceptor ...
func (config *Config) SetDisableTimeoutInterceptor(DisableTimeoutInterceptor bool) *Config {
	config.DisableTimeoutInterceptor = DisableTimeoutInterceptor
	return config
}

// SetDisableMetricInterceptor ...
func (config *Config) SetDisableMetricInterceptor(DisableMetricInterceptor bool) *Config {
	config.DisableMetricInterceptor = DisableMetricInterceptor
	return config
}

// SetDisableAccessInterceptor ...
func (config *Config) SetDisableAccessInterceptor(DisableAccessInterceptor bool) *Config {
	config.DisableAccessInterceptor = DisableAccessInterceptor
	return config
}

// SetAccessInterceptorLevel ...
func (config *Config) SetAccessInterceptorLevel(AccessInterceptorLevel string) *Config {
	config.AccessInterceptorLevel = AccessInterceptorLevel
	return config
}

// WithDialOption ...
func (config *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	if config.dialOptions == nil {
		config.dialOptions = make([]grpc.DialOption, 0)
	}
	config.dialOptions = append(config.dialOptions, opts...)
	return config
}

// Build ...
func (config *Config) Build() (*grpc.ClientConn, error) {
	if config.Debug {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(config.Address)),
		)
	}

	if !config.DisableAidInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(aidUnaryClientInterceptor()),
		)
	}

	if !config.DisableTimeoutInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(config.ReadTimeout, config.SlowThreshold)),
		)
	}

	if !config.DisableTraceInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	if !config.DisableAccessInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(config.Name, config.AccessInterceptorLevel)),
		)
	}

	if !config.DisableMetricInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(metricUnaryClientInterceptor(config.Name)),
		)
	}
	return newGRPCClient(config)
}

// Singleton returns a singleton client conn.
func (config *Config) Singleton() (*grpc.ClientConn, error) {
	if val, ok := singleton.Load(constant.ModuleClientGrpc, config.Name); ok && val != nil {
		return val.(*grpc.ClientConn), nil
	}
	cc, err := config.Build()
	if err != nil {
		return nil, err
	}

	singleton.Store(constant.ModuleClientGrpc, config.Name, cc)

	return cc, nil
}

// MustSingleton panics when error found.
func (config *Config) MustSingleton() *grpc.ClientConn {
	cc, err := config.Singleton()
	if err != nil {
		logger.Logger.Panic("client grpc build client conn panic")
	}
	return cc
}

// MustBuild ...
func (config *Config) MustBuild() *grpc.ClientConn {
	reg, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("client grpc build client failed: %v", err)
	}
	return reg
}
