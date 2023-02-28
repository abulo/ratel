package xmonitor

import (
	"net/http"
	"testing"
)

func TestServer_InitHandle(t *testing.T) {
	tests := []struct {
		name string
		s    *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.InitHandle()
		})
	}
}

func TestServer_HandleFunc(t *testing.T) {
	type args struct {
		pattern string
		handler http.HandlerFunc
	}
	tests := []struct {
		name string
		s    *Server
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.HandleFunc(tt.args.pattern, tt.args.handler)
		})
	}
}
