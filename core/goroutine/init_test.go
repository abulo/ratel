package goroutine

import "testing"

func Test_try(t *testing.T) {
	type args struct {
		fn      func() error
		cleaner func()
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
			if err := try(tt.args.fn, tt.args.cleaner); (err != nil) != tt.wantErr {
				t.Errorf("try() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_try2(t *testing.T) {
	type args struct {
		fn      func()
		cleaner func()
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
			if err := try2(tt.args.fn, tt.args.cleaner); (err != nil) != tt.wantErr {
				t.Errorf("try2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
