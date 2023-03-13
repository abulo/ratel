package sql

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"testing"
)

func TestNewSqlConn(t *testing.T) {
	type args struct {
		driverName string
		dns        string
		pool       *pool
		opts       []SqlOption
	}
	tests := []struct {
		name string
		args args
		want SqlConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSqlConn(tt.args.driverName, tt.args.dns, tt.args.pool, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSqlConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_getCtx(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		r    *commonSqlConn
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getCtx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commonSqlConn.getCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Close(t *testing.T) {
	tests := []struct {
		name    string
		db      *commonSqlConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Close(); (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_commonSqlConn_MultiInsert(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.MultiInsert(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.MultiInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Replace(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Replace(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_InsertUpdate(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.InsertUpdate(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.InsertUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Insert(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Insert(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Update(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Update(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Delete(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Delete(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Exec(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Exec(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("commonSqlConn.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_QueryRow(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name string
		db   *commonSqlConn
		args args
		want *Row
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.QueryRow(tt.args.ctx, tt.args.querySql, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commonSqlConn.QueryRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_QueryRows(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name string
		db   *commonSqlConn
		args args
		want *Rows
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.QueryRows(tt.args.ctx, tt.args.querySql, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commonSqlConn.QueryRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_ExecCtx(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name       string
		db         *commonSqlConn
		args       args
		wantResult sql.Result
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.db.ExecCtx(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.ExecCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("commonSqlConn.ExecCtx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_commonSqlConn_QueryCtx(t *testing.T) {
	type args struct {
		ctx      context.Context
		querySql string
		args     []any
	}
	tests := []struct {
		name       string
		db         *commonSqlConn
		args       args
		wantResult *sql.Rows
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.db.QueryCtx(tt.args.ctx, tt.args.querySql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.QueryCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("commonSqlConn.QueryCtx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_commonSqlConn_acceptable(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		db   *commonSqlConn
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.acceptable(tt.args.err); got != tt.want {
				t.Errorf("commonSqlConn.acceptable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commonSqlConn_Transact(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(context.Context, Session) error
	}
	tests := []struct {
		name    string
		db      *commonSqlConn
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Transact(tt.args.ctx, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("commonSqlConn.Transact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		query string
		args  []any
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
			got, err := Format(tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeValue(t *testing.T) {
	type args struct {
		buf *strings.Builder
		arg any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeValue(tt.args.buf, tt.args.arg)
		})
	}
}

func Test_escape(t *testing.T) {
	type args struct {
		input string
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
			if got := escape(tt.args.input); got != tt.want {
				t.Errorf("escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepr(t *testing.T) {
	type args struct {
		v any
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
			if got := Repr(tt.args.v); got != tt.want {
				t.Errorf("Repr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reprOfValue(t *testing.T) {
	type args struct {
		val reflect.Value
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
			if got := reprOfValue(tt.args.val); got != tt.want {
				t.Errorf("reprOfValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
