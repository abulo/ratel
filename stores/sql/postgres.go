package sql

// NewPostgres  returns a postgres connection.
func (c *Client) NewPostgres(opts ...SqlOption) SqlConn {
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
