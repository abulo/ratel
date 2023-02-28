package xgin

import (
	"reflect"
	"testing"

	rstatus "google.golang.org/genproto/googleapis/rpc/status"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
)

func TestEmptyMessage_Reset(t *testing.T) {
	tests := []struct {
		name string
		m    *EmptyMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Reset()
		})
	}
}

func TestEmptyMessage_String(t *testing.T) {
	tests := []struct {
		name string
		m    *EmptyMessage
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("EmptyMessage.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyMessage_ProtoMessage(t *testing.T) {
	tests := []struct {
		name string
		e    *EmptyMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.ProtoMessage()
		})
	}
}

func TestGRPCProxyMessage_Reset(t *testing.T) {
	tests := []struct {
		name string
		m    *GRPCProxyMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Reset()
		})
	}
}

func TestGRPCProxyMessage_String(t *testing.T) {
	tests := []struct {
		name string
		m    *GRPCProxyMessage
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("GRPCProxyMessage.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCProxyMessage_ProtoMessage(t *testing.T) {
	tests := []struct {
		name string
		g    *GRPCProxyMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.ProtoMessage()
		})
	}
}

func TestGRPCProxyMessage_MarshalJSONPB(t *testing.T) {
	type args struct {
		jsb *jsonpb.MarshalOptions
	}
	tests := []struct {
		name    string
		m       *GRPCProxyMessage
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.MarshalJSONPB(tt.args.jsb)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCProxyMessage.MarshalJSONPB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCProxyMessage.MarshalJSONPB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusErr_Proto(t *testing.T) {
	tests := []struct {
		name string
		e    *statusErr
		want *rstatus.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Proto(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statusErr.Proto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_statusFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  *statusErr
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := statusFromString(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statusFromString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("statusFromString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_createStatusErr(t *testing.T) {
	type args struct {
		code uint32
		msg  string
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
			if got := createStatusErr(tt.args.code, tt.args.msg); got != tt.want {
				t.Errorf("createStatusErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
