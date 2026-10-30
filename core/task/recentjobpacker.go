package task

import (
	"container/heap"
	"sync"
	"time"
)

// JobWithTime 记录任务名及其最近一次执行时间。
type JobWithTime struct {
	// JobName 为任务唯一名称。
	JobName string
	// RunningTime 为任务最近一次执行时间。
	RunningTime time.Time
}

// JobWithTimeHeap 按执行时间排序的任务堆，供 RecentJobPacker 维护时间窗口内的任务。
type JobWithTimeHeap []JobWithTime

// Pop 弹出堆顶元素，是 heap.Interface 的一部分。
func (jobHeap *JobWithTimeHeap) Pop() (ret any) {
	n := jobHeap.Len() - 1
	ret = (*jobHeap)[n]
	(*jobHeap) = (*jobHeap)[0:n]
	return
}

// Push 将元素加入堆，是 heap.Interface 的一部分。
func (jobHeap *JobWithTimeHeap) Push(x any) {
	(*jobHeap) = append((*jobHeap), x.(JobWithTime))
}

// Less 报告索引 i 处任务的执行时间是否早于索引 j 处。
func (jobHeap *JobWithTimeHeap) Less(i, j int) bool {
	return (*jobHeap)[i].RunningTime.Before((*jobHeap)[j].RunningTime)
}

// Len 返回堆中元素个数。
func (jobHeap *JobWithTimeHeap) Len() int {
	return len(*jobHeap)
}

// Swap 交换两个索引处的元素。
func (jobHeap *JobWithTimeHeap) Swap(i, j int) {
	t := (*jobHeap)[i]
	(*jobHeap)[i] = (*jobHeap)[j]
	(*jobHeap)[j] = t
}

// Index 返回索引 i 处的元素。
func (jobHeap *JobWithTimeHeap) Index(i int) any {
	return (*jobHeap)[i]
}

// RecentJobPacker 在给定时间窗口内缓存近期执行过的任务，供集群恢复稳定后回放；
// 实现自带互斥锁，并发安全。
type RecentJobPacker struct {
	sync.Mutex
	timeWindow time.Duration
	recentJobs JobWithTimeHeap
}

// AddJob 记录任务 jobName 在时间 t 的执行，随后剔除早于当前时间减去时间窗口的旧任务。
func (rjp *RecentJobPacker) AddJob(jobName string, t time.Time) (err error) {
	rjp.Lock()
	defer rjp.Unlock()
	heap.Push(&rjp.recentJobs, JobWithTime{
		JobName:     jobName,
		RunningTime: t,
	})

	unRencentTime := time.Now().Add(-rjp.timeWindow)
	for rjp.recentJobs.Len() > 0 && rjp.recentJobs.Index(0).(JobWithTime).RunningTime.Before(unRencentTime) {
		heap.Pop(&rjp.recentJobs)
	}
	return
}

// PopAllJobs 取出并清空全部已记录的任务名；无任务时返回空切片而非 nil。
func (rjp *RecentJobPacker) PopAllJobs() (jobNames []string) {
	rjp.Lock()
	defer rjp.Unlock()
	jobNames = make([]string, 0)
	for rjp.recentJobs.Len() > 0 {
		jobNames = append(jobNames, heap.Pop(&rjp.recentJobs).(JobWithTime).JobName)
	}
	rjp.recentJobs = make(JobWithTimeHeap, 0)
	heap.Init(&rjp.recentJobs)
	return
}

// NewRecentJobPacker 创建以 timeWindow 为记录窗口的 RecentJobPacker。
func NewRecentJobPacker(timeWindow time.Duration) IRecentJobPacker {
	return &RecentJobPacker{
		timeWindow: timeWindow,
		recentJobs: make([]JobWithTime, 0),
	}
}
