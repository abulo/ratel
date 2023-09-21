package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func (r *Client) HDel(ctx context.Context, key any, fields ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HDel(getCtx(ctx), r.k(key), fields...).Result()
		return err
	}, acceptable)
	return
}

// HExists 查看哈希表 key 中，给定域 field 是否存在。
func (r *Client) HExists(ctx context.Context, key any, field string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExists(getCtx(ctx), r.k(key), field).Result()
		return err
	}, acceptable)
	return
}

// HGet 返回哈希表 key 中给定域 field 的值。
func (r *Client) HGet(ctx context.Context, key any, field string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGet(getCtx(ctx), r.k(key), field).Result()
		return err
	}, acceptable)
	return
}

// HGetAll 返回哈希表 key 中，所有的域和值。
// 在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
func (r *Client) HGetAll(ctx context.Context, key any) (val map[string]string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGetAll(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HIncrBy 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (r *Client) HIncrBy(ctx context.Context, key any, field string, incr int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrBy(getCtx(ctx), r.k(key), field, incr).Result()
		return err
	}, acceptable)
	return
}

// HIncrByFloat 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Client) HIncrByFloat(ctx context.Context, key any, field string, incr float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrByFloat(getCtx(ctx), r.k(key), field, incr).Result()
		return err
	}, acceptable)
	return
}

// HKeys 返回哈希表 key 中的所有域。
func (r *Client) HKeys(ctx context.Context, key any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HKeys(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HLen 返回哈希表 key 中域的数量。
func (r *Client) HLen(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HLen(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HMGet 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (r *Client) HMGet(ctx context.Context, key any, fields ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMGet(getCtx(ctx), r.k(key), fields...).Result()
		return err
	}, acceptable)
	return
}

// HSet 将哈希表 key 中的域 field 的值设为 value 。
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Client) HSet(ctx context.Context, key any, value ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSet(getCtx(ctx), r.k(key), value...).Result()
		return err
	}, acceptable)
	return
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (r *Client) HMSet(ctx context.Context, key any, value ...any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMSet(getCtx(ctx), r.k(key), value...).Result()
		return err
	}, acceptable)
	return
}

// HSetNX 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
// 若域 field 已经存在，该操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func (r *Client) HSetNX(ctx context.Context, key any, field string, value any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSetNX(getCtx(ctx), r.k(key), field, value).Result()
		return err
	}, acceptable)
	return
}

// HScan 详细信息请参考 SCAN 命令。
func (r *Client) HScan(ctx context.Context, key any, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.HScan(getCtx(ctx), r.k(key), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// HVals 返回哈希表 key 中所有域的值。
func (r *Client) HVals(ctx context.Context, key any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HVals(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) HRandField(ctx context.Context, key any, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandField(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) HRandFieldWithValues(ctx context.Context, key any, count int) (val []redis.KeyValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandFieldWithValues(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
