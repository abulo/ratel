package registry

import (
	"encoding/json"
	"maps"

	"github.com/abulo/ratel/v3/server"
)

// Endpoints ...
type Endpoints struct {
	// 服务节点列表
	Nodes map[string]server.ServiceInfo

	// 路由配置
	RouteConfigs map[string]RouteConfig

	// 消费者元数据
	ConsumerConfigs map[string]ConsumerConfig

	// 服务元信息
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

// DeepCopy ...
func (in *Endpoints) DeepCopy() *Endpoints {
	if in == nil {
		return nil
	}

	out := newEndpoints()
	in.DeepCopyInfo(out)
	return out
}

// DeepCopyInfo ...
func (in *Endpoints) DeepCopyInfo(out *Endpoints) {
	maps.Copy(out.Nodes, in.Nodes)
	maps.Copy(out.RouteConfigs, in.RouteConfigs)
	maps.Copy(out.ConsumerConfigs, in.ConsumerConfigs)
	maps.Copy(out.ProviderConfigs, in.ProviderConfigs)
}

// ProviderConfig config of provider
// 通过这个配置，修改provider的属性
type ProviderConfig struct {
	ID         string            `json:"id"`
	Scheme     string            `json:"scheme"`
	Host       string            `json:"host"`
	Region     string            `json:"region"`
	Zone       string            `json:"zone"`
	Deployment string            `json:"deployment"`
	Metadata   map[string]string `json:"metadata"`
	Enable     bool              `json:"enable"`
}

// ConsumerConfig config of consumer
// 客户端调用app的配置
type ConsumerConfig struct {
	ID     string `json:"id"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

// RouteConfig ...
type RouteConfig struct {
	ID         string   `json:"id" toml:"id"`
	Scheme     string   `json:"scheme" toml:"scheme"`
	Host       string   `json:"host" toml:"host"`
	Deployment string   `json:"deployment"`
	URI        string   `json:"uri"`
	Upstream   Upstream `json:"upstream"`
}

// String ...
func (config RouteConfig) String() string {
	bs, _ := json.Marshal(config)
	return string(bs)
}

// Upstream represents upstream balancing config
type Upstream struct {
	Nodes  map[string]int `json:"nodes"`
	Groups map[string]int `json:"groups"`
}
