package query

import "testing"

func TestGrammar_compileSelect(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileSelect(); got != tt.want {
				t.Errorf("Grammar.compileSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileTable(t *testing.T) {
	type args struct {
		from bool
	}
	tests := []struct {
		name string
		g    Grammar
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileTable(tt.args.from); got != tt.want {
				t.Errorf("Grammar.compileTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileOrder(t *testing.T) {
	type args struct {
		isUnion bool
	}
	tests := []struct {
		name string
		g    Grammar
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileOrder(tt.args.isUnion); got != tt.want {
				t.Errorf("Grammar.compileOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileGroup(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileGroup(); got != tt.want {
				t.Errorf("Grammar.compileGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileLimit(t *testing.T) {
	type args struct {
		isUnion bool
	}
	tests := []struct {
		name string
		g    Grammar
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileLimit(tt.args.isUnion); got != tt.want {
				t.Errorf("Grammar.compileLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileDistinct(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileDistinct(); got != tt.want {
				t.Errorf("Grammar.compileDistinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileWhere(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileWhere(); got != tt.want {
				t.Errorf("Grammar.compileWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileJoin(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileJoin(); got != tt.want {
				t.Errorf("Grammar.compileJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileUnion(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileUnion(); got != tt.want {
				t.Errorf("Grammar.compileUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_Select(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Select(); got != tt.want {
				t.Errorf("Grammar.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_Insert(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Insert(); got != tt.want {
				t.Errorf("Grammar.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_Replace(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Replace(); got != tt.want {
				t.Errorf("Grammar.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileInsertValue(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileInsertValue(); got != tt.want {
				t.Errorf("Grammar.compileInsertValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_Delete(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Delete(); got != tt.want {
				t.Errorf("Grammar.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_compileUpdateValue(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.compileUpdateValue(); got != tt.want {
				t.Errorf("Grammar.compileUpdateValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_Update(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Update(); got != tt.want {
				t.Errorf("Grammar.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_InsertUpdate(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.InsertUpdate(); got != tt.want {
				t.Errorf("Grammar.InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrammar_ToSQL(t *testing.T) {
	tests := []struct {
		name string
		g    Grammar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.ToSQL(); got != tt.want {
				t.Errorf("Grammar.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
