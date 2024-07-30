package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SAdd 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
// 当 key 不是集合类型时，返回一个错误。
func (r *Client) SAdd(ctx context.Context, key any, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SAdd(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// SCard 返回集合 key 的基数(集合中元素的数量)。
func (r *Client) SCard(ctx context.Context, key any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SCard(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SDiff 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (r *Client) SDiff(ctx context.Context, keys ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiff(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SDiffStore 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SDiffStore(ctx context.Context, destination any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiffStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInter 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
// 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func (r *Client) SInter(ctx context.Context, keys ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInter(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterCard(ctx context.Context, limit int64, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterCard(getCtx(ctx), limit, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterStore(ctx context.Context, destination any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SIsMember 判断 member 元素是否集合 key 的成员。
func (r *Client) SIsMember(ctx context.Context, key any, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SIsMember(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) SMIsMember(ctx context.Context, key string, members ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMIsMember(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// SMembers 返回集合 key 中的所有成员。
// 不存在的 key 被视为空集合。
func (r *Client) SMembers(ctx context.Context, key any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembers(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SMembersMap -> SMembers
func (r *Client) SMembersMap(ctx context.Context, key any) (val map[string]struct{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembersMap(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SMove 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
func (r *Client) SMove(ctx context.Context, source, destination any, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMove(getCtx(ctx), r.k(source), r.k(destination), member).Result()
		return err
	}, acceptable)
	return
}

// SPop 移除并返回集合中的一个随机元素。
// 如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
func (r *Client) SPop(ctx context.Context, key any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPop(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SPopN -> SPop
func (r *Client) SPopN(ctx context.Context, key any, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPopN(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// SRandMember 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， SRANDMEMBER 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 SRANDMEMBER 则仅仅返回随机元素，而不对集合进行任何改动。
func (r *Client) SRandMember(ctx context.Context, key any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMember(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SRandMemberN -> SRandMember
func (r *Client) SRandMemberN(ctx context.Context, key any, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMemberN(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// SRem 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
func (r *Client) SRem(ctx context.Context, key any, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRem(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// SScan 详细信息请参考 SCAN 命令。
func (r *Client) SScan(ctx context.Context, key any, cursorIn uint64, match string, count int64) (val *redis.ScanCmd, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.SScan(getCtx(ctx), r.k(key), cursorIn, match, count)
		return nil
	}, acceptable)
	return
}

// SUnion 返回一个集合的全部成员，该集合是所有给定集合的并集。
// 不存在的 key 被视为空集。
func (r *Client) SUnion(ctx context.Context, keys ...any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnion(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SUnionStore 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SUnionStore(ctx context.Context, destination any, keys ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnionStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
