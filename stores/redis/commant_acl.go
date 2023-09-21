package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *Client) ACLDryRun(ctx context.Context, username string, command ...interface{}) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLDryRun(getCtx(ctx), username, command...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ACLLog(ctx context.Context, count int64) (val []*redis.ACLLogEntry, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLLog(getCtx(ctx), count).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ACLLogReset(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLLogReset(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
