package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// TFunctionLoad 将 Lua 函数库加载到 Redis 实例中 ctx: 上下文, lib: Lua库代码
func (r *Client) TFunctionLoad(ctx context.Context, lib string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionLoad(getCtx(ctx), lib).Result()
		return err
	}, acceptable)
	return
}

// TFunctionLoadArgs 将 Lua 函数库加载到 Redis 实例中 ctx: 上下文, lib: Lua库代码, options: 加载选项
func (r *Client) TFunctionLoadArgs(ctx context.Context, lib string, options *redis.TFunctionLoadOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionLoadArgs(getCtx(ctx), lib, options).Result()
		return err
	}, acceptable)
	return
}

// TFunctionDelete 从 Redis 实例中删除 Lua 函数库 ctx: 上下文, libName: 要删除的库名称
func (r *Client) TFunctionDelete(ctx context.Context, libName string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionDelete(getCtx(ctx), libName).Result()
		return err
	}, acceptable)
	return
}

// TFunctionList 列出 Redis 实例中的所有 Lua 函数库 ctx: 上下文
func (r *Client) TFunctionList(ctx context.Context) (val []map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionList(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// TFunctionListArgs 列出 Redis 实例中的所有 Lua 函数库 ctx: 上下文, options: 列表选项
func (r *Client) TFunctionListArgs(ctx context.Context, options *redis.TFunctionListOptions) (val []map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionListArgs(getCtx(ctx), options).Result()
		return err
	}, acceptable)
	return
}

// TFCall 调用 Redis 实例中的 Lua 函数 ctx: 上下文, libName: 库名称, funcName: 函数名称, numKeys: 键数量
func (r *Client) TFCall(ctx context.Context, libName, funcName string, numKeys int) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCall(getCtx(ctx), libName, funcName, numKeys).Result()
		return err
	}, acceptable)
	return
}

// TFCallArgs 调用 Redis 实例中的 Lua 函数 ctx: 上下文, libName: 库名称, funcName: 函数名称, numKeys: 键数量, options: 调用选项
func (r *Client) TFCallArgs(ctx context.Context, libName, funcName string, numKeys int, options *redis.TFCallOptions) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallArgs(getCtx(ctx), libName, funcName, numKeys, options).Result()
		return err
	}, acceptable)
	return
}

// TFCallASYNC 异步调用 Redis 实例中的 Lua 函数 ctx: 上下文, libName: 库名称, funcName: 函数名称, numKeys: 键数量
func (r *Client) TFCallASYNC(ctx context.Context, libName, funcName string, numKeys int) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallASYNC(getCtx(ctx), libName, funcName, numKeys).Result()
		return err
	}, acceptable)
	return
}

// TFCallASYNCArgs 异步调用 Redis 实例中的 Lua 函数 ctx: 上下文, libName: 库名称, funcName: 函数名称, numKeys: 键数量, options: 调用选项
func (r *Client) TFCallASYNCArgs(ctx context.Context, libName, funcName string, numKeys int, options *redis.TFCallOptions) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallASYNCArgs(getCtx(ctx), libName, funcName, numKeys, options).Result()
		return err
	}, acceptable)
	return
}
