package base

import "testing"

func Test_ratelHome(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ratelHome(); got != tt.want {
				t.Errorf("ratelHome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ratelHomeWithDir(t *testing.T) {
	type args struct {
		dir string
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
			if got := ratelHomeWithDir(tt.args.dir); got != tt.want {
				t.Errorf("ratelHomeWithDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	type args struct {
		src      string
		dst      string
		replaces []string
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
			if err := copyFile(tt.args.src, tt.args.dst, tt.args.replaces); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyDir(t *testing.T) {
	type args struct {
		src      string
		dst      string
		replaces []string
		ignores  []string
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
			if err := copyDir(tt.args.src, tt.args.dst, tt.args.replaces, tt.args.ignores); (err != nil) != tt.wantErr {
				t.Errorf("copyDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hasSets(t *testing.T) {
	type args struct {
		name string
		sets []string
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
			if got := hasSets(tt.args.name, tt.args.sets); got != tt.want {
				t.Errorf("hasSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree(t *testing.T) {
	type args struct {
		path string
		dir  string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Tree(tt.args.path, tt.args.dir)
		})
	}
}
