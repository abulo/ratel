package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Del 删除给定的一个或多个键，返回实际删除的数量；不存在的键被忽略。
func (r *Client) Del(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Del(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// Dump 序列化给定键的值，返回 RDB 格式的字符串；键不存在时返回 redis.Nil。
func (r *Client) Dump(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Dump(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Exists 判断给定键是否存在，存在键数量大于 0 时返回 true；支持一次传入多个键。
func (r *Client) Exists(ctx context.Context, key ...string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		valTmp, err := conn.Exists(getCtx(ctx), key...).Result()
		val = valTmp > 0
		return err
	}, acceptable)
	return
}

// Expire 为 key 设置生存时间 expiration，设置成功返回 true；键不存在时返回 false。
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Expire(getCtx(ctx), key, expiration).Result()
		return err
	}, acceptable)
	return
}

// ExpireAt 以 UNIX 时间戳 tm 设置 key 的过期时刻，成功返回 true。
func (r *Client) ExpireAt(ctx context.Context, key string, tm time.Time) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireAt(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireTime 返回 key 的剩余生存时间；键不存在或未设置过期时间时返回服务端约定的负值。
func (r *Client) ExpireTime(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireTime(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ExpireNX 仅当 key 当前未设置过期时间时设置生存时间 tm。
func (r *Client) ExpireNX(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireNX(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireXX 仅当 key 已设置过期时间时更新生存时间 tm。
func (r *Client) ExpireXX(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireXX(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireGT 仅当 key 当前剩余生存时间大于 tm 时将其更新为 tm。
func (r *Client) ExpireGT(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireGT(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// ExpireLT 仅当 key 当前剩余生存时间小于 tm 时将其更新为 tm。
func (r *Client) ExpireLT(ctx context.Context, key string, tm time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ExpireLT(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// Keys 返回匹配模式 pattern 的全部键；该命令在大库上会阻塞服务器，生产环境应改用 Scan。
func (r *Client) Keys(ctx context.Context, key string) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Keys(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Migrate 将 key 迁移到 host:port 指定的目标实例的 db 数据库，timeout 为迁移操作超时时间。
func (r *Client) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Migrate(getCtx(ctx), host, port, key, db, timeout).Result()
		return err
	}, acceptable)
	return
}

// Move 将 key 移动到当前实例的 db 数据库；目标库已存在同名键时移动失败。
func (r *Client) Move(ctx context.Context, key string, db int) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Move(getCtx(ctx), key, db).Result()
		return err
	}, acceptable)
	return
}

// ObjectFreq 返回 key 对应对象的对数访问频率（LFU 计数器）。
func (r *Client) ObjectFreq(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectFreq(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ObjectRefCount 返回 key 对应对象被引用的次数。
func (r *Client) ObjectRefCount(ctx context.Context, key string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectRefCount(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ObjectEncoding 返回 key 对应对象的内部编码格式。
func (r *Client) ObjectEncoding(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectEncoding(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// ObjectIdleTime 返回 key 对应对象自上次被访问以来的空闲时长。
func (r *Client) ObjectIdleTime(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ObjectIdleTime(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Persist 移除 key 的生存时间使其持久保留，操作成功返回 true。
func (r *Client) Persist(ctx context.Context, key string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Persist(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// PExpire 以毫秒为单位为 key 设置生存时间 expiration。
func (r *Client) PExpire(ctx context.Context, key string, expiration time.Duration) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpire(getCtx(ctx), key, expiration).Result()
		return err
	}, acceptable)
	return
}

// PExpireAt 以毫秒时间戳 tm 设置 key 的过期时刻。
func (r *Client) PExpireAt(ctx context.Context, key string, tm time.Time) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpireAt(getCtx(ctx), key, tm).Result()
		return err
	}, acceptable)
	return
}

// PExpireTime 以毫秒为单位返回 key 的剩余生存时间；键不存在或未设置过期时间时返回服务端约定的负值。
func (r *Client) PExpireTime(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PExpireTime(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// PTTL 以毫秒为单位返回 key 的剩余生存时间；键不存在或未设置过期时间时返回服务端约定的负值。
func (r *Client) PTTL(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.PTTL(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// RandomKey 随机返回当前数据库中的一个键；数据库为空时返回 redis.Nil。
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

// Rename 将 key 重命名为 newkey；原键不存在时报错，目标键已存在时被覆盖。
func (r *Client) Rename(ctx context.Context, key, newkey string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Rename(getCtx(ctx), key, newkey).Result()
		return err
	}, acceptable)
	return
}

// RenameNX 仅当 newkey 不存在时将 key 重命名为 newkey，成功返回 true。
func (r *Client) RenameNX(ctx context.Context, key, newkey string) (val bool, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RenameNX(getCtx(ctx), key, newkey).Result()
		return err
	}, acceptable)
	return
}

// Restore 用 Dump 产生的序列化值重建 key，ttl 为生存时间（0 表示不过期）。
func (r *Client) Restore(ctx context.Context, key string, ttl time.Duration, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Restore(getCtx(ctx), key, ttl, value).Result()
		return err
	}, acceptable)
	return
}

// RestoreReplace 同 Restore，但目标键已存在时先覆盖再写入。
func (r *Client) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.RestoreReplace(getCtx(ctx), key, ttl, value).Result()
		return err
	}, acceptable)
	return
}

// Sort 按 Sort 参数对 key 中的元素排序并返回结果。
func (r *Client) Sort(ctx context.Context, key string, sort *redis.Sort) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Sort(getCtx(ctx), key, sort).Result()
		return err
	}, acceptable)
	return
}

// SortRO 只读模式执行 Sort，排序参数中不允许使用 STORE。
func (r *Client) SortRO(ctx context.Context, key string, sort *redis.Sort) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortRO(getCtx(ctx), key, sort).Result()
		return err
	}, acceptable)
	return
}

// SortStore 按 Sort 参数排序并将结果存入 store 键，返回结果元素数量。
func (r *Client) SortStore(ctx context.Context, key, store string, sort *redis.Sort) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortStore(getCtx(ctx), key, store, sort).Result()
		return err
	}, acceptable)
	return
}

// SortInterfaces 同 Sort，排序结果以 any 切片返回。
func (r *Client) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) (val []any, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.SortInterfaces(getCtx(ctx), key, sort).Result()
		return err
	}, acceptable)
	return
}

// Touch 更新给定键的最后访问时间，返回实际被更新的键数量；不存在的键被忽略。
func (r *Client) Touch(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Touch(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// TTL 返回 key 的剩余生存时间（秒）；键不存在或未设置过期时间时返回服务端约定的负值。
func (r *Client) TTL(ctx context.Context, key string) (val time.Duration, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.TTL(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Type 返回 key 存储值的类型（string/list/set/zset/hash/stream 等）。
func (r *Client) Type(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Type(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// Copy 将 sourceKey 复制到 destKey 所在的 db 数据库；replace 指定目标键已存在时是否覆盖。
func (r *Client) Copy(ctx context.Context, sourceKey, destKey string, db int, replace bool) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Copy(getCtx(ctx), sourceKey, destKey, db, replace).Result()
		return err
	}, acceptable)
	return
}

// Scan 以游标 cursorIn 增量迭代键；match 为匹配模式，count 为建议批次大小，返回本次键与新游标。
func (r *Client) Scan(ctx context.Context, cursorIn uint64, match string, count int64) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.Scan(getCtx(ctx), cursorIn, match, count).Result()
		return err
	}, acceptable)
	return
}

// ScanIterator 返回 Scan 的迭代器，供调用方顺序遍历当前数据库中的键。
func (r *Client) ScanIterator(ctx context.Context, cursorIn uint64, match string, count int64) (val *redis.ScanIterator, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val = conn.Scan(getCtx(ctx), cursorIn, match, count).Iterator()
		return nil
	}, acceptable)
	return
}

// ScanType 增量迭代并返回指定类型 keyType 的键；match 为匹配模式，count 为建议批次大小。
func (r *Client) ScanType(ctx context.Context, cursorIn uint64, match string, count int64, keyType string) (val []string, cursor uint64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, cursor, err = conn.ScanType(getCtx(ctx), cursorIn, match, count, keyType).Result()
		return err
	}, acceptable)
	return
}
