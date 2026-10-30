package task

import "github.com/abulo/ratel/v3/core/task/cron"

// Job 是定时任务的最小执行单元，Run 在被调度到时执行任务逻辑。
type Job interface {
	Run()
}

// StableJob 是可在服务重启后恢复重建的任务，额外提供 cron 表达式读取与序列化、反序列化能力。
type StableJob interface {
	Job
	GetCron() string
	Serialize() ([]byte, error)
	UnSerialize([]byte) error
}

// JobWarpper 是注册到调度器的内部任务包装器，关联所属 Task、任务名、cron 表达式与原始任务。
type JobWarpper struct {
	// ID 为调度器分配的任务条目标识。
	ID cron.EntryID
	// Task 指向所属 Task。
	Task *Task
	// Name 为任务唯一名称。
	Name string
	// CronStr 为任务使用的 cron 表达式。
	CronStr string
	// Job 为被包装的原始任务。
	Job Job
}

// Run 先判定当前节点是否有权执行该任务，允许时再调用 Execute 执行。
func (job JobWarpper) Run() {
	if job.Task.allowThisNodeRun(job.Name) {
		job.Execute()
	}
}

// Execute 直接执行内部任务；任务为 nil 时不做任何操作。
func (job JobWarpper) Execute() {
	if job.Job != nil {
		job.Job.Run()
	}
}
