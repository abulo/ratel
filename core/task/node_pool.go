package task

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/driver"
)

// NodePool is a node pool
type NodePool struct {
	serviceName    string
	NodeID         string
	rwMut          sync.RWMutex
	nodes          *Map
	Driver         driver.Driver
	hashReplicas   int
	hashFn         Hash
	updateDuration time.Duration
	task           *Task
}

func newNodePool(serverName string, driver driver.Driver, task *Task, updateDuration time.Duration, hashReplicas int) (*NodePool, error) {

	err := driver.Ping()
	if err != nil {
		return nil, err
	}

	nodePool := &NodePool{
		Driver:         driver,
		serviceName:    serverName,
		task:           task,
		hashReplicas:   hashReplicas,
		updateDuration: updateDuration,
	}
	return nodePool, nil
}

// StartPool Start Service Watch Pool
func (np *NodePool) StartPool() error {
	var err error
	np.Driver.SetTimeout(np.updateDuration)
	np.NodeID, err = np.Driver.RegisterServiceNode(np.serviceName)
	if err != nil {
		return err
	}
	np.Driver.SetHeartBeat(np.NodeID)

	err = np.updatePool()
	if err != nil {
		return err
	}

	go np.tickerUpdatePool()
	return nil
}

func (np *NodePool) updatePool() error {
	nodes, err := np.Driver.GetServiceNodeList(np.serviceName)
	if err != nil {
		return err
	}

	np.rwMut.Lock()
	defer np.rwMut.Unlock()
	np.nodes = NewHash(np.hashReplicas, np.hashFn)
	for _, node := range nodes {
		np.nodes.Add(node)
	}
	return nil
}
func (np *NodePool) tickerUpdatePool() {
	tickers := time.NewTicker(np.updateDuration)
	for range tickers.C {
		if atomic.LoadInt32(&np.task.running) == taskRunning {
			err := np.updatePool()
			if err != nil {
				logger.Logger.Infof("update node pool error %+v", err)
			}
		} else {
			tickers.Stop()
			return
		}
	}
}

// PickNodeByJobName : 使用一致性hash算法根据任务名获取一个执行节点
func (np *NodePool) PickNodeByJobName(jobName string) string {
	np.rwMut.RLock()
	defer np.rwMut.RUnlock()
	if np.nodes.IsEmpty() {
		return ""
	}
	return np.nodes.Get(jobName)
}
