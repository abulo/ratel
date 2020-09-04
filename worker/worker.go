package worker

// Worker could scheduled by jupiter or customized scheduler
type Worker interface {
	Run() error
	Stop() error
}
