package redis

func OriginalClient(r *Client) (RedisNode, error) {
	return getClient(r)
}
