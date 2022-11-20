package module

import (
	"context"

	"github.com/abulo/ratel/v3/stores/query"
)

type Index struct {
	IndexName string `db:"INDEX_NAME"`
	Field     string `db:"FIELD"`
}

func QueryIndex(ctx context.Context, DbName, TableName string) ([]Index, error) {
	var res []Index
	err := Link.NewBuilder(ctx).Select("statistics.INDEX_NAME", "GROUP_CONCAT(CONCAT(statistics.COLUMN_NAME,':',`columns`.DATA_TYPE )) AS FIELD").Table("`information_schema`.`STATISTICS` AS statistics").LeftJoin("information_schema.`COLUMNS` AS `columns`", "statistics.COLUMN_NAME = `columns`.COLUMN_NAME").Where("statistics.TABLE_SCHEMA", DbName).Where("statistics.TABLE_NAME", TableName).Where("`columns`.TABLE_SCHEMA", DbName).Where("`columns`.TABLE_NAME", TableName).NotEqual("statistics.INDEX_NAME", "PRIMARY").GroupBy("statistics.TABLE_NAME", "statistics.INDEX_NAME").OrderBy("statistics.NON_UNIQUE", query.ASC).OrderBy("statistics.SEQ_IN_INDEX", query.ASC).Rows().ToStruct(&res)
	return res, err
}

type ModuleArg struct {
	PackageName  string
	TableName    string
	Mark         string
	FunctionList []Function
	PrimaryKey   string
	ColumnName   string
}

type Function struct {
	Type      string
	Name      string
	Argument  []Argument
	TableName string
	Mark      string
	Default   bool
}

type Argument struct {
	Field      string
	FieldType  string
	FieldInput string
}

// Column ...
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
}

//PrimaryKey

// QueryColumn 获取数据中表中字段的信息
func QueryColumnPrimaryKey(ctx context.Context, DbName, TableName string) (Column, error) {
	var res Column
	builder := Link.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).Where("COLUMN_KEY", "PRI").OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Row().ToStruct(&res)
	return res, err
}
