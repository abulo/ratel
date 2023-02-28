package etcdv3

import (
	"context"
	"reflect"
	"testing"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func TestNewClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetKeyValue(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantKv  *mvccpb.KeyValue
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			gotKv, err := client.GetKeyValue(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetKeyValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKv, tt.wantKv) {
				t.Errorf("Client.GetKeyValue() = %v, want %v", gotKv, tt.wantKv)
			}
		})
	}
}

func TestClient_GetPrefix(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			got, err := client.GetPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DelPrefix(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDeleted int64
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			gotDeleted, err := client.DelPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DelPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDeleted != tt.wantDeleted {
				t.Errorf("Client.DelPrefix() = %v, want %v", gotDeleted, tt.wantDeleted)
			}
		})
	}
}

func TestClient_GetValues(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			got, err := client.GetValues(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetLeaseSession(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []concurrency.SessionOption
	}
	tests := []struct {
		name             string
		client           *Client
		args             args
		wantLeaseSession *concurrency.Session
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeaseSession, err := tt.client.GetLeaseSession(tt.args.ctx, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLeaseSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLeaseSession, tt.wantLeaseSession) {
				t.Errorf("Client.GetLeaseSession() = %v, want %v", gotLeaseSession, tt.wantLeaseSession)
			}
		})
	}
}
