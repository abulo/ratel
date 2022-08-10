package mysql2struct

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
	"DATE":      {"time.Time", "query.NullDate"},
	"TIME":      {"time.Time", "query.NullTime"},
	"YEAR":      {"time.Time", "query.NullYear"},
	"DATETIME":  {"time.Time", "query.NullDateTime"},
	"TIMESTAMP": {"time.Time", "query.NullTime"},
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
