package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Del 删除给定的一个或多个 key 。
// 不存在的 key 会被忽略。
func (r *Client) Del(ctx context.Context, keys ...any) (val int64, err error) {
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

// Dump 序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (r *Client) Dump(ctx context.Context, key any) (val string, err error) {
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
func (r *Client) Exists(ctx context.Context, key ...any) (val bool, err error) {
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
func (r *Client) Expire(ctx context.Context, key any, expiration time.Duration) (val bool, err error) {
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
func (r *Client) ExpireAt(ctx context.Context, key any, tm time.Time) (val bool, err error) {
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

func (r *Client) ExpireTime(ctx context.Context, key any) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireTime(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// ExpireNX  ExpireNX 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireNX(ctx context.Context, key any, tm time.Duration) (val bool, err error) {
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
func (r *Client) ExpireXX(ctx context.Context, key any, tm time.Duration) (val bool, err error) {
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
func (r *Client) ExpireGT(ctx context.Context, key any, tm time.Duration) (val bool, err error) {
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
func (r *Client) ExpireLT(ctx context.Context, key any, tm time.Duration) (val bool, err error) {
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
func (r *Client) Keys(ctx context.Context, key any) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Keys(getCtx(ctx), r.k(key)).Result()
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
func (r *Client) Move(ctx context.Context, key any, db int) (val bool, err error) {
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
func (r *Client) ObjectRefCount(ctx context.Context, key any) (val int64, err error) {
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
func (r *Client) ObjectEncoding(ctx context.Context, key any) (val string, err error) {
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
func (r *Client) ObjectIdleTime(ctx context.Context, key any) (val time.Duration, err error) {
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
func (r *Client) Persist(ctx context.Context, key any) (val bool, err error) {
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
func (r *Client) PExpire(ctx context.Context, key any, expiration time.Duration) (val bool, err error) {
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
func (r *Client) PExpireAt(ctx context.Context, key any, tm time.Time) (val bool, err error) {
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

func (r *Client) PExpireTime(ctx context.Context, key any) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpireTime(getCtx(ctx), r.k(key)).Result()
		return err
	}, acceptable)
	return
}

// PTTL 这个命令类似于 TTL 命令，但它以毫秒为单位返回 key 的剩余生存时间，而不是像 TTL 命令那样，以秒为单位。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以毫秒为单位，返回 key 的剩余生存时间。
func (r *Client) PTTL(ctx context.Context, key any) (val time.Duration, err error) {
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
func (r *Client) Rename(ctx context.Context, key, newkey any) (val string, err error) {
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
func (r *Client) RenameNX(ctx context.Context, key, newkey any) (val bool, err error) {
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
func (r *Client) Restore(ctx context.Context, key any, ttl time.Duration, value string) (val string, err error) {
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
func (r *Client) RestoreReplace(ctx context.Context, key any, ttl time.Duration, value string) (val string, err error) {
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
func (r *Client) Sort(ctx context.Context, key any, sort *redis.Sort) (val []string, err error) {
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

func (r *Client) SortRO(ctx context.Context, key any, sort *redis.Sort) (val []string, err error) {
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
func (r *Client) SortStore(ctx context.Context, key, store any, sort *redis.Sort) (val int64, err error) {
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
func (r *Client) SortInterfaces(ctx context.Context, key any, sort *redis.Sort) (val []any, err error) {
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
func (r *Client) Touch(ctx context.Context, keys ...any) (val int64, err error) {
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
func (r *Client) TTL(ctx context.Context, key any) (val time.Duration, err error) {
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
func (r *Client) Type(ctx context.Context, key any) (val string, err error) {
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

// Copy
func (r *Client) Copy(ctx context.Context, sourceKey any, destKey any, db int, replace bool) (val int64, err error) {
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

// Scan 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素
func (r *Client) Scan(ctx context.Context, cursorIn uint64, match any, count int64) (val []string, cursor uint64, err error) {
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

func (r *Client) ScanIterator(ctx context.Context, cursorIn uint64, match any, count int64) (val *redis.ScanIterator, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.Scan(getCtx(ctx), cursorIn, r.k(match), count).Iterator()
		return nil
	}, acceptable)
	return
}

func (r *Client) ScanType(ctx context.Context, cursorIn uint64, match any, count int64, keyType string) (val []string, cursor uint64, err error) {
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
