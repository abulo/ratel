package etcd

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/abulo/ratel/v3/client/etcdv3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestNewEtcdDriver(t *testing.T) {
	type args struct {
		cli *etcdv3.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *EtcdDriver
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEtcdDriver(tt.args.cli)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEtcdDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEtcdDriver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_putKeyWithLease(t *testing.T) {
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		s       *EtcdDriver
		args    args
		want    clientv3.LeaseID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.putKeyWithLease(tt.args.key, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.putKeyWithLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EtcdDriver.putKeyWithLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_randNodeID(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name       string
		s          *EtcdDriver
		args       args
		wantNodeID string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNodeID := tt.s.randNodeID(tt.args.serviceName); gotNodeID != tt.wantNodeID {
				t.Errorf("EtcdDriver.randNodeID() = %v, want %v", gotNodeID, tt.wantNodeID)
			}
		})
	}
}

func TestEtcdDriver_watchService(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name    string
		s       *EtcdDriver
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.watchService(tt.args.serviceName); (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.watchService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getPrefix(t *testing.T) {
	type args struct {
		serviceName string
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
			if got := getPrefix(tt.args.serviceName); got != tt.want {
				t.Errorf("getPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_watcher(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		s    *EtcdDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.watcher(tt.args.serviceName)
		})
	}
}

func TestEtcdDriver_setServiceList(t *testing.T) {
	type args struct {
		serviceName string
		key         string
		val         string
	}
	tests := []struct {
		name string
		s    *EtcdDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.setServiceList(tt.args.serviceName, tt.args.key, tt.args.val)
		})
	}
}

func TestEtcdDriver_delServiceList(t *testing.T) {
	type args struct {
		serviceName string
		key         string
	}
	tests := []struct {
		name string
		s    *EtcdDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.delServiceList(tt.args.serviceName, tt.args.key)
		})
	}
}

func TestEtcdDriver_getServices(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		s    *EtcdDriver
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.getServices(tt.args.serviceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EtcdDriver.getServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_Ping(t *testing.T) {
	tests := []struct {
		name    string
		e       *EtcdDriver
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcdDriver_keepAlive(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID string
	}
	tests := []struct {
		name    string
		e       *EtcdDriver
		args    args
		want    <-chan *clientv3.LeaseKeepAliveResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.keepAlive(tt.args.ctx, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.keepAlive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EtcdDriver.keepAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_revoke(t *testing.T) {
	tests := []struct {
		name string
		e    *EtcdDriver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.revoke()
		})
	}
}

func TestEtcdDriver_SetHeartBeat(t *testing.T) {
	type args struct {
		nodeID string
	}
	tests := []struct {
		name string
		e    *EtcdDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.SetHeartBeat(tt.args.nodeID)
		})
	}
}

func TestEtcdDriver_SetTimeout(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		e    *EtcdDriver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.SetTimeout(tt.args.timeout)
		})
	}
}

func TestEtcdDriver_GetServiceNodeList(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name    string
		e       *EtcdDriver
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.GetServiceNodeList(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.GetServiceNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EtcdDriver.GetServiceNodeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdDriver_RegisterServiceNode(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name    string
		e       *EtcdDriver
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.RegisterServiceNode(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("EtcdDriver.RegisterServiceNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EtcdDriver.RegisterServiceNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
