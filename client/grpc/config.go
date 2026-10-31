package grpc

import (
	"time"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/singleton"
	"github.com/abulo/ratel/v3/registry/etcdv3"
	"github.com/abulo/ratel/v3/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config 保存 gRPC 客户端的全部配置项，由 New 初始化默认值，经链式 Set 方法填充后由 Build 构建连接。
//
// Config 为配置载体，应在 Build 前一次性写定，构建期间不保证并发安全。Build 要求 Etcd 已通过
// SetEtcd 配置，构建时会据此创建 etcd 解析器，缺失会导致空指针。
type Config struct {
	// Name 为客户端实例名称，用作单例缓存分组键以及访问日志、指标的服务名标签。
	Name string
	// BalancerName 为负载均衡策略名称，New 默认使用 roundrobin。
	BalancerName string
	// Address 为拨号目标地址或服务名，由配置的解析器解释。
	Address string
	// Block 为预留字段，当前拨号流程不读取。
	Block bool
	// DialTimeout 为拨号超时时间，New 默认 3 秒。
	DialTimeout time.Duration
	// ReadTimeout 为单次调用的超时时间，供超时拦截器使用，New 默认 1 秒。
	ReadTimeout time.Duration
	// Direct 为预留字段，当前寻址流程不读取。
	Direct bool
	// KeepAlive 为 gRPC keepalive 客户端参数；为 nil 时不单独设置。
	KeepAlive   *keepalive.ClientParameters
	dialOptions []grpc.DialOption
	// SlowThreshold 为慢调用判定阈值，调用耗时超过该值时记录错误日志，New 默认 600 毫秒。
	SlowThreshold time.Duration
	// EnableDebug 指定是否装配调试拦截器，在控制台打印目标地址、方法、响应与错误。
	EnableDebug bool
	// EnableTraceInterceptor 指定是否装配链路追踪拦截器。
	EnableTraceInterceptor bool
	// EnableAidInterceptor 指定是否装配应用 ID 元数据透传拦截器。
	EnableAidInterceptor bool
	// EnableTimeoutInterceptor 指定是否装配调用超时与慢调用统计拦截器。
	EnableTimeoutInterceptor bool
	// EnableMetricInterceptor 指定是否装配指标采集拦截器。
	EnableMetricInterceptor bool
	// EnableAccessInterceptor 指定是否装配访问日志拦截器。
	EnableAccessInterceptor bool
	// AccessInterceptorLevel 为访问日志级别，New 默认 "info"。
	AccessInterceptorLevel string
	// Etcd 为 etcd 服务发现配置，Build 前必须通过 SetEtcd 设置。
	Etcd *etcdv3.Config
}

// New 创建带默认值的 Config：负载均衡为 roundrobin，拨号超时 3 秒、
// 读取超时 1 秒、慢调用阈值 600 毫秒，访问日志级别 info，Block 为 true。
// 传输凭证统一由 getDialOptions 装配。
func New() *Config {
	return &Config{
		// PERF-5-5 [P2]：移除 New() 中重复预置的 insecure 凭证，凭证装配单点收敛到 getDialOptions，
		// 避免同一连接上出现两条 WithTransportCredentials 选项；拨号连接语义保持不变。
		dialOptions:            []grpc.DialOption{},
		BalancerName:           roundrobin.Name,
		DialTimeout:            time.Second * 3,
		ReadTimeout:            util.Duration("1s"),
		SlowThreshold:          util.Duration("600ms"),
		AccessInterceptorLevel: "info",
		Block:                  true,
	}
}

// SetEtcd 设置 etcd 服务发现配置，返回配置自身以便链式调用；Build 前必须调用本方法。
func (config *Config) SetEtcd(option *etcdv3.Config) *Config {
	config.Etcd = option
	return config
}

// SetName 设置客户端实例名称，返回配置自身以便链式调用。
func (config *Config) SetName(Name string) *Config {
	config.Name = Name
	return config
}

// SetBalancerName 设置负载均衡策略名称，返回配置自身以便链式调用。
func (config *Config) SetBalancerName(BalancerName string) *Config {
	config.BalancerName = BalancerName
	return config
}

// SetAddress 设置拨号目标地址或服务名，返回配置自身以便链式调用。
func (config *Config) SetAddress(Address string) *Config {
	config.Address = Address
	return config
}

// SetBlock 写入预留的 Block 字段，返回配置自身以便链式调用；该字段当前不参与拨号。
func (config *Config) SetBlock(Block bool) *Config {
	config.Block = Block
	return config
}

// SetDialTimeout 设置拨号超时时间，返回配置自身以便链式调用。
func (config *Config) SetDialTimeout(DialTimeout time.Duration) *Config {
	config.DialTimeout = DialTimeout
	return config
}

// SetReadTimeout 设置单次调用超时时间，返回配置自身以便链式调用。
func (config *Config) SetReadTimeout(ReadTimeout time.Duration) *Config {
	config.ReadTimeout = ReadTimeout
	return config
}

// SetDirect 写入预留的 Direct 字段，返回配置自身以便链式调用；该字段当前不参与寻址。
func (config *Config) SetDirect(Direct bool) *Config {
	config.Direct = Direct
	return config
}

// SetSlowThreshold 设置慢调用判定阈值，返回配置自身以便链式调用。
func (config *Config) SetSlowThreshold(SlowThreshold time.Duration) *Config {
	config.SlowThreshold = SlowThreshold
	return config
}

// SetEnableDebug 设置是否装配调试拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableDebug(Debug bool) *Config {
	config.EnableDebug = Debug
	return config
}

// SetEnableTraceInterceptor 设置是否装配链路追踪拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableTraceInterceptor(EnableTraceInterceptor bool) *Config {
	config.EnableTraceInterceptor = EnableTraceInterceptor
	return config
}

// SetEnableAidInterceptor 设置是否装配应用 ID 透传拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableAidInterceptor(EnableAidInterceptor bool) *Config {
	config.EnableAidInterceptor = EnableAidInterceptor
	return config
}

// SetEnableTimeoutInterceptor 设置是否装配调用超时与慢调用拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableTimeoutInterceptor(EnableTimeoutInterceptor bool) *Config {
	config.EnableTimeoutInterceptor = EnableTimeoutInterceptor
	return config
}

// SetEnableMetricInterceptor 设置是否装配指标采集拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableMetricInterceptor(EnableMetricInterceptor bool) *Config {
	config.EnableMetricInterceptor = EnableMetricInterceptor
	return config
}

// SetEnableAccessInterceptor 设置是否装配访问日志拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableAccessInterceptor(EnableAccessInterceptor bool) *Config {
	config.EnableAccessInterceptor = EnableAccessInterceptor
	return config
}

// SetAccessInterceptorLevel 设置访问日志级别，返回配置自身以便链式调用。
func (config *Config) SetAccessInterceptorLevel(AccessInterceptorLevel string) *Config {
	config.AccessInterceptorLevel = AccessInterceptorLevel
	return config
}

// WithDialOption 追加自定义 gRPC 拨号选项，返回配置自身以便链式调用。
func (config *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	if config.dialOptions == nil {
		config.dialOptions = make([]grpc.DialOption, 0)
	}
	config.dialOptions = append(config.dialOptions, opts...)
	return config
}

// Build 按已开启的开关装配对应单向拦截器，并依据配置创建 gRPC ClientConn；
// 拨号失败时返回错误，Etcd 缺失会导致空指针。
func (config *Config) Build() (*grpc.ClientConn, error) {
	if config.EnableDebug {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(config.Address)),
		)
	}

	if config.EnableAidInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(aidUnaryClientInterceptor()),
		)
	}

	if config.EnableTimeoutInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(config.ReadTimeout, config.SlowThreshold)),
		)
	}

	if config.EnableTraceInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	if config.EnableAccessInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(config.Name, config.AccessInterceptorLevel)),
		)
	}

	if config.EnableMetricInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(metricUnaryClientInterceptor(config.Name)),
		)
	}
	return newGRPCClient(config)
}

// Singleton 按 Name 复用单例 ClientConn：已缓存同名称连接时直接返回，否则构建后缓存；
// 构建失败时返回错误。Name 相同的配置须保持一致，避免覆盖已缓存连接。
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

// MustSingleton 单例方式构建 ClientConn，构建失败时直接 panic。
func (config *Config) MustSingleton() *grpc.ClientConn {
	cc, err := config.Singleton()
	if err != nil {
		logger.Logger.Panic("client grpc build client conn panic")
	}
	return cc
}

// MustBuild 构建 ClientConn，构建失败时直接 panic，适用于启动期配置确定、不允许失败的场景。
func (config *Config) MustBuild() *grpc.ClientConn {
	reg, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("client grpc build client failed: %v", err)
	}
	return reg
}
