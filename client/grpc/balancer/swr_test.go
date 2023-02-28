package balancer

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/balancer"
)

func Test_swrPickerBuilder_Build(t *testing.T) {
	type args struct {
		info PickerBuildInfo
	}
	tests := []struct {
		name string
		s    swrPickerBuilder
		args args
		want balancer.Picker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Build(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("swrPickerBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSWRPicker(t *testing.T) {
	type args struct {
		info PickerBuildInfo
	}
	tests := []struct {
		name string
		args args
		want *swrPicker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSWRPicker(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSWRPicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swrPicker_Pick(t *testing.T) {
	type args struct {
		info balancer.PickInfo
	}
	tests := []struct {
		name    string
		p       *swrPicker
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
				t.Errorf("swrPicker.Pick() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("swrPicker.Pick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_swrPicker_parseBuildInfo(t *testing.T) {
	type args struct {
		info PickerBuildInfo
	}
	tests := []struct {
		name string
		p    *swrPicker
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.parseBuildInfo(tt.args.info)
		})
	}
}
