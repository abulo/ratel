package sql

import (
	"reflect"
	"testing"
)

func TestClient_NewClickhouse(t *testing.T) {
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
			if got := tt.c.NewClickhouse(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewClickhouse() = %v, want %v", got, tt.want)
			}
		})
	}
}
