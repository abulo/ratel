package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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

// ClusterSlots 获取集群节点的映射数组
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

// ClusterNodes Get Cluster config for the node
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

// ClusterMeet Force a node cluster to handshake with another node
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

// ClusterForget Remove a node from the nodes table
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

// ClusterReplicate Reconfigure a node as a replica of the specified master node
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

// ClusterResetSoft Reset a Redis Cluster node
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

// ClusterResetHard Reset a Redis Cluster node
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

// ClusterInfo Provides info about Redis Cluster node state
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

// ClusterKeySlot Returns the hash slot of the specified key
func (r *Client) ClusterKeySlot(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterKeySlot(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ClusterGetKeysInSlot Return local key names in the specified hash slot
func (r *Client) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) (val []string, err error) {
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

// ClusterCountFailureReports Return the number of failure reports active for a given node
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

// ClusterCountKeysInSlot Return the number of local keys in the specified hash slot
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

// ClusterDelSlots Set hash slots as unbound in receiving node
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

// ClusterDelSlotsRange ->  ClusterDelSlots
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

// ClusterSaveConfig Forces the node to save cluster state on disk
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

// ClusterSlaves List replica nodes of the specified master node
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

// ClusterFailover Forces a replica to perform a manual failover of its master.
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

// ClusterAddSlots Assign new hash slots to receiving node
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

// ClusterAddSlotsRange -> ClusterAddSlots
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
