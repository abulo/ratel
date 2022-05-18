package mysql2struct

var DataTypeMap = map[string][]string{
	//整型
	"TINYINT":   {"int", "sql.NullInt64"},
	"SMALLINT":  {"int", "sql.NullInt64"},
	"MEDIUMINT": {"int", "sql.NullInt64"},
	"INT":       {"int", "sql.NullInt64"},
	"INTEGER":   {"int", "sql.NullInt64"},
	"BIGINT":    {"int", "sql.NullInt64"},
	//浮点数
	"FLOAT":   {"float64", "sql.NullFloat64"},
	"DOUBLE":  {"float64", "sql.NullFloat64"},
	"DECIMAL": {"float64", "sql.NullFloat64"},
	//时间
	"DATE":      {"time.Time", "sql.NullTime"},
	"TIME":      {"time.Time", "sql.NullTime"},
	"YEAR":      {"time.Time", "sql.NullTime"},
	"DATETIME":  {"time.Time", "sql.NullTime"},
	"TIMESTAMP": {"time.Time", "sql.NullTime"},
	//字符串
	"CHAR":       {"string", "sql.NullString"},
	"VARCHAR":    {"string", "sql.NullString"},
	"TINYBLOB":   {"string", "sql.NullString"},
	"TINYTEXT":   {"string", "sql.NullString"},
	"BLOB":       {"string", "sql.NullString"},
	"TEXT":       {"string", "sql.NullString"},
	"MEDIUMBLOB": {"string", "sql.NullString"},
	"MEDIUMTEXT": {"string", "sql.NullString"},
	"LONGBLOB":   {"string", "sql.NullString"},
	"LONGTEXT":   {"string", "sql.NullString"},
	"JSON":       {"string", "sql.NullString"},
}
