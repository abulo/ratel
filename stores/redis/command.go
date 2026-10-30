package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Pipeline 获取底层连接上的未执行管道对象，供调用方批量组装命令后统一提交。
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

// Pipelined 在管道中执行 fn 组装的全部命令并一次性提交，返回批量命令结果；
// fn 内通过传入的 Pipeliner 追加命令。
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

// TxPipelined 在事务管道（MULTI/EXEC）中执行 fn 组装的全部命令并一次性提交，
// 返回事务内的批量命令结果。
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

// TxPipeline 获取事务管道对象，供调用方批量组装命令后以事务方式统一提交。
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

// Command 返回服务器支持的全部命令及其详细信息。
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

// CommandList 按过滤器返回命令名称列表；filter 为 nil 时返回全部命令名称。
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

// CommandGetKeys 解析给定命令行并返回其涉及的键。
func (r *Client) CommandGetKeys(ctx context.Context, commands ...any) (val []string, err error) {
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

// CommandGetKeysAndFlags 解析给定命令行，返回其涉及的键及对应的读写访问标志。
func (r *Client) CommandGetKeysAndFlags(ctx context.Context, commands ...any) (val []redis.KeyFlags, err error) {
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

// ClientGetName 返回当前连接通过 CLIENT SETNAME 设置的客户端名称。
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

// Echo 原样返回传入的消息，用于连通性校验。
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

// Ping 探测与服务器的连通性，正常时返回 PONG。
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

// Quit 主动关闭当前连接。
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

// Unlink 异步删除给定键，删除由后台执行，命令返回时不保证内存已立即释放。
func (r *Client) Unlink(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.Unlink(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// BgRewriteAOF 异步重写 AOF 持久化文件。
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

// BgSave 异步将内存数据快照保存到磁盘，命令不等待保存完成。
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

// ClientKill 终止指定 ip:port 的客户端连接。
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

// ClientKillByFilter 按过滤条件批量终止匹配的客户端连接。
func (r *Client) ClientKillByFilter(ctx context.Context, keys ...string) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.ClientKillByFilter(getCtx(ctx), keys...).Result()
		return err
	}, acceptable)
	return
}

// ClientList 返回已连接客户端的原始信息字符串。
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

// ClientInfo 返回当前客户端的连接信息。
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

// ClientPause 在指定时长内暂停服务器处理客户端命令。
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

// ClientUnpause 恢复被 ClientPause 暂停的命令处理。
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

// ClientID 返回当前连接的客户端 ID。
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

// ClientUnblock 解除指定 ID 客户端在阻塞命令上的挂起状态。
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

// ClientUnblockWithError 解除指定 ID 客户端的阻塞，并以错误方式唤醒被阻塞的命令。
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

// ConfigGet 读取服务器指定配置项；parameter 支持通配匹配模式。
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

// ConfigResetStat 重置服务器各类运行统计计数器。
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

// ConfigSet 将服务器指定配置项修改为给定值。
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

// ConfigRewrite 将当前运行时配置写回服务器配置文件。
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

// DBSize 返回当前数据库中的键数量。
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

// FlushAll 删除所有数据库中的全部键。
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

// FlushAllAsync 异步删除所有数据库中的全部键。
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

// FlushDB 删除当前数据库中的全部键。
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

// FlushDBAsync 异步删除当前数据库中的全部键。
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

// Info 返回服务器统计信息；section 不传时返回全部模块，传入时限定指定模块。
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

// LastSave 返回最后一次成功快照保存的 UNIX 时间戳（秒）。
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

// Save 同步阻塞执行一次快照保存，完成后才返回。
func (r *Client) Save(ctx context.Context) (val string, err error) {
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

// Shutdown 关闭服务器，是否保存快照由服务端默认策略决定。
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

// ShutdownSave 先保存快照再关闭服务器。
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

// ShutdownNoSave 不保存快照直接关闭服务器，未持久化的数据将丢失。
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

// SlaveOf 将当前节点配置为指定主节点的从节点，host 与 port 为主节点地址。
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

// SlowLogGet 返回最近 num 条慢查询日志。
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

// Time 返回服务器当前时间。
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

// DebugObject 返回指定键的内部编码与引用等调试信息。
func (r *Client) DebugObject(ctx context.Context, key string) (val string, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.DebugObject(getCtx(ctx), key).Result()
		return err
	}, acceptable)
	return
}

// MemoryUsage 返回指定键占用的内存字节数；samples 控制估算采样精度，可选传入。
func (r *Client) MemoryUsage(ctx context.Context, key string, samples ...int) (val int64, err error) {
	err = r.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(r)
		if err != nil {
			return err
		}
		val, err = conn.MemoryUsage(getCtx(ctx), key, samples...).Result()
		return err
	}, acceptable)
	return
}

// ModuleLoadex 按给定配置向服务器加载模块。
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
// 			p.Get(ctx, cast.ToString(k))
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
