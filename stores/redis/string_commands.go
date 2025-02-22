package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Append 将字符串追加到key的值末尾 key: 键名 value: 要追加的字符串
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

// Decr 将key的值减1 key: 键名
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

// DecrBy 将key的值减去指定数值 key: 键名 value: 要减去的数值
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

// Get 获取key的值 key: 键名
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

// GetRange 获取key的值的子字符串 key: 键名 start: 起始位置 end: 结束位置
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

// GetSet 设置key的值并返回旧值 key: 键名 value: 新值
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

// GetEx 获取key的值并设置过期时间 key: 键名 ts: 过期时间
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

// GetDel 获取key的值并删除 key: 键名
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

// Incr 将key的值加1 key: 键名
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

// IncrBy 将key的值增加指定数值 key: 键名 value: 要增加的数值
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

// IncrByFloat 将key的值增加指定浮点数 key: 键名 value: 要增加的浮点数
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

// LCS 计算两个字符串的最长公共子序列 q: 查询参数
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

// MGet 获取多个key的值 keys: 键名列表
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

// MSet 设置多个key-value对 values: 键值对列表
func (r *Client) MSet(ctx context.Context, values ...any) (val string, err error) {
	// return getRedis(r).MSet(getCtx(ctx), values...)
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

// MSetNX 当所有key都不存在时设置多个key-value对 values: 键值对列表
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

// Set 设置key的值 key: 键名 value: 值 expiration: 过期时间
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

// SetArgs 设置key的值并指定额外参数 key: 键名 value: 值 a: 额外参数
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

// SetEx 设置key的值并指定过期时间 key: 键名 value: 值 expiration: 过期时间
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

// SetNX 当key不存在时设置值 key: 键名 value: 值 expiration: 过期时间
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

// SetXX 当key存在时设置值 key: 键名 value: 值 expiration: 过期时间
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

// SetRange 从指定位置覆写key的值 key: 键名 offset: 偏移量 value: 新值
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

// StrLen 获取key的值的长度 key: 键名
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
