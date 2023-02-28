package nlpword

import (
	"reflect"
	"testing"
)

func TestTrie_BuildFailureLinks(t *testing.T) {
	tests := []struct {
		name string
		tree *Trie
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tree.BuildFailureLinks()
		})
	}
}

func TestTrie_bfs(t *testing.T) {
	tests := []struct {
		name string
		tree *Trie
		want <-chan *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.bfs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie.bfs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTrie(t *testing.T) {
	tests := []struct {
		name string
		want *Trie
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrie(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_Add(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		tree *Trie
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tree.Add(tt.args.words...)
		})
	}
}

func TestTrie_add(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		tree *Trie
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tree.add(tt.args.word)
		})
	}
}

func TestTrie_Replace(t *testing.T) {
	type args struct {
		text      string
		character rune
	}
	tests := []struct {
		name string
		tree *Trie
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.Replace(tt.args.text, tt.args.character); got != tt.want {
				t.Errorf("Trie.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_Filter(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		tree *Trie
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.Filter(tt.args.text); got != tt.want {
				t.Errorf("Trie.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrie_Validate(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name  string
		tree  *Trie
		args  args
		want  bool
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tree.Validate(tt.args.text)
			if got != tt.want {
				t.Errorf("Trie.Validate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Trie.Validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTrie_FindIn(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name  string
		tree  *Trie
		args  args
		want  bool
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tree.FindIn(tt.args.text)
			if got != tt.want {
				t.Errorf("Trie.FindIn() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Trie.FindIn() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTrie_FindAll(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		tree *Trie
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tree.FindAll(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trie.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	type args struct {
		character rune
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.character); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRootNode(t *testing.T) {
	type args struct {
		character rune
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRootNode(tt.args.character); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRootNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_IsLeafNode(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.IsLeafNode(); got != tt.want {
				t.Errorf("Node.IsLeafNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_IsRootNode(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.IsRootNode(); got != tt.want {
				t.Errorf("Node.IsRootNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_IsPathEnd(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.IsPathEnd(); got != tt.want {
				t.Errorf("Node.IsPathEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
