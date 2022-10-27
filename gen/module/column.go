package module

// Index ...
type Index struct {
	NonUnique  int64  `db:"NON_UNIQUE"`
	SeqInIndex int64  `db:"SEQ_IN_INDEX"`
	IndexName  string `db:"INDEX_NAME"`
	IndexType  string `db:"INDEX_TYPE"`
	ColumnName string `db:"COLUMN_NAME"`
}

type Parse struct {
	Package   string
	Dao       string
	TableName string
	CondiTion []CondiTion
	Func      []Func
}

type Func struct {
	FuncName  string
	CondiTion []CondiTion
	NonUnique int64
}

type CondiTion struct {
	ItemType string
	ItemName string
}

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
	"TINYINT":   {"int64", "query.NullInt64", "inactive"},
	"SMALLINT":  {"int64", "query.NullInt64", "inactive"},
	"MEDIUMINT": {"int64", "query.NullInt64", "inactive"},
	"INT":       {"int64", "query.NullInt64", "inactive"},
	"INTEGER":   {"int64", "query.NullInt64", "inactive"},
	"BIGINT":    {"int64", "query.NullInt64", "inactive"},
	//浮点数
	"FLOAT":   {"float64", "query.NullFloat64", "inactive"},
	"DOUBLE":  {"float64", "query.NullFloat64", "inactive"},
	"DECIMAL": {"float64", "query.NullFloat64", "inactive"},
	//时间
	"DATE":      {"time", "query.NullDate", "inactive"},
	"TIME":      {"time", "query.NullTime", "inactive"},
	"YEAR":      {"time", "query.NullYear", "inactive"},
	"DATETIME":  {"time", "query.NullDateTime", "inactive"},
	"TIMESTAMP": {"time", "query.NullTime", "inactive"},
	//字符串
	"CHAR":       {"string", "query.NullString", "active"},
	"VARCHAR":    {"string", "query.NullString", "active"},
	"TINYBLOB":   {"string", "query.NullString", "inactive"},
	"TINYTEXT":   {"string", "query.NullString", "inactive"},
	"BLOB":       {"string", "query.NullString", "inactive"},
	"TEXT":       {"string", "query.NullString", "inactive"},
	"MEDIUMBLOB": {"string", "query.NullString", "inactive"},
	"MEDIUMTEXT": {"string", "query.NullString", "inactive"},
	"LONGBLOB":   {"string", "query.NullString", "inactive"},
	"LONGTEXT":   {"string", "query.NullString", "inactive"},
	"JSON":       {"string", "query.NullString", "inactive"},
}
