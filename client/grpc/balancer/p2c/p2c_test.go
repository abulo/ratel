package p2c

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func Test_newBuilder(t *testing.T) {
	tests := []struct {
		name string
		want balancer.Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_p2cPickerBuilder_Build(t *testing.T) {
	type args struct {
		info base.PickerBuildInfo
	}
	tests := []struct {
		name string
		p    *p2cPickerBuilder
		args args
		want balancer.Picker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Build(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("p2cPickerBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_p2cPicker_Pick(t *testing.T) {
	type args struct {
		opts balancer.PickInfo
	}
	tests := []struct {
		name    string
		p       *p2cPicker
		args    args
		want    balancer.PickResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Pick(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("p2cPicker.Pick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("p2cPicker.Pick() = %v, want %v", got, tt.want)
			}
		})
	}
}
