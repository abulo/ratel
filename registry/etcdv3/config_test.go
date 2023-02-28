package etcdv3

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/registry"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetNode(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetNode(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetNode(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.GetNode(); got != tt.want {
				t.Errorf("Config.GetNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetPrefix(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		config *Config
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetPrefix(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.SetPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    registry.Registry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_MustBuild(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   registry.Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.MustBuild(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.MustBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Singleton(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    registry.Registry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Singleton()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Singleton() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Singleton() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_MustSingleton(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   registry.Registry
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.MustSingleton(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.MustSingleton() = %v, want %v", got, tt.want)
			}
		})
	}
}
