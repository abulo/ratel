// Package elasticsearch 提供基于 olivere/elastic v7 的 Elasticsearch 客户端封装。
//
// 该包通过 Config 配置节点地址、基础认证与功能开关，NewClient 按节点地址缓存复用客户端；
// ESTracedTransport 在 HTTP 传输层挂载链路追踪与指标采集，MaxContentLength 限制参与
// 追踪记录的请求体大小。
package elasticsearch
