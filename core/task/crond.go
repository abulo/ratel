package task

import (
	"fmt"
	"log"
	"time"

	"github.com/abulo/ratel/v3/core/task/cron"
	"github.com/abulo/ratel/v3/core/task/driver"
)

const defaultDuration = 2 * time.Second

type Crond struct {
	jobs           map[string]*JobWrapper
	serviceName    string
	updateInterval time.Duration
	node           *Node
	cron           *cron.Cron
	opts           []cron.Option
	isRunning      bool
	lazyPick       bool
}

func NewCrond(serviceName string, driver driver.Driver, opts ...Option) *Crond {
	crond := &Crond{
		serviceName:    serviceName,
		updateInterval: defaultDuration,
		jobs:           make(map[string]*JobWrapper),
		opts:           make([]cron.Option, 0),
	}

	for _, opt := range opts {
		opt(crond)
	}

	crond.cron = cron.New(crond.opts...)

	crond.node = newNode(serviceName, driver, crond, crond.updateInterval)

	return crond
}

// AddJob add a job
func (c *Crond) AddJob(jobName, jobType, spec string, job cron.Job) error {
	return c.addJob(jobName, jobType, spec, nil, job)
}

// AddFunc add a cron func
func (c *Crond) AddFunc(jobName, jobType, spec string, cmd func()) error {
	return c.addJob(jobName, jobType, spec, cmd, nil)
}

func (c *Crond) addJob(jobName, jobType, spec string, cmd func(), job cron.Job) error {
	if _, ok := c.jobs[jobName]; ok {
		return fmt.Errorf("job[%s] already exists", jobName)
	}

	j := &JobWrapper{Job: job, Func: cmd, Crond: c, Name: jobName, Type: jobType}

	id, err := c.cron.AddJob(spec, j)
	if err != nil {
		return err
	}

	j.Id = id

	c.jobs[jobName] = j

	return nil
}

func (c *Crond) thisNodeRun(jobName string) bool {
	runNodeId, err := c.node.pickNode(jobName)
	if err != nil {
		log.Printf("error: pick node failed: [%+v]", err)
		return false
	}

	log.Printf("info: node[%s] will run job[%s]", runNodeId, jobName)

	return c.node.nodeId == runNodeId
}

// Start Crond
func (c *Crond) Start() {
	c.isRunning = true

	err := c.node.Start()
	if err != nil {
		c.isRunning = false
		log.Printf("error: crond start watch failed: [%+v]", err)
		return
	}

	log.Printf("info: crond started, nodeId[%s]", c.node.nodeId)

	c.cron.Start()
}

// Run Crond
func (c *Crond) Run() {
	c.isRunning = true

	err := c.node.Start()
	if err != nil {
		c.isRunning = false
		log.Printf("error: crond start watch failed: [%+v]", err)
		return
	}

	log.Printf("info: crond running, nodeId[%s]", c.node.nodeId)

	c.cron.Run()
}

// Stop Crond
func (c *Crond) Stop() {
	c.isRunning = false
	c.node.Stop()
	c.cron.Stop()
}

// Remove Job
func (c *Crond) Remove(jobName string) {
	if job, ok := c.jobs[jobName]; ok {
		delete(c.jobs, jobName)
		c.cron.Remove(job.Id)
	}
}
