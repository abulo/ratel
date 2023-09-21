package redis

import "context"

// Publish 将信息发送到指定的频道。
func (r *Client) Publish(ctx context.Context, channel any, message any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Publish(getCtx(ctx), r.k(channel), message).Result()
		return err
	}, acceptable)
	return
}

// SPublish 将信息发送到指定的频道。
func (r *Client) SPublish(ctx context.Context, channel any, message any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPublish(getCtx(ctx), r.k(channel), message).Result()
		return err
	}, acceptable)
	return
}

// PubSubChannels 订阅一个或多个符合给定模式的频道。
func (r *Client) PubSubChannels(ctx context.Context, pattern any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubChannels(getCtx(ctx), r.k(pattern)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) PubSubNumSub(ctx context.Context, channels ...any) (val map[string]int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubNumSub(getCtx(ctx), r.ks(channels...)...).Result()
		return err
	}, acceptable)
	return
}

// PubSubNumPat 用于获取redis订阅或者发布信息的状态
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

// PubSubShardChannels 订阅一个或多个符合给定模式的频道。
func (r *Client) PubSubShardChannels(ctx context.Context, pattern any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubShardChannels(getCtx(ctx), r.k(pattern)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) PubSubShardNumSub(ctx context.Context, channels ...any) (val map[string]int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PubSubShardNumSub(getCtx(ctx), r.ks(channels...)...).Result()
		return err
	}, acceptable)
	return
}
