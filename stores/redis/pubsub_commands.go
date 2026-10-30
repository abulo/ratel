package redis

import (
	"context"
)

// Publish 将 message 发布到 channel 频道，返回收到该消息的订阅者数量。
func (r *Client) Publish(ctx context.Context, channel, message string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Publish(getCtx(ctx), channel, message).Result()
		return err
	}, acceptable)
	return
}

// SPublish 将 message 发布到分片频道 channel，返回收到该消息的分片订阅者数量。
func (r *Client) SPublish(ctx context.Context, channel, message string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPublish(getCtx(ctx), channel, message).Result()
		return err
	}, acceptable)
	return
}

// PubSubChannels 返回匹配模式 pattern 的当前活跃频道；pattern 为空时返回全部活跃频道。
func (r *Client) PubSubChannels(ctx context.Context, pattern string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubChannels(getCtx(ctx), pattern).Result()
		return err
	}, acceptable)
	return
}

// PubSubNumSub 返回给定频道各自的订阅者数量，返回 map 的键为频道名。
func (r *Client) PubSubNumSub(ctx context.Context, channels ...string) (val map[string]int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubNumSub(getCtx(ctx), channels...).Result()
		return err
	}, acceptable)
	return
}

// PubSubNumPat 返回当前处于订阅状态的模式订阅（PSubscribe）数量。
func (r *Client) PubSubNumPat(ctx context.Context) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubNumPat(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// PubSubShardChannels 返回匹配模式 pattern 的活跃分片频道；pattern 为空时返回全部活跃分片频道。
func (r *Client) PubSubShardChannels(ctx context.Context, pattern string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubShardChannels(getCtx(ctx), pattern).Result()
		return err
	}, acceptable)
	return
}

// PubSubShardNumSub 返回给定分片频道各自的订阅者数量，返回 map 的键为分片频道名。
func (r *Client) PubSubShardNumSub(ctx context.Context, channels ...string) (val map[string]int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubShardNumSub(getCtx(ctx), channels...).Result()
		return err
	}, acceptable)
	return
}
