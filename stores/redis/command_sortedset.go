package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BZPopMax 是有序集合命令 ZPOPMAX带有阻塞功能的版本。
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...any) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMax(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BZPopMin 是有序集合命令 ZPOPMIN带有阻塞功能的版本。
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...any) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMin(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...any) (val string, ret []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, ret, err = conn.BZMPop(getCtx(ctx), timeout, order, count, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
// 如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。
// score 值可以是整数值或双精度浮点数。
// 如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZAdd(ctx context.Context, key any, members ...redis.Z) (val int64, err error) {
	// return getRedis(r).ZAdd(getCtx(ctx), r.k(key), members...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAdd(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZAddLT(ctx context.Context, key any, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddLT(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZAddGT(ctx context.Context, key any, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddGT(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddNX -> ZAdd
func (r *Client) ZAddNX(ctx context.Context, key any, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddNX(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddXX -> ZAdd
func (r *Client) ZAddXX(ctx context.Context, key any, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddXX(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZAddArgs(ctx context.Context, key any, args redis.ZAddArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgs(getCtx(ctx), r.k(key), args).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZAddArgsIncr(ctx context.Context, key any, args redis.ZAddArgs) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgsIncr(getCtx(ctx), r.k(key), args).Result()
		return err
	}, acceptable)
	return
}

// ZCard 返回有序集 key 的基数。
func (r *Client) ZCard(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCard(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 关于参数 min 和 max 的详细使用方法，请参考 ZRANGEBYSCORE 命令。
func (r *Client) ZCount(ctx context.Context, key any, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCount(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZLexCount -> ZCount
func (r *Client) ZLexCount(ctx context.Context, key any, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZLexCount(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
// 当 key 不是有序集类型时，返回一个错误。
// score 值可以是整数值或双精度浮点数。
func (r *Client) ZIncrBy(ctx context.Context, key any, increment float64, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZIncrBy(getCtx(ctx), r.k(key), increment, member).Result()
		return err
	}, acceptable)
	return
}

// ZInter 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
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
func (r *Client) ZInterCard(ctx context.Context, limit int64, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterCard(getCtx(ctx), limit, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ZInterStore 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
func (r *Client) ZInterStore(ctx context.Context, key any, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterStore(getCtx(ctx), r.k(key), store).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZMPop(ctx context.Context, order string, count int64, keys ...any) (key string, ret []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		key, ret, err = conn.ZMPop(getCtx(ctx), order, count, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZMScore(ctx context.Context, key any, members ...string) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZMScore(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMax 删除并返回有序集合key中的最多count个具有最高得分的成员。
func (r *Client) ZPopMax(ctx context.Context, key any, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMax(getCtx(ctx), r.k(key), count...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMin 删除并返回有序集合key中的最多count个具有最低得分的成员。
func (r *Client) ZPopMin(ctx context.Context, key any, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMin(getCtx(ctx), r.k(key), count...).Result()
		return err
	}, acceptable)
	return
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (r *Client) ZRange(ctx context.Context, key any, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRange(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeWithScores -> ZRange
func (r *Client) ZRangeWithScores(ctx context.Context, key any, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeWithScores(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
func (r *Client) ZRangeByScore(ctx context.Context, key any, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScore(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByLex -> ZRangeByScore
func (r *Client) ZRangeByLex(ctx context.Context, key any, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByLex(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScoreWithScores -> ZRangeByScore
func (r *Client) ZRangeByScoreWithScores(ctx context.Context, key any, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScoreWithScores(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}
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
func (r *Client) ZRangeStore(ctx context.Context, dst any, z redis.ZRangeArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeStore(getCtx(ctx), r.k(dst), z).Result()
		return err
	}, acceptable)
	return

}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
// 使用 ZREVRANK 命令可以获得成员按 score 值递减(从大到小)排列的排名。
func (r *Client) ZRank(ctx context.Context, key any, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRank(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRankWithScore(ctx context.Context, key any, member string) (val redis.RankScore, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRankWithScore(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZRem(ctx context.Context, key any, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRem(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
func (r *Client) ZRemRangeByRank(ctx context.Context, key any, start, stop int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByRank(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByScore 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 自版本2.1.6开始， score 值等于 min 或 max 的成员也可以不包括在内，详情请参见 ZRANGEBYSCORE 命令。
func (r *Client) ZRemRangeByScore(ctx context.Context, key any, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByScore(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByLex -> ZRemRangeByScore
func (r *Client) ZRemRangeByLex(ctx context.Context, key any, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByLex(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func (r *Client) ZRevRange(ctx context.Context, key any, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRange(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeWithScores -> ZRevRange
func (r *Client) ZRevRangeWithScores(ctx context.Context, key any, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeWithScores(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScore 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func (r *Client) ZRevRangeByScore(ctx context.Context, key any, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScore(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByLex -> ZRevRangeByScore
func (r *Client) ZRevRangeByLex(ctx context.Context, key any, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByLex(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScoreWithScores -> ZRevRangeByScore
func (r *Client) ZRevRangeByScoreWithScores(ctx context.Context, key any, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScoreWithScores(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (r *Client) ZRevRank(ctx context.Context, key any, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRank(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRevRankWithScore(ctx context.Context, key any, member string) (val redis.RankScore, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRankWithScore(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *Client) ZScore(ctx context.Context, key any, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZScore(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZUnionStore 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之 和 。
func (r *Client) ZUnionStore(ctx context.Context, dest any, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnionStore(getCtx(ctx), r.k(dest), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRandMember(ctx context.Context, key any, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMember(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRandMemberWithScores(ctx context.Context, key any, count int) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMemberWithScores(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
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
func (r *Client) ZDiff(ctx context.Context, keys ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiff(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZDiffWithScores(ctx context.Context, keys ...any) (val []redis.Z, err error) {
	// return getRedis(r).ZDiffWithScores(getCtx(ctx), r.ks(keys...)...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffWithScores(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZDiffStore(ctx context.Context, destination any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ZScan 详细信息请参考 SCAN 命令。
func (r *Client) ZScan(ctx context.Context, key any, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.ZScan(getCtx(ctx), r.k(key), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}
