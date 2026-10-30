package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// func (r *Client) ks(key ...any) []string {
// 	keys := make([]string, len(key))
// 	for i, k := range key {
// 		keys[i] = cast.ToString(k)
// 	}
// 	return keys
// }

// acceptable 判断错误是否属于可接受的失败：命令成功、键不存在（redis.Nil）或上下文已取消均视为可接受，
// 不计入熔断器的失败统计。
func acceptable(err error) bool {
	return err == nil || err == redis.Nil || err == context.Canceled
}

// getCtx 返回可用于执行命令的上下文；入参为 nil 或已取消时替换为 context.TODO()，
// 避免空上下文或已取消上下文导致命令直接失败。
func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}
