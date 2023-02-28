package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		opts Options
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Prefix(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Prefix(tt.args.key); got != tt.want {
				t.Errorf("Client.Prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_k(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.k(tt.args.key); got != tt.want {
				t.Errorf("Client.k() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ks(t *testing.T) {
	type args struct {
		key []string
	}
	tests := []struct {
		name string
		r    *Client
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ks(tt.args.key...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClient(t *testing.T) {
	tests := []struct {
		name    string
		r       *Client
		wantRes redis.Cmdable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.r.GetClient(); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Client.GetClient() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestClient_MGetByPipeline(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		r       *Client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.MGetByPipeline(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.MGetByPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.MGetByPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}
