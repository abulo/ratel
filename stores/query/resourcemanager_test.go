package query

import (
	"io"
	"reflect"
	"testing"
)

func TestNewResourceManager(t *testing.T) {
	tests := []struct {
		name string
		want *ResourceManager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceManager_Close(t *testing.T) {
	tests := []struct {
		name    string
		manager *ResourceManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.manager.Close(); (err != nil) != tt.wantErr {
				t.Errorf("ResourceManager.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceManager_GetResource(t *testing.T) {
	type args struct {
		key    string
		create func() (io.Closer, error)
	}
	tests := []struct {
		name    string
		manager *ResourceManager
		args    args
		want    io.Closer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.manager.GetResource(tt.args.key, tt.args.create)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceManager.GetResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceManager.GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
