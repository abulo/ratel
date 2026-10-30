package es

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/elasticsearch"
	"github.com/google/uuid"
)

// ExecCloser 抽象日志条目的后端写入器，由各实现负责把条目提交到对应存储。
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

type defaultExec struct {
	client   *elasticsearch.Client
	index    string
	canClose bool
}

// NewExec 基于给定的 Elasticsearch 客户端与索引名创建写入器。
func NewExec(client *elasticsearch.Client, index string) ExecCloser {
	return &defaultExec{
		client:   client,
		index:    index,
		canClose: true,
	}
}

// NewExecWithURL 基于给定的 Elasticsearch 客户端与索引名创建写入器，当前实现与 NewExec 行为一致。
func NewExecWithURL(client *elasticsearch.Client, index string) ExecCloser {
	return &defaultExec{
		client:   client,
		index:    index,
		canClose: true,
	}
}

func (e *defaultExec) Exec(entry *entry.Entry) error {
	ctx := context.Background()
	_, err := e.client.Index().Index(e.index).Id(uuid.New().String()).BodyJson(entry).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
