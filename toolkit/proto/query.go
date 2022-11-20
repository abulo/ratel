package proto

import (
	"context"
	"strings"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)

// Column ...
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
	ProtoType     string
	Seq           int64
}

type Table struct {
	TableName    string `db:"TABLE_NAME"`
	TableComment string `db:"TABLE_COMMENT"`
	Mark         string
}

type Index struct {
	IndexName string `db:"INDEX_NAME"`
	Field     string `db:"FIELD"`
}

type Proto struct {
	Table        Table
	Column       []Column
	PackageName  string
	Mark         string
	FunctionList []Function
}

// type ModuleArg struct {
// 	PackageName  string
// 	TableName    string
// 	Mark         string
// 	FunctionList []Function
// }

type Function struct {
	Type           string
	Name           string
	Argument       []Argument
	ArgumentNumber int
	TableName      string
	Mark           string
	Default        bool
}

type Argument struct {
	Field         string
	FieldType     string
	FieldInput    string
	ProtoType     string
	Seq           int64
	ColumnComment string
}

// DataTypeMap ...
var DataTypeMap = map[string]string{
	"json":       "string",
	"char":       "string",
	"varchar":    "string",
	"text":       "string",
	"longtext":   "string",
	"mediumtext": "string",
	"tinytext":   "string",
	"blob":       "bytes",
	"mediumblob": "bytes",
	"longblob":   "bytes",
	"varbinary":  "bytes",
	"binary":     "bytes",
	"date":       "int64",
	"time":       "int64",
	"datetime":   "int64",
	"timestamp":  "int64",
	"bool":       "bool",
	"bit":        "bool",
	"tinyint":    "int64",
	"smallint":   "int64",
	"int":        "int64",
	"mediumint":  "int64",
	"bigint":     "int64",
	"float":      "double",
	"decimal":    "double",
	"double":     "double",
}

// QueryColumn 获取数据中表中字段的信息
func QueryColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	builder := Link.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Rows().ToStruct(&res)
	if err == nil {
		for key, item := range res {
			res[key].ProtoType = ProtoType(item.DataType)
			newKey := key + 1
			res[key].Seq = cast.ToInt64(newKey)
		}
	}
	return res, err
}

// QueryTable 获取数据中表的信息
func QueryTable(ctx context.Context, DbName string, TableName string) (Table, error) {
	var res Table
	builder := Link.NewBuilder(ctx).Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName)
	err := builder.Row().ToStruct(&res)
	if err == nil {
		res.Mark = CamelStr(TableName)
	}
	return res, err
}

func QueryIndex(ctx context.Context, DbName, TableName string) ([]Index, error) {
	var res []Index
	err := Link.NewBuilder(ctx).Select("statistics.INDEX_NAME", "GROUP_CONCAT(CONCAT(statistics.COLUMN_NAME,':',`columns`.DATA_TYPE )) AS FIELD").Table("`information_schema`.`STATISTICS` AS statistics").LeftJoin("information_schema.`COLUMNS` AS `columns`", "statistics.COLUMN_NAME = `columns`.COLUMN_NAME").Where("statistics.TABLE_SCHEMA", DbName).Where("statistics.TABLE_NAME", TableName).Where("`columns`.TABLE_SCHEMA", DbName).Where("`columns`.TABLE_NAME", TableName).NotEqual("statistics.INDEX_NAME", "PRIMARY").GroupBy("statistics.TABLE_NAME", "statistics.INDEX_NAME").OrderBy("statistics.NON_UNIQUE", query.ASC).OrderBy("statistics.SEQ_IN_INDEX", query.ASC).Rows().ToStruct(&res)
	return res, err
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}

func Helper(name string) string {
	name = CamelStr(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

func ProtoType(dataType string) string {
	value, ok := DataTypeMap[dataType]
	if ok {
		return value
	}
	return "string"
}

func Add(numberOne, numberTwo interface{}) int {
	return cast.ToInt(numberOne) + cast.ToInt(numberTwo)
}
