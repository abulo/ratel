package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func transact(ctx context.Context, db *commonSqlConn, pool *pool, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	conn, err := db.connProv()
	if err != nil {
		return err
	}
	return transactOnConn(ctx, conn, pool, b, fn)
}

func transactOnConn(ctx context.Context, conn *sql.DB, pool *pool, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	var tx trans
	begin, err := conn.Begin()
	if err != nil {
		return err
	}
	tx = txSession{
		Tx:             begin,
		disableTrace:   pool.DisableTrace,
		disableMetric:  pool.DisableMetric,
		disablePrepare: pool.DisablePrepare,
		driverName:     pool.DriverName,
		dbName:         pool.DbName,
		addr:           pool.Addr,
	}

	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %w", p, e)
			} else {
				err = fmt.Errorf("recoveer from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %w", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	return fn(ctx, tx)
}

// getCtx returns a context with a timeout.
func (db txSession) getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}

// MultiInsert 批量插入
func (db txSession) MultiInsert(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Replace 替换
func (db txSession) Replace(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// InsertUpdate 插入或更新
func (db txSession) InsertUpdate(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Insert 插入
func (db txSession) Insert(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update 更新
func (db txSession) Update(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete 删除
func (db txSession) Delete(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Exec executes a query without returning any rows.
func (db txSession) Exec(ctx context.Context, querySql string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, querySql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryRow returns a single row from the database.
func (db txSession) QueryRow(ctx context.Context, querySql string, args ...any) *Row {
	res, err := db.QueryCtx(ctx, querySql, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{
		rows: &Rows{rows: res},
		err:  err,
	}
}

// QueryRow returns a single row from the database.
func (db txSession) QueryRows(ctx context.Context, querySql string, args ...any) *Rows {
	res, err := db.QueryCtx(ctx, querySql, args...)
	if err != nil {
		return &Rows{err: err}
	}
	return &Rows{
		rows: res,
		err:  err,
	}
}

func (db txSession) ExecCtx(ctx context.Context, querySql string, args ...any) (result sql.Result, err error) {
	ctx = db.getCtx(ctx)
	start := time.Now()
	conn := db.Tx
	if !db.disablePrepare {
		var stmt *sql.Stmt
		//添加预处理
		stmt, errStmt := conn.PrepareContext(ctx, querySql)
		if err != nil {
			err = errStmt
			return
		}
		defer stmt.Close()
		result, err = stmt.ExecContext(ctx, args...)
	} else {
		result, err = conn.ExecContext(ctx, querySql, args...)
	}
	if !db.disableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(db.driverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			ext.PeerAddress.Set(span, db.addr)
			ext.DBInstance.Set(span, db.dbName)
			ext.DBStatement.Set(span, db.driverName)
			// span.SetTag("call.func", query.Func)
			// span.SetTag("call.path", query.Path)
			span.LogFields(log.String("sql", querySql))
			span.LogFields(log.Object("args", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	if !db.disableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(db.driverName, db.dbName, db.addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(db.driverName, db.dbName, db.addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(db.driverName, db.dbName, db.addr).Observe(cost.Seconds())
	}
	return
}

func (db txSession) QueryCtx(ctx context.Context, querySql string, args ...any) (result *sql.Rows, err error) {
	ctx = db.getCtx(ctx)
	start := time.Now()
	conn := db.Tx
	if !db.disablePrepare {
		var stmt *sql.Stmt
		//添加预处理
		stmt, errStmt := conn.PrepareContext(ctx, querySql)
		if err != nil {
			err = errStmt
			return
		}
		defer stmt.Close()
		result, err = stmt.QueryContext(ctx, args...)
	} else {
		result, err = conn.QueryContext(ctx, querySql, args...)
	}
	defer result.Close()
	if !db.disableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(db.driverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerHostname.Set(span, hostName)
			ext.PeerAddress.Set(span, db.addr)
			ext.DBInstance.Set(span, db.dbName)
			ext.DBStatement.Set(span, db.driverName)
			// span.SetTag("call.func", query.Func)
			// span.SetTag("call.path", query.Path)
			span.LogFields(log.String("sql", querySql))
			span.LogFields(log.Object("args", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	if !db.disableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(db.driverName, db.dbName, db.addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(db.driverName, db.dbName, db.addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(db.driverName, db.dbName, db.addr).Observe(cost.Seconds())
	}
	return
}
