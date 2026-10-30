package task

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/consistenthash"
	"github.com/abulo/ratel/v3/core/task/driver"
)

const (
	// NodePoolStateSteady 表示集群成员列表稳定，当前节点可以运行任务。
	NodePoolStateSteady = "NodePoolStateSteady"
	// NodePoolStateUpgrade 表示集群成员列表发生变更，节点处于升级状态，暂不运行任务。
	NodePoolStateUpgrade = "NodePoolStateUpgrade"
)

// NodePool 基于一致性哈希维护集群成员列表，在稳定与升级两种状态间切换，并判定任务归属节点。
//
// 状态为 Steady 时当前节点允许运行分配到自己的任务；状态为 Upgrade 时节点不承接任务，
// 等待成员列表再次稳定。成员列表由读写锁保护，可并发查询。
type NodePool struct {
	serviceName    string
	nodeID         string
	rwMut          sync.RWMutex
	nodes          *consistenthash.Map
	driver         driver.DriverV2
	hashReplicas   int
	hashFn         consistenthash.Hash
	updateDuration time.Duration
	stopChan       chan int
	// preNodes 为上一次已知的节点列表，已排序，用于比对成员是否发生变更。
	preNodes            []string
	lastUpdateNodesTime atomic.Value
	state               atomic.Value
}

// NewNodePool 创建节点池；serviceName 为服务名，driver 为节点注册与发现驱动，
// updateDuration 为成员列表刷新间隔，hashReplicas 为一致性哈希虚拟副本数。
func NewNodePool(
	serviceName string,
	drv driver.DriverV2,
	updateDuration time.Duration,
	hashReplicas int,
) INodePool {
	np := &NodePool{
		serviceName:    serviceName,
		driver:         drv,
		hashReplicas:   hashReplicas,
		updateDuration: updateDuration,
		stopChan:       make(chan int, 1),
	}
	np.driver.Init(serviceName,
		driver.NewTimeoutOption(updateDuration))
	return np
}

// Start 启动驱动，拉取初始成员列表，并阻塞等待集群进入稳定状态后返回；
// 驱动启动或成员列表获取失败时返回错误。
func (np *NodePool) Start(ctx context.Context) (err error) {
	err = np.driver.Start(ctx)
	if err != nil {
		logger.Logger.Errorf("start pool error: %v", err)
		return
	}
	np.nodeID = np.driver.NodeID()
	nowNodes, err := np.driver.GetNodes(ctx)
	if err != nil {
		logger.Logger.Errorf("get nodes error: %v", err)
		return
	}
	np.state.Store(NodePoolStateUpgrade)
	np.updateHashRing(nowNodes)
	go np.waitingForHashRing()

	// 启动时先阻塞到集群状态稳定，避免新节点在成员列表尚未收敛时分发任务。
	for np.getState() != NodePoolStateSteady {
		<-time.After(np.updateDuration)
	}
	logger.Logger.Infof("nodepool started for serve, nodeID=%s", np.nodeID)

	return
}

// CheckJobAvailable 判断任务 jobName 是否被一致性哈希分配到当前节点。
//
// 哈希环未初始化时返回 ErrNodePoolIsNil；哈希环为空时返回 false 且无错误；
// 处于升级状态时返回 ErrNodePoolIsUpgrading。
func (np *NodePool) CheckJobAvailable(jobName string) (bool, error) {
	np.rwMut.RLock()
	defer np.rwMut.RUnlock()
	if np.nodes == nil {
		return false, ErrNodePoolIsNil
	}
	if np.nodes.IsEmpty() {
		return false, nil
	}
	if np.state.Load().(string) != NodePoolStateSteady {
		return false, ErrNodePoolIsUpgrading
	}
	targetNode := np.nodes.Get(jobName)
	newNpNodeId := np.nodeID

	return newNpNodeId == targetNode, nil
}

// Stop 停止成员列表刷新协程，并通知驱动下线当前节点。
func (np *NodePool) Stop(ctx context.Context) error {
	np.stopChan <- 1
	np.driver.Stop(ctx)
	np.preNodes = make([]string, 0)
	return nil
}

// GetNodeID 返回当前节点标识。
func (np *NodePool) GetNodeID() string {
	return np.nodeID
}

// GetLastNodesUpdateTime 返回最近一次成员列表发生变化并重建哈希环的时间；从未更新时为零值时间。
func (np *NodePool) GetLastNodesUpdateTime() time.Time {
	return np.lastUpdateNodesTime.Load().(time.Time)
}

func (np *NodePool) getState() string {
	return np.state.Load().(string)
}

func (np *NodePool) waitingForHashRing() {
	tick := time.NewTicker(np.updateDuration)
	for {
		select {
		case <-tick.C:
			nowNodes, err := np.driver.GetNodes(context.Background())
			if err != nil {
				logger.Logger.Errorf("get nodes error %v", err)
				continue
			}
			np.updateHashRing(nowNodes)
		case <-np.stopChan:
			return
		}
	}
}

func (np *NodePool) updateHashRing(nodes []string) {
	np.rwMut.Lock()
	defer np.rwMut.Unlock()
	if np.equalRing(nodes) {
		np.state.Store(NodePoolStateSteady)
		return
	}
	np.lastUpdateNodesTime.Store(time.Now())
	np.state.Store(NodePoolStateUpgrade)
	np.preNodes = make([]string, len(nodes))
	copy(np.preNodes, nodes)
	np.nodes = consistenthash.New(np.hashReplicas, np.hashFn)
	for _, v := range nodes {
		np.nodes.Add(v)
	}
}

// equalRing 比对新成员列表与上一次列表是否一致；注意会就地排序入参 a。
func (np *NodePool) equalRing(a []string) bool {
	if len(a) == len(np.preNodes) {
		la := len(a)
		sort.Strings(a)
		for i := range la {
			if a[i] != np.preNodes[i] {
				return false
			}
		}
		return true
	}
	return false
}
