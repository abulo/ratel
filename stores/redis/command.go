package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Pipeline 获取Redis管道 无参数
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

// Pipelined 执行管道操作 ctx: 上下文, fn: 管道操作函数
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

// TxPipelined 执行事务管道操作 ctx: 上下文, fn: 管道操作函数
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

// TxPipeline 获取事务管道 无参数
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

// Command 获取所有Redis命令的详细信息 ctx: 上下文
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

// CommandList 获取Redis命令列表 ctx: 上下文, filter: 过滤器
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

// CommandGetKeys 获取命令对应的键 ctx: 上下文, commands: 命令列表
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

// CommandGetKeysAndFlags 获取命令对应的键及标志 ctx: 上下文, commands: 命令列表
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

// ClientGetName 获取客户端连接名称 ctx: 上下文
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

// Echo 返回输入的消息 ctx: 上下文, message: 要返回的消息
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

// Ping 测试与Redis服务器的连接 ctx: 上下文
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

// Quit 关闭客户端连接 ctx: 上下文
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

// Unlink 异步删除键 ctx: 上下文, keys: 要删除的键列表
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

// BgRewriteAOF 异步重写AOF文件 ctx: 上下文
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

// BgSave 异步保存数据到磁盘 ctx: 上下文
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

// ClientKill 终止客户端连接 ctx: 上下文, ipPort: 客户端地址
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

// ClientKillByFilter 根据条件终止客户端连接 ctx: 上下文, keys: 过滤条件
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

// ClientList 获取客户端连接列表 ctx: 上下文
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

// ClientInfo 获取客户端详细信息 ctx: 上下文
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

// ClientPause 暂停处理客户端命令 ctx: 上下文, dur: 暂停时间
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

// ClientUnpause 恢复处理客户端命令 ctx: 上下文
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

// ClientID 获取客户端ID ctx: 上下文
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

// ClientUnblock 解除客户端阻塞状态 ctx: 上下文, id: 客户端ID
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

// ClientUnblockWithError 解除客户端阻塞状态并返回错误 ctx: 上下文, id: 客户端ID
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

// ConfigGet 获取配置参数 ctx: 上下文, parameter: 参数名
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

// ConfigResetStat 重置统计信息 ctx: 上下文
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

// ConfigSet 设置配置参数 ctx: 上下文, parameter: 参数名, value: 参数值
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

// ConfigRewrite 重写配置文件 ctx: 上下文
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

// DBSize 获取当前数据库键数量 ctx: 上下文
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

// FlushAll 删除所有数据库的所有键 ctx: 上下文
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

// FlushAllAsync 异步删除所有数据库的所有键 ctx: 上下文
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

// FlushDB 删除当前数据库的所有键 ctx: 上下文
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

// FlushDBAsync 异步删除当前数据库的所有键 ctx: 上下文
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

// Info 获取服务器信息 ctx: 上下文, section: 信息模块
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

// LastSave 获取最后一次保存时间 ctx: 上下文
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

// Save 同步保存数据到磁盘 ctx: 上下文
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

// Shutdown 关闭服务器 ctx: 上下文
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

// ShutdownSave 保存数据并关闭服务器 ctx: 上下文
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

// ShutdownNoSave 不保存数据直接关闭服务器 ctx: 上下文
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

// SlaveOf 设置主从复制 ctx: 上下文, host: 主机地址, port: 端口号
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

// SlowLogGet 获取慢查询日志 ctx: 上下文, num: 日志条数
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

// Time 获取服务器时间 ctx: 上下文
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

// DebugObject 获取键的调试信息 ctx: 上下文, key: 键名
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

// MemoryUsage 获取键的内存使用情况 ctx: 上下文, key: 键名, samples: 采样次数
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

// ModuleLoadex 加载Redis模块 ctx: 上下文, conf: 模块配置
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
