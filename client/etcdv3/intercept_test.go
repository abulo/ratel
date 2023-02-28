package etcdv3

import (
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

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

func Test_traceStreamClientInterceptor(t *testing.T) {
	tests := []struct {
		name string
		want grpc.StreamClientInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := traceStreamClientInterceptor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("traceStreamClientInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}
