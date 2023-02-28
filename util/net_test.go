package util

import (
	"reflect"
	"testing"
)

func TestGetHostName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostName()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHostName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostByName(t *testing.T) {
	type args struct {
		hostname string
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
			got, err := GetHostByName(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHostByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostByNameL(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostByNameL(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostByNameL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHostByNameL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostByAddr(t *testing.T) {
	type args struct {
		ipAddress string
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
			got, err := GetHostByAddr(tt.args.ipAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHostByAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIP2long(t *testing.T) {
	type args struct {
		ipAddress string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP2long(tt.args.ipAddress); got != tt.want {
				t.Errorf("IP2long() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLong2ip(t *testing.T) {
	type args struct {
		properAddress uint32
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
			if got := Long2ip(tt.args.properAddress); got != tt.want {
				t.Errorf("Long2ip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractIP(t *testing.T) {
	type args struct {
		addr string
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
			got, err := ExtractIP(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPrivateIP(t *testing.T) {
	type args struct {
		ipAddr string
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
			if got := isPrivateIP(tt.args.ipAddr); got != tt.want {
				t.Errorf("isPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
