package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// GetBit 返回 key 字符串在偏移量 offset 上的二进制位（0 或 1）。
//
// 当 offset 超出字符串长度或 key 不存在时，返回 0。
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

// SetBit 将 key 字符串在偏移量 offset 上的二进制位设置为 value（0 或 1），返回设置前该位的旧值。
//
// key 不存在时自动新建字符串；字符串会按偏移量伸展，空白位置以 0 填充；offset 必须非负且小于 2^32。
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

// BitCount 统计 key 字符串中被置为 1 的二进制位数量，可选通过 bitCount 指定字节范围 [start, end]。
//
// start、end 支持负数定位尾部；不存在的 key 按空字符串处理，计数结果为 0。
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

// BitOpAnd 将多个源键按位逻辑与的结果写入 destKey，返回 destKey 的字节长度。
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

// BitOpOr 将多个源键按位逻辑或的结果写入 destKey，返回 destKey 的字节长度。
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

// BitOpXor 将多个源键按位异或的结果写入 destKey，返回 destKey 的字节长度。
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

// BitOpNot 对单个源键按位取反，结果写入 destKey，返回 destKey 的字节长度。
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

// BitPos 返回 key 字符串中第一个被置为 bit（0 或 1）的位置；pos 可指定搜索起始位置。
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

// BitPosSpan 在 [start, end] 范围内查找第一个为 bit（0 或 1）的位置，span 指定字节或位范围类型。
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

// BitField 按 args 对位域执行读、写、自增等复合操作，返回各子命令的结果切片。
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

// BitFieldRO 只读模式执行位域读取操作，args 中不允许包含修改类子命令。
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
