package resource

import "sync"

type (
	// SingleFlight 将同一 key 下的并发调用合并为一次：首个调用实际执行 fn，
	// 其余并发到达的相同 key 调用等待并共享其返回结果；不同 key 的调用互不影响。
	SingleFlight interface {
		// Do 执行或共享 key 对应调用的结果，返回值与首个执行者完全一致。
		Do(key string, fn func() (any, error)) (any, error)

		// DoEx 行为同 Do，额外通过 fresh 标识当前调用是否为实际执行 fn 的首个调用：
		// fresh 为 true 表示本次执行了 fn，false 表示共享了他人的执行结果。
		DoEx(key string, fn func() (any, error)) (any, bool, error)
	}

	call struct {
		wg  sync.WaitGroup
		val any
		err error
	}

	// PERF-4-6 [P2]：按 key 哈希分片，不同 key 的并发查找落在不同互斥锁上，
	// 替代原先覆盖全部 key 的单一全局互斥锁（分片后锁竞争缓解程度预期收益需实测验证）。
	flightShard struct {
		calls map[string]*call
		lock  sync.Mutex
	}

	flightGroup struct {
		shards [flightShardCount]flightShard
	}
)

// flightShardCount 为单飞调用表的分片数量，取 2 的幂便于取模。
const flightShardCount = 16

// NewSingleFlight 创建一个并发安全的 SingleFlight 实例，内部按 key 跟踪进行中的调用。
func NewSingleFlight() SingleFlight {
	g := &flightGroup{}
	for i := range g.shards {
		g.shards[i].calls = make(map[string]*call)
	}
	return g
}

func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, true, c.err
}

// shard 返回 key 所属分片；内联 FNV-1a 哈希分布均匀，避免每次调用分配哈希对象。
func (g *flightGroup) shard(key string) *flightShard {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return &g.shards[h%flightShardCount]
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	sh := g.shard(key)
	sh.lock.Lock()
	if c, ok := sh.calls[key]; ok {
		sh.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	sh.calls[key] = c
	sh.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	sh := g.shard(key)
	defer func() {
		sh.lock.Lock()
		delete(sh.calls, key)
		sh.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
