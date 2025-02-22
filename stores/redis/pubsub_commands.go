package redis

import (
	"context"
)

// Publish 将信息发送到指定的频道 ctx: 上下文 channel: 频道名称 message: 要发送的消息
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

// SPublish 将信息发送到指定的频道 ctx: 上下文 channel: 频道名称 message: 要发送的消息
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

// PubSubChannels 获取匹配模式的活跃频道列表 ctx: 上下文 pattern: 频道匹配模式
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

// PubSubNumSub 获取指定频道的订阅数量 ctx: 上下文 channels: 要查询的频道列表
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

// PubSubNumPat 获取模式订阅的数量 ctx: 上下文
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

// PubSubShardChannels 获取匹配模式的活跃分片频道列表 ctx: 上下文 pattern: 频道匹配模式
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

// PubSubShardNumSub 获取指定分片频道的订阅数量 ctx: 上下文 channels: 要查询的分片频道列表
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
