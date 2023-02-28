package util

import (
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	type args struct {
		attempts int
		sleep    time.Duration
		f        func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.attempts, tt.args.sleep, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
