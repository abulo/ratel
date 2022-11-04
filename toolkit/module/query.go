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
