package xhertz

import (
	"reflect"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
)

func Test_metricServerInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want app.HandlerFunc
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
		want app.HandlerFunc
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
		want app.HandlerFunc
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

func Test_extractAID(t *testing.T) {
	type args struct {
		ctx *app.RequestContext
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

func TestHTTPHeader_Visit(t *testing.T) {
	type args struct {
		v func(k, v string)
	}
	tests := []struct {
		name string
		h    HTTPHeader
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Visit(tt.args.v)
		})
	}
}

func TestHTTPHeader_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		h    HTTPHeader
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestHTTPHeaderToCGIVariable(t *testing.T) {
	type args struct {
		key string
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
			if got := HTTPHeaderToCGIVariable(tt.args.key); got != tt.want {
				t.Errorf("HTTPHeaderToCGIVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCGIVariableToHTTPHeader(t *testing.T) {
	type args struct {
		key string
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
			if got := CGIVariableToHTTPHeader(tt.args.key); got != tt.want {
				t.Errorf("CGIVariableToHTTPHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
