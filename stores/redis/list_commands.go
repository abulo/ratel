package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BLPop 阻塞式从给定键列表的左侧弹出元素，返回被操作键与弹出元素；
// timeout 为最长阻塞等待时间，超时仍无元素时返回 redis.Nil。
func (r *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BLPop(getCtx(ctx), timeout, keys...).Result()
		return err
	}, acceptable)
	return
}

// BLMPop 阻塞式按 direction（left/right）从多个列表弹出 count 个元素，
// 返回被操作键与弹出元素；超时仍无元素时返回 redis.Nil。
func (r *Client) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) (key string, val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		key, val, err = conn.BLMPop(getCtx(ctx), timeout, direction, count, keys...).Result()
		return err
	}, acceptable)
	return
}

// BRPop 阻塞式从给定键列表的右侧弹出元素，返回被操作键与弹出元素；
// timeout 为最长阻塞等待时间，超时仍无元素时返回 redis.Nil。
func (r *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BRPop(getCtx(ctx), timeout, keys...).Result()
		return err
	}, acceptable)
	return
}

// BRPopLPush 阻塞式将 source 列表尾部元素弹出并推入 destination 列表头部；
// timeout 为最长阻塞等待时间，超时仍无元素时返回 redis.Nil。
func (r *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BRPopLPush(getCtx(ctx), source, destination, timeout).Result()
		return err
	}, acceptable)
	return
}

// LIndex 返回列表 key 中索引 index 处的元素；索引越界或列表不存在时返回 redis.Nil。
func (r *Client) LIndex(ctx context.Context, key string, index int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LIndex(getCtx(ctx), key, index).Result()
		return err
	}, acceptable)
	return
}

// LInsert 在列表 key 中按 op（BEFORE/AFTER）相对参考元素 pivot 插入 value，
// 返回插入后的列表长度；参考元素不存在时返回 -1。
func (r *Client) LInsert(ctx context.Context, key string, op string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsert(getCtx(ctx), key, op, pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LInsertBefore 在参考元素 pivot 之前插入 value，返回插入后的列表长度。
func (r *Client) LInsertBefore(ctx context.Context, key string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsertBefore(getCtx(ctx), key, pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LInsertAfter 在参考元素 pivot 之后插入 value，返回插入后的列表长度。
func (r *Client) LInsertAfter(ctx context.Context, key string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsertAfter(getCtx(ctx), key, pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LLen 返回列表 key 的长度；键不存在时返回 0。
func (r *Client) LLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LLen(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// LMPop 按 direction（left/right）从多个列表弹出 count 个元素，返回被操作键与弹出元素。
func (r *Client) LMPop(ctx context.Context, direction string, count int64, keys ...string) (key string, val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		key, val, err = conn.LMPop(getCtx(ctx), direction, count, keys...).Result()
		return err
	}, acceptable)
	return
}

// LPop 弹出列表 key 的左侧首元素；列表为空或不存在时返回 redis.Nil。
func (r *Client) LPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPop(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// LPopCount 从列表 key 左侧弹出最多 count 个元素，按弹出顺序返回。
func (r *Client) LPopCount(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPopCount(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// LPos 返回列表中 value 元素首次出现的索引；args 限定扫描范围与数量，未找到时返回 -1。
func (r *Client) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPos(getCtx(ctx), key, value, args).Result()
		return err
	}, acceptable)
	return
}

// LPosCount 返回列表中 value 元素匹配的前 count 个索引位置，扫描范围受 args 约束。
func (r *Client) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPosCount(getCtx(ctx), key, value, count, args).Result()
		return err
	}, acceptable)
	return
}

// LPush 在列表 key 左侧依次插入 values，返回插入后的列表长度；键不存在时先创建空列表。
func (r *Client) LPush(ctx context.Context, key string, values ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPush(getCtx(ctx), key, values...).Result()
		return err
	}, acceptable)
	return
}

// LPushX 仅当列表 key 已存在时在左侧插入 value，返回列表长度；键不存在时操作无效。
func (r *Client) LPushX(ctx context.Context, key string, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPushX(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// LRange 返回列表 key 在 [start, stop] 范围内的元素，索引支持负数定位尾部。
func (r *Client) LRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LRange(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// LRem 移除列表中值为 value 的元素，count 控制移除数量与扫描方向，返回实际移除数量。
func (r *Client) LRem(ctx context.Context, key string, count int64, value any) (val int64, err error) {

	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LRem(getCtx(ctx), key, count, value).Result()
		return err
	}, acceptable)
	return
}

// LSet 将列表 key 中索引 index 处的元素设置为 value；索引越界时命令报错。
func (r *Client) LSet(ctx context.Context, key string, index int64, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LSet(getCtx(ctx), key, index, value).Result()
		return err
	}, acceptable)
	return
}

// LTrim 修剪列表，仅保留 [start, stop] 范围内的元素，范围外元素被删除。
func (r *Client) LTrim(ctx context.Context, key string, start, stop int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LTrim(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// RPop 弹出列表 key 的右侧尾元素；列表为空或不存在时返回 redis.Nil。
func (r *Client) RPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPop(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// RPopCount 从列表 key 右侧弹出最多 count 个元素，按弹出顺序返回。
func (r *Client) RPopCount(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPopCount(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// RPopLPush 弹出 source 列表尾部元素并推入 destination 列表头部；源列表与目标列表可为同一列表。
func (r *Client) RPopLPush(ctx context.Context, source, destination string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPopLPush(getCtx(ctx), source, destination).Result()
		return err
	}, acceptable)
	return
}

// RPush 在列表 key 右侧依次插入 values，返回插入后的列表长度；键不存在时先创建空列表。
func (r *Client) RPush(ctx context.Context, key string, values ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPush(getCtx(ctx), key, values...).Result()
		return err
	}, acceptable)
	return
}

// RPushX 仅当列表 key 已存在时在右侧插入 value，返回列表长度；键不存在时操作无效。
func (r *Client) RPushX(ctx context.Context, key string, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPushX(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// LMove 将 source 列表由 srcpos 指定端的元素移动到 destination 列表的 destpos 端，
// 返回被移动的元素；源与目标可为同一列表。
func (r *Client) LMove(ctx context.Context, source, destination, srcpos, destpos string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LMove(getCtx(ctx), source, destination, srcpos, destpos).Result()
		return err
	}, acceptable)
	return
}

// BLMove 阻塞式执行 LMove，ts 为最长阻塞等待时间；超时仍无元素时返回 redis.Nil。
func (r *Client) BLMove(ctx context.Context, source, destination, srcpos, destpos string, ts time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BLMove(getCtx(ctx), source, destination, srcpos, destpos, ts).Result()
		return err
	}, acceptable)
	return
}
