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
	ErrJobExist     = errors.New("jobName already exist")
	ErrJobNotExist  = errors.New("jobName not exist")
	ErrJobWrongNode = errors.New("job is not running in this node")
)

type RecoverFuncType func(d *Task)

// Task is main struct
type Task struct {
	jobs               map[string]*JobWarpper
	jobsRWMut          sync.RWMutex
	ServerName         string
	nodePool           INodePool
	running            int32
	nodeUpdateDuration time.Duration
	hashReplicas       int
	cr                 *cron.Cron
	crOptions          []cron.Option
	RecoverFunc        RecoverFuncType
	recentJobs         IRecentJobPacker
	state              atomic.Value
	runningLocally     bool
}

// NewTask create a Task
func NewTask(serverName string, driver driver.DriverV2, cronOpts ...cron.Option) *Task {
	task := newTask(serverName)
	task.crOptions = cronOpts
	task.cr = cron.New(cronOpts...)
	task.running = taskStopped
	task.nodePool = NewNodePool(serverName, driver, task.nodeUpdateDuration, task.hashReplicas)
	return task
}

// NewTaskWithOption create a Task with Task Option
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

// AddJob  add a job
func (d *Task) AddJob(jobName, cronStr string, job Job) (err error) {
	return d.addJob(jobName, cronStr, job)
}

// AddFunc add a cron func
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

// Remove Job by jobName
func (d *Task) Remove(jobName string) {
	d.jobsRWMut.Lock()
	defer d.jobsRWMut.Unlock()

	if job, ok := d.jobs[jobName]; ok {
		delete(d.jobs, jobName)
		d.cr.Remove(job.ID)
	}
}

// Get job by jobName
// if this jobName not exist, will return error.
//
//	if `thisNodeOnly` is true
//		if this job is not available in this node, will return error.
//	otherwise return the struct of JobWarpper whose name is jobName.
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

// Get job list.
//
//	if `thisNodeOnly` is true
//		return all jobs available in this node.
//	otherwise return all jobs added to task.
//
// we never return nil. If there is no job.
// this func will return an empty slice.
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

// Start job
func (d *Task) Start() {
	// recover jobs before starting
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

// Run Job
func (d *Task) Run() {
	// recover jobs before starting
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

// Stop job
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

func (d *Task) NodeID() string {
	return d.nodePool.GetNodeID()
}
