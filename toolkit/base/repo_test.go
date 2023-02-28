package base

import (
	"context"
	"reflect"
	"testing"
)

func Test_repoDir(t *testing.T) {
	type args struct {
		url string
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
			if got := repoDir(tt.args.url); got != tt.want {
				t.Errorf("repoDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRepo(t *testing.T) {
	type args struct {
		url    string
		branch string
	}
	tests := []struct {
		name string
		args args
		want *Repo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepo(tt.args.url, tt.args.branch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Path(t *testing.T) {
	tests := []struct {
		name string
		r    *Repo
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Path(); got != tt.want {
				t.Errorf("Repo.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Pull(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Pull(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Repo.Pull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_Clone(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Clone(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Repo.Clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_CopyTo(t *testing.T) {
	type args struct {
		ctx     context.Context
		to      string
		modPath string
		ignores []string
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CopyTo(tt.args.ctx, tt.args.to, tt.args.modPath, tt.args.ignores); (err != nil) != tt.wantErr {
				t.Errorf("Repo.CopyTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_CopyToV2(t *testing.T) {
	type args struct {
		ctx      context.Context
		to       string
		modPath  string
		ignores  []string
		replaces []string
	}
	tests := []struct {
		name    string
		r       *Repo
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CopyToV2(tt.args.ctx, tt.args.to, tt.args.modPath, tt.args.ignores, tt.args.replaces); (err != nil) != tt.wantErr {
				t.Errorf("Repo.CopyToV2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
