package query

import (
	"reflect"
	"testing"
)

func TestRow_ToArray(t *testing.T) {
	tests := []struct {
		name       string
		r          *Row
		wantResult []string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.ToArray()
			if (err != nil) != tt.wantErr {
				t.Errorf("Row.ToArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Row.ToArray() = %v, want %v", gotResult, tt.wantResult)
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

func TestRow_ToInterface(t *testing.T) {
	tests := []struct {
		name       string
		r          *Row
		wantResult map[string]interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.r.ToInterface()
			if (err != nil) != tt.wantErr {
				t.Errorf("Row.ToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Row.ToInterface() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRow_ToStruct(t *testing.T) {
	type args struct {
		st interface{}
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

func TestRows_ToArray(t *testing.T) {
	tests := []struct {
		name     string
		r        *Rows
		wantData [][]string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.r.ToArray()
			if (err != nil) != tt.wantErr {
				t.Errorf("Rows.ToArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Rows.ToArray() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestRows_ToInterface(t *testing.T) {
	tests := []struct {
		name     string
		r        *Rows
		wantData []map[string]interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.r.ToInterface()
			if (err != nil) != tt.wantErr {
				t.Errorf("Rows.ToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Rows.ToInterface() = %v, want %v", gotData, tt.wantData)
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
		st interface{}
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
