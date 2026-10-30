// Package redis 提供基于 go-redis 的 Redis 客户端封装与命令操作集合。
//
// 该包通过 NewRedisClient 配合 Option 构造 Client 配置，支持单节点、集群、
// 哨兵与分片四种连接模式；Client 的各方法按数据类型组织为命令分组，统一经过
// 熔断器执行，并可按需挂载链路追踪与指标采集钩子。nil 或已取消的上下文会被
// 替换为 context.TODO()，键不存在等可接受错误由 acceptable 判定后不作为失败上报。
package redis
