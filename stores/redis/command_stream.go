package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// XAdd 将指定的流条目追加到指定key的流中。 如果key不存在，作为运行这个命令的副作用，将使用流的条目自动创建key。
func (r *Client) XAdd(ctx context.Context, a *redis.XAddArgs) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XAdd(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

// XDel 从指定流中移除指定的条目，并返回成功删除的条目的数量，在传递的ID不存在的情况下， 返回的数量可能与传递的ID数量不同。
func (r *Client) XDel(ctx context.Context, stream any, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XDel(getCtx(ctx), r.k(stream), ids...).Result()
		return err
	}, acceptable)
	return
}

// XLen 返回流中的条目数。如果指定的key不存在，则此命令返回0，就好像该流为空。 但是请注意，与其他的Redis类型不同，零长度流是可能的，所以你应该调用TYPE 或者 EXISTS 来检查一个key是否存在。
// 一旦内部没有任何的条目（例如调用XDEL后），流不会被自动删除，因为可能还存在与其相关联的消费者组。
func (r *Client) XLen(ctx context.Context, stream any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XLen(getCtx(ctx), r.k(stream)).Result()
		return err
	}, acceptable)
	return
}

// XRange 此命令返回流中满足给定ID范围的条目。范围由最小和最大ID指定。所有ID在指定的两个ID之间或与其中一个ID相等（闭合区间）的条目将会被返回。
func (r *Client) XRange(ctx context.Context, stream any, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRange(getCtx(ctx), r.k(stream), start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRangeN -> XRange
func (r *Client) XRangeN(ctx context.Context, stream any, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRangeN(getCtx(ctx), r.k(stream), start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRevRange 此命令与XRANGE完全相同，但显著的区别是以相反的顺序返回条目，并以相反的顺序获取开始-结束参数：在XREVRANGE中，你需要先指定结束ID，再指定开始ID，该命令就会从结束ID侧开始生成两个ID之间（或完全相同）的所有元素。
func (r *Client) XRevRange(ctx context.Context, stream any, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRange(getCtx(ctx), r.k(stream), start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRevRangeN -> XRevRange
func (r *Client) XRevRangeN(ctx context.Context, stream any, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRangeN(getCtx(ctx), r.k(stream), start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRead 从一个或者多个流中读取数据，仅返回ID大于调用者报告的最后接收ID的条目。此命令有一个阻塞选项，用于等待可用的项目，类似于BRPOP或者BZPOPMIN等等。
func (r *Client) XRead(ctx context.Context, a *redis.XReadArgs) (val []redis.XStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRead(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

// XReadStreams -> XRead
func (r *Client) XReadStreams(ctx context.Context, streams ...any) (val []redis.XStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XReadStreams(getCtx(ctx), r.ks(streams...)...).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreate command
func (r *Client) XGroupCreate(ctx context.Context, stream any, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreate(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreateMkStream command
func (r *Client) XGroupCreateMkStream(ctx context.Context, stream any, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreateMkStream(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupSetID command
func (r *Client) XGroupSetID(ctx context.Context, stream any, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupSetID(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupDestroy command
func (r *Client) XGroupDestroy(ctx context.Context, stream any, group string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDestroy(getCtx(ctx), r.k(stream), group).Result()
		return err
	}, acceptable)
	return
}

// XGroupDelConsumer command
func (r *Client) XGroupDelConsumer(ctx context.Context, stream any, group, consumer string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDelConsumer(getCtx(ctx), r.k(stream), group, consumer).Result()
		return err
	}, acceptable)
	return
}

// XReadGroup command
func (r *Client) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) (val []redis.XStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XReadGroup(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

// XAck command
func (r *Client) XAck(ctx context.Context, stream any, group string, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XAck(getCtx(ctx), r.k(stream), group, ids...).Result()
		return err
	}, acceptable)
	return
}

// XPending command
func (r *Client) XPending(ctx context.Context, stream any, group string) (val *redis.XPending, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XPending(getCtx(ctx), r.k(stream), group).Result()
		return err
	}, acceptable)
	return
}

// XPendingExt command
func (r *Client) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) (val []redis.XPendingExt, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XPendingExt(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

// XClaim command
func (r *Client) XClaim(ctx context.Context, a *redis.XClaimArgs) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XClaim(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

// XClaimJustID command
func (r *Client) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XClaimJustID(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XAutoClaim(ctx context.Context, a *redis.XAutoClaimArgs) (val []redis.XMessage, start string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, start, err = conn.XAutoClaim(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XAutoClaimJustID(ctx context.Context, a *redis.XAutoClaimArgs) (val []string, start string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, start, err = conn.XAutoClaimJustID(getCtx(ctx), a).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMaxLen(ctx context.Context, key any, maxLen int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLen(getCtx(ctx), r.k(key), maxLen).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMaxLenApprox(ctx context.Context, key any, maxLen, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLenApprox(getCtx(ctx), r.k(key), maxLen, limit).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMinID(ctx context.Context, key any, minID string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinID(getCtx(ctx), r.k(key), minID).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMinIDApprox(ctx context.Context, key any, minID string, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinIDApprox(getCtx(ctx), r.k(key), minID, limit).Result()
		return err
	}, acceptable)
	return
}

// XInfoGroups command
func (r *Client) XInfoGroups(ctx context.Context, key any) (val []redis.XInfoGroup, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoGroups(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoStream(ctx context.Context, key any) (val *redis.XInfoStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStream(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoStreamFull(ctx context.Context, key any, count int) (val *redis.XInfoStreamFull, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStreamFull(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoConsumers(ctx context.Context, key any, group string) (val []redis.XInfoConsumer, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoConsumers(getCtx(ctx), r.k(key), group).Result()
		return err
	}, acceptable)
	return
}
