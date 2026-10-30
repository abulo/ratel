package task

import "time"

// IRecentJobPacker 记录集群节点状态不稳定期间近期执行过的任务，待集群恢复稳定后供 Task 回放。
// 实现必须保证并发安全。
type IRecentJobPacker interface {
	// AddJob 将近期执行的任务连同其执行时间加入记录，超出时间窗口的旧任务会被自动剔除。
	AddJob(jobName string, t time.Time) error

	// PopAllJobs 取出并清空当前记录的全部任务名，返回顺序不做保证。
	PopAllJobs() (jobNames []string)
}
