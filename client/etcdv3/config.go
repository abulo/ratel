package etcdv3

import (
	"time"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/singleton"
	"github.com/abulo/ratel/v3/util"
)

// Config 是 etcd 客户端的全部配置项，由 New 创建并通过链式 Set 方法填充，最终由 Build 构建客户端。
//
// 零值字段不完整，须配合 Set 方法设置 Endpoints 等必填项；Config 为配置载体，应在构建前一次性写定，
// 构建期间不保证并发安全。
type Config struct {
	// Name 为客户端实例名称，作为单例缓存的分组键；为空时单例实例不做区分。
	Name string `json:"name"`
	// Endpoints 为 etcd 集群节点地址列表，格式 host:port，构建客户端时必填。
	Endpoints []string `json:"endpoints"`
	// CertFile 为客户端 TLS 证书文件路径，须与 KeyFile 成对配置。
	CertFile string `json:"certFile"`
	// KeyFile 为客户端 TLS 私钥文件路径，须与 CertFile 成对配置。
	KeyFile string `json:"keyFile"`
	// CaCert 为信任的 CA 证书文件路径，用于校验服务端证书。
	CaCert string `json:"caCert"`
	// BasicAuth 指定是否启用用户名密码基本认证。
	BasicAuth bool `json:"basicAuth"`
	// UserName 为基本认证用户名，仅在 BasicAuth 为 true 时生效。
	UserName string `json:"userName"`
	// Password 为基本认证密码，序列化时通过 "-" 忽略，避免泄露。
	Password string `json:"-"`
	// ConnectTimeout 为建立连接的超时时间，New 默认 5 秒。
	ConnectTimeout time.Duration `json:"connectTimeout"`
	// Secure 为 false 时追加 insecure 明文传输凭证；为 true 时不追加，TLS 由 CA 或证书配置决定。
	Secure bool `json:"secure"`
	// AutoSyncInterval 为自动同步集群成员列表的间隔时长，为 0 表示不自动同步。
	AutoSyncInterval time.Duration `json:"autoAsyncInterval"`
	// EnableTrace 指定是否为客户端 gRPC 调用装配链路追踪拦截器。
	EnableTrace bool `json:"enableTrace" toml:"enableTrace"`
	// TTL 为预留的租约存活时间配置（单位秒），当前构建流程未读取该字段。
	TTL int
}

// New 创建带默认值的 Config：不启用 BasicAuth 与 Secure，连接超时为 5 秒。
func New() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: util.Duration("5s"),
		Secure:         false,
	}
}

// Build 依据当前配置构建 etcd 客户端，配置非法或连接初始化失败时返回错误。
func (config *Config) Build() (*Client, error) {
	return NewClient(config)
}

// Singleton 按 Name 复用单例客户端：已缓存同名称实例时直接返回，否则构建后缓存；
// 构建失败时记录日志并返回错误。Name 相同的配置须保持一致，避免覆盖已缓存实例。
func (config *Config) Singleton() (*Client, error) {
	if client, ok := singleton.Load(constant.ModuleRegistryEtcd, config.Name); ok && client != nil {
		return client.(*Client), nil
	}

	client, err := config.Build()
	if err != nil {
		logger.Logger.Error("build etcd client failed", err)
		return nil, err
	}

	singleton.Store(constant.ModuleRegistryEtcd, config.Name, client)

	return client, nil
}

// MustBuild 构建客户端，失败时直接 panic，适用于启动期配置确定、不允许失败的场景。
func (config *Config) MustBuild() *Client {
	client, err := config.Build()
	if err != nil {
		logger.Logger.Panic("build etcd client failed", err)
	}
	return client
}

// MustSingleton 按单例方式构建客户端，失败时直接 panic。
func (config *Config) MustSingleton() *Client {
	client, err := config.Singleton()
	if err != nil {
		logger.Logger.Panic("build etcd client failed", err)
	}
	return client
}

// SetEndpoints 追加去重后的集群节点地址，并返回配置自身以便链式调用。
func (config *Config) SetEndpoints(endpoint []string) *Config {
	// PERF-5-19 [P2]：改用 seen 哈希集合做 O(N) 去重，替代反射式 util.InArray 的 O(N²) 遍历；
	// 保留原有追加顺序与去重语义不变。
	seen := make(map[string]struct{}, len(config.Endpoints)+len(endpoint))
	for _, existing := range config.Endpoints {
		seen[existing] = struct{}{}
	}
	for _, item := range endpoint {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			config.Endpoints = append(config.Endpoints, item)
		}
	}
	return config
}

// SetCertFile 设置客户端 TLS 证书文件路径，返回配置自身以便链式调用。
func (config *Config) SetCertFile(cert string) *Config {
	config.CertFile = cert
	return config
}

// SetKeyFile 设置客户端 TLS 私钥文件路径，返回配置自身以便链式调用。
func (config *Config) SetKeyFile(key string) *Config {
	config.KeyFile = key
	return config
}

// SetCaCert 设置信任的 CA 证书文件路径，返回配置自身以便链式调用。
func (config *Config) SetCaCert(ca string) *Config {
	config.CaCert = ca
	return config
}

// SetBasicAuth 设置是否启用基本认证，返回配置自身以便链式调用。
func (config *Config) SetBasicAuth(auth bool) *Config {
	config.BasicAuth = auth
	return config
}

// SetUserName 设置基本认证用户名，返回配置自身以便链式调用。
func (config *Config) SetUserName(userName string) *Config {
	config.UserName = userName
	return config
}

// SetPassword 设置基本认证密码，返回配置自身以便链式调用。
func (config *Config) SetPassword(pwd string) *Config {
	config.Password = pwd
	return config
}

// SetConnectTimeout 设置建立连接的超时时间，返回配置自身以便链式调用。
func (config *Config) SetConnectTimeout(timeout time.Duration) *Config {
	config.ConnectTimeout = timeout
	return config
}

// SetSecure 设置是否不追加 insecure 明文凭证，返回配置自身以便链式调用。
func (config *Config) SetSecure(secure bool) *Config {
	config.Secure = secure
	return config
}

// SetEnableTrace 设置是否装配链路追踪拦截器，返回配置自身以便链式调用。
func (config *Config) SetEnableTrace(enableTrace bool) *Config {
	config.EnableTrace = enableTrace
	return config
}

// SetAutoSyncInterval 设置自动同步集群成员列表的间隔时长，返回配置自身以便链式调用。
func (config *Config) SetAutoSyncInterval(autoSyncInterval time.Duration) *Config {
	config.AutoSyncInterval = autoSyncInterval
	return config
}

// SetName 设置单例缓存分组名称，返回配置自身以便链式调用。
func (config *Config) SetName(Name string) *Config {
	config.Name = Name
	return config
}

// GetName 返回当前配置的实例名称。
func (config *Config) GetName() string {
	return config.Name
}
