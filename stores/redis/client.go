package redis

// OriginalClient 获取单节点模式的底层 Redis 实例，内部复用 getClient 的连接复用逻辑；
// 调用方应保证传入的 Client 为单节点模式配置。
func OriginalClient(r *Client) (RedisNode, error) {
	return getClient(r)
}
