package redis

import (
	"context"
)

// PFAdd 将 els 加入 HyperLogLog key，返回 HLL 内部估算是否发生改变（1 表示已更新）。
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

// PFCount 返回给定 keys 合并去重后的近似基数（元素数量估计值）。
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

// PFMerge 将多个 HyperLogLog keys 合并到 dest，合并结果的近似基数等价于这些集合的并集。
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
