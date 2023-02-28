package xgin

import (
	"context"
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/server"
	"github.com/gin-gonic/gin"
)

func Test_newServer(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newServer(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Upgrade(t *testing.T) {
	type args struct {
		ws *WebSocket
	}
	tests := []struct {
		name string
		s    *Server
		args args
		want gin.IRoutes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Upgrade(tt.args.ws); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Upgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Serve(t *testing.T) {
	tests := []struct {
		name    string
		s       *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Serve(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Serve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Stop(t *testing.T) {
	tests := []struct {
		name    string
		s       *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Server.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_GracefulStop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		s       *Server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.GracefulStop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.GracefulStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Info(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want *server.ServiceInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Info(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Health(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Health(); got != tt.want {
				t.Errorf("Server.Health() = %v, want %v", got, tt.want)
			}
		})
	}
}
