package registry

import (
	"encoding/json"
	"maps"

	"github.com/abulo/ratel/v3/server"
)

// Endpoints 聚合一次服务发现结果中的节点清单与各类治理配置。
//
// 四个映射须由 newEndpoints 初始化后使用，零值直接写入会触发 nil map panic；
// 该结构体不保证并发安全，由接收方在监听回调中独占使用。
type Endpoints struct {
	// Nodes 为节点地址到服务信息的映射。
	Nodes map[string]server.ServiceInfo

	// RouteConfigs 为路由配置映射。
	RouteConfigs map[string]RouteConfig

	// ConsumerConfigs 为消费者元数据映射。
	ConsumerConfigs map[string]ConsumerConfig

	// ProviderConfigs 为提供者元数据映射。
	ProviderConfigs map[string]ProviderConfig
}

func newEndpoints() *Endpoints {
	return &Endpoints{
		Nodes:           make(map[string]server.ServiceInfo),
		RouteConfigs:    make(map[string]RouteConfig),
		ConsumerConfigs: make(map[string]ConsumerConfig),
		ProviderConfigs: make(map[string]ProviderConfig),
	}
}

// DeepCopy 返回 Endpoints 的独立副本；接收者为 nil 时返回 nil。
//
// 仅复制各映射的顶层条目，映射内部的嵌套元数据不会被递归克隆。
func (in *Endpoints) DeepCopy() *Endpoints {
	if in == nil {
		return nil
	}

	out := newEndpoints()
	in.DeepCopyInfo(out)
	return out
}

// DeepCopyInfo 将接收者各映射的条目复制到 out 中；同键条目以接收者的值覆盖。
func (in *Endpoints) DeepCopyInfo(out *Endpoints) {
	maps.Copy(out.Nodes, in.Nodes)
	maps.Copy(out.RouteConfigs, in.RouteConfigs)
	maps.Copy(out.ConsumerConfigs, in.ConsumerConfigs)
	maps.Copy(out.ProviderConfigs, in.ProviderConfigs)
}

// ProviderConfig 描述提供者节点的属性配置，可用于在线调整提供者的流量与归属。
type ProviderConfig struct {
	// ID 为提供者唯一标识。
	ID         string            `json:"id"`
	// Scheme 为访问协议。
	Scheme     string            `json:"scheme"`
	// Host 为提供者主机地址。
	Host       string            `json:"host"`
	// Region 为部署区域。
	Region     string            `json:"region"`
	// Zone 为部署可用区。
	Zone       string            `json:"zone"`
	// Deployment 为所属部署组。
	Deployment string            `json:"deployment"`
	// Metadata 为附加键值元数据。
	Metadata   map[string]string `json:"metadata"`
	// Enable 标记该提供者是否启用。
	Enable     bool              `json:"enable"`
}

// ConsumerConfig 描述消费者（客户端）调用应用时的配置。
type ConsumerConfig struct {
	// ID 为消费者唯一标识。
	ID     string `json:"id"`
	// Scheme 为访问协议。
	Scheme string `json:"scheme"`
	// Host 为消费者主机地址。
	Host   string `json:"host"`
}

// RouteConfig 描述一条路由规则及其上游负载均衡配置。
type RouteConfig struct {
	// ID 为路由唯一标识。
	ID         string   `json:"id" toml:"id"`
	// Scheme 为访问协议。
	Scheme     string   `json:"scheme" toml:"scheme"`
	// Host 为主机地址。
	Host       string   `json:"host" toml:"host"`
	// Deployment 为目标部署组。
	Deployment string   `json:"deployment"`
	// URI 为路由匹配的请求路径。
	URI        string   `json:"uri"`
	// Upstream 为上游负载均衡配置。
	Upstream   Upstream `json:"upstream"`
}

// String 将 RouteConfig 序列化为 JSON 字符串；序列化失败时返回空字符串。
func (config RouteConfig) String() string {
	bs, _ := json.Marshal(config)
	return string(bs)
}

// Upstream 描述上游节点与分组的权重负载均衡配置。
type Upstream struct {
	// Nodes 为节点名称到权重的映射。
	Nodes  map[string]int `json:"nodes"`
	// Groups 为分组名称到权重的映射。
	Groups map[string]int `json:"groups"`
}
