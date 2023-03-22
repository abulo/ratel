package sql

import (
	"reflect"
	"testing"
)

func TestClient_NewPostgres(t *testing.T) {
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
			if got := tt.c.NewPostgres(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewPostgres() = %v, want %v", got, tt.want)
			}
		})
	}
}
