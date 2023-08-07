package task

import "github.com/abulo/ratel/v3/core/task/cron"

const (
	JobLocal       = "Local"
	JobDistributed = "Distributed"
)

// JobWrapper is a job wrapper
type JobWrapper struct {
	Id   cron.EntryID
	Job  cron.Job
	Func func()

	Task *Task
	Name string
	Type string
}

func (job *JobWrapper) Run() {
	if job.Type == JobLocal {
		job.run()
	}

	if job.Type == JobDistributed &&
		job.Task.thisNodeRun(job.Name) {
		job.run()
	}
}

func (job *JobWrapper) run() {
	if job.Func != nil {
		job.Func()
	}

	if job.Job != nil {
		job.Job.Run()
	}
}
