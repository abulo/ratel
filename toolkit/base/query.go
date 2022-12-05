package base

// Table 表信息
type Table struct {
	TableName    string `db:"TABLE_NAME"`    // 表名
	TableComment string `db:"TABLE_COMMENT"` // 表注释
}

// Column 字段新
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`    // 字段名
	IsNullable    string `db:"IS_NULLABLE"`    // 是否为空
	DataType      string `db:"DATA_TYPE"`      // 字段类型
	ColumnKey     string `db:"COLUMN_KEY"`     // 是否索引
	ColumnComment string `db:"COLUMN_COMMENT"` // 字段描述
}

// DataType 字段信息
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

	res["binary"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["varbinary"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["tinyblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["blob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["mediumblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["longblob"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}

	res["text"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["json"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}
	res["enum"] = DataType{Default: "string", Empty: "query.NullString", Proto: "string"}

	res["time"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["date"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["datetime"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["timestamp"] = DataType{Default: "time.Time", Empty: "query.NullTime", Proto: "google.protobuf.Timestamp"}
	res["year"] = DataType{Default: "int32", Empty: "query.NullInt32", Proto: "int32"}

	res["bit"] = DataType{Default: "[]byte", Empty: "query.NullBytes", Proto: "bytes"}
	res["boolean"] = DataType{Default: "bool", Empty: "query.NullBool", Proto: "bool"}
	// "bit":        func(string) string { return "[]uint8" },
	// "boolean":    func(string) string { return "bool" },

	return res
}
