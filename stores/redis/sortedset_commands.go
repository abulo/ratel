package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// BZPopMax 阻塞式弹出有序集合中分数最高的成员，返回成员、分数及来源键；超时仍无成员时返回 redis.Nil。
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

// BZPopMin 阻塞式弹出有序集合中分数最低的成员，返回成员、分数及来源键；超时仍无成员时返回 redis.Nil。
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

// BZMPop 阻塞式按 order（MAX/MIN）从多个有序集合弹出 count 个成员，返回来源键与成员分数；超时返回 redis.Nil。
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

// ZAdd 向有序集合 key 添加 members（成员与分数），返回新增成员数量；已存在成员更新为给定分数。
func (r *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
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

// ZAddLT 仅当新分数小于当前分数时更新成员分数；成员不存在时新增。
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

// ZAddGT 仅当新分数大于当前分数时更新成员分数；成员不存在时新增。
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

// ZAddNX 仅当成员不存在时添加；已存在成员不更新分数。
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

// ZAddXX 仅当成员已存在时更新分数；不存在成员不添加。
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

// ZAddArgs 按 ZAddArgs 指定的条件（NX/XX、LT/GT、CH、INCR 等）向有序集合添加或更新成员。
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

// ZAddArgsIncr 按 ZAddArgs 对成员分数执行增量更新，返回更新后的分数。
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

// ZCard 返回有序集合 key 的成员数量；键不存在时返回 0。
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

// ZCount 返回分数落在 [min, max] 区间内的成员数量；min、max 可加 "(" 前缀表示开区间。
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

// ZLexCount 返回字典序落在 [min, max] 区间内的成员数量。
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

// ZIncrBy 为有序集合成员 member 的分数增加 increment，返回更新后的分数；成员不存在时先以 0 为基准。
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

// ZInter 计算多个有序集合的交集，返回成员名称；成员权重与聚合方式由 ZStore 指定。
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

// ZInterWithScores 计算多个有序集合的交集，并连同聚合分数一起返回。
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

// ZInterCard 返回多个有序集合交集的成员数量，limit 为参与计算的基数上限。
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

// ZInterStore 将多个有序集合的交集结果存入目标 key，返回结果成员数量。
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

// ZMPop 按 order（MAX/MIN）从多个有序集合弹出 count 个成员，返回来源键与成员分数。
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

// ZMScore 返回给定多个成员的分数切片，顺序与 members 一致；不存在成员对应位置为 nil。
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

// ZPopMax 弹出分数最高的成员，count 可选指定弹出数量。
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

// ZPopMin 弹出分数最低的成员，count 可选指定弹出数量。
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

// ZRange 按索引区间 [start, stop] 升序返回成员名称。
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

// ZRangeWithScores 按索引区间 [start, stop] 升序返回成员及其分数。
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

// ZRangeByScore 按 ZRangeBy 指定的分数区间升序返回成员名称。
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

// ZRangeByLex 按 ZRangeBy 指定的字典序区间升序返回成员名称。
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

// ZRangeByScoreWithScores 按分数区间升序返回成员及其分数。
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

// ZRangeArgs 按 ZRangeArgs 指定的条件返回成员名称。
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

// ZRangeArgsWithScores 按 ZRangeArgs 指定的条件返回成员及其分数。
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

// ZRangeStore 按 ZRangeArgs 取成员并存入目标集合 dst，返回写入的成员数量。
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

// ZRank 返回成员 member 按分数升序的排名（从 0 开始）；成员不存在时返回 redis.Nil。
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

// ZRankWithScore 返回成员 member 的升序排名及其分数。
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

// ZRem 移除有序集合中的 members，返回实际移除的成员数量。
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

// ZRemRangeByRank 移除排名落在 [start, stop] 区间内的成员，返回移除数量。
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

// ZRemRangeByScore 移除分数落在 [min, max] 区间内的成员，返回移除数量。
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

// ZRemRangeByLex 移除字典序落在 [min, max] 区间内的成员，返回移除数量。
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

// ZRevRange 按索引区间 [start, stop] 降序返回成员名称。
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

// ZRevRangeWithScores 按索引区间 [start, stop] 降序返回成员及其分数。
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

// ZRevRangeByScore 按 ZRangeBy 指定的分数区间降序返回成员名称。
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

// ZRevRangeByLex 按 ZRangeBy 指定的字典序区间降序返回成员名称。
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

// ZRevRangeByScoreWithScores 按分数区间降序返回成员及其分数。
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

// ZRevRank 返回成员 member 按分数降序的排名（从 0 开始）；成员不存在时返回 redis.Nil。
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

// ZRevRankWithScore 返回成员 member 的降序排名及其分数。
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

// ZScore 返回成员 member 的分数；成员不存在时返回 redis.Nil。
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

// ZUnionStore 将多个有序集合的并集存入 dest，返回结果成员数量；权重与聚合方式由 ZStore 指定。
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

// ZRandMember 随机返回有序集合中 count 个成员名称，不移除成员。
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

// ZRandMemberWithScores 随机返回有序集合中 count 个成员及其分数，不移除成员。
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

// ZUnion 计算多个有序集合的并集，返回成员名称；权重与聚合方式由 ZStore 指定。
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

// ZUnionWithScores 计算多个有序集合的并集，并连同聚合分数一起返回。
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

// ZDiff 计算多个有序集合的差集（第一个集合相对其余集合），返回成员名称。
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

// ZDiffWithScores 计算多个有序集合的差集，并连同分数一起返回。
func (r *Client) ZDiffWithScores(ctx context.Context, keys ...string) (val []redis.Z, err error) {
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

// ZDiffStore 将多个有序集合的差集结果存入 destination，返回结果成员数量。
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

// ZScan 以游标 cursorIn 迭代有序集合成员；match 为成员匹配模式，count 为建议批次大小，
// 返回本次成员与供下次调用使用的新游标。
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
