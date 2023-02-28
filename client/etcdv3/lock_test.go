package etcdv3

import (
	"reflect"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"
)

func TestClient_NewMutex(t *testing.T) {
	type args struct {
		key  string
		opts []concurrency.SessionOption
	}
	tests := []struct {
		name      string
		client    *Client
		args      args
		wantMutex *Mutex
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMutex, err := tt.client.NewMutex(tt.args.key, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NewMutex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMutex, tt.wantMutex) {
				t.Errorf("Client.NewMutex() = %v, want %v", gotMutex, tt.wantMutex)
			}
		})
	}
}

func TestMutex_Lock(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name    string
		mutex   *Mutex
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mutex.Lock(tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("Mutex.Lock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMutex_TryLock(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name    string
		mutex   *Mutex
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mutex.TryLock(tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("Mutex.TryLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMutex_Unlock(t *testing.T) {
	tests := []struct {
		name    string
		mutex   *Mutex
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mutex.Unlock(); (err != nil) != tt.wantErr {
				t.Errorf("Mutex.Unlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
