package stage

// Column ...
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
}

type Parse struct {
	Package    string
	Module     string
	Dao        string
	Table      string
	Pri        string
	Column     []Column
	View       string
	ListTotal  int64
	LayerTotal int64
	Title      string
}
