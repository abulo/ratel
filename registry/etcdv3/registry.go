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
	// defaultRetryTimes 为租约失效后重新批量注册的最大重试次数。
	defaultRetryTimes = 3
	// defaultRegisterTimeout 为单次重新注册等待与 keepalive 重连的超时时长。
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

// Kind 返回 etcdv3 注册器的类型标识 "etcdv3"。
func (reg *etcdv3Registry) Kind() string { return "etcdv3" }

// RegisterService 将服务业务信息与监控信息写入注册中心。
//
// 业务信息写入失败时直接返回错误；仅在业务注册成功后继续注册监控信息。
func (reg *etcdv3Registry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	err := reg.registerBiz(ctx, info)
	if err != nil {
		return err
	}
	return reg.registerMetric(ctx, info)
}

// UnregisterService 从注册中心删除指定服务节点的注册信息，删除成功时同步清理本地缓存。
func (reg *etcdv3Registry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.unregister(ctx, reg.registerKey(info))
}

// ListServices 按前缀查询注册中心中已注册的服务节点；已存在但反序列化失败的条目会被跳过。
func (reg *etcdv3Registry) ListServices(ctx context.Context, prefix string) (services []*server.ServiceInfo, err error) {
	getResp, getErr := reg.client.Get(ctx, prefix, clientv3.WithPrefix())
	if getErr != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err":    getErr,
			"prefix": prefix,
		}).Error(ecode.MsgWatchRequestErr, ecode.ErrKindRequestErr)
		return nil, getErr
	}

	// PERF-7-32 [P1] 两轮循环合并为一轮：同一 kv.Value 只反序列化一次到 registry.Update，
	// 优先取内嵌 Metadata 完整信息，缺失时退化为仅 Address 条目；返回切片结构与字段语义保持不变
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
		if service.Metadata != nil {
			services = append(services, service.Metadata)
		} else {
			services = append(services, &server.ServiceInfo{
				Address: service.Addr,
			})
		}
	}

	return
}

// WatchServices 监听指定前缀下的节点变更，返回携带当前全量快照的结果通道。
//
// 调用方在收到首个消息后持续从通道读取最新节点清单；监听建立失败时返回错误。
// 通道容量为 10，新事件写入失败时会被丢弃并记录告警，调用方应及时消费。
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
			// PERF-7-33 [P2] 当前事件先合并进累计快照，再非阻塞排空通道内已积压的同周期事件，
			// 一批事件只做一次全量 DeepCopy 与推送；快照内容与原逐事件推送的最终状态一致（事件频率/节点规模需实测验证，预期收益需实测验证）
			applyEvent := func(ev *mvccpb.KeyValue, op mvccpb.Event_EventType) {
				switch op {
				case mvccpb.PUT:
					updateAddrList(al, prefix, scheme, ev)
				case mvccpb.DELETE:
					deleteAddrList(al, prefix, scheme, ev)
				}
			}
			applyEvent(event.Kv, event.Type)
		drain:
			for {
				select {
				case ev, ok := <-watch.C():
					if !ok {
						break drain
					}
					applyEvent(ev.Kv, ev.Type)
				default:
					break drain
				}
			}

			out := al.DeepCopy()
			select {
			case addresses <- *out:
			default:
				logger.Logger.Warnf("invalid event")
			}
		}
	})

	return addresses, nil
}

// unregister 删除指定键；上下文未携带超时时间时按配置的读取超时补充超时控制。
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

// Close 取消监听上下文，并逐个注销本地缓存的服务节点；注销错误仅记录日志，不阻断其余节点清理。
func (reg *etcdv3Registry) Close() error {
	if reg.cancel != nil {
		reg.cancel()
	}
	var wg sync.WaitGroup
	// PERF-7-34 [P2] 信号量限制并发注销 goroutine 数量为 16，避免节点很多时一次性拉起等量 goroutine 同时向 etcd 发起批量 Delete
	sem := make(chan struct{}, 16)
	reg.kvs.Range(func(k, v any) bool {
		wg.Add(1)
		sem <- struct{}{}
		go func(k any) {
			defer wg.Done()
			defer func() { <-sem }()
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

// registerMetric 仅对监控类型服务写入 Prometheus 作业注册键，其余类型直接跳过。
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
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		// 每个服务实例复用同一个租约，减少向 etcd 重复申请租约的压力。
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
			// 使用 reg.ctx 在关闭注册器时统一取消 keepalive 循环。
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

// getOrGrantLeaseID 返回已创建的租约标识；尚未创建时按配置 TTL 向 etcd 申请新租约。
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

// doKeepalive 持续向 etcd 续租；租约失效或连接异常时重新批量注册本地缓存的键并重建租约，
// 直到注册器上下文被取消才退出循环。
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
				// 租约失效后重新批量注册本地缓存的键，并按固定次数重试。
				err := reg.registerAllKvs(cancelCtx)
				if err != nil {
					cancel()
					return
				}

				done <- struct{}{}
			}()

			select {
			case <-time.After(defaultRegisterTimeout):
				// 重新注册超时则取消上下文，并将租约标识置零进入下一轮重试。
				cancel()
				reg.setLeaseID(0)

				continue
			case <-done:
			}

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
				// keepalive 通道关闭时将租约标识置零，触发后续重新注册。
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

// registerAllKvs 遍历本地缓存的键重新写入注册中心，按固定次数重试；任一节点注册失败即停止继续遍历。
func (reg *etcdv3Registry) registerAllKvs(ctx context.Context) error {
	return util.Do(defaultRetryTimes, time.Second, func() error {
		var err error

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
			// 删除事件的 value 为空，无需解析元数据，直接移除对应节点条目。
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
