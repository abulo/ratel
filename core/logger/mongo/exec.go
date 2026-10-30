package mongo

import (
	"context"

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
	handler, err := e.client.NewCollection(e.collection)
	if err != nil {
		return err
	}
	_, err = handler.InsertOne(ctx, item)
	return err
}
