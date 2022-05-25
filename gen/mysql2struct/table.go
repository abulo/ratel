package mysql2struct

type Table struct {
	TableName    string `db:"TABLE_NAME"`
	TableComment string `db:"TABLE_COMMENT"`
}
