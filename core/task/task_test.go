package task

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/core/task/driver"
	"github.com/robfig/cron/v3"
)

func TestNewTask(t *testing.T) {
	type args struct {
		serverName string
		driver     driver.Driver
		cronOpts   []cron.Option
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTask(tt.args.serverName, tt.args.driver, tt.args.cronOpts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTaskWithOption(t *testing.T) {
	type args struct {
		serverName string
		driver     driver.Driver
		taskOpts   []Option
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskWithOption(tt.args.serverName, tt.args.driver, tt.args.taskOpts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskWithOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTask(t *testing.T) {
	type args struct {
		serverName string
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTask(tt.args.serverName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_AddJob(t *testing.T) {
	type args struct {
		jobName string
		cronStr string
		job     Job
	}
	tests := []struct {
		name    string
		d       *Task
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.AddJob(tt.args.jobName, tt.args.cronStr, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Task.AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_AddFunc(t *testing.T) {
	type args struct {
		jobName string
		cronStr string
		cmd     func()
	}
	tests := []struct {
		name    string
		d       *Task
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.AddFunc(tt.args.jobName, tt.args.cronStr, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Task.AddFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Total(t *testing.T) {
	tests := []struct {
		name string
		d    *Task
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Total(); got != tt.want {
				t.Errorf("Task.Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_addJob(t *testing.T) {
	type args struct {
		jobName string
		cronStr string
		cmd     func()
		job     Job
	}
	tests := []struct {
		name    string
		d       *Task
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.addJob(tt.args.jobName, tt.args.cronStr, tt.args.cmd, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Task.addJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Remove(t *testing.T) {
	type args struct {
		jobName string
	}
	tests := []struct {
		name string
		d    *Task
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Remove(tt.args.jobName)
		})
	}
}

func TestTask_allowThisNodeRun(t *testing.T) {
	type args struct {
		jobName string
	}
	tests := []struct {
		name string
		d    *Task
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.allowThisNodeRun(tt.args.jobName); got != tt.want {
				t.Errorf("Task.allowThisNodeRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Start(t *testing.T) {
	tests := []struct {
		name string
		d    *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Start()
		})
	}
}

func TestTask_Run(t *testing.T) {
	tests := []struct {
		name string
		d    *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Run()
		})
	}
}

func TestTask_startNodePool(t *testing.T) {
	tests := []struct {
		name    string
		d       *Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.startNodePool(); (err != nil) != tt.wantErr {
				t.Errorf("Task.startNodePool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Stop(t *testing.T) {
	tests := []struct {
		name string
		d    *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Stop()
		})
	}
}

func TestTask_WorkerStart(t *testing.T) {
	tests := []struct {
		name    string
		tr      *Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.WorkerStart(); (err != nil) != tt.wantErr {
				t.Errorf("Task.WorkerStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_WorkerStop(t *testing.T) {
	tests := []struct {
		name    string
		tr      *Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.WorkerStop(); (err != nil) != tt.wantErr {
				t.Errorf("Task.WorkerStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
