package util

import (
	"archive/zip"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestEcho(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Echo(tt.args.args...)
		})
	}
}

func TestUniqid(t *testing.T) {
	type args struct {
		prefix string
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
			if got := Uniqid(tt.args.prefix); got != tt.want {
				t.Errorf("Uniqid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExit(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Exit(tt.args.status)
		})
	}
}

func TestDie(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Die(tt.args.status)
		})
	}
}

func TestGetenv(t *testing.T) {
	type args struct {
		varname string
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
			if got := Getenv(tt.args.varname); got != tt.want {
				t.Errorf("Getenv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutenv(t *testing.T) {
	type args struct {
		setting string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Putenv(tt.args.setting); (err != nil) != tt.wantErr {
				t.Errorf("Putenv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryGetUsage(t *testing.T) {
	type args struct {
		realUsage bool
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MemoryGetUsage(tt.args.realUsage); got != tt.want {
				t.Errorf("MemoryGetUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryGetPeakUsage(t *testing.T) {
	type args struct {
		realUsage bool
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MemoryGetPeakUsage(tt.args.realUsage); got != tt.want {
				t.Errorf("MemoryGetPeakUsage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	type args struct {
		version1 string
		version2 string
		operator string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VersionCompare(tt.args.version1, tt.args.version2, tt.args.operator); got != tt.want {
				t.Errorf("VersionCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipOpen(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *zip.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ZipOpen(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZipOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPack(t *testing.T) {
	type args struct {
		order binary.ByteOrder
		data  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pack(tt.args.order, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	type args struct {
		order binary.ByteOrder
		data  string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.args.order, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernary(t *testing.T) {
	type args struct {
		condition bool
		trueVal   interface{}
		falseVal  interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.condition, tt.args.trueVal, tt.args.falseVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}
