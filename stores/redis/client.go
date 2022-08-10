package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/util"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Config 配置
type Config struct {
	Type          bool     //是否集群
	Hosts         []string //IP
	Password      string   //密码
	Database      int      //数据库
	PoolSize      int      //连接池大小
	KeyPrefix     string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
}

// New 新连接
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
	if !config.DisableMetric {
		opts.DisableMetric = config.DisableMetric
	}
	if !config.DisableTrace {
		opts.DisableTrace = config.DisableTrace
	}
	client := NewClient(opts)
	ctx := context.TODO()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Logger.Panic(err.Error())
	}
	return client
}

// RedisNil means nil reply, .e.g. when key does not exist.
const RedisNil = redis.Nil

// Client a struct representing the redis client
type Client struct {
	opts          Options
	client        *redis.Client
	clusterClient *redis.ClusterClient
	fmtString     string
	clientType    ClientType
}

// NewClient 新客户端
func NewClient(opts Options) *Client {
	r := &Client{opts: opts}
	switch opts.Type {
	// 群集客户端
	case ClientCluster:
		// NewClusterClient 返回一个 Redis 集群客户端
		tc := redis.NewClusterClient(opts.GetClusterConfig())
		ctx := context.TODO()
		if !opts.DisableTrace || !opts.DisableMetric {
			tc.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
				shard.AddHook(OpenTraceHook{
					DisableMetric: opts.DisableMetric,
					DisableTrace:  opts.DisableTrace,
					DB:            opts.Database,
					Addr:          util.Implode(";", opts.Hosts),
				})
				return nil
			})
		}
		r.clusterClient = tc
		r.clientType = ClientCluster
	// 标准客户端也是默认值
	case ClientNormal:
		fallthrough
	default:
		// NewClient 根据 Options 指定的 Redis Server 返回一个客户端。
		tc := redis.NewClient(opts.GetNormalConfig())
		if !opts.DisableTrace || !opts.DisableMetric {
			tc.AddHook(OpenTraceHook{
				DisableMetric: opts.DisableMetric,
				DisableTrace:  opts.DisableTrace,
				DB:            opts.Database,
				Addr:          util.Implode(";", opts.Hosts),
			})
		}
		r.client = tc
		r.clientType = ClientNormal
	}
	r.fmtString = opts.KeyPrefix + "%s"

	return r
}

// // IsCluster 判断是否集群
// func (r *Client) IsCluster() bool {
// 	if r.opts.Type == ClientCluster {
// 		return true
// 	}
// 	return false
// }

// Prefix 返回前缀+键
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
func (r *Client) GetClient() (res redis.Cmdable) {
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient
	case ClientNormal:
		res = r.client
	}
	return res
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

	if r.clientType == ClientCluster {
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
		logger.Logger.Debug("process cost: ", time.Since(start))
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
		logger.Logger.Debug("exec cost: ", time.Since(start))
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
func (r *Client) Pipeline() (res redis.Pipeliner) {
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Pipeline()
	case ClientNormal:
		res = r.client.Pipeline()
	}
	return res
}

// Pipelined 管道
func (r *Client) Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) (res []redis.Cmder, err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res, err = r.clusterClient.Pipelined(ctx, fn)
	case ClientNormal:
		res, err = r.client.Pipelined(ctx, fn)
	}
	return res, err
}

// TxPipeline 获取管道
func (r *Client) TxPipeline() (res redis.Pipeliner) {
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.TxPipeline()
	case ClientNormal:
		res = r.client.TxPipeline()
	}
	return res
}

// TxPipelined 管道
func (r *Client) TxPipelined(ctx context.Context, fn func(redis.Pipeliner) error) (res []redis.Cmder, err error) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res, err = r.clusterClient.TxPipelined(ctx, fn)
	case ClientNormal:
		res, err = r.client.TxPipelined(ctx, fn)
	}
	return res, err
}

// Command 返回有关所有Redis命令的详细信息的Array回复
func (r *Client) Command(ctx context.Context) (res *redis.CommandsInfoCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Command(ctx)
	case ClientNormal:
		res = r.client.Command(ctx)
	}
	return res
}

// ClientGetName returns the name of the connection.
func (r *Client) ClientGetName(ctx context.Context) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientGetName(ctx)
	case ClientNormal:
		res = r.client.ClientGetName(ctx)
	}
	return res
}

// Echo  批量字符串回复
func (r *Client) Echo(ctx context.Context, message interface{}) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Echo(ctx, message)
	case ClientNormal:
		res = r.client.Echo(ctx, message)
	}
	return res
}

// Ping 使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
// 通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值。
// 如果连接正常就返回一个 PONG ，否则返回一个连接错误。
func (r *Client) Ping(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Ping(ctx)
	case ClientNormal:
		res = r.client.Ping(ctx)
	}
	return res
}

// Quit 关闭连接
func (r *Client) Quit(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Quit(ctx)
	case ClientNormal:
		res = r.client.Quit(ctx)
	}
	return res
}

// Del 删除给定的一个或多个 key 。
// 不存在的 key 会被忽略。
func (r *Client) Del(ctx context.Context, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Del(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.Del(ctx, r.ks(keys...)...)
	}
	return res
}

// Unlink 这个命令非常类似于DEL：它删除指定的键。就像DEL键一样，如果它不存在，它将被忽略。但是，该命令在不同的线程中执行实际的内存回收，所以它不会阻塞，而DEL是。这是命令名称的来源：命令只是将键与键空间断开连接。实际删除将在以后异步发生。
func (r *Client) Unlink(ctx context.Context, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Unlink(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.Unlink(ctx, r.ks(keys...)...)
	}
	return res
}

// Dump 序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
// 如果 key 不存在，那么返回 nil 。
// 否则，返回序列化之后的值。
func (r *Client) Dump(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Dump(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Dump(ctx, r.k(key))
	}
	return res
}

// Exists 检查给定 key 是否存在。
// 若 key 存在，返回 1 ，否则返回 0 。
func (r *Client) Exists(ctx context.Context, key ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Exists(ctx, r.ks(key...)...)
	case ClientNormal:
		res = r.client.Exists(ctx, r.ks(key...)...)
	}
	return res
}

// Expire 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
// 设置成功返回 1 。
// 当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Expire(ctx, r.k(key), expiration)
	case ClientNormal:
		res = r.client.Expire(ctx, r.k(key), expiration)
	}
	return res
}

// ExpireAt  EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
// 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间
func (r *Client) ExpireAt(ctx context.Context, key string, tm time.Time) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ExpireAt(ctx, r.k(key), tm)
	case ClientNormal:
		res = r.client.ExpireAt(ctx, r.k(key), tm)
	}
	return res
}

// Keys 查找所有符合给定模式 pattern 的 key 。
func (r *Client) Keys(ctx context.Context, pattern string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Keys(ctx, r.k(pattern))
	case ClientNormal:
		res = r.client.Keys(ctx, r.k(pattern))
	}
	return res
}

// Migrate 将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
func (r *Client) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Migrate(ctx, host, port, r.k(key), db, timeout)
	case ClientNormal:
		res = r.client.Migrate(ctx, host, port, r.k(key), db, timeout)
	}
	return res
}

// Move 将当前数据库的 key 移动到给定的数据库 db 当中。
// 如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。
func (r *Client) Move(ctx context.Context, key string, db int) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Move(ctx, r.k(key), db)
	case ClientNormal:
		res = r.client.Move(ctx, r.k(key), db)
	}
	return res
}

// ObjectRefCount 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func (r *Client) ObjectRefCount(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ObjectRefCount(ctx, r.k(key))
	case ClientNormal:
		res = r.client.ObjectRefCount(ctx, r.k(key))
	}
	return res
}

// ObjectEncoding 返回给定 key 锁储存的值所使用的内部表示(representation)。
func (r *Client) ObjectEncoding(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ObjectEncoding(ctx, r.k(key))
	case ClientNormal:
		res = r.client.ObjectEncoding(ctx, r.k(key))
	}
	return res
}

// ObjectIdleTime 返回给定 key 自储存以来的空转时间(idle， 没有被读取也没有被写入)，以秒为单位。
func (r *Client) ObjectIdleTime(ctx context.Context, key string) (res *redis.DurationCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ObjectIdleTime(ctx, r.k(key))
	case ClientNormal:
		res = r.client.ObjectIdleTime(ctx, r.k(key))
	}
	return res
}

// Persist 移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。
// 当生存时间移除成功时，返回 1 .
// 如果 key 不存在或 key 没有设置生存时间，返回 0 。
func (r *Client) Persist(ctx context.Context, key string) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Persist(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Persist(ctx, r.k(key))
	}
	return res
}

// PExpire 毫秒为单位设置 key 的生存时间
func (r *Client) PExpire(ctx context.Context, key string, expiration time.Duration) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PExpire(ctx, r.k(key), expiration)
	case ClientNormal:
		res = r.client.PExpire(ctx, r.k(key), expiration)
	}
	return res
}

// PExpireAt 这个命令和 expireat 命令类似，但它以毫秒为单位设置 key 的过期 unix 时间戳，而不是像 expireat 那样，以秒为单位。
// 如果生存时间设置成功，返回 1 。 当 key 不存在或没办法设置生存时间时，返回 0
func (r *Client) PExpireAt(ctx context.Context, key string, tm time.Time) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PExpireAt(ctx, r.k(key), tm)
	case ClientNormal:
		res = r.client.PExpireAt(ctx, r.k(key), tm)
	}
	return res
}

// PTTL 这个命令类似于 TTL 命令，但它以毫秒为单位返回 key 的剩余生存时间，而不是像 TTL 命令那样，以秒为单位。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以毫秒为单位，返回 key 的剩余生存时间。
func (r *Client) PTTL(ctx context.Context, key string) (res *redis.DurationCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PTTL(ctx, r.k(key))
	case ClientNormal:
		res = r.client.PTTL(ctx, r.k(key))
	}
	return res
}

// RandomKey 从当前数据库中随机返回(不删除)一个 key 。
func (r *Client) RandomKey(ctx context.Context) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RandomKey(ctx)
	case ClientNormal:
		res = r.client.RandomKey(ctx)
	}
	return res
}

// Rename 将 key 改名为 newkey 。
// 当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。
// 当 newkey 已经存在时， RENAME 命令将覆盖旧值。
func (r *Client) Rename(ctx context.Context, key, newkey string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Rename(ctx, r.k(key), r.k(newkey))
	case ClientNormal:
		res = r.client.Rename(ctx, r.k(key), r.k(newkey))
	}
	return res
}

// RenameNX 当且仅当 newkey 不存在时，将 key 改名为 newkey 。
// 当 key 不存在时，返回一个错误。
func (r *Client) RenameNX(ctx context.Context, key, newkey string) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RenameNX(ctx, r.k(key), r.k(newkey))
	case ClientNormal:
		res = r.client.RenameNX(ctx, r.k(key), r.k(newkey))
	}
	return res
}

// Restore 反序列化给定的序列化值，并将它和给定的 key 关联。
// 参数 ttl 以毫秒为单位为 key 设置生存时间；如果 ttl 为 0 ，那么不设置生存时间。
// RESTORE 在执行反序列化之前会先对序列化值的 RDB 版本和数据校验和进行检查，如果 RDB 版本不相同或者数据不完整的话，那么 RESTORE 会拒绝进行反序列化，并返回一个错误。
func (r *Client) Restore(ctx context.Context, key string, ttl time.Duration, value string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Restore(ctx, r.k(key), ttl, value)
	case ClientNormal:
		res = r.client.Restore(ctx, r.k(key), ttl, value)
	}
	return res
}

// RestoreReplace -> Restore
func (r *Client) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RestoreReplace(ctx, r.k(key), ttl, value)
	case ClientNormal:
		res = r.client.RestoreReplace(ctx, r.k(key), ttl, value)
	}
	return res
}

// Sort 返回或保存给定列表、集合、有序集合 key 中经过排序的元素。
// 排序默认以数字作为对象，值被解释为双精度浮点数，然后进行比较。
func (r *Client) Sort(ctx context.Context, key string, sort *redis.Sort) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Sort(ctx, r.k(key), sort)
	case ClientNormal:
		res = r.client.Sort(ctx, r.k(key), sort)
	}
	return res
}

// SortStore -> Sort
func (r *Client) SortStore(ctx context.Context, key, store string, sort *redis.Sort) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SortStore(ctx, r.k(key), store, sort)
	case ClientNormal:
		res = r.client.SortStore(ctx, r.k(key), store, sort)
	}
	return res
}

// SortInterfaces -> Sort
func (r *Client) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) (res *redis.SliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SortInterfaces(ctx, r.k(key), sort)
	case ClientNormal:
		res = r.client.SortInterfaces(ctx, r.k(key), sort)
	}
	return res
}

// Touch 更改键的上次访问时间。返回指定的现有键的数量。
func (r *Client) Touch(ctx context.Context, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Touch(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.Touch(ctx, r.ks(keys...)...)
	}
	return res
}

// TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
// 当 key 不存在时，返回 -2 。
// 当 key 存在但没有设置剩余生存时间时，返回 -1 。
// 否则，以秒为单位，返回 key 的剩余生存时间。
func (r *Client) TTL(ctx context.Context, key string) (res *redis.DurationCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.TTL(ctx, r.k(key))
	case ClientNormal:
		res = r.client.TTL(ctx, r.k(key))
	}
	return res
}

// Type 返回 key 所储存的值的类型。
func (r *Client) Type(ctx context.Context, key string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Type(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Type(ctx, r.k(key))
	}
	return res
}

// Scan 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素
func (r *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) (res *redis.ScanCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Scan(ctx, cursor, r.k(match), count)
	case ClientNormal:
		res = r.client.Scan(ctx, cursor, r.k(match), count)
	}
	return res
}

// SScan 详细信息请参考 SCAN 命令。
func (r *Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) (res *redis.ScanCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SScan(ctx, r.k(key), cursor, match, count)
	case ClientNormal:
		res = r.client.SScan(ctx, r.k(key), cursor, match, count)
	}
	return res
}

// HScan 详细信息请参考 SCAN 命令。
func (r *Client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (res *redis.ScanCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HScan(ctx, r.k(key), cursor, match, count)
	case ClientNormal:
		res = r.client.HScan(ctx, r.k(key), cursor, match, count)
	}
	return res
}

// ZScan 详细信息请参考 SCAN 命令。
func (r *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) (res *redis.ScanCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZScan(ctx, r.k(key), cursor, match, count)
	case ClientNormal:
		res = r.client.ZScan(ctx, r.k(key), cursor, match, count)
	}
	return res
}

// Append 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (r *Client) Append(ctx context.Context, key, value string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Append(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.Append(ctx, r.k(key), value)
	}
	return res
}

// BitCount 计算给定字符串中，被设置为 1 的比特位的数量。
// 一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
// start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。
// 不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
func (r *Client) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitCount(ctx, r.k(key), bitCount)
	case ClientNormal:
		res = r.client.BitCount(ctx, r.k(key), bitCount)
	}
	return res
}

// BitOpAnd -> BitCount
func (r *Client) BitOpAnd(ctx context.Context, destKey string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitOpAnd(ctx, r.k(destKey), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BitOpAnd(ctx, r.k(destKey), r.ks(keys...)...)
	}
	return res
}

// BitOpOr -> BitCount
func (r *Client) BitOpOr(ctx context.Context, destKey string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitOpOr(ctx, r.k(destKey), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BitOpOr(ctx, r.k(destKey), r.ks(keys...)...)
	}
	return res
}

// BitOpXor -> BitCount
func (r *Client) BitOpXor(ctx context.Context, destKey string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitOpXor(ctx, r.k(destKey), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BitOpXor(ctx, r.k(destKey), r.ks(keys...)...)
	}
	return res
}

// BitOpNot -> BitCount
func (r *Client) BitOpNot(ctx context.Context, destKey string, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitOpXor(ctx, r.k(destKey), r.k(key))
	case ClientNormal:
		res = r.client.BitOpXor(ctx, r.k(destKey), r.k(key))
	}
	return res
}

// BitPos -> BitCount
func (r *Client) BitPos(ctx context.Context, key string, bit int64, pos ...int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitPos(ctx, r.k(key), bit, pos...)
	case ClientNormal:
		res = r.client.BitPos(ctx, r.k(key), bit, pos...)
	}
	return res
}

// BitField -> BitCount
func (r *Client) BitField(ctx context.Context, key string, args ...interface{}) (res *redis.IntSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BitField(ctx, r.k(key), args...)
	case ClientNormal:
		res = r.client.BitField(ctx, r.k(key), args...)
	}
	return res
}

// Decr 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 执行 DECR 命令之后 key 的值。
func (r *Client) Decr(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Decr(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Decr(ctx, r.k(key))
	}
	return res
}

// DecrBy 将 key 所储存的值减去减量 decrement 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于更多递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
// 减去 decrement 之后， key 的值。
func (r *Client) DecrBy(ctx context.Context, key string, value int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.DecrBy(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.DecrBy(ctx, r.k(key), value)
	}
	return res
}

// Get 返回 key 所关联的字符串值。
// 如果 key 不存在那么返回特殊值 nil 。
// 假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。
// 当 key 不存在时，返回 nil ，否则，返回 key 的值。
// 如果 key 不是字符串类型，那么返回一个错误。
func (r *Client) Get(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Get(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Get(ctx, r.k(key))
	}
	return res
}

// GetBit 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
// 字符串值指定偏移量上的位(bit)。
func (r *Client) GetBit(ctx context.Context, key string, offset int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GetBit(ctx, r.k(key), offset)
	case ClientNormal:
		res = r.client.GetBit(ctx, r.k(key), offset)
	}
	return res
}

// GetRange 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
// 负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
// GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
// 返回截取得出的子字符串。
func (r *Client) GetRange(ctx context.Context, key string, start, end int64) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GetRange(ctx, r.k(key), start, end)
	case ClientNormal:
		res = r.client.GetRange(ctx, r.k(key), start, end)
	}
	return res
}

// GetSet 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
// 返回给定 key 的旧值。
// 当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
func (r *Client) GetSet(ctx context.Context, key string, value interface{}) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GetSet(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.GetSet(ctx, r.k(key), value)
	}
	return res
}

// Incr 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 执行 INCR 命令之后 key 的值。
func (r *Client) Incr(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Incr(ctx, r.k(key))
	case ClientNormal:
		res = r.client.Incr(ctx, r.k(key))
	}
	return res
}

// IncrBy 将 key 所储存的值加上增量 increment 。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
// 关于递增(increment) / 递减(decrement)操作的更多信息，参见 INCR 命令。
// 加上 increment 之后， key 的值。
func (r *Client) IncrBy(ctx context.Context, key string, value int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.IncrBy(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.IncrBy(ctx, r.k(key), value)
	}
	return res
}

// IncrByFloat 为 key 中所储存的值加上浮点数增量 increment 。
// 如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
// 如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
func (r *Client) IncrByFloat(ctx context.Context, key string, value float64) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.IncrByFloat(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.IncrByFloat(ctx, r.k(key), value)
	}
	return res
}

// MGet 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
// 一个包含所有给定 key 的值的列表。
func (r *Client) MGet(ctx context.Context, keys ...string) (res *redis.SliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.MGet(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.MGet(ctx, r.ks(keys...)...)
	}
	return res
}

// MSet 同时设置一个或多个 key-value 对。
func (r *Client) MSet(ctx context.Context, values ...interface{}) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.MSet(ctx, values...)
	case ClientNormal:
		res = r.client.MSet(ctx, values...)
	}
	return res
}

// MSetNX 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
func (r *Client) MSetNX(ctx context.Context, values ...interface{}) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.MSetNX(ctx, values...)
	case ClientNormal:
		res = r.client.MSetNX(ctx, values...)
	}
	return res
}

// Set 将字符串值 value 关联到 key 。
// 如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
// 对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Set(ctx, r.k(key), value, expiration)
	case ClientNormal:
		res = r.client.Set(ctx, r.k(key), value, expiration)
	}
	return res
}

// SetEX ...
func (r *Client) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SetEX(ctx, r.k(key), value, expiration)
	case ClientNormal:
		res = r.client.SetEX(ctx, r.k(key), value, expiration)
	}
	return res
}

// SetBit 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 当 key 不存在时，自动生成一个新的字符串值。
// 字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
// offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
func (r *Client) SetBit(ctx context.Context, key string, offset int64, value int) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SetBit(ctx, r.k(key), offset, value)
	case ClientNormal:
		res = r.client.SetBit(ctx, r.k(key), offset, value)
	}
	return res
}

// SetNX 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
// SETNX 是『SET if Not eXists』(如果不存在，则 SET)的简写。
func (r *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SetNX(ctx, r.k(key), value, expiration)
	case ClientNormal:
		res = r.client.SetNX(ctx, r.k(key), value, expiration)
	}
	return res
}

// SetXX -> SetNX
func (r *Client) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SetXX(ctx, r.k(key), value, expiration)
	case ClientNormal:
		res = r.client.SetXX(ctx, r.k(key), value, expiration)
	}
	return res
}

// SetRange 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (r *Client) SetRange(ctx context.Context, key string, offset int64, value string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SetRange(ctx, r.k(key), offset, value)
	case ClientNormal:
		res = r.client.SetRange(ctx, r.k(key), offset, value)
	}
	return res
}

// StrLen 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。
func (r *Client) StrLen(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.StrLen(ctx, r.k(key))
	case ClientNormal:
		res = r.client.StrLen(ctx, r.k(key))
	}
	return res
}

// HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func (r *Client) HDel(ctx context.Context, key string, fields ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HDel(ctx, r.k(key), fields...)
	case ClientNormal:
		res = r.client.HDel(ctx, r.k(key), fields...)
	}
	return res
}

// HExists 查看哈希表 key 中，给定域 field 是否存在。
func (r *Client) HExists(ctx context.Context, key, field string) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HExists(ctx, r.k(key), field)
	case ClientNormal:
		res = r.client.HExists(ctx, r.k(key), field)
	}
	return res
}

// HGet 返回哈希表 key 中给定域 field 的值。
func (r *Client) HGet(ctx context.Context, key, field string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HGet(ctx, r.k(key), field)
	case ClientNormal:
		res = r.client.HGet(ctx, r.k(key), field)
	}
	return res
}

// HGetAll 返回哈希表 key 中，所有的域和值。
// 在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
func (r *Client) HGetAll(ctx context.Context, key string) (res *redis.StringStringMapCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HGetAll(ctx, r.k(key))
	case ClientNormal:
		res = r.client.HGetAll(ctx, r.k(key))
	}
	return res
}

// HIncrBy 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
// 对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
// 本操作的值被限制在 64 位(bit)有符号数字表示之内。
func (r *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HIncrBy(ctx, r.k(key), field, incr)
	case ClientNormal:
		res = r.client.HIncrBy(ctx, r.k(key), field, incr)
	}
	return res
}

// HIncrByFloat 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (r *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HIncrByFloat(ctx, r.k(key), field, incr)
	case ClientNormal:
		res = r.client.HIncrByFloat(ctx, r.k(key), field, incr)
	}
	return res
}

// HKeys 返回哈希表 key 中的所有域。
func (r *Client) HKeys(ctx context.Context, key string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HKeys(ctx, r.k(key))
	case ClientNormal:
		res = r.client.HKeys(ctx, r.k(key))
	}
	return res
}

// HLen 返回哈希表 key 中域的数量。
func (r *Client) HLen(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HLen(ctx, r.k(key))
	case ClientNormal:
		res = r.client.HLen(ctx, r.k(key))
	}
	return res
}

// HMGet 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (r *Client) HMGet(ctx context.Context, key string, fields ...string) (res *redis.SliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HMGet(ctx, r.k(key), fields...)
	case ClientNormal:
		res = r.client.HMGet(ctx, r.k(key), fields...)
	}
	return res
}

// HSet 将哈希表 key 中的域 field 的值设为 value 。
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (r *Client) HSet(ctx context.Context, key string, value ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HSet(ctx, r.k(key), value...)
	case ClientNormal:
		res = r.client.HSet(ctx, r.k(key), value...)
	}
	return res
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (r *Client) HMSet(ctx context.Context, key string, value ...interface{}) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HMSet(ctx, r.k(key), value...)
	case ClientNormal:
		res = r.client.HMSet(ctx, r.k(key), value...)
	}
	return res
}

// HSetNX 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
// 若域 field 已经存在，该操作无效。
// 如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
func (r *Client) HSetNX(ctx context.Context, key, field string, value interface{}) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HSetNX(ctx, r.k(key), field, value)
	case ClientNormal:
		res = r.client.HSetNX(ctx, r.k(key), field, value)
	}
	return res
}

// HVals 返回哈希表 key 中所有域的值。
func (r *Client) HVals(ctx context.Context, key string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.HVals(ctx, r.k(key))
	case ClientNormal:
		res = r.client.HVals(ctx, r.k(key))
	}
	return res
}

// BLPop 是列表的阻塞式(blocking)弹出原语。
// 它是 LPop 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPop 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的头元素。
func (r *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BLPop(ctx, timeout, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BLPop(ctx, timeout, r.ks(keys...)...)
	}
	return res
}

// BRPopLPush 是 RPOPLPUSH 的阻塞版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH 一样。
// 当列表 source 为空时， BRPOPLPUSH 命令将阻塞连接，直到等待超时，或有另一个客户端对 source 执行 LPUSH 或 RPUSH 命令为止。
func (r *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BRPopLPush(ctx, r.k(source), r.k(destination), timeout)
	case ClientNormal:
		res = r.client.BRPopLPush(ctx, r.k(source), r.k(destination), timeout)
	}
	return res
}

// LIndex 返回列表 key 中，下标为 index 的元素。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LIndex(ctx context.Context, key string, index int64) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LIndex(ctx, r.k(key), index)
	case ClientNormal:
		res = r.client.LIndex(ctx, r.k(key), index)
	}
	return res
}

// LInsert 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
// 当 pivot 不存在于列表 key 时，不执行任何操作。
// 当 key 不存在时， key 被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LInsert(ctx, r.k(key), op, pivot, value)
	case ClientNormal:
		res = r.client.LInsert(ctx, r.k(key), op, pivot, value)
	}
	return res
}

// LInsertAfter 同 LInsert
func (r *Client) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LInsertAfter(ctx, r.k(key), pivot, value)
	case ClientNormal:
		res = r.client.LInsertAfter(ctx, r.k(key), pivot, value)
	}
	return res
}

// LInsertBefore 同 LInsert
func (r *Client) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LInsertBefore(ctx, r.k(key), pivot, value)
	case ClientNormal:
		res = r.client.LInsertBefore(ctx, r.k(key), pivot, value)
	}
	return res
}

// LLen 返回列表 key 的长度。
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
// 如果 key 不是列表类型，返回一个错误。
func (r *Client) LLen(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LLen(ctx, r.k(key))
	case ClientNormal:
		res = r.client.LLen(ctx, r.k(key))
	}
	return res
}

// LPop 移除并返回列表 key 的头元素。
func (r *Client) LPop(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LPop(ctx, r.k(key))
	case ClientNormal:
		res = r.client.LPop(ctx, r.k(key))
	}
	return res
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头
// 如果 key 不存在，一个空列表会被创建并执行 LPush 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) LPush(ctx context.Context, key string, values ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LPush(ctx, r.k(key), values...)
	case ClientNormal:
		res = r.client.LPush(ctx, r.k(key), values...)
	}
	return res
}

// LPushX 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
// 和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
func (r *Client) LPushX(ctx context.Context, key string, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LPushX(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.LPushX(ctx, r.k(key), value)
	}
	return res
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
func (r *Client) LRange(ctx context.Context, key string, start, stop int64) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LRange(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.LRange(ctx, r.k(key), start, stop)
	}
	return res
}

// LRem 根据参数 count 的值，移除列表中与参数 value 相等的元素。
func (r *Client) LRem(ctx context.Context, key string, count int64, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LRem(ctx, r.k(key), count, value)
	case ClientNormal:
		res = r.client.LRem(ctx, r.k(key), count, value)
	}
	return res
}

// LSet 将列表 key 下标为 index 的元素的值设置为 value 。
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
// 关于列表下标的更多信息，请参考 LINDEX 命令。
func (r *Client) LSet(ctx context.Context, key string, index int64, value interface{}) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LSet(ctx, r.k(key), index, value)
	case ClientNormal:
		res = r.client.LSet(ctx, r.k(key), index, value)
	}
	return res
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 当 key 不是列表类型时，返回一个错误。
func (r *Client) LTrim(ctx context.Context, key string, start, stop int64) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LTrim(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.LTrim(ctx, r.k(key), start, stop)
	}
	return res
}

// BRPop 是列表的阻塞式(blocking)弹出原语。
// 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
// 当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的尾部元素。
// 关于阻塞操作的更多信息，请查看 BLPOP 命令， BRPOP 除了弹出元素的位置和 BLPOP 不同之外，其他表现一致。
func (r *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BRPop(ctx, timeout, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BRPop(ctx, timeout, r.ks(keys...)...)
	}
	return res
}

// RPopLPush 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
// 将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
// 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
// 举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作。
// 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
func (r *Client) RPopLPush(ctx context.Context, source, destination string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RPopLPush(ctx, r.k(source), r.k(destination))
	case ClientNormal:
		res = r.client.RPopLPush(ctx, r.k(source), r.k(destination))
	}
	return res
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
// 如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
// 当 key 存在但不是列表类型时，返回一个错误。
func (r *Client) RPush(ctx context.Context, key string, values ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RPush(ctx, r.k(key), values...)
	case ClientNormal:
		res = r.client.RPush(ctx, r.k(key), values...)
	}
	return res
}

// RPushX 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
// 和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。
func (r *Client) RPushX(ctx context.Context, key string, value interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RPushX(ctx, r.k(key), value)
	case ClientNormal:
		res = r.client.RPushX(ctx, r.k(key), value)
	}
	return res
}

// RPop 移除并返回列表 key 的尾元素。
func (r *Client) RPop(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.RPop(ctx, r.k(key))
	case ClientNormal:
		res = r.client.RPop(ctx, r.k(key))
	}
	return res
}

// SAdd 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
// 当 key 不是集合类型时，返回一个错误。
func (r *Client) SAdd(ctx context.Context, key string, members ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SAdd(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.SAdd(ctx, r.k(key), members...)
	}
	return res
}

// SCard 返回集合 key 的基数(集合中元素的数量)。
func (r *Client) SCard(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SCard(ctx, r.k(key))
	case ClientNormal:
		res = r.client.SCard(ctx, r.k(key))
	}
	return res
}

// SDiff 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (r *Client) SDiff(ctx context.Context, keys ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SDiff(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SDiff(ctx, r.ks(keys...)...)
	}
	return res
}

// SDiffStore 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SDiffStore(ctx context.Context, destination string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SDiffStore(ctx, r.k(destination), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SDiffStore(ctx, r.k(destination), r.ks(keys...)...)
	}
	return res
}

// SInter 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
// 当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
func (r *Client) SInter(ctx context.Context, keys ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SInter(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SInter(ctx, r.ks(keys...)...)
	}
	return res
}

// SInterStore 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SInterStore(ctx context.Context, destination string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SInterStore(ctx, r.k(destination), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SInterStore(ctx, r.k(destination), r.ks(keys...)...)
	}
	return res
}

// SIsMember 判断 member 元素是否集合 key 的成员。
func (r *Client) SIsMember(ctx context.Context, key string, member interface{}) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SIsMember(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.SIsMember(ctx, r.k(key), member)
	}
	return res
}

// SMembers 返回集合 key 中的所有成员。
// 不存在的 key 被视为空集合。
func (r *Client) SMembers(ctx context.Context, key string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SMembers(ctx, r.k(key))
	case ClientNormal:
		res = r.client.SMembers(ctx, r.k(key))
	}
	return res
}

// SMembersMap -> SMembers
func (r *Client) SMembersMap(ctx context.Context, key string) (res *redis.StringStructMapCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SMembersMap(ctx, r.k(key))
	case ClientNormal:
		res = r.client.SMembersMap(ctx, r.k(key))
	}
	return res
}

// SMove 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
func (r *Client) SMove(ctx context.Context, source, destination string, member interface{}) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SMove(ctx, r.k(source), r.k(destination), member)
	case ClientNormal:
		res = r.client.SMove(ctx, r.k(source), r.k(destination), member)
	}
	return res
}

// SPop 移除并返回集合中的一个随机元素。
// 如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
func (r *Client) SPop(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SPop(ctx, r.k(key))
	case ClientNormal:
		res = r.client.SPop(ctx, r.k(key))
	}
	return res
}

// SPopN -> SPop
func (r *Client) SPopN(ctx context.Context, key string, count int64) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SPopN(ctx, r.k(key), count)
	case ClientNormal:
		res = r.client.SPopN(ctx, r.k(key), count)
	}
	return res
}

// SRandMember 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 从 Redis 2.6 版本开始， SRANDMEMBER 命令接受可选的 count 参数：
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 SRANDMEMBER 则仅仅返回随机元素，而不对集合进行任何改动。
func (r *Client) SRandMember(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SRandMember(ctx, r.k(key))
	case ClientNormal:
		res = r.client.SRandMember(ctx, r.k(key))
	}
	return res
}

// SRandMemberN -> SRandMember
func (r *Client) SRandMemberN(ctx context.Context, key string, count int64) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SRandMemberN(ctx, r.k(key), count)
	case ClientNormal:
		res = r.client.SRandMemberN(ctx, r.k(key), count)
	}
	return res
}

// SRem 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
func (r *Client) SRem(ctx context.Context, key string, members ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SRem(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.SRem(ctx, r.k(key), members...)
	}
	return res
}

// SUnion 返回一个集合的全部成员，该集合是所有给定集合的并集。
// 不存在的 key 被视为空集。
func (r *Client) SUnion(ctx context.Context, keys ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SUnion(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SUnion(ctx, r.ks(keys...)...)
	}
	return res
}

// SUnionStore 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (r *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SUnionStore(ctx, r.k(destination), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.SUnionStore(ctx, r.k(destination), r.ks(keys...)...)
	}
	return res
}

// XAdd 将指定的流条目追加到指定key的流中。 如果key不存在，作为运行这个命令的副作用，将使用流的条目自动创建key。
func (r *Client) XAdd(ctx context.Context, a *redis.XAddArgs) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XAdd(ctx, a)
	case ClientNormal:
		res = r.client.XAdd(ctx, a)
	}
	return res
}

// XDel 从指定流中移除指定的条目，并返回成功删除的条目的数量，在传递的ID不存在的情况下， 返回的数量可能与传递的ID数量不同。
func (r *Client) XDel(ctx context.Context, stream string, ids ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XDel(ctx, stream, ids...)
	case ClientNormal:
		res = r.client.XDel(ctx, stream, ids...)
	}
	return res
}

// XLen 返回流中的条目数。如果指定的key不存在，则此命令返回0，就好像该流为空。 但是请注意，与其他的Redis类型不同，零长度流是可能的，所以你应该调用TYPE 或者 EXISTS 来检查一个key是否存在。
// 一旦内部没有任何的条目（例如调用XDEL后），流不会被自动删除，因为可能还存在与其相关联的消费者组。
func (r *Client) XLen(ctx context.Context, stream string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XLen(ctx, stream)
	case ClientNormal:
		res = r.client.XLen(ctx, stream)
	}
	return res
}

// XRange 此命令返回流中满足给定ID范围的条目。范围由最小和最大ID指定。所有ID在指定的两个ID之间或与其中一个ID相等（闭合区间）的条目将会被返回。
func (r *Client) XRange(ctx context.Context, stream, start, stop string) (res *redis.XMessageSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XRange(ctx, stream, start, stop)
	case ClientNormal:
		res = r.client.XRange(ctx, stream, start, stop)
	}
	return res
}

// XRangeN -> XRange
func (r *Client) XRangeN(ctx context.Context, stream, start, stop string, count int64) (res *redis.XMessageSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XRangeN(ctx, stream, start, stop, count)
	case ClientNormal:
		res = r.client.XRangeN(ctx, stream, start, stop, count)
	}
	return res
}

// XRevRange 此命令与XRANGE完全相同，但显著的区别是以相反的顺序返回条目，并以相反的顺序获取开始-结束参数：在XREVRANGE中，你需要先指定结束ID，再指定开始ID，该命令就会从结束ID侧开始生成两个ID之间（或完全相同）的所有元素。
func (r *Client) XRevRange(ctx context.Context, stream string, start, stop string) (res *redis.XMessageSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XRevRange(ctx, stream, start, stop)
	case ClientNormal:
		res = r.client.XRevRange(ctx, stream, start, stop)
	}
	return res
}

// XRevRangeN -> XRevRange
func (r *Client) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) (res *redis.XMessageSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XRevRangeN(ctx, stream, start, stop, count)
	case ClientNormal:
		res = r.client.XRevRangeN(ctx, stream, start, stop, count)
	}
	return res
}

// XRead 从一个或者多个流中读取数据，仅返回ID大于调用者报告的最后接收ID的条目。此命令有一个阻塞选项，用于等待可用的项目，类似于BRPOP或者BZPOPMIN等等。
func (r *Client) XRead(ctx context.Context, a *redis.XReadArgs) (res *redis.XStreamSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XRead(ctx, a)
	case ClientNormal:
		res = r.client.XRead(ctx, a)
	}
	return res
}

// XReadStreams -> XRead
func (r *Client) XReadStreams(ctx context.Context, streams ...string) (res *redis.XStreamSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XReadStreams(ctx, streams...)
	case ClientNormal:
		res = r.client.XReadStreams(ctx, streams...)
	}
	return res
}

// XGroupCreate command
func (r *Client) XGroupCreate(ctx context.Context, stream, group, start string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XGroupCreate(ctx, stream, group, start)
	case ClientNormal:
		res = r.client.XGroupCreate(ctx, stream, group, start)
	}
	return res
}

// XGroupCreateMkStream command
func (r *Client) XGroupCreateMkStream(ctx context.Context, stream, group, start string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XGroupCreateMkStream(ctx, stream, group, start)
	case ClientNormal:
		res = r.client.XGroupCreateMkStream(ctx, stream, group, start)
	}
	return res
}

// XGroupSetID command
func (r *Client) XGroupSetID(ctx context.Context, stream, group, start string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XGroupSetID(ctx, stream, group, start)
	case ClientNormal:
		res = r.client.XGroupSetID(ctx, stream, group, start)
	}
	return res
}

// XGroupDestroy command
func (r *Client) XGroupDestroy(ctx context.Context, stream, group string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XGroupDestroy(ctx, stream, group)
	case ClientNormal:
		res = r.client.XGroupDestroy(ctx, stream, group)
	}
	return res
}

// XGroupDelConsumer command
func (r *Client) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XGroupDelConsumer(ctx, stream, group, consumer)
	case ClientNormal:
		res = r.client.XGroupDelConsumer(ctx, stream, group, consumer)
	}
	return res
}

// XReadGroup command
func (r *Client) XReadGroup(ctx context.Context, a *redis.XReadGroupArgs) (res *redis.XStreamSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XReadGroup(ctx, a)
	case ClientNormal:
		res = r.client.XReadGroup(ctx, a)
	}
	return res
}

// XAck command
func (r *Client) XAck(ctx context.Context, stream, group string, ids ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XAck(ctx, stream, group, ids...)
	case ClientNormal:
		res = r.client.XAck(ctx, stream, group, ids...)
	}
	return res
}

// XPending command
func (r *Client) XPending(ctx context.Context, stream, group string) (res *redis.XPendingCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XPending(ctx, stream, group)
	case ClientNormal:
		res = r.client.XPending(ctx, stream, group)
	}
	return res
}

// XPendingExt command
func (r *Client) XPendingExt(ctx context.Context, a *redis.XPendingExtArgs) (res *redis.XPendingExtCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XPendingExt(ctx, a)
	case ClientNormal:
		res = r.client.XPendingExt(ctx, a)
	}
	return res
}

// XClaim command
func (r *Client) XClaim(ctx context.Context, a *redis.XClaimArgs) (res *redis.XMessageSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XClaim(ctx, a)
	case ClientNormal:
		res = r.client.XClaim(ctx, a)
	}
	return res
}

// XClaimJustID command
func (r *Client) XClaimJustID(ctx context.Context, a *redis.XClaimArgs) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XClaimJustID(ctx, a)
	case ClientNormal:
		res = r.client.XClaimJustID(ctx, a)
	}
	return res
}

// XTrim command
func (r *Client) XTrim(ctx context.Context, key string, maxLen int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XTrim(ctx, key, maxLen)
	case ClientNormal:
		res = r.client.XTrim(ctx, key, maxLen)
	}
	return res
}

// XTrimApprox command
func (r *Client) XTrimApprox(ctx context.Context, key string, maxLen int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XTrimApprox(ctx, key, maxLen)
	case ClientNormal:
		res = r.client.XTrimApprox(ctx, key, maxLen)
	}
	return res
}

// XInfoGroups command
func (r *Client) XInfoGroups(ctx context.Context, key string) (res *redis.XInfoGroupsCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.XInfoGroups(ctx, key)
	case ClientNormal:
		res = r.client.XInfoGroups(ctx, key)
	}
	return res
}

// BZPopMax 是有序集合命令 ZPOPMAX带有阻塞功能的版本。
func (r *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (res *redis.ZWithKeyCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BZPopMax(ctx, timeout, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BZPopMax(ctx, timeout, r.ks(keys...)...)
	}
	return res
}

// BZPopMin 是有序集合命令 ZPOPMIN带有阻塞功能的版本。
func (r *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (res *redis.ZWithKeyCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BZPopMin(ctx, timeout, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.BZPopMin(ctx, timeout, r.ks(keys...)...)
	}
	return res
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
// 如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。
// score 值可以是整数值或双精度浮点数。
// 如果 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAdd(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAdd(ctx, r.k(key), members...)
	}
	return res
}

// ZAddNX -> ZAdd
func (r *Client) ZAddNX(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAddNX(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAddNX(ctx, r.k(key), members...)
	}
	return res
}

// ZAddXX -> ZAdd
func (r *Client) ZAddXX(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAddXX(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAddXX(ctx, r.k(key), members...)
	}
	return res
}

// ZAddCh -> ZAdd
func (r *Client) ZAddCh(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAddCh(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAddCh(ctx, r.k(key), members...)
	}
	return res
}

// ZAddNXCh -> ZAdd
func (r *Client) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAddNXCh(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAddNXCh(ctx, r.k(key), members...)
	}
	return res
}

// ZAddXXCh -> ZAdd
func (r *Client) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZAddXXCh(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZAddXXCh(ctx, r.k(key), members...)
	}
	return res
}

// ZIncr Redis `ZADD key INCR score member` command.
func (r *Client) ZIncr(ctx context.Context, key string, member *redis.Z) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZIncr(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZIncr(ctx, r.k(key), member)
	}
	return res
}

// ZIncrNX Redis `ZADD key NX INCR score member` command.
func (r *Client) ZIncrNX(ctx context.Context, key string, member *redis.Z) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZIncrNX(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZIncrNX(ctx, r.k(key), member)
	}
	return res
}

// ZIncrXX Redis `ZADD key XX INCR score member` command.
func (r *Client) ZIncrXX(ctx context.Context, key string, member *redis.Z) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZIncrXX(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZIncrXX(ctx, r.k(key), member)
	}
	return res
}

// ZCard 返回有序集 key 的基数。
func (r *Client) ZCard(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZCard(ctx, r.k(key))
	case ClientNormal:
		res = r.client.ZCard(ctx, r.k(key))
	}
	return res
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
// 关于参数 min 和 max 的详细使用方法，请参考 ZRANGEBYSCORE 命令。
func (r *Client) ZCount(ctx context.Context, key, min, max string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZCount(ctx, r.k(key), min, max)
	case ClientNormal:
		res = r.client.ZCount(ctx, r.k(key), min, max)
	}
	return res
}

// ZLexCount -> ZCount
func (r *Client) ZLexCount(ctx context.Context, key, min, max string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZLexCount(ctx, r.k(key), min, max)
	case ClientNormal:
		res = r.client.ZLexCount(ctx, r.k(key), min, max)
	}
	return res
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
// 当 key 不是有序集类型时，返回一个错误。
// score 值可以是整数值或双精度浮点数。
func (r *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZIncrBy(ctx, r.k(key), increment, member)
	case ClientNormal:
		res = r.client.ZIncrBy(ctx, r.k(key), increment, member)
	}
	return res
}

// ZInterStore 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
// 关于 WEIGHTS 和 AGGREGATE 选项的描述，参见 ZUNIONSTORE 命令。
func (r *Client) ZInterStore(ctx context.Context, key string, store *redis.ZStore) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZInterStore(ctx, r.k(key), store)
	case ClientNormal:
		res = r.client.ZInterStore(ctx, r.k(key), store)
	}
	return res
}

// ZPopMax 删除并返回有序集合key中的最多count个具有最高得分的成员。
func (r *Client) ZPopMax(ctx context.Context, key string, count ...int64) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZPopMax(ctx, r.k(key), count...)
	case ClientNormal:
		res = r.client.ZPopMax(ctx, r.k(key), count...)
	}
	return res
}

// ZPopMin 删除并返回有序集合key中的最多count个具有最低得分的成员。
func (r *Client) ZPopMin(ctx context.Context, key string, count ...int64) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZPopMin(ctx, r.k(key), count...)
	case ClientNormal:
		res = r.client.ZPopMin(ctx, r.k(key), count...)
	}
	return res
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (r *Client) ZRange(ctx context.Context, key string, start, stop int64) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRange(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.ZRange(ctx, r.k(key), start, stop)
	}
	return res
}

// ZRangeWithScores -> ZRange
func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRangeWithScores(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.ZRangeWithScores(ctx, r.k(key), start, stop)
	}
	return res
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
func (r *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRangeByScore(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRangeByScore(ctx, r.k(key), opt)
	}
	return res
}

// ZRangeByLex -> ZRangeByScore
func (r *Client) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRangeByLex(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRangeByLex(ctx, r.k(key), opt)
	}
	return res
}

// ZRangeByScoreWithScores -> ZRangeByScore
func (r *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRangeByScoreWithScores(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRangeByScoreWithScores(ctx, r.k(key), opt)
	}
	return res
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
// 使用 ZREVRANK 命令可以获得成员按 score 值递减(从大到小)排列的排名。
func (r *Client) ZRank(ctx context.Context, key, member string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRank(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZRank(ctx, r.k(key), member)
	}
	return res
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// 当 key 存在但不是有序集类型时，返回一个错误。
func (r *Client) ZRem(ctx context.Context, key string, members ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRem(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.ZRem(ctx, r.k(key), members...)
	}
	return res
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推
func (r *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRemRangeByRank(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.ZRemRangeByRank(ctx, r.k(key), start, stop)
	}
	return res
}

// ZRemRangeByScore 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 自版本2.1.6开始， score 值等于 min 或 max 的成员也可以不包括在内，详情请参见 ZRANGEBYSCORE 命令。
func (r *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRemRangeByScore(ctx, r.k(key), min, max)
	case ClientNormal:
		res = r.client.ZRemRangeByScore(ctx, r.k(key), min, max)
	}
	return res
}

// ZRemRangeByLex -> ZRemRangeByScore
func (r *Client) ZRemRangeByLex(ctx context.Context, key, min, max string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRemRangeByLex(ctx, r.k(key), min, max)
	case ClientNormal:
		res = r.client.ZRemRangeByLex(ctx, r.k(key), min, max)
	}
	return res
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGE 命令的其他方面和 ZRANGE 命令一样。
func (r *Client) ZRevRange(ctx context.Context, key string, start, stop int64) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRange(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.ZRevRange(ctx, r.k(key), start, stop)
	}
	return res
}

// ZRevRangeWithScores -> ZRevRange
func (r *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRangeWithScores(ctx, r.k(key), start, stop)
	case ClientNormal:
		res = r.client.ZRevRangeWithScores(ctx, r.k(key), start, stop)
	}
	return res
}

// ZRevRangeByScore 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。
// 除了成员按 score 值递减的次序排列这一点外， ZREVRANGEBYSCORE 命令的其他方面和 ZRANGEBYSCORE 命令一样。
func (r *Client) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRangeByScore(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRevRangeByScore(ctx, r.k(key), opt)
	}
	return res
}

// ZRevRangeByLex -> ZRevRangeByScore
func (r *Client) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRangeByLex(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRevRangeByLex(ctx, r.k(key), opt)
	}
	return res
}

// ZRevRangeByScoreWithScores -> ZRevRangeByScore
func (r *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) (res *redis.ZSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRangeByScoreWithScores(ctx, r.k(key), opt)
	case ClientNormal:
		res = r.client.ZRevRangeByScoreWithScores(ctx, r.k(key), opt)
	}
	return res
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (r *Client) ZRevRank(ctx context.Context, key, member string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZRevRank(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZRevRank(ctx, r.k(key), member)
	}
	return res
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (r *Client) ZScore(ctx context.Context, key, member string) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZScore(ctx, r.k(key), member)
	case ClientNormal:
		res = r.client.ZScore(ctx, r.k(key), member)
	}
	return res
}

// ZUnionStore 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之 和 。
func (r *Client) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ZUnionStore(ctx, r.k(dest), store)
	case ClientNormal:
		res = r.client.ZUnionStore(ctx, r.k(dest), store)
	}
	return res
}

// PFAdd 将指定元素添加到HyperLogLog
func (r *Client) PFAdd(ctx context.Context, key string, els ...interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PFAdd(ctx, r.k(key), els...)
	case ClientNormal:
		res = r.client.PFAdd(ctx, r.k(key), els...)
	}
	return res
}

// PFCount 返回HyperlogLog观察到的集合的近似基数。
func (r *Client) PFCount(ctx context.Context, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PFCount(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.PFCount(ctx, r.ks(keys...)...)
	}
	return res
}

// PFMerge N个不同的HyperLogLog合并为一个。
func (r *Client) PFMerge(ctx context.Context, dest string, keys ...string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PFMerge(ctx, r.k(dest), r.ks(keys...)...)
	case ClientNormal:
		res = r.client.PFMerge(ctx, r.k(dest), r.ks(keys...)...)
	}
	return res
}

// BgRewriteAOF 异步重写附加文件
func (r *Client) BgRewriteAOF(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BgRewriteAOF(ctx)
	case ClientNormal:
		res = r.client.BgRewriteAOF(ctx)
	}
	return res
}

// BgSave 将数据集异步保存到磁盘
func (r *Client) BgSave(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.BgSave(ctx)
	case ClientNormal:
		res = r.client.BgSave(ctx)
	}
	return res
}

// ClientKill 杀掉客户端的连接
func (r *Client) ClientKill(ctx context.Context, ipPort string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientKill(ctx, ipPort)
	case ClientNormal:
		res = r.client.ClientKill(ctx, ipPort)
	}
	return res
}

// ClientKillByFilter is new style synx, while the ClientKill is old
// CLIENT KILL <option> [value] ... <option> [value]
func (r *Client) ClientKillByFilter(ctx context.Context, keys ...string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientKillByFilter(ctx, r.ks(keys...)...)
	case ClientNormal:
		res = r.client.ClientKillByFilter(ctx, r.ks(keys...)...)
	}
	return res
}

// ClientList 获取客户端连接列表
func (r *Client) ClientList(ctx context.Context) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientList(ctx)
	case ClientNormal:
		res = r.client.ClientList(ctx)
	}
	return res
}

// ClientPause 停止处理来自客户端的命令一段时间
func (r *Client) ClientPause(ctx context.Context, dur time.Duration) (res *redis.BoolCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientPause(ctx, dur)
	case ClientNormal:
		res = r.client.ClientPause(ctx, dur)
	}
	return res
}

// ClientID Returns the client ID for the current connection
func (r *Client) ClientID(ctx context.Context) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClientID(ctx)
	case ClientNormal:
		res = r.client.ClientID(ctx)
	}
	return res
}

// ConfigGet 获取指定配置参数的值
func (r *Client) ConfigGet(ctx context.Context, parameter string) (res *redis.SliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ConfigGet(ctx, parameter)
	case ClientNormal:
		res = r.client.ConfigGet(ctx, parameter)
	}
	return res
}

// ConfigResetStat 重置 INFO 命令中的某些统计数据
func (r *Client) ConfigResetStat(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ConfigResetStat(ctx)
	case ClientNormal:
		res = r.client.ConfigResetStat(ctx)
	}
	return res
}

// ConfigSet 修改 redis 配置参数，无需重启
func (r *Client) ConfigSet(ctx context.Context, parameter, value string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ConfigSet(ctx, parameter, value)
	case ClientNormal:
		res = r.client.ConfigSet(ctx, parameter, value)
	}
	return res
}

// ConfigRewrite 对启动 Redis 服务器时所指定的 redis.conf 配置文件进行改写
func (r *Client) ConfigRewrite(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ConfigRewrite(ctx)
	case ClientNormal:
		res = r.client.ConfigRewrite(ctx)
	}
	return res
}

// DBSize 返回当前数据库的 key 的数量
func (r *Client) DBSize(ctx context.Context) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.DBSize(ctx)
	case ClientNormal:
		res = r.client.DBSize(ctx)
	}
	return res
}

// FlushAll 删除所有数据库的所有key
func (r *Client) FlushAll(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.FlushAll(ctx)
	case ClientNormal:
		res = r.client.FlushAll(ctx)
	}
	return res
}

// FlushAllAsync 异步删除所有数据库的所有key
func (r *Client) FlushAllAsync(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.FlushAllAsync(ctx)
	case ClientNormal:
		res = r.client.FlushAllAsync(ctx)
	}
	return res
}

// FlushDB 删除当前数据库的所有key
func (r *Client) FlushDB(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.FlushDB(ctx)
	case ClientNormal:
		res = r.client.FlushDB(ctx)
	}
	return res
}

// FlushDBAsync 异步删除当前数据库的所有key
func (r *Client) FlushDBAsync(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.FlushDBAsync(ctx)
	case ClientNormal:
		res = r.client.FlushDBAsync(ctx)
	}
	return res
}

// Info 获取 Redis 服务器的各种信息和统计数值
func (r *Client) Info(ctx context.Context, section ...string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Info(ctx, section...)
	case ClientNormal:
		res = r.client.Info(ctx, section...)
	}
	return res
}

// LastSave 返回最近一次 Redis 成功将数据保存到磁盘上的时间，以 UNIX 时间戳格式表示
func (r *Client) LastSave(ctx context.Context) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.LastSave(ctx)
	case ClientNormal:
		res = r.client.LastSave(ctx)
	}
	return res
}

// Save 异步保存数据到硬盘
func (r *Client) Save(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Save(ctx)
	case ClientNormal:
		res = r.client.Save(ctx)
	}
	return res
}

// Shutdown 关闭服务器
func (r *Client) Shutdown(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Shutdown(ctx)
	case ClientNormal:
		res = r.client.Shutdown(ctx)
	}
	return res
}

// ShutdownSave 异步保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownSave(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ShutdownSave(ctx)
	case ClientNormal:
		res = r.client.ShutdownSave(ctx)
	}
	return res
}

// ShutdownNoSave 不保存数据到硬盘，并关闭服务器
func (r *Client) ShutdownNoSave(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ShutdownNoSave(ctx)
	case ClientNormal:
		res = r.client.ShutdownNoSave(ctx)
	}
	return res
}

// SlaveOf 将当前服务器转变为指定服务器的从属服务器(slave server)
func (r *Client) SlaveOf(ctx context.Context, host, port string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.SlaveOf(ctx, host, port)
	case ClientNormal:
		res = r.client.SlaveOf(ctx, host, port)
	}
	return res
}

// Time 返回当前服务器时间
func (r *Client) Time(ctx context.Context) (res *redis.TimeCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Time(ctx)
	case ClientNormal:
		res = r.client.Time(ctx)
	}
	return res
}

// Eval 执行 Lua 脚本。
func (r *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (res *redis.Cmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Eval(ctx, script, r.ks(keys...), args...)
	case ClientNormal:
		res = r.client.Eval(ctx, script, r.ks(keys...), args...)
	}
	return res
}

// EvalSha 执行 Lua 脚本。
func (r *Client) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (res *redis.Cmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.EvalSha(ctx, sha1, r.ks(keys...), args...)
	case ClientNormal:
		res = r.client.EvalSha(ctx, sha1, r.ks(keys...), args...)
	}
	return res
}

// ScriptExists 查看指定的脚本是否已经被保存在缓存当中。
func (r *Client) ScriptExists(ctx context.Context, hashes ...string) (res *redis.BoolSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ScriptExists(ctx, hashes...)
	case ClientNormal:
		res = r.client.ScriptExists(ctx, hashes...)
	}
	return res
}

// ScriptFlush 从脚本缓存中移除所有脚本。
func (r *Client) ScriptFlush(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ScriptFlush(ctx)
	case ClientNormal:
		res = r.client.ScriptFlush(ctx)
	}
	return res
}

// ScriptKill 杀死当前正在运行的 Lua 脚本。
func (r *Client) ScriptKill(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ScriptKill(ctx)
	case ClientNormal:
		res = r.client.ScriptKill(ctx)
	}
	return res

}

// ScriptLoad 将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本。
func (r *Client) ScriptLoad(ctx context.Context, script string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ScriptLoad(ctx, script)
	case ClientNormal:
		res = r.client.ScriptLoad(ctx, script)
	}
	return res
}

// DebugObject 获取 key 的调试信息
func (r *Client) DebugObject(ctx context.Context, key string) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.DebugObject(ctx, r.k(key))
	case ClientNormal:
		res = r.client.DebugObject(ctx, r.k(key))
	}
	return res
}

// Publish 将信息发送到指定的频道。
func (r *Client) Publish(ctx context.Context, channel string, message interface{}) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Publish(ctx, r.k(channel), message)
	case ClientNormal:
		res = r.client.Publish(ctx, r.k(channel), message)
	}
	return res
}

// PubSubChannels 订阅一个或多个符合给定模式的频道。
func (r *Client) PubSubChannels(ctx context.Context, pattern string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PubSubChannels(ctx, r.k(pattern))
	case ClientNormal:
		res = r.client.PubSubChannels(ctx, r.k(pattern))
	}
	return res
}

// PubSubNumSub 查看订阅与发布系统状态。
func (r *Client) PubSubNumSub(ctx context.Context, channels ...string) (res *redis.StringIntMapCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PubSubNumSub(ctx, r.ks(channels...)...)
	case ClientNormal:
		res = r.client.PubSubNumSub(ctx, r.ks(channels...)...)
	}
	return res
}

// PubSubNumPat 用于获取redis订阅或者发布信息的状态
func (r *Client) PubSubNumPat(ctx context.Context) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PubSubNumPat(ctx)
	case ClientNormal:
		res = r.client.PubSubNumPat(ctx)
	}
	return res
}

// ClusterSlots 获取集群节点的映射数组
func (r *Client) ClusterSlots(ctx context.Context) (res *redis.ClusterSlotsCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterSlots(ctx)
	case ClientNormal:
		res = r.client.ClusterSlots(ctx)
	}
	return res
}

// ClusterNodes Get Cluster config for the node
func (r *Client) ClusterNodes(ctx context.Context) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterNodes(ctx)
	case ClientNormal:
		res = r.client.ClusterNodes(ctx)
	}
	return res
}

// ClusterMeet Force a node cluster to handshake with another node
func (r *Client) ClusterMeet(ctx context.Context, host, port string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterMeet(ctx, host, port)
	case ClientNormal:
		res = r.client.ClusterMeet(ctx, host, port)
	}
	return res
}

// ClusterForget Remove a node from the nodes table
func (r *Client) ClusterForget(ctx context.Context, nodeID string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterForget(ctx, nodeID)
	case ClientNormal:
		res = r.client.ClusterForget(ctx, nodeID)
	}
	return res
}

// ClusterReplicate Reconfigure a node as a replica of the specified master node
func (r *Client) ClusterReplicate(ctx context.Context, nodeID string) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterReplicate(ctx, nodeID)
	case ClientNormal:
		res = r.client.ClusterReplicate(ctx, nodeID)
	}
	return res

}

// ClusterResetSoft Reset a Redis Cluster node
func (r *Client) ClusterResetSoft(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterResetSoft(ctx)
	case ClientNormal:
		res = r.client.ClusterResetSoft(ctx)
	}
	return res
}

// ClusterResetHard Reset a Redis Cluster node
func (r *Client) ClusterResetHard(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterResetHard(ctx)
	case ClientNormal:
		res = r.client.ClusterResetHard(ctx)
	}
	return res
}

// ClusterInfo Provides info about Redis Cluster node state
func (r *Client) ClusterInfo(ctx context.Context) (res *redis.StringCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterInfo(ctx)
	case ClientNormal:
		res = r.client.ClusterInfo(ctx)
	}
	return res
}

// ClusterKeySlot Returns the hash slot of the specified key
func (r *Client) ClusterKeySlot(ctx context.Context, key string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterKeySlot(ctx, r.k(key))
	case ClientNormal:
		res = r.client.ClusterKeySlot(ctx, r.k(key))
	}
	return res
}

// ClusterGetKeysInSlot Return local key names in the specified hash slot
func (r *Client) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterGetKeysInSlot(ctx, slot, count)
	case ClientNormal:
		res = r.client.ClusterGetKeysInSlot(ctx, slot, count)
	}
	return res
}

// ClusterCountFailureReports Return the number of failure reports active for a given node
func (r *Client) ClusterCountFailureReports(ctx context.Context, nodeID string) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterCountFailureReports(ctx, nodeID)
	case ClientNormal:
		res = r.client.ClusterCountFailureReports(ctx, nodeID)
	}
	return res
}

// ClusterCountKeysInSlot Return the number of local keys in the specified hash slot
func (r *Client) ClusterCountKeysInSlot(ctx context.Context, slot int) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterCountKeysInSlot(ctx, slot)
	case ClientNormal:
		res = r.client.ClusterCountKeysInSlot(ctx, slot)
	}
	return res
}

// ClusterDelSlots Set hash slots as unbound in receiving node
func (r *Client) ClusterDelSlots(ctx context.Context, slots ...int) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterDelSlots(ctx, slots...)
	case ClientNormal:
		res = r.client.ClusterDelSlots(ctx, slots...)
	}
	return res
}

// ClusterDelSlotsRange ->  ClusterDelSlots
func (r *Client) ClusterDelSlotsRange(ctx context.Context, min, max int) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterDelSlotsRange(ctx, min, max)
	case ClientNormal:
		res = r.client.ClusterDelSlotsRange(ctx, min, max)
	}
	return res
}

// ClusterSaveConfig Forces the node to save cluster state on disk
func (r *Client) ClusterSaveConfig(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterSaveConfig(ctx)
	case ClientNormal:
		res = r.client.ClusterSaveConfig(ctx)
	}
	return res
}

// ClusterSlaves List replica nodes of the specified master node
func (r *Client) ClusterSlaves(ctx context.Context, nodeID string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterSlaves(ctx, nodeID)
	case ClientNormal:
		res = r.client.ClusterSlaves(ctx, nodeID)
	}
	return res
}

// ClusterFailover Forces a replica to perform a manual failover of its master.
func (r *Client) ClusterFailover(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterFailover(ctx)
	case ClientNormal:
		res = r.client.ClusterFailover(ctx)
	}
	return res
}

// ClusterAddSlots Assign new hash slots to receiving node
func (r *Client) ClusterAddSlots(ctx context.Context, slots ...int) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterAddSlots(ctx, slots...)
	case ClientNormal:
		res = r.client.ClusterAddSlots(ctx, slots...)
	}
	return res
}

// ClusterAddSlotsRange -> ClusterAddSlots
func (r *Client) ClusterAddSlotsRange(ctx context.Context, min, max int) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ClusterAddSlotsRange(ctx, min, max)
	case ClientNormal:
		res = r.client.ClusterAddSlotsRange(ctx, min, max)
	}
	return res
}

// GeoAdd 将指定的地理空间位置（纬度、经度、名称）添加到指定的key中
func (r *Client) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoAdd(ctx, r.k(key), geoLocation...)
	case ClientNormal:
		res = r.client.GeoAdd(ctx, r.k(key), geoLocation...)
	}
	return res
}

// GeoPos 从key里返回所有给定位置元素的位置（经度和纬度）
func (r *Client) GeoPos(ctx context.Context, key string, members ...string) (res *redis.GeoPosCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoPos(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.GeoPos(ctx, r.k(key), members...)
	}
	return res
}

// GeoRadius 以给定的经纬度为中心， 找出某一半径内的元素
func (r *Client) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (res *redis.GeoLocationCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoRadius(ctx, r.k(key), longitude, latitude, query)
	case ClientNormal:
		res = r.client.GeoRadius(ctx, r.k(key), longitude, latitude, query)
	}
	return res
}

// GeoRadiusStore -> GeoRadius
func (r *Client) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoRadiusStore(ctx, r.k(key), longitude, latitude, query)
	case ClientNormal:
		res = r.client.GeoRadiusStore(ctx, r.k(key), longitude, latitude, query)
	}
	return res
}

// GeoRadiusByMember -> GeoRadius
func (r *Client) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (res *redis.GeoLocationCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoRadiusByMember(ctx, r.k(key), member, query)
	case ClientNormal:
		res = r.client.GeoRadiusByMember(ctx, r.k(key), member, query)
	}
	return res
}

// GeoRadiusByMemberStore 找出位于指定范围内的元素，中心点是由给定的位置元素决定
func (r *Client) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoRadiusByMemberStore(ctx, r.k(key), member, query)
	case ClientNormal:
		res = r.client.GeoRadiusByMemberStore(ctx, r.k(key), member, query)
	}
	return res
}

// GeoDist 返回两个给定位置之间的距离
func (r *Client) GeoDist(ctx context.Context, key string, member1, member2, unit string) (res *redis.FloatCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoDist(ctx, r.k(key), member1, member2, unit)
	case ClientNormal:
		res = r.client.GeoDist(ctx, r.k(key), member1, member2, unit)
	}
	return res
}

// GeoHash 返回一个或多个位置元素的 Geohash 表示
func (r *Client) GeoHash(ctx context.Context, key string, members ...string) (res *redis.StringSliceCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.GeoHash(ctx, r.k(key), members...)
	case ClientNormal:
		res = r.client.GeoHash(ctx, r.k(key), members...)
	}
	return res
}

// ReadOnly Enables read queries for a connection to a cluster replica node
func (r *Client) ReadOnly(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ReadOnly(ctx)
	case ClientNormal:
		res = r.client.ReadOnly(ctx)
	}
	return res
}

// ReadWrite Disables read queries for a connection to a cluster replica node
func (r *Client) ReadWrite(ctx context.Context) (res *redis.StatusCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.ReadWrite(ctx)
	case ClientNormal:
		res = r.client.ReadWrite(ctx)
	}
	return res
}

// MemoryUsage Estimate the memory usage of a key
func (r *Client) MemoryUsage(ctx context.Context, key string, samples ...int) (res *redis.IntCmd) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.MemoryUsage(ctx, r.k(key), samples...)
	case ClientNormal:
		res = r.client.MemoryUsage(ctx, r.k(key), samples...)
	}
	return res
}

// Subscribe subscribes the client to the specified channels.
// Channels can be omitted to create empty subscription.
func (r *Client) Subscribe(ctx context.Context, channels ...string) (res *redis.PubSub) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.Subscribe(ctx, channels...)
	case ClientNormal:
		res = r.client.Subscribe(ctx, channels...)
	}
	return res
}

// PSubscribe subscribes the client to the given patterns.
// Patterns can be omitted to create empty subscription.
func (r *Client) PSubscribe(ctx context.Context, channels ...string) (res *redis.PubSub) {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	switch r.clientType {
	case ClientCluster:
		res = r.clusterClient.PSubscribe(ctx, channels...)
	case ClientNormal:
		res = r.client.PSubscribe(ctx, channels...)
	}
	return res
}

// ErrNotImplemented not implemented error
var ErrNotImplemented = errors.New("Not implemented")
