package elasticsearch

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestESTraceServerInterceptor(t *testing.T) {
	type args struct {
		DisableMetric bool
		DisableTrace  bool
		Addr          string
	}
	tests := []struct {
		name string
		args args
		want *http.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ESTraceServerInterceptor(tt.args.DisableMetric, tt.args.DisableTrace, tt.args.Addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ESTraceServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestESTracedTransport_RoundTrip(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name     string
		tr       *ESTracedTransport
		args     args
		wantResp *http.Response
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.tr.RoundTrip(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ESTracedTransport.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ESTracedTransport.RoundTrip() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
