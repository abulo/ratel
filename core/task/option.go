package task

import (
	"time"

	"github.com/abulo/ratel/v3/core/task/cron"
)

// Option is Task Option
type Option func(*Task)

// WithNodeUpdateDuration set node update duration
func WithNodeUpdateDuration(d time.Duration) Option {
	return func(task *Task) {
		task.nodeUpdateDuration = d
	}
}

// WithHashReplicas set hashReplicas
func WithHashReplicas(d int) Option {
	return func(task *Task) {
		task.hashReplicas = d
	}
}

// CronOptionLocation is warp cron with location
func CronOptionLocation(loc *time.Location) Option {
	return func(task *Task) {
		f := cron.WithLocation(loc)
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionSeconds is warp cron with seconds
func CronOptionSeconds() Option {
	return func(task *Task) {
		f := cron.WithSeconds()
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionParser is warp cron with schedules.
func CronOptionParser(p cron.ScheduleParser) Option {
	return func(task *Task) {
		f := cron.WithParser(p)
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionChain is Warp cron with chain
func CronOptionChain(wrappers ...cron.JobWrapper) Option {
	return func(task *Task) {
		f := cron.WithChain(wrappers...)
		task.crOptions = append(task.crOptions, f)
	}
}

// You can defined yourself recover function to make the
// job will be added to your task when the process restart
func WithRecoverFunc(recoverFunc RecoverFuncType) Option {
	return func(task *Task) {
		task.RecoverFunc = recoverFunc
	}
}

// You can use this option to start the recent jobs rerun
// after the cluster upgrading.
func WithClusterStable(timeWindow time.Duration) Option {
	return func(d *Task) {
		d.recentJobs = NewRecentJobPacker(timeWindow)
	}
}

func RunningLocally() Option {
	return func(d *Task) {
		d.runningLocally = true
	}
}
