package tdengine

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Migrator struct {
	migrator.Migrator
	d Dialect
}

type Column struct {
	name              string
	nullable          sql.NullString
	datatype          string
	maxlen            sql.NullInt64
	precision         sql.NullInt64
	scale             sql.NullInt64
	datetimeprecision sql.NullInt64
}

func (c Column) Name() string {
	return c.name
}

func (c Column) DatabaseTypeName() string {
	return c.datatype
}

func (c Column) Length() (length int64, ok bool) {
	ok = c.maxlen.Valid
	if ok {
		length = c.maxlen.Int64
	} else {
		length = 0
	}
	return
}

func (c Column) Nullable() (nullable bool, ok bool) {
	if c.nullable.Valid {
		nullable, ok = c.nullable.String == "YES", true
	} else {
		nullable, ok = false, false
	}
	return
}

func (c Column) DecimalSize() (precision int64, scale int64, ok bool) {
	return
}

func (m Migrator) FullDataTypeOf(field *schema.Field) (expr clause.Expr) {
	expr.SQL = m.d.DataTypeOf(field)
	return
}

func (m Migrator) AlterColumn(value any, field string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(field); field != nil {
			return m.DB.Exec(
				"ALTER TABLE ? MODIFY COLUMN ? ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: field.DBName}, m.FullDataTypeOf(field),
			).Error
		}
		return fmt.Errorf("failed to look up field with name: %s", field)
	})
}

func (m Migrator) RenameColumn(value any, oldName, newName string) error {
	return errors.New("RenameColumn not support")
}

func (m Migrator) RenameIndex(value any, oldName, newName string) error {
	return errors.New("RenameIndex not support")
}

func (m Migrator) DropConstraint(value any, name string) error {
	return errors.New("DropConstraint not support")
}
func (m Migrator) AutoMigrate(values ...any) error {
	return errors.New("AutoMigrate not support")
}
