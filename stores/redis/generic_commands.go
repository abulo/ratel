package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Del 删除给定的一个或多个key keys: 要删除的key列表
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

// Dump 序列化给定key key: 要序列化的key
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

// Exists 检查给定key是否存在 key: 要检查的key
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

// Expire 为给定key设置生存时间 key: 要设置的key expiration: 生存时间
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

// ExpireAt 以UNIX时间戳格式设置key的过期时间 key: 要设置的key tm: 过期时间
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

// ExpireTime 返回key的剩余生存时间 key: 要查询的key
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

// ExpireNX 仅在key不存在时设置生存时间 key: 要设置的key tm: 生存时间
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

// ExpireXX 仅在key存在时设置生存时间 key: 要设置的key tm: 生存时间
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

// ExpireGT 仅在当前生存时间大于给定值时设置 key: 要设置的key tm: 生存时间
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

// ExpireLT 仅在当前生存时间小于给定值时设置 key: 要设置的key tm: 生存时间
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

// Keys 查找符合模式的key key: 要匹配的模式
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

// Migrate 将key迁移到其他实例 host: 目标主机 port: 目标端口 key: 要迁移的key db: 目标数据库 timeout: 超时时间
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

// Move 将key移动到其他数据库 key: 要移动的key db: 目标数据库
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

// ObjectFreq 返回key引用对象的频率 key: 要查询的key
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

// ObjectRefCount 返回key的引用计数 key: 要查询的key
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

// ObjectEncoding 返回key的内部编码格式 key: 要查询的key
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

// ObjectIdleTime 返回key的空转时间 key: 要查询的key
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

// Persist 移除key的生存时间 key: 要操作的key
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

// PExpire 以毫秒为单位设置key生存时间 key: 要设置的key expiration: 生存时间
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

// PExpireAt 以毫秒为单位设置key过期时间戳 key: 要设置的key tm: 过期时间
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

// PTTL 以毫秒为单位返回key剩余生存时间 key: 要查询的key
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

// RandomKey 随机返回一个key
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

// Rename 重命名key key: 原key名 newkey: 新key名
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

// RenameNX 仅在newkey不存在时重命名key key: 原key名 newkey: 新key名
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

// Restore 反序列化值并与key关联 key: 目标key ttl: 生存时间 value: 序列化值
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

// RestoreReplace -> Restore
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

// Sort 对key中的元素进行排序 key: 要排序的key sort: 排序参数
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

// SortStore -> Sort
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

// SortInterfaces -> Sort
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

// Touch 更改key的最后访问时间 keys: 要操作的key列表
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

// TTL 返回key的剩余生存时间 key: 要查询的key
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

// Type 返回key存储值的类型 key: 要查询的key
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

// Copy 复制key sourceKey: 源key destKey: 目标key db: 目标数据库 replace: 是否替换
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

// Scan 增量迭代key cursorIn: 起始游标 match: 匹配模式 count: 每次返回数量
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

// ScanIterator 增量迭代key并返回迭代器 cursorIn: 起始游标 match: 匹配模式 count: 每次返回数量
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

// ScanType 增量迭代key并返回指定类型的key cursorIn: 起始游标 match: 匹配模式 count: 每次返回数量 keyType: key类型
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
