package session

import (
	"context"
	"reflect"
	"testing"
)

func TestSession_Put(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   string
		value any
	}
	tests := []struct {
		name    string
		session *Session
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.session.Put(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Session.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		session *Session
		args    args
		want    any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.Get(tt.args.ctx, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		session *Session
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.session.Remove(tt.args.ctx, tt.args.key)
		})
	}
}

func Test_setSliceMap(t *testing.T) {
	type args struct {
		m     map[string]any
		keys  []string
		value any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setSliceMap(tt.args.m, tt.args.keys, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setSliceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSliceMap(t *testing.T) {
	type args struct {
		m    map[string]any
		keys []string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSliceMap(tt.args.m, tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSliceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delSliceMap(t *testing.T) {
	type args struct {
		m    map[string]any
		keys []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delSliceMap(tt.args.m, tt.args.keys)
		})
	}
}

func TestSession_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		session *Session
		args    args
		want    int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.Destroy(tt.args.ctx); got != tt.want {
				t.Errorf("Session.Destroy() = %v, want %v", got, tt.want)
			}
		})
	}
}
