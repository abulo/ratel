package dao

import (
	"context"

	"github.com/abulo/ratel/v3/stores/query"
)

type Table struct {
	TableName    string `db:"TABLE_NAME"`
	TableComment string `db:"TABLE_COMMENT"`
}

// Column ...
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
}

// DataTypeMap ...
var DataTypeMap = map[string][]string{
	//整型
	"TINYINT":   {"int64", "query.NullInt64"},
	"SMALLINT":  {"int64", "query.NullInt64"},
	"MEDIUMINT": {"int64", "query.NullInt64"},
	"INT":       {"int64", "query.NullInt64"},
	"INTEGER":   {"int64", "query.NullInt64"},
	"BIGINT":    {"int64", "query.NullInt64"},
	//浮点数
	"FLOAT":   {"float64", "query.NullFloat64"},
	"DOUBLE":  {"float64", "query.NullFloat64"},
	"DECIMAL": {"float64", "query.NullFloat64"},
	//时间

	"DATE":      {"query.NullDate", "query.NullDate"},
	"TIME":      {"query.NullTime", "query.NullTime"},
	"YEAR":      {"query.NullYear", "query.NullYear"},
	"DATETIME":  {"query.NullDateTime", "query.NullDateTime"},
	"TIMESTAMP": {"query.NullTimeStamp", "query.NullTimeStamp"},

	//字符串
	"CHAR":       {"string", "query.NullString"},
	"VARCHAR":    {"string", "query.NullString"},
	"TINYBLOB":   {"string", "query.NullString"},
	"TINYTEXT":   {"string", "query.NullString"},
	"BLOB":       {"string", "query.NullString"},
	"TEXT":       {"string", "query.NullString"},
	"MEDIUMBLOB": {"string", "query.NullString"},
	"MEDIUMTEXT": {"string", "query.NullString"},
	"LONGBLOB":   {"string", "query.NullString"},
	"LONGTEXT":   {"string", "query.NullString"},
	"JSON":       {"string", "query.NullString"},
}

// QueryTable 获取数据中表的信息
func QueryTable(ctx context.Context, DbName string) ([]Table, error) {
	var res []Table
	builder := Link.NewBuilder(ctx).Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName)
	err := builder.Rows().ToStruct(&res)
	return res, err
}

// QueryColumn 获取数据中表中字段的信息
func QueryColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	builder := Link.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Rows().ToStruct(&res)
	return res, err
}
