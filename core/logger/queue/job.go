package queue

// Jober 表示可被队列异步执行的任务，Job 为任务执行入口，由 worker 协程调用。
type Jober interface {
	Job()
}

// SyncJober 表示可等待执行结果的同步任务，在 Jober 基础上补充结果通道与错误获取能力。
type SyncJober interface {
	Jober
	Wait() <-chan any
	Error() error
}

type job struct {
	v        any
	callback func(any)
}

// NewJob 创建绑定数据与回调的异步任务；callback 在 worker 协程中以 v 为入参执行。
func NewJob(v any, fn func(any)) Jober {
	return &job{
		v:        v,
		callback: fn,
	}
}

func (j *job) Job() {
	j.callback(j.v)
}

type syncJob struct {
	err      error
	result   chan any
	v        any
	callback func(any) (any, error)
}

// NewSyncJob 创建可等待结果的同步任务；callback 返回执行结果与错误，任务执行结束后写入 Wait 返回的通道。
func NewSyncJob(v any, fn func(any) (any, error)) SyncJober {
	return &syncJob{
		result:   make(chan any, 1),
		v:        v,
		callback: fn,
	}
}

func (j *syncJob) Job() {
	result, err := j.callback(j.v)
	if err != nil {
		j.err = err
		close(j.result)
		return
	}

	j.result <- result

	close(j.result)
}

func (j *syncJob) Wait() <-chan any {
	return j.result
}

func (j *syncJob) Error() error {
	return j.err
}
