package tdengine

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/abulo/ratel/v3/stores/sql/plugin/tdengine/utils"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var _ gorm.Dialector = Dialect{}

// DefaultDriverName is the default driver name for TDengine.
const DefaultDriverName = "taosSql"

type Dialect struct {
	DriverName string
	DSN        string
	Conn       gorm.ConnPool
}

func (Dialect) Name() string {
	return "tdengine"
}

func (dialect Dialect) Initialize(db *gorm.DB) (err error) {
	if dialect.DriverName == "" {
		dialect.DriverName = DefaultDriverName
	}
	db.SkipDefaultTransaction = true
	db.DisableNestedTransaction = true
	db.DisableAutomaticPing = true
	db.DisableForeignKeyConstraintWhenMigrating = true
	if dialect.Conn != nil {
		db.ConnPool = dialect.Conn
	} else {
		db.ConnPool, err = sql.Open(dialect.DriverName, dialect.DSN)
		if err != nil {
			return err
		}
	}
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		LastInsertIDReversed: true,
		QueryClauses:         []string{"SELECT", "FROM", "WHERE", "WINDOW", "FILL", "GROUP BY", "ORDER BY", "SLIMIT", "LIMIT"},
		CreateClauses:        []string{"CREATE TABLE", "INSERT", "USING", "VALUES", "ON CONFLICT", "RETURNING"},
	})

	for k, v := range dialect.ClauseBuilders() {
		db.ClauseBuilders[k] = v
	}
	return nil
}

func (Dialect) ClauseBuilders() map[string]clause.ClauseBuilder {
	return map[string]clause.ClauseBuilder{
		"INSERT": func(c clause.Clause, builder clause.Builder) {
			if _, ok := c.Expression.(clause.Insert); ok {
				if stmt, ok := builder.(*gorm.Statement); ok {
					_, containsCreateTable := stmt.Clauses["CREATE TABLE"]
					if containsCreateTable {
						return
					}
				}
			}
			c.Build(builder)
		},
		"FOR": func(c clause.Clause, builder clause.Builder) {
			if _, ok := c.Expression.(clause.Locking); ok {
				return
			}
			c.Build(builder)
		},
		"VALUES": func(c clause.Clause, builder clause.Builder) {
			if _, ok := c.Expression.(clause.Values); ok {
				if stmt, ok := builder.(*gorm.Statement); ok {
					_, containsCreateTable := stmt.Clauses["CREATE TABLE"]
					if containsCreateTable {
						return
					}
				}
			}
			c.Build(builder)
		},
	}
}
func (Dialect) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "NULL"}
}

func (dialect Dialect) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialect,
		CreateIndexAfterCreateTable: false,
	}}, dialect}
}

func (Dialect) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v any) {
	writer.WriteByte('?')
}

func (Dialect) QuoteTo(writer clause.Writer, str string) {
	utils.QuoteTo(writer, str)
}

func (Dialect) Explain(sql string, vars ...any) string {
	return logger.ExplainSQL(sql, nil, "'", vars...)
}

func (Dialect) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "bool"
	case schema.Int, schema.Uint:
		constraint := func(sqlType string) string {
			if field.DataType == schema.Uint {
				sqlType += " unsigned"
			}
			return sqlType
		}
		switch {
		case field.Size <= 8:
			return constraint("tinyint")
		case field.Size <= 16:
			return constraint("smallint")
		case field.Size <= 32:
			return constraint("int")
		default:
			return constraint("bigint")
		}
	case schema.Float:
		// TODO: Decimal暂未支持
		// if field.Precision > 0 {
		// 	return fmt.Sprintf("decimal(%d, %d)", field.Precision, field.Scale)
		// }
		if field.Size <= 32 {
			return "float"
		} else {
			return "double"
		}
	case schema.String:
		size := field.Size
		if size == 0 {
			size = 64
		}
		return fmt.Sprintf("NCHAR(%d)", size)
	case schema.Time:
		return "TIMESTAMP"
	case schema.Bytes:
		size := field.Size
		if size == 0 {
			size = 64
		}
		return fmt.Sprintf("BINARY(%d)", size)
	default:
		return string(field.DataType)
	}
}

func (Dialect) SavePoint(tx *gorm.DB, name string) error {
	return errors.New("not support transaction")
}

func (Dialect) RollbackTo(tx *gorm.DB, name string) error {
	return errors.New("not support transaction")
}
