package worker

// Worker could scheduled by jupiter or customized scheduler
type Worker interface {
	WorkerStart() error
	WorkerStop() error
}
