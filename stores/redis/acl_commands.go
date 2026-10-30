package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ACLDryRun 不实际执行命令，仅检查用户 username 是否有权执行 command，返回权限判定结果。
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

// ACLLog 返回最近 count 条 ACL 拒绝与安全事件日志记录。
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

// ACLLogReset 清空当前服务器内存中的 ACL 日志。
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

// ACLSetUser 按 rules 配置用户 username 的 ACL 规则（密码、命令权限、键权限等），规则语法由服务器定义。
func (r *Client) ACLSetUser(ctx context.Context, username string, rules ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLSetUser(getCtx(ctx), username, rules...).Result()
		return err
	}, acceptable)
	return
}

// ACLDelUser 删除指定用户 username，返回实际被删除的用户数量。
func (r *Client) ACLDelUser(ctx context.Context, username string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLDelUser(getCtx(ctx), username).Result()
		return err
	}, acceptable)
	return
}

// ACLList 返回所有用户 ACL 规则的原始文本列表。
func (r *Client) ACLList(ctx context.Context) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLList(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ACLCat 返回服务器全部命令类别名称。
func (r *Client) ACLCat(ctx context.Context) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLCat(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ACLCatArgs 按 ACLCatArgs 查询指定命令类别下的命令列表。
func (r *Client) ACLCatArgs(ctx context.Context, options *redis.ACLCatArgs) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ACLCatArgs(getCtx(ctx), options).Result()
		return err
	}, acceptable)
	return
}
