package queue

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// NewListQueue 创建一个列表队列，指定工作线程的数量
func NewListQueue(maxThread int) *ListQueue {
	return NewListQueueWithMaxLen(maxThread, 0)
}

// NewListQueueWithMaxLen 创建一个列表队列，指定工作线程数和最大元素数
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

// ListQueue 是一个列表任务队列，用于在高并发情况下缓解服务器压力并改进任务处理
type ListQueue struct {
	maxLen     int
	maxWorker  int
	workers    []*worker
	workerPool chan chan Jober
	list       *list.List
	lock       *sync.RWMutex
	wg         *sync.WaitGroup
	running    uint32
}

// Run 开始运行队列
func (q *ListQueue) Run() {
	if atomic.LoadUint32(&q.running) == 1 {
		return
	}
	atomic.StoreUint32(&q.running, 1)

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i] = newWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *ListQueue) dispatcher() {
	for {
		q.lock.RLock()
		if atomic.LoadUint32(&q.running) != 1 && q.list.Len() == 0 {
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

// Push 将可执行任务推入队列
func (q *ListQueue) Push(job Jober) {
	if atomic.LoadUint32(&q.running) != 1 {
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

// Terminate 终止队列以接收任务并释放资源
func (q *ListQueue) Terminate() {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	atomic.StoreUint32(&q.running, 0)
	q.wg.Wait()

	for i := 0; i < q.maxWorker; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}
