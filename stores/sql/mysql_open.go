package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func mysqlOpen(driverName, dns string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dns)
	if err != nil {
		return nil, errors.Wrapf(err, "open mysql connection failed, dns: %s", dns)
	}
	return db, nil
}
