package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Eval 执行Lua脚本 ctx:上下文 script:Lua脚本 keys:键列表 args:参数列表
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

// EvalSha 通过SHA1执行Lua脚本 ctx:上下文 sha1:脚本SHA1值 keys:键列表 args:参数列表
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

// EvalRO 以只读模式执行Lua脚本 ctx:上下文 script:Lua脚本 keys:键列表 args:参数列表
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

// EvalShaRO 以只读模式通过SHA1执行Lua脚本 ctx:上下文 sha1:脚本SHA1值 keys:键列表 args:参数列表
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

// ScriptExists 检查脚本是否存在于缓存中 ctx:上下文 hashes:脚本SHA1值列表
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

// ScriptFlush 清空脚本缓存 ctx:上下文
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

// ScriptKill 终止正在运行的Lua脚本 ctx:上下文
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

// ScriptLoad 加载脚本到缓存但不执行 ctx:上下文 script:Lua脚本
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

// FunctionLoad 加载Redis函数 ctx:上下文 code:函数代码
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

// FunctionLoadReplace 替换已存在的Redis函数 ctx:上下文 code:函数代码
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

// FunctionDelete 删除Redis函数 ctx:上下文 libName:函数库名称
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

// FunctionFlush 清空所有Redis函数 ctx:上下文
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

// FunctionKill 终止正在运行的Redis函数 ctx:上下文
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

// FunctionFlushAsync 异步清空所有Redis函数 ctx:上下文
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

// FunctionList 列出所有Redis函数 ctx:上下文 q:查询条件
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

// FunctionDump 导出Redis函数 ctx:上下文
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

// FunctionRestore 恢复Redis函数 ctx:上下文 libDump:函数库数据
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

// FunctionStats 获取Redis函数统计信息 ctx:上下文
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

// FCall 调用Redis函数 ctx:上下文 function:函数名称 keys:键列表 args:参数列表
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

// FCallRo 以只读模式调用Redis函数 ctx:上下文 function:函数名称 keys:键列表 args:参数列表
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

// FCallRO 以只读模式调用Redis函数 ctx:上下文 function:函数名称 keys:键列表 args:参数列表
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
