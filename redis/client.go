package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

//Config 配置
type Config struct {
	Type      bool     //是否集群
	Hosts     []string //IP
	Password  string   //密码
	Database  int      //数据库
	PoolSize  int      //连接池大小
	KeyPrefix string
}

//New 新连接
func New(config *Config) *Client {
	opts := Options{}
	if config.Type {
		opts.Type = ClientCluster
	} else {
		opts.Type = ClientNormal
	}
	opts.Hosts = config.Hosts
	opts.KeyPrefix = config.KeyPrefix

	if config.PoolSize > 0 {
		opts.PoolSize = config.PoolSize
	} else {
		opts.PoolSize = 64
	}
	if config.Database > 0 {
		opts.Database = config.Database
	} else {
		opts.Database = 0
	}
	if config.Password != "" {
		opts.Password = config.Password
	}
	client := NewClient(opts)
	ctx := context.TODO()
	if err := client.Ping(ctx).Err(); err != nil {
		_logger.Panic(err.Error())
	}
	return client
}

// RedisNil means nil reply, .e.g. when key does not exist.
const RedisNil = redis.Nil

// Client a struct representing the redis client
type Client struct {
	opts      Options
	client    redis.Cmdable
	fmtString string
}

// NewClient 新客户端
func NewClient(opts Options) *Client {
	r := &Client{opts: opts}

	switch opts.Type {
	// 群集客户端
	case ClientCluster:

		tc := redis.NewClusterClient(opts.GetClusterConfig())
		if trace {
			ctx := context.TODO()
			tc.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
				shard.AddHook(OpenTelemetryHook{})
				return nil
			})
		}
		r.client = tc
	// 标准客户端也是默认值
	case ClientNormal:
		fallthrough
	default:
		tc := redis.NewClient(opts.GetNormalConfig())
		if trace {
			tc.AddHook(OpenTelemetryHook{})
		}
		r.client = tc
	}
	r.fmtString = opts.KeyPrefix + "%s"
	return r
}

// IsCluster 判断是否集群
func (r *Client) IsCluster() bool {
	if r.opts.Type == ClientCluster {
		return true
	}
	return false
}

//Prefix 返回前缀+键
func (r *Client) Prefix(key string) string {
	return fmt.Sprintf(r.fmtString, key)
}

// k 格式化并返回带前缀的密钥
func (r *Client) k(key string) string {
	return fmt.Sprintf(r.fmtString, key)
}

// ks 使用前缀格式化并返回一组键
func (r *Client) ks(key ...string) []string {
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = r.k(k)
	}
	return keys
}

// GetClient 返回客户端
func (r *Client) GetClient() redis.Cmdable {
	return r.client
}

// MGetByPipeline gets multiple values from keys,Pipeline is used when
// redis is a cluster,This means higher IO performance
// params: keys ...string
// return: []string, error
func (r *Client) MGetByPipeline(ctx context.Context, keys ...string) ([]string, error) {

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	var res []string

	if r.IsCluster() {
		start := time.Now()
		pipeLineLen := 100
		pipeCount := len(keys)/pipeLineLen + 1
		pipes := make([]redis.Pipeliner, pipeCount)
		for i := 0; i < pipeCount; i++ {
			pipes[i] = r.client.Pipeline()
		}
		for i, k := range keys {
			p := pipes[i%pipeCount]
			p.Get(ctx, r.k(k))
		}
		_logger.Debug("process cost: %v", time.Since(start))
		start = time.Now()
		var wg sync.WaitGroup
		var lock sync.Mutex
		errors := make(chan error, pipeCount)
		for _, p := range pipes {
			p := p
			wg.Add(1)
			go func() {
				defer wg.Done()
				cmders, err := p.Exec(ctx)
				if err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}
				lock.Lock()
				defer lock.Unlock()
				for _, cmder := range cmders {
					result, _ := cmder.(*redis.StringCmd).Result()
					res = append(res, result)
				}
			}()
		}
		wg.Wait()
		_logger.Debug("exec cost: %v", time.Since(start))

		if len(errors) > 0 {
			return nil, <-errors
		}

		return res, nil
	}

	vals, err := r.client.MGet(ctx, keys...).Result()

	if redis.Nil != err && nil != err {
		return nil, err
	}

	for _, item := range vals {
		res = append(res, fmt.Sprintf("%s", item))
	}

	return res, err
}

// Pipeline 获取管道
func (r *Client) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

//Pipelined 管道
func (r *Client) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.Pipelined(ctx, fn)
}

//TxPipeline 获取管道
func (r *Client) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

//TxPipelined 管道
func (r *Client) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.TxPipelined(ctx, fn)
}

//Command 返回有关所有Redis命令的详细信息的Array回复
func (r *Client) Command(ctx context.Context) *redis.CommandsInfoCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Command(ctx)
}

// ClientGetName returns the name of the connection.
func (r *Client) ClientGetName(ctx context.Context) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientGetName(ctx)
}

// Echo  批量字符串回复
func (r *Client) Echo(ctx context.Context, message interface{}) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Echo(ctx, message)
}

// Ping 使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
// 通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值。
// 如果连接正常就返回一个 PONG ，否则返回一个连接错误。
func (r *Client) Ping(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Ping(ctx)
}

//Quit 关闭连接
func (r *Client) Quit(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Quit(ctx)
}

// Del 删除给定的一个或多个 key 。
// 不存在的 key 会被忽略。
func (r *Client) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Del(ctx, r.ks(keys...)...)
}

// Unlink 这个命令非常类似于DEL：它删除指定的键。就像DEL键一样，如果它不存在，它将被忽略。但是，该命令在不同的线程中执行实际的内存回收，所以它不会阻塞，而DEL是。这是命令名称的来源：命令只是将键与键空间断开连接。实际删除将在以后异步发生。
func (r *Client) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Unlink(ctx, r.ks(keys...)...)
}

// Dump 序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (r *Client) Dump(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Dump(ctx, r.k(key))
}

// Exists 检查给定 key 是否存在。
// 若 key 存在，返回 1 ，否则返回 0 。
func (r *Client) Exists(ctx context.Context, key ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Exists(ctx, r.ks(key...)...)
}

// Expire 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
// 设置成功返回 1 。
// 当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Expire(ctx, r.k(key), expiration)
}

// ExpireAt  EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ExpireAt(ctx, r.k(key), tm)
}

// Keys 查找所有符合给定模式 pattern 的 key 。
func (r *Client) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Keys(ctx, r.k(pattern))
}

//Migrate 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
func (r *Client) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Migrate(ctx, host, port, r.k(key), db, timeout)
}

// Move 将当前数据库的 key 移动到给定的数据库 db 当中。
// 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。
func (r *Client) Move(ctx context.Context, key string, db int) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Move(ctx, r.k(key), db)
}

//ObjectRefCount 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func (r *Client) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ObjectRefCount(ctx, r.k(key))
}

//ObjectEncoding 返回给定 key 锁储存的值所使用的内部表示(representation)。
func (r *Client) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ObjectEncoding(ctx, r.k(key))
}

//ObjectIdleTime 返回给定 key 自储存以来的空转时间(idle， 没有被读取也没有被写入)，以秒为单位。
func (r *Client) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ObjectIdleTime(ctx, r.k(key))
}

// Persist 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。
// 当生存时间移除成功时，返回 1 .
// 如果 key 不存在或 key 没有设置生存时间，返回 0 。
func (r *Client) Persist(ctx context.Context, key string) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Persist(ctx, r.k(key))
}

// PExpire 毫秒为单位设置 key 的生存时间
func (r *Client) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PExpire(ctx, r.k(key), expiration)
}

// PExpireAt 这个命令和 expireat 命令类似，但它以毫秒为单位设置 key 的过期 unix 时间戳，而不是像 expireat 那样，以秒为单位。
// 如果生存时间设置成功，返回 1 。 当 key 不存在或没办法设置生存时间时，返回 0
func (r *Client) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PExpireAt(ctx, r.k(key), tm)
}

// PTTL 这个命令类似于 TTL 命令，但它以毫秒为单位返回 key 的剩余生存时间，而不是像 TTL 命令那样，以秒为单位。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以毫秒为单位，返回 key 的剩余生存时间。
func (r *Client) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PTTL(ctx, r.k(key))
}

// RandomKey 从当前数据库中随机返回(不删除)一个 key 。
func (r *Client) RandomKey(ctx context.Context) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RandomKey(ctx)
}

// Rename 将 key 改名为 newkey 。
// 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。
// 当 newkey 已经存在时， RENAME 命令将覆盖旧值。
func (r *Client) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Rename(ctx, r.k(key), r.k(newkey))
}

// RenameNX 当且仅当 newkey 不存在时，将 key 改名为 newkey 。
// 当 key 不存在时，返回一个错误。
func (r *Client) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RenameNX(ctx, r.k(key), r.k(newkey))
}

// Restore 反序列化给定的序列化值，并将它和给定的 key 关联。
// 参数 ttl 以毫秒为单位为 key 设置生存时间；如果 ttl 为 0 ，那么不设置生存时间。
// RESTORE 在执行反序列化之前会先对序列化值的 RDB 版本和数据校验和进行检查，如果 RDB 版本不相同或者数据不完整的话，那么 RESTORE 会拒绝进行反序列化，并返回一个错误。
func (r *Client) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Restore(ctx, r.k(key), ttl, value)
}

// RestoreReplace -> Restore
func (r *Client) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RestoreReplace(ctx, r.k(key), ttl, value)
}

// Sort 返回或保存给定列表、集合、有序集合 key 中经过排序的元素。
// 排序默认以数字作为对象，值被解释为双精度浮点数，然后进行比较。
func (r *Client) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Sort(ctx, r.k(key), sort)
}

//SortStore -> Sort
func (r *Client) SortStore(ctx context.Context, key, store string, sort *redis.Sort) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SortStore(ctx, r.k(key), store, sort)
}

//SortInterfaces -> Sort
func (r *Client) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) *redis.SliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SortInterfaces(ctx, r.k(key), sort)
}

// Touch 更改键的上次访问时间。返回指定的现有键的数量。
func (r *Client) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Touch(ctx, r.ks(keys...)...)
}

// TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以秒为单位，返回 key 的剩余生存时间。
func (r *Client) TTL(ctx context.Context, key string) *redis.DurationCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.TTL(ctx, r.k(key))
}

// Type 返回 key 所储存的值的类型。
func (r *Client) Type(ctx context.Context, key string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Type(ctx, r.k(key))
}

// Scan 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素
func (r *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Scan(ctx, cursor, r.k(match), count)
}

// SScan 详细信息请参考 SCAN 命令。
func (r *Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SScan(ctx, r.k(key), cursor, match, count)
}

// HScan 详细信息请参考 SCAN 命令。
func (r *Client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HScan(ctx, r.k(key), cursor, match, count)
}

// ZScan 详细信息请参考 SCAN 命令。
func (r *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZScan(ctx, r.k(key), cursor, match, count)
}

// Append 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (r *Client) Append(ctx context.Context, key, value string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Append(ctx, r.k(key), value)
}

// BitCount 计算给定字符串中，被设置为 1 的比特位的数量。
// 一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
// start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。
// 不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
func (r *Client) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitCount(ctx, r.k(key), bitCount)
}

// BitOpAnd -> BitCount
func (r *Client) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitOpAnd(ctx, r.k(destKey), r.ks(keys...)...)
}

// BitOpOr -> BitCount
func (r *Client) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitOpOr(ctx, r.k(destKey), r.ks(keys...)...)
}

// BitOpXor -> BitCount
func (r *Client) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitOpXor(ctx, r.k(destKey), r.ks(keys...)...)
}

// BitOpNot -> BitCount
func (r *Client) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitOpXor(ctx, r.k(destKey), r.k(key))
}

// BitPos -> BitCount
func (r *Client) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitPos(ctx, r.k(key), bit, pos...)
}

// BitField -> BitCount
func (r *Client) BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BitField(ctx, r.k(key), args...)
}

// Decr 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 执行 DECR 命令之后 key 的值。
func (r *Client) Decr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Decr(ctx, r.k(key))
}

// DecrBy 将 key 所储存的值减去减量 decrement 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于更多递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 减去 decrement 之后， key 的值。
func (r *Client) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.DecrBy(ctx, r.k(key), value)
}

// Get 返回 key 所关联的字符串值。
// 如果 key 不存在那么返回特殊值 nil 。
// 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。
// 当 key 不存在时，返回 nil ，否则，返回 key 的值。
// 如果 key 不是字符串类型，那么返回一个错误。
func (r *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Get(ctx, r.k(key))
}

// GetBit 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
// 字符串值指定偏移量上的位(bit)。
func (r *Client) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GetBit(ctx, r.k(key), offset)
}

// GetRange 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
// 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
// GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
// 返回截取得出的子字符串。
func (r *Client) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GetRange(ctx, r.k(key), start, end)
}

// GetSet 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 返回给定 key 的旧值。
// 当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
func (r *Client) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GetSet(ctx, r.k(key), value)
}

// Incr 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 执行 INCR 命令之后 key 的值。
func (r *Client) Incr(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Incr(ctx, r.k(key))
}

// IncrBy 将 key 所储存的值加上增量 increment 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，参见 INCR 命令。
// 加上 increment 之后， key 的值。
func (r *Client) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.IncrBy(ctx, r.k(key), value)
}

// IncrByFloat 为 key 中所储存的值加上浮点数增量 increment 。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
// 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
func (r *Client) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.IncrByFloat(ctx, r.k(key), value)
}

// MGet 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
// 一个包含所有给定 key 的值的列表。
func (r *Client) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.MGet(ctx, r.ks(keys...)...)
}

// MSet 同时设置一个或多个 key-value 对。
func (r *Client) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.MSet(ctx, values...)
}

// MSetNX 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
func (r *Client) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.MSetNX(ctx, values...)
}

// Set 将字符串值 value 关联到 key 。
// 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
// 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Set(ctx, r.k(key), value, expiration)
}

// SetBit 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 当 key 不存在时，自动生成一个新的字符串值。
// 字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
// offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
func (r *Client) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SetBit(ctx, r.k(key), offset, value)
}

// SetNX 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
// SETNX 是『SET if Not eXists』(如果不存在，则 SET)的简写。
func (r *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SetNX(ctx, r.k(key), value, expiration)
}

// SetXX -> SetNX
func (r *Client) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SetXX(ctx, r.k(key), value, expiration)
}

// SetRange 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (r *Client) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SetRange(ctx, r.k(key), offset, value)
}

// StrLen 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。
func (r *Client) StrLen(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.StrLen(ctx, r.k(key))
}

// HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func (r *Client) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HDel(ctx, r.k(key), fields...)
}

//HExists 查看哈希表 key 中，给定域 field 是否存在。
func (r *Client) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HExists(ctx, r.k(key), field)
}

// HGet 返回哈希表 key 中给定域 field 的值。
func (r *Client) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HGet(ctx, r.k(key), field)
}

// HGetAll 返回哈希表 key 中，所有的域和值。
// 在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
func (r *Client) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HGetAll(ctx, r.k(key))
}

// HIncrBy 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (r *Client) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HIncrBy(ctx, r.k(key), field, incr)
}

// HIncrByFloat 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HIncrByFloat(ctx, r.k(key), field, incr)
}

// HKeys 返回哈希表 key 中的所有域。
func (r *Client) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HKeys(ctx, r.k(key))
}

//HLen 返回哈希表 key 中域的数量。
func (r *Client) HLen(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HLen(ctx, r.k(key))
}

// HMGet 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HMGet(ctx, r.k(key), fields...)
}

// HSet 将哈希表 key 中的域 field 的值设为 value 。
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Client) HSet(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HSet(ctx, r.k(key), value...)
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (r *Client) HMSet(ctx context.Context, key string, value ...interface{}) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HMSet(ctx, r.k(key), value...)
}

// HSetNX 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
// 若域 field 已经存在，该操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func (r *Client) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HSetNX(ctx, r.k(key), field, value)
}

// HVals 返回哈希表 key 中所有域的值。
func (r *Client) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.HVals(ctx, r.k(key))
}

// BLPop 是列表的阻塞式(blocking)弹出原语。
// 它是 LPop 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPop 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的头元素。
func (r *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BLPop(ctx, timeout, r.ks(keys...)...)
}

// BRPopLPush 是 RPOPLPUSH 的阻塞版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH 一样。
// 当列表 source 为空时， BRPOPLPUSH 命令将阻塞连接，直到等待超时，或有另一个客户端对 source 执行 LPUSH 或 RPUSH 命令为止。
func (r *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BRPopLPush(ctx, r.k(source), r.k(destination), timeout)
}

// LIndex 返回列表 key 中，下标为 index 的元素。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LIndex(ctx, r.k(key), index)
}

// LInsert 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
// 当 pivot 不存在于列表 key 时，不执行任何操作。
// 当 key 不存在时， key 被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LInsert(ctx, r.k(key), op, pivot, value)
}

// LInsertAfter 同 LInsert
func (r *Client) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LInsertAfter(ctx, r.k(key), pivot, value)
}

// LInsertBefore 同 LInsert
func (r *Client) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LInsertBefore(ctx, r.k(key), pivot, value)
}

// LLen 返回列表 key 的长度。
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LLen(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LLen(ctx, r.k(key))
}

// LPop 移除并返回列表 key 的头元素。
func (r *Client) LPop(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LPop(ctx, r.k(key))
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头
// 如果 key 不存在，一个空列表会被创建并执行 LPush 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LPush(ctx, r.k(key), values...)
}

// LPushX 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
// 和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
func (r *Client) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LPushX(ctx, r.k(key), value)
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (r *Client) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LRange(ctx, r.k(key), start, stop)
}

// LRem 根据参数 count 的值，移除列表中与参数 value 相等的元素。
func (r *Client) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LRem(ctx, r.k(key), count, value)
}

// LSet 将列表 key 下标为 index 的元素的值设置为 value 。
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
// 关于列表下标的更多信息，请参考 LINDEX 命令。
func (r *Client) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LSet(ctx, r.k(key), index, value)
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 当 key 不是列表类型时，返回一个错误。
func (r *Client) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LTrim(ctx, r.k(key), start, stop)
}

// BRPop 是列表的阻塞式(blocking)弹出原语。
// 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的尾部元素。
// 关于阻塞操作的更多信息，请查看 BLPOP 命令， BRPOP 除了弹出元素的位置和 BLPOP 不同之外，其他表现一致。
func (r *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BRPop(ctx, timeout, r.ks(keys...)...)
}

// RPopLPush 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
// 将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
// 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
// 举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作。
// 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
func (r *Client) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RPopLPush(ctx, r.k(source), r.k(destination))
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
// 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RPush(ctx, r.k(key), values...)
}

// RPushX 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
// 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。
func (r *Client) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RPushX(ctx, r.k(key), value)
}

// RPop 移除并返回列表 key 的尾元素。
func (r *Client) RPop(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.RPop(ctx, r.k(key))
}

// SAdd 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
// 当 key 不是集合类型时，返回一个错误。
func (r *Client) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SAdd(ctx, r.k(key), members...)
}

// SCard 返回集合 key 的基数(集合中元素的数量)。
func (r *Client) SCard(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SCard(ctx, r.k(key))
}

// SDiff 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (r *Client) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SDiff(ctx, r.ks(keys...)...)
}

// SDiffStore 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SDiffStore(ctx, r.k(destination), r.ks(keys...)...)
}

// SInter 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
// 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func (r *Client) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SInter(ctx, r.ks(keys...)...)
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SInterStore(ctx, r.k(destination), r.ks(keys...)...)
}

// SIsMember 判断 member 元素是否集合 key 的成员。
func (r *Client) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SIsMember(ctx, r.k(key), member)
}

// SMembers 返回集合 key 中的所有成员。
// 不存在的 key 被视为空集合。
func (r *Client) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SMembers(ctx, r.k(key))
}

// SMembersMap -> SMembers
func (r *Client) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SMembersMap(ctx, r.k(key))
}

// SMove 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
func (r *Client) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SMove(ctx, r.k(source), r.k(destination), member)
}

// SPop 移除并返回集合中的一个随机元素。
// 如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
func (r *Client) SPop(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SPop(ctx, r.k(key))
}

// SPopN -> SPop
func (r *Client) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SPopN(ctx, r.k(key), count)
}

// SRandMember 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， SRANDMEMBER 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 SRANDMEMBER 则仅仅返回随机元素，而不对集合进行任何改动。
func (r *Client) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SRandMember(ctx, r.k(key))
}

// SRandMemberN -> SRandMember
func (r *Client) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SRandMemberN(ctx, r.k(key), count)
}

// SRem 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
func (r *Client) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SRem(ctx, r.k(key), members...)
}

// SUnion 返回一个集合的全部成员，该集合是所有给定集合的并集。
// 不存在的 key 被视为空集。
func (r *Client) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SUnion(ctx, r.ks(keys...)...)
}

// SUnionStore 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SUnionStore(ctx, r.k(destination), r.ks(keys...)...)
}

// XAdd 将指定的流条目追加到指定key的流中。 如果key不存在，作为运行这个命令的副作用，将使用流的条目自动创建key。
func (r *Client) XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XAdd(ctx, a)
}

// XDel 从指定流中移除指定的条目，并返回成功删除的条目的数量，在传递的ID不存在的情况下， 返回的数量可能与传递的ID数量不同。
func (r *Client) XDel(ctx context.Context, stream string, ids ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XDel(ctx, stream, ids...)
}

// XLen 返回流中的条目数。如果指定的key不存在，则此命令返回0，就好像该流为空。 但是请注意，与其他的Redis类型不同，零长度流是可能的，所以你应该调用TYPE 或者 EXISTS 来检查一个key是否存在。
// 一旦内部没有任何的条目（例如调用XDEL后），流不会被自动删除，因为可能还存在与其相关联的消费者组。
func (r *Client) XLen(ctx context.Context, stream string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XLen(ctx, stream)
}

// XRange 此命令返回流中满足给定ID范围的条目。范围由最小和最大ID指定。所有ID在指定的两个ID之间或与其中一个ID相等（闭合区间）的条目将会被返回。
func (r *Client) XRange(ctx context.Context, stream, start, stop string) *redis.XMessageSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XRange(ctx, stream, start, stop)
}

// XRangeN -> XRange
func (r *Client) XRangeN(ctx context.Context, stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XRangeN(ctx, stream, start, stop, count)
}

// XRevRange 此命令与XRANGE完全相同，但显著的区别是以相反的顺序返回条目，并以相反的顺序获取开始-结束参数：在XREVRANGE中，你需要先指定结束ID，再指定开始ID，该命令就会从结束ID侧开始生成两个ID之间（或完全相同）的所有元素。
func (r *Client) XRevRange(ctx context.Context, stream string, start, stop string) *redis.XMessageSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XRevRange(ctx, stream, start, stop)
}

// XRevRangeN -> XRevRange
func (r *Client) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XRevRangeN(ctx, stream, start, stop, count)
}

// XRead 从一个或者多个流中读取数据，仅返回ID大于调用者报告的最后接收ID的条目。此命令有一个阻塞选项，用于等待可用的项目，类似于BRPOP或者BZPOPMIN等等。
func (r *Client) XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XRead(ctx, a)
}

//XReadStreams -> XRead
func (r *Client) XReadStreams(ctx context.Context, streams ...string) *redis.XStreamSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XReadStreams(ctx, streams...)
}

// XGroupCreate command
func (r *Client) XGroupCreate(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XGroupCreate(ctx, stream, group, start)
}

// XGroupCreateMkStream command
func (r *Client) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XGroupCreateMkStream(ctx, stream, group, start)
}

// XGroupSetID command
func (r *Client) XGroupSetID(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XGroupSetID(ctx, stream, group, start)
}

// XGroupDestroy command
func (r *Client) XGroupDestroy(ctx context.Context, stream, group string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XGroupDestroy(ctx, stream, group)
}

// XGroupDelConsumer command
func (r *Client) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XGroupDelConsumer(ctx, stream, group, consumer)
}

// XReadGroup command
func (r *Client) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XReadGroup(ctx, a)
}

// XAck command
func (r *Client) XAck(ctx context.Context, stream, group string, ids ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XAck(ctx, stream, group, ids...)
}

// XPending command
func (r *Client) XPending(ctx context.Context, stream, group string) *redis.XPendingCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XPending(ctx, stream, group)
}

// XPendingExt command
func (r *Client) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XPendingExt(ctx, a)
}

// XClaim command
func (r *Client) XClaim(ctx context.Context, a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XClaim(ctx, a)
}

// XClaimJustID command
func (r *Client) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XClaimJustID(ctx, a)
}

// XTrim command
func (r *Client) XTrim(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XTrim(ctx, key, maxLen)
}

// XTrimApprox command
func (r *Client) XTrimApprox(ctx context.Context, key string, maxLen int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XTrimApprox(ctx, key, maxLen)
}

// XInfoGroups command
func (r *Client) XInfoGroups(ctx context.Context, key string) *redis.XInfoGroupsCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.XInfoGroups(ctx, key)
}

// BZPopMax 是有序集合命令 ZPOPMAX带有阻塞功能的版本。
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BZPopMax(ctx, timeout, r.ks(keys...)...)
}

// BZPopMin 是有序集合命令 ZPOPMIN带有阻塞功能的版本。
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BZPopMin(ctx, timeout, r.ks(keys...)...)
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
// 如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。
// score 值可以是整数值或双精度浮点数。
// 如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAdd(ctx, r.k(key), members...)
}

// ZAddNX -> ZAdd
func (r *Client) ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAddNX(ctx, r.k(key), members...)
}

// ZAddXX -> ZAdd
func (r *Client) ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAddXX(ctx, r.k(key), members...)
}

// ZAddCh -> ZAdd
func (r *Client) ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAddCh(ctx, r.k(key), members...)
}

// ZAddNXCh -> ZAdd
func (r *Client) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAddNXCh(ctx, r.k(key), members...)
}

// ZAddXXCh -> ZAdd
func (r *Client) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZAddXXCh(ctx, r.k(key), members...)
}

// ZIncr Redis `ZADD key INCR score member` command.
func (r *Client) ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZIncr(ctx, r.k(key), member)
}

// ZIncrNX Redis `ZADD key NX INCR score member` command.
func (r *Client) ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZIncrNX(ctx, r.k(key), member)
}

// ZIncrXX Redis `ZADD key XX INCR score member` command.
func (r *Client) ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZIncrXX(ctx, r.k(key), member)
}

// ZCard 返回有序集 key 的基数。
func (r *Client) ZCard(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZCard(ctx, r.k(key))
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 关于参数 min 和 max 的详细使用方法，请参考 ZRANGEBYSCORE 命令。
func (r *Client) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZCount(ctx, r.k(key), min, max)
}

//ZLexCount -> ZCount
func (r *Client) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZLexCount(ctx, r.k(key), min, max)
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
// 当 key 不是有序集类型时，返回一个错误。
// score 值可以是整数值或双精度浮点数。
func (r *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZIncrBy(ctx, r.k(key), increment, member)
}

// ZInterStore 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
func (r *Client) ZInterStore(ctx context.Context, key string, store *redis.ZStore) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZInterStore(ctx, r.k(key), store)
}

// ZPopMax 删除并返回有序集合key中的最多count个具有最高得分的成员。
func (r *Client) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZPopMax(ctx, r.k(key), count...)
}

// ZPopMin 删除并返回有序集合key中的最多count个具有最低得分的成员。
func (r *Client) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZPopMin(ctx, r.k(key), count...)
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (r *Client) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRange(ctx, r.k(key), start, stop)
}

// ZRangeWithScores -> ZRange
func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRangeWithScores(ctx, r.k(key), start, stop)
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
func (r *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRangeByScore(ctx, r.k(key), opt)
}

// ZRangeByLex -> ZRangeByScore
func (r *Client) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRangeByLex(ctx, r.k(key), opt)
}

// ZRangeByScoreWithScores -> ZRangeByScore
func (r *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRangeByScoreWithScores(ctx, r.k(key), opt)
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
// 使用 ZREVRANK 命令可以获得成员按 score 值递减(从大到小)排列的排名。
func (r *Client) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRank(ctx, r.k(key), member)
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRem(ctx, r.k(key), members...)
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
func (r *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRemRangeByRank(ctx, r.k(key), start, stop)
}

// ZRemRangeByScore 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 自版本2.1.6开始， score 值等于 min 或 max 的成员也可以不包括在内，详情请参见 ZRANGEBYSCORE 命令。
func (r *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRemRangeByScore(ctx, r.k(key), min, max)
}

//ZRemRangeByLex -> ZRemRangeByScore
func (r *Client) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRemRangeByLex(ctx, r.k(key), min, max)
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func (r *Client) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRange(ctx, r.k(key), start, stop)
}

//ZRevRangeWithScores -> ZRevRange
func (r *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRangeWithScores(ctx, r.k(key), start, stop)
}

// ZRevRangeByScore 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func (r *Client) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRangeByScore(ctx, r.k(key), opt)
}

// ZRevRangeByLex -> ZRevRangeByScore
func (r *Client) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRangeByLex(ctx, r.k(key), opt)
}

// ZRevRangeByScoreWithScores -> ZRevRangeByScore
func (r *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRangeByScoreWithScores(ctx, r.k(key), opt)
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (r *Client) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZRevRank(ctx, r.k(key), member)
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *Client) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZScore(ctx, r.k(key), member)
}

// ZUnionStore 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之 和 。
func (r *Client) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ZUnionStore(ctx, r.k(dest), store)
}

// PFAdd 将指定元素添加到HyperLogLog
func (r *Client) PFAdd(ctx context.Context, key string, els ...interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PFAdd(ctx, r.k(key), els...)
}

// PFCount 返回HyperlogLog观察到的集合的近似基数。
func (r *Client) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PFCount(ctx, r.ks(keys...)...)
}

// PFMerge N个不同的HyperLogLog合并为一个。
func (r *Client) PFMerge(ctx context.Context, dest string, keys ...string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PFMerge(ctx, r.k(dest), r.ks(keys...)...)
}

// BgRewriteAOF 异步重写附加文件
func (r *Client) BgRewriteAOF(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BgRewriteAOF(ctx)
}

// BgSave 将数据集异步保存到磁盘
func (r *Client) BgSave(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.BgSave(ctx)
}

// ClientKill 杀掉客户端的连接
func (r *Client) ClientKill(ctx context.Context, ipPort string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientKill(ctx, ipPort)
}

// ClientKillByFilter is new style synx, while the ClientKill is old
// CLIENT KILL <option> [value] ... <option> [value]
func (r *Client) ClientKillByFilter(ctx context.Context, keys ...string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientKillByFilter(ctx, r.ks(keys...)...)
}

// ClientList 获取客户端连接列表
func (r *Client) ClientList(ctx context.Context) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientList(ctx)
}

// ClientPause 停止处理来自客户端的命令一段时间
func (r *Client) ClientPause(ctx context.Context, dur time.Duration) *redis.BoolCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientPause(ctx, dur)
}

// ClientID Returns the client ID for the current connection
func (r *Client) ClientID(ctx context.Context) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClientID(ctx)
}

// ConfigGet 获取指定配置参数的值
func (r *Client) ConfigGet(ctx context.Context, parameter string) *redis.SliceCmd {
	return r.client.ConfigGet(ctx, parameter)
}

// ConfigResetStat 重置 INFO 命令中的某些统计数据
func (r *Client) ConfigResetStat(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ConfigResetStat(ctx)
}

// ConfigSet 修改 redis 配置参数，无需重启
func (r *Client) ConfigSet(ctx context.Context, parameter, value string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ConfigSet(ctx, parameter, value)
}

// ConfigRewrite 对启动 Redis 服务器时所指定的 redis.conf 配置文件进行改写
func (r *Client) ConfigRewrite(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ConfigRewrite(ctx)
}

// DBSize 返回当前数据库的 key 的数量
func (r *Client) DBSize(ctx context.Context) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.DBSize(ctx)
}

// FlushAll 删除所有数据库的所有key
func (r *Client) FlushAll(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.FlushAll(ctx)
}

// FlushAllAsync 异步删除所有数据库的所有key
func (r *Client) FlushAllAsync(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.FlushAllAsync(ctx)
}

// FlushDB 删除当前数据库的所有key
func (r *Client) FlushDB(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.FlushDB(ctx)
}

// FlushDBAsync 异步删除当前数据库的所有key
func (r *Client) FlushDBAsync(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.FlushDBAsync(ctx)
}

// Info 获取 Redis 服务器的各种信息和统计数值
func (r *Client) Info(ctx context.Context, section ...string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Info(ctx, section...)
}

// LastSave 返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示
func (r *Client) LastSave(ctx context.Context) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.LastSave(ctx)
}

//Save 异步保存数据到硬盘
func (r *Client) Save(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Save(ctx)
}

// Shutdown 关闭服务器
func (r *Client) Shutdown(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Shutdown(ctx)
}

// ShutdownSave 异步保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownSave(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ShutdownSave(ctx)
}

// ShutdownNoSave 不保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownNoSave(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ShutdownNoSave(ctx)
}

// SlaveOf 将当前服务器转变为指定服务器的从属服务器(slave server)
func (r *Client) SlaveOf(ctx context.Context, host, port string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.SlaveOf(ctx, host, port)
}

// Time 返回当前服务器时间
func (r *Client) Time(ctx context.Context) *redis.TimeCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Time(ctx)
}

//Eval 执行 Lua 脚本。
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Eval(ctx, script, r.ks(keys...), args...)
}

//EvalSha 执行 Lua 脚本。
func (r *Client) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.EvalSha(ctx, sha1, r.ks(keys...), args...)
}

// ScriptExists 查看指定的脚本是否已经被保存在缓存当中。
func (r *Client) ScriptExists(ctx context.Context, hashes ...string) *redis.BoolSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ScriptExists(ctx, hashes...)
}

// ScriptFlush 从脚本缓存中移除所有脚本。
func (r *Client) ScriptFlush(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ScriptFlush(ctx)
}

// ScriptKill 杀死当前正在运行的 Lua 脚本。
func (r *Client) ScriptKill(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ScriptKill(ctx)

}

// ScriptLoad 将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本。
func (r *Client) ScriptLoad(ctx context.Context, script string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ScriptLoad(ctx, script)
}

// DebugObject 获取 key 的调试信息
func (r *Client) DebugObject(ctx context.Context, key string) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.DebugObject(ctx, r.k(key))
}

//Publish 将信息发送到指定的频道。
func (r *Client) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.Publish(ctx, r.k(channel), message)
}

//PubSubChannels 订阅一个或多个符合给定模式的频道。
func (r *Client) PubSubChannels(ctx context.Context, pattern string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PubSubChannels(ctx, r.k(pattern))
}

// PubSubNumSub 查看订阅与发布系统状态。
func (r *Client) PubSubNumSub(ctx context.Context, channels ...string) *redis.StringIntMapCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PubSubNumSub(ctx, r.ks(channels...)...)
}

// PubSubNumPat 用于获取redis订阅或者发布信息的状态
func (r *Client) PubSubNumPat(ctx context.Context) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.PubSubNumPat(ctx)
}

//ClusterSlots 获取集群节点的映射数组
func (r *Client) ClusterSlots(ctx context.Context) *redis.ClusterSlotsCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterSlots(ctx)
}

// ClusterNodes Get Cluster config for the node
func (r *Client) ClusterNodes(ctx context.Context) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterNodes(ctx)
}

// ClusterMeet Force a node cluster to handshake with another node
func (r *Client) ClusterMeet(ctx context.Context, host, port string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterMeet(ctx, host, port)
}

// ClusterForget Remove a node from the nodes table
func (r *Client) ClusterForget(ctx context.Context, nodeID string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterForget(ctx, nodeID)
}

// ClusterReplicate Reconfigure a node as a replica of the specified master node
func (r *Client) ClusterReplicate(ctx context.Context, nodeID string) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterReplicate(ctx, nodeID)

}

// ClusterResetSoft Reset a Redis Cluster node
func (r *Client) ClusterResetSoft(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterResetSoft(ctx)
}

// ClusterResetHard Reset a Redis Cluster node
func (r *Client) ClusterResetHard(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterResetHard(ctx)
}

// ClusterInfo Provides info about Redis Cluster node state
func (r *Client) ClusterInfo(ctx context.Context) *redis.StringCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterInfo(ctx)
}

// ClusterKeySlot Returns the hash slot of the specified key
func (r *Client) ClusterKeySlot(ctx context.Context, key string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterKeySlot(ctx, r.k(key))
}

// ClusterGetKeysInSlot Return local key names in the specified hash slot
func (r *Client) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterGetKeysInSlot(ctx, slot, count)
}

// ClusterCountFailureReports Return the number of failure reports active for a given node
func (r *Client) ClusterCountFailureReports(ctx context.Context, nodeID string) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterCountFailureReports(ctx, nodeID)
}

// ClusterCountKeysInSlot Return the number of local keys in the specified hash slot
func (r *Client) ClusterCountKeysInSlot(ctx context.Context, slot int) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterCountKeysInSlot(ctx, slot)
}

// ClusterDelSlots Set hash slots as unbound in receiving node
func (r *Client) ClusterDelSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterDelSlots(ctx, slots...)
}

// ClusterDelSlotsRange ->  ClusterDelSlots
func (r *Client) ClusterDelSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterDelSlotsRange(ctx, min, max)
}

// ClusterSaveConfig Forces the node to save cluster state on disk
func (r *Client) ClusterSaveConfig(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterSaveConfig(ctx)
}

// ClusterSlaves List replica nodes of the specified master node
func (r *Client) ClusterSlaves(ctx context.Context, nodeID string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterSlaves(ctx, nodeID)
}

// ClusterFailover Forces a replica to perform a manual failover of its master.
func (r *Client) ClusterFailover(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterFailover(ctx)
}

// ClusterAddSlots Assign new hash slots to receiving node
func (r *Client) ClusterAddSlots(ctx context.Context, slots ...int) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterAddSlots(ctx, slots...)
}

// ClusterAddSlotsRange -> ClusterAddSlots
func (r *Client) ClusterAddSlotsRange(ctx context.Context, min, max int) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ClusterAddSlotsRange(ctx, min, max)
}

//GeoAdd 将指定的地理空间位置（纬度、经度、名称）添加到指定的key中
func (r *Client) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoAdd(ctx, r.k(key), geoLocation...)
}

// GeoPos 从key里返回所有给定位置元素的位置（经度和纬度）
func (r *Client) GeoPos(ctx context.Context, key string, members ...string) *redis.GeoPosCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoPos(ctx, r.k(key), members...)
}

// GeoRadius 以给定的经纬度为中心， 找出某一半径内的元素
func (r *Client) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoRadius(ctx, r.k(key), longitude, latitude, query)
}

// GeoRadiusStore -> GeoRadius
func (r *Client) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoRadiusStore(ctx, r.k(key), longitude, latitude, query)
}

// GeoRadiusByMember -> GeoRadius
func (r *Client) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoRadiusByMember(ctx, r.k(key), member, query)
}

//GeoRadiusByMemberStore 找出位于指定范围内的元素，中心点是由给定的位置元素决定
func (r *Client) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoRadiusByMemberStore(ctx, r.k(key), member, query)
}

// GeoDist 返回两个给定位置之间的距离
func (r *Client) GeoDist(ctx context.Context, key string, member1, member2, unit string) *redis.FloatCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoDist(ctx, r.k(key), member1, member2, unit)
}

// GeoHash 返回一个或多个位置元素的 Geohash 表示
func (r *Client) GeoHash(ctx context.Context, key string, members ...string) *redis.StringSliceCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.GeoHash(ctx, r.k(key), members...)
}

// ReadOnly Enables read queries for a connection to a cluster replica node
func (r *Client) ReadOnly(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ReadOnly(ctx)
}

// ReadWrite Disables read queries for a connection to a cluster replica node
func (r *Client) ReadWrite(ctx context.Context) *redis.StatusCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.ReadWrite(ctx)
}

// MemoryUsage Estimate the memory usage of a key
func (r *Client) MemoryUsage(ctx context.Context, key string, samples ...int) *redis.IntCmd {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return r.client.MemoryUsage(ctx, r.k(key), samples...)
}

// ErrNotImplemented not implemented error
var ErrNotImplemented = errors.New("Not implemented")
