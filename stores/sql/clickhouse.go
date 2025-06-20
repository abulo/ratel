package sql

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// clickhouse Dsn
func (c *Client) ClickHouseDsn() string {
	link := "clickhouse://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Database + "?dial_timeout=" + c.DialTimeOut + "&read_timeout=" + c.ReadTimeOut
	return link
}

// clickhouse open
func (c *Client) ClickHouseDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return clickhouse.Open(dsn), nil
}
