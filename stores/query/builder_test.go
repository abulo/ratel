package query

import (
	"reflect"
	"testing"
)

func TestBuilder_Table(t *testing.T) {
	type args struct {
		tableName []string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Table(tt.args.tableName...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Select(t *testing.T) {
	type args struct {
		columns []string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Select(tt.args.columns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Where(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Where(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrWhere(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrWhere(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Equal(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Equal(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrEqual(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrEqual(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_NotEqual(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.NotEqual(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.NotEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Greater(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Greater(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Greater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_GreaterEqual(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.GreaterEqual(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.GreaterEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Less(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Less(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_LessEqual(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.LessEqual(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.LessEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrNotEqual(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrNotEqual(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrNotEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Between(t *testing.T) {
	type args struct {
		column string
		value1 any
		value2 any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Between(tt.args.column, tt.args.value1, tt.args.value2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Between() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrBetween(t *testing.T) {
	type args struct {
		column string
		value1 any
		value2 any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrBetween(tt.args.column, tt.args.value1, tt.args.value2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_NotBetween(t *testing.T) {
	type args struct {
		column string
		value1 any
		value2 any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.NotBetween(tt.args.column, tt.args.value1, tt.args.value2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.NotBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_NotOrBetween(t *testing.T) {
	type args struct {
		column string
		value1 any
		value2 any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.NotOrBetween(tt.args.column, tt.args.value1, tt.args.value2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.NotOrBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_In(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.In(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrIn(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrIn(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_NotIn(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.NotIn(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.NotIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrNotIn(t *testing.T) {
	type args struct {
		column string
		value  []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrNotIn(tt.args.column, tt.args.value...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrNotIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_IsNULL(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.IsNULL(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.IsNULL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrIsNULL(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrIsNULL(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrIsNULL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_IsNotNULL(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.IsNotNULL(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.IsNotNULL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrIsNotNULL(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrIsNotNULL(tt.args.column); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrIsNotNULL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Like(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Like(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Like() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrLike(t *testing.T) {
	type args struct {
		column string
		value  any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrLike(tt.args.column, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrLike() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Join(t *testing.T) {
	type args struct {
		tablename string
		on        string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Join(tt.args.tablename, tt.args.on); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_InnerJoin(t *testing.T) {
	type args struct {
		tablename string
		on        string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.InnerJoin(tt.args.tablename, tt.args.on); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.InnerJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_LeftJoin(t *testing.T) {
	type args struct {
		tablename string
		on        string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.LeftJoin(tt.args.tablename, tt.args.on); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.LeftJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_RightJoin(t *testing.T) {
	type args struct {
		tablename string
		on        string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.RightJoin(tt.args.tablename, tt.args.on); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.RightJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Union(t *testing.T) {
	type args struct {
		unions []Builder
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Union(tt.args.unions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_UnionOffset(t *testing.T) {
	type args struct {
		offset int64
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.UnionOffset(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.UnionOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_UnionLimit(t *testing.T) {
	type args struct {
		limit int64
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.UnionLimit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.UnionLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_UnionOrderBy(t *testing.T) {
	type args struct {
		column    string
		direction string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.UnionOrderBy(tt.args.column, tt.args.direction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.UnionOrderBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_UnionAll(t *testing.T) {
	type args struct {
		unions []Builder
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.UnionAll(tt.args.unions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.UnionAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Distinct(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Distinct(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_GroupBy(t *testing.T) {
	type args struct {
		groups []string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.GroupBy(tt.args.groups...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_OrderBy(t *testing.T) {
	type args struct {
		column    string
		direction string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.OrderBy(tt.args.column, tt.args.direction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.OrderBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Offset(t *testing.T) {
	type args struct {
		offset int64
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Offset(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Skip(t *testing.T) {
	type args struct {
		offset int64
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Skip(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Limit(t *testing.T) {
	type args struct {
		limit int64
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Limit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_ToSQL(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.ToSQL(tt.args.method); got != tt.want {
				t.Errorf("Builder.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_toWhere(t *testing.T) {
	type args struct {
		column   string
		operator string
		valueNum int64
		do       string
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.toWhere(tt.args.column, tt.args.operator, tt.args.valueNum, tt.args.do); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.toWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_addArg(t *testing.T) {
	type args struct {
		value []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.addArg(tt.args.value...)
		})
	}
}

func TestBuilder_beforeArg(t *testing.T) {
	type args struct {
		value []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.beforeArg(tt.args.value...)
		})
	}
}

func TestBuilder_setData(t *testing.T) {
	type args struct {
		data []map[string]any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.builder.setData(tt.args.data...)
		})
	}
}

func TestBuilder_getInsertMap(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name        string
		builder     *Builder
		args        args
		wantColumns []string
		wantValues  map[string][]any
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotColumns, gotValues, err := tt.builder.getInsertMap(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.getInsertMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotColumns, tt.wantColumns) {
				t.Errorf("Builder.getInsertMap() gotColumns = %v, want %v", gotColumns, tt.wantColumns)
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Builder.getInsertMap() gotValues = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func TestBuilder_IsZero(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.IsZero(tt.args.v); got != tt.want {
				t.Errorf("Builder.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_MultiInsert(t *testing.T) {
	type args struct {
		datas []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.MultiInsert(tt.args.datas...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.MultiInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_MultiInsertSQL(t *testing.T) {
	type args struct {
		datas []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.MultiInsertSQL(tt.args.datas...); got != tt.want {
				t.Errorf("Builder.MultiInsertSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Replace(t *testing.T) {
	type args struct {
		datas []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Replace(tt.args.datas...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_ReplaceSQL(t *testing.T) {
	type args struct {
		datas []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.ReplaceSQL(tt.args.datas...); got != tt.want {
				t.Errorf("Builder.ReplaceSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_InsertUpdate(t *testing.T) {
	type args struct {
		insert any
		update any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.InsertUpdate(tt.args.insert, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.InsertUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_InsertUpdateSQL(t *testing.T) {
	type args struct {
		insert any
		update any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.InsertUpdateSQL(tt.args.insert, tt.args.update); got != tt.want {
				t.Errorf("Builder.InsertUpdateSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Insert(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Insert(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_InsertSQL(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.InsertSQL(tt.args.data); got != tt.want {
				t.Errorf("Builder.InsertSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Update(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Update(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_UpdateSQL(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.UpdateSQL(tt.args.data); got != tt.want {
				t.Errorf("Builder.UpdateSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Delete(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Delete()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_DeleteSQL(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.DeleteSQL(); got != tt.want {
				t.Errorf("Builder.DeleteSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Count(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Exec(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.builder.Exec(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_ExecSQL(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.ExecSQL(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("Builder.ExecSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_QueryRows(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Rows
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.QueryRows(tt.args.sql, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.QueryRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_QueryRowsSQL(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.QueryRowsSQL(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("Builder.QueryRowsSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_QueryRowSQL(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.QueryRowSQL(tt.args.sql, tt.args.args...); got != tt.want {
				t.Errorf("Builder.QueryRowSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_QueryRow(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		builder *Builder
		args    args
		want    *Row
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.QueryRow(tt.args.sql, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.QueryRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Row(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    *Row
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Row(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Row() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_RowSQL(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.RowSQL(); got != tt.want {
				t.Errorf("Builder.RowSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_RowsSQL(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.RowsSQL(); got != tt.want {
				t.Errorf("Builder.RowsSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_Rows(t *testing.T) {
	tests := []struct {
		name    string
		builder *Builder
		want    *Rows
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Rows(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.Rows() = %v, want %v", got, tt.want)
			}
		})
	}
}
