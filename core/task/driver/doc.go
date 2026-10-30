// Package driver 提供分布式任务节点注册与成员发现的驱动抽象，内置 Redis 普通键、
// Redis 有序集合与 etcd 租约三种实现。
//
// Task 通过 DriverV2 接口完成节点注册、心跳保活与在线节点列表获取，
// 不同中间件实现遵循同一套启动、停止与查询契约。
package driver
