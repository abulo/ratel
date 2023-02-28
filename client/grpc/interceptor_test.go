package grpc

import (
	"context"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func Test_debugUnaryClientInterceptor(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := debugUnaryClientInterceptor(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("debugUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aidUnaryClientInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want grpc.UnaryClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aidUnaryClientInterceptor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aidUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeoutUnaryClientInterceptor(t *testing.T) {
	type args struct {
		timeout       time.Duration
		slowThreshold time.Duration
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeoutUnaryClientInterceptor(tt.args.timeout, tt.args.slowThreshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeoutUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_traceUnaryClientInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want grpc.UnaryClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := traceUnaryClientInterceptor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("traceUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggerUnaryClientInterceptor(t *testing.T) {
	type args struct {
		name                   string
		accessInterceptorLevel string
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggerUnaryClientInterceptor(tt.args.name, tt.args.accessInterceptorLevel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loggerUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metricUnaryClientInterceptor(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := metricUnaryClientInterceptor(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metricUnaryClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}
