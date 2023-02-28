package redis

import (
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestOptions_GetClusterConfig(t *testing.T) {
	tests := []struct {
		name string
		o    Options
		want *redis.ClusterOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetClusterConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.GetClusterConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_GetNormalConfig(t *testing.T) {
	tests := []struct {
		name string
		o    Options
		want *redis.Options
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetNormalConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.GetNormalConfig() = %v, want %v", got, tt.want)
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
