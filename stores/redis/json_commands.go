package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// JSONArrAppend 向JSON数组追加元素 ctx: 上下文, key: Redis键名, path: JSON路径, values: 要追加的值
func (r *Client) JSONArrAppend(ctx context.Context, key, path string, values ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrAppend(getCtx(ctx), key, path, values...).Result()
		return err
	}, acceptable)
	return
}

// JSONArrIndex 查找JSON数组中元素的索引 ctx: 上下文, key: Redis键名, path: JSON路径, value: 要查找的值
func (r *Client) JSONArrIndex(ctx context.Context, key, path string, value ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrIndex(getCtx(ctx), key, path, value...).Result()
		return err
	}, acceptable)
	return
}

// JSONArrIndexWithArgs 带参数查找JSON数组中元素的索引 ctx: 上下文, key: Redis键名, path: JSON路径, options: 查找选项, value: 要查找的值
func (r *Client) JSONArrIndexWithArgs(ctx context.Context, key, path string, options *redis.JSONArrIndexArgs, value ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrIndexWithArgs(getCtx(ctx), key, path, options, value...).Result()
		return err
	}, acceptable)
	return
}

// JSONArrInsert 向JSON数组插入元素 ctx: 上下文, key: Redis键名, path: JSON路径, index: 插入位置, values: 要插入的值
func (r *Client) JSONArrInsert(ctx context.Context, key, path string, index int64, values ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrInsert(getCtx(ctx), key, path, index, values...).Result()
		return err
	}, acceptable)
	return
}

// JSONArrLen 获取JSON数组长度 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONArrLen(ctx context.Context, key, path string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrLen(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONArrPop 从JSON数组弹出元素 ctx: 上下文, key: Redis键名, path: JSON路径, index: 弹出位置
func (r *Client) JSONArrPop(ctx context.Context, key, path string, index int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrPop(getCtx(ctx), key, path, index).Result()
		return err
	}, acceptable)
	return
}

// JSONArrTrim 修剪JSON数组 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONArrTrim(ctx context.Context, key, path string) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrTrim(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONArrTrimWithArgs 带参数修剪JSON数组 ctx: 上下文, key: Redis键名, path: JSON路径, options: 修剪选项
func (r *Client) JSONArrTrimWithArgs(ctx context.Context, key, path string, options *redis.JSONArrTrimArgs) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONArrTrimWithArgs(getCtx(ctx), key, path, options).Result()
		return err
	}, acceptable)
	return
}

// JSONClear 清除JSON值 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONClear(ctx context.Context, key, path string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONClear(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONDebugMemory 调试JSON内存使用 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONDebugMemory(ctx context.Context, key, path string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONDebugMemory(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONDel 删除JSON值 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONDel(ctx context.Context, key, path string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONDel(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONForget 删除JSON值（别名） ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONForget(ctx context.Context, key, path string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONForget(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONGet 获取JSON值 ctx: 上下文, key: Redis键名, paths: JSON路径
func (r *Client) JSONGet(ctx context.Context, key string, paths ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONGet(getCtx(ctx), key, paths...).Result()
		return err
	}, acceptable)
	return
}

// JSONGetWithArgs 带参数获取JSON值 ctx: 上下文, key: Redis键名, options: 获取选项, paths: JSON路径
func (r *Client) JSONGetWithArgs(ctx context.Context, key string, options *redis.JSONGetArgs, paths ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONGetWithArgs(getCtx(ctx), key, options, paths...).Result()
		return err
	}, acceptable)
	return
}

// JSONMerge 合并JSON值 ctx: 上下文, key: Redis键名, path: JSON路径, value: 要合并的值
func (r *Client) JSONMerge(ctx context.Context, key, path string, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.JSONMerge(getCtx(ctx), key, path, value).Result()
		return err
	}, acceptable)
	return
}

// JSONMSetArgs 批量设置JSON值（带参数） ctx: 上下文, docs: JSON设置参数列表
func (r *Client) JSONMSetArgs(ctx context.Context, docs []redis.JSONSetArgs) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONMSetArgs(getCtx(ctx), docs).Result()
		return err
	}, acceptable)
	return
}

// JSONMSet 批量设置JSON值 ctx: 上下文, params: 设置参数
func (r *Client) JSONMSet(ctx context.Context, params ...any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONMSet(getCtx(ctx), params...).Result()
		return err
	}, acceptable)
	return
}

// JSONMGet 批量获取JSON值 ctx: 上下文, path: JSON路径, keys: Redis键名列表
func (r *Client) JSONMGet(ctx context.Context, path string, keys ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONMGet(getCtx(ctx), path, keys...).Result()
		return err
	}, acceptable)
	return
}

// JSONNumIncrBy 增加JSON数值 ctx: 上下文, key: Redis键名, path: JSON路径, value: 增加值
func (r *Client) JSONNumIncrBy(ctx context.Context, key, path string, value float64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONNumIncrBy(getCtx(ctx), key, path, value).Result()
		return err
	}, acceptable)
	return
}

// JSONObjKeys 获取JSON对象键名列表 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONObjKeys(ctx context.Context, key, path string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONObjKeys(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}

// JSONObjLen 获取JSON对象长度 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONObjLen(ctx context.Context, key, path string) (val []*int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONObjLen(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return

}

// JSONSet 设置JSON值 ctx: 上下文, key: Redis键名, path: JSON路径, value: 要设置的值
func (r *Client) JSONSet(ctx context.Context, key, path string, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONSet(getCtx(ctx), key, path, value).Result()
		return err
	}, acceptable)
	return
}

// JSONSetMode 带模式设置JSON值 ctx: 上下文, key: Redis键名, path: JSON路径, value: 要设置的值, mode: 设置模式
func (r *Client) JSONSetMode(ctx context.Context, key, path string, value any, mode string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONSetMode(getCtx(ctx), key, path, value, mode).Result()
		return err
	}, acceptable)
	return
}

// JSONStrAppend 向JSON字符串追加内容 ctx: 上下文, key: Redis键名, path: JSON路径, value: 要追加的值
func (r *Client) JSONStrAppend(ctx context.Context, key, path, value string) (val []*int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONStrAppend(getCtx(ctx), key, path, value).Result()
		return err
	}, acceptable)
	return

}

// JSONStrLen 获取JSON字符串长度 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONStrLen(ctx context.Context, key, path string) (val []*int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONStrLen(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return

}

// JSONToggle 切换JSON布尔值 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONToggle(ctx context.Context, key, path string) (val []*int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONToggle(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return

}

// JSONType 获取JSON值类型 ctx: 上下文, key: Redis键名, path: JSON路径
func (r *Client) JSONType(ctx context.Context, key, path string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.JSONType(getCtx(ctx), key, path).Result()
		return err
	}, acceptable)
	return
}
