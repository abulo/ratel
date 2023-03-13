package sql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func postgresOpen(driverName, dns string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dns)
	if err != nil {
		return nil, errors.Wrapf(err, "open postgres connection failed, dns: %s", dns)
	}
	return db, nil
}
