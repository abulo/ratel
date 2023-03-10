package xgrpc

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func TestStreamInterceptorChain(t *testing.T) {
	type args struct {
		interceptors []grpc.StreamServerInterceptor
	}
	tests := []struct {
		name string
		args args
		want grpc.StreamServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StreamInterceptorChain(tt.args.interceptors...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StreamInterceptorChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnaryInterceptorChain(t *testing.T) {
	type args struct {
		interceptors []grpc.UnaryServerInterceptor
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnaryInterceptorChain(tt.args.interceptors...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryInterceptorChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prometheusUnaryServerInterceptor(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prometheusUnaryServerInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("prometheusUnaryServerInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prometheusUnaryServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prometheusStreamServerInterceptor(t *testing.T) {
	type args struct {
		srv     any
		ss      grpc.ServerStream
		info    *grpc.StreamServerInfo
		handler grpc.StreamHandler
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
			if err := prometheusStreamServerInterceptor(tt.args.srv, tt.args.ss, tt.args.info, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("prometheusStreamServerInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_traceUnaryServerInterceptor(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := traceUnaryServerInterceptor(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("traceUnaryServerInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("traceUnaryServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contextedServerStream_Context(t *testing.T) {
	tests := []struct {
		name string
		css  contextedServerStream
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.css.Context(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("contextedServerStream.Context() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_traceStreamServerInterceptor(t *testing.T) {
	type args struct {
		srv     any
		ss      grpc.ServerStream
		info    *grpc.StreamServerInfo
		handler grpc.StreamHandler
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
			if err := traceStreamServerInterceptor(tt.args.srv, tt.args.ss, tt.args.info, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("traceStreamServerInterceptor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractAID(t *testing.T) {
	type args struct {
		ctx context.Context
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
