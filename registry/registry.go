package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/abulo/ratel/v3/server"
)

// Event 表示服务发现推送的事件类型。
type Event uint8

const (
	// EventUnknown 为未识别或未知的事件类型。
	EventUnknown Event = iota
	// EventUpdate 表示节点信息发生变更。
	EventUpdate
	// EventDelete 表示节点被移除。
	EventDelete
)

// Kind 区分注册中心中服务节点的角色类型。
type Kind uint8

const (
	// KindUnknown 为未识别或未知的节点角色。
	KindUnknown Kind = iota
	// KindProvider 为服务提供者节点。
	KindProvider
	// KindConfigurator 为配置中心节点。
	KindConfigurator
	// KindConsumer 为服务消费者节点。
	KindConsumer
)

// String 返回 Kind 对应的注册中心分组名称字符串；未知类型返回 "unknown"。
func (kind Kind) String() string {
	switch kind {
	case KindProvider:
		return "providers"
	case KindConfigurator:
		return "configurators"
	case KindConsumer:
		return "consumers"
	default:
		return "unknown"
	}
}

// ToKind 将分组名称字符串还原为 Kind；无法识别的字符串返回 KindUnknown。
func ToKind(kindStr string) Kind {
	switch kindStr {
	case "providers":
		return KindProvider
	case "configurators":
		return KindConfigurator
	case "consumers":
		return KindConsumer
	default:
		return KindUnknown
	}
}

// ServerInstance 描述一个可被发现的服务实例的网络位置与标签。
type ServerInstance struct {
	// Scheme 为访问协议，如 http、grpc。
	Scheme string
	// IP 为实例监听的 IP 地址。
	IP     string
	// Port 为实例监听端口。
	Port   int
	// Labels 为实例附带的标签键值对。
	Labels map[string]string
}

// EventMessage 封装一次服务发现事件的角色、身份信息与附带数据。
type EventMessage struct {
	Event
	Kind
	// Name 为事件关联的服务名称。
	Name    string
	// Scheme 为事件关联的访问协议。
	Scheme  string
	// Address 为事件关联的节点地址。
	Address string
	// Message 为事件附带的原始数据，由具体实现约定其结构。
	Message any
}

// Registry 是服务注册与发现的统一抽象。
//
// 实现须自行控制 RPC 超时与重试；WatchServices 返回的通道在服务下线或上下文
// 取消后应被实现负责关闭；Close 须可重复调用并释放底层连接与监听资源。
type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	UnregisterService(context.Context, *server.ServiceInfo) error
	ListServices(context.Context, string) ([]*server.ServiceInfo, error)
	WatchServices(context.Context, string) (chan Endpoints, error)
	Kind() string
	io.Closer
}

// GetServiceKey 根据前缀与服务信息拼接注册中心中的唯一键，
// 格式为 /前缀/服务名/角色/协议://地址。
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", prefix, s.Name, s.Kind.String(), s.Scheme, s.Address)
}

// GetServiceValue 将服务信息序列化为 JSON 字符串，作为注册中心中的存储值；
// 序列化失败时返回空字符串。
func GetServiceValue(s *server.ServiceInfo) string {
	val, _ := json.Marshal(s)
	return string(val)
}

// GetService 将注册中心存储的 JSON 字符串反序列化为 ServiceInfo；
// 解析失败时返回零值 ServiceInfo 的指针。
func GetService(s string) *server.ServiceInfo {
	var si server.ServiceInfo
	_ = json.Unmarshal([]byte(s), &si)
	return &si
}

// Configuration 保存服务的客户端路由策略与服务端标签分组配置。
type Configuration struct {
	// Routes 为客户端路由策略列表。
	Routes []Route           `json:"routes"`
	// Labels 为服务端标签，用于分组。
	Labels map[string]string `json:"labels"`
}

// Route 描述单条路由方法的流量分配规则。
type Route struct {
	// Method 为路由方法名。
	Method string `json:"method" toml:"method"`
	// WeightGroups 为权重组列表，按各组比率分配流量。
	WeightGroups []WeightGroup `json:"weightGroups" toml:"weightGroups"`
	// Deployment 为目标部署组，流量将被导入该部署组。
	Deployment string `json:"deployment" toml:"deployment"`
}

// WeightGroup 描述单个流量组及其分配权重。
type WeightGroup struct {
	// Group 为流量组名称。
	Group  string `json:"group" toml:"group"`
	// Weight 为该组在总流量中所占权重。
	Weight int    `json:"weight" toml:"weight"`
}
