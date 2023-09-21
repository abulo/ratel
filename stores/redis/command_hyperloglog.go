package redis

import "context"

// PFAdd 将指定元素添加到HyperLogLog
func (r *Client) PFAdd(ctx context.Context, key any, els ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFAdd(getCtx(ctx), r.k(key), els...).Result()
		return err
	}, acceptable)
	return
}

// PFCount 返回HyperlogLog观察到的集合的近似基数。
func (r *Client) PFCount(ctx context.Context, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFCount(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// PFMerge N个不同的HyperLogLog合并为一个。
func (r *Client) PFMerge(ctx context.Context, dest any, keys ...any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFMerge(getCtx(ctx), r.k(dest), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
