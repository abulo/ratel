package constant

// ServiceKind service kind
type ServiceKind uint8

// ServiceUnknown ...
const (
	//ServiceUnknown service non-name
	ServiceUnknown ServiceKind = iota
	//ServiceProvider service provider
	ServiceProvider
	//ServiceMonitor service monitor
	ServiceMonitor
	//ServiceConsumer service consumer
	ServiceConsumer
)

var serviceKinds = make(map[ServiceKind]string)

func init() {
	serviceKinds[ServiceUnknown] = "unknown"
	serviceKinds[ServiceProvider] = "providers"
	serviceKinds[ServiceMonitor] = "monitor"
	serviceKinds[ServiceConsumer] = "consumers"
}

// String ...
func (sk ServiceKind) String() string {
	if s, ok := serviceKinds[sk]; ok {
		return s
	}
	return "unknown"
}
