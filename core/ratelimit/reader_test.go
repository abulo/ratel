package ratelimit

import (
	"io"
	"reflect"
	"testing"
)

func TestReader(t *testing.T) {
	type args struct {
		r      io.Reader
		bucket *Bucket
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reader(tt.args.r, tt.args.bucket); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_Read(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		r       *reader
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Read(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("reader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writer_Write(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		w       *writer
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Write(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("writer.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
