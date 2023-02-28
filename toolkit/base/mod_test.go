package base

import "testing"

func TestModulePath(t *testing.T) {
	type args struct {
		filename string
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
			got, err := ModulePath(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ModulePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleVersion(t *testing.T) {
	type args struct {
		path string
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
			got, err := ModuleVersion(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModuleVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ModuleVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRatelMod(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RatelMod(); got != tt.want {
				t.Errorf("RatelMod() = %v, want %v", got, tt.want)
			}
		})
	}
}
