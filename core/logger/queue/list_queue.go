package queue

import (
	"container/list"
	"sync"
	"sync/atomic"
)

// NewListQueue 创建无长度上限的链表队列，maxThread 为消费 worker 数量。
func NewListQueue(maxThread int) *ListQueue {
	return NewListQueueWithMaxLen(maxThread, 0)
}

// NewListQueueWithMaxLen 创建链表队列并限定待缓冲任务上限 maxLen；maxLen 为 0 表示不限制长度。
func NewListQueueWithMaxLen(maxThread, maxLen int) *ListQueue {
	// PERF-3-1 [P1]：cond 与同一把锁绑定，入队后唤醒等待中的 dispatcher，替代空队列固定 sleep 轮询。
	lock := new(sync.RWMutex)
	return &ListQueue{
		maxLen:     maxLen,
		maxWorker:  maxThread,
		workers:    make([]*worker, maxThread),
		workerPool: make(chan chan Jober, maxThread),
		list:       list.New(),
		lock:       lock,
		cond:       sync.NewCond(lock),
		wg:         new(sync.WaitGroup),
	}
}

// ListQueue 基于双向链表与 worker 协程池实现任务队列，适合待缓冲任务数量需要动态伸缩的场景。
//
// 零值不可直接使用，须经 NewListQueue 或 NewListQueueWithMaxLen 构造并调用 Run 启动。
type ListQueue struct {
	maxLen     int
	maxWorker  int
	workers    []*worker
	workerPool chan chan Jober
	list       *list.List
	lock       *sync.RWMutex
	// PERF-3-1 [P1]：任务到达时 Signal 唤醒阻塞在 Wait 上的 dispatcher，消除 10ms 固定轮询延迟。
	cond    *sync.Cond
	wg      *sync.WaitGroup
	running atomic.Uint32
}

// Run 启动队列；重复调用时直接返回，不会重复创建 worker。
func (q *ListQueue) Run() {
	if q.running.Load() == 1 {
		return
	}
	q.running.Store(1)

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i] = newWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *ListQueue) dispatcher() {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		// PERF-3-1 [P1]：空队列时阻塞在 cond.Wait 上等待唤醒；只有在已停止且队列排空后才退出，
		// 消费顺序仍为链表 FIFO，与原轮询实现一致。
		for q.list.Len() == 0 {
			if q.running.Load() != 1 {
				return
			}
			q.cond.Wait()
		}
		ele := q.list.Front()
		q.lock.Unlock()

		worker := <-q.workerPool
		worker <- ele.Value.(Jober)

		q.lock.Lock()
		q.list.Remove(ele)
	}
}

// Push 将任务追加到链表尾部；队列未运行或已超出 maxLen 上限时任务被静默丢弃。
func (q *ListQueue) Push(job Jober) {
	if q.running.Load() != 1 {
		return
	}

	// PERF-3-1 [P1]：长度检查与入队合并为同一个 Lock 临界区，去掉原先先 RLock 查长度再 Lock 入队的两段式加锁；
	// 长度上限语义保持不变（maxLen 为 0 时不限制；默认无界上限的收紧属行为变更，留给业务确认，本次不实施）。
	q.lock.Lock()
	if q.maxLen > 0 && q.list.Len() > q.maxLen {
		q.lock.Unlock()
		return
	}

	q.wg.Add(1)
	q.list.PushBack(job)
	q.cond.Signal()
	q.lock.Unlock()
}

// Terminate 停止接收新任务，等待已入队任务执行完毕后停止 worker；队列未运行时直接返回。
func (q *ListQueue) Terminate() {
	if q.running.Load() != 1 {
		return
	}

	q.running.Store(0)
	// PERF-3-1 [P1]：停止信号通过 cond.Broadcast 唤醒可能正阻塞在 Wait 上的 dispatcher，使其退出事件循环。
	q.cond.Broadcast()
	q.wg.Wait()

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}
