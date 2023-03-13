package sql

import (
	"reflect"
	"testing"
)

func TestNewEpr(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want Epr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEpr(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpr_ToString(t *testing.T) {
	tests := []struct {
		name string
		e    Epr
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ToString(); got != tt.want {
				t.Errorf("Epr.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
