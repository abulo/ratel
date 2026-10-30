package driver

import (
	"time"

	"github.com/google/uuid"
)

// GlobalKeyPrefix 为分布式任务在 Redis 与 etcd 中使用的全局键前缀。
const GlobalKeyPrefix = "Distributed-Cron:"

// GetKeyPre 返回指定服务名对应的存储键前缀。
func GetKeyPre(serviceName string) string {
	return GlobalKeyPrefix + serviceName + ":"
}

// GetNodeId 生成全局唯一的节点标识，由服务名前缀与 UUID 组成。
func GetNodeId(serviceName string) string {
	return GetKeyPre(serviceName) + uuid.New().String()
}

// TimePre 返回时间 t 向前推移 preDuration 后的 Unix 时间戳，用于过滤已过期节点。
func TimePre(t time.Time, preDuration time.Duration) int64 {
	return t.Add(-preDuration).Unix()
}
