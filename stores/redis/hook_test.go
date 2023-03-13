package redis

import (
	"context"
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestOpenTraceHook_DialHook(t *testing.T) {
	type args struct {
		hook redis.DialHook
	}
	tests := []struct {
		name string
		op   OpenTraceHook
		args args
		want redis.DialHook
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.DialHook(tt.args.hook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenTraceHook.DialHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTraceHook_ProcessHook(t *testing.T) {
	type args struct {
		hook redis.ProcessHook
	}
	tests := []struct {
		name string
		op   OpenTraceHook
		args args
		want redis.ProcessHook
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.ProcessHook(tt.args.hook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenTraceHook.ProcessHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTraceHook_ProcessPipelineHook(t *testing.T) {
	type args struct {
		hook redis.ProcessPipelineHook
	}
	tests := []struct {
		name string
		op   OpenTraceHook
		args args
		want redis.ProcessPipelineHook
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.ProcessPipelineHook(tt.args.hook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenTraceHook.ProcessPipelineHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTraceHook_BeforeProcess(t *testing.T) {
	type args struct {
		ctx context.Context
		cmd redis.Cmder
	}
	tests := []struct {
		name    string
		op      OpenTraceHook
		args    args
		want    context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.op.BeforeProcess(tt.args.ctx, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenTraceHook.BeforeProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenTraceHook.BeforeProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTraceHook_AfterProcess(t *testing.T) {
	type args struct {
		ctx context.Context
		cmd redis.Cmder
	}
	tests := []struct {
		name    string
		op      OpenTraceHook
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.op.AfterProcess(tt.args.ctx, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("OpenTraceHook.AfterProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenTraceHook_BeforeProcessPipeline(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmds []redis.Cmder
	}
	tests := []struct {
		name    string
		op      OpenTraceHook
		args    args
		want    context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.op.BeforeProcessPipeline(tt.args.ctx, tt.args.cmds)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenTraceHook.BeforeProcessPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenTraceHook.BeforeProcessPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTraceHook_AfterProcessPipeline(t *testing.T) {
	type args struct {
		ctx  context.Context
		cmds []redis.Cmder
	}
	tests := []struct {
		name    string
		op      OpenTraceHook
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.op.AfterProcessPipeline(tt.args.ctx, tt.args.cmds); (err != nil) != tt.wantErr {
				t.Errorf("OpenTraceHook.AfterProcessPipeline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		b []byte
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
			if got := String(tt.args.b); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendCmd(t *testing.T) {
	type args struct {
		b   []byte
		cmd redis.Cmder
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendCmd(tt.args.b, tt.args.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendArg(t *testing.T) {
	type args struct {
		b []byte
		v any
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendArg(tt.args.b, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendUTF8String(t *testing.T) {
	type args struct {
		b []byte
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendUTF8String(tt.args.b, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendUTF8String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendRune(t *testing.T) {
	type args struct {
		b []byte
		r rune
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendRune(tt.args.b, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appendRune() = %v, want %v", got, tt.want)
			}
		})
	}
}
