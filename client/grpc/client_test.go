package grpc

import (
	"reflect"
	"testing"

	"google.golang.org/grpc"
)

func Test_newGRPCClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newGRPCClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newGRPCClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGRPCClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDialOptions(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want []grpc.DialOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDialOptions(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDialOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
