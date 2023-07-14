package driver

import "time"

const (
	PrefixKey = "crond"
	JAR       = ":"
	SPL       = "/"
	REG       = "*"
)

//Driver is a driver interface
type Driver interface {
	// Ping 检查驱动是否可用
	Ping() error

	// SetKeepaliveInterval 设置心跳间隔
	SetKeepaliveInterval(interval time.Duration)

	// Keepalive 维持nodeId的心跳
	Keepalive(nodeId string)

	// GetServiceNodeList 获取serviceName所有node节点
	GetServiceNodeList(serviceName string) (nodeIds []string, err error)

	// RegisterServiceNode 向serviceName注册node节点
	RegisterServiceNode(serviceName string) (nodeId string, err error)

	// UnRegisterServiceNode 取消注册node节点
	UnRegisterServiceNode()
}
