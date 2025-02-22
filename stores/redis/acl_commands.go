package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ACLDryRun 用于在不实际执行命令的情况下检查用户是否具有执行给定命令的权限。ctx:上下文 username:用户名 command:要检查的命令
func (r *Client) ACLDryRun(ctx context.Context, username string, command ...any) (val string, err error) {
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

// ACLLog 获取 ACL 日志。ctx:上下文 count:要获取的日志条数
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

// ACLLogReset 重置 ACL 日志。ctx:上下文
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
