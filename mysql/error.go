package mysql

type DbError struct {
	msg string
	sql Sql
}

type Epr struct {
	value string
}

func NewEpr(value string) Epr {
	return Epr{value: value}
}
func (e Epr) ToString() string {
	return e.value
}

func NewDBError(msg string, sql Sql) DbError {
	return DbError{msg: msg, sql: sql}
}

func (e DbError) Error() string {
	return "DBError:" + e.msg + " " + e.sql.ToJson()
}
