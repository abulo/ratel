package grpc

import (
	"time"

	"github.com/abulo/ratel/v3/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
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
	OnDialError               string // panic | error
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
}

func New() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 3,
		ReadTimeout:            util.Duration("1s"),
		SlowThreshold:          util.Duration("600ms"),
		OnDialError:            "panic",
		AccessInterceptorLevel: "info",
		Block:                  true,
	}
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
func (config *Config) Build() *grpc.ClientConn {
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
