// Package session 提供基于 Redis 的会话存储，支持以点分键路径读写嵌套会话数据。
//
// Session 将整份会话序列化为 JSON 存放在单个 Redis 键下；Put、Get、Remove 按点分路径
// 操作嵌套字段，Destroy 删除整份会话，TTL 以秒为单位控制会话过期时间。
package session
