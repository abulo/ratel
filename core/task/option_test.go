package task

import (
	"reflect"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestWithNodeUpdateDuration(t *testing.T) {
	type args struct {
		d time.Duration
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
			if got := WithNodeUpdateDuration(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithNodeUpdateDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithHashReplicas(t *testing.T) {
	type args struct {
		d int
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
			if got := WithHashReplicas(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithHashReplicas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronOptionLocation(t *testing.T) {
	type args struct {
		loc *time.Location
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
			if got := CronOptionLocation(tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CronOptionLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronOptionSeconds(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CronOptionSeconds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CronOptionSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronOptionParser(t *testing.T) {
	type args struct {
		p cron.ScheduleParser
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
			if got := CronOptionParser(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CronOptionParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronOptionChain(t *testing.T) {
	type args struct {
		wrappers []cron.JobWrapper
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
			if got := CronOptionChain(tt.args.wrappers...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CronOptionChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRecoverFunc(t *testing.T) {
	type args struct {
		recoverFunc RecoverFuncType
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
			if got := WithRecoverFunc(tt.args.recoverFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithRecoverFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
