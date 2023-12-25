package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func (r *Client) CommandList(ctx context.Context, filter *redis.FilterBy) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CommandList(getCtx(ctx), filter).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) CommandGetKeys(ctx context.Context, commands ...interface{}) (val []string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CommandGetKeys(getCtx(ctx), commands...).Result()
		return err
	}, acceptable)
	return
}
func (r *Client) CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) (val []redis.KeyFlags, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.CommandGetKeysAndFlags(getCtx(ctx), commands...).Result()
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

// Unlink 这个命令非常类似于DEL：它删除指定的键。就像DEL键一样，如果它不存在，它将被忽略。但是，该命令在不同的线程中执行实际的内存回收，所以它不会阻塞，而DEL是。这是命令名称的来源：命令只是将键与键空间断开连接。实际删除将在以后异步发生。
func (r *Client) Unlink(ctx context.Context, keys ...any) (val int64, err error) {
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
func (r *Client) ClientKillByFilter(ctx context.Context, keys ...any) (val int64, err error) {
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

func (r *Client) ClientInfo(ctx context.Context) (val *redis.ClientInfo, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientInfo(getCtx(ctx)).Result()
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

func (r *Client) DebugObject(ctx context.Context, key any) (val string, err error) {
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

func (r *Client) MemoryUsage(ctx context.Context, key any, samples ...int) (val int64, err error) {
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

func (r *Client) ModuleLoadex(ctx context.Context, conf *redis.ModuleLoadexConfig) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ModuleLoadex(getCtx(ctx), conf).Result()
		return err
	}, acceptable)
	return
}

// func (r *Client) Process(ctx context.Context, cmd redis.Cmder) error {
// 	return r.brk.DoWithAcceptable(func() error {
// 		conn, err := getRedis(r)
// 		if err != nil {
// 			return err
// 		}
// 		return conn.Process(getCtx(ctx), cmd)
// 	}, acceptable)
// }

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
