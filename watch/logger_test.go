package watch

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_loggerWriter_Write(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name    string
		w       *loggerWriter
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Write(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("loggerWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("loggerWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_consoleLogger_Output(t *testing.T) {
	type args struct {
		t   int8
		msg string
	}
	tests := []struct {
		name string
		c    *consoleLogger
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Output(tt.args.t, tt.args.msg)
		})
	}
}

func TestNewConsoleLogger(t *testing.T) {
	type args struct {
		showIgnore bool
	}
	tests := []struct {
		name    string
		args    args
		want    Logger
		wantErr string
		wantOut string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &bytes.Buffer{}
			out := &bytes.Buffer{}
			if got := NewConsoleLogger(tt.args.showIgnore, err, out); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsoleLogger() = %v, want %v", got, tt.want)
			}
			if gotErr := err.String(); gotErr != tt.wantErr {
				t.Errorf("NewConsoleLogger() = %v, want %v", gotErr, tt.wantErr)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("NewConsoleLogger() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
