package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GetBit 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。key: 键名, offset: 位偏移量。当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
func (r *Client) GetBit(ctx context.Context, key string, offset int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetBit(getCtx(ctx), key, offset).Result()
		return err
	}, acceptable)
	return
}

// SetBit 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。key: 键名, offset: 位偏移量, value: 要设置的值(0或1)。当 key 不存在时，自动生成一个新的字符串值。字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
func (r *Client) SetBit(ctx context.Context, key string, offset int64, value int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetBit(getCtx(ctx), key, offset, value).Result()
		return err
	}, acceptable)
	return
}

// BitCount 计算给定字符串中，被设置为 1 的比特位的数量。key: 键名, bitCount: 计数范围(start/end)。一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
func (r *Client) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitCount(getCtx(ctx), key, bitCount).Result()
		return err
	}, acceptable)
	return
}

// BitOpAnd 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destKey 上。destKey: 目标键名, keys: 源键名列表。
func (r *Client) BitOpAnd(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpAnd(getCtx(ctx), destKey, keys...).Result()
		return err
	}, acceptable)
	return
}

// BitOpOr 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destKey 上。destKey: 目标键名, keys: 源键名列表。
func (r *Client) BitOpOr(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpOr(getCtx(ctx), destKey, keys...).Result()
		return err
	}, acceptable)
	return
}

// BitOpXor 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destKey 上。destKey: 目标键名, keys: 源键名列表。
func (r *Client) BitOpXor(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpXor(getCtx(ctx), destKey, keys...).Result()
		return err
	}, acceptable)
	return
}

// BitOpNot 对一个保存二进制位的字符串 key 进行位非操作，并将结果保存到 destKey 上。destKey: 目标键名, key: 源键名。
func (r *Client) BitOpNot(ctx context.Context, destKey string, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpNot(getCtx(ctx), destKey, key).Result()
		return err
	}, acceptable)
	return
}

// BitPos 返回字符串里面第一个被设置为1或者0的bit位。key: 键名, bit: 要查找的位值(0或1), pos: 可选起始位置。
func (r *Client) BitPos(ctx context.Context, key string, bit int64, pos ...int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitPos(getCtx(ctx), key, bit, pos...).Result()
		return err
	}, acceptable)
	return
}

// BitPosSpan 返回字符串里面第一个被设置为1或者0的bit位。key: 键名, bit: 要查找的位值(0或1), start: 起始位置, end: 结束位置, span: 范围类型。
func (r *Client) BitPosSpan(ctx context.Context, key string, bit int8, start, end int64, span string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitPosSpan(getCtx(ctx), key, bit, start, end, span).Result()
		return err
	}, acceptable)
	return
}

// BitField 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。key: 键名, args: 位域操作参数列表。
func (r *Client) BitField(ctx context.Context, key string, args ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitField(getCtx(ctx), key, args...).Result()
		return err
	}, acceptable)
	return
}

// BitFieldRO 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。key: 键名, args: 位域操作参数列表。
func (r *Client) BitFieldRO(ctx context.Context, key string, args ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitFieldRO(getCtx(ctx), key, args...).Result()
		return err
	}, acceptable)
	return
}
