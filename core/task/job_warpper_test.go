package task

import "testing"

func TestJobWarpper_Run(t *testing.T) {
	tests := []struct {
		name string
		job  JobWarpper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.job.Run()
		})
	}
}
