package query

import "testing"

func TestBatchError_Add(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		be   *BatchError
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.be.Add(tt.args.err)
		})
	}
}

func TestBatchError_Err(t *testing.T) {
	tests := []struct {
		name    string
		be      *BatchError
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.be.Err(); (err != nil) != tt.wantErr {
				t.Errorf("BatchError.Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBatchError_NotNil(t *testing.T) {
	tests := []struct {
		name string
		be   *BatchError
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.be.NotNil(); got != tt.want {
				t.Errorf("BatchError.NotNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorArray_Error(t *testing.T) {
	tests := []struct {
		name string
		ea   errorArray
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ea.Error(); got != tt.want {
				t.Errorf("errorArray.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
