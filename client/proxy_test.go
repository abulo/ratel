package client

import (
	"reflect"
	"testing"
)

func TestNewClientProxy(t *testing.T) {
	tests := []struct {
		name string
		want *Proxy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientProxy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}
