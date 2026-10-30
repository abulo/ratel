package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// HDel 删除哈希 key 中给定的一个或多个域，返回实际删除的域数量；不存在的域被忽略。
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

// HExists 判断哈希 key 中给定域 field 是否存在，存在返回 true。
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

// HGet 返回哈希 key 中给定域 field 的值；域不存在时返回 redis.Nil。
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

// HGetAll 返回哈希 key 中全部域与值组成的 map；空哈希或不存在的 key 返回空 map。
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

// HIncrBy 为哈希 key 中域 field 的整数值加上增量 incr，返回更新后的值。
//
// 增量可为负数，等价于减法；key 或域不存在时先将域初始化为 0 再运算；
// 对存储非整数数值的域执行该命令会报错；结果限制在 64 位有符号整数范围内。
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

// HIncrByFloat 为哈希 key 中域 field 的值加上浮点增量 incr，返回更新后的值。
//
// 域或 key 不存在时先将域置为 0，再执行加法。
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

// HKeys 返回哈希 key 中全部域名称。
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

// HLen 返回哈希 key 中域的数量。
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

// HMGet 返回哈希 key 中给定多个域的值，切片顺序与域顺序一致；不存在的域对应位置为 nil。
//
// 不存在的 key 按空哈希处理，返回值只包含 nil 占位项。
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

// HSet 将哈希 key 中域的值设为给定值，参数按域、值成对传入，返回本次新增的域数量。
//
// key 不存在时新建哈希；已存在的域旧值被覆盖。
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

// HMSet 同时将多组域值写入哈希 key，参数按域、值成对传入，成功返回 true。
//
// 已存在的域被覆盖；key 不存在时新建空哈希。
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

// HSetNX 仅当域 field 不存在时将其设为给定值，成功返回 true。
//
// 域已存在时操作无效并返回 false；key 不存在时新建哈希。
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

// HScan 以游标 cursorIn 迭代哈希 key 中的域值对；match 为域匹配模式，count 为建议批次大小，
// 返回本次元素与供下次调用使用的新游标。
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

// HScanNoValues 以游标 cursorIn 迭代哈希 key 中的域名称（不返回值）；match 为匹配模式，
// count 为建议批次大小，返回本次域名称与新游标。
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

// HVals 返回哈希 key 中全部域的值。
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

// HRandField 随机返回哈希 key 中 count 个域名称。
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

// HRandFieldWithValues 随机返回哈希 key 中 count 个域及其对应的值。
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

// HExpire 为哈希 key 中给定域设置秒级过期时间，返回各域的状态码；不存在的域被忽略。
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

// HExpireWithArgs 按给定 HExpireArgs 为哈希 key 中给定域设置秒级过期时间。
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

// HPExpire 以毫秒为单位为哈希 key 中给定域设置过期时间。
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

// HPExpireWithArgs 按给定 HExpireArgs 以毫秒为单位为哈希域设置过期时间。
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

// HExpireAt 为哈希 key 中给定域设置由 tm 指定的秒级过期时刻。
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

// HExpireAtWithArgs 按给定 HExpireArgs 为哈希域设置由 tm 指定的秒级过期时刻。
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

// HPExpireAt 以毫秒时间戳为哈希 key 中给定域设置过期时刻。
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

// HPExpireAtWithArgs 按给定 HExpireArgs 以毫秒时间戳为哈希域设置过期时刻。
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

// HPersist 移除哈希 key 中给定域的过期时间，使其持久保留。
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

// HExpireTime 返回哈希 key 中给定域的秒级过期时刻，状态码与返回切片一一对应。
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

// HPExpireTime 以毫秒为单位返回哈希 key 中给定域的过期时刻。
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

// HTTL 返回哈希 key 中给定域的剩余生存时间（秒）。
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

// HPTTL 以毫秒为单位返回哈希 key 中给定域的剩余生存时间。
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
