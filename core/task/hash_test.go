package task

import (
	"reflect"
	"testing"
)

func TestNewHash(t *testing.T) {
	type args struct {
		replicas int
		fn       Hash
	}
	tests := []struct {
		name string
		args args
		want *Map
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHash(tt.args.replicas, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		m    *Map
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsEmpty(); got != tt.want {
				t.Errorf("Map.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Add(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		m    *Map
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Add(tt.args.keys...)
		})
	}
}

func TestMap_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *Map
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.key); got != tt.want {
				t.Errorf("Map.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
