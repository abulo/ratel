package task

import (
	"fmt"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/task/cron"
	"github.com/abulo/ratel/v3/core/task/driver"
)

const defaultDuration = 2 * time.Second

type Task struct {
	jobs           map[string]*JobWrapper
	serviceName    string
	updateInterval time.Duration
	node           *Node
	cron           *cron.Cron
	opts           []cron.Option
	isRunning      bool
	lazyPick       bool
}

func NewTask(serviceName string, driver driver.Driver, opts ...Option) *Task {
	Task := &Task{
		serviceName:    serviceName,
		updateInterval: defaultDuration,
		jobs:           make(map[string]*JobWrapper),
		opts:           make([]cron.Option, 0),
	}

	for _, opt := range opts {
		opt(Task)
	}

	Task.cron = cron.New(Task.opts...)

	Task.node = newNode(serviceName, driver, Task, Task.updateInterval)

	return Task
}

// AddJob add a job
func (c *Task) AddJob(jobName, jobType, spec string, job cron.Job) error {
	return c.addJob(jobName, jobType, spec, nil, job)
}

// AddFunc add a cron func
func (c *Task) AddFunc(jobName, jobType, spec string, cmd func()) error {
	return c.addJob(jobName, jobType, spec, cmd, nil)
}

func (c *Task) addJob(jobName, jobType, spec string, cmd func(), job cron.Job) error {
	if _, ok := c.jobs[jobName]; ok {
		return fmt.Errorf("job[%s] already exists", jobName)
	}

	j := &JobWrapper{Job: job, Func: cmd, Task: c, Name: jobName, Type: jobType}

	id, err := c.cron.AddJob(spec, j)
	if err != nil {
		return err
	}

	j.Id = id

	c.jobs[jobName] = j

	return nil
}

func (c *Task) thisNodeRun(jobName string) bool {
	runNodeId, err := c.node.pickNode(jobName)
	if err != nil {
		logger.Logger.Printf("error: pick node failed: [%+v]", err)
		return false
	}
	return c.node.nodeId == runNodeId
}

// Start Task
func (c *Task) Start() {
	c.isRunning = true
	err := c.node.Start()
	if err != nil {
		c.isRunning = false
		logger.Logger.Printf("error: Task start watch failed: [%+v]", err)
		return
	}
	logger.Logger.Printf("info: Task started, nodeId[%s]", c.node.nodeId)
	c.cron.Start()
}

// Run Task
func (c *Task) Run() {
	c.isRunning = true
	err := c.node.Start()
	if err != nil {
		c.isRunning = false
		logger.Logger.Printf("error: Task start watch failed: [%+v]", err)
		return
	}
	logger.Logger.Printf("info: Task running, nodeId[%s]", c.node.nodeId)
	c.cron.Run()
}

// Stop Task
func (c *Task) Stop() {
	c.isRunning = false
	c.node.Stop()
	c.cron.Stop()
}

// Remove Job
func (c *Task) Remove(jobName string) {
	if job, ok := c.jobs[jobName]; ok {
		delete(c.jobs, jobName)
		c.cron.Remove(job.Id)
	}
}
