package query

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/metric"
	"github.com/abulo/ratel/v3/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// Connection 链接
type Connection interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	NewQuery(ctx context.Context) *QueryBuilder
	SqlRaw() string
	LastSql(query string, args ...interface{})
}

// Sql sql语句
type Sql struct {
	Sql      string
	Args     []interface{}
	CostTime time.Duration
}

// QueryDb mysql 配置
type QueryDb struct {
	DB            *sql.DB
	Sql           Sql
	DriverName    string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	Prepare       bool
	DBName        string
	Addr          string
}

//QueryTx 事务
type QueryTx struct {
	TX            *sql.Tx
	Sql           Sql
	DriverName    string
	DisableMetric bool // 关闭指标采集
	DisableTrace  bool // 关闭链路追踪
	Prepare       bool
	DBName        string
	Addr          string
}

//NewQuery 生成一个新的查询构造器
func (querydb *QueryDb) NewQuery(ctx context.Context) *QueryBuilder {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return &QueryBuilder{connection: querydb, ctx: ctx}
}

//Begin 开启一个事务
func (querydb *QueryDb) Begin() (*QueryTx, error) {
	tx, err := querydb.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &QueryTx{TX: tx, DriverName: querydb.DriverName, DisableTrace: querydb.DisableTrace, DisableMetric: querydb.DisableMetric, Prepare: querydb.Prepare, DBName: querydb.DBName, Addr: querydb.Addr}, nil
}

//Exec 复用执行语句
func (querydb *QueryDb) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if querydb.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querydb.Sql.Sql = query
	querydb.Sql.Args = args
	start := time.Now()
	defer func() {
		querydb.Sql.CostTime = time.Since(start)
	}()
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if !querydb.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querydb.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, querydb.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, querydb.DBName)
			ext.DBStatement.Set(span, querydb.DriverName)
			span.LogFields(log.String("sql", querydb.SqlRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res sql.Result
	var err error
	var stmt *sql.Stmt
	if querydb.Prepare {
		//添加预处理
		stmt, err = querydb.DB.PrepareContext(ctx, query)
		if err != nil {
			querydb.DB.PingContext(ctx)
			return nil, err
		}
		defer stmt.Close()
		res, err = stmt.ExecContext(ctx, args...)
		querydb.DB.PingContext(ctx)
	} else {
		res, err = querydb.DB.ExecContext(ctx, query, args...)
		querydb.DB.PingContext(ctx)
	}

	if !querydb.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(querydb.DriverName, querydb.DBName, querydb.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(querydb.DriverName, querydb.DBName, querydb.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(querydb.DriverName, querydb.DBName, querydb.Addr).Observe(cost.Seconds())
	}

	return res, err
}

//Query 复用查询语句
func (querydb *QueryDb) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if querydb.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querydb.Sql.Sql = query
	querydb.Sql.Args = args
	start := time.Now()
	defer func() {
		querydb.Sql.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !querydb.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querydb.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, querydb.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, querydb.DBName)
			ext.DBStatement.Set(span, querydb.DriverName)
			span.LogFields(log.String("sql", querydb.SqlRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	var res *sql.Rows
	var err error
	var stmt *sql.Stmt

	if querydb.Prepare {
		//添加预处理
		stmt, err = querydb.DB.PrepareContext(ctx, query)
		if err != nil {
			querydb.DB.PingContext(ctx)
			return nil, err
		}
		defer stmt.Close()
		res, err = stmt.QueryContext(ctx, args...)
		querydb.DB.PingContext(ctx)
	} else {
		res, err = querydb.DB.QueryContext(ctx, query, args...)
		querydb.DB.PingContext(ctx)
	}

	if !querydb.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(querydb.DriverName, querydb.DBName, querydb.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(querydb.DriverName, querydb.DBName, querydb.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(querydb.DriverName, querydb.DBName, querydb.Addr).Observe(cost.Seconds())
	}

	return res, err
}

// Commit 事务提交
func (querytx *QueryTx) Commit() error {
	return querytx.TX.Commit()
}

// Rollback 事务回滚
func (querytx *QueryTx) Rollback() error {
	return querytx.TX.Rollback()
}

// NewQuery 生成一个新的查询构造器
func (querytx *QueryTx) NewQuery(ctx context.Context) *QueryBuilder {
	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	return &QueryBuilder{connection: querytx, ctx: ctx}
}

//Exec 复用执行语句
func (querytx *QueryTx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if querytx.TX == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querytx.Sql.Sql = query
	querytx.Sql.Args = args
	start := time.Now()
	defer func() {
		querytx.Sql.CostTime = time.Since(start)

	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !querytx.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querytx.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, querytx.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, querytx.DBName)
			ext.DBStatement.Set(span, querytx.DriverName)
			span.LogFields(log.String("sql", querytx.SqlRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res sql.Result
	var err error
	var stmt *sql.Stmt
	if querytx.Prepare {
		//添加预处理
		stmt, err = querytx.TX.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		res, err = stmt.ExecContext(ctx, args...)
	} else {
		res, err = querytx.TX.ExecContext(ctx, query, args...)
	}

	if !querytx.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(querytx.DriverName, querytx.DBName, querytx.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(querytx.DriverName, querytx.DBName, querytx.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(querytx.DriverName, querytx.DBName, querytx.Addr).Observe(cost.Seconds())
	}

	return res, err

}

//Query 复用查询语句
func (querytx *QueryTx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if querytx.TX == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querytx.Sql.Sql = query
	querytx.Sql.Args = args
	start := time.Now()
	defer func() {
		querytx.Sql.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if !querytx.DisableTrace {
		if parentSpan := trace.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querytx.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			hostName, err := os.Hostname()
			if err != nil {
				hostName = "unknown"
			}
			ext.PeerAddress.Set(span, querytx.Addr)
			ext.PeerHostname.Set(span, hostName)
			ext.DBInstance.Set(span, querytx.DBName)
			ext.DBStatement.Set(span, querytx.DriverName)
			span.LogFields(log.String("sql", querytx.SqlRaw()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res *sql.Rows
	var err error
	var stmt *sql.Stmt
	if querytx.Prepare {
		//添加预处理
		stmt, err = querytx.TX.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		res, err = stmt.QueryContext(ctx, args...)
	} else {
		res, err = querytx.TX.QueryContext(ctx, query, args...)
	}
	if !querytx.DisableMetric {
		cost := time.Since(start)
		if err != nil {
			metric.LibHandleCounter.WithLabelValues(querytx.DriverName, querytx.DBName, querytx.Addr, "ERR").Inc()
		} else {
			metric.LibHandleCounter.Inc(querytx.DriverName, querytx.DBName, querytx.Addr, "OK")
		}
		metric.LibHandleHistogram.WithLabelValues(querytx.DriverName, querytx.DBName, querytx.Addr).Observe(cost.Seconds())
	}
	return res, err
}

func (querytx *QueryTx) SqlRaw() string {
	return querytx.Sql.ToString()
}

func (querydb *QueryDb) SqlRaw() string {
	return querydb.Sql.ToString()
}

func (querytx *QueryTx) LastSql(query string, args ...interface{}) {
	querytx.Sql.Sql = query
	querytx.Sql.Args = args
}
func (querydb *QueryDb) LastSql(query string, args ...interface{}) {
	querydb.Sql.Sql = query
	querydb.Sql.Args = args
}

// ToString sql语句转出string
func (sqlRaw Sql) ToString() string {
	s := sqlRaw.Sql
	for _, v := range sqlRaw.Args {
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
		} else {
			return strings.Replace(s, "?", strconv.Quote(fmt.Sprintf("%v", v)), 1)
		}
	}
	return strings.Replace(s, "?", fmt.Sprintf("%v", v), 1)
}

// ToJson sql语句转出json
func (sql Sql) ToJson() string {
	return fmt.Sprintf(`{"sql":%s,"costtime":"%s"}`, strconv.Quote(sql.ToString()), sql.CostTime)
}
