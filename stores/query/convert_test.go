package query

import (
	"reflect"
	"testing"
)

func Test_toString(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name    string
		args    args
		wantDst string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDst, err := toString(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("toString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDst != tt.wantDst {
				t.Errorf("toString() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

func Test_extractTagInfo(t *testing.T) {
	type args struct {
		st reflect.Value
	}
	tests := []struct {
		name        string
		args        args
		wantTagList map[string]reflect.Value
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTagList, err := extractTagInfo(tt.args.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTagInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTagList, tt.wantTagList) {
				t.Errorf("extractTagInfo() = %v, want %v", gotTagList, tt.wantTagList)
			}
		})
	}
}
