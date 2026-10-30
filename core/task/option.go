package task

import (
	"time"

	"github.com/abulo/ratel/v3/core/task/cron"
)

// Option 是 Task 的配置选项函数，采用函数式选项模式在构造 Task 时调整其行为。
type Option func(*Task)

// WithNodeUpdateDuration 设置节点列表的刷新间隔，默认 3 秒。
func WithNodeUpdateDuration(d time.Duration) Option {
	return func(task *Task) {
		task.nodeUpdateDuration = d
	}
}

// WithHashReplicas 设置一致性哈希环中每个节点对应的虚拟副本数，默认 50。
func WithHashReplicas(d int) Option {
	return func(task *Task) {
		task.hashReplicas = d
	}
}

// CronOptionLocation 指定 cron 调度器解析时间表时使用的时区。
func CronOptionLocation(loc *time.Location) Option {
	return func(task *Task) {
		f := cron.WithLocation(loc)
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionSeconds 让 cron 解析器在表达式首位包含秒字段。
func CronOptionSeconds() Option {
	return func(task *Task) {
		f := cron.WithSeconds()
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionParser 使用自定义调度表解析器替换默认解析器。
func CronOptionParser(p cron.ScheduleParser) Option {
	return func(task *Task) {
		f := cron.WithParser(p)
		task.crOptions = append(task.crOptions, f)
	}
}

// CronOptionChain 为该 Task 添加一组任务包装器，统一应用到后续注册的所有任务上。
func CronOptionChain(wrappers ...cron.JobWrapper) Option {
	return func(task *Task) {
		f := cron.WithChain(wrappers...)
		task.crOptions = append(task.crOptions, f)
	}
}

// WithRecoverFunc 注册任务启动前执行的恢复回调，用于在进程重启时重建任务状态。
func WithRecoverFunc(recoverFunc RecoverFuncType) Option {
	return func(task *Task) {
		task.RecoverFunc = recoverFunc
	}
}

// WithClusterStable 开启近期任务回放：在集群节点变更导致状态不稳定期间记录近期执行过的任务，
// 待集群恢复稳定后重新触发；timeWindow 为记录近期任务的时间窗口。
func WithClusterStable(timeWindow time.Duration) Option {
	return func(d *Task) {
		d.recentJobs = NewRecentJobPacker(timeWindow)
	}
}

// RunningLocally 让 Task 以本地单机模式运行，不创建节点池、不参与集群一致性哈希分配。
func RunningLocally() Option {
	return func(d *Task) {
		d.runningLocally = true
	}
}
