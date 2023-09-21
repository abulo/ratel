package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (r *Client) TFunctionLoad(ctx context.Context, lib string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionLoad(getCtx(ctx), lib).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFunctionLoadArgs(ctx context.Context, lib string, options *redis.TFunctionLoadOptions) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionLoadArgs(getCtx(ctx), lib, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFunctionDelete(ctx context.Context, libName string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionDelete(getCtx(ctx), libName).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFunctionList(ctx context.Context) (val []map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionList(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFunctionListArgs(ctx context.Context, options *redis.TFunctionListOptions) (val []map[string]any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFunctionListArgs(getCtx(ctx), options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFCall(ctx context.Context, libName string, funcName string, numKeys int) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCall(getCtx(ctx), libName, funcName, numKeys).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFCallArgs(ctx context.Context, libName string, funcName string, numKeys int, options *redis.TFCallOptions) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallArgs(getCtx(ctx), libName, funcName, numKeys, options).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFCallASYNC(ctx context.Context, libName string, funcName string, numKeys int) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallASYNC(getCtx(ctx), libName, funcName, numKeys).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) TFCallASYNCArgs(ctx context.Context, libName string, funcName string, numKeys int, options *redis.TFCallOptions) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TFCallASYNCArgs(getCtx(ctx), libName, funcName, numKeys, options).Result()
		return err
	}, acceptable)
	return
}
