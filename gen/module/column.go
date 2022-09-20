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

	"DATE":      {"time", "query.NullDate"},
	"TIME":      {"time", "query.NullTime"},
	"YEAR":      {"time", "query.NullYear"},
	"DATETIME":  {"time", "query.NullDateTime"},
	"TIMESTAMP": {"time", "query.NullTime"},

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
