package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *Client) TSAdd(ctx context.Context, key string, timestamp any, value float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAdd(getCtx(ctx), r.k(key), timestamp, value).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) TSAddWithArgs(ctx context.Context, key string, timestamp any, value float64, options *redis.TSOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAddWithArgs(getCtx(ctx), r.k(key), timestamp, value, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSCreate(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreate(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSCreateWithArgs(ctx context.Context, key string, options *redis.TSOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSCreateWithArgs(getCtx(ctx), r.k(key), options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSAlter(ctx context.Context, key string, options *redis.TSAlterOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSAlter(getCtx(ctx), r.k(key), options).Result()
		return err
	}, acceptable)
	return
}

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

func (r *Client) TSIncrBy(ctx context.Context, key string, timestamp float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSIncrBy(getCtx(ctx), r.k(key), timestamp).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSIncrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSIncrByWithArgs(getCtx(ctx), r.k(key), timestamp, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSDecrBy(ctx context.Context, key string, timestamp float64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSDecrBy(getCtx(ctx), r.k(key), timestamp).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSDecrByWithArgs(ctx context.Context, key string, timestamp float64, options *redis.TSIncrDecrOptions) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}

		val, err = conn.TSDecrByWithArgs(getCtx(ctx), r.k(key), timestamp, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSDel(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSDel(getCtx(ctx), r.k(key), fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

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

func (r *Client) TSGet(ctx context.Context, key string) (val redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSGet(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSGetWithArgs(ctx context.Context, key string, options *redis.TSGetOptions) (val redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSGetWithArgs(getCtx(ctx), r.k(key), options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSInfo(ctx context.Context, key string) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSInfo(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSInfoWithArgs(ctx context.Context, key string, options *redis.TSInfoOptions) (val map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSInfoWithArgs(getCtx(ctx), r.k(key), options).Result()
		return err
	}, acceptable)
	return
}

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

func (r *Client) TSRevRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRevRange(getCtx(ctx), r.k(key), fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSRevRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRevRangeOptions) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRevRangeWithArgs(getCtx(ctx), r.k(key), fromTimestamp, toTimestamp, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRange(getCtx(ctx), r.k(key), fromTimestamp, toTimestamp).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TSRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *redis.TSRangeOptions) (val []redis.TSTimestampValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TSRangeWithArgs(getCtx(ctx), r.k(key), fromTimestamp, toTimestamp, options).Result()
		return err
	}, acceptable)
	return
}

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
