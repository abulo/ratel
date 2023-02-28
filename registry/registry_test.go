package registry

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/server"
)

func TestKind_String(t *testing.T) {
	tests := []struct {
		name string
		kind Kind
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("Kind.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToKind(t *testing.T) {
	type args struct {
		kindStr string
	}
	tests := []struct {
		name string
		args args
		want Kind
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToKind(tt.args.kindStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceKey(t *testing.T) {
	type args struct {
		prefix string
		s      *server.ServiceInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceKey(tt.args.prefix, tt.args.s); got != tt.want {
				t.Errorf("GetServiceKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceValue(t *testing.T) {
	type args struct {
		s *server.ServiceInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceValue(tt.args.s); got != tt.want {
				t.Errorf("GetServiceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *server.ServiceInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetService(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService() = %v, want %v", got, tt.want)
			}
		})
	}
}
