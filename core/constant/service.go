package constant

// ServiceKind 标识服务在框架中的角色类型。
type ServiceKind uint8

// 服务角色枚举。
const (
	// ServiceUnknown 为未识别或未指定的服务角色。
	ServiceUnknown ServiceKind = iota
	// ServiceProvider 为对外提供服务的服务端角色。
	ServiceProvider
	// ServiceMonitor 为监控组件角色。
	ServiceMonitor
	// ServiceConsumer 为消费服务的客户端角色。
	ServiceConsumer
)

var serviceKinds = make(map[ServiceKind]string)

func init() {
	serviceKinds[ServiceUnknown] = "unknown"
	serviceKinds[ServiceProvider] = "providers"
	serviceKinds[ServiceMonitor] = "monitor"
	serviceKinds[ServiceConsumer] = "consumers"
}

// String 返回服务角色对应的小写展示名称；取值未登记时回退为 "unknown"。
func (sk ServiceKind) String() string {
	if s, ok := serviceKinds[sk]; ok {
		return s
	}
	return "unknown"
}
