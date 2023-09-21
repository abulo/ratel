package redis

import "context"

// DBSize 返回当前数据库的 key 的数量
func (r *Client) DBSize(ctx context.Context) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DBSize(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
