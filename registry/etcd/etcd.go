package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/abulo/ratel/client/etcd"
	"github.com/abulo/ratel/goroutine"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/registry"
	"github.com/abulo/ratel/server"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// Config ...
type Config struct {
	*etcd.Config
	ReadTimeout time.Duration
	ServiceTTL  time.Duration
}

type etcdRegistry struct {
	client *etcd.Client
	kvs    sync.Map
	*Config
	Prefix   string
	cancel   context.CancelFunc
	rmu      *sync.RWMutex
	sessions map[string]*concurrency.Session
}

//Build ...
func Build(config *Config) *etcdRegistry {
	reg := &etcdRegistry{
		client:   config.Config.Build(),
		Config:   config,
		kvs:      sync.Map{},
		rmu:      &sync.RWMutex{},
		sessions: make(map[string]*concurrency.Session),
	}
	return reg
}

// RegisterService register service to registry
func (reg *etcdRegistry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	err := reg.registerBiz(ctx, info)
	if err != nil {
		return err
	}
	return reg.registerMetric(ctx, info)
}

// UnregisterService unregister service from registry
func (reg *etcdRegistry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.unregister(ctx, reg.registerKey(info))
}

// ListServices list service registered in registry with name `name`
func (reg *etcdRegistry) ListServices(ctx context.Context, name string, scheme string) (services []*server.ServiceInfo, err error) {
	target := fmt.Sprintf("/%s/%s/providers/%s://", reg.Prefix, name, scheme)
	getResp, getErr := reg.client.Get(ctx, target, clientv3.WithPrefix())
	if getErr != nil {
		logger.Logger.Error("watch request err ", getErr, target)
		return nil, getErr
	}

	for _, kv := range getResp.Kvs {
		var service server.ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			logger.Logger.Info("invalid service", err)
			continue
		}
		services = append(services, &service)
	}

	return
}

// WatchServices watch service change event, then return address list
func (reg *etcdRegistry) WatchServices(ctx context.Context, name string, scheme string) (chan registry.Endpoints, error) {
	prefix := fmt.Sprintf("/%s/%s/", reg.Prefix, name)
	watch, err := reg.client.WatchPrefix(context.Background(), prefix)
	if err != nil {
		return nil, err
	}

	var addresses = make(chan registry.Endpoints, 10)
	var al = &registry.Endpoints{
		Nodes: make(map[string]server.ServiceInfo),
	}

	for _, kv := range watch.IncipientKeyValues() {
		updateAddrList(al, prefix, scheme, kv)
	}

	addresses <- *al

	goroutine.Go(func() {
		for event := range watch.C() {
			al2 := reg.cloneEndPoints(al)
			switch event.Type {
			case mvccpb.PUT:
				updateAddrList(al2, prefix, scheme, event.Kv)
			case mvccpb.DELETE:
				deleteAddrList(al2, prefix, scheme, event.Kv)
			}

			select {
			case addresses <- *al2:
			default:
				logger.Logger.Warnf("invalid")
			}
		}
	})

	return addresses, nil
}

// Close ...
func (reg *etcdRegistry) Close() error {
	if reg.cancel != nil {
		reg.cancel()
	}
	var wg sync.WaitGroup
	reg.kvs.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(k interface{}) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := reg.unregister(ctx, k.(string))
			if err != nil {
				logger.Logger.Error("unregister service", err, k, v)
			} else {
				logger.Logger.Info("unregister service", k, v)
			}
			cancel()
		}(k)
		return true
	})
	wg.Wait()
	return nil
}

func updateAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		switch {
		// 解析服务注册键
		case strings.HasPrefix(addr, "providers/"+scheme):
			addr = strings.TrimPrefix(addr, "providers/")
			uri, err := url.Parse(addr)
			if err != nil {
				logger.Logger.Error("parse uri", err, kv.Key)
				continue
			}
			var serviceInfo server.ServiceInfo
			if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
				logger.Logger.Error("parse uri", err, kv.Key)
				continue
			}
			al.Nodes[uri.String()] = serviceInfo
		}
	}
}

func deleteAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		if strings.HasPrefix(addr, "providers/"+scheme) {
			// 解析服务注册键
			addr = strings.TrimPrefix(addr, "providers/")
			if addr == "" {
				continue
			}
			uri, err := url.Parse(addr)
			if err != nil {
				logger.Logger.Error("parse uri", err, kv.Key)
				continue
			}
			delete(al.Nodes, uri.String())
		}
	}
}

func (reg *etcdRegistry) registerBiz(ctx context.Context, info *server.ServiceInfo) error {
	var readCtx context.Context
	var readCancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		readCtx, readCancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer readCancel()
	}
	key := reg.registerKey(info)
	val := reg.registerValue(info)

	opOptions := make([]clientv3.OpOption, 0)
	// opOptions = append(opOptions, clientv3.WithSerializable())
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		//todo ctx without timeout for same as service life?
		sess, err := reg.getSession(key, concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	}

	_, err := reg.client.Put(readCtx, key, val, opOptions...)
	if err != nil {
		logger.Logger.Error("register service", err, info)
		return err
	}
	logger.Logger.Info("register service", key, val)
	reg.kvs.Store(key, val)
	return nil

}

func (reg *etcdRegistry) registerMetric(ctx context.Context, info *server.ServiceInfo) error {
	metric := "/prometheus/job/%s/%s/%s"

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer cancel()
	}
	key := fmt.Sprintf(metric, info.Name, info.Scheme, info.Address)
	opOptions := make([]clientv3.OpOption, 0)
	// opOptions = append(opOptions, clientv3.WithSerializable())
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		//todo ctx without timeout for same as service life?
		sess, err := reg.getSession(key, concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	}
	_, err := reg.client.Put(ctx, key, info.Address, opOptions...)
	if err != nil {
		logger.Logger.Error("register service", err, key, info)
		return err
	}

	logger.Logger.Info("register service", key, info.Address)
	reg.kvs.Store(key, info.Address)
	return nil
}

func (reg *etcdRegistry) registerKey(info *server.ServiceInfo) string {
	return registry.GetServiceKey(reg.Prefix, info)
}

func (reg *etcdRegistry) registerValue(info *server.ServiceInfo) string {
	return registry.GetServiceValue(info)
}

func (reg *etcdRegistry) getSession(k string, opts ...concurrency.SessionOption) (*concurrency.Session, error) {
	reg.rmu.RLock()
	sess, ok := reg.sessions[k]
	reg.rmu.RUnlock()
	if ok {
		return sess, nil
	}
	sess, err := concurrency.NewSession(reg.client.Client)
	if err != nil {
		return sess, err
	}
	reg.rmu.Lock()
	reg.sessions[k] = sess
	reg.rmu.Unlock()
	return sess, nil
}
func (reg *etcdRegistry) delSession(k string) error {
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		reg.rmu.RLock()
		sess, ok := reg.sessions[k]
		reg.rmu.RUnlock()
		if ok {
			reg.rmu.Lock()
			delete(reg.sessions, k)
			reg.rmu.Unlock()
			if err := sess.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (reg *etcdRegistry) unregister(ctx context.Context, key string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer cancel()
	}

	if err := reg.delSession(key); err != nil {
		return err
	}

	_, err := reg.client.Delete(ctx, key)
	if err == nil {
		reg.kvs.Delete(key)
	}
	return err
}

func (reg *etcdRegistry) cloneEndPoints(src *registry.Endpoints) *registry.Endpoints {
	dst := &registry.Endpoints{
		Nodes: make(map[string]server.ServiceInfo),
	}
	for k, v := range src.Nodes {
		dst.Nodes[k] = v
	}

	return dst
}
