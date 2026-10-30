package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/redis/go-redis/v9"
)

// redisDefaultTimeout 为未显式配置超时时的默认节点过期时间。
const (
	redisDefaultTimeout = 5 * time.Second
)

// RedisDriver 使用普通 Redis 键值对维护在线节点：节点以带过期时间的键注册，心跳定时续期；
// 驱动通过内部互斥锁保证 Start、Stop 的并发安全。
type RedisDriver struct {
	c           redis.UniversalClient
	serviceName string
	nodeID      string
	timeout     time.Duration
	started     bool
	// runtimeCtx 定义驱动自身生命周期，Stop 取消该上下文以通知心跳协程退出。
	runtimeCtx    context.Context
	runtimeCancel context.CancelFunc

	sync.Mutex
}

func newRedisDriver(redisClient redis.UniversalClient) *RedisDriver {
	rd := &RedisDriver{
		c:       redisClient,
		timeout: redisDefaultTimeout,
	}
	rd.started = false
	return rd
}

// Init 实现 DriverV2：记录服务名、生成节点标识并应用给定配置选项。
func (rd *RedisDriver) Init(serviceName string, opts ...Option) {
	rd.serviceName = serviceName
	rd.nodeID = GetNodeId(rd.serviceName)

	for _, opt := range opts {
		rd.WithOption(opt)
	}
}

// NodeID 返回当前节点标识。
func (rd *RedisDriver) NodeID() string {
	return rd.nodeID
}

// Start 实现 DriverV2：重复启动返回错误；注册节点成功后启动心跳协程。
func (rd *RedisDriver) Start(ctx context.Context) (err error) {
	rd.Lock()
	defer rd.Unlock()
	if rd.started {
		err = errors.New("this driver is started")
		return
	}
	rd.runtimeCtx, rd.runtimeCancel = context.WithCancel(context.TODO())
	rd.started = true
	err = rd.registerServiceNode()
	if err != nil {
		return
	}
	go rd.heartBeat()
	return
}

// Stop 实现 DriverV2：取消驱动生命周期上下文，使心跳协程注销节点并退出，随后标记已停止。
func (rd *RedisDriver) Stop(ctx context.Context) (err error) {
	rd.Lock()
	defer rd.Unlock()
	rd.runtimeCancel()
	rd.started = false
	return
}

// GetNodes 实现 DriverV2：按服务名前缀扫描全部在线节点键。
func (rd *RedisDriver) GetNodes(ctx context.Context) (nodes []string, err error) {
	mathStr := fmt.Sprintf("%s*", GetKeyPre(rd.serviceName))
	return rd.scan(ctx, mathStr)
}

// heartBeat 按超时时长一半的间隔续期节点键；收到停止信号时删除节点键后退出。
func (rd *RedisDriver) heartBeat() {
	tick := time.NewTicker(rd.timeout / 2)
	for {
		select {
		case <-tick.C:
			{
				if err := rd.registerServiceNode(); err != nil {
					logger.Logger.Errorf("register service node error %+v", err)
				}
			}
		case <-rd.runtimeCtx.Done():
			{
				if err := rd.c.Del(context.Background(), rd.nodeID, rd.nodeID).Err(); err != nil {
					logger.Logger.Errorf("unregister service node error %+v", err)
				}
				return
			}
		}
	}
}

func (rd *RedisDriver) registerServiceNode() error {
	return rd.c.SetEx(context.Background(), rd.nodeID, rd.nodeID, rd.timeout).Err()
}

func (rd *RedisDriver) scan(ctx context.Context, matchStr string) ([]string, error) {
	ret := make([]string, 0)
	iter := rd.c.Scan(ctx, 0, matchStr, -1).Iterator()
	for iter.Next(ctx) {
		err := iter.Err()
		if err != nil {
			return nil, err
		}
		ret = append(ret, iter.Val())
	}
	return ret, nil
}

// WithOption 实现 DriverV2：按选项类型应用配置，当前仅支持超时选项。
func (rd *RedisDriver) WithOption(opt Option) (err error) {
	switch opt.Type() {
	case OptionTypeTimeout:
		{
			rd.timeout = opt.(TimeoutOption).timeout
		}
	}
	return
}
