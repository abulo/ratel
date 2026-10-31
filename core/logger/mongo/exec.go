package mongo

import (
	"context"
	"sync"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// ExecCloser 抽象日志条目的后端写入器，由各实现负责把条目提交到对应存储。
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

type defaultExec struct {
	client     *mongodb.MongoDB
	collection string
	canClose   bool

	// PERF-3-4 [P1]：首次写入时惰性获取并缓存集合句柄，避免每条日志都走 NewCollection（内部 Ping 一次 RTT）。
	// DecoratedCollection 的链式条件字段由 reset() 在写操作后就地改写，不保证多 goroutine 并发安全，
	// 因此用 handlerMu 保护共享句柄上的 InsertOne；两后台 worker 的写入在此串行化，
	// 锁内只做一次 InsertOne 网络操作，不再包含 Ping；锁粒度与并发度的实际收益需实测验证。
	handlerMu sync.Mutex
	handler   *mongodb.DecoratedCollection
}

// NewExec 基于给定的 MongoDB 客户端与集合名创建写入器。
func NewExec(client *mongodb.MongoDB, collection string) ExecCloser {
	return &defaultExec{
		client:     client,
		collection: collection,
		canClose:   true,
	}
}

// NewExecWithURL 基于给定的 MongoDB 客户端与集合名创建写入器，当前实现与 NewExec 行为一致。
func NewExecWithURL(client *mongodb.MongoDB, collection string) ExecCloser {
	return &defaultExec{
		client:     client,
		collection: collection,
		canClose:   true,
	}
}

func (e *defaultExec) Exec(entry *entry.Entry) error {
	item := make(bson.M)
	item["host"] = entry.Host
	item["timestamp"] = entry.Timestamp
	item["file"] = entry.File
	item["func"] = entry.Func
	item["message"] = entry.Message
	item["level"] = entry.Level
	data := bson.M(entry.Data)
	item["data"] = data
	ctx := context.Background()

	// PERF-3-4 [P1]：句柄在首个日志时创建并缓存；若 NewCollection 失败则保持 handler=nil，
	// 后续日志重试获取，与原实现“每条都重新创建、失败即返回错误”的错误行为保持一致。
	e.handlerMu.Lock()
	defer e.handlerMu.Unlock()
	if e.handler == nil {
		handler, err := e.client.NewCollection(e.collection)
		if err != nil {
			return err
		}
		e.handler = handler
	}
	_, err := e.handler.InsertOne(ctx, item)
	return err
}
