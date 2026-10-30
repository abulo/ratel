package etcdv3

import (
	"context"

	"github.com/abulo/ratel/v3/core/ecode"
	"github.com/abulo/ratel/v3/core/goroutine"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watch 封装对某个前缀的变更监听，监听开始时的既有键值通过 IncipientKeyValues 取得，
// 后续变更事件通过 C 返回的通道投递；由 Client.WatchPrefix 创建，使用结束后须调用 Close 停止监听。
type Watch struct {
	revision  int64
	cancel    context.CancelFunc
	eventChan chan *clientv3.Event

	incipientKVs []*mvccpb.KeyValue
}

// C 返回用于接收变更事件的通道；通道缓冲为 100，消费者须及时读取，积压时新事件会被丢弃。
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

// IncipientKeyValues 返回监听建立时该前缀下已存在的键值快照，仅在创建 Watch 时填充一次。
func (w *Watch) IncipientKeyValues() []*mvccpb.KeyValue {
	return w.incipientKVs
}

// WatchPrefix 监听以 prefix 为前缀的全部键的后续变更。
//
// 方法先按前缀读取当前修订版本作为基线，既有键值存入返回的 Watch；随后在后台协程中持续监听，
// 监听通道因压缩或中断断开后会按最新修订版本自动重建监听。读取基线失败时返回错误。
func (client *Client) WatchPrefix(ctx context.Context, prefix string) (*Watch, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var w = &Watch{
		revision:     resp.Header.Revision,
		eventChan:    make(chan *clientv3.Event, 100),
		incipientKVs: resp.Kvs,
	}

	goroutine.Go(func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.Client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					logger.Logger.WithFields(logrus.Fields{
						"err":      err,
						"register": ecode.ErrKindRegisterErr,
						"prefix":   prefix,
					}).Error(ecode.MsgWatchRequestErr)
					continue
				}
				for _, ev := range n.Events {
					select {
					case w.eventChan <- ev:
					default:
						// 事件通道缓冲已满时直接丢弃当前事件，避免阻塞后台监听循环。
						logger.Logger.Error("watch etcd with prefix err block event chan, drop event message")
					}
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel = cancel
			if w.revision > 0 {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
			} else {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	})

	return w, nil
}

// Close 取消后台监听协程；重复调用或未初始化 cancel 时安全返回，不返回错误。
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
