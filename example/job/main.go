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

type Engine struct {
	app.Application
}

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

type JobRunner struct {
	JobName string
}

func NewJobRunner() *JobRunner {
	return &JobRunner{
		JobName: "jobrunner",
	}
}

func (j *JobRunner) Run() {
	fmt.Println("i am job runner")
}

func (j *JobRunner) GetJobName() string {
	return j.JobName
}
