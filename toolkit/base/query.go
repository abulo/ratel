package base

import (
	"context"
)

// DaoParam 承载生成 DAO 代码所需的表信息与字段列表，零值不可直接使用，须先填充 Table 与 TableColumn。
type DaoParam struct {
	Table       Table
	TableColumn []Column
}

// Table 描述数据库中一张表的元信息。
type Table struct {
	// TableName 为表名。
	TableName string `gorm:"column:TABLE_NAME"`
	// TableComment 为表注释。
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

// Column 描述表中一个字段的元信息及其在各目标语言下的类型映射。
type Column struct {
	// ColumnName 为字段名。
	ColumnName string `gorm:"column:COLUMN_NAME"`
	// IsNullable 表示字段是否允许为 NULL，取值 YES 或 NO。
	IsNullable string `gorm:"column:IS_NULLABLE"`
	// DataType 为归一化后的数据库字段类型。
	DataType string `gorm:"column:DATA_TYPE"`
	// ColumnKey 为索引键类型，主键为 PRI，唯一索引为 UNIQUE。
	ColumnKey string `gorm:"column:COLUMN_KEY"`
	// ColumnComment 为字段注释。
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`
	// PosiTion 为字段在生成代码中的位置序号。
	PosiTion int64 `gorm:"-"`
	// DataTypeMap 为该字段对应的多语言类型映射。
	DataTypeMap DataType `gorm:"-"`
	// AlisaColumnName 为字段别名，主键字段固定为 id。
	AlisaColumnName string `gorm:"-"`
	// CamelName 为 ColumnName 的大驼峰形式，元数据加载阶段预计算一次，供各渲染 helper 复用（PERF-8-4 [P1]）。
	CamelName string `gorm:"-"`
	// LowerName 为 ColumnName 的小驼峰形式，元数据加载阶段预计算一次，供各渲染 helper 复用（PERF-8-4 [P1]）。
	LowerName string `gorm:"-"`
}

// DataType 描述单个数据库类型在 Go、Proto、TypeScript、JSON 各场景下的类型表示。
type DataType struct {
	// Default 为非空场景使用的 Go 类型。
	Default string
	// Empty 为允许为空场景使用的 null 包类型。
	Empty string
	// Proto 为 gRPC Proto 类型名。
	Proto string
	// OptionProto 表示是否以 optional 形式生成 Proto 字段。
	OptionProto bool
	// TypeScript 为 TypeScript 类型名。
	TypeScript string
	// Json 为 JSON 序列化时使用的默认值。
	Json string
}

// Index 描述表上一个索引的元信息。
type Index struct {
	// IndexName 为索引名称。
	IndexName string `gorm:"column:INDEX_NAME"`
	// Field 为索引覆盖的字段，多个字段以逗号分隔。
	Field string `gorm:"column:FIELD"`
}

// ModuleParam 承载生成一个模块所需的表、方法与路径上下文，由命令行交互流程填充。
type ModuleParam struct {
	// Table 为表信息。
	Table Table
	// TableColumn 为表字段列表。
	TableColumn []Column
	// Method 为待生成的方法列表。
	Method []Method
	// Pkg 为生成代码的包名（路径末级）。
	Pkg string
	// PkgPath 为包在模块内的相对路径。
	PkgPath string
	// Primary 为主键字段信息。
	Primary Column
	// ModName 为当前 go.mod 的模块路径。
	ModName string
	// Page 表示方法是否需要分页参数。
	Page bool
	// SoftDelete 表示是否使用软删除。
	SoftDelete bool
	// Condition 为构造查询条件所需的字段列表。
	Condition []Column
}

// Method 描述待生成的一个数据访问方法及其绑定的表上下文。
type Method struct {
	// Type 为方法类型，取值 Create、Update、Delete、Show、Item、List、Recover、Drop。
	Type string
	// Name 为生成的函数名。
	Name string
	// Alias 为函数别名。
	Alias string
	// Condition 为该方法查询条件的字段列表。
	Condition []Column
	// ConditionTotal 为查询条件字段数量。
	ConditionTotal int
	// Table 为绑定的表信息。
	Table Table
	// TableColumn 为表结构信息。
	TableColumn []Column
	// Pkg 为包名。
	Pkg string
	// PkgPath 为包在模块内的相对路径。
	PkgPath string
	// Primary 为主键字段信息。
	Primary Column
	// ModName 为 go.mod 的模块路径。
	ModName string
	// Page 表示是否分页。
	Page bool
	// SoftDelete 表示是否软删除。
	SoftDelete bool
}

// SqlResultConvert 记录需要归一化的数据库类型映射，键为原始类型，值为归一化后的类型。
var SqlResultConvert = map[string]string{
	"int8": "bigint",
	"int4": "int",
}

// sharedDataTypeMap 为全流程共享的只读类型映射表；类型集合在运行期间恒定，包级构造一次，
// 避免每张表重复重建 33 项映射（PERF-8-8 [P2]，NewDataType 返回结果按原约定不可被调用方修改）。
var sharedDataTypeMap = NewDataType()

// NewDataType 返回数据库类型到各目标语言类型映射的完整对照表，调用方不应修改返回结果。
func NewDataType() map[string]DataType {
	res := make(map[string]DataType)
	res["numeric"] = DataType{Default: "float64", Empty: "null.Float64", Proto: "double", OptionProto: true, TypeScript: "number", Json: "0"}
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
	res["timestamptz"] = DataType{Default: "null.Time", Empty: "null.Time", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["date"] = DataType{Default: "null.Date", Empty: "null.Date", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["datetime"] = DataType{Default: "null.DateTime", Empty: "null.DateTime", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["timestamp"] = DataType{Default: "null.TimeStamp", Empty: "null.TimeStamp", Proto: "google.protobuf.Timestamp", OptionProto: false, TypeScript: "string", Json: ""}
	res["year"] = DataType{Default: "int32", Empty: "null.Int32", Proto: "int32", OptionProto: true, TypeScript: "number", Json: "0"}

	res["bit"] = DataType{Default: "null.Bytes", Empty: "null.Bytes", Proto: "bytes", OptionProto: true, TypeScript: "string", Json: ""}
	res["boolean"] = DataType{Default: "bool", Empty: "null.Bool", Proto: "bool", OptionProto: true, TypeScript: "boolean", Json: "false"}

	return res
}

// TableList 查询指定数据库下全部数据表的元信息。
//
// 依据当前 Driver 选择 MySQL 或 PostgreSQL 的系统表执行查询；SQL 执行失败时返回对应错误。
func TableList(ctx context.Context, DbName string) ([]Table, error) {
	var res []Table
	var query string
	var args []any
	var err error
	switch *Driver {
	case "mysql":
		query = "SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA=?"
		args = []any{DbName}
	case "postgres":
		query = "SELECT C.relname AS \"TABLE_NAME\",obj_description(C.OID) AS \"TABLE_COMMENT\" FROM pg_catalog.pg_class C WHERE C.relkind='r' AND C.relnamespace=(SELECT OID FROM pg_catalog.pg_namespace WHERE nspname='public')"
	}
	handle := Query.Raw(query, args...)
	handle.Scan(&res)
	err = handle.Error
	return res, err
}

// TableItem 查询指定数据库中某一张表的元信息；表不存在时返回零值 Table 与查询错误。
func TableItem(ctx context.Context, DbName, TableName string) (Table, error) {
	var res Table
	var query string
	var args []any
	var err error
	switch *Driver {
	case "mysql":
		query = "SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA=? AND TABLE_NAME=?"
		args = []any{DbName, TableName}
	case "postgres":
		query, args = "SELECT C.relname AS \"TABLE_NAME\",obj_description(C.OID) AS \"TABLE_COMMENT\" FROM pg_catalog.pg_class C WHERE C.relkind='r' AND C.relnamespace=(SELECT OID FROM pg_catalog.pg_namespace WHERE nspname='public') AND C.relname=?", []any{TableName}
	}
	handle := Query.Raw(query, args...)
	handle.Scan(&res)
	err = handle.Error
	return res, err
}

// TableColumn 查询指定表的全部字段信息，并补全类型映射与字段位置序号；SQL 执行失败时返回错误。
func TableColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	var query string
	var args []any
	var err error
	switch *Driver {
	case "mysql":
		query = "SELECT COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY ORDINAL_POSITION ASC"
		args = []any{DbName, TableName}
	case "postgres":
		query, args = "SELECT C.COLUMN_NAME AS \"COLUMN_NAME\",C.is_nullable AS \"IS_NULLABLE\",C.udt_name AS \"DATA_TYPE\",CASE WHEN EXISTS (SELECT 1 FROM pg_constraint con JOIN pg_attribute A ON A.attnum=con.conkey[1] AND A.attrelid=con.conrelid WHERE con.contype='p' AND con.conrelid=cl.OID AND A.attname=C.COLUMN_NAME) THEN 'PRI' WHEN EXISTS (SELECT 1 FROM pg_constraint con JOIN pg_attribute A ON A.attnum=con.conkey[1] AND A.attrelid=con.conrelid WHERE con.contype='u' AND con.conrelid=cl.OID AND A.attname=C.COLUMN_NAME) THEN 'UNIQUE' ELSE 'NONE' END AS \"COLUMN_KEY\",col_description (cl.OID,A.attnum) AS \"COLUMN_COMMENT\" FROM information_schema.COLUMNS C LEFT JOIN pg_catalog.pg_class cl ON cl.relname=C.TABLE_NAME LEFT JOIN pg_catalog.pg_attribute A ON A.attrelid=cl.OID AND A.attnum=C.ordinal_position WHERE C.table_schema='public' AND C.TABLE_NAME=? order by C.ordinal_position ASC;", []any{TableName}
	}

	handle := Query.Raw(query, args...)
	handle.Scan(&res)
	err = handle.Error
	if err == nil {
		// PERF-8-9 [P2] 类型归一化、类型映射回填、位置序号与驼峰预计算合并为单趟循环；
		// PERF-8-8 [P2] 复用包级只读映射表；PERF-8-4 [P1] 在此一次性预计算驼峰名称供渲染 helper 复用。
		for k := range res {
			if val, ok := SqlResultConvert[res[k].DataType]; ok {
				res[k].DataType = val
			}
			res[k].DataTypeMap = sharedDataTypeMap[res[k].DataType]
			res[k].PosiTion = int64(k + 1)
			res[k].CamelName = CamelStr(res[k].ColumnName)
			res[k].LowerName = Helper(res[k].ColumnName)
		}
	}
	return res, err
}

// TableIndex 查询指定表的索引信息，主键索引不在结果内；SQL 执行失败时返回错误。
func TableIndex(ctx context.Context, DbName, TableName string) ([]Index, error) {
	var res []Index
	var query string
	var args []any
	var err error
	switch *Driver {
	case "mysql":
		query = "SELECT statistics.INDEX_NAME,GROUP_CONCAT(CONCAT(statistics.COLUMN_NAME) ORDER BY statistics.NON_UNIQUE ASC,statistics.SEQ_IN_INDEX ASC) AS FIELD FROM `information_schema`.`STATISTICS` AS statistics LEFT JOIN `information_schema`.`COLUMNS` AS `columns` ON statistics.COLUMN_NAME = `columns`.COLUMN_NAME WHERE statistics.TABLE_SCHEMA=? AND statistics.TABLE_NAME=? AND `columns`.TABLE_SCHEMA=? AND `columns`.TABLE_NAME=? AND statistics.INDEX_NAME!='PRIMARY' GROUP BY statistics.TABLE_NAME,statistics.INDEX_NAME ORDER BY statistics.NON_UNIQUE ASC,statistics.SEQ_IN_INDEX ASC"
		args = []any{DbName, TableName, DbName, TableName}
	case "postgres":
		query, args = "WITH index_info AS (SELECT i.relname AS index_name,ARRAY_AGG(A.attname ORDER BY ia.attnum) AS fields,ci.indisprimary FROM pg_class T JOIN pg_index ci ON T.OID=ci.indrelid JOIN pg_class i ON ci.indexrelid=i.OID JOIN UNNEST(ci.indkey) WITH ORDINALITY AS ia (attnum,ORDINALITY) ON ia.attnum> 0 JOIN pg_attribute A ON A.attnum=ia.attnum AND A.attrelid=T.OID WHERE T.relname=? AND T.relkind='r' GROUP BY i.relname,ci.indisprimary) SELECT index_name AS \"INDEX_NAME\",array_to_string(fields,',') AS \"FIELD\" FROM index_info WHERE NOT indisprimary ORDER BY index_name;", []any{TableName}
	}
	handle := Query.Raw(query, args...)
	handle.Scan(&res)
	err = handle.Error
	return res, err
}

// TablePrimary 查询指定表的主键字段信息，并将字段别名固定为 id；表没有主键或查询失败时返回错误。
func TablePrimary(ctx context.Context, DbName, TableName string) (Column, error) {
	var res Column
	var query string
	var args []any
	var err error
	switch *Driver {
	case "mysql":
		query = "SELECT COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME=? AND COLUMN_KEY='PRI' ORDER BY ORDINAL_POSITION ASC"
		args = []any{DbName, TableName}
	case "postgres":
		query, args = "SELECT * FROM (SELECT C.COLUMN_NAME AS \"COLUMN_NAME\",C.is_nullable AS \"IS_NULLABLE\",C.udt_name AS \"DATA_TYPE\",CASE WHEN EXISTS (SELECT 1 FROM pg_constraint con JOIN pg_attribute A ON A.attnum=con.conkey[1] AND A.attrelid=con.conrelid WHERE con.contype='p' AND con.conrelid=cl.OID AND A.attname=C.COLUMN_NAME) THEN 'PRI' WHEN EXISTS (SELECT 1 FROM pg_constraint con JOIN pg_attribute A ON A.attnum=con.conkey[1] AND A.attrelid=con.conrelid WHERE con.contype='u' AND con.conrelid=cl.OID AND A.attname=C.COLUMN_NAME) THEN 'UNIQUE' ELSE 'NONE' END AS \"COLUMN_KEY\",col_description (cl.OID,A.attnum) AS \"COLUMN_COMMENT\" FROM information_schema.COLUMNS C LEFT JOIN pg_catalog.pg_class cl ON cl.relname=C.TABLE_NAME LEFT JOIN pg_catalog.pg_attribute A ON A.attrelid=cl.OID AND A.attnum=C.ordinal_position WHERE C.table_schema='public' AND C.TABLE_NAME=?) AS TT WHERE TT.\"COLUMN_KEY\"='PRI'", []any{TableName}
	}
	handle := Query.Raw(query, args...)
	handle.Scan(&res)
	err = handle.Error
	if err == nil {
		if val, ok := SqlResultConvert[res.DataType]; ok {
			res.DataType = val
		}
		// PERF-8-8 [P2] 复用包级只读映射表；PERF-8-4 [P1] 预计算驼峰名称。
		res.DataTypeMap = sharedDataTypeMap[res.DataType]
		res.PosiTion = 1
		res.AlisaColumnName = "id"
		res.CamelName = CamelStr(res.ColumnName)
		res.LowerName = Helper(res.ColumnName)
	}
	return res, err
}
