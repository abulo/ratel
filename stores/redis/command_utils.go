package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// Prefix 返回前缀+键
func (r *Client) Prefix(key string) string {
	return fmt.Sprintf(r.KeyPrefix, key)
}

// k 格式化并返回带前缀的密钥
func (r *Client) k(key any) string {
	return fmt.Sprintf(r.KeyPrefix, cast.ToString(key))
}

// ks 使用前缀格式化并返回一组键
func (r *Client) ks(key ...any) []string {
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = r.k(k)
	}
	return keys
}

func acceptable(err error) bool {
	return err == nil || err == redis.Nil || err == context.Canceled
}

func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}
