package registry

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/server"
)

func TestService(t *testing.T) {
	type args struct {
		srv server.Server
	}
	tests := []struct {
		name string
		args args
		want server.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Service(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Serve(t *testing.T) {
	tests := []struct {
		name    string
		s       *service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Serve(); (err != nil) != tt.wantErr {
				t.Errorf("service.Serve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
