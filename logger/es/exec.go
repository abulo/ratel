package es

import (
	"context"

	"github.com/abulo/ratel/v3/logger/entry"
	"github.com/abulo/ratel/v3/stores/elasticsearch"
	"github.com/google/uuid"
)

// ExecCloser 将logrus条目写入数据库并关闭数据库
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

type defaultExec struct {
	client   *elasticsearch.Client
	canClose bool
}

// NewExec create an exec instance
func NewExec(client *elasticsearch.Client) ExecCloser {
	return &defaultExec{
		client:   client,
		canClose: true,
	}
}

// NewExecWithURL create an exec instance
func NewExecWithURL(client *elasticsearch.Client) ExecCloser {
	return &defaultExec{
		client:   client,
		canClose: true,
	}
}

// Exec ...
func (e *defaultExec) Exec(entry *entry.Entry) error {
	ctx := context.Background()
	_, err := e.client.Index().Index("logger_entry").Id(uuid.New().String()).BodyJson(entry).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
