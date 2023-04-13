package sql

import (
	"context"
	"database/sql"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/metric"
	"github.com/abulo/ratel/v3/core/resource"
	"github.com/abulo/ratel/v3/core/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/cast"
)

type (

	// Session stands for raw connections or transaction sessions
	Session interface {
		MultiInsert(ctx context.Context, query string, args ...any) (int64, error)
		Replace(ctx context.Context, query string, args ...any) (int64, error)
		InsertUpdate(ctx context.Context, query string, args ...any) (int64, error)
		Insert(ctx context.Context, query string, args ...any) (int64, error)
		Update(ctx context.Context, query string, args ...any) (int64, error)
		Delete(ctx context.Context, query string, args ...any) (int64, error)
		Exec(ctx context.Context, query string, args ...any) (int64, error)
		Count(ctx context.Context, query string, args ...any) (int64, error)
		QueryRow(ctx context.Context, query string, args ...any) *Row
		QueryRows(ctx context.Context, query string, args ...any) *Rows
		ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
		QueryCtx(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}

	// SqlConn only stands for raw connections, so Transact method can be called.
	SqlConn interface {
		Session
		Transact(ctx context.Context, fn func(context.Context, Session) error) error
	}

	connProvider func() (*sql.DB, error)

	// SqlOption defines the method to customize a sql connection.
	SqlOption     func(*commonSqlConn)
	commonSqlConn struct {
		connProv connProvider
		// onError  func(context.Context, error)
		beginTx        beginnable
		brk            resource.Breaker
		accept         func(error) bool
		disableMetric  bool   // 关闭指标采集
		disableTrace   bool   // 关闭链路追踪
		disablePrepare bool   // 关闭预处理
		driverName     string // 驱动
		dbName         string
		addr           string
	}

	beginnable func() (trans, error)

	trans interface {
		Session
		Commit() error
		Rollback() error
	}

	txSession struct {
		*sql.Tx
		disableMetric  bool   // 关闭指标采集
		disableTrace   bool   // 关闭链路追踪
		disablePrepare bool   // 关闭预处理
		driverName     string // 驱动
		dbName         string
		addr           string
	}

	// sessionConn interface {
	// 	Exec(query string, args ...any) (sql.Result, error)
	// 	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// 	Query(query string, args ...any) (*sql.Rows, error)
	// 	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// }
)

// NewSqlConn returns a SqlConn with given driver name and dns.
func NewSqlConn(driverName, dns string, pool *pool, opts ...SqlOption) SqlConn {
	conn := &commonSqlConn{
		disableTrace:  pool.DisableTrace,
		disableMetric: pool.DisableMetric,
		driverName:    pool.DriverName,
		dbName:        pool.DbName,
		addr:          pool.Addr,
		connProv: func() (*sql.DB, error) {
			return getSqlConn(driverName, dns, pool)
		},
		// 事务处理
		beginTx: func() (trans, error) {
			db, err := getSqlConn(driverName, dns, pool)
			if err != nil {
				return nil, err
			}
			tx, err := db.Begin()
			if err != nil {
				return nil, err
			}
			return txSession{
				Tx:             tx,
				disableTrace:   pool.DisableTrace,
				disableMetric:  pool.DisableMetric,
				disablePrepare: pool.DisablePrepare,
				driverName:     pool.DriverName,
				dbName:         pool.DbName,
				addr:           pool.Addr,
			}, nil
		},
		brk: resource.NewBreaker(),
	}
	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// getCtx returns a context with a timeout.
func getCtx(ctx context.Context) context.Context {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return ctx
}

func ResultAccept(err error) error {
	if err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone || err == context.Canceled {
		return nil
	}
	return err
}

// Close closes the connection.
func (db *commonSqlConn) Close() error {
	return nil
}

// MultiInsert 批量插入
func (db *commonSqlConn) MultiInsert(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Replace 替换
func (db *commonSqlConn) Replace(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// InsertUpdate 插入或更新
func (db *commonSqlConn) InsertUpdate(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Insert 插入
func (db *commonSqlConn) Insert(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update 更新
func (db *commonSqlConn) Update(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete 删除
func (db *commonSqlConn) Delete(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Exec executes a query without returning any rows.
func (db *commonSqlConn) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryRow returns a single row from the database.
func (db *commonSqlConn) QueryRow(ctx context.Context, query string, args ...any) *Row {
	res, err := db.QueryCtx(ctx, query, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{
		rows: &Rows{rows: res},
		err:  err,
	}
}

func (db *commonSqlConn) Count(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.QueryCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	result := &Row{
		rows: &Rows{rows: res},
		err:  err,
	}
	d, err := result.ToMap()
	if err != nil || d == nil {
		return 0, err
	}
	if len(d) < 1 {
		return 0, nil
	}
	v := d["_C"]
	return strconv.ParseInt(v, 10, 0)
}

// QueryRow returns a single row from the database.
func (db *commonSqlConn) QueryRows(ctx context.Context, query string, args ...any) *Rows {
	res, err := db.QueryCtx(ctx, query, args...)
	if err != nil {
		return &Rows{err: err}
	}
	return &Rows{
		rows: res,
		err:  err,
	}
}

// ExecCtx executes a query without returning any rows.
func (db *commonSqlConn) ExecCtx(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	ctx = getCtx(ctx)
	start := time.Now()
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = db.connProv()
		if err != nil {
			return err
		}
		if !db.disableTrace {
			call := Caller(7)
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
				span.LogFields(log.String("sql", query))
				span.LogFields(log.Object("args", args))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)
			}
		}
		if !db.disablePrepare {
			var stmt *sql.Stmt
			//添加预处理
			stmt, err = conn.PrepareContext(ctx, query)
			if err != nil {
				return err
			}
			defer func() {
				if err := stmt.Close(); err != nil {
					logger.Logger.Error("Error closing stmt: ", err)
				}
			}()
			result, err = stmt.ExecContext(ctx, args...)
		} else {
			result, err = conn.ExecContext(ctx, query, args...)
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
		return err
	}, db.acceptable)
	return
}

// QueryCtx executes a query that returns rows, typically a SELECT.
func (db *commonSqlConn) QueryCtx(ctx context.Context, query string, args ...any) (result *sql.Rows, err error) {
	ctx = getCtx(ctx)
	start := time.Now()
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = db.connProv()
		if err != nil {
			return err
		}
		if !db.disableTrace {
			call := Caller(7)
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
				span.LogFields(log.String("sql", query))
				span.LogFields(log.Object("args", args))
				span.LogFields(log.Object("call", call))
				defer span.Finish()
				ctx = opentracing.ContextWithSpan(ctx, span)

			}
		}
		if !db.disablePrepare {
			var stmt *sql.Stmt
			//添加预处理
			stmt, err = conn.PrepareContext(ctx, query)
			if err != nil {
				return err
			}
			defer func() {
				if err := stmt.Close(); err != nil {
					logger.Logger.Error("Error closing stmt: ", err)
				}
			}()
			result, err = stmt.QueryContext(ctx, args...)
		} else {
			result, err = conn.QueryContext(ctx, query, args...)
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
		return err
	}, db.acceptable)
	return
}

// acceptable returns true if the error is acceptable.
func (db *commonSqlConn) acceptable(err error) bool {
	ok := err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone || err == context.Canceled
	if db.accept == nil {
		return ok
	}
	return ok || db.accept(err)
}

func (db *commonSqlConn) Transact(ctx context.Context, fn func(context.Context, Session) error) (err error) {
	ctx = getCtx(ctx)
	err = db.brk.DoWithAcceptable(func() error {
		pool := &pool{
			DisableTrace:   db.disableTrace,
			DisableMetric:  db.disableMetric,
			DisablePrepare: db.disablePrepare,
			DriverName:     db.driverName,
			DbName:         db.dbName,
			Addr:           db.addr,
		}
		return transact(ctx, db, pool, db.beginTx, fn)
	}, db.acceptable)
	return
}

func Caller(skip int) map[string]string {
	pc, file, lineNo, _ := runtime.Caller(skip)
	name := runtime.FuncForPC(pc).Name()
	return map[string]string{
		"path": file + ":" + cast.ToString(lineNo),
		"func": name,
	}
}
