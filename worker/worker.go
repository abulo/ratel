package worker

// Worker could scheduled by ratel or customized scheduler
type Worker interface {
	WorkerStart() error
	WorkerStop() error
}
