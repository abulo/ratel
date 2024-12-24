package sql

import (
	"database/sql"

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

func PgxOpen(driverName, dns string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dns)
	if err != nil {
		return nil, errors.Wrapf(err, "open postgres connection failed, dns: %s", dns)
	}
	return db, nil
}
