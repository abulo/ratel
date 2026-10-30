// Package resolver 提供基于 etcd 服务发现的 gRPC 解析器 Builder。
//
// 解析器按目标服务名持续监听注册中心的节点列表，以及路由、提供方与消费方配置；节点集合变更时
// 将最新地址与属性写入 resolver.State 上报给负载均衡器，Close 后停止后台监听。
package resolver
