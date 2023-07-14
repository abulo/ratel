package etcd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/task/driver"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	defaultTTL = 5 // 5 second ttl
	ctxTimeout = 5 * time.Second
)

type driverEtcd struct {
	driver *etcdv3.Client
	ctx    context.Context
	ttl    int64

	cancelCh chan struct{}

	aliveCh <-chan *clientv3.LeaseKeepAliveResponse
	leaseId clientv3.LeaseID
}

func NewDriver(driver *etcdv3.Client) driver.Driver {
	return &driverEtcd{driver: driver, ctx: context.Background(), cancelCh: make(chan struct{})}
}

func (e *driverEtcd) Ping() error { return nil }

func (e *driverEtcd) SetKeepaliveInterval(interval time.Duration) {
	ttl := int64(interval.Seconds())

	if ttl < defaultTTL {
		ttl = defaultTTL
	}

	e.ttl = ttl
}

func (e *driverEtcd) Keepalive(nodeId string) { go e.keepalive(nodeId) }

func (e *driverEtcd) GetServiceNodeList(serviceName string) (nodeIds []string, err error) {
	ctx, cancel := context.WithTimeout(e.ctx, ctxTimeout)
	defer cancel()

	//  /crond/{serviceName}/
	var prefix = driver.SPL + driver.PrefixKey + driver.SPL + serviceName + driver.SPL

	gets, err := e.driver.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return
	}

	for _, kv := range gets.Kvs {
		nodeIds = append(nodeIds, string(kv.Key))
	}

	return
}

func (e *driverEtcd) RegisterServiceNode(serviceName string) (nodeId string, err error) {
	//         /crond/{serviceName}/{uuid}
	nodeId = driver.SPL + driver.PrefixKey + driver.SPL + serviceName + driver.SPL + uuid.NewString()

	return nodeId, e.register(nodeId)
}

func (e *driverEtcd) UnRegisterServiceNode() { e.cancelCh <- struct{}{} }

func (e *driverEtcd) register(nodeId string) error {
	ctx, cancel := context.WithTimeout(e.ctx, ctxTimeout)
	defer cancel()

	grant, err := e.driver.Grant(ctx, e.ttl)
	if err != nil {
		return err
	}

	e.leaseId = grant.ID

	_, err = e.driver.Put(ctx, nodeId, "ok", clientv3.WithLease(e.leaseId))
	if err != nil {
		return err
	}

	e.aliveCh, err = e.driver.KeepAlive(context.Background(), e.leaseId)
	if err != nil {
		return err
	}

	return nil
}

func (e *driverEtcd) unregister(nodeId string) (err error) {
	_, er1 := e.driver.Delete(e.ctx, nodeId)

	_, er2 := e.driver.Revoke(e.ctx, e.leaseId)

	return CombineError(er1, er2)
}

func (e *driverEtcd) keepalive(nodeId string) {
	ticker := time.NewTicker(time.Duration(e.ttl) * time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-e.cancelCh:
			err := e.unregister(nodeId)
			if err != nil {
				log.Printf("error: node[%s] unregister failed: [%+v]", nodeId, err)
			}
			return

		case <-ticker.C:
			if e.aliveCh == nil {
				err := e.register(nodeId)
				if err != nil {
					log.Printf("error: node[%s] register failed: [%+v]", nodeId, err)
				}
			}

		case _, ok := <-e.aliveCh:
			if !ok {
				err := e.register(nodeId)
				if err != nil {
					log.Printf("error: node[%s] register failed: [%+v]", nodeId, err)
				}
			}
		}
	}
}

func CombineError(errs ...error) error {
	var errStr string
	for _, err := range errs {
		if err != nil {
			if errStr == "" {
				errStr += err.Error()
			} else {
				errStr += "; " + err.Error()
			}
		}
	}
	if errStr == "" {
		return nil
	}
	return errors.New(errStr)
}
