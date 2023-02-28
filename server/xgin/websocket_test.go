package xgin

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWebSocket_Upgrade(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		ws   *WebSocket
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ws.Upgrade(tt.args.w, tt.args.r)
		})
	}
}

func TestWebSocketOptions(t *testing.T) {
	type args struct {
		pattern string
		handler WebSocketFunc
		opts    []WebSocketOption
	}
	tests := []struct {
		name string
		args args
		want *WebSocket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WebSocketOptions(tt.args.pattern, tt.args.handler, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebSocketOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
