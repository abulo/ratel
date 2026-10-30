package queue

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// NewListQueue 创建无长度上限的链表队列，maxThread 为消费 worker 数量。
func NewListQueue(maxThread int) *ListQueue {
	return NewListQueueWithMaxLen(maxThread, 0)
}

// NewListQueueWithMaxLen 创建链表队列并限定待缓冲任务上限 maxLen；maxLen 为 0 表示不限制长度。
func NewListQueueWithMaxLen(maxThread, maxLen int) *ListQueue {
	return &ListQueue{
		maxLen:     maxLen,
		maxWorker:  maxThread,
		workers:    make([]*worker, maxThread),
		workerPool: make(chan chan Jober, maxThread),
		list:       list.New(),
		lock:       new(sync.RWMutex),
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
	wg         *sync.WaitGroup
	running    atomic.Uint32
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
	for {
		q.lock.RLock()
		if q.running.Load() != 1 && q.list.Len() == 0 {
			q.lock.RUnlock()
			break
		}
		ele := q.list.Front()
		q.lock.RUnlock()

		if ele == nil {
			time.Sleep(time.Millisecond * 10)
			continue
		}

		worker := <-q.workerPool
		worker <- ele.Value.(Jober)

		q.lock.Lock()
		q.list.Remove(ele)
		q.lock.Unlock()
	}
}

// Push 将任务追加到链表尾部；队列未运行或已超出 maxLen 上限时任务被静默丢弃。
func (q *ListQueue) Push(job Jober) {
	if q.running.Load() != 1 {
		return
	}

	if q.maxLen > 0 {
		q.lock.RLock()
		if q.list.Len() > q.maxLen {
			q.lock.RUnlock()
			return
		}
		q.lock.RUnlock()
	}

	q.wg.Add(1)
	q.lock.Lock()
	q.list.PushBack(job)
	q.lock.Unlock()
}

// Terminate 停止接收新任务，等待已入队任务执行完毕后停止 worker；队列未运行时直接返回。
func (q *ListQueue) Terminate() {
	if q.running.Load() != 1 {
		return
	}

	q.running.Store(0)
	q.wg.Wait()

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}
