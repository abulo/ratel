package base

import (
	"context"

	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/spf13/cast"
)

type DaoParam struct {
	Table       Table
	TableColumn []Column
}

// Table 表信息
type Table struct {
	TableName    string `db:"TABLE_NAME"`    // 表名
	TableComment string `db:"TABLE_COMMENT"` // 表注释
}

// Column 字段新
type Column struct {
	ColumnName      string   `db:"COLUMN_NAME"`    // 字段名
	IsNullable      string   `db:"IS_NULLABLE"`    // 是否为空
	DataType        string   `db:"DATA_TYPE"`      // 字段类型
	ColumnKey       string   `db:"COLUMN_KEY"`     // 是否索引
	ColumnComment   string   `db:"COLUMN_COMMENT"` // 字段描述
	PosiTion        int64    // 排序信息
	DataTypeMap     DataType // 字段类型信息
	AlisaColumnName string   // 字段名
}

// DataType 字段类型信息
type DataType struct {
	Default     string // 不空时
	Empty       string // 为空时
	Proto       string // Grpc 协议
	OptionProto bool   // Grpc 协议
	TypeScript  string // TypeScript
	Json        string // Json
}

// Index 索引信息
type Index struct {
	IndexName string `db:"INDEX_NAME"` // 索引名称
	Field     string `db:"FIELD"`      // 索引作用字段
}

type ModuleParam struct {
	Table       Table
	TableColumn []Column
	Method      []Method // 方法
	Pkg         string   // 包名
	PkgPath     string   // 包名路径
	Primary     Column   // 主键信息
	ModName     string   // go.mod 信息
	Page        bool     // 是否分页
	SoftDelete  bool     // 是否删除
	Condition   []Column // 函数需要的条件信息
}

// Method 构造的函数
type Method struct {
	Type           string   // 方法类型(List多个/Item单条/Create新增/Update修改/Delete删除/Only单个)
	Name           string   // 函数名称
	Alias          string   // 函数别名
	Condition      []Column // 函数需要的条件信息
	ConditionTotal int      // 条件数量
	Table          Table    // 表信息
	TableColumn    []Column // 表结构信息
	Pkg            string   // 包名
	PkgPath        string   // 包名路径
	Primary        Column   // 主键信息
	ModName        string   // go.mod 信息
	Page           bool     // 是否分页
	SoftDelete     bool     // 是否删除
}

func NewDataType() map[string]DataType {
	res := make(map[string]DataType)
	res["numeric"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["integer"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["int"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["smallint"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["mediumint"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["tinyint"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}
	res["bigint"] = DataType{Default: "int64", Empty: "null.Int64", Proto: "int64", OptionProto: true, TypeScript: "number", Json: "0"}

	res["float"] = DataType{Default: "float32", Empty: "null.Float32", Proto: "float", OptionProto: true, TypeScript: "number", Json: "0"}
	res["real"] = DataType{Default: "float64", Empty: "null.Float64", Proto: "double", OptionProto: true, TypeScript: "number", Json: "0"}
	res["double"] = DataType{Default: "float64", Empty: "null.Float64", Proto: "double", OptionProto: true, TypeScript: "number", Json: "0"}
	res["decimal"] = DataType{Default: "float64", Empty: "null.Float64", Proto: "double", OptionProto: true, TypeScript: "number", Json: "0"}

	res["char"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}
	res["varchar"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}
	res["tinytext"] = DataType{Default: "string", Empty: "null.String", Proto: "string", TypeScript: "string", Json: ""}
	res["mediumtext"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}
	res["longtext"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}
	res["text"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}
	res["json"] = DataType{Default: "null.JSON", Empty: "null.JSON", Proto: "bytes", OptionProto: true, TypeScript: "any", Json: "{}"}
	res["enum"] = DataType{Default: "string", Empty: "null.String", Proto: "string", OptionProto: true, TypeScript: "string", Json: ""}

	res["binary"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["varbinary"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["tinyblob"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["blob"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["mediumblob"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["longblob"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}

	res["time"] = DataType{Default: "null.CTime", Empty: "null.CTime", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["date"] = DataType{Default: "null.Date", Empty: "null.Date", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["datetime"] = DataType{Default: "null.DateTime", Empty: "null.DateTime", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["timestamp"] = DataType{Default: "null.TimeStamp", Empty: "null.TimeStamp", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["year"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}

	res["bit"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["boolean"] = DataType{Default: "bool", Empty: "null.Bool", Proto: "bool", OptionProto: true, TypeScript: "boolean", Json: "false"}

	return res
}

// TableList 获取数据中表的信息
func TableList(ctx context.Context, DbName string) ([]Table, error) {
	var res []Table
	builder := sql.NewBuilder()
	query, args, err := builder.Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName).Rows()
	if err != nil {
		return res, err
	}
	err = Query.QueryRows(ctx, query, args...).ToStruct(&res)
	return res, err
}

// TableItem 获取数据中表的信息
func TableItem(ctx context.Context, DbName, TableName string) (Table, error) {
	var res Table
	builder := sql.NewBuilder()
	query, args, err := builder.Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).Row()
	if err != nil {
		return res, err
	}
	err = Query.QueryRow(ctx, query, args...).ToStruct(&res)
	return res, err
}

// TableColumn 获取数据中表中字段的信息
func TableColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	builder := sql.NewBuilder()
	query, args, err := builder.Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).OrderBy("ORDINAL_POSITION", sql.ASC).Rows()
	if err != nil {
		return res, err
	}
	err = Query.QueryRows(ctx, query, args...).ToStruct(&res)
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

// TableIndex 获取表的索引信息
func TableIndex(ctx context.Context, DbName, TableName string) ([]Index, error) {
	var res []Index
	builder := sql.NewBuilder()
	query, args, err := builder.Select("statistics.INDEX_NAME", "GROUP_CONCAT(CONCAT(statistics.COLUMN_NAME) ORDER BY statistics.NON_UNIQUE ASC,statistics.SEQ_IN_INDEX ASC) AS FIELD").Table("`information_schema`.`STATISTICS` AS statistics").LeftJoin("information_schema.`COLUMNS` AS `columns`", "statistics.COLUMN_NAME = `columns`.COLUMN_NAME").Where("statistics.TABLE_SCHEMA", DbName).Where("statistics.TABLE_NAME", TableName).Where("`columns`.TABLE_SCHEMA", DbName).Where("`columns`.TABLE_NAME", TableName).NotEqual("statistics.INDEX_NAME", "PRIMARY").GroupBy("statistics.TABLE_NAME", "statistics.INDEX_NAME").OrderBy("statistics.NON_UNIQUE", sql.ASC).OrderBy("statistics.SEQ_IN_INDEX", sql.ASC).Rows()
	if err != nil {
		return res, err
	}
	err = Query.QueryRows(ctx, query, args...).ToStruct(&res)
	return res, err
}

// TablePrimary 获取主键
func TablePrimary(ctx context.Context, DbName, TableName string) (Column, error) {
	var res Column
	builder := sql.NewBuilder()
	query, args, err := builder.Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).Where("COLUMN_KEY", "PRI").OrderBy("ORDINAL_POSITION", sql.ASC).Row()
	if err != nil {
		return res, err
	}
	err = Query.QueryRow(ctx, query, args...).ToStruct(&res)
	if err == nil {
		dataType := NewDataType()
		res.DataTypeMap = dataType[res.DataType]
		res.PosiTion = 1
		res.AlisaColumnName = "id"
	}
	return res, err
}
