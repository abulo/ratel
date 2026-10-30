package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// TSAdd 向时间序列 key 追加 timestamp 时刻的 value 数据点，返回写入后的时间戳。
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

// TSAddWithArgs 按 TSOptions 向时间序列 key 追加数据点，选项可控制是否创建序列与去重策略。
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

// TSCreate 创建新的空时间序列 key；键已存在时报错。
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

// TSCreateWithArgs 按 TSOptions 创建时间序列 key，可配置保留策略、标签与压缩规则。
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

// TSAlter 修改已存在时间序列 key 的保留策略、标签等属性。
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

// TSCreateRule 将 sourceKey 按 aggregator 聚合、bucketDuration 为时间桶生成规则写入 destKey。
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

// TSCreateRuleWithArgs 按 TSCreateRuleOptions 创建聚合规则，可指定对齐与优先级等选项。
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

// TSIncrBy 将时间序列 key 在 timestamp 时刻的值累加，返回更新后的时间戳。
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

// TSIncrByWithArgs 按 TSIncrDecrOptions 累加时间序列值，选项可控制是否创建序列。
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

// TSDecrBy 将时间序列 key 在 timestamp 时刻的值递减，返回更新后的时间戳。
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

// TSDecrByWithArgs 按 TSIncrDecrOptions 递减时间序列值。
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

// TSDel 删除时间序列 key 中 [fromTimestamp, toTimestamp] 区间内的数据点，返回删除数量。
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

// TSDeleteRule 移除 sourceKey 指向 destKey 的聚合规则。
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

// TSGet 返回时间序列 key 的最新数据点（时间戳与值）。
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

// TSGetWithArgs 按 TSGetOptions 返回时间序列 key 的最新数据点。
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

// TSInfo 返回时间序列 key 的样本数量、内存占用、标签、规则等统计信息。
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

// TSInfoWithArgs 按 TSInfoOptions 返回时间序列信息，可限定返回的信息类别。
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

// TSMAdd 批量写入多组 键-时间戳-值 数据点，返回各序列写入后的时间戳。
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

// TSQueryIndex 返回满足 filterExpr 标签过滤表达式的全部时间序列键名。
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

// TSRevRange 按时间倒序返回 key 在 [fromTimestamp, toTimestamp] 区间内的数据点。
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

// TSRevRangeWithArgs 按 TSRevRangeOptions 倒序查询区间数据点，选项支持聚合、过滤与数量限制。
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

// TSRange 按时间正序返回 key 在 [fromTimestamp, toTimestamp] 区间内的数据点。
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

// TSRangeWithArgs 按 TSRangeOptions 正序查询区间数据点，选项支持聚合、过滤与数量限制。
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

// TSMRange 跨满足 filterExpr 的多个时间序列，查询 [fromTimestamp, toTimestamp] 区间数据点。
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

// TSMRangeWithArgs 按 TSMRangeOptions 跨序列查询区间数据点，选项支持聚合、标签过滤与分组。
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

// TSMRevRange 跨满足 filterExpr 的多个时间序列，按时间倒序查询区间数据点。
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

// TSMRevRangeWithArgs 按 TSMRevRangeOptions 跨序列倒序查询区间数据点。
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

// TSMGet 批量返回满足 filters 过滤条件的各时间序列最新数据点。
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

// TSMGetWithArgs 按 TSMGetOptions 批量获取各序列最新数据点。
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
