package base

import (
	"context"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/spf13/cast"
)

// Table 表信息
type Table struct {
	TableName    string `db:"TABLE_NAME"`    // 表名
	TableComment string `db:"TABLE_COMMENT"` // 表注释
}

// Column 字段新
type Column struct {
	ColumnName    string   `db:"COLUMN_NAME"`    // 字段名
	IsNullable    string   `db:"IS_NULLABLE"`    // 是否为空
	DataType      string   `db:"DATA_TYPE"`      // 字段类型
	ColumnKey     string   `db:"COLUMN_KEY"`     // 是否索引
	ColumnComment string   `db:"COLUMN_COMMENT"` // 字段描述
	PosiTion      int64    // 排序信息
	DataTypeMap   DataType //字段类型信息
}

// DataType 字段类型信息
type DataType struct {
	Default string // 不空时
	Empty   string // 为空时
	Proto   string // Grpc 协议
}

func NewDataType() map[string]DataType {
	res := make(map[string]DataType)
	res["numeric"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["integer"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["int"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["smallint"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["mediumint"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["tinyint"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}
	res["bigint"] = DataType{Default: "int64", Empty: "query.NullInt64", Proto: "int64"}

	res["float"] = DataType{Default: "float32", Empty: "query.NullFloat32", Proto: "float"}
	res["real"] = DataType{Default: "float64", Empty: "query.NullFloat64", Proto: "double"}
	res["double"] = DataType{Default: "float64", Empty: "query.NullFloat64", Proto: "double"}
	res["decimal"] = DataType{Default: "float64", Empty: "query.NullFloat64", Proto: "double"}

	res["char"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["varchar"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["tinytext"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["mediumtext"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["longtext"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["text"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["json"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["enum"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}

	// res["binary"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["varbinary"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["tinyblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["blob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["mediumblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["longblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}

	res["binary"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["varbinary"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["tinyblob"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["blob"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["mediumblob"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["longblob"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}

	// res["time"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	// res["date"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	// res["datetime"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	// res["timestamp"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}

	res["time"] = DataType{Default: "query.NullTime", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["date"] = DataType{Default: "query.NullTime", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["datetime"] = DataType{Default: "query.NullTime", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["timestamp"] = DataType{Default: "query.NullTime", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["year"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}

	// res["bit"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	// res["boolean"] = DataType{Default: "bool", Empty: "query.NullBool", Proto: "bool"}

	res["bit"] = DataType{Default: "query.NullBytes", Empty: "query.NullBytes", Proto: "bytes"}
	res["boolean"] = DataType{Default: "query.NullBool", Empty: "query.NullBool", Proto: "bool"}

	return res
}

// TableList 获取数据中表的信息
func TableList(ctx context.Context, DbName string) ([]Table, error) {
	var res []Table
	builder := Query.NewBuilder(ctx).Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName)
	err := builder.Rows().ToStruct(&res)
	return res, err
}

// TableColumn 获取数据中表中字段的信息
func TableColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	builder := Query.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Rows().ToStruct(&res)
	if err == nil {
		dataType := NewDataType()
		for key, item := range res {
			res[key].DataTypeMap = dataType[item.DataType]
			newKey := key + 1
			res[key].PosiTion = cast.ToInt64(newKey)
		}
	}
	return res, err
}
