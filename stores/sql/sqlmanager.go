package sql

import (
	"database/sql"
	"io"

	"github.com/abulo/ratel/v3/core/resource"
	"github.com/pkg/errors"
)

var connManager = resource.NewResourceManager()

func getSqlConn(driverName, server string, pool *pool) (*sql.DB, error) {
	conn, err := getCachedSqlConn(driverName, server, pool)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getCachedSqlConn(driverName, server string, pool *pool) (*sql.DB, error) {
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		conn, err := newDBConnection(driverName, server, pool)
		if err != nil {
			return nil, err
		}

		return conn, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*sql.DB), nil
}

func newDBConnection(driverName, dns string, pool *pool) (*sql.DB, error) {
	conn, err := dbConnection(driverName, dns)
	if err != nil {
		return nil, err
	}
	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxOpenConns(pool.MaxOpenConns)
	conn.SetMaxIdleConns(pool.MaxIdleConns)
	conn.SetConnMaxLifetime(pool.MaxLifetime)
	conn.SetConnMaxIdleTime(pool.MaxIdleTime)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

func dbConnection(driverName, dns string) (*sql.DB, error) {
	switch driverName {
	case driverMysql:
		return mysqlOpen(driverName, dns)
	case driverClickhouse:
		return clickhouseOpen(driverName, dns)
	case driverPostgres:
		return postgresOpen(driverName, dns)
	default:
		return nil, errors.Errorf("unsupported driver: %s", driverName)
	}
}
