package registry

import (
	"reflect"
	"testing"
)

func Test_newEndpoints(t *testing.T) {
	tests := []struct {
		name string
		want *Endpoints
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEndpoints(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoints_DeepCopy(t *testing.T) {
	tests := []struct {
		name string
		in   *Endpoints
		want *Endpoints
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.DeepCopy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoints.DeepCopy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoints_DeepCopyInfo(t *testing.T) {
	type args struct {
		out *Endpoints
	}
	tests := []struct {
		name string
		in   *Endpoints
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.DeepCopyInfo(tt.args.out)
		})
	}
}

func TestRouteConfig_String(t *testing.T) {
	tests := []struct {
		name   string
		config RouteConfig
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.String(); got != tt.want {
				t.Errorf("RouteConfig.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
