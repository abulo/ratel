package module

// Column ...
type Column struct {
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
	CondiTion []string
	Func      []Func
}

type Func struct {
	FuncName  string
	CondiTion []string
	NonUnique int64
}
