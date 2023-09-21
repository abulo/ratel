package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GetBit 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
// 字符串值指定偏移量上的位(bit)。
func (r *Client) GetBit(ctx context.Context, key any, offset int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetBit(getCtx(ctx), r.k(key), offset).Result()
		return err
	}, acceptable)
	return
}

// SetBit 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 当 key 不存在时，自动生成一个新的字符串值。
// 字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
// offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
func (r *Client) SetBit(ctx context.Context, key any, offset int64, value int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetBit(getCtx(ctx), r.k(key), offset, value).Result()
		return err
	}, acceptable)
	return
}

// BitCount 计算给定字符串中，被设置为 1 的比特位的数量。
// 一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
// start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。
// 不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
func (r *Client) BitCount(ctx context.Context, key any, bitCount *redis.BitCount) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitCount(getCtx(ctx), r.k(key), bitCount).Result()
		return err
	}, acceptable)
	return
}

// BitOpAnd -> BitCount
func (r *Client) BitOpAnd(ctx context.Context, destKey any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpAnd(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpOr -> BitCount
func (r *Client) BitOpOr(ctx context.Context, destKey any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpOr(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpXor -> BitCount
func (r *Client) BitOpXor(ctx context.Context, destKey any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpXor(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpNot -> BitCount
func (r *Client) BitOpNot(ctx context.Context, destKey any, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpNot(getCtx(ctx), r.k(destKey), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// BitPos -> BitCount
func (r *Client) BitPos(ctx context.Context, key any, bit int64, pos ...int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitPos(getCtx(ctx), r.k(key), bit, pos...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) BitPosSpan(ctx context.Context, key any, bit int8, start, end int64, span string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitPosSpan(getCtx(ctx), r.k(key), bit, start, end, span).Result()
		return err
	}, acceptable)
	return
}

// BitField -> BitCount
func (r *Client) BitField(ctx context.Context, key any, args ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitField(getCtx(ctx), r.k(key), args...).Result()
		return err
	}, acceptable)
	return
}
