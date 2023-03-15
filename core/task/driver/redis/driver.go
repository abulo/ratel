package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/driver"
	"github.com/abulo/ratel/v3/stores/redis"
)

// RedisDriver is redisDriver
type RedisDriver struct {
	client  *redis.Client
	timeout time.Duration
	Key     string
}

// NewDriver return a redis driver
func NewDriver(client *redis.Client) (*RedisDriver, error) {
	return &RedisDriver{
		client: client,
	}, nil
}

// Ping is check redis valid
func (rd *RedisDriver) Ping() error {
	reply, err := rd.client.Ping(context.Background())
	if err != nil {
		return err
	}
	if reply != "PONG" {
		return fmt.Errorf("Ping received is error, %s", string(reply))
	}
	return err
}

// SetTimeout set redis timeout
func (rd *RedisDriver) SetTimeout(timeout time.Duration) {
	rd.timeout = timeout
}

// SetHeartBeat set heatbeat
func (rd *RedisDriver) SetHeartBeat(nodeID string) {
	go rd.heartBeat(nodeID)
}
func (rd *RedisDriver) heartBeat(nodeID string) {

	//每间隔timeout/2设置一次key的超时时间为timeout
	key := nodeID
	tickers := time.NewTicker(rd.timeout / 2)
	for range tickers.C {
		keyExist, err := rd.client.Expire(context.Background(), key, rd.timeout)
		if err != nil {
			logger.Logger.Errorf("redis expire error %+v", err)
			continue
		}
		if !keyExist {
			if err := rd.registerServiceNode(nodeID); err != nil {
				logger.Logger.Errorf("register service node error %+v", err)
			}
		}
	}
}

// GetServiceNodeList get a serveice node  list
func (rd *RedisDriver) GetServiceNodeList(serviceName string) ([]string, error) {
	mathStr := fmt.Sprintf("%s*", driver.GetKeyPre(serviceName))
	return rd.scan(mathStr)
}

// RegisterServiceNode  register a service node
func (rd *RedisDriver) RegisterServiceNode(serviceName string) (nodeID string, err error) {
	nodeID = driver.GetNodeId(serviceName)
	if err := rd.registerServiceNode(nodeID); err != nil {
		return "", err
	}
	return nodeID, nil
}

func (rd *RedisDriver) registerServiceNode(nodeID string) error {
	_, err := rd.client.Set(context.Background(), nodeID, nodeID, rd.timeout)
	return err
}

func (rd *RedisDriver) scan(matchStr string) ([]string, error) {
	ret := make([]string, 0)
	ctx := context.Background()
	iter, err := rd.client.ScanIterator(ctx, 0, matchStr, -1)
	if err != nil {
		return nil, err
	}
	for iter.Next(ctx) {
		err := iter.Err()
		if err != nil {
			return nil, err
		}
		ret = append(ret, iter.Val())
	}
	return ret, nil
}
