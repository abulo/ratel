package cycle

import (
	"sync"
	"sync/atomic"
)

// Cycle 管理一组后台任务的并发生命周期；须由 NewCycle 创建，零值不可直接使用。
//
// 任务的追加、完成等待与关闭操作均为并发安全。
type Cycle struct {
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	done    chan struct{}
	quit    chan error
	closing uint32
	waiting uint32
}

// NewCycle 创建一个空的生命周期控制器。
func NewCycle() *Cycle {
	return &Cycle{
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		quit:    make(chan error),
		closing: 0,
		waiting: 0,
	}
}

// Run 将 fn 作为后台任务纳入当前生命周期；fn 返回的错误会发送到 Wait 返回的退出通道。
//
// TODO(owner): 等待前补充任务状态与 panic 检查。
func (c *Cycle) Run(fn func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wg.Add(1)
	go func(c *Cycle) {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}(c)
}

// Done 返回在全部已登记任务结束后关闭的通道；多次调用返回同一通道，等待协程只启动一次。
func (c *Cycle) Done() <-chan struct{} {
	if atomic.CompareAndSwapUint32(&c.waiting, 0, 1) {
		go func(c *Cycle) {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.wg.Wait()
			close(c.done)
		}(c)
	}
	return c.done
}

// DoneAndClose 等待全部任务结束后关闭当前 Cycle。
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

// Close 关闭退出通道；重复调用安全，关闭操作只执行一次。
func (c *Cycle) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
		close(c.quit)
	}
}

// Wait 返回用于接收任务错误的通道；任一后台任务返回非空错误时，错误会送达该通道。
func (c *Cycle) Wait() <-chan error {
	return c.quit
}
