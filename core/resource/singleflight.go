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

	flightGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

// NewSingleFlight 创建一个并发安全的 SingleFlight 实例，内部按 key 跟踪进行中的调用。
func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
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

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
