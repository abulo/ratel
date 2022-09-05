package query

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/metric"
	"github.com/abulo/ratel/v3/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Connection 链接
type Connection interface {
	Exec(ctx context.Context, querySQL string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, querySQL string, args ...interface{}) (*sql.Rows, error)
	NewBuilder(ctx context.Context) *Builder
	SQLRaw() string
	LastSQL(querySQL string, args ...interface{})
}

// SQL sql语句
type SQL struct {
	SQL      string
	Args     []interface{}
	CostTime time.Duration
}

// Query 构造查询
type Query struct {
	DB            *sql.DB
	SQL           SQL
	DriverName    string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	Prepare       bool
	DBName        string
	Addr          string
	Func          string // 上个函数调用的函数名称
	Path          string // 上个函数调用的函数位置
}

// Transaction 事务
type Transaction struct {
	Transaction   *sql.Tx
	SQL           SQL
	DriverName    string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	Prepare       bool
	DBName        string
	Addr          string
	Func          string // 上个函数调用的函数名称
	Path          string // 上个函数调用的函数位置
}

// NewBuilder 生成一个新的查询构造器
func (query *Query) NewBuilder(ctx context.Context) *Builder {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !query.DisableTrace {
		pc, file, lineNo, _ := runtime.Caller(1)
		name := runtime.FuncForPC(pc).Name()
		query.Path = file + ":" + cast.ToString(lineNo)
		query.Func = name
	}
	return &Builder{connection: query, ctx: ctx}
}

// Begin 开启一个事务
func (query *Query) Begin() (*Transaction, error) {
	transaction, err := query.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{Transaction: transaction, DriverName: query.DriverName, DisableTrace: query.DisableTrace, DisableMetric: query.DisableMetric, Prepare: query.Prepare, DBName: query.DBName, Addr: query.Addr, Func: query.Func, Path: query.Path}, nil
}

// Exec 复用执行语句
func (query *Query) Exec(ctx context.Context, querySQL string, args ...interface{}) (sql.Result, error) {
	if query.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	query.SQL.SQL = querySQL
	query.SQL.Args = args
	start := time.Now()
	defer func() {
		query.SQL.CostTime = time.Since(start)
	}()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if !query.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(query.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, query.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, query.DBName)
			ext.DBStatement.Set(span, query.DriverName)
			span.SetTag("call.func", query.Func)
			span.SetTag("call.path", query.Path)
			span.LogFields(log.String("sql", query.SQLRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res sql.Result
	var err error
	var stmt *sql.Stmt
	if query.Prepare {
		//添加预处理
		stmt, err = query.DB.PrepareContext(ctx, querySQL)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				logger.Logger.Error("Error closing stmt: ", err)
			}
		}()
		res, err = stmt.ExecContext(ctx, args...)
	} else {
		res, err = query.DB.ExecContext(ctx, querySQL, args...)
	}

	if !query.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(query.DriverName, query.DBName, query.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(query.DriverName, query.DBName, query.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(query.DriverName, query.DBName, query.Addr).Observe(cost.Seconds())
	}

	return res, err
}

// Query 复用查询语句
func (query *Query) Query(ctx context.Context, querySQL string, args ...interface{}) (*sql.Rows, error) {
	if query.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	query.SQL.SQL = querySQL
	query.SQL.Args = args
	start := time.Now()
	defer func() {
		query.SQL.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !query.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(query.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, query.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, query.DBName)
			ext.DBStatement.Set(span, query.DriverName)
			span.SetTag("call.func", query.Func)
			span.SetTag("call.path", query.Path)
			span.LogFields(log.String("sql", query.SQLRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	var res *sql.Rows
	var err error
	var stmt *sql.Stmt

	if query.Prepare {
		//添加预处理
		stmt, err = query.DB.PrepareContext(ctx, querySQL)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				logger.Logger.Error("Error closing stmt: ", err)
			}
		}()
		res, err = stmt.QueryContext(ctx, args...)
	} else {
		res, err = query.DB.QueryContext(ctx, querySQL, args...)
	}

	if !query.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(query.DriverName, query.DBName, query.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(query.DriverName, query.DBName, query.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(query.DriverName, query.DBName, query.Addr).Observe(cost.Seconds())
	}

	return res, err
}

// Commit 事务提交
func (query *Transaction) Commit() error {
	return query.Transaction.Commit()
}

// Rollback 事务回滚
func (query *Transaction) Rollback() error {
	return query.Transaction.Rollback()
}

// NewBuilder 生成一个新的查询构造器
func (query *Transaction) NewBuilder(ctx context.Context) *Builder {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return &Builder{connection: query, ctx: ctx}
}

// Exec 复用执行语句
func (query *Transaction) Exec(ctx context.Context, querySQL string, args ...interface{}) (sql.Result, error) {
	if query.Transaction == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	query.SQL.SQL = querySQL
	query.SQL.Args = args
	start := time.Now()
	defer func() {
		query.SQL.CostTime = time.Since(start)

	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !query.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(query.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, query.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, query.DBName)
			ext.DBStatement.Set(span, query.DriverName)
			span.SetTag("call.func", query.Func)
			span.SetTag("call.path", query.Path)
			span.LogFields(log.String("sql", query.SQLRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res sql.Result
	var err error
	var stmt *sql.Stmt
	if query.Prepare {
		//添加预处理
		stmt, err = query.Transaction.PrepareContext(ctx, querySQL)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				logger.Logger.Error("Error closing stmt: ", err)
			}
		}()
		res, err = stmt.ExecContext(ctx, args...)
	} else {
		res, err = query.Transaction.ExecContext(ctx, querySQL, args...)
	}

	if !query.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(query.DriverName, query.DBName, query.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(query.DriverName, query.DBName, query.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(query.DriverName, query.DBName, query.Addr).Observe(cost.Seconds())
	}

	return res, err

}

// Query 复用查询语句
func (query *Transaction) Query(ctx context.Context, querySQL string, args ...interface{}) (*sql.Rows, error) {
	if query.Transaction == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	query.SQL.SQL = querySQL
	query.SQL.Args = args
	start := time.Now()
	defer func() {
		query.SQL.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !query.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(query.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, query.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, query.DBName)
			ext.DBStatement.Set(span, query.DriverName)
			span.LogFields(log.String("sql", query.SQLRaw()))
			span.SetTag("call.func", query.Func)
			span.SetTag("call.path", query.Path)
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res *sql.Rows
	var err error
	var stmt *sql.Stmt
	if query.Prepare {
		//添加预处理
		stmt, err = query.Transaction.PrepareContext(ctx, querySQL)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				logger.Logger.Error("Error closing stmt: ", err)
			}
		}()
		res, err = stmt.QueryContext(ctx, args...)
	} else {
		res, err = query.Transaction.QueryContext(ctx, querySQL, args...)
	}
	if !query.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(query.DriverName, query.DBName, query.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(query.DriverName, query.DBName, query.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(query.DriverName, query.DBName, query.Addr).Observe(cost.Seconds())
	}
	return res, err
}

// SQLRaw ...
func (query *Transaction) SQLRaw() string {
	return query.SQL.ToString()
}

// SQLRaw ...
func (query *Query) SQLRaw() string {
	return query.SQL.ToString()
}

// LastSQL ...
func (query *Transaction) LastSQL(querySQL string, args ...interface{}) {
	query.SQL.SQL = querySQL
	query.SQL.Args = args
}

// LastSQL ...
func (query *Query) LastSQL(querySQL string, args ...interface{}) {
	query.SQL.SQL = querySQL
	query.SQL.Args = args
}

// ToString sql语句转出string
func (SQLRaw SQL) ToString() string {
	s := SQLRaw.SQL
	for _, v := range SQLRaw.Args {
		if isNilFixed(v) {
			v = "NULL"
		} else {
			switch reflect.ValueOf(v).Interface().(type) {
			case NullString:
				v = v.(NullString).Result()
			case NullInt64:
				v = v.(NullInt64).Result()
			case NullInt32:
				v = v.(NullInt32).Result()
			case NullFloat64:
				v = v.(NullFloat64).Result()
			case NullBool:
				v = v.(NullBool).Result()
			case NullDateTime:
				v = v.(NullDateTime).Result()
			case NullDate:
				v = v.(NullDate).Result()
			}
		}
		s = convert(s, v)
	}
	return s
}

func isNilFixed(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func convert(s string, v interface{}) string {
	switch reflect.ValueOf(v).Interface().(type) {
	case string:
		if val := fmt.Sprintf("%v", v); val == "NULL" {
			return strings.Replace(s, "?", fmt.Sprintf("%v", v), 1)
		}
		return strings.Replace(s, "?", strconv.Quote(fmt.Sprintf("%v", v)), 1)
	}
	return strings.Replace(s, "?", fmt.Sprintf("%v", v), 1)
}

// ToJSON sql语句转出json
func (sql SQL) ToJSON() string {
	return fmt.Sprintf(`{"sql":%s,"costtime":"%s"}`, strconv.Quote(sql.ToString()), sql.CostTime)
}
