package registry

import (
	"log"
	"sync"

	"github.com/abulo/ratel/v3/core/logger"
)

// PERF-7-35 [P2] registryBuilder 由裸 map 改为读写锁保护：注册走写锁、查找走读锁，避免后续运行期动态注册触发 map 并发读写 fatal；当前启动期顺序注册行为不变
var (
	registryBuilder   = make(map[string]Builder)
	registryBuilderMu sync.RWMutex
)

// Config 是注册器配置集合，按键名保存各注册器的 ConfigLab 配置。
type Config map[string]ConfigLab

// ConfigLab 描述单个注册器的初始化配置。
type ConfigLab struct {
	// Kind 为底层注册器类型，例如 etcdv3、consul；为空时按 etcdv3 处理。
	Kind string `json:"kind" description:"底层注册器类型, eg: etcdv3, consul"`
	// ConfigKey 为构造该注册器时使用的配置键。
	ConfigKey string `json:"configKey" description:"底册注册器的配置键"`
	// DeplaySeconds 为延迟注册的秒数。
	DeplaySeconds int `json:"deplaySeconds" description:"延迟注册"`
}

// DefaultRegisterer 是包级默认注册器，初始为 Local 本地实现，
// 可由 InitDefaultRegister 在遍历配置后覆盖；调用方不应在并发环境中替换它。
var DefaultRegisterer Registry = &Local{}

// Builder 根据配置键构造对应的 Registry 实例。
type Builder func(string) Registry

// BuildFunc 根据配置键构造 Registry 实例，构造失败时通过第二个返回值返回错误。
type BuildFunc func(string) (Registry, error)

// RegisterBuilder 注册指定类型名称对应的注册器构造函数。
//
// 同一类型重复注册会直接 panic，避免静默覆盖已有实现。
func RegisterBuilder(kind string, build Builder) {
	registryBuilderMu.Lock()
	defer registryBuilderMu.Unlock()
	if _, ok := registryBuilder[kind]; ok {
		log.Panicf("duplicate register registry builder: %s", kind)
	}
	registryBuilder[kind] = build
}

// New 创建一个空的注册器配置集合。
func New() Config {
	var config Config
	return config
}

// Lab 向配置集合追加指定名称的注册器配置，返回更新后的配置集合以支持链式调用。
func (config Config) Lab(name string, lab ConfigLab) Config {
	// FIXME(owner): 此处应使用入参 name 作为键，当前硬编码为 "ddd"，不同配置会互相覆盖。
	config["ddd"] = lab
	return config
}

// InitDefaultRegister 遍历配置集合，按各配置的类型查找已注册的构造函数并构造默认注册器。
//
// 配置类型为空时按 etcdv3 处理；类型未注册构造函数时仅记录日志并跳过，不影响其余配置；
// 循环内每次构造成功都会覆盖 DefaultRegisterer，最终默认注册器为最后一个成功构造的实例。
func (config Config) InitDefaultRegister() {
	for name, item := range config {
		var itemKind = item.Kind
		if itemKind == "" {
			itemKind = "etcdv3"
		}
		registryBuilderMu.RLock()
		build, ok := registryBuilder[itemKind]
		registryBuilderMu.RUnlock()
		if !ok {
			logger.Logger.Printf("invalid registry kind: %s", itemKind)
			continue
		}
		DefaultRegisterer = build(item.ConfigKey)
		logger.Logger.Printf("build registry %s with config: %s", name, item.ConfigKey)
	}
}
