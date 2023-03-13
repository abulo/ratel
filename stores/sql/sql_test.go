package sql

import (
	"reflect"
	"testing"
)

func TestNewSqlClient(t *testing.T) {
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
			got, err := NewSqlClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSqlClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSqlClient() = %v, want %v", got, tt.want)
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
