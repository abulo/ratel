package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Eval 执行 Lua 脚本 script，keys 为脚本可见的键列表，args 为附加参数，返回脚本执行结果。
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Eval(getCtx(ctx), script, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// EvalSha 按脚本 SHA1 校验和执行已缓存的 Lua 脚本；脚本不在缓存中时返回 NOSCRIPT 错误。
func (r *Client) EvalSha(ctx context.Context, sha1 string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalSha(getCtx(ctx), sha1, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// EvalRO 以只读模式执行 Lua 脚本，脚本内不允许执行写命令。
func (r *Client) EvalRO(ctx context.Context, script string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalRO(getCtx(ctx), script, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// EvalShaRO 以只读模式按 SHA1 执行已缓存 Lua 脚本，脚本内不允许执行写命令。
func (r *Client) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalShaRO(getCtx(ctx), sha1, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// ScriptExists 批量检查给定 SHA1 脚本是否已缓存在服务器，返回布尔切片顺序与 hashes 一致。
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

// ScriptFlush 清空服务器内存中的 Lua 脚本缓存。
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

// ScriptKill 终止当前正在运行的 Lua 脚本；无运行中脚本时返回错误。
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

// ScriptLoad 将 Lua 脚本加载到服务器缓存但不执行，返回 SHA1 校验和供 EvalSha 使用。
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

// FunctionLoad 加载 Redis 函数库代码 code，返回函数库名称；已存在同名函数库时报错。
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

// FunctionLoadReplace 加载函数库代码 code，并替换服务器上已存在的同名函数库。
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

// FunctionDelete 删除指定的 Redis 函数库 libName。
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

// FunctionFlush 清空服务器上全部已加载的 Redis 函数库。
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

// FunctionKill 终止当前正在运行的 Redis 函数；无运行中函数时返回错误。
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

// FunctionFlushAsync 异步清空全部 Redis 函数库，命令不等待清理完成即返回。
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

// FunctionList 按 FunctionListQuery 列出函数库及其包含函数的定义信息。
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

// FunctionDump 导出全部已加载函数库的二进制快照，可配合 FunctionRestore 在其他实例恢复。
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

// FunctionRestore 用 FunctionDump 产生的快照恢复函数库。
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

// FunctionStats 返回当前 Redis 函数的运行状态与统计信息。
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

// FCall 调用指定 Redis 函数，keys 为脚本可见的键列表，args 为附加参数，返回函数返回值。
func (r *Client) FCall(ctx context.Context, function string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCall(getCtx(ctx), function, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// FCallRo 以只读模式调用 Redis 函数，函数内不允许执行写命令。
func (r *Client) FCallRo(ctx context.Context, function string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCallRo(getCtx(ctx), function, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// FCallRO 以只读模式调用 Redis 函数，函数内不允许执行写命令。
func (r *Client) FCallRO(ctx context.Context, function string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FCallRO(getCtx(ctx), function, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// ScriptRunner 包装 go-redis 的 *redis.Script：固定脚本在服务端缓存 SHA1，运行期优先走
// EVALSHA，收到 NOSCRIPT 时由 go-redis 自动重新加载并回退 EVAL，避免每次调用都重复传输整段
// Lua 源码。实例本身可并发使用。
//
// PERF-1-5 [P2] 复用型脚本的网络传输与服务端解析优化；收益取决于脚本大小与调用频次（预期收益需实测验证）。
type ScriptRunner struct {
	client *Client
	script *redis.Script
}

// NewScriptRunner 基于脚本源码 src 创建可复用的脚本执行器；同一脚本建议只创建一次并长期持有。
func (r *Client) NewScriptRunner(src string) *ScriptRunner {
	return &ScriptRunner{
		client: r,
		script: redis.NewScript(src),
	}
}

// Run 以写模式执行脚本，keys 为脚本可见的键列表，args 为附加参数，返回脚本执行结果。
func (sr *ScriptRunner) Run(ctx context.Context, keys []string, args ...any) (val any, err error) {
	err = sr.client.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(sr.client)
		if err != nil {
			return err
		}
		val, err = sr.script.Run(getCtx(ctx), conn, keys, args...).Result()
		return err
	}, acceptable)
	return
}

// RunRO 以只读模式执行脚本，脚本内不允许执行写命令。
func (sr *ScriptRunner) RunRO(ctx context.Context, keys []string, args ...any) (val any, err error) {
	err = sr.client.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(sr.client)
		if err != nil {
			return err
		}
		val, err = sr.script.RunRO(getCtx(ctx), conn, keys, args...).Result()
		return err
	}, acceptable)
	return
}
