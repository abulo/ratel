package metric

import (
	"reflect"
	"testing"
)

func TestGaugeVecOpts_Build(t *testing.T) {
	tests := []struct {
		name string
		opts GaugeVecOpts
		want *gaugeVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GaugeVecOpts.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGaugeVec(t *testing.T) {
	type args struct {
		name   string
		labels []string
	}
	tests := []struct {
		name string
		args args
		want *gaugeVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGaugeVec(tt.args.name, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGaugeVec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gaugeVec_Inc(t *testing.T) {
	type args struct {
		labels []string
	}
	tests := []struct {
		name string
		gv   *gaugeVec
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gv.Inc(tt.args.labels...)
		})
	}
}

func Test_gaugeVec_Add(t *testing.T) {
	type args struct {
		v      float64
		labels []string
	}
	tests := []struct {
		name string
		gv   *gaugeVec
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gv.Add(tt.args.v, tt.args.labels...)
		})
	}
}

func Test_gaugeVec_Set(t *testing.T) {
	type args struct {
		v      float64
		labels []string
	}
	tests := []struct {
		name string
		gv   *gaugeVec
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gv.Set(tt.args.v, tt.args.labels...)
		})
	}
}
