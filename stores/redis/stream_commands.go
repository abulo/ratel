package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// XAdd 按 XAddArgs 将一条消息追加到流，返回新生成的消息 ID。
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

// XDel 从流 stream 中删除给定 ids 的消息，返回实际删除的条目数量。
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

// XLen 返回流 stream 中的消息条目数量。
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

// XRange 返回流中 ID 落在 [start, stop] 区间内的消息，按 ID 升序排列。
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

// XRangeN 同 XRange，最多返回 count 条消息。
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

// XRevRange 以 ID 降序返回流中 [start, stop] 区间内的消息。
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

// XRevRangeN 以 ID 降序返回区间内消息，最多返回 count 条。
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

// XRead 按 XReadArgs 从一个或多个流读取新消息，可按参数配置阻塞等待。
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

// XReadStreams 按 streams 顺序从多个流读取消息，参数成对给出流名与最后读取 ID。
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

// XGroupCreate 在流 stream 上创建消费者组 group，起始位置为 start。
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

// XGroupCreateMkStream 同 XGroupCreate，但流不存在时先自动创建空流。
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

// XGroupSetID 将消费者组 group 的最后递送 ID 重置为 start。
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

// XGroupDestroy 删除流 stream 上的消费者组 group，返回被删除的消费者组数量。
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

// XGroupCreateConsumer 在消费者组 group 中创建消费者 consumer，返回受影响消费者数量。
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

// XGroupDelConsumer 从消费者组 group 中删除消费者 consumer，返回其遗留的待处理消息数量。
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

// XReadGroup 按 XReadGroupArgs 以消费者组身份读取消息，读取到的消息进入待处理状态。
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

// XAck 确认消费者组 group 已处理 ids 消息，返回被确认的消息数量。
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

// XPending 返回消费者组 group 的待处理消息汇总信息（总数、ID 区间、消费者分布等）。
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

// XPendingExt 按 XPendingExtArgs 返回待处理消息的明细列表。
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

// XClaim 按 XClaimArgs 将超时未确认的待处理消息认领给指定消费者，返回被认领的消息。
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

// XClaimJustID 同 XClaim，仅返回被认领消息的 ID。
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

// XAutoClaim 按 XAutoClaimArgs 自动批量认领超时消息，返回消息列表与下次认领起始 ID。
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

// XAutoClaimJustID 同 XAutoClaim，仅返回被认领消息 ID 与下次认领起始 ID。
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

// XTrimMaxLen 将流裁剪至最多 maxLen 条消息，返回被裁剪删除的条目数量。
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

// XTrimMaxLenApprox 近似按 maxLen 裁剪流，limit 控制单次修剪步长以降低命令开销。
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

// XTrimMinID 删除流中 ID 小于 minID 的消息，返回被裁剪删除的条目数量。
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

// XTrimMinIDApprox 近似按 minID 裁剪流，limit 控制单次修剪步长以降低命令开销。
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

// XInfoGroups 返回流上全部消费者组的统计信息。
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

// XInfoStream 返回流的基本统计信息（长度、首尾消息 ID、消费者组数量等）。
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

// XInfoStreamFull 返回流的完整信息，count 限制随附返回的消息条目数量。
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

// XInfoConsumers 返回消费者组 group 下各消费者的状态与待处理消息统计。
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
