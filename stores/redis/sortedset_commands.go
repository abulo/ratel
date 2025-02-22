package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BZPopMax 阻塞式获取有序集合中分数最大的元素 ctx: 上下文 timeout: 超时时间 keys: 一个或多个有序集合键
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMax(getCtx(ctx), timeout, keys...).Result()
		return err
	}, acceptable)
	return
}

// BZPopMin 阻塞式获取有序集合中分数最小的元素 ctx: 上下文 timeout: 超时时间 keys: 一个或多个有序集合键
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMin(getCtx(ctx), timeout, keys...).Result()
		return err
	}, acceptable)
	return
}

// BZMPop 阻塞式从有序集合中弹出元素 ctx: 上下文 timeout: 超时时间 order: 排序方式 count: 弹出数量 keys: 一个或多个有序集合键
func (r *Client) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) (val string, ret []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, ret, err = conn.BZMPop(getCtx(ctx), timeout, order, count, keys...).Result()
		return err
	}, acceptable)
	return
}

// ZAdd 向有序集合添加元素 ctx: 上下文 key: 有序集合键 members: 要添加的元素及其分数
func (r *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	// return getRedis(r).ZAdd(getCtx(ctx), key, members...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAdd(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddLT 仅当新分数小于当前分数时更新元素 ctx: 上下文 key: 有序集合键 members: 要添加的元素及其分数
func (r *Client) ZAddLT(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddLT(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddGT 仅当新分数大于当前分数时更新元素 ctx: 上下文 key: 有序集合键 members: 要添加的元素及其分数
func (r *Client) ZAddGT(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddGT(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddNX 仅当元素不存在时添加元素 ctx: 上下文 key: 有序集合键 members: 要添加的元素及其分数
func (r *Client) ZAddNX(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddNX(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddXX 仅当元素存在时更新元素 ctx: 上下文 key: 有序集合键 members: 要添加的元素及其分数
func (r *Client) ZAddXX(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddXX(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddArgs 使用ZAddArgs参数向有序集合添加元素 ctx: 上下文 key: 有序集合键 args: ZAdd参数
func (r *Client) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgs(getCtx(ctx), key, args).Result()
		return err
	}, acceptable)
	return
}

// ZAddArgsIncr 使用ZAddArgs参数增加元素的分数 ctx: 上下文 key: 有序集合键 args: ZAdd参数
func (r *Client) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgsIncr(getCtx(ctx), key, args).Result()
		return err
	}, acceptable)
	return
}

// ZCard 获取有序集合的元素数量 ctx: 上下文 key: 有序集合键
func (r *Client) ZCard(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCard(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ZCount 统计分数在指定范围内的元素数量 ctx: 上下文 key: 有序集合键 min: 最小分数 max: 最大分数
func (r *Client) ZCount(ctx context.Context, key string, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCount(getCtx(ctx), key, min, max).Result()
		return err
	}, acceptable)
	return
}

// ZLexCount 统计字典序在指定范围内的元素数量 ctx: 上下文 key: 有序集合键 min: 最小字典序 max: 最大字典序
func (r *Client) ZLexCount(ctx context.Context, key string, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZLexCount(getCtx(ctx), key, min, max).Result()
		return err
	}, acceptable)
	return
}

// ZIncrBy 增加有序集合中元素的分数 ctx: 上下文 key: 有序集合键 increment: 增量 member: 元素
func (r *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZIncrBy(getCtx(ctx), key, increment, member).Result()
		return err
	}, acceptable)
	return
}

// ZInter 计算多个有序集合的交集 ctx: 上下文 store: ZStore参数
func (r *Client) ZInter(ctx context.Context, store *redis.ZStore) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInter(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}

// ZInterWithScores 计算多个有序集合的交集并返回分数 ctx: 上下文 store: ZStore参数
func (r *Client) ZInterWithScores(ctx context.Context, store *redis.ZStore) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterWithScores(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}

// ZInterCard 计算多个有序集合的交集基数 ctx: 上下文 limit: 最大返回数量 keys: 一个或多个有序集合键
func (r *Client) ZInterCard(ctx context.Context, limit int64, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterCard(getCtx(ctx), limit, keys...).Result()
		return err
	}, acceptable)
	return
}

// ZInterStore 计算多个有序集合的交集并存储结果 ctx: 上下文 key: 目标有序集合键 store: ZStore参数
func (r *Client) ZInterStore(ctx context.Context, key string, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterStore(getCtx(ctx), key, store).Result()
		return err
	}, acceptable)
	return
}

// ZMPop 从有序集合中弹出元素 ctx: 上下文 order: 排序方式 count: 弹出数量 keys: 一个或多个有序集合键
func (r *Client) ZMPop(ctx context.Context, order string, count int64, keys ...string) (key string, ret []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		key, ret, err = conn.ZMPop(getCtx(ctx), order, count, keys...).Result()
		return err
	}, acceptable)
	return
}

// ZMScore 获取多个元素的分数 ctx: 上下文 key: 有序集合键 members: 要查询的元素
func (r *Client) ZMScore(ctx context.Context, key string, members ...string) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZMScore(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMax 删除并返回分数最高的元素 ctx: 上下文 key: 有序集合键 count: 要删除的元素数量
func (r *Client) ZPopMax(ctx context.Context, key string, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMax(getCtx(ctx), key, count...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMin 删除并返回分数最低的元素 ctx: 上下文 key: 有序集合键 count: 要删除的元素数量
func (r *Client) ZPopMin(ctx context.Context, key string, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMin(getCtx(ctx), key, count...).Result()
		return err
	}, acceptable)
	return
}

// ZRange 获取有序集合指定范围内的元素 ctx: 上下文 key: 有序集合键 start: 起始索引 stop: 结束索引
func (r *Client) ZRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRange(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeWithScores 获取有序集合指定范围内的元素及其分数 ctx: 上下文 key: 有序集合键 start: 起始索引 stop: 结束索引
func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeWithScores(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScore 获取分数在指定范围内的元素 ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScore(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByLex 获取字典序在指定范围内的元素 ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByLex(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScoreWithScores 获取分数在指定范围内的元素及其分数 ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScoreWithScores(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeArgs 使用ZRangeArgs参数获取有序集合元素 ctx: 上下文 z: ZRangeArgs参数
func (r *Client) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeArgs(getCtx(ctx), z).Result()
		return err
	}, acceptable)
	return
}

// ZRangeArgsWithScores 使用ZRangeArgs参数获取有序集合元素及其分数 ctx: 上下文 z: ZRangeArgs参数
func (r *Client) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeArgsWithScores(getCtx(ctx), z).Result()
		return err
	}, acceptable)
	return
}

// ZRangeStore 使用ZRangeArgs参数获取有序集合元素并存储到目标集合 ctx: 上下文 dst: 目标集合键 z: ZRangeArgs参数
func (r *Client) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeStore(getCtx(ctx), dst, z).Result()
		return err
	}, acceptable)
	return

}

// ZRank 获取元素的排名 ctx: 上下文 key: 有序集合键 member: 元素
func (r *Client) ZRank(ctx context.Context, key string, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRank(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// ZRankWithScore 获取元素的排名及其分数 ctx: 上下文 key: 有序集合键 member: 元素
func (r *Client) ZRankWithScore(ctx context.Context, key string, member string) (val redis.RankScore, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRankWithScore(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// ZRem 移除有序集合中的元素 ctx: 上下文 key: 有序集合键 members: 要移除的元素
func (r *Client) ZRem(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRem(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByRank 移除指定排名范围内的元素 ctx: 上下文 key: 有序集合键 start: 起始排名 stop: 结束排名
func (r *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByRank(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByScore 移除指定分数范围内的元素 ctx: 上下文 key: 有序集合键 min: 最小分数 max: 最大分数
func (r *Client) ZRemRangeByScore(ctx context.Context, key string, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByScore(getCtx(ctx), key, min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByLex 移除指定字典序范围内的元素 ctx: 上下文 key: 有序集合键 min: 最小字典序 max: 最大字典序
func (r *Client) ZRemRangeByLex(ctx context.Context, key string, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByLex(getCtx(ctx), key, min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRevRange 获取指定范围内的元素（倒序） ctx: 上下文 key: 有序集合键 start: 起始索引 stop: 结束索引
func (r *Client) ZRevRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRange(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeWithScores 获取指定范围内的元素及其分数（倒序） ctx: 上下文 key: 有序集合键 start: 起始索引 stop: 结束索引
func (r *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeWithScores(getCtx(ctx), key, start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScore 获取指定分数范围内的元素（倒序） ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScore(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByLex 获取指定字典序范围内的元素（倒序） ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByLex(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScoreWithScores 获取指定分数范围内的元素及其分数（倒序） ctx: 上下文 key: 有序集合键 opt: ZRangeBy参数
func (r *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScoreWithScores(getCtx(ctx), key, opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRank 获取元素的倒序排名 ctx: 上下文 key: 有序集合键 member: 元素
func (r *Client) ZRevRank(ctx context.Context, key string, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRank(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// ZRevRankWithScore 获取元素的倒序排名及其分数 ctx: 上下文 key: 有序集合键 member: 元素
func (r *Client) ZRevRankWithScore(ctx context.Context, key string, member string) (val redis.RankScore, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRankWithScore(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// ZScore 获取元素的分数 ctx: 上下文 key: 有序集合键 member: 元素
func (r *Client) ZScore(ctx context.Context, key string, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZScore(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// ZUnionStore 计算多个有序集合的并集并存储结果 ctx: 上下文 dest: 目标集合键 store: ZStore参数
func (r *Client) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnionStore(getCtx(ctx), dest, store).Result()
		return err
	}, acceptable)
	return
}

// ZRandMember 随机获取有序集合中的元素 ctx: 上下文 key: 有序集合键 count: 要获取的元素数量
func (r *Client) ZRandMember(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMember(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// ZRandMemberWithScores 随机获取有序集合中的元素及其分数 ctx: 上下文 key: 有序集合键 count: 要获取的元素数量
func (r *Client) ZRandMemberWithScores(ctx context.Context, key string, count int) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMemberWithScores(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// ZUnion 计算多个有序集合的并集 ctx: 上下文 store: ZStore参数
func (r *Client) ZUnion(ctx context.Context, store redis.ZStore) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnion(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}

// ZUnionWithScores 计算多个有序集合的并集并返回分数 ctx: 上下文 store: ZStore参数
func (r *Client) ZUnionWithScores(ctx context.Context, store redis.ZStore) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnionWithScores(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}

// ZDiff 计算多个有序集合的差集 ctx: 上下文 keys: 一个或多个有序集合键
func (r *Client) ZDiff(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiff(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// ZDiffWithScores 计算多个有序集合的差集并返回分数 ctx: 上下文 keys: 一个或多个有序集合键
func (r *Client) ZDiffWithScores(ctx context.Context, keys ...string) (val []redis.Z, err error) {
	// return getRedis(r).ZDiffWithScores(getCtx(ctx), keys...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffWithScores(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// ZDiffStore 计算多个有序集合的差集并存储结果 ctx: 上下文 destination: 目标集合键 keys: 一个或多个有序集合键
func (r *Client) ZDiffStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffStore(getCtx(ctx), destination, keys...).Result()
		return err
	}, acceptable)
	return
}

// ZScan 扫描有序集合中的元素 ctx: 上下文 key: 有序集合键 cursorIn: 游标 match: 匹配模式 count: 返回数量
func (r *Client) ZScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.ZScan(getCtx(ctx), key, cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}
