package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// XAdd 将指定的流条目追加到指定key的流中 ctx: 上下文 a: 流条目参数
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

// XDel 从指定流中移除指定的条目 ctx: 上下文 stream: 流名称 ids: 要删除的条目ID列表
func (r *Client) XDel(ctx context.Context, stream string, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XDel(getCtx(ctx), stream, ids...).Result()
		return err
	}, acceptable)
	return
}

// XLen 返回流中的条目数 ctx: 上下文 stream: 流名称
func (r *Client) XLen(ctx context.Context, stream string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XLen(getCtx(ctx), stream).Result()
		return err
	}, acceptable)
	return
}

// XRange 返回流中满足给定ID范围的条目 ctx: 上下文 stream: 流名称 start: 起始ID stop: 结束ID
func (r *Client) XRange(ctx context.Context, stream string, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRange(getCtx(ctx), stream, start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRangeN 返回流中满足给定ID范围的条目，限制返回数量 ctx: 上下文 stream: 流名称 start: 起始ID stop: 结束ID count: 最大返回数量
func (r *Client) XRangeN(ctx context.Context, stream string, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRangeN(getCtx(ctx), stream, start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRevRange 以相反顺序返回流中满足给定ID范围的条目 ctx: 上下文 stream: 流名称 start: 起始ID stop: 结束ID
func (r *Client) XRevRange(ctx context.Context, stream string, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRange(getCtx(ctx), stream, start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRevRangeN 以相反顺序返回流中满足给定ID范围的条目，限制返回数量 ctx: 上下文 stream: 流名称 start: 起始ID stop: 结束ID count: 最大返回数量
func (r *Client) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRangeN(getCtx(ctx), stream, start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRead 从一个或多个流中读取数据 ctx: 上下文 a: 读取参数
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

// XReadStreams 从多个流中读取数据 ctx: 上下文 streams: 流名称列表
func (r *Client) XReadStreams(ctx context.Context, streams ...string) (val []redis.XStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XReadStreams(getCtx(ctx), streams...).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreate 创建消费者组 ctx: 上下文 stream: 流名称 group: 消费者组名称 start: 起始ID
func (r *Client) XGroupCreate(ctx context.Context, stream string, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreate(getCtx(ctx), stream, group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreateMkStream 创建消费者组并创建流（如果不存在） ctx: 上下文 stream: 流名称 group: 消费者组名称 start: 起始ID
func (r *Client) XGroupCreateMkStream(ctx context.Context, stream string, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreateMkStream(getCtx(ctx), stream, group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupSetID 设置消费者组的最后递送ID ctx: 上下文 stream: 流名称 group: 消费者组名称 start: 新的起始ID
func (r *Client) XGroupSetID(ctx context.Context, stream string, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupSetID(getCtx(ctx), stream, group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupDestroy 删除消费者组 ctx: 上下文 stream: 流名称 group: 消费者组名称
func (r *Client) XGroupDestroy(ctx context.Context, stream string, group string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDestroy(getCtx(ctx), stream, group).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreateConsumer 创建消费者 ctx: 上下文 stream: 流名称 group: 消费者组名称 consumer: 消费者名称
func (r *Client) XGroupCreateConsumer(ctx context.Context, stream string, group, consumer string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreateConsumer(getCtx(ctx), stream, group, consumer).Result()
		return err
	}, acceptable)
	return
}

// XGroupDelConsumer 删除消费者 ctx: 上下文 stream: 流名称 group: 消费者组名称 consumer: 消费者名称
func (r *Client) XGroupDelConsumer(ctx context.Context, stream string, group, consumer string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDelConsumer(getCtx(ctx), stream, group, consumer).Result()
		return err
	}, acceptable)
	return
}

// XReadGroup 从消费者组读取数据 ctx: 上下文 a: 读取参数
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

// XAck 确认消息已被处理 ctx: 上下文 stream: 流名称 group: 消费者组名称 ids: 消息ID列表
func (r *Client) XAck(ctx context.Context, stream string, group string, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XAck(getCtx(ctx), stream, group, ids...).Result()
		return err
	}, acceptable)
	return
}

// XPending 获取待处理消息信息 ctx: 上下文 stream: 流名称 group: 消费者组名称
func (r *Client) XPending(ctx context.Context, stream string, group string) (val *redis.XPending, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XPending(getCtx(ctx), stream, group).Result()
		return err
	}, acceptable)
	return
}

// XPendingExt 获取待处理消息的详细信息 ctx: 上下文 a: 查询参数
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

// XClaim 认领待处理消息 ctx: 上下文 a: 认领参数
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

// XClaimJustID 认领待处理消息并仅返回ID ctx: 上下文 a: 认领参数
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

// XAutoClaim 自动认领待处理消息 ctx: 上下文 a: 自动认领参数
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

// XAutoClaimJustID 自动认领待处理消息并仅返回ID ctx: 上下文 a: 自动认领参数
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

// XTrimMaxLen 根据最大长度修剪流 ctx: 上下文 key: 流名称 maxLen: 最大长度
func (r *Client) XTrimMaxLen(ctx context.Context, key string, maxLen int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLen(getCtx(ctx), key, maxLen).Result()
		return err
	}, acceptable)
	return
}

// XTrimMaxLenApprox 根据最大长度近似修剪流 ctx: 上下文 key: 流名称 maxLen: 最大长度 limit: 修剪限制
func (r *Client) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLenApprox(getCtx(ctx), key, maxLen, limit).Result()
		return err
	}, acceptable)
	return
}

// XTrimMinID 根据最小ID修剪流 ctx: 上下文 key: 流名称 minID: 最小ID
func (r *Client) XTrimMinID(ctx context.Context, key string, minID string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinID(getCtx(ctx), key, minID).Result()
		return err
	}, acceptable)
	return
}

// XTrimMinIDApprox 根据最小ID近似修剪流 ctx: 上下文 key: 流名称 minID: 最小ID limit: 修剪限制
func (r *Client) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinIDApprox(getCtx(ctx), key, minID, limit).Result()
		return err
	}, acceptable)
	return
}

// XInfoGroups 获取流消费者组信息 ctx: 上下文 key: 流名称
func (r *Client) XInfoGroups(ctx context.Context, key string) (val []redis.XInfoGroup, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoGroups(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// XInfoStream 获取流信息 ctx: 上下文 key: 流名称
func (r *Client) XInfoStream(ctx context.Context, key string) (val *redis.XInfoStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStream(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// XInfoStreamFull 获取流的完整信息 ctx: 上下文 key: 流名称 count: 返回条目数量
func (r *Client) XInfoStreamFull(ctx context.Context, key string, count int) (val *redis.XInfoStreamFull, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStreamFull(getCtx(ctx), key, count).Result()
		return err
	}, acceptable)
	return
}

// XInfoConsumers 获取消费者信息 ctx: 上下文 key: 流名称 group: 消费者组名称
func (r *Client) XInfoConsumers(ctx context.Context, key string, group string) (val []redis.XInfoConsumer, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoConsumers(getCtx(ctx), key, group).Result()
		return err
	}, acceptable)
	return
}
