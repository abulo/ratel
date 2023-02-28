package xgin

import (
	"reflect"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	type args struct {
		code int
		msg  []string
	}
	tests := []struct {
		name string
		args args
		want *HTTPError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPError(tt.args.code, tt.args.msg...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    HTTPError
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("HTTPError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
