package metric

import (
	"reflect"
	"testing"
)

func TestHistogramVecOpts_Build(t *testing.T) {
	tests := []struct {
		name string
		opts HistogramVecOpts
		want *histogramVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HistogramVecOpts.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_histogramVec_Observe(t *testing.T) {
	type args struct {
		v      float64
		labels []string
	}
	tests := []struct {
		name      string
		histogram *histogramVec
		args      args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.histogram.Observe(tt.args.v, tt.args.labels...)
		})
	}
}
