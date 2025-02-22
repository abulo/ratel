package redis

import (
	"context"
)

// PFAdd 将指定元素添加到HyperLogLog key: HyperLogLog键名, els: 要添加的元素列表
func (r *Client) PFAdd(ctx context.Context, key string, els ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFAdd(getCtx(ctx), key, els...).Result()
		return err
	}, acceptable)
	return
}

// PFCount 返回HyperlogLog观察到的集合的近似基数。 keys: 要查询的HyperLogLog键名列表
func (r *Client) PFCount(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFCount(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// PFMerge N个不同的HyperLogLog合并为一个。 dest: 目标HyperLogLog键名, keys: 要合并的HyperLogLog键名列表
func (r *Client) PFMerge(ctx context.Context, dest string, keys ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFMerge(getCtx(ctx), dest, keys...).Result()
		return err
	}, acceptable)
	return
}
