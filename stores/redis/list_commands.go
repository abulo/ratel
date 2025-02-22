package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BLPop 阻塞式弹出列表元素 ctx: 上下文, timeout: 超时时间, keys: 要弹出的键列表
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

// BLMPop 阻塞式弹出多个列表元素 ctx: 上下文, timeout: 超时时间, direction: 弹出方向(left/right), count: 弹出数量, keys: 要弹出的键列表
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

// BRPop 阻塞式弹出列表尾部元素 ctx: 上下文, timeout: 超时时间, keys: 要弹出的键列表
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

// BRPopLPush 阻塞式弹出并插入元素 ctx: 上下文, source: 源列表, destination: 目标列表, timeout: 超时时间
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

// LIndex 获取列表指定位置的元素 ctx: 上下文, key: 列表键, index: 元素索引
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

// LInsert 在列表中插入元素 ctx: 上下文, key: 列表键, op: 插入位置(BEFORE/AFTER), pivot: 参考元素, value: 要插入的值
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

// LInsertBefore 在参考元素前插入值 ctx: 上下文, key: 列表键, pivot: 参考元素, value: 要插入的值
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

// LInsertAfter 在参考元素后插入值 ctx: 上下文, key: 列表键, pivot: 参考元素, value: 要插入的值
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

// LLen 获取列表长度 ctx: 上下文, key: 列表键
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

// LMPop 弹出多个列表元素 ctx: 上下文, direction: 弹出方向(left/right), count: 弹出数量, keys: 要弹出的键列表
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

// LPop 弹出列表头元素 ctx: 上下文, key: 列表键
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

// LPopCount 弹出多个列表头元素 ctx: 上下文, key: 列表键, count: 弹出数量
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

// LPush 在列表头部插入元素 ctx: 上下文, key: 列表键, values: 要插入的值
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

// LPushX 仅在列表存在时在头部插入元素 ctx: 上下文, key: 列表键, value: 要插入的值
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

// LRange 获取列表指定范围的元素 ctx: 上下文, key: 列表键, start: 起始索引, stop: 结束索引
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

// LRem 移除列表中指定值的元素 ctx: 上下文, key: 列表键, count: 移除数量, value: 要移除的值
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

// LSet 设置列表指定位置的元素 ctx: 上下文, key: 列表键, index: 元素索引, value: 要设置的值
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

// LTrim 修剪列表只保留指定范围的元素 ctx: 上下文, key: 列表键, start: 起始索引, stop: 结束索引
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

// RPop 弹出列表尾部元素 ctx: 上下文, key: 列表键
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

// RPopCount 弹出多个列表尾部元素 ctx: 上下文, key: 列表键, count: 弹出数量
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

// RPopLPush 弹出尾部元素并插入到另一个列表 ctx: 上下文, source: 源列表, destination: 目标列表
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

// RPush 在列表尾部插入元素 ctx: 上下文, key: 列表键, values: 要插入的值
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

// RPushX 仅在列表存在时在尾部插入元素 ctx: 上下文, key: 列表键, value: 要插入的值
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
