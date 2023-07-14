package task

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/task/driver"
	"github.com/abulo/ratel/v3/core/task/hash"
)

type Node struct {
	mu sync.Mutex

	serviceName    string
	nodeId         string
	updateInterval time.Duration

	driver driver.Driver
	crond  *Crond

	nodes *hash.ConsistentHash
}

func newNode(serviceName string, driver driver.Driver, crond *Crond, updateInterval time.Duration) *Node {
	err := driver.Ping()
	if err != nil {
		panic(err)
	}

	return &Node{
		serviceName:    serviceName,
		driver:         driver,
		crond:          crond,
		updateInterval: updateInterval,
	}
}

func (n *Node) Start() error {
	n.driver.SetKeepaliveInterval(n.updateInterval)

	nodeId, err := n.driver.RegisterServiceNode(n.serviceName)
	if err != nil {
		return err
	}

	n.driver.Keepalive(nodeId)

	n.nodeId = nodeId

	if n.crond.lazyPick {
		return nil
	}

	err = n.ttl()
	if err != nil {
		return err
	}

	go n.tickerTTL()

	return nil
}

func (n *Node) Stop() {
	n.mu.Lock()
	n.driver.UnRegisterServiceNode()
	n.mu.Unlock()
}

func (n *Node) ttl() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	nodeIds, err := n.driver.GetServiceNodeList(n.serviceName)
	if err != nil {
		return err
	}

	n.nodes = hash.NewConsistentHash().Add(nodeIds...)

	return nil
}

func (n *Node) tickerTTL() {
	ticker := time.NewTicker(n.updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		if n.crond.isRunning {
			err := n.ttl()
			if err != nil {
				log.Printf("error: update node failed: [%+v]", err)
			}
		} else {
			return
		}
	}
}

func (n *Node) pickNode(jobName string) (string, error) {
	if n.crond.lazyPick {
		return n.pickLazy(jobName)
	}

	return n.pickHung(jobName)
}

func (n *Node) pickHung(jobName string) (string, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.nodes.IsEmpty() {
		return "", errors.New("empty nodes")
	}

	return n.nodes.Get(jobName), nil
}

func (n *Node) pickLazy(jobName string) (string, error) {
	nodeIds, err := n.driver.GetServiceNodeList(n.serviceName)
	if err != nil {
		return "", err
	}

	if len(nodeIds) <= 0 {
		return "", errors.New("empty nodes")
	}

	nodes := hash.NewConsistentHash().Add(nodeIds...)

	return nodes.Get(jobName), nil
}
