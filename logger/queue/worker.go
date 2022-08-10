package queue

import (
	"sync"
)

// 创建一个worker队列
func newWorker(pool chan chan Jober, wg *sync.WaitGroup) *worker {
	return &worker{
		pool:    pool,
		wg:      wg,
		jobChan: make(chan Jober),
		quit:    make(chan struct{}),
	}
}

// worker队列结构体
type worker struct {
	pool    chan chan Jober
	wg      *sync.WaitGroup
	jobChan chan Jober
	quit    chan struct{}
}

// 开始这个worker
func (w *worker) Start() {
	w.pool <- w.jobChan
	go w.dispatcher()
}

// 分发执行
func (w *worker) dispatcher() {
	for {
		select {
		case j := <-w.jobChan:
			j.Job()
			w.pool <- w.jobChan
			w.wg.Done()
		case <-w.quit:
			<-w.pool
			close(w.jobChan)
			return
		}
	}
}

// 停止 这个worker
func (w *worker) Stop() {
	close(w.quit)
}
