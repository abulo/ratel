package etcdv3

import (
	"context"
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/registry"
	"github.com/abulo/ratel/v3/server"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Test_newETCDRegistry(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *etcdv3Registry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newETCDRegistry(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newETCDRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newETCDRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_Kind(t *testing.T) {
	tests := []struct {
		name string
		reg  *etcdv3Registry
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reg.Kind(); got != tt.want {
				t.Errorf("etcdv3Registry.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_RegisterService(t *testing.T) {
	type args struct {
		ctx  context.Context
		info *server.ServiceInfo
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.RegisterService(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.RegisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_UnregisterService(t *testing.T) {
	type args struct {
		ctx  context.Context
		info *server.ServiceInfo
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.UnregisterService(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.UnregisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_ListServices(t *testing.T) {
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name         string
		reg          *etcdv3Registry
		args         args
		wantServices []*server.ServiceInfo
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServices, err := tt.reg.ListServices(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.ListServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotServices, tt.wantServices) {
				t.Errorf("etcdv3Registry.ListServices() = %v, want %v", gotServices, tt.wantServices)
			}
		})
	}
}

func Test_etcdv3Registry_WatchServices(t *testing.T) {
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		want    chan registry.Endpoints
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.WatchServices(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.WatchServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("etcdv3Registry.WatchServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_unregister(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.unregister(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.unregister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_Close(t *testing.T) {
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.Close(); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_registerMetric(t *testing.T) {
	type args struct {
		ctx  context.Context
		info *server.ServiceInfo
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.registerMetric(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.registerMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_registerBiz(t *testing.T) {
	type args struct {
		ctx  context.Context
		info *server.ServiceInfo
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.registerBiz(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.registerBiz() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_registerKV(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		val string
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.registerKV(tt.args.ctx, tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.registerKV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_etcdv3Registry_getOrGrantLeaseID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		want    clientv3.LeaseID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reg.getOrGrantLeaseID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.getOrGrantLeaseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("etcdv3Registry.getOrGrantLeaseID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_getLeaseID(t *testing.T) {
	tests := []struct {
		name string
		reg  *etcdv3Registry
		want clientv3.LeaseID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reg.getLeaseID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("etcdv3Registry.getLeaseID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_setLeaseID(t *testing.T) {
	type args struct {
		leaseId clientv3.LeaseID
	}
	tests := []struct {
		name string
		reg  *etcdv3Registry
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reg.setLeaseID(tt.args.leaseId)
		})
	}
}

func Test_etcdv3Registry_doKeepalive(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		reg  *etcdv3Registry
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.reg.doKeepalive(tt.args.ctx)
		})
	}
}

func Test_etcdv3Registry_registerKey(t *testing.T) {
	type args struct {
		info *server.ServiceInfo
	}
	tests := []struct {
		name string
		reg  *etcdv3Registry
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reg.registerKey(tt.args.info); got != tt.want {
				t.Errorf("etcdv3Registry.registerKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_registerValue(t *testing.T) {
	type args struct {
		info *server.ServiceInfo
	}
	tests := []struct {
		name string
		reg  *etcdv3Registry
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reg.registerValue(tt.args.info); got != tt.want {
				t.Errorf("etcdv3Registry.registerValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_etcdv3Registry_registerAllKvs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		reg     *etcdv3Registry
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reg.registerAllKvs(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("etcdv3Registry.registerAllKvs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteAddrList(t *testing.T) {
	type args struct {
		al     *registry.Endpoints
		prefix string
		scheme string
		kvs    []*mvccpb.KeyValue
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteAddrList(tt.args.al, tt.args.prefix, tt.args.scheme, tt.args.kvs...)
		})
	}
}

func Test_updateAddrList(t *testing.T) {
	type args struct {
		al     *registry.Endpoints
		prefix string
		scheme string
		kvs    []*mvccpb.KeyValue
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateAddrList(tt.args.al, tt.args.prefix, tt.args.scheme, tt.args.kvs...)
		})
	}
}

func Test_isIPPort(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIPPort(tt.args.addr); got != tt.want {
				t.Errorf("isIPPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getScheme(t *testing.T) {
	type args struct {
		prefix string
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
			if got := getScheme(tt.args.prefix); got != tt.want {
				t.Errorf("getScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
