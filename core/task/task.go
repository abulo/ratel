package task

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/cron"
	"github.com/abulo/ratel/v3/core/task/driver"
)

const (
	defaultReplicas = 50
	defaultDuration = 3 * time.Second

	taskRunning = 1
	taskStopped = 0

	taskStateSteady  = "taskStateSteady"
	taskStateUpgrade = "taskStateUpgrade"
)

var (
	// ErrJobExist 表示注册任务时任务名已存在。
	ErrJobExist = errors.New("jobName already exist")
	// ErrJobNotExist 表示按任务名查找任务时任务不存在。
	ErrJobNotExist = errors.New("jobName not exist")
	// ErrJobWrongNode 表示该任务被集群分配到其他节点，不在当前节点运行。
	ErrJobWrongNode = errors.New("job is not running in this node")
)

// RecoverFuncType 是 Task 启动前执行的恢复回调签名，参数为当前 Task 实例。
type RecoverFuncType func(d *Task)

// Task 是分布式定时任务管理器，负责注册、调度与移除任务，并通过节点池在集群内分配任务。
//
// 零值不可直接使用，须通过 NewTask 或 NewTaskWithOption 构造；任务注册表由内部读写锁保护，
// Start、Run、Stop 可重复调用，重复启动为无操作。
type Task struct {
	jobs      map[string]*JobWarpper
	jobsRWMut sync.RWMutex
	// ServerName 为当前服务在集群中的名称。
	ServerName         string
	nodePool           INodePool
	running            int32
	nodeUpdateDuration time.Duration
	hashReplicas       int
	cr                 *cron.Cron
	crOptions          []cron.Option
	// RecoverFunc 为启动前执行的恢复回调，可为 nil。
	RecoverFunc    RecoverFuncType
	recentJobs     IRecentJobPacker
	state          atomic.Value
	runningLocally bool
}

// NewTask 使用指定服务名与分布式驱动创建 Task，并附加 cron 调度器选项；启动前可继续注册任务。
//
// driver 用于节点注册与集群成员发现，不可为 nil；cronOpts 配置时区、解析器、包装器等调度行为。
func NewTask(serverName string, driver driver.DriverV2, cronOpts ...cron.Option) *Task {
	task := newTask(serverName)
	task.crOptions = cronOpts
	task.cr = cron.New(cronOpts...)
	task.running = taskStopped
	task.nodePool = NewNodePool(serverName, driver, task.nodeUpdateDuration, task.hashReplicas)
	return task
}

// NewTaskWithOption 在 NewTask 基础上接受 Task 级配置选项，可指定本地运行模式、节点刷新间隔、
// 哈希副本数以及 cron 相关选项。
func NewTaskWithOption(serverName string, driver driver.DriverV2, taskOpts ...Option) *Task {
	task := newTask(serverName)
	for _, opt := range taskOpts {
		opt(task)
	}

	task.cr = cron.New(task.crOptions...)
	if !task.runningLocally {
		task.nodePool = NewNodePool(serverName, driver, task.nodeUpdateDuration, task.hashReplicas)
	}
	return task
}

func newTask(serverName string) *Task {
	return &Task{
		ServerName:         serverName,
		jobs:               make(map[string]*JobWarpper),
		crOptions:          make([]cron.Option, 0),
		nodeUpdateDuration: defaultDuration,
		hashReplicas:       defaultReplicas,
	}
}

// AddJob 按 cronStr 注册任务 jobName；任务名重复时返回 ErrJobExist，表达式解析失败时返回解析错误。
func (d *Task) AddJob(jobName, cronStr string, job Job) (err error) {
	return d.addJob(jobName, cronStr, job)
}

// AddFunc 将无参函数 cmd 按 cronStr 注册为定时任务，错误行为约束同 AddJob。
func (d *Task) AddFunc(jobName, cronStr string, cmd func()) (err error) {
	return d.addJob(jobName, cronStr, cron.FuncJob(cmd))
}

func (d *Task) addJob(jobName, cronStr string, job Job) (err error) {
	logger.Logger.Infof("addJob '%s' : %s", jobName, cronStr)

	d.jobsRWMut.Lock()
	defer d.jobsRWMut.Unlock()
	if _, ok := d.jobs[jobName]; ok {
		return ErrJobExist
	}
	innerJob := &JobWarpper{
		Name:    jobName,
		CronStr: cronStr,
		Job:     job,
		Task:    d,
	}
	entryID, err := d.cr.AddJob(cronStr, innerJob)
	if err != nil {
		return err
	}
	innerJob.ID = entryID
	d.jobs[jobName] = innerJob
	return nil
}

// Remove 按任务名移除任务并从调度器注销；任务不存在时静默返回。
func (d *Task) Remove(jobName string) {
	d.jobsRWMut.Lock()
	defer d.jobsRWMut.Unlock()

	if job, ok := d.jobs[jobName]; ok {
		delete(d.jobs, jobName)
		d.cr.Remove(job.ID)
	}
}

// GetJob 按任务名返回内部任务包装器；任务不存在时返回 ErrJobNotExist。
//
// thisNodeOnly 为 true 时，仅当任务被集群分配到当前节点才返回任务，否则返回 ErrJobWrongNode；
// 节点池查询失败时返回对应驱动错误。
func (d *Task) GetJob(jobName string, thisNodeOnly bool) (*JobWarpper, error) {
	d.jobsRWMut.RLock()
	defer d.jobsRWMut.RUnlock()

	job, ok := d.jobs[jobName]
	if !ok {
		logger.Logger.Warnf("job: %s, not exist", jobName)
		return nil, ErrJobNotExist
	}
	if !thisNodeOnly {
		return job, nil
	}
	isRunningHere, err := d.nodePool.CheckJobAvailable(jobName)
	if err != nil {
		return nil, err
	}
	if !isRunningHere {
		return nil, ErrJobWrongNode
	}
	return job, nil
}

// GetJobs 返回当前已注册任务的包装器切片；无任务时返回空切片而非 nil。
//
// thisNodeOnly 为 true 时仅返回分配到当前节点的任务；节点池查询失败的任务会被跳过。
func (d *Task) GetJobs(thisNodeOnly bool) []*JobWarpper {
	d.jobsRWMut.RLock()
	defer d.jobsRWMut.RUnlock()

	ret := make([]*JobWarpper, 0)
	for _, v := range d.jobs {
		var (
			isRunningHere bool
			ok            bool = true
			err           error
		)
		if thisNodeOnly {
			isRunningHere, err = d.nodePool.CheckJobAvailable(v.Name)
			if err != nil {
				continue
			}
			ok = isRunningHere
		}
		if ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// allowThisNodeRun 判定当前节点是否可执行任务；集群不稳定时记录该任务供后续回放。
func (d *Task) allowThisNodeRun(jobName string) (ok bool) {
	if d.runningLocally {
		return true
	}
	ok, err := d.nodePool.CheckJobAvailable(jobName)
	if err != nil {
		logger.Logger.Errorf("allow this node run error, err=%v", err)
		ok = false
		d.state.Store(taskStateUpgrade)
	} else {
		d.state.Store(taskStateSteady)
		if d.recentJobs != nil {
			go d.reRunRecentJobs(d.recentJobs.PopAllJobs())
		}
	}
	if d.recentJobs != nil {
		if d.state.Load().(string) == taskStateUpgrade {
			d.recentJobs.AddJob(jobName, time.Now())
		}
	}
	return
}

// Start 在独立 goroutine 中启动调度器与节点池；重复启动为无操作。
//
// 启动前先执行已注册的 RecoverFunc；节点池启动失败时回退为停止状态并直接返回。
func (d *Task) Start() {
	if d.RecoverFunc != nil {
		d.RecoverFunc(d)
	}
	if atomic.CompareAndSwapInt32(&d.running, taskStopped, taskRunning) {
		if !d.runningLocally {
			if err := d.startNodePool(); err != nil {
				atomic.StoreInt32(&d.running, taskStopped)
				return
			}
			logger.Logger.Infof("task started, nodeID is %s", d.nodePool.GetNodeID())
		}
		d.cr.Start()
	} else {
		logger.Logger.Infof("task have started")
	}
}

// Run 在调用 goroutine 中阻塞运行调度器与节点池；重复启动为无操作。
//
// 启动前先执行已注册的 RecoverFunc；节点池启动失败时回退为停止状态并直接返回。
func (d *Task) Run() {
	if d.RecoverFunc != nil {
		d.RecoverFunc(d)
	}
	if atomic.CompareAndSwapInt32(&d.running, taskStopped, taskRunning) {
		if !d.runningLocally {
			if err := d.startNodePool(); err != nil {
				atomic.StoreInt32(&d.running, taskStopped)
				return
			}
			logger.Logger.Infof("task running, nodeID is %s", d.nodePool.GetNodeID())
		}
		d.cr.Run()
	} else {
		logger.Logger.Infof("task already running")
	}
}

func (d *Task) startNodePool() error {
	if err := d.nodePool.Start(context.Background()); err != nil {
		logger.Logger.Errorf("task start node pool error %+v", err)
		return err
	}
	return nil
}

// Stop 停止节点池并通知调度器停止后续调度；已触发执行中的任务不受影响。
func (d *Task) Stop() {
	tick := time.NewTicker(time.Millisecond)
	if !d.runningLocally {
		d.nodePool.Stop(context.Background())
	}
	for range tick.C {
		if atomic.CompareAndSwapInt32(&d.running, taskRunning, taskStopped) {
			d.cr.Stop()
			logger.Logger.Infof("task stopped")
			return
		}
	}
}

func (d *Task) reRunRecentJobs(jobNames []string) {
	logger.Logger.Infof("reRunRecentJobs: length=%d", len(jobNames))
	for _, jobName := range jobNames {
		if job, ok := d.jobs[jobName]; ok {
			if ok, _ := d.nodePool.CheckJobAvailable(jobName); ok {
				job.Execute()
			}
		}
	}
}

// NodeID 返回当前节点在集群中的标识。
//
// 注意：以 RunningLocally 模式构造的 Task 未创建节点池，调用本方法会触发空指针解引用。
func (d *Task) NodeID() string {
	return d.nodePool.GetNodeID()
}
