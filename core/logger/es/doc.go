// Package es 提供基于 Elasticsearch 的 logrus 日志钩子实现。
//
// 该包定义了日志索引的映射与初始化流程，并通过 logrus Hook 将日志条目
// 异步写入 Elasticsearch；索引创建失败时构造钩子会直接 panic。
package es

import (
	"context"

	"github.com/abulo/ratel/v3/stores/elasticsearch"
	"github.com/pkg/errors"
)

// IndexString 返回日志索引 logger_entry 的映射（mapping）定义 JSON 字符串。
//
// 映射包含 file、func、message、level 四个 keyword 字段，data 为 object
// 字段，host 为 keyword 字段，timestamp 为 date 字段，供 CreateIndex 建索引使用。
func IndexString() string {
	body := `{
		"settings": {
			"index": {
				"number_of_shards": 1,
				"number_of_replicas": 0
			}
		},
		"mappings": {
			"dynamic": false,
			"properties": {
				"file": {
					"type": "keyword"
				},
				"func": {
					"type": "keyword"
				},
				"message": {
					"type": "keyword"
				},
				"level": {
					"type": "keyword"
				},
				"data": {
					"type": "object"
				},
				"host": {
					"type": "keyword"
				},
				"timestamp": {
					"type": "date"
				}
			}
		}
	}`
	return body
}

// CreateIndex 确保日志索引 logger_entry 已存在，不存在时按 IndexString 的映射创建。
//
// 创建请求未获得服务端确认时返回错误，索引已存在时直接返回 nil。
func CreateIndex(client *elasticsearch.Client) error {
	ctx := context.Background()
	exists, err := client.IndexExists("logger_entry").Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := client.CreateIndex("logger_entry").BodyJson(IndexString()).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Not acknowledged")
		}
	}
	return nil
}
