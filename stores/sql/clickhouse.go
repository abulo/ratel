package sql

// NewClickhouse returns a clickhouse connection.
func (c *Client) NewClickhouse(opts ...SqlOption) SqlConn {
	poolOpt := &pool{
		MaxLifetime:    c.MaxLifetime,
		MaxIdleTime:    c.MaxIdleTime,
		MaxOpenConns:   c.MaxOpenConns,
		MaxIdleConns:   c.MaxIdleConns,
		DisableMetric:  c.DisableMetric,
		DisableTrace:   c.DisableTrace,
		DisablePrepare: c.DisablePrepare,
		DriverName:     c.DriverName,
	}
	return NewSqlConn(c.DriverName, c.dns(), poolOpt, opts...)
}
