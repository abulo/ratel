package driver

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	redisClient "github.com/abulo/ratel/v3/stores/redis"
	"github.com/redis/go-redis/v9"
)

type RedisZSetDriver struct {
	c           *redisClient.Client
	serviceName string
	nodeID      string
	timeout     time.Duration
	started     bool

	// this context is used to define
	// the lifetime of this driver.
	runtimeCtx    context.Context
	runtimeCancel context.CancelFunc

	sync.Mutex
}

func newRedisZSetDriver(redisClient *redisClient.Client) *RedisZSetDriver {
	rd := &RedisZSetDriver{
		c:       redisClient,
		timeout: redisDefaultTimeout,
	}
	rd.started = false
	return rd
}

func (rd *RedisZSetDriver) Init(serviceName string, opts ...Option) {
	rd.serviceName = serviceName
	rd.nodeID = GetNodeId(serviceName)
	for _, opt := range opts {
		rd.WithOption(opt)
	}
}

func (rd *RedisZSetDriver) NodeID() string {
	return rd.nodeID
}

func (rd *RedisZSetDriver) GetNodes(ctx context.Context) (nodes []string, err error) {
	rd.Lock()
	defer rd.Unlock()
	val, err := rd.c.ZRangeByScore(ctx, GetKeyPre(rd.serviceName), &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", TimePre(time.Now(), rd.timeout)),
		Max: "+inf",
	})
	if err != nil {
		return nil, err
	}
	nodes = make([]string, len(val))
	copy(nodes, val)
	logger.Logger.Infof("nodes=%v", nodes)
	return
}
func (rd *RedisZSetDriver) Start(ctx context.Context) (err error) {
	rd.Lock()
	defer rd.Unlock()
	if rd.started {
		err = errors.New("this driver is started")
		return
	}
	rd.runtimeCtx, rd.runtimeCancel = context.WithCancel(context.TODO())
	rd.started = true
	// register
	err = rd.registerServiceNode()
	if err != nil {
		logger.Logger.Errorf("register service error=%v", err)
		return
	}
	// heartbeat timer
	go rd.heartBeat()
	return
}
func (rd *RedisZSetDriver) Stop(ctx context.Context) (err error) {
	rd.Lock()
	defer rd.Unlock()
	rd.runtimeCancel()
	rd.started = false
	return
}

func (rd *RedisZSetDriver) WithOption(opt Option) (err error) {
	switch opt.Type() {
	case OptionTypeTimeout:
		{
			rd.timeout = opt.(TimeoutOption).timeout
		}
	}
	return
}

// private function

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
				if _, err := rd.c.Del(context.Background(), rd.nodeID, rd.nodeID); err != nil {
					logger.Logger.Errorf("unregister service node error %+v", err)
				}
				return
			}
		}
	}
}

func (rd *RedisZSetDriver) registerServiceNode() error {
	_, err := rd.c.ZAdd(context.Background(), GetKeyPre(rd.serviceName), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: rd.nodeID,
	})
	return err
}
