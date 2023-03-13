package sql

import (
	"reflect"
	"testing"
)

func TestClient_NewMysql(t *testing.T) {
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
			if got := tt.c.NewMysql(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewMysql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mysqlAcceptable(t *testing.T) {
	type args struct {
		err error
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
			if got := mysqlAcceptable(tt.args.err); got != tt.want {
				t.Errorf("mysqlAcceptable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withMysqlAcceptable(t *testing.T) {
	tests := []struct {
		name string
		want SqlOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withMysqlAcceptable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withMysqlAcceptable() = %v, want %v", got, tt.want)
			}
		})
	}
}
