package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Append 将 value 追加到 key 字符串值的末尾，返回追加后的字符串长度；
// key 不存在时按空串新建。
func (r *Client) Append(ctx context.Context, key string, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Append(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// Decr 将 key 的数值减 1，返回减后的值；key 不存在时先置 0 再减。
func (r *Client) Decr(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Decr(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// DecrBy 将 key 的数值减去给定 value，返回结果；value 为负数时等价于加法。
func (r *Client) DecrBy(ctx context.Context, key string, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DecrBy(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// Get 返回 key 对应的字符串值；键不存在时返回 redis.Nil 错误。
func (r *Client) Get(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Get(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// GetRange 返回 key 字符串值在 [start, end] 区间内的子串，索引支持负数定位尾部。
func (r *Client) GetRange(ctx context.Context, key string, start, end int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetRange(getCtx(ctx), key, start, end).Result()
		return err
	}, acceptable)
	return
}

// GetSet 将 key 设置为 value 并返回设置前的旧值；键不存在时旧值返回 redis.Nil。
func (r *Client) GetSet(ctx context.Context, key string, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetSet(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// GetEx 读取 key 值的同时为其设置过期时间，返回读取到的值。
func (r *Client) GetEx(ctx context.Context, key string, ts time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetEx(getCtx(ctx), key, ts).Result()
		return err
	}, acceptable)
	return
}

// GetDel 读取 key 值并在返回后删除该键；键不存在时返回 redis.Nil。
func (r *Client) GetDel(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetDel(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Incr 将 key 的数值加 1，返回加后的值；key 不存在时先置 0 再加。
func (r *Client) Incr(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Incr(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// IncrBy 将 key 的数值增加给定 value，返回结果；value 为负数时等价于减法。
func (r *Client) IncrBy(ctx context.Context, key string, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrBy(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// IncrByFloat 将 key 的数值增加给定浮点数增量，返回更新后的结果。
func (r *Client) IncrByFloat(ctx context.Context, key string, value float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrByFloat(getCtx(ctx), key, value).Result()
		return err
	}, acceptable)
	return
}

// LCS 按查询条件计算两个字符串的最长公共子序列（LCS）。
func (r *Client) LCS(ctx context.Context, q *redis.LCSQuery) (val *redis.LCSMatch, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LCS(getCtx(ctx), q).Result()
		return err
	}, acceptable)
	return
}

// MGet 批量返回多个键的值，返回切片顺序与传入键顺序一致；不存在的键对应位置为 nil。
func (r *Client) MGet(ctx context.Context, keys ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MGet(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// MSet 批量设置给定的键值对，参数按 key、value 成对传入。
func (r *Client) MSet(ctx context.Context, values ...any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MSet(getCtx(ctx), values...).Result()
		return err
	}, acceptable)
	return

}

// MSetNX 仅当所有给定键都不存在时批量设置键值对；任一键已存在则整个操作不执行，返回 false。
func (r *Client) MSetNX(ctx context.Context, values ...any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MSetNX(getCtx(ctx), values...).Result()
		return err
	}, acceptable)
	return
}

// Set 设置 key 的值为 value；expiration 为 0 表示不设置过期时间。
func (r *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Set(getCtx(ctx), key, value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetArgs 按给定 SetArgs（含 NX/XX 条件、过期方式等）设置 key 的值。
func (r *Client) SetArgs(ctx context.Context, key string, value any, a redis.SetArgs) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetArgs(getCtx(ctx), key, value, a).Result()
		return err
	}, acceptable)
	return
}

// SetEx 设置 key 的值并指定过期时间，过期时间必须大于 0。
func (r *Client) SetEx(ctx context.Context, key string, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetEx(getCtx(ctx), key, value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetNX 仅当 key 不存在时设置值，成功返回 true；键已存在时不修改并返回 false。
func (r *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetNX(getCtx(ctx), key, value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetXX 仅当 key 已存在时设置值；键不存在时不创建并返回 false。
func (r *Client) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetXX(getCtx(ctx), key, value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetRange 从偏移量 offset 开始覆写 key 字符串的内容，返回修改后的字符串长度。
func (r *Client) SetRange(ctx context.Context, key string, offset int64, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetRange(getCtx(ctx), key, offset, value).Result()
		return err
	}, acceptable)
	return
}

// StrLen 返回 key 字符串值的字节长度；键不存在时返回 0。
func (r *Client) StrLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.StrLen(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}
