package create

type TableType int

const (
	STable TableType = iota + 1
	CTable
)

type ColumnType string

const (
	Timestamp        ColumnType = "TIMESTAMP"
	Int              ColumnType = "INT"
	IntUnsigned      ColumnType = "INT UNSIGNED"
	BigInt           ColumnType = "BIGINT"
	BigIntUnsigned   ColumnType = "BIGINT UNSIGNED"
	Float            ColumnType = "FLOAT"
	Double           ColumnType = "DOUBLE"
	Binary           ColumnType = "BINARY"
	SmallInt         ColumnType = "SMALLINT"
	SmallIntUnsigned ColumnType = "SMALLINT UNSIGNED"
	TinyInt          ColumnType = "TINYINT"
	TinyIntUnsigned  ColumnType = "TINYINT UNSIGNED"
	Bool             ColumnType = "BOOL"
	NChar            ColumnType = "NCHAR"
	JSON             ColumnType = "JSON"
	VarChar          ColumnType = "VARCHAR"
	Geometry         ColumnType = "GEOMETRY" // TODO: not support yet.
	VarBinary        ColumnType = "VARBINARY"
)
