package sql

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func Test_transact(t *testing.T) {
	type args struct {
		ctx  context.Context
		db   *commonSqlConn
		pool *pool
		b    beginnable
		fn   func(context.Context, Session) error
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
			if err := transact(tt.args.ctx, tt.args.db, tt.args.pool, tt.args.b, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("transact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_transactOnConn(t *testing.T) {
	type args struct {
		ctx  context.Context
		conn *sql.DB
		pool *pool
		b    beginnable
		fn   func(context.Context, Session) error
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
			if err := transactOnConn(tt.args.ctx, tt.args.conn, tt.args.pool, tt.args.b, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("transactOnConn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_txSession_MultiInsert(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.MultiInsert(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.MultiInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Replace(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Replace(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_InsertUpdate(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.InsertUpdate(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.InsertUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Insert(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Insert(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Update(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Update(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Delete(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Delete(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Exec(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Exec(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_QueryRow(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name string
		db   txSession
		args args
		want *Row
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.QueryRow(tt.args.ctx, tt.args.query, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("txSession.QueryRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_Count(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name    string
		db      txSession
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Count(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("txSession.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_QueryRows(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name string
		db   txSession
		args args
		want *Rows
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.db.QueryRows(tt.args.ctx, tt.args.query, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("txSession.QueryRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_txSession_ExecCtx(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name       string
		db         txSession
		args       args
		wantResult sql.Result
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.db.ExecCtx(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.ExecCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("txSession.ExecCtx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_txSession_QueryCtx(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
		args  []any
	}
	tests := []struct {
		name       string
		db         txSession
		args       args
		wantResult *sql.Rows
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.db.QueryCtx(tt.args.ctx, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("txSession.QueryCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("txSession.QueryCtx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
