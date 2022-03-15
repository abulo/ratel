package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	GetLastSql() Sql
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
	DB         *sql.DB
	lastsql    Sql
	DriverName string
	Trace      bool
	Prepare    bool
}

//QueryTx 事务
type QueryTx struct {
	TX         *sql.Tx
	lastsql    Sql
	DriverName string
	Trace      bool
	Prepare    bool
}

func (querydb *QueryDb) SetTrace(t bool) {
	querydb.Trace = t
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
	return &QueryTx{TX: tx, DriverName: querydb.DriverName, Trace: querydb.Trace, Prepare: querydb.Prepare}, nil
}

//Exec 复用执行语句
func (querydb *QueryDb) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if querydb.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querydb.lastsql.Sql = query
	querydb.lastsql.Args = args
	start := time.Now()
	defer func() {
		querydb.lastsql.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if querydb.Trace {

		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querydb.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, querydb.DriverName)
			span.LogFields(log.String("sql", query))
			span.LogFields(log.Object("param", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}

	var res sql.Result
	var err error

	if querydb.Prepare {
		//添加预处理
		stmt, err := querydb.DB.PrepareContext(ctx, query)
		if err != nil {
			querydb.DB.PingContext(ctx)
			return res, err
		}
		defer stmt.Close()
		res, err = stmt.ExecContext(ctx, args...)
		querydb.DB.PingContext(ctx)
	} else {
		res, err = querydb.DB.ExecContext(ctx, query, args...)
		querydb.DB.PingContext(ctx)
	}

	return res, err
}

//Query 复用查询语句
func (querydb *QueryDb) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if querydb.DB == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querydb.lastsql.Sql = query
	querydb.lastsql.Args = args
	start := time.Now()
	defer func() {
		querydb.lastsql.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if querydb.Trace {

		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querydb.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, querydb.DriverName)
			span.LogFields(log.String("sql", query))
			span.LogFields(log.Object("param", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	var res *sql.Rows
	var err error

	if querydb.Prepare {
		//添加预处理
		stmt, err := querydb.DB.PrepareContext(ctx, query)
		if err != nil {
			querydb.DB.PingContext(ctx)
			return res, err
		}
		defer stmt.Close()
		res, err = stmt.QueryContext(ctx, args...)
		querydb.DB.PingContext(ctx)
	} else {
		res, err = querydb.DB.QueryContext(ctx, query, args...)
		querydb.DB.PingContext(ctx)
	}
	return res, err
}

//GetLastSql 获取sql语句
func (querydb *QueryDb) GetLastSql() Sql {
	return querydb.lastsql
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
	querytx.lastsql.Sql = query
	querytx.lastsql.Args = args
	start := time.Now()
	defer func() {
		querytx.lastsql.CostTime = time.Since(start)

	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}

	if querytx.Trace {

		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querytx.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, querytx.DriverName)
			span.LogFields(log.String("sql", query))
			span.LogFields(log.Object("param", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}

	}
	var res sql.Result
	var err error
	if querytx.Prepare {
		//添加预处理
		stmt, err := querytx.TX.PrepareContext(ctx, query)
		if err != nil {
			return res, err
		}
		res, err = stmt.ExecContext(ctx, args...)
	} else {
		res, err = querytx.TX.ExecContext(ctx, query, args...)
	}
	return res, err

}

//Query 复用查询语句
func (querytx *QueryTx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if querytx.TX == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}
	querytx.lastsql.Sql = query
	querytx.lastsql.Args = args
	start := time.Now()
	defer func() {
		querytx.lastsql.CostTime = time.Since(start)
	}()

	if ctx == nil || ctx.Err() != nil {
		ctx = context.TODO()
	}
	if querytx.Trace {
		if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
			parentCtx := parentSpan.Context()
			span := opentracing.StartSpan(querytx.DriverName, opentracing.ChildOf(parentCtx))
			ext.SpanKindRPCClient.Set(span)
			ext.PeerService.Set(span, querytx.DriverName)
			span.LogFields(log.String("sql", query))
			span.LogFields(log.Object("param", args))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	var res *sql.Rows
	var err error
	if querytx.Prepare {
		res, err = querytx.TX.QueryContext(ctx, query, args...)
	} else {
		//添加预处理
		stmt, err := querytx.TX.PrepareContext(ctx, query)
		if err != nil {
			return res, err
		}
		res, err = stmt.QueryContext(ctx, args...)
	}
	return res, err
}

//GetLastSql 获取sql语句
func (querytx *QueryTx) GetLastSql() Sql {
	return querytx.lastsql
}

func (querytx *QueryTx) LastSql(query string, args ...interface{}) {
	querytx.lastsql.Sql = query
	querytx.lastsql.Args = args
}

func (querydb *QueryDb) LastSql(query string, args ...interface{}) {
	querydb.lastsql.Sql = query
	querydb.lastsql.Args = args
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
