package task

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/driver"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
)

const (
	defaultReplicas = 50
	defaultDuration = time.Second
)

const (
	taskRunning = 1
	taskStopped = 0
)

type RecoverFuncType func(d *Task)

// Task is main struct
type Task struct {
	jobs               map[string]*JobWarpper
	jobsRWMut          sync.Mutex
	ServerName         string
	nodePool           *NodePool
	running            int32
	nodeUpdateDuration time.Duration
	hashReplicas       int
	cr                 *cron.Cron
	crOptions          []cron.Option
	RecoverFunc        RecoverFuncType
}

// NewTask create a Task
func NewTask(serverName string, driver driver.Driver, cronOpts ...cron.Option) *Task {
	task := newTask(serverName)
	task.crOptions = cronOpts
	task.cr = cron.New(cronOpts...)
	task.running = taskStopped
	var err error
	task.nodePool, err = newNodePool(serverName, driver, task, task.nodeUpdateDuration, task.hashReplicas)
	if err != nil {
		logger.Logger.Errorf("ERR: %s", err.Error())
		return nil
	}
	return task
}

// NewTaskWithOption create a Task with Task Option
func NewTaskWithOption(serverName string, driver driver.Driver, taskOpts ...Option) *Task {
	task := newTask(serverName)
	for _, opt := range taskOpts {
		opt(task)
	}

	task.cr = cron.New(task.crOptions...)
	var err error
	task.nodePool, err = newNodePool(serverName, driver, task, task.nodeUpdateDuration, task.hashReplicas)
	if err != nil {
		logger.Logger.Errorf("ERR: %s", err.Error())
		return nil
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
	return d.addJob(jobName, cronStr, nil, job)
}

// AddFunc add a cron func
func (d *Task) AddFunc(jobName, cronStr string, cmd func()) (err error) {
	return d.addJob(jobName, cronStr, cmd, nil)
}

// Total get total job
func (d *Task) Total() int64 {
	return cast.ToInt64(d.cr.Entries())
}

func (d *Task) addJob(jobName, cronStr string, cmd func(), job Job) (err error) {
	logger.Logger.Infof("addJob '%s' :  %s", jobName, cronStr)

	d.jobsRWMut.Lock()
	defer d.jobsRWMut.Unlock()
	if _, ok := d.jobs[jobName]; ok {
		return errors.New("jobName already exist")
	}
	innerJob := JobWarpper{
		Name:    jobName,
		CronStr: cronStr,
		Func:    cmd,
		Job:     job,
		Task:    d,
	}
	entryID, err := d.cr.AddJob(cronStr, innerJob)
	if err != nil {
		return err
	}
	innerJob.ID = entryID
	d.jobs[jobName] = &innerJob
	return nil
}

// Remove Job
func (d *Task) Remove(jobName string) {
	d.jobsRWMut.Lock()
	defer d.jobsRWMut.Unlock()

	if job, ok := d.jobs[jobName]; ok {
		delete(d.jobs, jobName)
		d.cr.Remove(job.ID)
	}
}

func (d *Task) allowThisNodeRun(jobName string) bool {
	allowRunNode := d.nodePool.PickNodeByJobName(jobName)
	logger.Logger.Infof("job '%s' running in node %s", jobName, allowRunNode)
	if allowRunNode == "" {
		logger.Logger.Errorf("node pool is empty")
		return false
	}
	return d.nodePool.NodeID == allowRunNode
}

// Start job
func (d *Task) Start() {
	// recover jobs before starting
	if d.RecoverFunc != nil {
		d.RecoverFunc(d)
	}

	if atomic.CompareAndSwapInt32(&d.running, taskStopped, taskRunning) {
		if err := d.startNodePool(); err != nil {
			atomic.StoreInt32(&d.running, taskStopped)
			return
		}
		d.cr.Start()
		logger.Logger.Infof("task started , nodeID is %s", d.nodePool.NodeID)
	} else {
		logger.Logger.Infof("task have started")
	}
}

// Run Job
func (d *Task) Run() {
	if atomic.CompareAndSwapInt32(&d.running, taskStopped, taskRunning) {
		if err := d.startNodePool(); err != nil {
			atomic.StoreInt32(&d.running, taskStopped)
			return
		}

		logger.Logger.Infof("task running nodeID is %s", d.nodePool.NodeID)
		d.cr.Run()
	} else {
		logger.Logger.Infof("task already running")
	}
}

func (d *Task) startNodePool() error {
	if err := d.nodePool.StartPool(); err != nil {
		logger.Logger.Errorf("task start node pool error %+v", err)
		return err
	}
	return nil
}

// Stop job
func (d *Task) Stop() {
	tick := time.NewTicker(time.Millisecond)
	for range tick.C {
		if atomic.CompareAndSwapInt32(&d.running, taskRunning, taskStopped) {
			d.cr.Stop()
			logger.Logger.Infof("task stopped")
			return
		}
	}
}

// WorkerStart ...
func (t *Task) WorkerStart() error {
	t.Start()
	return nil
}

// WorkerStop ...
func (t *Task) WorkerStop() error {
	t.Stop()
	return nil
}
