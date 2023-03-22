package sql

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
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
			got, err := NewClient(tt.args.opts...)
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

func TestClient_dns(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.dns(); got != tt.want {
				t.Errorf("Client.dns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_NewSqlClient(t *testing.T) {
	type args struct {
		opts []SqlOption
	}
	tests := []struct {
		name string
		c    *Client
		args args
		want SqlConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.NewSqlClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewSqlClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
