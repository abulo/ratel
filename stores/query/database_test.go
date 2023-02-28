package query

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestQuery_NewBuilder(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		query *Query
		args  args
		want  *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.query.NewBuilder(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.NewBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Begin(t *testing.T) {
	tests := []struct {
		name    string
		query   *Query
		want    *Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.Begin()
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.Begin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Begin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Exec(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name    string
		query   *Query
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.Exec(tt.args.ctx, tt.args.querySQL, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Query(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name    string
		query   *Query
		args    args
		want    *sql.Rows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.Query(tt.args.ctx, tt.args.querySQL, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_Commit(t *testing.T) {
	tests := []struct {
		name    string
		query   *Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.query.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_Rollback(t *testing.T) {
	tests := []struct {
		name    string
		query   *Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.query.Rollback(); (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_NewBuilder(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		query *Transaction
		args  args
		want  *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.query.NewBuilder(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transaction.NewBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_Exec(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name    string
		query   *Transaction
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.Exec(tt.args.ctx, tt.args.querySQL, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transaction.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_Query(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name    string
		query   *Transaction
		args    args
		want    *sql.Rows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.Query(tt.args.ctx, tt.args.querySQL, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transaction.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_SQLRaw(t *testing.T) {
	tests := []struct {
		name  string
		query *Transaction
		want  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.query.SQLRaw(); got != tt.want {
				t.Errorf("Transaction.SQLRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_SQLRaw(t *testing.T) {
	tests := []struct {
		name  string
		query *Query
		want  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.query.SQLRaw(); got != tt.want {
				t.Errorf("Query.SQLRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_LastSQL(t *testing.T) {
	type args struct {
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name  string
		query *Transaction
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.LastSQL(tt.args.querySQL, tt.args.args...)
		})
	}
}

func TestQuery_LastSQL(t *testing.T) {
	type args struct {
		querySQL string
		args     []interface{}
	}
	tests := []struct {
		name  string
		query *Query
		args  args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.LastSQL(tt.args.querySQL, tt.args.args...)
		})
	}
}

func TestSQL_ToString(t *testing.T) {
	tests := []struct {
		name   string
		SQLRaw SQL
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.SQLRaw.ToString(); got != tt.want {
				t.Errorf("SQL.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNilFixed(t *testing.T) {
	type args struct {
		i interface{}
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
			if got := isNilFixed(tt.args.i); got != tt.want {
				t.Errorf("isNilFixed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert(t *testing.T) {
	type args struct {
		s string
		v interface{}
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
			if got := convert(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQL_ToJSON(t *testing.T) {
	tests := []struct {
		name string
		sql  SQL
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sql.ToJSON(); got != tt.want {
				t.Errorf("SQL.ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
