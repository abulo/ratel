package etcdv3

import (
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/singleton"
	"github.com/abulo/ratel/v3/registry"
)

// Config 保存 etcdv3 注册器的配置，内嵌底层 etcd 客户端配置。
type Config struct {
	*etcdv3.Config
	// ReadTimeout 为查询与注销等单次操作的超时时间，默认 3 秒。
	ReadTimeout time.Duration
	// Node 为节点标识，用作单例缓存的键。
	Node        string
	// Prefix 为注册中心键名前缀，默认 "ratel"。
	Prefix      string
	// ServiceTTL 为注册会话的租约时长；为 0 时不使用租约，写入的键永久保留。
	ServiceTTL  time.Duration
}

// New 创建带默认配置的 etcdv3 Config：读取超时 3 秒、键前缀 ratel、不启用租约。
func New() *Config {
	return &Config{
		Config:      etcdv3.New(),
		ReadTimeout: time.Second * 3,
		Prefix:      "ratel",
		ServiceTTL:  0,
	}
}

// SetNode 设置节点标识并返回配置本身，支持链式调用。
func (config *Config) SetNode(prefix string) *Config {
	config.Node = prefix
	return config
}

// GetNode 返回当前节点标识。
func (config *Config) GetNode() string {
	return config.Node
}

// SetPrefix 设置注册中心键名前缀并返回配置本身，支持链式调用。
func (config *Config) SetPrefix(prefix string) *Config {
	config.Prefix = prefix
	return config
}

// Build 依据配置构造 etcdv3 注册器；底层客户端连接失败或配置无效时返回错误。
func (config *Config) Build() (registry.Registry, error) {
	return newETCDRegistry(config)
}

// MustBuild 构造 etcdv3 注册器；构造失败时记录错误并直接 panic。
func (config *Config) MustBuild() registry.Registry {
	reg, err := config.Build()
	if err != nil {
		logger.Logger.Panicf("build registry failed: %v", err)
	}
	return reg
}

// Singleton 按节点标识返回缓存的注册器实例。
//
// 首次构造成功后将实例缓存到单例管理器，同一节点标识后续复用同一实例；
// 构造失败时返回错误且不写入缓存。
func (config *Config) Singleton() (registry.Registry, error) {
	if val, ok := singleton.Load(constant.ModuleClientEtcd, config.Node); ok {
		return val.(registry.Registry), nil
	}

	reg, err := config.Build()
	if err != nil {
		return nil, err
	}

	singleton.Store(constant.ModuleClientEtcd, config.Node, reg)

	return reg, nil
}

// MustSingleton 单例方式构造 etcdv3 注册器；构造失败时记录错误并直接 panic。
func (config *Config) MustSingleton() registry.Registry {
	reg, err := config.Singleton()
	if err != nil {
		logger.Logger.Panicf("build registry failed: %v", err)
	}

	return reg
}
