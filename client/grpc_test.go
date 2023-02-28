package client

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/client/grpc"
	"github.com/abulo/ratel/v3/registry/etcdv3"
)

func TestNewGrpcConfig(t *testing.T) {
	tests := []struct {
		name string
		want *GrpcConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGrpcConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGrpcConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcConfig_Store(t *testing.T) {
	type args struct {
		client *grpc.Config
	}
	tests := []struct {
		name  string
		proxy *GrpcConfig
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.proxy.Store(tt.args.client)
		})
	}
}

func TestProxy_StoreGrpc(t *testing.T) {
	type args struct {
		group string
		proxy *GrpcConfig
	}
	tests := []struct {
		name      string
		proxypool *Proxy
		args      args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.proxypool.StoreGrpc(tt.args.group, tt.args.proxy)
		})
	}
}

func TestProxy_LoadGrpc(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		name      string
		proxypool *Proxy
		args      args
		want      *grpc.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.proxypool.LoadGrpc(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.LoadGrpc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEtcdConfig(t *testing.T) {
	tests := []struct {
		name string
		want *EtcdConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEtcdConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEtcdConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdConfig_Store(t *testing.T) {
	type args struct {
		client *etcdv3.Config
	}
	tests := []struct {
		name  string
		proxy *EtcdConfig
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.proxy.Store(tt.args.client)
		})
	}
}

func TestProxy_StoreEtcd(t *testing.T) {
	type args struct {
		group string
		proxy *EtcdConfig
	}
	tests := []struct {
		name      string
		proxypool *Proxy
		args      args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.proxypool.StoreEtcd(tt.args.group, tt.args.proxy)
		})
	}
}

func TestProxy_LoadEtcd(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		name      string
		proxypool *Proxy
		args      args
		want      *etcdv3.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.proxypool.LoadEtcd(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Proxy.LoadEtcd() = %v, want %v", got, tt.want)
			}
		})
	}
}
