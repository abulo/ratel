package mysql

import (
	"strconv"
	"strings"
	"unsafe"
)

//Grammar sql 语法
type Grammar struct {
	builder *QueryBuilder
	method  string
}

func (g Grammar) compileSelect() string {
	if len(g.builder.columns) < 1 {
		return "*"
	}
	return strings.Join(g.builder.columns, ",")
}
func (g Grammar) compileTable(from bool) string {
	if len(g.builder.table) < 1 {
		return ""
	}
	if from {
		return " FROM " + strings.Join(g.builder.table, ",")
	} else {
		return strings.Join(g.builder.table, ",")
	}

}
func (g Grammar) compileOrder(isUnion bool) string {
	orders := g.builder.orders
	if isUnion {
		orders = g.builder.unOrders
	}
	if len(orders) < 1 {
		return ""
	}
	return " ORDER BY " + strings.Join(orders, ",")
}

func (g Grammar) compileGroup() string {
	if len(g.builder.groups) < 1 {
		return ""
	}
	return " GROUP BY " + strings.Join(g.builder.groups, ",")
}

func (g Grammar) compileLimit(isUnion bool) string {
	limit := g.builder.limit
	offset := g.builder.offset
	if isUnion {
		limit = g.builder.unLimit
		offset = g.builder.offset
	}
	if limit > 0 {
		return " LIMIT " + strconv.FormatInt(offset, 10) + "," + strconv.FormatInt(limit, 10)
	} else {
		return ""
	}
}

func (g Grammar) compileDistinct() string {
	if g.builder.distinct {
		return " DISTINCT "
	}
	return ""
}
func (g Grammar) compileWhere() string {
	len := len(g.builder.where)
	if len < 1 {
		return ""
	}
	w := g.builder.where
	sql := " WHERE "
	for i := 0; i < len; i++ {
		if i > 0 {
			sql += " " + w[i].do
		}
		sql += " " + w[i].column
		if w[i].operator != "" {
			switch w[i].operator {
			case BETWEEN, NOTBETWEEN:
				sql += " " + w[i].operator + " ? AND ?"
			case IN, NOTIN:
				int64_num := w[i].valuenum - 1
				int_num := *(*int)(unsafe.Pointer(&int64_num))
				sql += " " + w[i].operator + "(?" + strings.Repeat(",?", int_num) + ")"
			case ISNULL, ISNOTNULL:
				sql += " " + w[i].operator
				break
			default:
				sql += " " + w[i].operator + " ?"
			}
		}

	}
	return sql
}
func (g Grammar) compileJoin() string {
	len := len(g.builder.joins)
	if len < 1 {
		return ""
	}
	sql := ""
	joins := g.builder.joins
	for i := 0; i < len; i++ {
		sql += " " + joins[i].operator + " " + joins[i].table + " ON " + joins[i].on
	}
	return sql
}
func (g Grammar) compileUnion() string {
	len := len(g.builder.unions)
	if len < 1 {
		return ""
	}
	sql := ""
	unions := g.builder.unions
	var g1 Grammar

	for i := 0; i < len; i++ {
		g1.builder = &unions[i].query

		sql += " " + unions[i].operator
		sql += " (" + g1.Select() + ")"
	}

	return sql
}

//Select 构造select
func (g Grammar) Select() string {
	s1, s2 := "", ""
	if len(g.builder.unions) > 0 {
		s1 = "("
		s2 = ")"
	}
	sql := s1 + "SELECT "
	sql += g.compileDistinct()
	sql += g.compileSelect()
	sql += g.compileTable(true)
	sql += g.compileJoin()
	sql += g.compileWhere()
	sql += g.compileGroup()
	sql += g.compileOrder(false)
	sql += g.compileLimit(false)
	sql += s2
	sql += g.compileUnion()
	sql += g.compileLimit(true)
	sql += g.compileOrder(true)

	return sql
}
func (g Grammar) Insert() string {
	sql := "INSERT INTO "
	sql += g.compileTable(false)
	sql += " " + g.compileInsertValue()
	return sql
}
func (g Grammar) Replace() string {
	sql := "REPLACE INTO "
	sql += g.compileTable(false)
	sql += g.compileInsertValue()
	return sql
}
func (g Grammar) compileInsertValue() string {
	sql := " ("
	for k, v := range g.builder.data {
		for kv, _ := range v {
			if k == 0 { //取第一列
				g.builder.columns = append(g.builder.columns, kv)
			}
		}
	}
	columns := g.builder.columns
	columnsLen := len(g.builder.columns)
	for index := 0; index < len(g.builder.data); index++ {
		d := g.builder.data[index]
		for i := 0; i < columnsLen; i++ {
			field := columns[i]
			g.builder.addArg(d[field])
		}
	}
	sql += strings.Join(g.builder.columns, ",")
	collen := len(g.builder.columns)
	sql += ") VALUES (?" + strings.Repeat(",?", collen-1) + ")"
	len := len(g.builder.data)
	if len > 1 {
		for i := 1; i < len; i++ {
			sql += " ,(?" + strings.Repeat(",?", collen-1) + ")"
		}
	}
	return sql
}
func (g Grammar) Delete() string {
	sql := "DELETE "
	sql += g.compileTable(true)
	sql += g.compileWhere()
	sql += g.compileOrder(false)
	if g.builder.limit > 0 {
		sql += " LIMIT " + strconv.FormatInt(g.builder.limit, 10)
	}
	return sql
}
func (g Grammar) compileUpdateValue() string {
	sql := ""
	data := g.builder.data[0] //取一个
	for k, v := range data {
		switch vv := v.(type) {
		case Epr:
			sql += k + " = " + vv.ToString() + ","
		default:
			sql += k + " = ?,"
			g.builder.beforeArg(vv)
		}
	}
	sql = strings.Trim(sql, ",")
	return sql
}
func (g Grammar) Update() string {
	sql := "UPDATE "
	sql += g.compileTable(false)
	sql += " SET "
	sql += g.compileUpdateValue()
	// sql += strings.Join(g.builder.columns, " = ?,") + " = ?"
	sql += g.compileWhere()
	sql += g.compileOrder(false)
	if g.builder.limit > 0 {
		sql += " LIMIT " + strconv.FormatInt(g.builder.limit, 10)
	}
	return sql
}
func (g Grammar) InsertUpdate() string {
	old := g.builder.data
	//insert
	g.builder.data = old[:1]
	sql := "INSERT INTO "
	sql += g.compileTable(false)
	sql += " " + g.compileInsertValue()
	sql += " ON DUPLICATE KEY UPDATE "
	g.builder.data = old[1:]
	sql += g.compileUpdateValue()
	return sql
}
func (g Grammar) ToSql() string {
	g.method = strings.ToUpper(g.method)
	switch g.method {
	case "INSERT":
		return g.Insert()
	case "DELETE":
		return g.Delete()
	case "UPDATE":
		return g.Update()
	case "REPLACE":
		return g.Replace()
	case "INSERTUPDATE":
		return g.Replace()
	default:
		return g.Select()
	}
}
