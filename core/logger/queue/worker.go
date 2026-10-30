package queue

import (
	"sync"
)

func newWorker(pool chan chan Jober, wg *sync.WaitGroup) *worker {
	return &worker{
		pool:    pool,
		wg:      wg,
		jobChan: make(chan Jober),
		quit:    make(chan struct{}),
	}
}

type worker struct {
	pool    chan chan Jober
	wg      *sync.WaitGroup
	jobChan chan Jober
	quit    chan struct{}
}

func (w *worker) Start() {
	w.pool <- w.jobChan
	go w.dispatcher()
}

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

func (w *worker) Stop() {
	close(w.quit)
}
