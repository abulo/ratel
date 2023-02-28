package xgrpc

import (
	"reflect"
	"testing"

	"google.golang.org/grpc"
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

func TestConfig_WithServerOption(t *testing.T) {
	type args struct {
		options []grpc.ServerOption
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
			if got := tt.config.WithServerOption(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithServerOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithStreamInterceptor(t *testing.T) {
	type args struct {
		intes []grpc.StreamServerInterceptor
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
			if got := tt.config.WithStreamInterceptor(tt.args.intes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithStreamInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithUnaryInterceptor(t *testing.T) {
	type args struct {
		intes []grpc.UnaryServerInterceptor
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
			if got := tt.config.WithUnaryInterceptor(tt.args.intes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithUnaryInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_MustBuild(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *Server
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

func TestConfig_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    *Server
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

func TestConfig_Address(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.Address(); got != tt.want {
				t.Errorf("Config.Address() = %v, want %v", got, tt.want)
			}
		})
	}
}
