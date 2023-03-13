package sql

import (
	"reflect"
	"testing"
)

func TestRow_ToAny(t *testing.T) {
	tests := []struct {
		name       string
		r          *Row
		wantResult map[string]any
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.ToAny()
			if (err != nil) != tt.wantErr {
				t.Errorf("Row.ToAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Row.ToAny() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRow_ToMap(t *testing.T) {
	tests := []struct {
		name       string
		r          *Row
		wantResult map[string]string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.ToMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Row.ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Row.ToMap() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRow_ToStruct(t *testing.T) {
	type args struct {
		st any
	}
	tests := []struct {
		name    string
		r       *Row
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.ToStruct(tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("Row.ToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRows_ToAny(t *testing.T) {
	tests := []struct {
		name     string
		r        *Rows
		wantData []map[string]any
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.r.ToAny()
			if (err != nil) != tt.wantErr {
				t.Errorf("Rows.ToAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Rows.ToAny() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestRows_ToMap(t *testing.T) {
	tests := []struct {
		name     string
		r        *Rows
		wantData []map[string]string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.r.ToMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Rows.ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Rows.ToMap() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestRows_ToStruct(t *testing.T) {
	type args struct {
		st any
	}
	tests := []struct {
		name    string
		r       *Rows
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.ToStruct(tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("Rows.ToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
