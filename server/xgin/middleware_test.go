package xgin

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_extractAID(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractAID(tt.args.ctx); got != tt.want {
				t.Errorf("extractAID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metricServerInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metricServerInterceptor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metricServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_traceServerInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := traceServerInterceptor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("traceServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_recoverMiddleware(t *testing.T) {
	type args struct {
		slowQueryThresholdInMilli int64
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := recoverMiddleware(tt.args.slowQueryThresholdInMilli); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("recoverMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stack(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stack(tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_function(t *testing.T) {
	type args struct {
		pc uintptr
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := function(tt.args.pc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("function() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_source(t *testing.T) {
	type args struct {
		lines [][]byte
		n     int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := source(tt.args.lines, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source() = %v, want %v", got, tt.want)
			}
		})
	}
}
