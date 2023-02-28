package util

import (
	"reflect"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomBytes(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandomBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandomBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomInt(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomInt(tt.args.min, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandomInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RandomInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrPad(t *testing.T) {
	type args struct {
		str1 string
		str2 string
		i    int
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
			if got := StrPad(tt.args.str1, tt.args.str2, tt.args.i); got != tt.want {
				t.Errorf("StrPad() = %v, want %v", got, tt.want)
			}
		})
	}
}
