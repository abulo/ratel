package main

import (
	"fmt"

	"github.com/abulo/ratel/v3/app"
	"github.com/abulo/ratel/v3/logger"
)

// go run main.go --job=jobrunner
func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		logger.Logger.Error(err.Error())
	}
}

// Engine ...
type Engine struct {
	app.Application
}

// NewEngine ...
func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.initJob,
	); err != nil {
		logger.Logger.Panic("startup", err)
	}
	return eng
}

func (e *Engine) initJob() error {
	return e.Job(NewJobRunner())
}

// JobRunner ...
type JobRunner struct {
	JobName string
}

// NewJobRunner ...
func NewJobRunner() *JobRunner {
	return &JobRunner{
		JobName: "jobrunner",
	}
}

// Run ...
func (j *JobRunner) Run() {
	fmt.Println("i am job runner")
}

// GetJobName ...
func (j *JobRunner) GetJobName() string {
	return j.JobName
}
