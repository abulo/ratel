package queue

// Jober 是一个可以执行的异步任务
type Jober interface {
	Job()
}

// SyncJober 可以执行的同步任务
type SyncJober interface {
	Jober
	Wait() <-chan any
	Error() error
}

type job struct {
	v        any
	callback func(any)
}

// NewJob 创建一个异步任务
func NewJob(v any, fn func(any)) Jober {
	return &job{
		v:        v,
		callback: fn,
	}
}

// Job ...
func (j *job) Job() {
	j.callback(j.v)
}

type syncJob struct {
	err      error
	result   chan any
	v        any
	callback func(any) (any, error)
}

// NewSyncJob 创建同步任务
func NewSyncJob(v any, fn func(any) (any, error)) SyncJober {
	return &syncJob{
		result:   make(chan any, 1),
		v:        v,
		callback: fn,
	}
}

// Job ...
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

// Wait ...
func (j *syncJob) Wait() <-chan any {
	return j.result
}

// Error ...
func (j *syncJob) Error() error {
	return j.err
}
