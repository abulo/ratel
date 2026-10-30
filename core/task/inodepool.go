package task

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNodePoolIsUpgrading 表示节点池处于集群成员变更的升级状态，此时任务不可在当前节点调度。
	ErrNodePoolIsUpgrading = errors.New("nodePool is upgrading")
	// ErrNodePoolIsNil 表示节点池尚未初始化，例如本地模式或节点池未启动。
	ErrNodePoolIsNil = errors.New("nodePool is nil")
)

// INodePool 抽象集群节点池，负责节点注册、成员发现与任务归属判定。
//
// 实现须保证 Start 与 Stop 可被重复调用，GetNodes 在超时时返回错误。
type INodePool interface {
	// Start 启动节点池，完成节点注册并开始维护集群成员列表。
	Start(ctx context.Context) error
	// CheckJobAvailable 判断任务 jobName 是否被一致性哈希分配到当前节点。
	CheckJobAvailable(jobName string) (bool, error)
	// Stop 停止节点池并注销当前节点。
	Stop(ctx context.Context) error

	// GetNodeID 返回当前节点标识。
	GetNodeID() string
	// GetLastNodesUpdateTime 返回最近一次集群成员列表更新时间。
	GetLastNodesUpdateTime() time.Time
}
