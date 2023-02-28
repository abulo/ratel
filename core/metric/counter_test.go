package metric

import (
	"reflect"
	"testing"
)

func TestCounterVecOpts_Build(t *testing.T) {
	tests := []struct {
		name string
		opts CounterVecOpts
		want *counterVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CounterVecOpts.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCounterVec(t *testing.T) {
	type args struct {
		name   string
		labels []string
	}
	tests := []struct {
		name string
		args args
		want *counterVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounterVec(tt.args.name, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounterVec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counterVec_Inc(t *testing.T) {
	type args struct {
		labels []string
	}
	tests := []struct {
		name    string
		counter *counterVec
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.counter.Inc(tt.args.labels...)
		})
	}
}

func Test_counterVec_Add(t *testing.T) {
	type args struct {
		v      float64
		labels []string
	}
	tests := []struct {
		name    string
		counter *counterVec
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.counter.Add(tt.args.v, tt.args.labels...)
		})
	}
}
