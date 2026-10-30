package cron

// JobWrapper 为给定 Job 叠加一层横切行为（如异常恢复、并发互斥、日志记录）。
type JobWrapper func(Job) Job

// Chain 是一组 JobWrapper 的有序序列，用于统一为提交的任务叠加横切行为。
type Chain struct {
	wrappers []JobWrapper
}

// NewChain 由给定的包装器构造一条包装器链。
func NewChain(c ...JobWrapper) Chain {
	return Chain{c}
}

// Then 用链上全部包装器包装给定任务。
//
// NewChain(m1, m2, m3).Then(job) 等价于 m1(m2(m3(job)))，
// 即链首包装器最外层、链尾包装器最内层。
func (c Chain) Then(j Job) Job {
	for i := range c.wrappers {
		j = c.wrappers[len(c.wrappers)-i-1](j)
	}
	return j
}
