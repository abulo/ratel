package mongo

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// ExecCloser 将logrus条目写入数据库并关闭数据库
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

type defaultExec struct {
	client     *mongodb.MongoDB
	collection string
	canClose   bool
}

// NewExec create an exec instance
func NewExec(client *mongodb.MongoDB, collection string) ExecCloser {
	return &defaultExec{
		client:     client,
		collection: collection,
		canClose:   true,
	}
}

// NewExecWithURL create an exec instance
func NewExecWithURL(client *mongodb.MongoDB, collection string) ExecCloser {
	return &defaultExec{
		client:     client,
		collection: collection,
		canClose:   true,
	}
}

// Exec ...
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
