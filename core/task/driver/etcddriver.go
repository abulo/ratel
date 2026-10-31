package driver

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 驱动默认参数：最小租约时长、拨号超时与业务操作超时。
const (
	etcdDefaultLease    = 5
	etcdDialTimeout     = 3 * time.Second
	etcdBusinessTimeout = 5 * time.Second
)

// EtcdDriver 基于 etcd 租约与前缀监听维护在线节点列表：当前节点携带租约心跳注册，
// 其他节点通过 Watch 前缀实时增删本地节点缓存。
type EtcdDriver struct {
	nodeID      string
	serviceName string
	cli         *etcdv3.Client
	nodes       *sync.Map
	lease       int64
	leaseID     clientv3.LeaseID
	leaseCh     <-chan *clientv3.LeaseKeepAliveResponse
	ctx         context.Context
	cancel      context.CancelFunc
	// PERF-3-9 [P1]：started 配合 mu 防止 Stop 后再次 Start 时重复起 watcher/心跳协程。
	mu      sync.Mutex
	started bool
}

func newEtcdDriver(cli *etcdv3.Client) *EtcdDriver {
	ser := &EtcdDriver{
		cli:   cli,
		nodes: &sync.Map{},
	}

	return ser
}

// putKeyWithLease 申请租约并以该租约写入 key/value；租约时长小于默认值时按默认值处理。
func (e *EtcdDriver) putKeyWithLease(ctx context.Context, key, val string) (clientv3.LeaseID, error) {
	if e.lease < etcdDefaultLease {
		e.lease = etcdDefaultLease
	}

	subCtx, cancel := context.WithTimeout(ctx, etcdBusinessTimeout)
	defer cancel()
	resp, err := e.cli.Grant(subCtx, e.lease)
	if err != nil {
		return 0, err
	}
	_, err = e.cli.Put(subCtx, key, val, clientv3.WithLease(resp.ID))
	if err != nil {
		return 0, err
	}

	return resp.ID, nil
}

// watchService 先按前缀载入现有节点，再启动协程监听前缀变更以维护节点缓存。
func (e *EtcdDriver) watchService(ctx context.Context, serviceName string) error {
	prefix := GetKeyPre(serviceName)
	resp, err := e.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		e.setServiceList(string(ev.Key), string(ev.Value))
	}

	go e.watcher(serviceName)
	return nil
}

// watcher 持续监听服务前缀，PUT 事件写入缓存，DELETE 事件移除缓存。
func (e *EtcdDriver) watcher(serviceName string) {
	prefix := GetKeyPre(serviceName)
	// PERF-3-9 [P1]：Watch 绑定驱动生命周期上下文 e.ctx；Stop 调用 cancel 后该 range 随之退出，
	// 不再因 context.Background() 导致 watcher goroutine 与底层 watch 流残留泄漏。
	rch := e.cli.Watch(e.ctx, prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				e.setServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				e.delServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (e *EtcdDriver) setServiceList(key, val string) {
	e.nodes.Store(key, val)
}

func (e *EtcdDriver) delServiceList(key string) {
	e.nodes.Delete(key)
}

func (e *EtcdDriver) getServices() []string {
	addrs := make([]string, 0)
	e.nodes.Range(func(key, _ any) bool {
		addrs = append(addrs, key.(string))
		return true
	})
	return addrs
}

func (e *EtcdDriver) keepAlive(ctx context.Context, nodeID string) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	var err error
	e.leaseID, err = e.putKeyWithLease(ctx, nodeID, nodeID)
	if err != nil {
		logger.Logger.Errorf("putKeyWithLease error: %v", err)
		return nil, err
	}

	return e.cli.KeepAlive(ctx, e.leaseID)
}

func (e *EtcdDriver) revoke(ctx context.Context) {
	_, err := e.cli.Lease.Revoke(ctx, e.leaseID)
	if err != nil {
		logger.Logger.Infof("lease revoke error: %v", err)
	}
}

func (e *EtcdDriver) startHeartBeat(ctx context.Context) {
	var err error
	e.leaseCh, err = e.keepAlive(ctx, e.nodeID)
	if err != nil {
		logger.Logger.Errorf("keep alive error, %v", err)
		return
	}
}

// keepHeartBeat 监听驱动停止信号与租约通道：任一结束即退出协程，不再保活。
func (e *EtcdDriver) keepHeartBeat() {
	for {
		select {
		case <-e.ctx.Done():
			{
				logger.Logger.Warnf("driver stopped")
				return
			}
		case _, ok := <-e.leaseCh:
			{
				if !ok {
					logger.Logger.Warnf("lease channel stop, driver stopped")
					return
				}
			}
		}
	}
}

// Init 实现 DriverV2：记录服务名、生成节点标识并应用给定配置选项。
func (e *EtcdDriver) Init(serverName string, opts ...Option) {
	e.serviceName = serverName
	e.nodeID = GetNodeId(serverName)
	for _, opt := range opts {
		e.WithOption(opt)
	}
}

// NodeID 返回当前节点标识。
func (e *EtcdDriver) NodeID() string {
	return e.nodeID
}

// GetNodes 实现 DriverV2：返回本地缓存的全部在线节点标识。
func (e *EtcdDriver) GetNodes(ctx context.Context) (nodes []string, err error) {
	return e.getServices(), nil
}

// Start 实现 DriverV2：每次启动重建驱动级上下文，开启租约保活与前缀监听，随后启动保活协程；
// 重复 Start 在驱动未 Stop 时直接返回错误，避免重复注册 watcher 与心跳协程。
func (e *EtcdDriver) Start(ctx context.Context) (err error) {
	e.mu.Lock()
	if e.started {
		e.mu.Unlock()
		return errors.New("this driver is started")
	}
	e.ctx, e.cancel = context.WithCancel(context.TODO())
	e.started = true
	e.mu.Unlock()

	e.startHeartBeat(ctx)
	err = e.watchService(ctx, e.serviceName)
	if err != nil {
		// PERF-3-9 [P1]：监听启动失败时回滚上下文与启动标志，避免残留半启动状态。
		e.cancel()
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
		return
	}
	go e.keepHeartBeat()
	return
}

// Stop 实现 DriverV2：撤销租约使节点键失效，并取消驱动上下文停止监听与保活协程。
func (e *EtcdDriver) Stop(ctx context.Context) (err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.revoke(ctx)
	if e.cancel != nil {
		e.cancel()
	}
	e.started = false
	return nil
}

// WithOption 实现 DriverV2：按选项类型应用配置，当前仅支持超时选项（映射为租约时长）。
func (e *EtcdDriver) WithOption(opt Option) (err error) {
	switch opt.Type() {
	case OptionTypeTimeout:
		{
			e.lease = int64(opt.(TimeoutOption).timeout.Seconds())
		}
	}
	return
}
