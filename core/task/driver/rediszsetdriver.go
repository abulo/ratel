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

// RedisZSetDriver 使用 Redis 有序集合维护在线节点：成员为节点标识，分数为注册时间戳，
// 过期节点由 GetNodes 按时间戳下限自动剔除；驱动通过内部互斥锁保证 Start、Stop 的并发安全。
type RedisZSetDriver struct {
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

func newRedisZSetDriver(redisClient redis.UniversalClient) *RedisZSetDriver {
	rd := &RedisZSetDriver{
		c:       redisClient,
		timeout: redisDefaultTimeout,
	}
	rd.started = false
	return rd
}

// Init 实现 DriverV2：记录服务名、生成节点标识并应用给定配置选项。
func (rd *RedisZSetDriver) Init(serviceName string, opts ...Option) {
	rd.serviceName = serviceName
	rd.nodeID = GetNodeId(serviceName)
	for _, opt := range opts {
		rd.WithOption(opt)
	}
}

// NodeID 返回当前节点标识。
func (rd *RedisZSetDriver) NodeID() string {
	return rd.nodeID
}

// GetNodes 实现 DriverV2：仅返回分数不早于当前时间减去超时时长的在线节点。
func (rd *RedisZSetDriver) GetNodes(ctx context.Context) (nodes []string, err error) {
	rd.Lock()
	defer rd.Unlock()
	sliceCmd := rd.c.ZRangeByScore(ctx, GetKeyPre(rd.serviceName), &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", TimePre(time.Now(), rd.timeout)),
		Max: "+inf",
	})
	if err = sliceCmd.Err(); err != nil {
		return nil, err
	} else {
		nodes = make([]string, len(sliceCmd.Val()))
		copy(nodes, sliceCmd.Val())
	}
	return
}

// Start 实现 DriverV2：重复启动返回错误；注册节点成功后启动心跳协程。
func (rd *RedisZSetDriver) Start(ctx context.Context) (err error) {
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
		logger.Logger.Errorf("register service error=%v", err)
		return
	}
	go rd.heartBeat()
	return
}

// Stop 实现 DriverV2：取消驱动生命周期上下文，使心跳协程注销节点并退出，随后标记已停止。
func (rd *RedisZSetDriver) Stop(ctx context.Context) (err error) {
	rd.Lock()
	defer rd.Unlock()
	rd.runtimeCancel()
	rd.started = false
	return
}

// WithOption 实现 DriverV2：按选项类型应用配置，当前仅支持超时选项。
func (rd *RedisZSetDriver) WithOption(opt Option) (err error) {
	switch opt.Type() {
	case OptionTypeTimeout:
		{
			rd.timeout = opt.(TimeoutOption).timeout
		}
	}
	return
}

// heartBeat 按超时时长一半的间隔刷新节点分数；收到停止信号时尝试清理节点键后退出协程。
func (rd *RedisZSetDriver) heartBeat() {
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

func (rd *RedisZSetDriver) registerServiceNode() error {

	return rd.c.ZAdd(context.Background(), GetKeyPre(rd.serviceName), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: rd.nodeID,
	}).Err()
}
