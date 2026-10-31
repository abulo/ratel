package es

import (
	"context"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/elasticsearch"
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
	// PERF-3-3 [P1]：去掉每条日志自行生成的 UUID 文档 ID，由 ES 自动生成 _id，省掉 uuid.New() 的系统调用与随机数开销。
	// 批量聚合（BulkIndexer / 定时+条数双阈值 flush）涉及队列与失败重试语义调整，风险较高，留待二期结合线上日志量级实测后实施。
	_, err := e.client.Index().Index(e.index).BodyJson(entry).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
