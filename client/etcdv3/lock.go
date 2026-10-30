package etcdv3

import (
	"context"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"
)

// Mutex 是基于 etcd 租约与事务实现的分布式互斥锁，由 Client.NewMutex 创建。
//
// 零值不可直接使用；锁与绑定的租约会话一一对应，Unlock 后会话被关闭，不可重复使用。
type Mutex struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

// NewMutex 为给定 key 创建分布式互斥锁并建立绑定的租约会话；会话创建失败时返回错误。
// 未通过 SessionOption 指定 TTL 时，租约 TTL 使用 etcd concurrency 包默认值 60 秒。
func (client *Client) NewMutex(key string, opts ...concurrency.SessionOption) (mutex *Mutex, err error) {
	mutex = &Mutex{}
	mutex.s, err = concurrency.NewSession(client.Client, opts...)
	if err != nil {
		return
	}
	mutex.m = concurrency.NewMutex(mutex.s, key)
	return
}

// Lock 在 timeout 时长内阻塞等待获取锁；超时或获取失败时返回错误。
func (mutex *Mutex) Lock(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}

// TryLock 在 timeout 时长内获取锁；当前实现与 Lock 行为一致，同样阻塞等待，不做非阻塞探测。
func (mutex *Mutex) TryLock(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}

// Unlock 释放持有的锁并关闭绑定的租约会话；释放失败时直接返回错误，不再关闭会话。
func (mutex *Mutex) Unlock() (err error) {
	err = mutex.m.Unlock(context.TODO())
	if err != nil {
		return
	}
	return mutex.s.Close()
}
