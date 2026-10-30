package queue

import (
	"sync"
	"sync/atomic"
)

// NewQueue 创建带缓冲通道队列；maxCapacity 为待缓冲任务上限，maxThread 为消费 worker 数量。
func NewQueue(maxCapacity, maxThread int) *Queue {
	return &Queue{
		jobQueue:   make(chan Jober, maxCapacity),
		maxWorkers: maxThread,
		workerPool: make(chan chan Jober, maxThread),
		workers:    make([]*worker, maxThread),
		wg:         new(sync.WaitGroup),
	}
}

// Queue 基于带缓冲通道与 worker 协程池实现任务队列。
//
// 零值不可直接使用，须经 NewQueue 构造并调用 Run 启动；同一实例重复调用 Run 不会重复创建 worker。
type Queue struct {
	maxWorkers int
	jobQueue   chan Jober
	workerPool chan chan Jober
	workers    []*worker
	running    atomic.Uint32
	wg         *sync.WaitGroup
}

// Run 启动队列；重复调用时直接返回，不会重复创建 worker。
func (q *Queue) Run() {
	if q.running.Load() == 1 {
		return
	}

	q.running.Store(1)
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

// Terminate 停止接收新任务，等待已入队任务执行完毕后关闭通道与 worker；队列未运行时直接返回。
func (q *Queue) Terminate() {
	if q.running.Load() != 1 {
		return
	}

	q.running.Store(0)
	q.wg.Wait()

	close(q.jobQueue)
	for i := 0; i < q.maxWorkers; i++ {
		q.workers[i].Stop()
	}
	close(q.workerPool)
}

// Push 将任务放入队列；队列未运行时直接返回、任务被丢弃；缓冲已满时调用方阻塞等待空位。
func (q *Queue) Push(job Jober) {
	if q.running.Load() != 1 {
		return
	}

	q.wg.Add(1)
	q.jobQueue <- job
}
