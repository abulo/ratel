package nlpword

import (
	"reflect"
	"testing"
)

func Test_ac_fail(t *testing.T) {
	type args struct {
		node *Node
		c    rune
	}
	tests := []struct {
		name string
		ac   *ac
		args args
		want *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ac.fail(tt.args.node, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ac.fail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ac_next(t *testing.T) {
	type args struct {
		node *Node
		c    rune
	}
	tests := []struct {
		name string
		ac   *ac
		args args
		want *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ac.next(tt.args.node, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ac.next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ac_output(t *testing.T) {
	type args struct {
		node     *Node
		runes    []rune
		position int
	}
	tests := []struct {
		name string
		ac   *ac
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ac.output(tt.args.node, tt.args.runes, tt.args.position)
		})
	}
}

func Test_ac_firstOutput(t *testing.T) {
	type args struct {
		node     *Node
		runes    []rune
		position int
	}
	tests := []struct {
		name string
		ac   *ac
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ac.firstOutput(tt.args.node, tt.args.runes, tt.args.position); got != tt.want {
				t.Errorf("ac.firstOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ac_replace(t *testing.T) {
	type args struct {
		node     *Node
		runes    []rune
		position int
		replace  rune
	}
	tests := []struct {
		name string
		ac   *ac
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ac.replace(tt.args.node, tt.args.runes, tt.args.position, tt.args.replace)
		})
	}
}
