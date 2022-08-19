package query

import (
	"context"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// BETWEEN ...
const (
	BETWEEN    = "BETWEEN"
	NOTBETWEEN = "NOT BETWEEN"
	IN         = "IN"
	NOTIN      = "NOT IN"
	AND        = "AND"
	OR         = "OR"
	ISNULL     = "IS NULL"
	ISNOTNULL  = "IS NOT NULL"
	EQUAL      = "="
	NOTEQUAL   = "!="
	LIKE       = "LIKE"
	JOIN       = "JOIN"
	INNERJOIN  = "INNER JOIN"
	LEFTJOIN   = "LEFT JOIN"
	RIGHTJOIN  = "RIGHT JOIN"
	UNION      = "UNION"
	UNIONALL   = "UNION ALL"
	DESC       = "DESC"
	ASC        = "ASC"
)

// Query 查询构造器
type Query struct {
	ctx        context.Context
	connection Connection
	table      []string
	columns    []string
	where      []w
	orders     []string
	groups     []string
	limit      int64
	offset     int64
	distinct   bool
	// binds      []string
	joins     []join
	unions    []union
	unLimit   int64
	unOffset  int64
	unOrders  []string
	args      []interface{}
	whereArgs []interface{}
	data      []map[string]interface{}
}
type join struct {
	table    string
	on       string
	operator string
}
type union struct {
	query    Query
	operator string
}
type w struct {
	column   string
	operator string
	valueNum int64
	do       string
}

// Table 设置操作的表名称
func (query *Query) Table(tableName ...string) *Query {
	query.table = tableName
	return query
}

// Select 查询字段
func (query *Query) Select(columns ...string) *Query {
	query.columns = columns
	return query
}

// Where 构造条件语句
func (query *Query) Where(column string, value ...interface{}) *Query {
	if len(value) == 0 { //一个参数直接where
		query.toWhere(column, "", 0, AND)
	} else if len(value) == 1 { //2个参数直接where =
		query.toWhere(column, EQUAL, 1, AND)
		query.addArg(value[0])
	} else { //3个参数
		switch v := value[0].(type) {
		case string:
			query.toWhere(column, v, 1, AND)
			query.addArg(value[1])
		}
	}
	return query
}

// OrWhere 构造OR条件
func (query *Query) OrWhere(column string, value ...interface{}) *Query {
	if len(value) == 0 { //一个参数直接where
		query.toWhere(column, "", 0, OR)
	} else if len(value) == 1 { //2个参数直接where =
		query.toWhere(column, EQUAL, 1, OR)
		query.addArg(value[0])
	} else {
		switch v := value[0].(type) {
		case string:
			query.toWhere(column, v, 1, OR)
			query.addArg(value[1])
		}
	}
	return query
}

// Equal 构造等于
func (query *Query) Equal(column string, value interface{}) *Query {
	query.toWhere(column, EQUAL, 1, AND)
	query.addArg(value)
	return query
}

// OrEqual 构造或者等于
func (query *Query) OrEqual(column string, value interface{}) *Query {
	query.toWhere(column, EQUAL, 1, OR)
	query.addArg(value)
	return query
}

// NotEqual 构造不等于
func (query *Query) NotEqual(column string, value interface{}) *Query {
	query.toWhere(column, NOTEQUAL, 1, AND)
	query.addArg(value)
	return query
}

// OrNotEqual 构造或者不等于
func (query *Query) OrNotEqual(column string, value interface{}) *Query {
	query.toWhere(column, NOTEQUAL, 1, OR)
	query.addArg(value)
	return query
}

// Between 构造Between
func (query *Query) Between(column string, value1 interface{}, value2 interface{}) *Query {
	query.toWhere(column, BETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}

// OrBetween 构造 或者 Between
func (query *Query) OrBetween(column string, value1 interface{}, value2 interface{}) *Query {
	query.toWhere(column, BETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}

// NotBetween 构造不Not Between
func (query *Query) NotBetween(column string, value1 interface{}, value2 interface{}) *Query {
	query.toWhere(column, NOTBETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}

// NotOrBetween 构造 Not Between  OR Not Between
func (query *Query) NotOrBetween(column string, value1 interface{}, value2 interface{}) *Query {
	query.toWhere(column, NOTBETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}

// In 构造 in语句
func (query *Query) In(column string, value ...interface{}) *Query {
	query.toWhere(column, IN, int64(len(value)), AND)
	query.addArg(value...)
	return query
}

// OrIn orin语句
func (query *Query) OrIn(column string, value ...interface{}) *Query {
	query.toWhere(column, IN, int64(len(value)), OR)
	query.addArg(value...)
	return query
}

// NotIn .
func (query *Query) NotIn(column string, value ...interface{}) *Query {
	query.toWhere(column, NOTIN, int64(len(value)), AND)
	query.addArg(value...)
	return query
}

// OrNotIn .
func (query *Query) OrNotIn(column string, value ...interface{}) *Query {
	query.toWhere(column, NOTIN, int64(len(value)), OR)
	query.addArg(value...)
	return query
}

// IsNULL .
func (query *Query) IsNULL(column string) *Query {
	query.toWhere(column, ISNULL, 0, AND)
	return query
}

// OrIsNULL .
func (query *Query) OrIsNULL(column string) *Query {
	query.toWhere(column, ISNULL, 0, OR)
	return query
}

// IsNotNULL .
func (query *Query) IsNotNULL(column string) *Query {
	query.toWhere(column, ISNOTNULL, 0, AND)
	return query
}

// OrIsNotNULL .
func (query *Query) OrIsNotNULL(column string) *Query {
	query.toWhere(column, ISNOTNULL, 0, OR)
	return query
}

// Like .
func (query *Query) Like(column string, value interface{}) *Query {
	query.toWhere(column, LIKE, 1, AND)
	query.addArg(value)
	return query
}

// OrLike .
func (query *Query) OrLike(column string, value interface{}) *Query {
	query.toWhere(column, LIKE, 1, OR)
	query.addArg(value)
	return query
}

// Join .
func (query *Query) Join(tablename string, on string) *Query {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: JOIN})
	return query
}

// InnerJoin ...
func (query *Query) InnerJoin(tablename string, on string) *Query {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: INNERJOIN})
	return query
}

// LeftJoin .
func (query *Query) LeftJoin(tablename string, on string) *Query {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: LEFTJOIN})
	return query
}

// RightJoin .
func (query *Query) RightJoin(tablename string, on string) *Query {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: RIGHTJOIN})
	return query
}

// Union .
func (query *Query) Union(unions ...Query) *Query {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNION})
		query.addArg(unions[i].args...)
	}
	return query
}

// UnionOffset .
func (query *Query) UnionOffset(offset int64) *Query {
	query.unOffset = offset
	return query
}

// UnionLimit .
func (query *Query) UnionLimit(limit int64) *Query {
	query.unLimit = limit
	return query
}

// UnionOrderBy .
func (query *Query) UnionOrderBy(column string, direction string) *Query {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.unOrders = append(query.unOrders, column)
	return query
}

// UnionAll .
func (query *Query) UnionAll(unions ...Query) *Query {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNIONALL})
		query.addArg(unions[i].args...)
	}
	return query
}

// Distinct .
func (query *Query) Distinct() *Query {
	query.distinct = true
	return query
}

// GroupBy .
func (query *Query) GroupBy(groups ...string) *Query {
	query.groups = groups
	return query
}

// OrderBy .
func (query *Query) OrderBy(column string, direction string) *Query {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.orders = append(query.orders, column)
	return query
}

// Offset .
func (query *Query) Offset(offset int64) *Query {
	query.offset = offset
	return query
}

// Skip .
func (query *Query) Skip(offset int64) *Query {
	query.offset = offset
	return query
}

// Limit .
func (query *Query) Limit(limit int64) *Query {
	query.limit = limit
	return query
}

// ToSQL 输出SQL语句
func (query *Query) ToSQL(method string) string {
	grammar := Grammar{query: query, method: method}
	return grammar.ToSQL()
}
func (query *Query) toWhere(column string, operator string, valueNum int64, do string) *Query {
	query.where = append(
		query.where,
		w{column: column, operator: operator, valueNum: valueNum, do: do})
	return query
}
func (query *Query) addArg(value ...interface{}) {
	query.args = append(query.args, value...)
}

func (query *Query) beforeArg(value ...interface{}) {
	query.whereArgs = append(query.whereArgs, value...)
}

func (query *Query) setData(data ...map[string]interface{}) {
	query.data = data
}

func (query *Query) getInsertMap(data interface{}) (columns []string, values map[string][]interface{}, err error) {
	stValue := reflect.Indirect(reflect.ValueOf(data))

	values = make(map[string][]interface{}, 0)
	switch stValue.Kind() {
	case reflect.Struct:
		var ignore bool
		for i := 0; i < stValue.NumField(); i++ {

			v := reflect.Indirect(stValue.Field(i))

			//处理嵌套的struct中的db映射字段
			if v.Kind() == reflect.Struct {

				var ignore bool

				switch v.Interface().(type) {
				case time.Time:
					ignore = true
				case NullDateTime:
					ignore = true
				case NullString:
					ignore = true
				case NullBool:
					ignore = true
				case NullInt64:
					ignore = true
				case NullInt32:
					ignore = true
				case NullFloat64:
					ignore = true
				case NullDate:
					ignore = true
				}

				if !ignore {
					cols, vals, err := query.getInsertMap(v.Interface())
					if err != nil {
						return nil, nil, err
					}

					for _, column := range cols {
						if _, ok := values[column]; !ok {
							columns = append(columns, column)
						}
					}

					for column, v := range vals {
						if _, ok := values[column]; ok {
							values[column] = append(values[column], v...)
						} else {
							values[column] = v
						}
					}
				}
			}

			tag := stValue.Type().Field(i).Tag.Get("db")
			attrList := strings.Split(tag, ",")
			ignore = false

			if len(attrList) > 0 {
				for _, attr := range attrList {
					if attr == "-" {
						ignore = true
						break
					}
				}
			}

			if ignore {
				continue
			}

			column := attrList[0]
			if column != "" && !query.IsZero(v) {
				if _, ok := values[column]; ok {
					values[column] = append(values[column], v.Interface())
				} else {
					columns = append(columns, column)
					values[column] = []interface{}{v.Interface()}
				}
			}
		}
	case reflect.Map:
		keys := stValue.MapKeys()
		for _, k := range keys {
			column := k.String()
			if _, ok := values[column]; ok {
				values[column] = append(values[column], stValue.MapIndex(k).Interface())
			} else {
				columns = append(columns, column)
				values[column] = []interface{}{stValue.MapIndex(k).Interface()}
			}
		}
	case reflect.Slice:
		n := stValue.Len()
		for i := 0; i < n; i++ {

			item := stValue.Index(i)
			cols, vals, err := query.getInsertMap(item.Interface())

			if err != nil {
				return nil, nil, err
			}

			for _, column := range cols {
				if _, ok := values[column]; !ok {
					columns = append(columns, column)
				}
			}

			for column, v := range vals {
				if _, ok := values[column]; ok {
					values[column] = append(values[column], v...)
				} else {
					values[column] = v
				}
			}
		}
	}
	return
}

// IsZero ...
func (query *Query) IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !v.Index(i).IsZero() {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				return false
			}
		}
		return true
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(&reflect.ValueError{
			Method: "reflect.Value.IsZero",
			Kind:   v.Kind(),
		})
	}
}

// MultiInsert 批量插入
func (query *Query) MultiInsert(datas ...interface{}) (int64, error) {

	stVal := reflect.ValueOf(datas)
	if stVal.Kind() != reflect.Slice {
		return 0, errors.New("data is not []interface{} type")
	}
	n := stVal.Len()
	if n > 0 {
		columns, values, err := query.getInsertMap(datas)
		if err != nil {
			return 0, err
		}
		bindingsArr := make([]map[string]interface{}, n)
		for i := 0; i < n; i++ {
			bindings := make(map[string]interface{}, 0)
			for _, column := range columns {
				bindings[column] = values[column][i]
			}
			bindingsArr[i] = bindings
		}
		query.setData(bindingsArr...)
		grammar := Grammar{query: query}
		sql := grammar.Insert()
		if len(query.columns) < 1 {
			return 0, errors.New("insert data cannot be empty")
		}
		result, err := query.connection.Exec(query.ctx, sql, query.args...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New("insert data cannot be empty")

}

// MultiInsertSQL 批量插入
func (query *Query) MultiInsertSQL(datas ...interface{}) string {
	stVal := reflect.ValueOf(datas)
	if stVal.Kind() != reflect.Slice {
		return ""
	}
	n := stVal.Len()
	if n > 0 {
		columns, values, err := query.getInsertMap(datas)
		if err != nil {
			return ""
		}
		bindingsArr := make([]map[string]interface{}, n)
		for i := 0; i < n; i++ {
			bindings := make(map[string]interface{}, 0)
			for _, column := range columns {
				bindings[column] = values[column][i]
			}
			bindingsArr[i] = bindings
		}
		query.setData(bindingsArr...)
		grammar := Grammar{query: query}
		sql := grammar.Insert()
		if len(query.columns) < 1 {
			return ""
		}
		query.connection.LastSQL(sql, query.args...)
		return query.connection.SQLRaw()
	}
	return ""
}

// Replace 替换
func (query *Query) Replace(datas ...interface{}) (int64, error) {

	stVal := reflect.ValueOf(datas)
	if stVal.Kind() != reflect.Slice {
		return 0, errors.New("data is not []interface{} type")
	}
	n := stVal.Len()
	if n > 0 {
		columns, values, err := query.getInsertMap(datas)
		if err != nil {
			return 0, err
		}
		bindingsArr := make([]map[string]interface{}, n)
		for i := 0; i < n; i++ {
			bindings := make(map[string]interface{}, 0)
			for _, column := range columns {
				bindings[column] = values[column][i]
			}
			bindingsArr[i] = bindings
		}
		query.setData(bindingsArr...)
		grammar := Grammar{query: query}
		sql := grammar.Replace()
		if len(query.columns) < 1 {
			return 0, errors.New("insert data cannot be empty")
		}
		result, err := query.connection.Exec(query.ctx, sql, query.args...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New("insert data cannot be empty")
}

// ReplaceSQL 替换
func (query *Query) ReplaceSQL(datas ...interface{}) string {

	stVal := reflect.ValueOf(datas)
	if stVal.Kind() != reflect.Slice {
		return ""
	}
	n := stVal.Len()
	if n > 0 {
		columns, values, err := query.getInsertMap(datas)
		if err != nil {
			return ""
		}
		bindingsArr := make([]map[string]interface{}, n)
		for i := 0; i < n; i++ {
			bindings := make(map[string]interface{}, 0)
			for _, column := range columns {
				bindings[column] = values[column][i]
			}
			bindingsArr[i] = bindings
		}
		query.setData(bindingsArr...)
		grammar := Grammar{query: query}
		sql := grammar.Replace()
		if len(query.columns) < 1 {
			return ""
		}
		query.connection.LastSQL(sql, query.args...)
		return query.connection.SQLRaw()
	}
	return ""
}

// InsertUpdate ...
func (query *Query) InsertUpdate(insert interface{}, update interface{}) (int64, error) {

	columns, values, err := query.getInsertMap(insert)
	if err != nil {
		return 0, err
	}
	bindingsInsert := map[string]interface{}{}
	for _, column := range columns {
		bindingsInsert[column] = values[column][0]
	}

	columnsup, valuesup, errup := query.getInsertMap(update)
	if errup != nil {
		return 0, errup
	}
	bindingsUpdate := map[string]interface{}{}
	for _, column := range columnsup {
		bindingsUpdate[column] = valuesup[column][0]
	}

	query.setData(bindingsInsert, bindingsUpdate)
	grammar := Grammar{query: query}
	sql := grammar.InsertUpdate()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertUpdateSQL ...
func (query *Query) InsertUpdateSQL(insert interface{}, update interface{}) string {

	columns, values, err := query.getInsertMap(insert)
	if err != nil {
		return err.Error()
	}
	bindingsInsert := map[string]interface{}{}
	for _, column := range columns {
		bindingsInsert[column] = values[column][0]
	}

	columnsup, valuesup, errup := query.getInsertMap(update)
	if errup != nil {
		return errup.Error()
	}
	bindingsUpdate := map[string]interface{}{}
	for _, column := range columnsup {
		bindingsUpdate[column] = valuesup[column][0]
	}

	query.setData(bindingsInsert, bindingsUpdate)
	grammar := Grammar{query: query}
	sql := grammar.InsertUpdate()
	query.connection.LastSQL(sql, query.args...)
	return query.connection.SQLRaw()
}

// Insert 插入数据
func (query *Query) Insert(data interface{}) (int64, error) {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return 0, err
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{query: query}
	sql := grammar.Insert()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertSQL 获取SQL语句
func (query *Query) InsertSQL(data interface{}) string {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return ""
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{query: query}
	sql := grammar.Insert()
	query.connection.LastSQL(sql, query.args...)
	return query.connection.SQLRaw()
}

// Update 更新
func (query *Query) Update(data interface{}) (int64, error) {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return 0, err
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{query: query}
	sql := grammar.Update()
	args := append(query.whereArgs, query.args...)

	result, err := query.connection.Exec(query.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateSQL 更新
func (query *Query) UpdateSQL(data interface{}) string {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return ""
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{query: query}
	sql := grammar.Update()
	args := append(query.whereArgs, query.args...)
	query.connection.LastSQL(sql, args...)
	return query.connection.SQLRaw()
}

// Delete .
func (query *Query) Delete() (int64, error) {
	grammar := Grammar{query: query}
	sql := grammar.Delete()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteSQL .
func (query *Query) DeleteSQL() string {
	grammar := Grammar{query: query}
	sql := grammar.Delete()
	query.connection.LastSQL(sql, query.args...)
	return query.connection.SQLRaw()
}

// Count ...
func (query *Query) Count() (int64, error) {
	query.Select("COUNT(1) AS _C")
	d, err := query.Row().ToMap()
	if err != nil || d == nil {
		return 0, err
	}
	if len(d) < 1 {
		return 0, nil
	}
	v := d["_C"]
	return strconv.ParseInt(v, 10, 0)
}

// Exec 原始SQl语句执行
func (query *Query) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := query.connection.Exec(query.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ExecSQL 原始SQl语句执行
func (query *Query) ExecSQL(sql string, args ...interface{}) string {
	query.connection.LastSQL(sql, args...)
	return query.connection.SQLRaw()
}

// QueryRows ...
func (query *Query) QueryRows(sql string, args ...interface{}) *Rows {
	rows, err := query.connection.Query(query.ctx, sql, args...)
	if err != nil {
		return &Rows{rs: nil, lastError: err}
	}
	return &Rows{rs: rows, lastError: err}
}

// QueryRowsSQL ...
func (query *Query) QueryRowsSQL(sql string, args ...interface{}) string {
	query.connection.LastSQL(sql, args...)
	return query.connection.SQLRaw()
}

// QueryRowSQL ...
func (query *Query) QueryRowSQL(sql string, args ...interface{}) string {
	query.connection.LastSQL(sql, args...)
	return query.connection.SQLRaw()
}

// QueryRow ...
func (query *Query) QueryRow(sql string, args ...interface{}) *Row {
	rs := query.QueryRows(sql, args...)
	r := new(Row)
	r.rs = rs
	return r
}

// Row 获取一条记录
func (query *Query) Row() *Row {
	query.offset = 0
	query.limit = 1
	rs := query.Rows()
	r := new(Row)
	r.rs = rs
	return r
}

// RowSQL ...
func (query *Query) RowSQL() string {
	grammar := Grammar{query: query}
	sql := grammar.Select()

	query.connection.LastSQL(sql, query.args...)
	return query.connection.SQLRaw()
}

// RowsSQL ...
func (query *Query) RowsSQL() string {
	grammar := Grammar{query: query}
	sql := grammar.Select()

	query.connection.LastSQL(sql, query.args...)
	return query.connection.SQLRaw()
}

// Rows 获取多条记录
func (query *Query) Rows() *Rows {
	grammar := Grammar{query: query}
	sql := grammar.Select()
	rows, err := query.connection.Query(query.ctx, sql, query.args...)
	if err != nil {
		return &Rows{rs: nil, lastError: err}
	}
	return &Rows{rs: rows, lastError: err}
}
