package registry

import (
	"reflect"
	"testing"
)

func TestRegisterBuilder(t *testing.T) {
	type args struct {
		kind  string
		build Builder
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterBuilder(tt.args.kind, tt.args.build)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Config
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

func TestConfig_Lab(t *testing.T) {
	type args struct {
		name string
		lab  ConfigLab
	}
	tests := []struct {
		name   string
		config Config
		args   args
		want   Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.Lab(tt.args.name, tt.args.lab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Lab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_InitDefaultRegister(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.InitDefaultRegister()
		})
	}
}
