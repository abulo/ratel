package sql

import (
	"github.com/go-sql-driver/mysql"
)

const (
	duplicateEntryCode uint16 = 1062
)

// NewMysql returns a mysql connection.
func (c *Client) NewMysql(opts ...SqlOption) SqlConn {
	opts = append(opts, withMysqlAcceptable())
	poolOpt := &pool{
		MaxLifetime:    c.MaxLifetime,
		MaxIdleTime:    c.MaxIdleTime,
		MaxOpenConns:   c.MaxOpenConns,
		MaxIdleConns:   c.MaxIdleConns,
		DisableMetric:  c.DisableMetric,
		DisableTrace:   c.DisableTrace,
		DisablePrepare: c.DisablePrepare,
		DriverName:     c.DriverName,
		DbName:         c.Database,
		Addr:           c.Host,
	}
	return NewSqlConn(c.DriverName, c.dns(), poolOpt, opts...)
}

func mysqlAcceptable(err error) bool {
	if err == nil {
		return true
	}
	myerr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	switch myerr.Number {
	case duplicateEntryCode:
		return true
	default:
		return false
	}
}

func withMysqlAcceptable() SqlOption {
	return func(conn *commonSqlConn) {
		conn.accept = mysqlAcceptable
	}
}
