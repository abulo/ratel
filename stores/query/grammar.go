package query

import (
	"strconv"
	"strings"
	"unsafe"
)

// Grammar sql 语法
type Grammar struct {
	query  *Query
	method string
}

func (g Grammar) compileSelect() string {
	if len(g.query.columns) < 1 {
		return "*"
	}
	return strings.Join(g.query.columns, ",")
}
func (g Grammar) compileTable(from bool) string {
	if len(g.query.table) < 1 {
		return ""
	}
	if from {
		return " FROM " + strings.Join(g.query.table, ",")
	}
	return strings.Join(g.query.table, ",")

}
func (g Grammar) compileOrder(isUnion bool) string {
	orders := g.query.orders
	if isUnion {
		orders = g.query.unOrders
	}
	if len(orders) < 1 {
		return ""
	}
	return " ORDER BY " + strings.Join(orders, ",")
}

func (g Grammar) compileGroup() string {
	if len(g.query.groups) < 1 {
		return ""
	}
	return " GROUP BY " + strings.Join(g.query.groups, ",")
}

func (g Grammar) compileLimit(isUnion bool) string {
	limit := g.query.limit
	offset := g.query.offset
	if isUnion {
		limit = g.query.unLimit
		offset = g.query.offset
	}
	if limit > 0 {
		return " LIMIT " + strconv.FormatInt(offset, 10) + "," + strconv.FormatInt(limit, 10)
	}
	return ""

}

func (g Grammar) compileDistinct() string {
	if g.query.distinct {
		return " DISTINCT "
	}
	return ""
}
func (g Grammar) compileWhere() string {
	len := len(g.query.where)
	if len < 1 {
		return ""
	}
	w := g.query.where
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
				int64Num := w[i].valueNum - 1
				intNum := *(*int)(unsafe.Pointer(&int64Num))
				sql += " " + w[i].operator + "(?" + strings.Repeat(",?", intNum) + ")"
			case ISNULL, ISNOTNULL:
				sql += " " + w[i].operator
				// break
			default:
				sql += " " + w[i].operator + " ?"
			}
		}

	}
	return sql
}
func (g Grammar) compileJoin() string {
	len := len(g.query.joins)
	if len < 1 {
		return ""
	}
	sql := ""
	joins := g.query.joins
	for i := 0; i < len; i++ {
		sql += " " + joins[i].operator + " " + joins[i].table + " ON " + joins[i].on
	}
	return sql
}
func (g Grammar) compileUnion() string {
	len := len(g.query.unions)
	if len < 1 {
		return ""
	}
	sql := ""
	unions := g.query.unions
	var g1 Grammar

	for i := 0; i < len; i++ {
		g1.query = &unions[i].query
		sql += " " + unions[i].operator
		sql += " (" + g1.Select() + ")"
	}

	return sql
}

// Select 构造select
func (g Grammar) Select() string {
	s1, s2 := "", ""
	if len(g.query.unions) > 0 {
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

// Insert ...
func (g Grammar) Insert() string {
	sql := "INSERT INTO "
	sql += g.compileTable(false)
	sql += " " + g.compileInsertValue()
	return sql
}

// Replace ...
func (g Grammar) Replace() string {
	sql := "REPLACE INTO "
	sql += g.compileTable(false)
	sql += g.compileInsertValue()
	return sql
}
func (g Grammar) compileInsertValue() string {
	sql := " ("
	for k, v := range g.query.data {
		for kv := range v {
			if k == 0 { //取第一列
				g.query.columns = append(g.query.columns, kv)
			}
		}
	}
	columns := g.query.columns
	columnsLen := len(g.query.columns)
	for index := 0; index < len(g.query.data); index++ {
		d := g.query.data[index]
		for i := 0; i < columnsLen; i++ {
			field := columns[i]
			g.query.addArg(d[field])
		}
	}
	sql += strings.Join(g.query.columns, ",")
	colLen := len(g.query.columns)
	sql += ") VALUES (?" + strings.Repeat(",?", colLen-1) + ")"
	len := len(g.query.data)
	if len > 1 {
		for i := 1; i < len; i++ {
			sql += " ,(?" + strings.Repeat(",?", colLen-1) + ")"
		}
	}
	return sql
}

// Delete ...
func (g Grammar) Delete() string {
	sql := "DELETE "
	sql += g.compileTable(true)
	sql += g.compileWhere()
	sql += g.compileOrder(false)
	if g.query.limit > 0 {
		sql += " LIMIT " + strconv.FormatInt(g.query.limit, 10)
	}
	return sql
}
func (g Grammar) compileUpdateValue() string {
	sql := ""
	data := g.query.data[0] //取一个
	for k, v := range data {
		switch vv := v.(type) {
		case Epr:
			sql += k + " = " + vv.ToString() + ","
		default:
			sql += k + " = ?,"
			g.query.beforeArg(vv)
		}
	}
	sql = strings.Trim(sql, ",")
	return sql
}

// Update ...
func (g Grammar) Update() string {
	sql := "UPDATE "
	sql += g.compileTable(false)
	sql += " SET "
	sql += g.compileUpdateValue()
	// sql += strings.Join(g.query.columns, " = ?,") + " = ?"
	sql += g.compileWhere()
	sql += g.compileOrder(false)
	if g.query.limit > 0 {
		sql += " LIMIT " + strconv.FormatInt(g.query.limit, 10)
	}
	return sql
}

// InsertUpdate ...
func (g Grammar) InsertUpdate() string {
	old := g.query.data
	//insert
	g.query.data = old[:1]
	sql := "INSERT INTO "
	sql += g.compileTable(false)
	sql += " " + g.compileInsertValue()
	sql += " ON DUPLICATE KEY UPDATE "
	g.query.data = old[1:]
	sql += g.compileUpdateValue()
	return sql
}

// ToSQL ...
func (g Grammar) ToSQL() string {
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
