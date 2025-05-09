package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ClusterMyShardID 获取当前节点的ID ctx:上下文
func (r *Client) ClusterMyShardID(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterMyShardID(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterMyID(ctx context.Context) *StringCmd
func (r *Client) ClusterMyID(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterMyID(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterSlots 获取集群节点的映射数组 ctx:上下文
func (r *Client) ClusterSlots(ctx context.Context) (val []redis.ClusterSlot, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSlots(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterShards 获取集群分片信息 ctx:上下文
func (r *Client) ClusterShards(ctx context.Context) (val []redis.ClusterShard, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterShards(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterLinks 获取集群节点间的链接信息 ctx:上下文
func (r *Client) ClusterLinks(ctx context.Context) (val []redis.ClusterLink, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterLinks(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterNodes 获取集群节点信息 ctx:上下文
func (r *Client) ClusterNodes(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterNodes(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterMeet 将新节点加入集群 ctx:上下文 host:主机地址 port:端口号
func (r *Client) ClusterMeet(ctx context.Context, host, port string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterMeet(getCtx(ctx), host, port).Result()
		return err
	}, acceptable)
	return
}

// ClusterForget 从集群中移除指定节点 ctx:上下文 nodeID:节点ID
func (r *Client) ClusterForget(ctx context.Context, nodeID string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterForget(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterReplicate 将当前节点配置为指定主节点的从节点 ctx:上下文 nodeID:主节点ID
func (r *Client) ClusterReplicate(ctx context.Context, nodeID string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterReplicate(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterResetSoft 软重置集群节点（保留数据） ctx:上下文
func (r *Client) ClusterResetSoft(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterResetSoft(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterResetHard 硬重置集群节点（清除数据） ctx:上下文
func (r *Client) ClusterResetHard(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterResetHard(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterInfo 获取集群节点状态信息 ctx:上下文
func (r *Client) ClusterInfo(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterInfo(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterKeySlot 返回指定key的哈希槽 ctx:上下文 key:键名
func (r *Client) ClusterKeySlot(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterKeySlot(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ClusterGetKeysInSlot 返回指定哈希槽中的key列表 ctx:上下文 slot:哈希槽 count:返回key的数量
func (r *Client) ClusterGetKeysInSlot(ctx context.Context, slot, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterGetKeysInSlot(getCtx(ctx), slot, count).Result()
		return err
	}, acceptable)
	return
}

// ClusterCountFailureReports 返回指定节点的故障报告数量 ctx:上下文 nodeID:节点ID
func (r *Client) ClusterCountFailureReports(ctx context.Context, nodeID string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterCountFailureReports(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterCountKeysInSlot 返回指定哈希槽中的key数量 ctx:上下文 slot:哈希槽
func (r *Client) ClusterCountKeysInSlot(ctx context.Context, slot int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterCountKeysInSlot(getCtx(ctx), slot).Result()
		return err
	}, acceptable)
	return
}

// ClusterDelSlots 删除指定哈希槽的绑定 ctx:上下文 slots:哈希槽列表
func (r *Client) ClusterDelSlots(ctx context.Context, slots ...int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterDelSlots(getCtx(ctx), slots...).Result()
		return err
	}, acceptable)
	return
}

// ClusterDelSlotsRange 删除指定范围内的哈希槽绑定 ctx:上下文 min:起始哈希槽 max:结束哈希槽
func (r *Client) ClusterDelSlotsRange(ctx context.Context, min, max int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterDelSlotsRange(getCtx(ctx), min, max).Result()
		return err
	}, acceptable)
	return
}

// ClusterSaveConfig 强制将集群状态保存到磁盘 ctx:上下文
func (r *Client) ClusterSaveConfig(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSaveConfig(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterSlaves 列出指定主节点的从节点列表 ctx:上下文 nodeID:主节点ID
func (r *Client) ClusterSlaves(ctx context.Context, nodeID string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSlaves(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterFailover 强制从节点执行手动故障转移 ctx:上下文
func (r *Client) ClusterFailover(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterFailover(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterAddSlots 为当前节点分配新的哈希槽 ctx:上下文 slots:哈希槽列表
func (r *Client) ClusterAddSlots(ctx context.Context, slots ...int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterAddSlots(getCtx(ctx), slots...).Result()
		return err
	}, acceptable)
	return
}

// ClusterAddSlotsRange 为当前节点分配指定范围内的哈希槽 ctx:上下文 min:起始哈希槽 max:结束哈希槽
func (r *Client) ClusterAddSlotsRange(ctx context.Context, min, max int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterAddSlotsRange(getCtx(ctx), min, max).Result()
		return err
	}, acceptable)
	return
}

// ReadOnly 将当前连接设置为只读模式（针对从节点） ctx:上下文
func (r *Client) ReadOnly(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ReadOnly(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ReadWrite 将当前连接设置为读写模式（针对从节点） ctx:上下文
func (r *Client) ReadWrite(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ReadWrite(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
