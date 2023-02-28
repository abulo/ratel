package goroutine

import (
	"reflect"
	"testing"
	"time"
)

func TestSerial(t *testing.T) {
	type args struct {
		fns []func()
	}
	tests := []struct {
		name string
		args args
		want func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serial(tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParallel(t *testing.T) {
	type args struct {
		fns []func()
	}
	tests := []struct {
		name string
		args args
		want func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parallel(tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parallel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictParallel(t *testing.T) {
	type args struct {
		restrict int
		fns      []func()
	}
	tests := []struct {
		name string
		args args
		want func()
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RestrictParallel(tt.args.restrict, tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestrictParallel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoDirect(t *testing.T) {
	type args struct {
		fn   interface{}
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoDirect(tt.args.fn, tt.args.args...)
		})
	}
}

func TestGo(t *testing.T) {
	type args struct {
		fn func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Go(tt.args.fn)
		})
	}
}

func TestDelayGo(t *testing.T) {
	type args struct {
		delay time.Duration
		fn    func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DelayGo(tt.args.delay, tt.args.fn)
		})
	}
}

func TestSafeGo(t *testing.T) {
	type args struct {
		fn  func()
		rec func(error)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SafeGo(tt.args.fn, tt.args.rec)
		})
	}
}
