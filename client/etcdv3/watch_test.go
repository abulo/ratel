package etcdv3

import (
	"context"
	"reflect"
	"testing"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestWatch_C(t *testing.T) {
	tests := []struct {
		name string
		w    *Watch
		want chan *clientv3.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.C(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Watch.C() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatch_IncipientKeyValues(t *testing.T) {
	tests := []struct {
		name string
		w    *Watch
		want []*mvccpb.KeyValue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.IncipientKeyValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Watch.IncipientKeyValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_WatchPrefix(t *testing.T) {
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    *Watch
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.WatchPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.WatchPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.WatchPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatch_Close(t *testing.T) {
	tests := []struct {
		name    string
		w       *Watch
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Watch.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
