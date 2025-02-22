package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// TSAdd 向时间序列中添加数据点 ctx: 上下文, key: 时间序列键名, timestamp: 时间戳, value: 数据值
func (r *Client) TSAdd(ctx context.Context, key string, timestamp any, value float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAdd(getCtx(ctx), key, timestamp, value).Result()
		return err
	}, acceptable)
	return
}

// TSAddWithArgs 向时间序列中添加数据点（带选项） ctx: 上下文, key: 时间序列键名, timestamp: 时间戳, value: 数据值, options: 时间序列选项
func (r *Client) TSAddWithArgs(ctx context.Context, key string, timestamp any, value float64, options *redis.TSOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAddWithArgs(getCtx(ctx), key, timestamp, value, options).Result()
		return err
	}, acceptable)
	return
}

// TSCreate 创建新的时间序列 ctx: 上下文, key: 时间序列键名
func (r *Client) TSCreate(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreate(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TSCreateWithArgs 创建新的时间序列（带选项） ctx: 上下文, key: 时间序列键名, options: 时间序列选项
func (r *Client) TSCreateWithArgs(ctx context.Context, key string, options *redis.TSOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreateWithArgs(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// TSAlter 修改时间序列属性 ctx: 上下文, key: 时间序列键名, options: 修改选项
func (r *Client) TSAlter(ctx context.Context, key string, options *redis.TSAlterOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAlter(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// TSCreateRule 创建时间序列聚合规则 ctx: 上下文, sourceKey: 源时间序列键名, destKey: 目标时间序列键名, aggregator: 聚合器类型, bucketDuration: 时间桶大小
func (r *Client) TSCreateRule(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreateRule(getCtx(ctx), sourceKey, destKey, aggregator, bucketDuration).Result()
		return err
	}, acceptable)
	return
}

// TSCreateRuleWithArgs 创建时间序列聚合规则（带选项） ctx: 上下文, sourceKey: 源时间序列键名, destKey: 目标时间序列键名, aggregator: 聚合器类型, bucketDuration: 时间桶大小, options: 创建规则选项
func (r *Client) TSCreateRuleWithArgs(ctx context.Context, sourceKey string, destKey string, aggregator redis.Aggregator, bucketDuration int, options *redis.TSCreateRuleOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreateRuleWithArgs(getCtx(ctx), sourceKey, destKey, aggregator, bucketDuration, options).Result()
		return err
	}, acceptable)
	return
}

// TSIncrBy 增加时间序列值 ctx: 上下文, key: 时间序列键名, timestamp: 时间戳
func (r *Client) TSIncrBy(ctx context.Context, key string, timestamp float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSIncrBy(getCtx(ctx), key, timestamp).Result()
		return err
	}, acceptable)
	return
}

// TSIncrByWithArgs 增加时间序列值（带选项） ctx: 上下文, key: 时间序列键名, timestamp: 时间戳, options: 增减选项
func (r *Client) TSIncrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSIncrByWithArgs(getCtx(ctx), key, timestamp, options).Result()
		return err
	}, acceptable)
	return
}

// TSDecrBy 减少时间序列值 ctx: 上下文, key: 时间序列键名, timestamp: 时间戳
func (r *Client) TSDecrBy(ctx context.Context, key string, timestamp float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSDecrBy(getCtx(ctx), key, timestamp).Result()
		return err
	}, acceptable)
	return
}

// TSDecrByWithArgs 减少时间序列值（带选项） ctx: 上下文, key: 时间序列键名, timestamp: 时间戳, options: 增减选项
func (r *Client) TSDecrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSDecrByWithArgs(getCtx(ctx), key, timestamp, options).Result()
		return err
	}, acceptable)
	return
}

// TSDel 删除时间序列数据点 ctx: 上下文, key: 时间序列键名, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳
func (r *Client) TSDel(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSDel(getCtx(ctx), key, fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

// TSDeleteRule 删除时间序列聚合规则 ctx: 上下文, sourceKey: 源时间序列键名, destKey: 目标时间序列键名
func (r *Client) TSDeleteRule(ctx context.Context, sourceKey string, destKey string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSDeleteRule(getCtx(ctx), sourceKey, destKey).Result()
		return err
	}, acceptable)
	return
}

// TSGet 获取时间序列最新数据点 ctx: 上下文, key: 时间序列键名
func (r *Client) TSGet(ctx context.Context, key string) (val redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSGet(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TSGetWithArgs 获取时间序列最新数据点（带选项） ctx: 上下文, key: 时间序列键名, options: 获取选项
func (r *Client) TSGetWithArgs(ctx context.Context, key string, options *redis.TSGetOptions) (val redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSGetWithArgs(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// TSInfo 获取时间序列信息 ctx: 上下文, key: 时间序列键名
func (r *Client) TSInfo(ctx context.Context, key string) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSInfo(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// TSInfoWithArgs 获取时间序列信息（带选项） ctx: 上下文, key: 时间序列键名, options: 信息获取选项
func (r *Client) TSInfoWithArgs(ctx context.Context, key string, options *redis.TSInfoOptions) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSInfoWithArgs(getCtx(ctx), key, options).Result()
		return err
	}, acceptable)
	return
}

// TSMAdd 批量添加时间序列数据点 ctx: 上下文, ktvSlices: 键-时间戳-值切片数组
func (r *Client) TSMAdd(ctx context.Context, ktvSlices [][]any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSMAdd(getCtx(ctx), ktvSlices).Result()
		return err
	}, acceptable)
	return
}

// TSQueryIndex 查询匹配过滤器的时间序列键名 ctx: 上下文, filterExpr: 过滤表达式数组
func (r *Client) TSQueryIndex(ctx context.Context, filterExpr []string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSQueryIndex(getCtx(ctx), filterExpr).Result()
		return err
	}, acceptable)
	return
}

// TSRevRange 反向查询时间序列数据点 ctx: 上下文, key: 时间序列键名, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳
func (r *Client) TSRevRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRevRange(getCtx(ctx), key, fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

// TSRevRangeWithArgs 反向查询时间序列数据点（带选项） ctx: 上下文, key: 时间序列键名, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, options: 查询选项
func (r *Client) TSRevRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRevRangeOptions) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRevRangeWithArgs(getCtx(ctx), key, fromTimestamp, toTimestamp, options).Result()
		return err
	}, acceptable)
	return
}

// TSRange 查询时间序列数据点 ctx: 上下文, key: 时间序列键名, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳
func (r *Client) TSRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRange(getCtx(ctx), key, fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

// TSRangeWithArgs 查询时间序列数据点（带选项） ctx: 上下文, key: 时间序列键名, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, options: 查询选项
func (r *Client) TSRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRangeOptions) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRangeWithArgs(getCtx(ctx), key, fromTimestamp, toTimestamp, options).Result()
		return err
	}, acceptable)
	return
}

// TSMRange 批量查询时间序列数据点 ctx: 上下文, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, filterExpr: 过滤表达式数组
func (r *Client) TSMRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSMRange(getCtx(ctx), fromTimestamp, toTimestamp, filterExpr).Result()
		return err
	}, acceptable)
	return
}

// TSMRangeWithArgs 批量查询时间序列数据点（带选项） ctx: 上下文, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, filterExpr: 过滤表达式数组, options: 查询选项
func (r *Client) TSMRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRangeOptions) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSMRangeWithArgs(getCtx(ctx), fromTimestamp, toTimestamp, filterExpr, options).Result()
		return err
	}, acceptable)
	return
}

// TSMRevRange 批量反向查询时间序列数据点 ctx: 上下文, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, filterExpr: 过滤表达式数组
func (r *Client) TSMRevRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.TSMRevRange(getCtx(ctx), fromTimestamp, toTimestamp, filterExpr).Result()
		return err
	}, acceptable)
	return
}

// TSMRevRangeWithArgs 批量反向查询时间序列数据点（带选项） ctx: 上下文, fromTimestamp: 起始时间戳, toTimestamp: 结束时间戳, filterExpr: 过滤表达式数组, options: 查询选项
func (r *Client) TSMRevRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *redis.TSMRevRangeOptions) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.TSMRevRangeWithArgs(getCtx(ctx), fromTimestamp, toTimestamp, filterExpr, options).Result()
		return err
	}, acceptable)
	return
}

// TSMGet 批量获取时间序列最新数据点 ctx: 上下文, filters: 过滤表达式数组
func (r *Client) TSMGet(ctx context.Context, filters []string) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.TSMGet(getCtx(ctx), filters).Result()
		return err
	}, acceptable)
	return
}

// TSMGetWithArgs 批量获取时间序列最新数据点（带选项） ctx: 上下文, filters: 过滤表达式数组, options: 获取选项
func (r *Client) TSMGetWithArgs(ctx context.Context, filters []string, options *redis.TSMGetOptions) (val map[string][]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {

			return err
		}
		val, err = conn.TSMGetWithArgs(getCtx(ctx), filters, options).Result()
		return err
	}, acceptable)
	return
}
