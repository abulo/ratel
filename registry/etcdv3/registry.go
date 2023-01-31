package etcdv3

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/constant"
	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/registry"
	"github.com/abulo/ratel/v3/server"
	"github.com/abulo/ratel/v3/util"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdv3Registry struct {
	ctx    context.Context
	client *etcdv3.Client
	kvs    sync.Map
	*Config
	cancel  context.CancelFunc
	rmu     *sync.RWMutex
	leaseID clientv3.LeaseID
	once    sync.Once
}

const (
	// defaultRetryTimes default retry times
	defaultRetryTimes = 3
	// defaultKeepAliveTimeout is the default timeout for keepalive requests.
	defaultRegisterTimeout = 5 * time.Second
)

func newETCDRegistry(config *Config) (*etcdv3Registry, error) {
	etcdv3Client, err := config.Config.Build()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	reg := &etcdv3Registry{
		ctx:    ctx,
		cancel: cancel,
		client: etcdv3Client,
		Config: config,
		kvs:    sync.Map{},
		rmu:    &sync.RWMutex{},
	}
	return reg, nil
}

// Kind ...
func (reg *etcdv3Registry) Kind() string { return "etcdv3" }

// RegisterService register service to registry
func (reg *etcdv3Registry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	err := reg.registerBiz(ctx, info)
	if err != nil {
		return err
	}
	return reg.registerMetric(ctx, info)
}

// UnregisterService unregister service from registry
func (reg *etcdv3Registry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.unregister(ctx, reg.registerKey(info))
}

// ListServices list service registered in registry with name `name`
func (reg *etcdv3Registry) ListServices(ctx context.Context, prefix string) (services []*server.ServiceInfo, err error) {
	getResp, getErr := reg.client.Get(ctx, prefix, clientv3.WithPrefix())
	if getErr != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err":    getErr,
			"prefix": prefix,
		}).Error(ecode.MsgWatchRequestErr, ecode.ErrKindRequestErr)
		return nil, getErr
	}

	for _, kv := range getResp.Kvs {
		var service server.ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("invalid service")
			continue
		}
		services = append(services, &service)
	}

	for _, kv := range getResp.Kvs {
		var service registry.Update
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
				"key": kv.Key,
				"val": kv.Value,
			}).Error("invalid service")
			continue
		}
		services = append(services, &server.ServiceInfo{
			Address: service.Addr,
		})
	}

	return
}

// WatchServices watch service change event, then return address list
func (reg *etcdv3Registry) WatchServices(ctx context.Context, prefix string) (chan registry.Endpoints, error) {
	watch, err := reg.client.WatchPrefix(context.Background(), prefix)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err":    err,
			"prefix": prefix,
			"val":    ecode.MsgWatchRequestErr,
		}).Error("reg.client.WatchPrefix failed")
		return nil, err
	}

	var addresses = make(chan registry.Endpoints, 10)
	var al = &registry.Endpoints{
		Nodes:           make(map[string]server.ServiceInfo),
		RouteConfigs:    make(map[string]registry.RouteConfig),
		ConsumerConfigs: make(map[string]registry.ConsumerConfig),
		ProviderConfigs: make(map[string]registry.ProviderConfig),
	}

	scheme := getScheme(prefix)

	for _, kv := range watch.IncipientKeyValues() {
		updateAddrList(al, prefix, scheme, kv)
	}
	addresses <- *al.DeepCopy()

	goroutine.Go(func() {
		for event := range watch.C() {
			switch event.Type {
			case mvccpb.PUT:
				updateAddrList(al, prefix, scheme, event.Kv)
			case mvccpb.DELETE:
				deleteAddrList(al, prefix, scheme, event.Kv)
			}

			// var snapshot registry.Endpoints
			out := al.DeepCopy()
			select {
			// case addresses <- snapshot:
			case addresses <- *out:
			default:
				logger.Logger.Warnf("invalid event")
			}
		}
	})

	return addresses, nil
}

func (reg *etcdv3Registry) unregister(ctx context.Context, key string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer cancel()
	}

	_, err := reg.client.Delete(ctx, key)
	if err == nil {
		reg.kvs.Delete(key)
	}
	return err
}

// Close ...
func (reg *etcdv3Registry) Close() error {
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
				logger.Logger.WithFields(logrus.Fields{
					"err": err,
					"key": k,
					"val": v,
				}).Error("unregister service")
			} else {
				logger.Logger.WithFields(logrus.Fields{
					"key": k,
					"val": v,
				}).Info("unregister service")
			}
			cancel()
		}(k)
		return true
	})
	wg.Wait()
	return nil
}

func (reg *etcdv3Registry) registerMetric(ctx context.Context, info *server.ServiceInfo) error {
	if info.Kind != constant.ServiceMonitor {
		return nil
	}
	metric := "/prometheus/job/%s/%s"
	val := info.Address
	key := fmt.Sprintf(metric, info.Name, env.HostName())
	return reg.registerKV(ctx, key, val)
}

func (reg *etcdv3Registry) registerBiz(ctx context.Context, info *server.ServiceInfo) error {
	key := reg.registerKey(info)
	val := reg.registerValue(info)

	return reg.registerKV(ctx, key, val)
}

func (reg *etcdv3Registry) registerKV(ctx context.Context, key, val string) error {
	opOptions := make([]clientv3.OpOption, 0)
	// opOptions = append(opOptions, clientv3.WithSerializable())
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		// 这里基于应用名为key做缓存，每个服务实例应该只需要创建一个lease，降低etcd的压力
		lease, err := reg.getOrGrantLeaseID(ctx)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"code": ecode.ErrKindRegisterErr,
				"err":  err,
				"val":  val,
				"key":  key,
			}).Error("getSession failed")
			return err
		}

		reg.once.Do(func() {
			// we use reg.ctx to manully cancel lease keepalive loop
			go reg.doKeepalive(reg.ctx)
		})
		opOptions = append(opOptions, clientv3.WithLease(lease))
	}
	_, err := reg.client.Put(ctx, key, val, opOptions...)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"code": ecode.ErrKindRegisterErr,
			"err":  err,
			"key":  key,
		}).Error("register service")
		return err
	}
	logger.Logger.WithFields(logrus.Fields{
		"val": val,
		"key": key,
	}).Info("register service")
	reg.kvs.Store(key, val)
	return nil
}

func (reg *etcdv3Registry) getOrGrantLeaseID(ctx context.Context) (clientv3.LeaseID, error) {
	reg.rmu.Lock()
	defer reg.rmu.Unlock()

	if reg.leaseID != 0 {
		return reg.leaseID, nil
	}
	grant, err := reg.client.Grant(ctx, int64(reg.ServiceTTL.Seconds()))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"code": ecode.ErrKindRegisterErr,
			"err":  err,
		}).Error("reg.client.Grant failed")
		return 0, err
	}
	reg.leaseID = grant.ID
	return grant.ID, nil
}

func (reg *etcdv3Registry) getLeaseID() clientv3.LeaseID {
	reg.rmu.RLock()
	defer reg.rmu.RUnlock()

	return reg.leaseID
}

func (reg *etcdv3Registry) setLeaseID(leaseId clientv3.LeaseID) {
	reg.rmu.Lock()
	defer reg.rmu.Unlock()

	reg.leaseID = leaseId
}

// doKeepAlive periodically sends keep alive requests to etcd server.
// when the keep alive request fails or timeout, it will try to re-establish the lease.
func (reg *etcdv3Registry) doKeepalive(ctx context.Context) {
	logger.Logger.Debug("start keepalive...")
	kac, err := reg.client.KeepAlive(ctx, reg.getLeaseID())
	if err != nil {
		reg.setLeaseID(0)
		logger.Logger.WithFields(logrus.Fields{
			"code": ecode.ErrKindRegisterErr,
			"err":  err,
		}).Error("reg.client.Grant failed")
	}
	for {
		if reg.getLeaseID() == 0 {
			cancelCtx, cancel := context.WithCancel(ctx)
			done := make(chan struct{}, 1)
			go func() {
				// do register again, and retry 3 times
				err := reg.registerAllKvs(cancelCtx)
				if err != nil {
					cancel()
					return
				}

				done <- struct{}{}
			}()

			// wait registerAllKvs success
			select {
			case <-time.After(defaultRegisterTimeout):
				// when timeout happens
				// we should cancel the context and retry again
				cancel()
				// mark leaseID as 0 to retry register
				reg.setLeaseID(0)

				continue
			case <-done:
				// when done happens, we just receive the kac channel
				// or wait the registry context done
			}

			// try do keepalive again
			// when error or timeout happens, just continue and try again
			kac, err = reg.client.KeepAlive(ctx, reg.getLeaseID())
			if err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"code": ecode.ErrKindRegisterErr,
					"err":  err,
				}).Error("reg.client.KeepAlive failed")

				time.Sleep(defaultRegisterTimeout)
				continue
			}
			logger.Logger.WithFields(logrus.Fields{
				"leaseid": reg.getLeaseID(),
			}).Debug("reg.client.KeepAlive finished")
		}

		select {
		case data, ok := <-kac:
			if !ok {
				// when error happens
				// mark leaseID as 0 to retry register
				reg.setLeaseID(0)
				logger.Logger.WithFields(logrus.Fields{
					"leaseid": reg.getLeaseID(),
				}).Debug("need to retry registration")
				continue
			}
			logger.Logger.WithFields(logrus.Fields{
				"leaseid": reg.getLeaseID(),
				"data":    data,
			}).Debug("do keepalive")
		case <-reg.ctx.Done():
			logger.Logger.Debug("exit keepalive")
			return
		}
	}
}

func (reg *etcdv3Registry) registerKey(info *server.ServiceInfo) string {
	return info.RegistryName()
}

func (reg *etcdv3Registry) registerValue(info *server.ServiceInfo) string {
	update := registry.Update{
		Op:       registry.Add,
		Addr:     info.Address,
		Metadata: info,
	}

	val, _ := json.Marshal(update)

	return string(val)
}

func (reg *etcdv3Registry) registerAllKvs(ctx context.Context) error {
	// do register again, and retry 3 times
	return util.Do(defaultRetryTimes, time.Second, func() error {
		var err error

		// all kvs stored in reg.kvs, and we can range this map to register again
		reg.kvs.Range(func(key, value any) bool {
			err = reg.registerKV(ctx, key.(string), value.(string))
			if err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"key":   key,
					"value": value,
					"code":  ecode.ErrKindRegisterErr,
					"err":   err,
				}).Error("registerKV failed")
			}
			return err == nil
		})
		return err
	})
}

func deleteAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)

		if isIPPort(addr) {
			// 直接删除addr 因为Delete操作的value值为空
			delete(al.Nodes, addr)
			delete(al.RouteConfigs, addr)
		}
	}
}

func updateAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		if isIPPort(addr) {
			var meta registry.Update
			if err := json.Unmarshal(kv.Value, &meta); err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"key":   kv.Key,
					"value": kv.Value,
					"err":   err,
				}).Error("unmarshal meta")
				continue
			}

			switch meta.Op {
			case registry.Add:
				al.Nodes[addr] = server.ServiceInfo{
					Address: addr,
				}
			case registry.Delete:
				delete(al.Nodes, addr)
			}
		}
	}
}

func isIPPort(addr string) bool {
	_, _, err := net.SplitHostPort(addr)
	return err == nil
}

func getScheme(prefix string) string {
	return strings.Split(prefix, ":")[0]
}
