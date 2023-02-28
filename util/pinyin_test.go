package util

import "testing"

func TestZhCharToFirstPinyin(t *testing.T) {
	type args struct {
		p string
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
			if got := ZhCharToFirstPinyin(tt.args.p); got != tt.want {
				t.Errorf("ZhCharToFirstPinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZhCharToPinyin(t *testing.T) {
	type args struct {
		p string
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
			if got := ZhCharToPinyin(tt.args.p); got != tt.want {
				t.Errorf("ZhCharToPinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}
