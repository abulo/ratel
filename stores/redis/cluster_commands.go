package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ClusterMyShardID 返回当前节点所在分片的 ID。
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

// ClusterMyID 返回当前节点在集群中的节点 ID。
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

// ClusterSlots 返回集群哈希槽到服务节点的映射信息。
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

// ClusterShards 返回集群各分片及其节点与槽位信息。
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

// ClusterLinks 返回当前节点与集群其他节点之间的连接链路信息。
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

// ClusterNodes 返回集群全部节点的原始状态信息字符串。
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

// ClusterMeet 引导 host:port 节点加入当前集群。
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

// ClusterForget 将 nodeID 指定的节点从集群中移除。
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

// ClusterReplicate 将当前节点配置为 nodeID 主节点的从节点。
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

// ClusterResetSoft 软重置当前节点的集群状态，保留本节点已有数据。
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

// ClusterResetHard 硬重置当前节点的集群状态，并清空本节点数据。
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

// ClusterInfo 返回当前节点的集群状态统计信息。
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

// ClusterKeySlot 返回 key 经哈希计算得到的哈希槽编号。
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

// ClusterGetKeysInSlot 返回 slot 槽中最多 count 个键的列表。
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

// ClusterCountFailureReports 返回 nodeID 节点当前尚未清除的故障报告数量。
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

// ClusterCountKeysInSlot 返回 slot 槽中存储的键数量。
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

// ClusterDelSlots 解除当前节点对给定 slots 哈希槽的绑定。
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

// ClusterDelSlotsRange 解除当前节点对 [min, max] 区间内哈希槽的绑定。
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

// ClusterSaveConfig 强制将当前节点的集群配置持久化到磁盘。
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

// ClusterSlaves 列出 nodeID 主节点的全部从节点信息。
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

// ClusterFailover 在当前从节点上手动触发故障转移，将其提升为主节点。
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

// ClusterAddSlots 将给定 slots 哈希槽分配给当前节点。
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

// ClusterAddSlotsRange 将 [min, max] 区间内的哈希槽分配给当前节点。
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

// ReadOnly 将当前连接切换为只读模式，允许从节点响应读命令。
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

// ReadWrite 将当前连接从只读模式恢复为读写模式。
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
