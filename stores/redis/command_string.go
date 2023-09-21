package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Append 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (r *Client) Append(ctx context.Context, key any, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Append(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// Decr 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 执行 DECR 命令之后 key 的值。
func (r *Client) Decr(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Decr(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// DecrBy 将 key 所储存的值减去减量 decrement 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于更多递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 减去 decrement 之后， key 的值。
func (r *Client) DecrBy(ctx context.Context, key any, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DecrBy(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// Get 返回 key 所关联的字符串值。
// 如果 key 不存在那么返回特殊值 nil 。
// 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。
// 当 key 不存在时，返回 nil ，否则，返回 key 的值。
// 如果 key 不是字符串类型，那么返回一个错误。
func (r *Client) Get(ctx context.Context, key any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Get(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// GetRange 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
// 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
// GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
// 返回截取得出的子字符串。
func (r *Client) GetRange(ctx context.Context, key any, start, end int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetRange(getCtx(ctx), r.k(key), start, end).Result()
		return err
	}, acceptable)
	return
}

// GetSet 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 返回给定 key 的旧值。
// 当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
func (r *Client) GetSet(ctx context.Context, key any, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetSet(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// GetEx
func (r *Client) GetEx(ctx context.Context, key any, ts time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetEx(getCtx(ctx), r.k(key), ts).Result()
		return err
	}, acceptable)
	return
}

// GetEx
func (r *Client) GetDel(ctx context.Context, key any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetDel(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Incr 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 执行 INCR 命令之后 key 的值。
func (r *Client) Incr(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Incr(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// IncrBy 将 key 所储存的值加上增量 increment 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，参见 INCR 命令。
// 加上 increment 之后， key 的值。
func (r *Client) IncrBy(ctx context.Context, key any, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrBy(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// IncrByFloat 为 key 中所储存的值加上浮点数增量 increment 。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
// 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
func (r *Client) IncrByFloat(ctx context.Context, key any, value float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrByFloat(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

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

// MGet 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
// 一个包含所有给定 key 的值的列表。
func (r *Client) MGet(ctx context.Context, keys ...any) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MGet(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// MSet 同时设置一个或多个 key-value 对。
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

// MSetNX 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
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

// Set 将字符串值 value 关联到 key 。
// 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
// 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
func (r *Client) Set(ctx context.Context, key any, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Set(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) SetArgs(ctx context.Context, key any, value any, a redis.SetArgs) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetArgs(getCtx(ctx), r.k(key), value, a).Result()
		return err
	}, acceptable)
	return
}

// SetEX ...
func (r *Client) SetEx(ctx context.Context, key any, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetEx(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetNX 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
// SETNX 是『SET if Not eXists』(如果不存在，则 SET)的简写。
func (r *Client) SetNX(ctx context.Context, key any, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetNX(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetXX -> SetNX
func (r *Client) SetXX(ctx context.Context, key any, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetXX(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetRange 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (r *Client) SetRange(ctx context.Context, key any, offset int64, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetRange(getCtx(ctx), r.k(key), offset, value).Result()
		return err
	}, acceptable)
	return
}

// StrLen 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。
func (r *Client) StrLen(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.StrLen(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}
