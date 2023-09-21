package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Eval 执行 Lua 脚本。
func (r *Client) Eval(ctx context.Context, script string, keys []any, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Eval(getCtx(ctx), script, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalSha 执行 Lua 脚本。
func (r *Client) EvalSha(ctx context.Context, sha1 string, keys []any, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalSha(getCtx(ctx), sha1, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalRO 执行 Lua 脚本。
func (r *Client) EvalRO(ctx context.Context, script string, keys []any, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalRO(getCtx(ctx), script, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalShaRO 执行 Lua 脚本。
func (r *Client) EvalShaRO(ctx context.Context, sha1 string, keys []any, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalShaRO(getCtx(ctx), sha1, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// ScriptExists 查看指定的脚本是否已经被保存在缓存当中。
func (r *Client) ScriptExists(ctx context.Context, hashes ...string) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptExists(getCtx(ctx), hashes...).Result()
		return err
	}, acceptable)
	return
}

// ScriptFlush 从脚本缓存中移除所有脚本。
func (r *Client) ScriptFlush(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptFlush(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ScriptKill 杀死当前正在运行的 Lua 脚本。
func (r *Client) ScriptKill(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptKill(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ScriptLoad 将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本。
func (r *Client) ScriptLoad(ctx context.Context, script string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptLoad(getCtx(ctx), script).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) FunctionLoad(ctx context.Context, code string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionLoad(getCtx(ctx), code).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionLoadReplace(ctx context.Context, code string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionLoadReplace(getCtx(ctx), code).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionDelete(ctx context.Context, libName string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionDelete(getCtx(ctx), libName).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionFlush(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionFlush(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) FunctionKill(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionKill(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionFlushAsync(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionFlushAsync(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionList(ctx context.Context, q redis.FunctionListQuery) (val []redis.Library, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionList(getCtx(ctx), q).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionDump(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionDump(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionRestore(ctx context.Context, libDump string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionRestore(getCtx(ctx), libDump).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FunctionStats(ctx context.Context) (val redis.FunctionStats, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FunctionStats(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FCall(ctx context.Context, function string, keys []any, args ...interface{}) (val interface{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCall(getCtx(ctx), function, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FCallRo(ctx context.Context, function string, keys []any, args ...interface{}) (val interface{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCallRo(getCtx(ctx), function, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) FCallRO(ctx context.Context, function string, keys []any, args ...interface{}) (val interface{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCallRO(getCtx(ctx), function, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}
