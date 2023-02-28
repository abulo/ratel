package xmonitor

import (
	"reflect"
	"testing"
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

func TestConfig_WithHost(t *testing.T) {
	type args struct {
		host string
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
			if got := tt.config.WithHost(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithPort(t *testing.T) {
	type args struct {
		port int
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
			if got := tt.config.WithPort(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithDeployment(t *testing.T) {
	type args struct {
		deployment string
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
			if got := tt.config.WithDeployment(tt.args.deployment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithMode(t *testing.T) {
	type args struct {
		mode string
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
			if got := tt.config.WithMode(tt.args.mode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_WithServiceAddress(t *testing.T) {
	type args struct {
		serviceAddress string
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
			if got := tt.config.WithServiceAddress(tt.args.serviceAddress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.WithServiceAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Build(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.Build(); !reflect.DeepEqual(got, tt.want) {
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
