package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/abulo/ratel/server"
)

// Event ...
type Event uint8

// EventUnknown ...
const (
	// EventUnknown ...
	EventUnknown Event = iota
	// EventUpdate ...
	EventUpdate
	// EventDelete ...
	EventDelete
)

// Kind ...
type Kind uint8

// KindUnknown ...
const (
	// KindUnknown ...
	KindUnknown Kind = iota
	// KindProvider ...
	KindProvider
	// KindConfigurator ...
	KindConfigurator
	// KindConsumer ...
	KindConsumer
)

// String ...
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

// ToKind ...
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

// ServerInstance ...
type ServerInstance struct {
	Scheme string
	IP     string
	Port   int
	Labels map[string]string
}

// EventMessage ...
type EventMessage struct {
	Event
	Kind
	Name    string
	Scheme  string
	Address string
	Message interface{}
}

// Registry register/unregister service
// registry impl should control rpc timeout
type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	UnregisterService(context.Context, *server.ServiceInfo) error
	ListServices(context.Context, string) ([]*server.ServiceInfo, error)
	WatchServices(context.Context, string) (chan Endpoints, error)
	Kind() string
	io.Closer
}

// GetServiceKey ..
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", prefix, s.Name, s.Kind.String(), s.Scheme, s.Address)
}

// GetServiceValue ..
func GetServiceValue(s *server.ServiceInfo) string {
	val, _ := json.Marshal(s)
	return string(val)
}

// GetService ..
func GetService(s string) *server.ServiceInfo {
	var si server.ServiceInfo
	_ = json.Unmarshal([]byte(s), &si)
	return &si
}

// Configuration ...
type Configuration struct {
	Routes []Route           `json:"routes"` // 配置客户端路由策略
	Labels map[string]string `json:"labels"` // 配置服务端标签: 分组
}

// Route represents route configuration
type Route struct {
	// 路由方法名
	Method string `json:"method" toml:"method"`
	// 路由权重组, 按比率在各个权重组中分配流量
	WeightGroups []WeightGroup `json:"weightGroups" toml:"weightGroups"`
	// 路由部署组, 将流量导入部署组
	Deployment string `json:"deployment" toml:"deployment"`
}

// WeightGroup ...
type WeightGroup struct {
	Group  string `json:"group" toml:"group"`
	Weight int    `json:"weight" toml:"weight"`
}
