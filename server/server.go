package server

import (
	"context"
	"fmt"

	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/env"
)

// Option 是 ServiceInfo 的配置选项函数，以函数式选项模式在 ApplyOptions 中写入服务信息。
type Option func(c *ServiceInfo)

// ConfigInfo 保存服务对外暴露的路由清单，供服务注册与治理使用。
type ConfigInfo struct {
	Routes []Route
}

// ServiceInfo 描述一个可注册服务的身份、地址、流量归属与环境元信息。
//
// 零值不可直接使用，须经 ApplyOptions 结合各 With 选项构造；该结构体为注册信息载体，
// 构造完成后不应在并发环境下继续修改。
type ServiceInfo struct {
	// Name 为服务名称。
	Name       string               `json:"name"`
	// AppID 为应用唯一标识。
	AppID      string               `json:"appId"`
	// Scheme 为访问协议，如 http、grpc、asynq。
	Scheme     string               `json:"scheme"`
	// Address 为服务对外暴露的监听地址 host:port。
	Address    string               `json:"address"`
	// Weight 为负载均衡权重，默认 100。
	Weight     float64              `json:"weight"`
	// Enable 标记服务是否允许被调度调用。
	Enable     bool                 `json:"enable"`
	// Healthy 标记服务当前是否健康。
	Healthy    bool                 `json:"healthy"`
	// Metadata 为随服务一起注册的键值元数据。
	Metadata   map[string]string    `json:"metadata"`
	// Region 为服务部署所在区域。
	Region     string               `json:"region"`
	// Zone 为服务部署所在可用区。
	Zone       string               `json:"zone"`
	// Kind 为服务类型（提供者、监控、消费者等）。
	Kind       constant.ServiceKind `json:"kind"`
	// Deployment 为部署组，不同组之间流量相互隔离，例如内部调用与第三方调用可配置不同部署组。
	Deployment string               `json:"deployment"`
	// Group 为流量组，负载均衡时按 Group 在组之间分配流量。
	Group      string               `json:"group"`
	// Services 记录该服务对外提供的子服务清单。
	Services   map[string]*Service  `json:"services" toml:"services"`
	// Version 为服务版本号，默认 v1。
	Version    string               `json:"version"`
	// Mode 为应用运行模式（如 dev、release）。
	Mode       string               `json:"mode"`
	// Hostname 为运行所在主机名。
	Hostname   string               `json:"hostname"`
}

// Service 描述注册到注册中心的单个子服务，包含其命名空间、名称、标签与方法清单。
type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

// RegistryName 返回服务在注册中心使用的唯一标识，格式为 协议:服务名:v1:运行模式/地址。
func (s *ServiceInfo) RegistryName() string {
	return fmt.Sprintf("%s:%s:%s:%s/%s", s.Scheme, s.Name, "v1", env.AppMode(), s.Address)
}

// ServicePrefix 返回服务路由前缀，格式为 协议:服务名:v1:运行模式/，用于按前缀匹配服务。
func (s *ServiceInfo) ServicePrefix() string {
	return fmt.Sprintf("%s:%s:%s:%s/", s.Scheme, s.Name, "v1", env.AppMode())
}

// Label 返回 协议://地址 形式的服务定位标签。
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

// Equal 判断两个 ServiceInfo 是否在注册标识上等价。
//
// 仅比较 Name、Address、Kind、Group 四个区分身份的字段；实现该方法是为了兼容
// grpc-go resolver 在服务发现与负载均衡中对地址属性做等值比较的契约。
func (si ServiceInfo) Equal(o any) bool {
	oa, ok := o.(ServiceInfo)
	return ok &&
		oa.Name == si.Name &&
		oa.Address == si.Address &&
		oa.Kind == si.Kind &&
		oa.Group == si.Group
}

// Server 抽象各类服务端（HTTP、gRPC、异步任务等）统一的生命周期与查询接口。
//
// 实现须保证 Serve 在后台阻塞运行，Stop 立即关闭，GracefulStop 在指定上下文超时前优雅退出；
// Info 返回当前服务的注册信息，Health 用于健康探活。
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
	Health() bool
}

// Route 描述一条对外暴露的路由及其权重组配置。
type Route struct {
	// WeightGroups 为该路由参与负载均衡的权重组列表。
	WeightGroups []WeightGroup
	// Method 为路由对应的请求方法名。
	Method string
}

// WeightGroup 描述某个流量组及其在路由上分配的权重。
type WeightGroup struct {
	Group  string
	Weight int
}

// ApplyOptions 以默认服务信息为基础，依次套用传入的 Option 并返回构造完成的 ServiceInfo。
//
// 默认信息中的 Metadata 已初始化，WithMetaData 可直接写入；调用方未传入任何选项时返回默认值。
func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

// WithMetaData 向服务元数据写入一个键值对；Metadata 须已由 ApplyOptions 初始化。
func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

// WithScheme 设置服务访问协议，如 http、grpc、asynq。
func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

// WithAddress 设置服务对外暴露的监听地址 host:port。
func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

// WithKind 设置服务类型（提供者、监控、消费者等）。
func WithKind(kind constant.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       env.Name(),
		AppID:      env.AppID(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     env.AppRegion(),
		Zone:       env.AppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
		Mode:       env.AppMode(),
		Hostname:   env.AppHost(),
		Version:    "v1",
	}
	si.Metadata["appMode"] = env.AppMode()
	si.Metadata["appHost"] = env.AppHost()
	si.Metadata["startTime"] = env.StartTime()
	si.Metadata["buildTime"] = env.BuildTime()
	si.Metadata["buildVersion"] = env.BuildVersion()
	si.Metadata["ratelVersion"] = env.RatelVersion()
	return si
}
