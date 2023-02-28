package util

import "testing"

func TestExec(t *testing.T) {
	type args struct {
		command   string
		output    *[]string
		returnVar *int
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
			if got := Exec(tt.args.command, tt.args.output, tt.args.returnVar); got != tt.want {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem(t *testing.T) {
	type args struct {
		command   string
		returnVar *int
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
			if got := System(tt.args.command, tt.args.returnVar); got != tt.want {
				t.Errorf("System() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassThru(t *testing.T) {
	type args struct {
		command   string
		returnVar *int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PassThru(tt.args.command, tt.args.returnVar)
		})
	}
}
