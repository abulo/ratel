package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// JSONArrAppend 在 key 的 path 对应 JSON 数组末尾追加 values，返回各路径更新后的数组长度。
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

// JSONArrIndex 返回 path 数组中 value 首次出现的索引；未找到时对应位置为 nil。
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

// JSONArrIndexWithArgs 按 JSONArrIndexArgs（起止位置、范围）查找 value 在数组中的索引。
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

// JSONArrInsert 在 path 数组的 index 位置插入 values，返回各路径更新后的数组长度。
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

// JSONArrLen 返回 path 对应 JSON 数组的长度切片，非数组路径对应位置为 nil。
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

// JSONArrPop 弹出 path 数组中 index 位置的元素并返回；index 为负表示从尾部计数。
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

// JSONArrTrim 将 path 数组修剪为合法长度；本封装使用默认起止位置，返回各路径修剪后长度。
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

// JSONArrTrimWithArgs 按 JSONArrTrimArgs 指定的起止位置修剪数组，返回修剪后长度。
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

// JSONClear 清空 path 下的 JSON 对象与数组内容（标量置空），返回被清除的值数量。
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

// JSONDebugMemory 返回 key 的 path 下 JSON 值占用的内存字节数。
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

// JSONDel 删除 key 的 path 下的 JSON 值，返回被删除的路径数量。
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

// JSONForget 是 JSONDel 的别名，删除指定 path 下的 JSON 值。
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

// JSONGet 返回 key 在 paths 路径下的 JSON 文本；多路径时结果为按路径组织的 JSON。
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

// JSONGetWithArgs 按 JSONGetArgs（换行、缩进、转义等）读取指定 paths 的 JSON 值。
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

// JSONMerge 将 value 合并到 key 的 path 下，对象按键递归合并，标量与数组直接覆盖。
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

// JSONMSetArgs 按 JSONSetArgs 列表批量写入多个键的 JSON 值。
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

// JSONMSet 批量写入 JSON 值，参数按 键、路径、值 成对传入。
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

// JSONMGet 批量读取多个 keys 在同一 path 下的 JSON 值，切片顺序与 keys 一致。
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

// JSONNumIncrBy 将 key 的 path 数值增加 value，返回更新后的 JSON 数值文本；非数值路径报错。
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

// JSONObjKeys 返回 path 对应 JSON 对象的顶层键名列表切片。
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

// JSONObjLen 返回 path 对应 JSON 对象的成员数量切片，非对象路径对应位置为 nil。
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

// JSONSet 将 key 的 path 设为 value；路径不存在时自动创建对象或数组层级。
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

// JSONSetMode 按 mode（NX/XX）条件设置 JSON 值，仅在路径不存在或已存在时写入。
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

// JSONStrAppend 在 key 的 path 字符串末尾追加 value，返回各路径更新后的字符串长度。
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

// JSONStrLen 返回 path 对应 JSON 字符串的长度切片，非字符串路径对应位置为 nil。
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

// JSONToggle 翻转 path 处布尔值（true/false 互换），返回翻转后布尔状态切片。
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

// JSONType 返回各 path 下 JSON 值的类型（object/array/string/number/boolean/null）切片。
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
