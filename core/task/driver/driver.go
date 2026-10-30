package driver

import (
	"context"

	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/redis/go-redis/v9"
)

// DriverV2 抽象分布式任务的节点注册与成员发现驱动，一个 Task 对应一个驱动实例。
//
// 实现约束：Start 与 Stop 必须支持多次调用；GetNodes 在超时时必须返回错误。
type DriverV2 interface {
	// Init 使用服务名与驱动选项完成初始化，须在 Start 之前调用。
	Init(serviceName string, opts ...Option)
	// NodeID 返回当前节点标识。
	NodeID() string
	// GetNodes 返回当前在线节点标识列表；超时或查询失败时返回错误。
	GetNodes(ctx context.Context) (nodes []string, err error)

	// Start 向远程存储（etcd、Redis 等）注册当前节点，并启动协程维持心跳连接，随后返回。
	Start(ctx context.Context) (err error)

	// Stop 停止维持心跳的协程并注销当前节点。
	Stop(ctx context.Context) (err error)

	// WithOption 应用单个驱动配置选项；无法识别的选项类型静默忽略。
	WithOption(opt Option) (err error)
}

// NewRedisDriver 基于给定 Redis 客户端创建普通键值驱动实现。
func NewRedisDriver(redisClient redis.UniversalClient) DriverV2 {
	return newRedisDriver(redisClient)
}

// NewEtcdDriver 基于给定 etcd 客户端创建租约驱动实现。
func NewEtcdDriver(etcdCli *etcdv3.Client) DriverV2 {
	return newEtcdDriver(etcdCli)
}

// NewRedisZSetDriver 基于给定 Redis 客户端创建有序集合驱动实现，按时间戳区分在线节点。
func NewRedisZSetDriver(redisClient redis.UniversalClient) DriverV2 {
	return newRedisZSetDriver(redisClient)
}
