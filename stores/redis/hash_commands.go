package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// HDel 删除哈希表key中的一个或多个指定域，不存在的域将被忽略 ctx: 上下文 key: 哈希表键 fields: 要删除的域
func (r *Client) HDel(ctx context.Context, key string, fields ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HDel(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HExists 查看哈希表key中给定域是否存在 ctx: 上下文 key: 哈希表键 field: 要检查的域
func (r *Client) HExists(ctx context.Context, key, field string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExists(getCtx(ctx), key, field).Result()
		return err
	}, acceptable)
	return
}

// HGet 返回哈希表key中给定域的值 ctx: 上下文 key: 哈希表键 field: 要获取值的域
func (r *Client) HGet(ctx context.Context, key, field string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGet(getCtx(ctx), key, field).Result()
		return err
	}, acceptable)
	return
}

// HGetAll 返回哈希表key中所有的域和值 ctx: 上下文 key: 哈希表键
// 在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
func (r *Client) HGetAll(ctx context.Context, key string) (val map[string]string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGetAll(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// HIncrBy 为哈希表key中的域的值加上增量 ctx: 上下文 key: 哈希表键 field: 要增加值的域 incr: 增量值
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (r *Client) HIncrBy(ctx context.Context, key string, field string, incr int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrBy(getCtx(ctx), key, field, incr).Result()
		return err
	}, acceptable)
	return
}

// HIncrByFloat 为哈希表key中的域加上浮点数增量 ctx: 上下文 key: 哈希表键 field: 要增加值的域 incr: 浮点数增量
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Client) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrByFloat(getCtx(ctx), key, field, incr).Result()
		return err
	}, acceptable)
	return
}

// HKeys 返回哈希表key中的所有域 ctx: 上下文 key: 哈希表键
func (r *Client) HKeys(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HKeys(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// HLen 返回哈希表key中域的数量 ctx: 上下文 key: 哈希表键
func (r *Client) HLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HLen(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// HMGet 返回哈希表key中一个或多个给定域的值 ctx: 上下文 key: 哈希表键 fields: 要获取值的域
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMGet(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HSet 将哈希表key中的域的值设为指定值 ctx: 上下文 key: 哈希表键 value: 要设置的值
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Client) HSet(ctx context.Context, key string, value ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSet(getCtx(ctx), key, value...).Result()
		return err
	}, acceptable)
	return
}

// HMSet 同时将多个field-value对设置到哈希表key中 ctx: 上下文 key: 哈希表键 value: 要设置的键值对
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (r *Client) HMSet(ctx context.Context, key string, value ...any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMSet(getCtx(ctx), key, value...).Result()
		return err
	}, acceptable)
	return
}

// HSetNX 当域不存在时将哈希表key中的域的值设置为指定值 ctx: 上下文 key: 哈希表键 field: 要设置的域 value: 要设置的值
// 若域 field 已经存在，该操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func (r *Client) HSetNX(ctx context.Context, key string, field string, value any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSetNX(getCtx(ctx), key, field, value).Result()
		return err
	}, acceptable)
	return
}

// HScan 迭代哈希表中的键值对 ctx: 上下文 key: 哈希表键 cursorIn: 游标 match: 匹配模式 count: 每次迭代返回的数量
func (r *Client) HScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.HScan(getCtx(ctx), key, cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// HScanNoValues 迭代哈希表中的键（不返回值） ctx: 上下文 key: 哈希表键 cursor: 游标 match: 匹配模式 count: 每次迭代返回的数量
func (r *Client) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) (val []string, cursorOut uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursorOut, err = conn.HScanNoValues(getCtx(ctx), key, cursor, match, count).Result()
		return err
	}, acceptable)
	return
}

// HVals 返回哈希表key中所有域的值 ctx: 上下文 key: 哈希表键
func (r *Client) HVals(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HVals(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// HRandField 从哈希表key中随机获取指定数量的域 ctx: 上下文 key: 哈希表键 count: 要获取的域数量
func (r *Client) HRandField(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandField(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// HRandFieldWithValues 从哈希表key中随机获取指定数量的域及其值 ctx: 上下文 key: 哈希表键 count: 要获取的域数量
func (r *Client) HRandFieldWithValues(ctx context.Context, key string, count int) (val []redis.KeyValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandFieldWithValues(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// HExpire 设置哈希表key中一个或多个域的过期时间 ctx: 上下文 key: 哈希表键 expiration: 过期时间 fields: 要设置过期时间的域
// 如果指定的域不存在，那么命令将被忽略。
func (r *Client) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExpire(getCtx(ctx), key, expiration, fields...).Result()
		return err
	}, acceptable)
	return
}

// HExpireWithArgs 使用额外参数设置哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 expiration: 过期时间 expirationArgs: 额外参数 fields: 要设置过期时间的域
func (r *Client) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExpireWithArgs(getCtx(ctx), key, expiration, expirationArgs, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPExpire 以毫秒为单位设置哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 expiration: 过期时间 fields: 要设置过期时间的域
func (r *Client) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPExpire(getCtx(ctx), key, expiration, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPExpireWithArgs 使用额外参数以毫秒为单位设置哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 expiration: 过期时间 expirationArgs: 额外参数 fields: 要设置过期时间的域
func (r *Client) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs redis.HExpireArgs, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPExpireWithArgs(getCtx(ctx), key, expiration, expirationArgs, fields...).Result()
		return err
	}, acceptable)
	return
}

// HExpireAt 设置哈希表key中域的过期时间戳 ctx: 上下文 key: 哈希表键 tm: 过期时间戳 fields: 要设置过期时间的域
func (r *Client) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExpireAt(getCtx(ctx), key, tm, fields...).Result()
		return err
	}, acceptable)
	return
}

// HExpireAtWithArgs 使用额外参数设置哈希表key中域的过期时间戳 ctx: 上下文 key: 哈希表键 tm: 过期时间戳 expirationArgs: 额外参数 fields: 要设置过期时间的域
func (r *Client) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExpireAtWithArgs(getCtx(ctx), key, tm, expirationArgs, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPExpireAt 以毫秒为单位设置哈希表key中域的过期时间戳 ctx: 上下文 key: 哈希表键 tm: 过期时间戳 fields: 要设置过期时间的域
func (r *Client) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPExpireAt(getCtx(ctx), key, tm, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPExpireAtWithArgs 使用额外参数以毫秒为单位设置哈希表key中域的过期时间戳 ctx: 上下文 key: 哈希表键 tm: 过期时间戳 expirationArgs: 额外参数 fields: 要设置过期时间的域
func (r *Client) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs redis.HExpireArgs, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPExpireAtWithArgs(getCtx(ctx), key, tm, expirationArgs, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPersist 移除哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 fields: 要移除过期时间的域
func (r *Client) HPersist(ctx context.Context, key string, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPersist(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HExpireTime 返回哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 fields: 要查询的域
func (r *Client) HExpireTime(ctx context.Context, key string, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExpireTime(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPExpireTime 以毫秒为单位返回哈希表key中域的过期时间 ctx: 上下文 key: 哈希表键 fields: 要查询的域
func (r *Client) HPExpireTime(ctx context.Context, key string, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPExpireTime(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HTTL 返回哈希表key中域的剩余生存时间 ctx: 上下文 key: 哈希表键 fields: 要查询的域
func (r *Client) HTTL(ctx context.Context, key string, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HTTL(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}

// HPTTL 以毫秒为单位返回哈希表key中域的剩余生存时间 ctx: 上下文 key: 哈希表键 fields: 要查询的域
func (r *Client) HPTTL(ctx context.Context, key string, fields ...string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HPTTL(getCtx(ctx), key, fields...).Result()
		return err
	}, acceptable)
	return
}
