package queue

import (
	"sync"
	"sync/atomic"
)

// NewQueue 创建一个队列，指定缓冲区数和工作线程数
func NewQueue(maxCapacity, maxThread int) *Queue {
	return &Queue{
		jobQueue:   make(chan Jober, maxCapacity),
		maxWorkers: maxThread,
		workerPool: make(chan chan Jober, maxThread),
		workers:    make([]*worker, maxThread),
		wg:         new(sync.WaitGroup),
	}
}

// Queue 任务队列排队以在高并发情况下减轻服务器压力并改进任务处理
type Queue struct {
	maxWorkers int
	jobQueue   chan Jober
	workerPool chan chan Jober
	workers    []*worker
	running    uint32
	wg         *sync.WaitGroup
}

// Run 开始运行队列
func (q *Queue) Run() {
	if atomic.LoadUint32(&q.running) == 1 {
		return
	}

	atomic.StoreUint32(&q.running, 1)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i] = newWorker(q.workerPool, q.wg)
		q.workers[i].Start()
	}

	go q.dispatcher()
}

func (q *Queue) dispatcher() {
	for job := range q.jobQueue {
		worker := <-q.workerPool
		worker <- job
	}
}

// Terminate 终止队列以接收任务并释放资源
func (q *Queue) Terminate() {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	atomic.StoreUint32(&q.running, 0)
	q.wg.Wait()

	close(q.jobQueue)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}

// Push 将可执行任务放入队列
func (q *Queue) Push(job Jober) {
	if atomic.LoadUint32(&q.running) != 1 {
		return
	}

	q.wg.Add(1)
	q.jobQueue <- job
}
