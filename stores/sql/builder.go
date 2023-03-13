package sql

import (
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/stores/null"
	"github.com/pkg/errors"
)

// BETWEEN ...
const (
	BETWEEN      = "BETWEEN"
	NOTBETWEEN   = "NOT BETWEEN"
	IN           = "IN"
	NOTIN        = "NOT IN"
	AND          = "AND"
	OR           = "OR"
	ISNULL       = "IS NULL"
	ISNOTNULL    = "IS NOT NULL"
	EQUAL        = "="
	NOTEQUAL     = "!="
	GREATER      = ">"
	GREATEREQUAL = ">="
	LESS         = "<"
	LESSEQUAL    = "<="
	LIKE         = "LIKE"
	JOIN         = "JOIN"
	INNERJOIN    = "INNER JOIN"
	LEFTJOIN     = "LEFT JOIN"
	RIGHTJOIN    = "RIGHT JOIN"
	UNION        = "UNION"
	UNIONALL     = "UNION ALL"
	DESC         = "DESC"
	ASC          = "ASC"
)

// Builder 查询构造器
type Builder struct {
	table     []string
	columns   []string
	where     []w
	orders    []string
	groups    []string
	limit     int64
	offset    int64
	distinct  bool
	joins     []join
	unions    []union
	unLimit   int64
	unOffset  int64
	unOrders  []string
	args      []any
	whereArgs []any
	data      []map[string]any
}
type join struct {
	table    string
	on       string
	operator string
}
type union struct {
	builder  Builder
	operator string
}
type w struct {
	column   string
	operator string
	valueNum int64
	do       string
}

// NewBuilder 创建查询构造器
func NewBuilder() *Builder {
	return &Builder{}
}

// Table 设置操作的表名称
func (builder *Builder) Table(tableName ...string) *Builder {
	builder.table = tableName
	return builder
}

// Select 查询字段
func (builder *Builder) Select(columns ...string) *Builder {
	builder.columns = columns
	return builder
}

// Where 构造条件语句
func (builder *Builder) Where(column string, value ...any) *Builder {
	if len(value) == 0 { //一个参数直接where
		builder.toWhere(column, "", 0, AND)
	} else if len(value) == 1 { //2个参数直接where =
		builder.toWhere(column, EQUAL, 1, AND)
		builder.addArg(value[0])
	} else { //3个参数
		switch v := value[0].(type) {
		case string:
			builder.toWhere(column, v, 1, AND)
			builder.addArg(value[1])
		}
	}
	return builder
}

// OrWhere 构造OR条件
func (builder *Builder) OrWhere(column string, value ...any) *Builder {
	if len(value) == 0 { //一个参数直接where
		builder.toWhere(column, "", 0, OR)
	} else if len(value) == 1 { //2个参数直接where =
		builder.toWhere(column, EQUAL, 1, OR)
		builder.addArg(value[0])
	} else {
		switch v := value[0].(type) {
		case string:
			builder.toWhere(column, v, 1, OR)
			builder.addArg(value[1])
		}
	}
	return builder
}

// Equal 构造等于
func (builder *Builder) Equal(column string, value any) *Builder {
	builder.toWhere(column, EQUAL, 1, AND)
	builder.addArg(value)
	return builder
}

// OrEqual 构造或者等于
func (builder *Builder) OrEqual(column string, value any) *Builder {
	builder.toWhere(column, EQUAL, 1, OR)
	builder.addArg(value)
	return builder
}

// NotEqual 构造不等于
func (builder *Builder) NotEqual(column string, value any) *Builder {
	builder.toWhere(column, NOTEQUAL, 1, AND)
	builder.addArg(value)
	return builder
}

// Greater 构造大于
func (builder *Builder) Greater(column string, value any) *Builder {
	builder.toWhere(column, GREATER, 1, AND)
	builder.addArg(value)
	return builder
}

// Greater 构造大于等于
func (builder *Builder) GreaterEqual(column string, value any) *Builder {
	builder.toWhere(column, GREATEREQUAL, 1, AND)
	builder.addArg(value)
	return builder
}

// Greater 构造小于
func (builder *Builder) Less(column string, value any) *Builder {
	builder.toWhere(column, LESS, 1, AND)
	builder.addArg(value)
	return builder
}

// Greater 构造小于等于
func (builder *Builder) LessEqual(column string, value any) *Builder {
	builder.toWhere(column, LESSEQUAL, 1, AND)
	builder.addArg(value)
	return builder
}

// OrNotEqual 构造或者不等于
func (builder *Builder) OrNotEqual(column string, value any) *Builder {
	builder.toWhere(column, NOTEQUAL, 1, OR)
	builder.addArg(value)
	return builder
}

// Between 构造Between
func (builder *Builder) Between(column string, value1 any, value2 any) *Builder {
	builder.toWhere(column, BETWEEN, 2, AND)
	builder.addArg(value1, value2)
	return builder
}

// OrBetween 构造 或者 Between
func (builder *Builder) OrBetween(column string, value1 any, value2 any) *Builder {
	builder.toWhere(column, BETWEEN, 2, OR)
	builder.addArg(value1, value2)
	return builder
}

// NotBetween 构造不Not Between
func (builder *Builder) NotBetween(column string, value1 any, value2 any) *Builder {
	builder.toWhere(column, NOTBETWEEN, 2, AND)
	builder.addArg(value1, value2)
	return builder
}

// NotOrBetween 构造 Not Between  OR Not Between
func (builder *Builder) NotOrBetween(column string, value1 any, value2 any) *Builder {
	builder.toWhere(column, NOTBETWEEN, 2, OR)
	builder.addArg(value1, value2)
	return builder
}

// In 构造 in语句
func (builder *Builder) In(column string, value ...any) *Builder {
	builder.toWhere(column, IN, int64(len(value)), AND)
	builder.addArg(value...)
	return builder
}

// OrIn orin语句
func (builder *Builder) OrIn(column string, value ...any) *Builder {
	builder.toWhere(column, IN, int64(len(value)), OR)
	builder.addArg(value...)
	return builder
}

// NotIn .
func (builder *Builder) NotIn(column string, value ...any) *Builder {
	builder.toWhere(column, NOTIN, int64(len(value)), AND)
	builder.addArg(value...)
	return builder
}

// OrNotIn .
func (builder *Builder) OrNotIn(column string, value ...any) *Builder {
	builder.toWhere(column, NOTIN, int64(len(value)), OR)
	builder.addArg(value...)
	return builder
}

// IsNULL .
func (builder *Builder) IsNULL(column string) *Builder {
	builder.toWhere(column, ISNULL, 0, AND)
	return builder
}

// OrIsNULL .
func (builder *Builder) OrIsNULL(column string) *Builder {
	builder.toWhere(column, ISNULL, 0, OR)
	return builder
}

// IsNotNULL .
func (builder *Builder) IsNotNULL(column string) *Builder {
	builder.toWhere(column, ISNOTNULL, 0, AND)
	return builder
}

// OrIsNotNULL .
func (builder *Builder) OrIsNotNULL(column string) *Builder {
	builder.toWhere(column, ISNOTNULL, 0, OR)
	return builder
}

// Like .
func (builder *Builder) Like(column string, value any) *Builder {
	builder.toWhere(column, LIKE, 1, AND)
	builder.addArg(value)
	return builder
}

// OrLike .
func (builder *Builder) OrLike(column string, value any) *Builder {
	builder.toWhere(column, LIKE, 1, OR)
	builder.addArg(value)
	return builder
}

// Join .
func (builder *Builder) Join(tablename string, on string) *Builder {
	builder.joins = append(builder.joins, join{table: tablename, on: on, operator: JOIN})
	return builder
}

// InnerJoin ...
func (builder *Builder) InnerJoin(tablename string, on string) *Builder {
	builder.joins = append(builder.joins, join{table: tablename, on: on, operator: INNERJOIN})
	return builder
}

// LeftJoin .
func (builder *Builder) LeftJoin(tablename string, on string) *Builder {
	builder.joins = append(builder.joins, join{table: tablename, on: on, operator: LEFTJOIN})
	return builder
}

// RightJoin .
func (builder *Builder) RightJoin(tablename string, on string) *Builder {
	builder.joins = append(builder.joins, join{table: tablename, on: on, operator: RIGHTJOIN})
	return builder
}

// Union .
func (builder *Builder) Union(unions ...Builder) *Builder {
	for i, len := 0, len(unions); i < len; i++ {
		builder.unions = append(builder.unions, union{builder: unions[i], operator: UNION})
		builder.addArg(unions[i].args...)
	}
	return builder
}

// UnionOffset .
func (builder *Builder) UnionOffset(offset int64) *Builder {
	builder.unOffset = offset
	return builder
}

// UnionLimit .
func (builder *Builder) UnionLimit(limit int64) *Builder {
	builder.unLimit = limit
	return builder
}

// UnionOrderBy .
func (builder *Builder) UnionOrderBy(column string, direction string) *Builder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	builder.unOrders = append(builder.unOrders, column)
	return builder
}

// UnionAll .
func (builder *Builder) UnionAll(unions ...Builder) *Builder {
	for i, len := 0, len(unions); i < len; i++ {
		builder.unions = append(builder.unions, union{builder: unions[i], operator: UNIONALL})
		builder.addArg(unions[i].args...)
	}
	return builder
}

// Distinct .
func (builder *Builder) Distinct() *Builder {
	builder.distinct = true
	return builder
}

// GroupBy .
func (builder *Builder) GroupBy(groups ...string) *Builder {
	builder.groups = groups
	return builder
}

// OrderBy .
func (builder *Builder) OrderBy(column string, direction string) *Builder {
	if strings.ToUpper(direction) == DESC {
		column += " " + DESC
	} else {
		column += " " + ASC
	}
	builder.orders = append(builder.orders, column)
	return builder
}

// Offset .
func (builder *Builder) Offset(offset int64) *Builder {
	builder.offset = offset
	return builder
}

// Skip .
func (builder *Builder) Skip(offset int64) *Builder {
	builder.offset = offset
	return builder
}

// Limit .
func (builder *Builder) Limit(limit int64) *Builder {
	builder.limit = limit
	return builder
}

// ToSQL 输出SQL语句
func (builder *Builder) ToSQL(method string) string {
	grammar := Grammar{builder: builder, method: method}
	return grammar.ToSQL()
}
func (builder *Builder) toWhere(column string, operator string, valueNum int64, do string) *Builder {
	builder.where = append(
		builder.where,
		w{column: column, operator: operator, valueNum: valueNum, do: do})
	return builder
}
func (builder *Builder) addArg(value ...any) {
	builder.args = append(builder.args, value...)
}

func (builder *Builder) beforeArg(value ...any) {
	builder.whereArgs = append(builder.whereArgs, value...)
}

func (builder *Builder) setData(data ...map[string]any) {
	builder.data = data
}

func (builder *Builder) getInsertMap(data any) (columns []string, values map[string][]any, err error) {
	stValue := reflect.Indirect(reflect.ValueOf(data))

	values = make(map[string][]any, 0)
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
				case null.Bool:
					ignore = true
				case null.Byte:
					ignore = true
				case null.Bytes:
					ignore = true
				case null.CTime:
					ignore = true
				case null.Date:
					ignore = true
				case null.DateTime:
					ignore = true
				case null.Float32:
					ignore = true
				case null.Float64:
					ignore = true
				case null.Int:
					ignore = true
				case null.Int8:
					ignore = true
				case null.Int16:
					ignore = true
				case null.Int32:
					ignore = true
				case null.Int64:
					ignore = true
				case null.JSON:
					ignore = true
				case null.String:
					ignore = true
				case null.Time:
					ignore = true
				case null.TimeStamp:
					ignore = true
				case null.Uint:
					ignore = true
				case null.Uint8:
					ignore = true
				case null.Uint16:
					ignore = true
				case null.Uint32:
					ignore = true
				case null.Uint64:
					ignore = true
				}

				if !ignore {
					cols, vals, err := builder.getInsertMap(v.Interface())
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
			if column != "" && !builder.IsZero(v) {
				if _, ok := values[column]; ok {
					values[column] = append(values[column], v.Interface())
				} else {
					columns = append(columns, column)
					values[column] = []any{v.Interface()}
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
				values[column] = []any{stValue.MapIndex(k).Interface()}
			}
		}
	case reflect.Slice:
		n := stValue.Len()
		for i := 0; i < n; i++ {

			item := stValue.Index(i)
			cols, vals, err := builder.getInsertMap(item.Interface())

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
func (builder *Builder) IsZero(v reflect.Value) bool {
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
func (builder *Builder) MultiInsert(data ...any) (sql string, args []any, err error) {
	stVal := reflect.ValueOf(data)
	if stVal.Kind() != reflect.Slice {
		err = errors.New("data is not []any type")
		return
	}
	n := stVal.Len()
	if n < 1 {
		err = errors.New("insert data cannot be empty")
		return
	}
	columns, values, errBuild := builder.getInsertMap(data)
	if errBuild != nil {
		err = errBuild
		return
	}
	bindingsArr := make([]map[string]any, n)
	for i := 0; i < n; i++ {
		bindings := make(map[string]any, 0)
		for _, column := range columns {
			bindings[column] = values[column][i]
		}
		bindingsArr[i] = bindings
	}
	builder.setData(bindingsArr...)
	grammar := Grammar{builder: builder}
	sql = grammar.Insert()
	if len(builder.columns) < 1 {
		err = errors.New("insert data cannot be empty")
		return
	}
	args, err = builder.args, nil
	return
}

// Replace 替换
func (builder *Builder) Replace(data ...any) (sql string, args []any, err error) {

	stVal := reflect.ValueOf(data)
	if stVal.Kind() != reflect.Slice {
		err = errors.New("data is not []any type")
		return
	}
	n := stVal.Len()
	if n < 1 {
		err = errors.New("insert data cannot be empty")
		return
	}
	columns, values, errBuild := builder.getInsertMap(data)
	if errBuild != nil {
		err = errBuild
		return
	}
	bindingsArr := make([]map[string]any, n)
	for i := 0; i < n; i++ {
		bindings := make(map[string]any, 0)
		for _, column := range columns {
			bindings[column] = values[column][i]
		}
		bindingsArr[i] = bindings
	}
	builder.setData(bindingsArr...)
	grammar := Grammar{builder: builder}
	sql = grammar.Replace()
	if len(builder.columns) < 1 {
		err = errors.New("insert data cannot be empty")
		return
	}
	args, err = builder.args, nil
	return
}

// InsertUpdate ...
func (builder *Builder) InsertUpdate(insert any, update any) (sql string, args []any, err error) {
	columns, values, err := builder.getInsertMap(insert)
	if err != nil {
		return
	}
	bindingsInsert := map[string]any{}
	for _, column := range columns {
		bindingsInsert[column] = values[column][0]
	}
	columns, values, errBuild := builder.getInsertMap(update)
	if errBuild != nil {
		err = errBuild
		return
	}
	bindingsUpdate := map[string]any{}
	for _, column := range columns {
		bindingsUpdate[column] = values[column][0]
	}
	builder.setData(bindingsInsert, bindingsUpdate)
	grammar := Grammar{builder: builder}
	sql = grammar.InsertUpdate()
	args, err = builder.args, nil
	return
}

// Insert 插入数据
func (builder *Builder) Insert(data any) (sql string, args []any, err error) {
	columns, values, errBuild := builder.getInsertMap(data)
	if errBuild != nil {
		err = errBuild
		return
	}
	bindings := map[string]any{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	builder.setData(bindings)
	grammar := Grammar{builder: builder}
	sql = grammar.Insert()
	args, err = builder.args, nil
	return
}

// Update 更新
func (builder *Builder) Update(data any) (sql string, args []any, err error) {
	columns, values, errBuild := builder.getInsertMap(data)
	if errBuild != nil {
		err = errBuild
		return
	}
	bindings := map[string]any{}
	for _, column := range columns {
		bindings[column] = values[column][0]
	}
	builder.setData(bindings)
	grammar := Grammar{builder: builder}
	sql = grammar.Update()
	args = append(builder.whereArgs, builder.args...)
	err = nil
	return
}

// Delete .
func (builder *Builder) Delete() (sql string, args []any, err error) {
	grammar := Grammar{builder: builder}
	sql = grammar.Delete()
	args, err = builder.args, nil
	return
}

// Count ...
func (builder *Builder) Count() (sql string, args []any, err error) {
	builder.Select("COUNT(1) AS _C")
	grammar := Grammar{builder: builder}
	sql = grammar.Delete()
	args, err = builder.args, nil
	return
}

// Rows
func (builder *Builder) Rows() (sql string, args []any, err error) {
	grammar := Grammar{builder: builder}
	sql = grammar.Select()
	args, err = builder.args, nil
	return
}

// Row
func (builder *Builder) Row() (sql string, args []any, err error) {
	builder.offset = 0
	builder.limit = 1
	grammar := Grammar{builder: builder}
	sql = grammar.Select()
	args, err = builder.args, nil
	return
}

// // QueryRowsSQL ...
// func (builder *Builder) QueryRowsSQL(sql string, args ...any) string {
// 	builder.connection.LastSQL(sql, args...)
// 	return builder.connection.SQLRaw()
// }

// // QueryRowSQL ...
// func (builder *Builder) QueryRowSQL(sql string, args ...any) string {
// 	builder.connection.LastSQL(sql, args...)
// 	return builder.connection.SQLRaw()
// }

// // QueryRow ...
// func (builder *Builder) QueryRow(sql string, args ...any) *Row {
// 	rs := builder.QueryRows(sql, args...)
// 	r := new(Row)
// 	r.rs = rs
// 	return r
// }

// // Row 获取一条记录
// func (builder *Builder) Row() *Row {
// 	builder.offset = 0
// 	builder.limit = 1
// 	rs := builder.Rows()
// 	r := new(Row)
// 	r.rs = rs
// 	return r
// }

// // RowSQL ...
// func (builder *Builder) RowSQL() string {
// 	grammar := Grammar{builder: builder}
// 	sql := grammar.Select()

// 	builder.connection.LastSQL(sql, builder.args...)
// 	return builder.connection.SQLRaw()
// }

// // RowsSQL ...
// func (builder *Builder) RowsSQL() string {
// 	grammar := Grammar{builder: builder}
// 	sql := grammar.Select()

// 	builder.connection.LastSQL(sql, builder.args...)
// 	return builder.connection.SQLRaw()
// }

// // Rows 获取多条记录
// func (builder *Builder) Rows() *Rows {
// 	grammar := Grammar{builder: builder}
// 	sql := grammar.Select()
// 	rows, err := builder.connection.Query(builder.ctx, sql, builder.args...)
// 	if err != nil {
// 		return &Rows{rs: nil, lastError: err}
// 	}
// 	return &Rows{rs: rows, lastError: err}
// }
