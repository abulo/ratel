package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SAdd 向集合 key 添加 members，返回本次新增的元素数量；已存在的元素被忽略。
func (r *Client) SAdd(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SAdd(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// SCard 返回集合 key 的元素数量；键不存在时返回 0。
func (r *Client) SCard(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SCard(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// SDiff 返回给定多个集合之间的差集，即第一个集合中有而其余集合中没有的元素。
func (r *Client) SDiff(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiff(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// SDiffStore 将给定集合的差集写入 destination 集合，返回结果中的元素数量。
func (r *Client) SDiffStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiffStore(getCtx(ctx), destination, keys...).Result()
		return err
	}, acceptable)
	return
}

// SInter 返回给定多个集合的交集元素。
func (r *Client) SInter(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInter(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// SInterCard 返回给定多个集合交集的元素数量，limit 为参与计算的集合基数上限。
func (r *Client) SInterCard(ctx context.Context, limit int64, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterCard(getCtx(ctx), limit, keys...).Result()
		return err
	}, acceptable)
	return
}

// SInterStore 将给定集合的交集写入 destination 集合，返回结果中的元素数量。
func (r *Client) SInterStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterStore(getCtx(ctx), destination, keys...).Result()
		return err
	}, acceptable)
	return
}

// SIsMember 判断 member 是否属于集合 key，属于返回 true。
func (r *Client) SIsMember(ctx context.Context, key string, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SIsMember(getCtx(ctx), key, member).Result()
		return err
	}, acceptable)
	return
}

// SMIsMember 批量判断 members 是否属于集合 key，返回的布尔切片顺序与 members 顺序一致。
func (r *Client) SMIsMember(ctx context.Context, key string, members ...any) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMIsMember(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// SMembers 返回集合 key 的全部元素。
func (r *Client) SMembers(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembers(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// SMembersMap 返回集合全部元素组成的 map（值为空结构体），便于调用方做 O(1) 成员判断。
func (r *Client) SMembersMap(ctx context.Context, key string) (val map[string]struct{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembersMap(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// SMove 将 member 从 source 集合移动到 destination 集合，移动成功返回 true。
func (r *Client) SMove(ctx context.Context, source, destination string, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMove(getCtx(ctx), source, destination, member).Result()
		return err
	}, acceptable)
	return
}

// SPop 随机移除并返回集合 key 的一个元素；集合为空或不存在时返回 redis.Nil。
func (r *Client) SPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPop(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// SPopN 随机移除并返回集合 key 中最多 count 个元素。
func (r *Client) SPopN(ctx context.Context, key string, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPopN(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// SRandMember 随机返回集合 key 的一个元素，不移除；集合为空时返回 redis.Nil。
func (r *Client) SRandMember(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMember(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// SRandMemberN 随机返回集合 key 中 count 个元素，不移除成员。
func (r *Client) SRandMemberN(ctx context.Context, key string, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMemberN(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// SRem 从集合 key 中移除 members，返回实际移除的元素数量；不存在的成员被忽略。
func (r *Client) SRem(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRem(getCtx(ctx), key, members...).Result()
		return err
	}, acceptable)
	return
}

// SScan 以游标 cursorIn 迭代集合元素；match 为成员匹配模式，count 为建议批次大小，
// 返回本次元素与供下次调用使用的新游标。
func (r *Client) SScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val *redis.ScanCmd, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.SScan(getCtx(ctx), key, cursorIn, match, count)
		return nil
	}, acceptable)
	return
}

// SUnion 返回给定多个集合的并集元素。
func (r *Client) SUnion(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnion(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// SUnionStore 将给定集合的并集写入 destination 集合，返回结果中的元素数量。
func (r *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnionStore(getCtx(ctx), destination, keys...).Result()
		return err
	}, acceptable)
	return
}
