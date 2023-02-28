package metric

import (
	"reflect"
	"testing"
)

func TestSummaryVecOpts_Build(t *testing.T) {
	tests := []struct {
		name string
		opts SummaryVecOpts
		want *summaryVec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SummaryVecOpts.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_summaryVec_Observe(t *testing.T) {
	type args struct {
		v      float64
		labels []string
	}
	tests := []struct {
		name    string
		summary *summaryVec
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.summary.Observe(tt.args.v, tt.args.labels...)
		})
	}
}
