package clickhouse

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

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

// QueryBuilder 查询构造器
type QueryBuilder struct {
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
	binds      []string
	joins      []join
	unions     []union
	unlimmit   int64
	unoffset   int64
	unorders   []string

	args      []interface{}
	whereArgs []interface{}
	datas     []map[string]interface{}
}
type join struct {
	table    string
	on       string
	operator string
}
type union struct {
	query    QueryBuilder
	operator string
}
type w struct {
	column   string
	operator string
	valuenum int64
	do       string
}

//Table 设置操作的表名称
func (query *QueryBuilder) Table(tablename ...string) *QueryBuilder {
	query.table = tablename
	return query
}

//Select 查询字段
func (query *QueryBuilder) Select(columns ...string) *QueryBuilder {
	query.columns = columns
	return query
}

//Where 构造条件语句
func (query *QueryBuilder) Where(column string, value ...interface{}) *QueryBuilder {
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

//OrWhere 构造OR条件
func (query *QueryBuilder) OrWhere(column string, value ...interface{}) *QueryBuilder {
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

//Equal 构造等于
func (query *QueryBuilder) Equal(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, EQUAL, 1, AND)
	query.addArg(value)
	return query
}

// OrEqual 构造或者等于
func (query *QueryBuilder) OrEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, EQUAL, 1, OR)
	query.addArg(value)
	return query
}

//NotEqual 构造不等于
func (query *QueryBuilder) NotEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, NOTEQUAL, 1, AND)
	query.addArg(value)
	return query
}

//OrNotEqual 构造或者不等于
func (query *QueryBuilder) OrNotEqual(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, NOTEQUAL, 1, OR)
	query.addArg(value)
	return query
}

//Between 构造Between
func (query *QueryBuilder) Between(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, BETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}

//OrBetween 构造 或者 Between
func (query *QueryBuilder) OrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, BETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}

// NotBetween 构造不Not Between
func (query *QueryBuilder) NotBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, NOTBETWEEN, 2, AND)
	query.addArg(value1, value2)
	return query
}

// NotOrBetween 构造 Not Between  OR Not Between
func (query *QueryBuilder) NotOrBetween(column string, value1 interface{}, value2 interface{}) *QueryBuilder {
	query.toWhere(column, NOTBETWEEN, 2, OR)
	query.addArg(value1, value2)
	return query
}

// In 构造 in语句
func (query *QueryBuilder) In(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, IN, int64(len(value)), AND)
	query.addArg(value...)
	return query
}

// OrIn orin语句
func (query *QueryBuilder) OrIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, IN, int64(len(value)), OR)
	query.addArg(value...)
	return query
}

//NotIn .
func (query *QueryBuilder) NotIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, NOTIN, int64(len(value)), AND)
	query.addArg(value...)
	return query
}

//OrNotIn .
func (query *QueryBuilder) OrNotIn(column string, value ...interface{}) *QueryBuilder {
	query.toWhere(column, NOTIN, int64(len(value)), OR)
	query.addArg(value...)
	return query
}

//IsNULL .
func (query *QueryBuilder) IsNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNULL, 0, AND)
	return query
}

//OrIsNULL .
func (query *QueryBuilder) OrIsNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNULL, 0, OR)
	return query
}

//IsNotNULL .
func (query *QueryBuilder) IsNotNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNOTNULL, 0, AND)
	return query
}

//OrIsNotNULL .
func (query *QueryBuilder) OrIsNotNULL(column string) *QueryBuilder {
	query.toWhere(column, ISNOTNULL, 0, OR)
	return query
}

//Like .
func (query *QueryBuilder) Like(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, LIKE, 1, AND)
	query.addArg(value)
	return query
}

//OrLike .
func (query *QueryBuilder) OrLike(column string, value interface{}) *QueryBuilder {
	query.toWhere(column, LIKE, 1, OR)
	query.addArg(value)
	return query
}

//Join .
func (query *QueryBuilder) Join(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: JOIN})
	return query
}

func (query *QueryBuilder) InnerJoin(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: INNERJOIN})
	return query
}

//LeftJoin .
func (query *QueryBuilder) LeftJoin(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: LEFTJOIN})
	return query
}

//RightJoin .
func (query *QueryBuilder) RightJoin(tablename string, on string) *QueryBuilder {
	query.joins = append(query.joins, join{table: tablename, on: on, operator: RIGHTJOIN})
	return query
}

//Union .
func (query *QueryBuilder) Union(unions ...QueryBuilder) *QueryBuilder {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNION})
		query.addArg(unions[i].args...)
	}
	return query
}

//UnionOffset .
func (query *QueryBuilder) UnionOffset(offset int64) *QueryBuilder {
	query.unoffset = offset
	return query
}

//UnionLimit .
func (query *QueryBuilder) UnionLimit(limit int64) *QueryBuilder {
	query.unlimmit = limit
	return query
}

//UnionOrderBy .
func (query *QueryBuilder) UnionOrderBy(column string, direction string) *QueryBuilder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.unorders = append(query.unorders, column)
	return query
}

//UnionAll .
func (query *QueryBuilder) UnionAll(unions ...QueryBuilder) *QueryBuilder {
	for i, len := 0, len(unions); i < len; i++ {
		query.unions = append(query.unions, union{query: unions[i], operator: UNIONALL})
		query.addArg(unions[i].args...)
	}
	return query
}

// Distinct .
func (query *QueryBuilder) Distinct() *QueryBuilder {
	query.distinct = true
	return query
}

//GroupBy .
func (query *QueryBuilder) GroupBy(groups ...string) *QueryBuilder {
	query.groups = groups
	return query
}

//OrderBy .
func (query *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	query.orders = append(query.orders, column)
	return query
}

//Offset .
func (query *QueryBuilder) Offset(offset int64) *QueryBuilder {
	query.offset = offset
	return query
}

//Skip .
func (query *QueryBuilder) Skip(offset int64) *QueryBuilder {
	query.offset = offset
	return query
}

//Limit .
func (query *QueryBuilder) Limit(limit int64) *QueryBuilder {
	query.limit = limit
	return query
}

//ToSql 输出SQL语句
func (query *QueryBuilder) ToSql(method string) string {
	grammar := Grammar{builder: query, method: method}
	return grammar.ToSql()
}
func (query *QueryBuilder) toWhere(column string, operator string, valuenum int64, do string) *QueryBuilder {
	query.where = append(
		query.where,
		w{column: column, operator: operator, valuenum: valuenum, do: do})
	return query
}
func (query *QueryBuilder) addArg(value ...interface{}) {
	query.args = append(query.args, value...)
}

func (query *QueryBuilder) beforeArg(value ...interface{}) {
	query.whereArgs = append(query.whereArgs, value...)
}

func (query *QueryBuilder) setData(datas ...map[string]interface{}) {
	query.datas = datas
}

func (b *QueryBuilder) getInsertMap(data interface{}) (columns []string, values map[string][]interface{}, err error) {
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
					cols, vals, err := b.getInsertMap(v.Interface())
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
			if column != "" {
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
			cols, vals, err := b.getInsertMap(item.Interface())

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

//MultiInsert 批量插入
func (query *QueryBuilder) MultiInsert(datas ...interface{}) (int64, error) {

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
		grammar := Grammar{builder: query}
		sql := grammar.Insert()
		if len(query.columns) < 1 {
			return 0, errors.New("insert data cannot be empty")
		}
		result, err := query.connection.Exec(query.ctx, sql, query.args...)
		if err != nil {
			err = NewDBError(err.Error(), query.connection.GetLastSql())
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New("insert data cannot be empty")

}

//MultiInsertSQL 批量插入
func (query *QueryBuilder) MultiInsertSQL(datas ...interface{}) string {
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
		grammar := Grammar{builder: query}
		sql := grammar.Insert()
		if len(query.columns) < 1 {
			return ""
		}
		query.connection.LastSql(sql, query.args...)
		return query.connection.GetLastSql().ToString()
	}
	return ""
}

//Replace 替换
func (query *QueryBuilder) Replace(datas ...interface{}) (int64, error) {

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
		grammar := Grammar{builder: query}
		sql := grammar.Replace()
		if len(query.columns) < 1 {
			return 0, errors.New("insert data cannot be empty")
		}
		result, err := query.connection.Exec(query.ctx, sql, query.args...)
		if err != nil {
			err = NewDBError(err.Error(), query.connection.GetLastSql())
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New("insert data cannot be empty")
}

//ReplaceSQL 替换
func (query *QueryBuilder) ReplaceSQL(datas ...interface{}) string {

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
		grammar := Grammar{builder: query}
		sql := grammar.Replace()
		if len(query.columns) < 1 {
			return ""
		}
		query.connection.LastSql(sql, query.args...)
		return query.connection.GetLastSql().ToString()
	}
	return ""
}

//InsertUpdate ...
func (query *QueryBuilder) InsertUpdate(insert interface{}, update interface{}) (int64, error) {

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
	grammar := Grammar{builder: query}
	sql := grammar.InsertUpdate()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return 0, err
	}
	return result.RowsAffected()
}

//InsertUpdateSQL ...
func (query *QueryBuilder) InsertUpdateSQL(insert interface{}, update interface{}) string {

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
	grammar := Grammar{builder: query}
	sql := grammar.InsertUpdate()
	query.connection.LastSql(sql, query.args...)
	return query.connection.GetLastSql().ToString()
}

//Insert 插入数据
func (query *QueryBuilder) Insert(data interface{}) (int64, error) {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return 0, err
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{builder: query}
	sql := grammar.Insert()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return 0, err
	}
	return result.LastInsertId()
}

//InsertSQL 获取SQL语句
func (query *QueryBuilder) InsertSQL(data interface{}) string {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return ""
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{builder: query}
	sql := grammar.Insert()
	query.connection.LastSql(sql, query.args...)
	return query.connection.GetLastSql().ToString()
}

//Update 更新
func (query *QueryBuilder) Update(data interface{}) (int64, error) {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return 0, err
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{builder: query}
	sql := grammar.Update()
	args := append(query.whereArgs, query.args...)

	result, err := query.connection.Exec(query.ctx, sql, args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return 0, err
	}
	return result.RowsAffected()
}

//UpdateSQL 更新
func (query *QueryBuilder) UpdateSQL(data interface{}) string {
	columns, values, err := query.getInsertMap(data)
	if err != nil {
		return ""
	}
	bindings := map[string]interface{}{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	query.setData(bindings)
	grammar := Grammar{builder: query}
	sql := grammar.Update()
	args := append(query.whereArgs, query.args...)
	query.connection.LastSql(sql, args...)
	return query.connection.GetLastSql().ToString()
}

//Delete .
func (query *QueryBuilder) Delete() (int64, error) {
	grammar := Grammar{builder: query}
	sql := grammar.Delete()
	result, err := query.connection.Exec(query.ctx, sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return 0, err
	}
	return result.RowsAffected()
}

//DeleteSQL .
func (query *QueryBuilder) DeleteSQL() string {
	grammar := Grammar{builder: query}
	sql := grammar.Delete()
	query.connection.LastSql(sql, query.args...)
	return query.connection.GetLastSql().ToString()
}

//Count ...
func (query *QueryBuilder) Count() (int64, error) {
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

//Exec 原始SQl语句执行
func (query *QueryBuilder) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := query.connection.Exec(query.ctx, sql, args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return 0, err
	}
	return result.RowsAffected()
}

//ExecSQL 原始SQl语句执行
func (query *QueryBuilder) ExecSQL(sql string, args ...interface{}) string {
	query.connection.LastSql(sql, args...)
	return query.connection.GetLastSql().ToString()
}

// QueryRows ...
func (query *QueryBuilder) QueryRows(sql string, args ...interface{}) *Rows {
	rows, err := query.connection.Query(query.ctx, sql, args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return &Rows{rs: nil, lastError: err}
	}
	return &Rows{rs: rows, lastError: err}
}

//QueryRowsSQL ...
func (query *QueryBuilder) QueryRowsSQL(sql string, args ...interface{}) string {
	query.connection.LastSql(sql, args...)
	return query.connection.GetLastSql().ToString()
}

//QueryRowSQL ...
func (query *QueryBuilder) QueryRowSQL(sql string, args ...interface{}) string {
	query.connection.LastSql(sql, args...)
	return query.connection.GetLastSql().ToString()
}

// QueryRow ...
func (query *QueryBuilder) QueryRow(sql string, args ...interface{}) *Row {
	rs := query.QueryRows(sql, args...)
	r := new(Row)
	r.rs = rs
	return r
}

//Row 获取一条记录
func (query *QueryBuilder) Row() *Row {
	query.offset = 0
	query.limit = 1
	rs := query.Rows()
	r := new(Row)
	r.rs = rs
	return r
}

// RowSQL ...
func (query *QueryBuilder) RowSQL() string {
	grammar := Grammar{builder: query}
	sql := grammar.Select()

	query.connection.LastSql(sql, query.args...)
	return query.connection.GetLastSql().ToString()
}

// RowsSQL ...
func (query *QueryBuilder) RowsSQL() string {
	grammar := Grammar{builder: query}
	sql := grammar.Select()

	query.connection.LastSql(sql, query.args...)
	return query.connection.GetLastSql().ToString()
}

//Rows 获取多条记录
func (query *QueryBuilder) Rows() *Rows {
	grammar := Grammar{builder: query}
	sql := grammar.Select()
	rows, err := query.connection.Query(query.ctx, sql, query.args...)
	if err != nil {
		err = NewDBError(err.Error(), query.connection.GetLastSql())
		return &Rows{rs: nil, lastError: err}
	}
	return &Rows{rs: rows, lastError: err}
}
