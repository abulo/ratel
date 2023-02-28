package registry

import (
	"context"
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/server"
)

func TestLocal_ListServices(t *testing.T) {
	type args struct {
		ctx context.Context
		s   string
	}
	tests := []struct {
		name    string
		n       Local
		args    args
		want    []*server.ServiceInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.ListServices(tt.args.ctx, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.ListServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Local.ListServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocal_WatchServices(t *testing.T) {
	type args struct {
		ctx context.Context
		s   string
	}
	tests := []struct {
		name    string
		n       Local
		args    args
		want    chan Endpoints
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.WatchServices(tt.args.ctx, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.WatchServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Local.WatchServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocal_RegisterService(t *testing.T) {
	type args struct {
		ctx context.Context
		si  *server.ServiceInfo
	}
	tests := []struct {
		name    string
		n       Local
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.RegisterService(tt.args.ctx, tt.args.si); (err != nil) != tt.wantErr {
				t.Errorf("Local.RegisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_UnregisterService(t *testing.T) {
	type args struct {
		ctx context.Context
		si  *server.ServiceInfo
	}
	tests := []struct {
		name    string
		n       Local
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnregisterService(tt.args.ctx, tt.args.si); (err != nil) != tt.wantErr {
				t.Errorf("Local.UnregisterService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_Close(t *testing.T) {
	tests := []struct {
		name    string
		n       Local
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Local.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocal_Kind(t *testing.T) {
	tests := []struct {
		name string
		n    Local
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Kind(); got != tt.want {
				t.Errorf("Local.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}
