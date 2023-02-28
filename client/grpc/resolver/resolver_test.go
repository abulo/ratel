package resolver

import (
	"reflect"
	"testing"

	"github.com/abulo/ratel/v3/registry"
	"google.golang.org/grpc/resolver"
)

func TestNewEtcdBuilder(t *testing.T) {
	type args struct {
		name     string
		registry registry.Registry
	}
	tests := []struct {
		name string
		args args
		want resolver.Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEtcdBuilder(tt.args.name, tt.args.registry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEtcdBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseBuilder_Build(t *testing.T) {
	type args struct {
		target resolver.Target
		cc     resolver.ClientConn
		opts   resolver.BuildOptions
	}
	tests := []struct {
		name    string
		b       *baseBuilder
		args    args
		want    resolver.Resolver
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Build(tt.args.target, tt.args.cc, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseBuilder_Scheme(t *testing.T) {
	tests := []struct {
		name string
		b    baseBuilder
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Scheme(); got != tt.want {
				t.Errorf("baseBuilder.Scheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseResolver_ResolveNow(t *testing.T) {
	type args struct {
		options resolver.ResolveNowOptions
	}
	tests := []struct {
		name string
		b    *baseResolver
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.ResolveNow(tt.args.options)
		})
	}
}

func Test_baseResolver_Close(t *testing.T) {
	tests := []struct {
		name string
		b    *baseResolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Close()
		})
	}
}
