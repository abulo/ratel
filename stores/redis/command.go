package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Prefix 返回前缀+键
func (r *Client) Prefix(key string) string {
	return fmt.Sprintf(r.KeyPrefix, key)
}

// k 格式化并返回带前缀的密钥
func (r *Client) k(key string) string {
	return fmt.Sprintf(r.KeyPrefix, key)
}

// ks 使用前缀格式化并返回一组键
func (r *Client) ks(key ...string) []string {
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = r.k(k)
	}
	return keys
}

func acceptable(err error) bool {
	return err == nil || err == redis.Nil || err == context.Canceled
}

// Pipeline 获取管道
func (r *Client) Pipeline() (val redis.Pipeliner, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.Pipeline()
		return nil
	}, acceptable)
	return
}
func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}

// Pipelined 管道
func (r *Client) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) (val []redis.Cmder, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Pipelined(getCtx(ctx), fn)
		return err
	}, acceptable)
	return
}

// TxPipelined 管道
func (r *Client) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) (val []redis.Cmder, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TxPipelined(getCtx(ctx), fn)
		return err
	}, acceptable)
	return
}

// TxPipeline 获取管道
func (r *Client) TxPipeline() (val redis.Pipeliner, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.TxPipeline()
		return nil
	}, acceptable)
	return
}

// Command 返回有关所有Redis命令的详细信息的Array回复
func (r *Client) Command(ctx context.Context) (val map[string]*redis.CommandInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Command(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClientGetName returns the name of the connection.
func (r *Client) ClientGetName(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientGetName(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Echo  批量字符串回复
func (r *Client) Echo(ctx context.Context, message any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Echo(getCtx(ctx), message).Result()
		return err
	}, acceptable)
	return
}

// Ping 使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
// 通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值。
// 如果连接正常就返回一个 PONG ，否则返回一个连接错误。
func (r *Client) Ping(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Ping(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Quit 关闭连接
func (r *Client) Quit(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Quit(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Del 删除给定的一个或多个 key 。
// 不存在的 key 会被忽略。
func (r *Client) Del(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Del(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// Unlink 这个命令非常类似于DEL：它删除指定的键。就像DEL键一样，如果它不存在，它将被忽略。但是，该命令在不同的线程中执行实际的内存回收，所以它不会阻塞，而DEL是。这是命令名称的来源：命令只是将键与键空间断开连接。实际删除将在以后异步发生。
func (r *Client) Unlink(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Unlink(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// Dump 序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (r *Client) Dump(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Dump(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Exists 检查给定 key 是否存在。
// 若 key 存在，返回 1 ，否则返回 0 。
func (r *Client) Exists(ctx context.Context, key ...string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		valTmp, err := conn.Exists(getCtx(ctx), r.ks(key...)...).Result()
		val = valTmp > 0
		return err
	}, acceptable)
	return
}

// Expire 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
// 设置成功返回 1 。
// 当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Expire(getCtx(ctx), r.k(key), expiration).Result()
		return err
	}, acceptable)
	return
}

// ExpireAt  EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireAt(ctx context.Context, key string, tm time.Time) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireAt(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireNX  ExpireNX 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireNX(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireNX(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireXX  ExpireXX 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireXX(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireXX(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireGT  ExpireGT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireGT(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireGT(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireLT  ExpireLT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireLT(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireLT(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// Keys 查找所有符合给定模式 pattern 的 key 。
func (r *Client) Keys(ctx context.Context, pattern string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Keys(getCtx(ctx), r.k(pattern)).Result()
		return err
	}, acceptable)
	return
}

// Migrate 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
func (r *Client) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Migrate(getCtx(ctx), host, port, r.k(key), db, timeout).Result()
		return err
	}, acceptable)
	return
}

// Move 将当前数据库的 key 移动到给定的数据库 db 当中。
// 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。
func (r *Client) Move(ctx context.Context, key string, db int) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Move(getCtx(ctx), r.k(key), db).Result()
		return err
	}, acceptable)
	return
}

// ObjectRefCount 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func (r *Client) ObjectRefCount(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectRefCount(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ObjectEncoding 返回给定 key 锁储存的值所使用的内部表示(representation)。
func (r *Client) ObjectEncoding(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectEncoding(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ObjectIdleTime 返回给定 key 自储存以来的空转时间(idle， 没有被读取也没有被写入)，以秒为单位。
func (r *Client) ObjectIdleTime(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectIdleTime(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Persist 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。
// 当生存时间移除成功时，返回 1 .
// 如果 key 不存在或 key 没有设置生存时间，返回 0 。
func (r *Client) Persist(ctx context.Context, key string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Persist(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// PExpire 毫秒为单位设置 key 的生存时间
func (r *Client) PExpire(ctx context.Context, key string, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpire(getCtx(ctx), r.k(key), expiration).Result()
		return err
	}, acceptable)
	return
}

// PExpireAt 这个命令和 expireat 命令类似，但它以毫秒为单位设置 key 的过期 unix 时间戳，而不是像 expireat 那样，以秒为单位。
// 如果生存时间设置成功，返回 1 。 当 key 不存在或没办法设置生存时间时，返回 0
func (r *Client) PExpireAt(ctx context.Context, key string, tm time.Time) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpireAt(getCtx(ctx), r.k(key), tm).Result()
		return err
	}, acceptable)
	return
}

// PTTL 这个命令类似于 TTL 命令，但它以毫秒为单位返回 key 的剩余生存时间，而不是像 TTL 命令那样，以秒为单位。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以毫秒为单位，返回 key 的剩余生存时间。
func (r *Client) PTTL(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PTTL(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// RandomKey 从当前数据库中随机返回(不删除)一个 key 。
func (r *Client) RandomKey(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RandomKey(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Rename 将 key 改名为 newkey 。
// 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。
// 当 newkey 已经存在时， RENAME 命令将覆盖旧值。
func (r *Client) Rename(ctx context.Context, key, newkey string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Rename(getCtx(ctx), r.k(key), r.k(newkey)).Result()
		return err
	}, acceptable)
	return
}

// RenameNX 当且仅当 newkey 不存在时，将 key 改名为 newkey 。
// 当 key 不存在时，返回一个错误。
func (r *Client) RenameNX(ctx context.Context, key, newkey string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RenameNX(getCtx(ctx), r.k(key), r.k(newkey)).Result()
		return err
	}, acceptable)
	return
}

// Restore 反序列化给定的序列化值，并将它和给定的 key 关联。
// 参数 ttl 以毫秒为单位为 key 设置生存时间；如果 ttl 为 0 ，那么不设置生存时间。
// RESTORE 在执行反序列化之前会先对序列化值的 RDB 版本和数据校验和进行检查，如果 RDB 版本不相同或者数据不完整的话，那么 RESTORE 会拒绝进行反序列化，并返回一个错误。
func (r *Client) Restore(ctx context.Context, key string, ttl time.Duration, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Restore(getCtx(ctx), r.k(key), ttl, value).Result()
		return err
	}, acceptable)
	return
}

// RestoreReplace -> Restore
func (r *Client) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RestoreReplace(getCtx(ctx), r.k(key), ttl, value).Result()
		return err
	}, acceptable)
	return
}

// Sort 返回或保存给定列表、集合、有序集合 key 中经过排序的元素。
// 排序默认以数字作为对象，值被解释为双精度浮点数，然后进行比较。
func (r *Client) Sort(ctx context.Context, key string, sort *redis.Sort) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Sort(getCtx(ctx), r.k(key), sort).Result()
		return err
	}, acceptable)
	return
}

// SortRO -> Sort
func (r *Client) SortRO(ctx context.Context, key string, sort *redis.Sort) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortRO(getCtx(ctx), r.k(key), sort).Result()
		return err
	}, acceptable)
	return
}

// SortStore -> Sort
func (r *Client) SortStore(ctx context.Context, key, store string, sort *redis.Sort) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortStore(getCtx(ctx), r.k(key), r.k(store), sort).Result()
		return err
	}, acceptable)
	return
}

// SortInterfaces -> Sort
func (r *Client) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortInterfaces(getCtx(ctx), r.k(key), sort).Result()
		return err
	}, acceptable)
	return
}

// Touch 更改键的上次访问时间。返回指定的现有键的数量。
func (r *Client) Touch(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Touch(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以秒为单位，返回 key 的剩余生存时间。
func (r *Client) TTL(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TTL(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Type 返回 key 所储存的值的类型。
func (r *Client) Type(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Type(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Append 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (r *Client) Append(ctx context.Context, key, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Append(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// Decr 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 执行 DECR 命令之后 key 的值。
func (r *Client) Decr(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Decr(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// DecrBy 将 key 所储存的值减去减量 decrement 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于更多递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 减去 decrement 之后， key 的值。
func (r *Client) DecrBy(ctx context.Context, key string, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DecrBy(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// Get 返回 key 所关联的字符串值。
// 如果 key 不存在那么返回特殊值 nil 。
// 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。
// 当 key 不存在时，返回 nil ，否则，返回 key 的值。
// 如果 key 不是字符串类型，那么返回一个错误。
func (r *Client) Get(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Get(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// GetRange 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
// 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
// GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
// 返回截取得出的子字符串。
func (r *Client) GetRange(ctx context.Context, key string, start, end int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetRange(getCtx(ctx), r.k(key), start, end).Result()
		return err
	}, acceptable)
	return
}

// GetSet 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 返回给定 key 的旧值。
// 当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
func (r *Client) GetSet(ctx context.Context, key string, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetSet(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// GetEx
func (r *Client) GetEx(ctx context.Context, key string, ts time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetEx(getCtx(ctx), r.k(key), ts).Result()
		return err
	}, acceptable)
	return
}

// GetEx
func (r *Client) GetDel(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetDel(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Incr 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 执行 INCR 命令之后 key 的值。
func (r *Client) Incr(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Incr(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// IncrBy 将 key 所储存的值加上增量 increment 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，参见 INCR 命令。
// 加上 increment 之后， key 的值。
func (r *Client) IncrBy(ctx context.Context, key string, value int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrBy(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// IncrByFloat 为 key 中所储存的值加上浮点数增量 increment 。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
// 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
func (r *Client) IncrByFloat(ctx context.Context, key string, value float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.IncrByFloat(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// MGet 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
// 一个包含所有给定 key 的值的列表。
func (r *Client) MGet(ctx context.Context, keys ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MGet(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// MSet 同时设置一个或多个 key-value 对。
func (r *Client) MSet(ctx context.Context, values ...any) (val string, err error) {
	// return getRedis(r).MSet(getCtx(ctx), values...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MSet(getCtx(ctx), values...).Result()
		return err
	}, acceptable)
	return

}

// MSetNX 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
func (r *Client) MSetNX(ctx context.Context, values ...any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MSetNX(getCtx(ctx), values...).Result()
		return err
	}, acceptable)
	return
}

// Set 将字符串值 value 关联到 key 。
// 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
// 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
func (r *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Set(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) SetArgs(ctx context.Context, key string, value any, a redis.SetArgs) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetArgs(getCtx(ctx), r.k(key), value, a).Result()
		return err
	}, acceptable)
	return
}

// SetEX ...
func (r *Client) SetEx(ctx context.Context, key string, value any, expiration time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetEx(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetNX 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
// SETNX 是『SET if Not eXists』(如果不存在，则 SET)的简写。
func (r *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetNX(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetXX -> SetNX
func (r *Client) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetXX(getCtx(ctx), r.k(key), value, expiration).Result()
		return err
	}, acceptable)
	return
}

// SetRange 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (r *Client) SetRange(ctx context.Context, key string, offset int64, value string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetRange(getCtx(ctx), r.k(key), offset, value).Result()
		return err
	}, acceptable)
	return
}

// StrLen 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。
func (r *Client) StrLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.StrLen(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// Copy
func (r *Client) Copy(ctx context.Context, sourceKey string, destKey string, db int, replace bool) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Copy(getCtx(ctx), r.k(sourceKey), r.k(destKey), db, replace).Result()
		return err
	}, acceptable)
	return
}

// GetBit 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
// 字符串值指定偏移量上的位(bit)。
func (r *Client) GetBit(ctx context.Context, key string, offset int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GetBit(getCtx(ctx), r.k(key), offset).Result()
		return err
	}, acceptable)
	return
}

// SetBit 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 当 key 不存在时，自动生成一个新的字符串值。
// 字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
// offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
func (r *Client) SetBit(ctx context.Context, key string, offset int64, value int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SetBit(getCtx(ctx), r.k(key), offset, value).Result()
		return err
	}, acceptable)
	return
}

// BitCount 计算给定字符串中，被设置为 1 的比特位的数量。
// 一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
// start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。
// 不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
func (r *Client) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitCount(getCtx(ctx), r.k(key), bitCount).Result()
		return err
	}, acceptable)
	return
}

// BitOpAnd -> BitCount
func (r *Client) BitOpAnd(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpAnd(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpOr -> BitCount
func (r *Client) BitOpOr(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpOr(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpXor -> BitCount
func (r *Client) BitOpXor(ctx context.Context, destKey string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpXor(getCtx(ctx), r.k(destKey), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BitOpNot -> BitCount
func (r *Client) BitOpNot(ctx context.Context, destKey string, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitOpNot(getCtx(ctx), r.k(destKey), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// BitPos -> BitCount
func (r *Client) BitPos(ctx context.Context, key string, bit int64, pos ...int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitPos(getCtx(ctx), r.k(key), bit, pos...).Result()
		return err
	}, acceptable)
	return
}

// BitField -> BitCount
func (r *Client) BitField(ctx context.Context, key string, args ...any) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BitField(getCtx(ctx), r.k(key), args...).Result()
		return err
	}, acceptable)
	return
}

// Scan 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素
func (r *Client) Scan(ctx context.Context, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	// return getRedis(r).Scan(getCtx(ctx), cursor, r.k(match), count)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.Scan(getCtx(ctx), cursorIn, r.k(match), count).Result()
		return err
	}, acceptable)
	return
}

// Scan 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素
func (r *Client) ScanIterator(ctx context.Context, cursorIn uint64, match string, count int64) (val *redis.ScanIterator, err error) {
	// return getRedis(r).Scan(getCtx(ctx), cursor, r.k(match), count)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.Scan(getCtx(ctx), cursorIn, r.k(match), count).Iterator()
		return err
	}, acceptable)
	return
}

func (r *Client) ScanType(ctx context.Context, cursorIn uint64, match string, count int64, keyType string) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.ScanType(getCtx(ctx), cursorIn, r.k(match), count, keyType).Result()
		return err
	}, acceptable)
	return
}

// SScan 详细信息请参考 SCAN 命令。
func (r *Client) SScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.SScan(getCtx(ctx), r.k(key), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// HScan 详细信息请参考 SCAN 命令。
func (r *Client) HScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.HScan(getCtx(ctx), r.k(key), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// ZScan 详细信息请参考 SCAN 命令。
func (r *Client) ZScan(ctx context.Context, key string, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.ZScan(getCtx(ctx), r.k(key), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func (r *Client) HDel(ctx context.Context, key string, fields ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HDel(getCtx(ctx), r.k(key), fields...).Result()
		return err
	}, acceptable)
	return
}

// HExists 查看哈希表 key 中，给定域 field 是否存在。
func (r *Client) HExists(ctx context.Context, key, field string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HExists(getCtx(ctx), r.k(key), field).Result()
		return err
	}, acceptable)
	return
}

// HGet 返回哈希表 key 中给定域 field 的值。
func (r *Client) HGet(ctx context.Context, key, field string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGet(getCtx(ctx), r.k(key), field).Result()
		return err
	}, acceptable)
	return
}

// HGetAll 返回哈希表 key 中，所有的域和值。
// 在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
func (r *Client) HGetAll(ctx context.Context, key string) (val map[string]string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HGetAll(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HIncrBy 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (r *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrBy(getCtx(ctx), r.k(key), field, incr).Result()
		return err
	}, acceptable)
	return
}

// HIncrByFloat 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HIncrByFloat(getCtx(ctx), r.k(key), field, incr).Result()
		return err
	}, acceptable)
	return
}

// HKeys 返回哈希表 key 中的所有域。
func (r *Client) HKeys(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HKeys(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HLen 返回哈希表 key 中域的数量。
func (r *Client) HLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HLen(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// HMGet 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMGet(getCtx(ctx), r.k(key), fields...).Result()
		return err
	}, acceptable)
	return
}

// HSet 将哈希表 key 中的域 field 的值设为 value 。
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Client) HSet(ctx context.Context, key string, value ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSet(getCtx(ctx), r.k(key), value...).Result()
		return err
	}, acceptable)
	return
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (r *Client) HMSet(ctx context.Context, key string, value ...any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HMSet(getCtx(ctx), r.k(key), value...).Result()
		return err
	}, acceptable)
	return
}

// HSetNX 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
// 若域 field 已经存在，该操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func (r *Client) HSetNX(ctx context.Context, key, field string, value any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HSetNX(getCtx(ctx), r.k(key), field, value).Result()
		return err
	}, acceptable)
	return
}

// HVals 返回哈希表 key 中所有域的值。
func (r *Client) HVals(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HVals(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) HRandField(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandField(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) HRandFieldWithValues(ctx context.Context, key string, count int) (val []redis.KeyValue, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.HRandFieldWithValues(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// BLPop 是列表的阻塞式(blocking)弹出原语。
// 它是 LPop 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPop 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的头元素。
func (r *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BLPop(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BRPop 是列表的阻塞式(blocking)弹出原语。
// 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的尾部元素。
// 关于阻塞操作的更多信息，请查看 BLPOP 命令， BRPOP 除了弹出元素的位置和 BLPOP 不同之外，其他表现一致。
func (r *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BRPop(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BRPopLPush 是 RPOPLPUSH 的阻塞版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH 一样。
// 当列表 source 为空时， BRPOPLPUSH 命令将阻塞连接，直到等待超时，或有另一个客户端对 source 执行 LPUSH 或 RPUSH 命令为止。
func (r *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BRPopLPush(getCtx(ctx), r.k(source), r.k(destination), timeout).Result()
		return err
	}, acceptable)
	return
}

// LIndex 返回列表 key 中，下标为 index 的元素。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LIndex(ctx context.Context, key string, index int64) (val string, err error) {
	// return getRedis(r).LIndex(getCtx(ctx), r.k(key), index)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LIndex(getCtx(ctx), r.k(key), index).Result()
		return err
	}, acceptable)
	return
}

// LInsert 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
// 当 pivot 不存在于列表 key 时，不执行任何操作。
// 当 key 不存在时， key 被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LInsert(ctx context.Context, key, op string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsert(getCtx(ctx), r.k(key), op, pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LInsertBefore 同 LInsert
func (r *Client) LInsertBefore(ctx context.Context, key string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsertBefore(getCtx(ctx), r.k(key), pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LInsertAfter 同 LInsert
func (r *Client) LInsertAfter(ctx context.Context, key string, pivot, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LInsertAfter(getCtx(ctx), r.k(key), pivot, value).Result()
		return err
	}, acceptable)
	return
}

// LLen 返回列表 key 的长度。
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LLen(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LLen(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// LPop 移除并返回列表 key 的头元素。
func (r *Client) LPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPop(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// LPopCount 移除并返回列表 key 的头元素。
func (r *Client) LPopCount(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPopCount(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) LPos(ctx context.Context, key string, value string, args redis.LPosArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPos(getCtx(ctx), r.k(key), value, args).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) LPosCount(ctx context.Context, key string, value string, count int64, args redis.LPosArgs) (val []int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPosCount(getCtx(ctx), r.k(key), value, count, args).Result()
		return err
	}, acceptable)
	return
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头
// 如果 key 不存在，一个空列表会被创建并执行 LPush 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) LPush(ctx context.Context, key string, values ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPush(getCtx(ctx), r.k(key), values...).Result()
		return err
	}, acceptable)
	return
}

// LPushX 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
// 和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
func (r *Client) LPushX(ctx context.Context, key string, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LPushX(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (r *Client) LRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LRange(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// LRem 根据参数 count 的值，移除列表中与参数 value 相等的元素。
func (r *Client) LRem(ctx context.Context, key string, count int64, value any) (val int64, err error) {

	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LRem(getCtx(ctx), r.k(key), count, value).Result()
		return err
	}, acceptable)
	return
}

// LSet 将列表 key 下标为 index 的元素的值设置为 value 。
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
// 关于列表下标的更多信息，请参考 LINDEX 命令。
func (r *Client) LSet(ctx context.Context, key string, index int64, value any) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LSet(getCtx(ctx), r.k(key), index, value).Result()
		return err
	}, acceptable)
	return
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 当 key 不是列表类型时，返回一个错误。
func (r *Client) LTrim(ctx context.Context, key string, start, stop int64) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LTrim(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// RPop 移除并返回列表 key 的头元素。
func (r *Client) RPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPop(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// RPopCount 移除并返回列表 key 的头元素。
func (r *Client) RPopCount(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPopCount(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// RPopLPush 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
// 将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
// 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
// 举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作。
// 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
func (r *Client) RPopLPush(ctx context.Context, source, destination string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPopLPush(getCtx(ctx), r.k(source), r.k(destination)).Result()
		return err
	}, acceptable)
	return
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
// 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) RPush(ctx context.Context, key string, values ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPush(getCtx(ctx), r.k(key), values...).Result()
		return err
	}, acceptable)
	return
}

// RPushX 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
// 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。
func (r *Client) RPushX(ctx context.Context, key string, value any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RPushX(getCtx(ctx), r.k(key), value).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) LMove(ctx context.Context, source, destination, srcpos, destpos string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LMove(getCtx(ctx), r.k(source), r.k(destination), r.k(srcpos), r.k(destpos)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) BLMove(ctx context.Context, source, destination, srcpos, destpos string, ts time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BLMove(getCtx(ctx), r.k(source), r.k(destination), r.k(srcpos), r.k(destpos), ts).Result()
		return err
	}, acceptable)
	return
}

// SAdd 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
// 当 key 不是集合类型时，返回一个错误。
func (r *Client) SAdd(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SAdd(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// SCard 返回集合 key 的基数(集合中元素的数量)。
func (r *Client) SCard(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SCard(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SDiff 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (r *Client) SDiff(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiff(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SDiffStore 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SDiffStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SDiffStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInter 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
// 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func (r *Client) SInter(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInter(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterCard(ctx context.Context, limit int64, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterCard(getCtx(ctx), limit, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SInterStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SIsMember 判断 member 元素是否集合 key 的成员。
func (r *Client) SIsMember(ctx context.Context, key string, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SIsMember(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// SMembers 返回集合 key 中的所有成员。
// 不存在的 key 被视为空集合。
func (r *Client) SMembers(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembers(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SMembersMap -> SMembers
func (r *Client) SMembersMap(ctx context.Context, key string) (val map[string]struct{}, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMembersMap(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SMove 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
func (r *Client) SMove(ctx context.Context, source, destination string, member any) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SMove(getCtx(ctx), r.k(source), r.k(destination), member).Result()
		return err
	}, acceptable)
	return
}

// SPop 移除并返回集合中的一个随机元素。
// 如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
func (r *Client) SPop(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPop(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SPopN -> SPop
func (r *Client) SPopN(ctx context.Context, key string, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SPopN(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// SRandMember 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， SRANDMEMBER 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 SRANDMEMBER 则仅仅返回随机元素，而不对集合进行任何改动。
func (r *Client) SRandMember(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMember(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// SRandMemberN -> SRandMember
func (r *Client) SRandMemberN(ctx context.Context, key string, count int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRandMemberN(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

// SRem 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
func (r *Client) SRem(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SRem(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// SUnion 返回一个集合的全部成员，该集合是所有给定集合的并集。
// 不存在的 key 被视为空集。
func (r *Client) SUnion(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnion(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// SUnionStore 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SUnionStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// XAdd 将指定的流条目追加到指定key的流中。 如果key不存在，作为运行这个命令的副作用，将使用流的条目自动创建key。
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

// XDel 从指定流中移除指定的条目，并返回成功删除的条目的数量，在传递的ID不存在的情况下， 返回的数量可能与传递的ID数量不同。
func (r *Client) XDel(ctx context.Context, stream string, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XDel(getCtx(ctx), r.k(stream), ids...).Result()
		return err
	}, acceptable)
	return
}

// XLen 返回流中的条目数。如果指定的key不存在，则此命令返回0，就好像该流为空。 但是请注意，与其他的Redis类型不同，零长度流是可能的，所以你应该调用TYPE 或者 EXISTS 来检查一个key是否存在。
// 一旦内部没有任何的条目（例如调用XDEL后），流不会被自动删除，因为可能还存在与其相关联的消费者组。
func (r *Client) XLen(ctx context.Context, stream string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XLen(getCtx(ctx), r.k(stream)).Result()
		return err
	}, acceptable)
	return
}

// XRange 此命令返回流中满足给定ID范围的条目。范围由最小和最大ID指定。所有ID在指定的两个ID之间或与其中一个ID相等（闭合区间）的条目将会被返回。
func (r *Client) XRange(ctx context.Context, stream, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRange(getCtx(ctx), r.k(stream), start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRangeN -> XRange
func (r *Client) XRangeN(ctx context.Context, stream, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRangeN(getCtx(ctx), r.k(stream), start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRevRange 此命令与XRANGE完全相同，但显著的区别是以相反的顺序返回条目，并以相反的顺序获取开始-结束参数：在XREVRANGE中，你需要先指定结束ID，再指定开始ID，该命令就会从结束ID侧开始生成两个ID之间（或完全相同）的所有元素。
func (r *Client) XRevRange(ctx context.Context, stream string, start, stop string) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRange(getCtx(ctx), r.k(stream), start, stop).Result()
		return err
	}, acceptable)
	return
}

// XRevRangeN -> XRevRange
func (r *Client) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) (val []redis.XMessage, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XRevRangeN(getCtx(ctx), r.k(stream), start, stop, count).Result()
		return err
	}, acceptable)
	return
}

// XRead 从一个或者多个流中读取数据，仅返回ID大于调用者报告的最后接收ID的条目。此命令有一个阻塞选项，用于等待可用的项目，类似于BRPOP或者BZPOPMIN等等。
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

// XReadStreams -> XRead
func (r *Client) XReadStreams(ctx context.Context, streams ...string) (val []redis.XStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XReadStreams(getCtx(ctx), r.ks(streams...)...).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreate command
func (r *Client) XGroupCreate(ctx context.Context, stream, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreate(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupCreateMkStream command
func (r *Client) XGroupCreateMkStream(ctx context.Context, stream, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupCreateMkStream(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupSetID command
func (r *Client) XGroupSetID(ctx context.Context, stream, group, start string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupSetID(getCtx(ctx), r.k(stream), group, start).Result()
		return err
	}, acceptable)
	return
}

// XGroupDestroy command
func (r *Client) XGroupDestroy(ctx context.Context, stream, group string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDestroy(getCtx(ctx), r.k(stream), group).Result()
		return err
	}, acceptable)
	return
}

// XGroupDelConsumer command
func (r *Client) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XGroupDelConsumer(getCtx(ctx), r.k(stream), group, consumer).Result()
		return err
	}, acceptable)
	return
}

// XReadGroup command
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

// XAck command
func (r *Client) XAck(ctx context.Context, stream, group string, ids ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XAck(getCtx(ctx), r.k(stream), group, ids...).Result()
		return err
	}, acceptable)
	return
}

// XPending command
func (r *Client) XPending(ctx context.Context, stream, group string) (val *redis.XPending, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XPending(getCtx(ctx), r.k(stream), group).Result()
		return err
	}, acceptable)
	return
}

// XPendingExt command
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

// XClaim command
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

// XClaimJustID command
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

func (r *Client) XTrimMaxLen(ctx context.Context, key string, maxLen int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLen(getCtx(ctx), r.k(key), maxLen).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMaxLenApprox(ctx context.Context, key string, maxLen, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMaxLenApprox(getCtx(ctx), r.k(key), maxLen, limit).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMinID(ctx context.Context, key string, minID string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinID(getCtx(ctx), r.k(key), minID).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XTrimMinIDApprox(getCtx(ctx), r.k(key), minID, limit).Result()
		return err
	}, acceptable)
	return
}

// XInfoGroups command
func (r *Client) XInfoGroups(ctx context.Context, key string) (val []redis.XInfoGroup, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoGroups(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoStream(ctx context.Context, key string) (val *redis.XInfoStream, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStream(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoStreamFull(ctx context.Context, key string, count int) (val *redis.XInfoStreamFull, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoStreamFull(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) XInfoConsumers(ctx context.Context, key string, group string) (val []redis.XInfoConsumer, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.XInfoConsumers(getCtx(ctx), r.k(key), group).Result()
		return err
	}, acceptable)
	return
}

// BZPopMax 是有序集合命令 ZPOPMAX带有阻塞功能的版本。
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMax(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// BZPopMin 是有序集合命令 ZPOPMIN带有阻塞功能的版本。
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (val *redis.ZWithKey, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BZPopMin(getCtx(ctx), timeout, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
// 如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。
// score 值可以是整数值或双精度浮点数。
// 如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	// return getRedis(r).ZAdd(getCtx(ctx), r.k(key), members...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAdd(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddNX -> ZAdd
func (r *Client) ZAddNX(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddNX(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZAddXX -> ZAdd
func (r *Client) ZAddXX(ctx context.Context, key string, members ...redis.Z) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddXX(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgs(getCtx(ctx), r.k(key), args).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZAddArgsIncr(ctx context.Context, key string, args redis.ZAddArgs) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZAddArgsIncr(getCtx(ctx), r.k(key), args).Result()
		return err
	}, acceptable)
	return
}

// ZCard 返回有序集 key 的基数。
func (r *Client) ZCard(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCard(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 关于参数 min 和 max 的详细使用方法，请参考 ZRANGEBYSCORE 命令。
func (r *Client) ZCount(ctx context.Context, key, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZCount(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZLexCount -> ZCount
func (r *Client) ZLexCount(ctx context.Context, key, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZLexCount(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
// 当 key 不是有序集类型时，返回一个错误。
// score 值可以是整数值或双精度浮点数。
func (r *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZIncrBy(getCtx(ctx), r.k(key), increment, member).Result()
		return err
	}, acceptable)
	return
}

// ZInter 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
func (r *Client) ZInter(ctx context.Context, store *redis.ZStore) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInter(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZInterWithScores(ctx context.Context, store *redis.ZStore) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterWithScores(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZInterCard(ctx context.Context, limit int64, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterCard(getCtx(ctx), limit, r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ZInterStore 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
func (r *Client) ZInterStore(ctx context.Context, key string, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZInterStore(getCtx(ctx), r.k(key), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZMScore(ctx context.Context, key string, members ...string) (val []float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZMScore(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMax 删除并返回有序集合key中的最多count个具有最高得分的成员。
func (r *Client) ZPopMax(ctx context.Context, key string, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMax(getCtx(ctx), r.k(key), count...).Result()
		return err
	}, acceptable)
	return
}

// ZPopMin 删除并返回有序集合key中的最多count个具有最低得分的成员。
func (r *Client) ZPopMin(ctx context.Context, key string, count ...int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZPopMin(getCtx(ctx), r.k(key), count...).Result()
		return err
	}, acceptable)
	return
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (r *Client) ZRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRange(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeWithScores -> ZRange
func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeWithScores(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
func (r *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScore(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByLex -> ZRangeByScore
func (r *Client) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByLex(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRangeByScoreWithScores -> ZRangeByScore
func (r *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeByScoreWithScores(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRangeArgs(ctx context.Context, z redis.ZRangeArgs) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeArgs(getCtx(ctx), z).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRangeArgsWithScores(ctx context.Context, z redis.ZRangeArgs) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeArgsWithScores(getCtx(ctx), z).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRangeStore(ctx context.Context, dst string, z redis.ZRangeArgs) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRangeStore(getCtx(ctx), r.k(dst), z).Result()
		return err
	}, acceptable)
	return

}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
// 使用 ZREVRANK 命令可以获得成员按 score 值递减(从大到小)排列的排名。
func (r *Client) ZRank(ctx context.Context, key, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRank(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZRem(ctx context.Context, key string, members ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRem(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
func (r *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByRank(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByScore 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 自版本2.1.6开始， score 值等于 min 或 max 的成员也可以不包括在内，详情请参见 ZRANGEBYSCORE 命令。
func (r *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByScore(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRemRangeByLex -> ZRemRangeByScore
func (r *Client) ZRemRangeByLex(ctx context.Context, key, min, max string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRemRangeByLex(getCtx(ctx), r.k(key), min, max).Result()
		return err
	}, acceptable)
	return
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func (r *Client) ZRevRange(ctx context.Context, key string, start, stop int64) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRange(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeWithScores -> ZRevRange
func (r *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeWithScores(getCtx(ctx), r.k(key), start, stop).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScore 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func (r *Client) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScore(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByLex -> ZRevRangeByScore
func (r *Client) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByLex(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRangeByScoreWithScores -> ZRevRangeByScore
func (r *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRangeByScoreWithScores(getCtx(ctx), r.k(key), opt).Result()
		return err
	}, acceptable)
	return
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (r *Client) ZRevRank(ctx context.Context, key, member string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRevRank(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *Client) ZScore(ctx context.Context, key, member string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZScore(getCtx(ctx), r.k(key), member).Result()
		return err
	}, acceptable)
	return
}

// ZUnionStore 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之 和 。
func (r *Client) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnionStore(getCtx(ctx), r.k(dest), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRandMember(ctx context.Context, key string, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMember(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZRandMemberWithScores(ctx context.Context, key string, count int) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZRandMemberWithScores(getCtx(ctx), r.k(key), count).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZUnion(ctx context.Context, store redis.ZStore) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnion(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZUnionWithScores(ctx context.Context, store redis.ZStore) (val []redis.Z, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZUnionWithScores(getCtx(ctx), store).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZDiff(ctx context.Context, keys ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiff(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZDiffWithScores(ctx context.Context, keys ...string) (val []redis.Z, err error) {
	// return getRedis(r).ZDiffWithScores(getCtx(ctx), r.ks(keys...)...)
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffWithScores(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ZDiffStore(ctx context.Context, destination string, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ZDiffStore(getCtx(ctx), r.k(destination), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// PFAdd 将指定元素添加到HyperLogLog
func (r *Client) PFAdd(ctx context.Context, key string, els ...any) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFAdd(getCtx(ctx), r.k(key), els...).Result()
		return err
	}, acceptable)
	return
}

// PFCount 返回HyperlogLog观察到的集合的近似基数。
func (r *Client) PFCount(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFCount(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// PFMerge N个不同的HyperLogLog合并为一个。
func (r *Client) PFMerge(ctx context.Context, dest string, keys ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PFMerge(getCtx(ctx), r.k(dest), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ConfigGet(ctx context.Context, parameter string) (val map[string]string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ConfigGet(getCtx(ctx), parameter).Result()
		return err
	}, acceptable)
	return
}

// ConfigResetStat 重置 INFO 命令中的某些统计数据
func (r *Client) ConfigResetStat(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ConfigResetStat(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ConfigSet 修改 redis 配置参数，无需重启
func (r *Client) ConfigSet(ctx context.Context, parameter, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ConfigSet(getCtx(ctx), parameter, value).Result()
		return err
	}, acceptable)
	return
}

// ConfigRewrite 对启动 Redis 服务器时所指定的 redis.conf 配置文件进行改写
func (r *Client) ConfigRewrite(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ConfigRewrite(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// BgRewriteAOF 异步重写附加文件
func (r *Client) BgRewriteAOF(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BgRewriteAOF(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// BgSave 将数据集异步保存到磁盘
func (r *Client) BgSave(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.BgSave(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClientKill 杀掉客户端的连接
func (r *Client) ClientKill(ctx context.Context, ipPort string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientKill(getCtx(ctx), ipPort).Result()
		return err
	}, acceptable)
	return
}

// ClientKillByFilter is new style synx, while the ClientKill is old
// CLIENT KILL <option> [value] ... <option> [value]
func (r *Client) ClientKillByFilter(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientKillByFilter(getCtx(ctx), r.ks(keys...)...).Result()
		return err
	}, acceptable)
	return
}

// ClientList 获取客户端连接列表
func (r *Client) ClientList(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientList(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClientPause 停止处理来自客户端的命令一段时间
func (r *Client) ClientPause(ctx context.Context, dur time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientPause(getCtx(ctx), dur).Result()
		return err
	}, acceptable)
	return
}

// ClientPause 停止处理来自客户端的命令一段时间
func (r *Client) ClientUnpause(ctx context.Context) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientUnpause(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClientID Returns the client ID for the current connection
func (r *Client) ClientID(ctx context.Context) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientID(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ClientUnblock(ctx context.Context, id int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientUnblock(getCtx(ctx), id).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ClientUnblockWithError(ctx context.Context, id int64) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientUnblockWithError(getCtx(ctx), id).Result()
		return err
	}, acceptable)
	return
}

// DBSize 返回当前数据库的 key 的数量
func (r *Client) DBSize(ctx context.Context) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DBSize(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// FlushAll 删除所有数据库的所有key
func (r *Client) FlushAll(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FlushAll(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// FlushAllAsync 异步删除所有数据库的所有key
func (r *Client) FlushAllAsync(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FlushAllAsync(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// FlushDB 删除当前数据库的所有key
func (r *Client) FlushDB(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FlushDB(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// FlushDBAsync 异步删除当前数据库的所有key
func (r *Client) FlushDBAsync(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.FlushDBAsync(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Info 获取 Redis 服务器的各种信息和统计数值
func (r *Client) Info(ctx context.Context, section ...string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Info(getCtx(ctx), section...).Result()
		return err
	}, acceptable)
	return
}

// LastSave 返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示
func (r *Client) LastSave(ctx context.Context) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.LastSave(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Save 异步保存数据到硬盘
func (r *Client) Save(ctx context.Context) (val string, err error) {
	// return getRedis(r).Save(getCtx(ctx))
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Save(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// Shutdown 关闭服务器
func (r *Client) Shutdown(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Shutdown(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ShutdownSave 异步保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownSave(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ShutdownSave(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ShutdownNoSave 不保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownNoSave(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ShutdownNoSave(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// SlaveOf 将当前服务器转变为指定服务器的从属服务器(slave server)
func (r *Client) SlaveOf(ctx context.Context, host, port string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SlaveOf(getCtx(ctx), host, port).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) SlowLogGet(ctx context.Context, num int64) (val []redis.SlowLog, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SlowLogGet(getCtx(ctx), num).Result()
		return err
	}, acceptable)
	return
}

// Time 返回当前服务器时间
func (r *Client) Time(ctx context.Context) (val time.Time, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Time(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) DebugObject(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DebugObject(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

func (r *Client) ReadOnly(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ReadOnly(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) ReadWrite(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ReadWrite(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) MemoryUsage(ctx context.Context, key string, samples ...int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MemoryUsage(getCtx(ctx), r.k(key), samples...).Result()
		return err
	}, acceptable)
	return
}

// Eval 执行 Lua 脚本。
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Eval(getCtx(ctx), script, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalSha 执行 Lua 脚本。
func (r *Client) EvalSha(ctx context.Context, sha1 string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalSha(getCtx(ctx), sha1, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalRO 执行 Lua 脚本。
func (r *Client) EvalRO(ctx context.Context, script string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalRO(getCtx(ctx), script, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// EvalShaRO 执行 Lua 脚本。
func (r *Client) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...any) (val any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.EvalShaRO(getCtx(ctx), sha1, r.ks(keys...), args...).Result()
		return err
	}, acceptable)
	return
}

// ScriptExists 查看指定的脚本是否已经被保存在缓存当中。
func (r *Client) ScriptExists(ctx context.Context, hashes ...string) (val []bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptExists(getCtx(ctx), hashes...).Result()
		return err
	}, acceptable)
	return
}

// ScriptFlush 从脚本缓存中移除所有脚本。
func (r *Client) ScriptFlush(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptFlush(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ScriptKill 杀死当前正在运行的 Lua 脚本。
func (r *Client) ScriptKill(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptKill(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ScriptLoad 将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本。
func (r *Client) ScriptLoad(ctx context.Context, script string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ScriptLoad(getCtx(ctx), script).Result()
		return err
	}, acceptable)
	return
}

// Publish 将信息发送到指定的频道。
func (r *Client) Publish(ctx context.Context, channel string, message any) (val int64, err error) {
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
func (r *Client) SPublish(ctx context.Context, channel string, message any) (val int64, err error) {
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
func (r *Client) PubSubChannels(ctx context.Context, pattern string) (val []string, err error) {
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

func (r *Client) PubSubNumSub(ctx context.Context, channels ...string) (val map[string]int64, err error) {
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
func (r *Client) PubSubShardChannels(ctx context.Context, pattern string) (val []string, err error) {
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
func (r *Client) PubSubShardNumSub(ctx context.Context, channels ...string) (val map[string]int64, err error) {
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

// ClusterSlots 获取集群节点的映射数组
func (r *Client) ClusterSlots(ctx context.Context) (val []redis.ClusterSlot, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSlots(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterNodes Get Cluster config for the node
func (r *Client) ClusterNodes(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterNodes(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterMeet Force a node cluster to handshake with another node
func (r *Client) ClusterMeet(ctx context.Context, host, port string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterMeet(getCtx(ctx), host, port).Result()
		return err
	}, acceptable)
	return
}

// ClusterForget Remove a node from the nodes table
func (r *Client) ClusterForget(ctx context.Context, nodeID string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterForget(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterReplicate Reconfigure a node as a replica of the specified master node
func (r *Client) ClusterReplicate(ctx context.Context, nodeID string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterReplicate(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterResetSoft Reset a Redis Cluster node
func (r *Client) ClusterResetSoft(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterResetSoft(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterResetHard Reset a Redis Cluster node
func (r *Client) ClusterResetHard(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterResetHard(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterInfo Provides info about Redis Cluster node state
func (r *Client) ClusterInfo(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterInfo(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterKeySlot Returns the hash slot of the specified key
func (r *Client) ClusterKeySlot(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterKeySlot(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ClusterGetKeysInSlot Return local key names in the specified hash slot
func (r *Client) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterGetKeysInSlot(getCtx(ctx), slot, count).Result()
		return err
	}, acceptable)
	return
}

// ClusterCountFailureReports Return the number of failure reports active for a given node
func (r *Client) ClusterCountFailureReports(ctx context.Context, nodeID string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterCountFailureReports(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterCountKeysInSlot Return the number of local keys in the specified hash slot
func (r *Client) ClusterCountKeysInSlot(ctx context.Context, slot int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterCountKeysInSlot(getCtx(ctx), slot).Result()
		return err
	}, acceptable)
	return
}

// ClusterDelSlots Set hash slots as unbound in receiving node
func (r *Client) ClusterDelSlots(ctx context.Context, slots ...int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterDelSlots(getCtx(ctx), slots...).Result()
		return err
	}, acceptable)
	return
}

// ClusterDelSlotsRange ->  ClusterDelSlots
func (r *Client) ClusterDelSlotsRange(ctx context.Context, min, max int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterDelSlotsRange(getCtx(ctx), min, max).Result()
		return err
	}, acceptable)
	return
}

// ClusterSaveConfig Forces the node to save cluster state on disk
func (r *Client) ClusterSaveConfig(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSaveConfig(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterSlaves List replica nodes of the specified master node
func (r *Client) ClusterSlaves(ctx context.Context, nodeID string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterSlaves(getCtx(ctx), nodeID).Result()
		return err
	}, acceptable)
	return
}

// ClusterFailover Forces a replica to perform a manual failover of its master.
func (r *Client) ClusterFailover(ctx context.Context) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterFailover(getCtx(ctx)).Result()
		return err
	}, acceptable)
	return
}

// ClusterAddSlots Assign new hash slots to receiving node
func (r *Client) ClusterAddSlots(ctx context.Context, slots ...int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterAddSlots(getCtx(ctx), slots...).Result()
		return err
	}, acceptable)
	return
}

// ClusterAddSlotsRange -> ClusterAddSlots
func (r *Client) ClusterAddSlotsRange(ctx context.Context, min, max int) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClusterAddSlotsRange(getCtx(ctx), min, max).Result()
		return err
	}, acceptable)
	return
}

// GeoAdd 将指定的地理空间位置（纬度、经度、名称）添加到指定的key中
func (r *Client) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoAdd(getCtx(ctx), r.k(key), geoLocation...).Result()
		return err
	}, acceptable)
	return
}

// GeoPos 从key里返回所有给定位置元素的位置（经度和纬度）
func (r *Client) GeoPos(ctx context.Context, key string, members ...string) (val []*redis.GeoPos, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoPos(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// GeoRadius 以给定的经纬度为中心， 找出某一半径内的元素
func (r *Client) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadius(getCtx(ctx), r.k(key), longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusStore -> GeoRadius
func (r *Client) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusStore(getCtx(ctx), r.k(key), longitude, latitude, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMember -> GeoRadius
func (r *Client) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMember(getCtx(ctx), r.k(key), member, query).Result()
		return err
	}, acceptable)
	return
}

// GeoRadiusByMemberStore 找出位于指定范围内的元素，中心点是由给定的位置元素决定
func (r *Client) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoRadiusByMemberStore(getCtx(ctx), r.k(key), member, query).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearch(ctx context.Context, key string, q *redis.GeoSearchQuery) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearch(getCtx(ctx), r.k(key), q).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearchLocation(ctx context.Context, key string, q *redis.GeoSearchLocationQuery) (val []redis.GeoLocation, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchLocation(getCtx(ctx), r.k(key), q).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) GeoSearchStore(ctx context.Context, key, store string, q *redis.GeoSearchStoreQuery) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoSearchStore(getCtx(ctx), r.k(key), store, q).Result()
		return err
	}, acceptable)
	return
}

// GeoDist 返回两个给定位置之间的距离
func (r *Client) GeoDist(ctx context.Context, key string, member1, member2, unit string) (val float64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoDist(getCtx(ctx), r.k(key), member1, member2, unit).Result()
		return err
	}, acceptable)
	return

}

// GeoHash 返回一个或多个位置元素的 Geohash 表示
func (r *Client) GeoHash(ctx context.Context, key string, members ...string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.GeoHash(getCtx(ctx), r.k(key), members...).Result()
		return err
	}, acceptable)
	return
}

// MGetByPipeline gets multiple values from keys,Pipeline is used when
// redis is a cluster,This means higher IO performance
// params: keys ...string
// return: []string, error
// func (r *Client) MGetByPipeline(ctx context.Context, keys ...string) ([]string, error) {
// 	var res []string
// 	if r.ClientType == ClientCluster {
// 		start := time.Now()
// 		pipeLineLen := 100
// 		pipeCount := len(keys)/pipeLineLen + 1
// 		pipes := make([]redis.Pipeliner, pipeCount)
// 		for i := 0; i < pipeCount; i++ {

// 			pipes[i] = getRedis(r).Pipeline()
// 		}
// 		for i, k := range keys {
// 			p := pipes[i%pipeCount]
// 			p.Get(ctx, r.k(k))
// 		}
// 		logger.Logger.Debug("process cost: ", time.Since(start))
// 		start = time.Now()
// 		var wg sync.WaitGroup
// 		var lock sync.Mutex
// 		errors := make(chan error, pipeCount)
// 		for _, p := range pipes {
// 			p := p
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				cmders, err := p.Exec(ctx)
// 				if err != nil {
// 					select {
// 					case errors <- err:
// 					default:
// 					}
// 					return
// 				}
// 				lock.Lock()
// 				defer lock.Unlock()
// 				for _, cmder := range cmders {
// 					result, _ := cmder.(*redis.StringCmd).Result()
// 					res = append(res, result)
// 				}
// 			}()
// 		}
// 		wg.Wait()
// 		logger.Logger.Debug("exec cost: ", time.Since(start))
// 		if len(errors) > 0 {
// 			return nil, <-errors
// 		}
// 		return res, nil
// 	}
// 	vals, err := getRedis(r).MGet(ctx, keys...).Result()
// 	if redis.Nil != err && nil != err {
// 		return nil, err
// 	}
// 	for _, item := range vals {
// 		res = append(res, fmt.Sprintf("%s", item))
// 	}
// 	return res, err
// }
