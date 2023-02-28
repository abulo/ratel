//go:build go1.19
// +build go1.19

package base

import "testing"

func TestGoInstall(t *testing.T) {
	type args struct {
		path []string
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
			if err := GoInstall(tt.args.path...); (err != nil) != tt.wantErr {
				t.Errorf("GoInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
