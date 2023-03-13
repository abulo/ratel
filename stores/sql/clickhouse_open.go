package sql

import (
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pkg/errors"
)

func clickhouseOpen(driverName, dns string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dns)
	if err != nil {
		return nil, errors.Wrapf(err, "open clickhouse connection failed, dns: %s", dns)
	}
	return db, nil
}
