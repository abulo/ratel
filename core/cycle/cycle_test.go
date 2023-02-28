package cycle

import (
	"reflect"
	"testing"
)

func TestNewCycle(t *testing.T) {
	tests := []struct {
		name string
		want *Cycle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCycle(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCycle_Run(t *testing.T) {
	type args struct {
		fn func() error
	}
	tests := []struct {
		name string
		c    *Cycle
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Run(tt.args.fn)
		})
	}
}

func TestCycle_Done(t *testing.T) {
	tests := []struct {
		name string
		c    *Cycle
		want <-chan struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Done(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cycle.Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCycle_DoneAndClose(t *testing.T) {
	tests := []struct {
		name string
		c    *Cycle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.DoneAndClose()
		})
	}
}

func TestCycle_Close(t *testing.T) {
	tests := []struct {
		name string
		c    *Cycle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Close()
		})
	}
}

func TestCycle_Wait(t *testing.T) {
	tests := []struct {
		name string
		c    *Cycle
		want <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Wait(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cycle.Wait() = %v, want %v", got, tt.want)
			}
		})
	}
}
