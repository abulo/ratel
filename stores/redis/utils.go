package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ks 使用前缀格式化并返回一组键 key: 要格式化的键
// func (r *Client) ks(key ...any) []string {
// 	keys := make([]string, len(key))
// 	for i, k := range key {
// 		keys[i] = cast.ToString(k)
// 	}
// 	return keys
// }

// acceptable 判断错误是否可接受 err: 要判断的错误
func acceptable(err error) bool {
	return err == nil || err == redis.Nil || err == context.Canceled
}

// getCtx 获取有效的上下文 ctx: 要检查的上下文
func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}
