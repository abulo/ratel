package registry

import "github.com/abulo/ratel/v3/server"

// Operation 表示服务节点变更事件的操作类型。
type Operation uint8

const (
	// Add 表示新增一个服务节点。
	Add Operation = iota
	// Delete 表示删除一个已存在的服务节点。
	Delete
)

// Update 描述一次服务发现变更通知；同一 Update 中 Addr 为空字符串且 Metadata 为 nil
// 属于无效状态，调用方不应构造该组合。
type Update struct {
	// Op 为本次变更的操作类型。
	Op Operation
	// Addr 为变更后的节点地址；未更新地址时为空字符串。
	Addr string
	// Metadata 为变更后的节点元数据；未更新元数据时为 nil，自定义命名实现可不提供该字段。
	Metadata *server.ServiceInfo
}
