package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SAdd 向集合添加元素 key:集合键名 members:要添加的元素
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

// SCard 获取集合元素数量 key:集合键名
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

// SDiff 获取集合差集 keys:要比较的集合键名
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

// SDiffStore 获取并存储集合差集 destination:目标集合键名 keys:要比较的集合键名
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

// SInter 获取集合交集 keys:要比较的集合键名
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

// SInterCard 获取集合交集元素数量 limit:最大返回数量 keys:要比较的集合键名
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

// SInterStore 获取并存储集合交集 destination:目标集合键名 keys:要比较的集合键名
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

// SIsMember 判断元素是否在集合中 key:集合键名 member:要判断的元素
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

// SMIsMember 批量判断元素是否在集合中 key:集合键名 members:要判断的元素列表
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

// SMembers 获取集合所有元素 key:集合键名
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

// SMembersMap 获取集合所有元素到map key:集合键名
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

// SMove 移动集合元素 source:源集合键名 destination:目标集合键名 member:要移动的元素
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

// SPop 移除并返回随机元素 key:集合键名
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

// SPopN 移除并返回多个随机元素 key:集合键名 count:要移除的元素数量
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

// SRandMember 获取随机元素 key:集合键名
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

// SRandMemberN 获取多个随机元素 key:集合键名 count:要获取的元素数量
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

// SRem 移除集合元素 key:集合键名 members:要移除的元素
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

// SScan 扫描集合元素 key:集合键名 cursorIn:游标 match:匹配模式 count:每次扫描数量
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

// SUnion 获取集合并集 keys:要合并的集合键名
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

// SUnionStore 获取并存储集合并集 destination:目标集合键名 keys:要合并的集合键名
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
