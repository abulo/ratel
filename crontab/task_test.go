package crontab

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestNewTaskManager(t *testing.T) {
	tests := []struct {
		name string
		want *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	type args struct {
		tname string
		spec  string
		f     TaskFunc
		opts  []Option
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
			if got := NewTask(tt.args.tname, tt.args.spec, tt.args.f, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_GetSpec(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetSpec(tt.args.in0); got != tt.want {
				t.Errorf("Task.GetSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_GetStatus(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetStatus(tt.args.in0); got != tt.want {
				t.Errorf("Task.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		tr      *Task
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Task.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_SetNext(t *testing.T) {
	type args struct {
		ctx context.Context
		now time.Time
	}
	tests := []struct {
		name string
		tr   *Task
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetNext(tt.args.ctx, tt.args.now)
		})
	}
}

func TestTask_GetNext(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetNext(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.GetNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_SetPrev(t *testing.T) {
	type args struct {
		ctx context.Context
		now time.Time
	}
	tests := []struct {
		name string
		tr   *Task
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetPrev(tt.args.ctx, tt.args.now)
		})
	}
}

func TestTask_GetPrev(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetPrev(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.GetPrev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_GetTimeout(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetTimeout(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.GetTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionFunc_apply(t *testing.T) {
	type args struct {
		t *Task
	}
	tests := []struct {
		name string
		f    optionFunc
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.apply(tt.args.t)
		})
	}
}

func TestTimeoutOption(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeoutOption(tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeoutOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_SetCron(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name string
		tr   *Task
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetCron(tt.args.spec)
		})
	}
}

func TestTask_parse(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want *Schedule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.parse(tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_parseSpec(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name string
		tr   *Task
		args args
		want *Schedule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.parseSpec(tt.args.spec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.parseSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_Next(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		s    *Schedule
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Next(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Schedule.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dayMatches(t *testing.T) {
	type args struct {
		s *Schedule
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dayMatches(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("dayMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskManager_WorkerStart(t *testing.T) {
	tests := []struct {
		name    string
		tr      *TaskManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.WorkerStart(); (err != nil) != tt.wantErr {
				t.Errorf("TaskManager.WorkerStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskManager_WorkerStop(t *testing.T) {
	tests := []struct {
		name    string
		tr      *TaskManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.WorkerStop(); (err != nil) != tt.wantErr {
				t.Errorf("TaskManager.WorkerStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskManager_StartTask(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.StartTask()
		})
	}
}

func TestTaskManager_run(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.run()
		})
	}
}

func TestTaskManager_setTasksStartTime(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		name string
		tr   *TaskManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.setTasksStartTime(tt.args.now)
		})
	}
}

func TestTaskManager_markManagerStop(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.markManagerStop()
		})
	}
}

func TestTaskManager_runNextTasks(t *testing.T) {
	type args struct {
		sortList  *MapSorter
		effective time.Time
	}
	tests := []struct {
		name string
		tr   *TaskManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.runNextTasks(tt.args.sortList, tt.args.effective)
		})
	}
}

func TestTaskManager_StopTask(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.StopTask()
		})
	}
}

func TestTaskManager_GracefulShutdown(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
		want <-chan struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GracefulShutdown(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskManager.GracefulShutdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskManager_AddTask(t *testing.T) {
	type args struct {
		taskname string
		task     Tasker
	}
	tests := []struct {
		name string
		tr   *TaskManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.AddTask(tt.args.taskname, tt.args.task)
		})
	}
}

func TestTaskManager_Len(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Len(); got != tt.want {
				t.Errorf("TaskManager.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskManager_DeleteTask(t *testing.T) {
	type args struct {
		taskname string
	}
	tests := []struct {
		name string
		tr   *TaskManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.DeleteTask(tt.args.taskname)
		})
	}
}

func TestTaskManager_ClearTask(t *testing.T) {
	tests := []struct {
		name string
		tr   *TaskManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.ClearTask()
		})
	}
}

func TestNewMapSorter(t *testing.T) {
	type args struct {
		m map[string]Tasker
	}
	tests := []struct {
		name string
		args args
		want *MapSorter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapSorter(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMapSorter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapSorter_Sort(t *testing.T) {
	tests := []struct {
		name string
		ms   *MapSorter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.Sort()
		})
	}
}

func TestMapSorter_Len(t *testing.T) {
	tests := []struct {
		name string
		ms   *MapSorter
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.Len(); got != tt.want {
				t.Errorf("MapSorter.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapSorter_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		ms   *MapSorter
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("MapSorter.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapSorter_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		ms   *MapSorter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.Swap(tt.args.i, tt.args.j)
		})
	}
}

func Test_getField(t *testing.T) {
	type args struct {
		field string
		r     bounds
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getField(tt.args.field, tt.args.r); got != tt.want {
				t.Errorf("getField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRange(t *testing.T) {
	type args struct {
		expr string
		r    bounds
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRange(tt.args.expr, tt.args.r); got != tt.want {
				t.Errorf("getRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIntOrName(t *testing.T) {
	type args struct {
		expr  string
		names map[string]uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIntOrName(tt.args.expr, tt.args.names); got != tt.want {
				t.Errorf("parseIntOrName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mustParseInt(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParseInt(tt.args.expr); got != tt.want {
				t.Errorf("mustParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBits(t *testing.T) {
	type args struct {
		min  uint
		max  uint
		step uint
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBits(tt.args.min, tt.args.max, tt.args.step); got != tt.want {
				t.Errorf("getBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_all(t *testing.T) {
	type args struct {
		r bounds
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := all(tt.args.r); got != tt.want {
				t.Errorf("all() = %v, want %v", got, tt.want)
			}
		})
	}
}
