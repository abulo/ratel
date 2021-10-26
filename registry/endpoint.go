package registry

import (
	"github.com/abulo/ratel/v2/server"
)

// Endpoints ...
type Endpoints struct {
	// 服务节点列表
	Nodes map[string]server.ServiceInfo
}
