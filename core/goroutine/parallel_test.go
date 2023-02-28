package goroutine

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

var (
	fn1     = func() error { return nil }
	fn2     = func() error { return errors.New("BOOM") }
	timeout = time.After(2 * time.Second)
)

func testRun(t *testing.T) {
	var count int
	err := ParallelWithErrorChan(fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimit(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(2, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimitWithNegativeConcurrencyValue(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(-1, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func testRunLimitWithConcurrencyGreaterThanPassedFunctions(t *testing.T) {
	var count int
	err := RestrictParallelWithErrorChan(3, fn1, fn2)
outer:
	for {
		select {
		case <-err:
			count++
			if count == 2 {
				break outer
			}
		case <-timeout:
			t.Errorf("parallel.Run() failed, got timeout error")
			break outer
		}
	}

	if count != 2 {
		t.Errorf("parallel.Run() failed, got '%v', expected '%v'", count, 2)
	}
}

func TestParallelWithError(t *testing.T) {
	type args struct {
		fns []func() error
	}
	tests := []struct {
		name string
		args args
		want func() error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelWithError(tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelWithError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParallelWithErrorChan(t *testing.T) {
	type args struct {
		fns []func() error
	}
	tests := []struct {
		name string
		args args
		want chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelWithErrorChan(tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParallelWithErrorChan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictParallelWithErrorChan(t *testing.T) {
	type args struct {
		concurrency int
		fns         []func() error
	}
	tests := []struct {
		name string
		args args
		want chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RestrictParallelWithErrorChan(tt.args.concurrency, tt.args.fns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestrictParallelWithErrorChan() = %v, want %v", got, tt.want)
			}
		})
	}
}
