package redis

import (
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestClient_GetRingClientConfig(t *testing.T) {
	tests := []struct {
		name string
		o    *Client
		want *redis.RingOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetRingClientConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetRingClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetFailoverClientConfig(t *testing.T) {
	tests := []struct {
		name string
		o    *Client
		want *redis.FailoverOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetFailoverClientConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetFailoverClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClusterClientConfig(t *testing.T) {
	tests := []struct {
		name string
		o    *Client
		want *redis.ClusterOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetClusterClientConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetClusterClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetClientConfig(t *testing.T) {
	tests := []struct {
		name string
		o    *Client
		want *redis.Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetClientConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRWType_IsReadOnly(t *testing.T) {
	tests := []struct {
		name string
		rw   *RWType
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rw.IsReadOnly(); got != tt.want {
				t.Errorf("RWType.IsReadOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRWType_FmtSuffix(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		rw   *RWType
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rw.FmtSuffix(tt.args.key); got != tt.want {
				t.Errorf("RWType.FmtSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
