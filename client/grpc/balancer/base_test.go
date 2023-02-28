package balancer

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/resolver"
)

func TestNewBalancerBuilderV2(t *testing.T) {
	type args struct {
		name   string
		pb     PickerBuilder
		config base.Config
	}
	tests := []struct {
		name string
		args args
		want balancer.Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBalancerBuilderV2(tt.args.name, tt.args.pb, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBalancerBuilderV2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseBuilder_Build(t *testing.T) {
	type args struct {
		cc  balancer.ClientConn
		opt balancer.BuildOptions
	}
	tests := []struct {
		name string
		bb   *baseBuilder
		args args
		want balancer.Balancer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb.Build(tt.args.cc, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseBuilder_Name(t *testing.T) {
	tests := []struct {
		name string
		bb   *baseBuilder
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bb.Name(); got != tt.want {
				t.Errorf("baseBuilder.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseBalancer_HandleResolvedAddrs(t *testing.T) {
	type args struct {
		addrs []resolver.Address
		err   error
	}
	tests := []struct {
		name string
		b    *baseBalancer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.HandleResolvedAddrs(tt.args.addrs, tt.args.err)
		})
	}
}

func Test_baseBalancer_ResolverError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		b    *baseBalancer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.ResolverError(tt.args.err)
		})
	}
}

func Test_baseBalancer_UpdateClientConnState(t *testing.T) {
	type args struct {
		s balancer.ClientConnState
	}
	tests := []struct {
		name    string
		b       *baseBalancer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UpdateClientConnState(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("baseBalancer.UpdateClientConnState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_baseBalancer_regeneratePicker(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		b    *baseBalancer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.regeneratePicker(tt.args.err)
		})
	}
}

func Test_baseBalancer_HandleSubConnStateChange(t *testing.T) {
	type args struct {
		sc balancer.SubConn
		s  connectivity.State
	}
	tests := []struct {
		name string
		b    *baseBalancer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.HandleSubConnStateChange(tt.args.sc, tt.args.s)
		})
	}
}

func Test_baseBalancer_UpdateSubConnState(t *testing.T) {
	type args struct {
		sc    balancer.SubConn
		state balancer.SubConnState
	}
	tests := []struct {
		name string
		b    *baseBalancer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.UpdateSubConnState(tt.args.sc, tt.args.state)
		})
	}
}

func Test_baseBalancer_Close(t *testing.T) {
	tests := []struct {
		name string
		b    *baseBalancer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Close()
		})
	}
}

func TestNewErrPickerV2(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want balancer.Picker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrPickerV2(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrPickerV2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errPickerV2_Pick(t *testing.T) {
	type args struct {
		info balancer.PickInfo
	}
	tests := []struct {
		name    string
		p       *errPickerV2
		args    args
		want    balancer.PickResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Pick(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("errPickerV2.Pick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errPickerV2.Pick() = %v, want %v", got, tt.want)
			}
		})
	}
}
