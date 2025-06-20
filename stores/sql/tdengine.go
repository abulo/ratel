package sql

import (
	"github.com/abulo/ratel/v3/stores/sql/plugin/tdengine"
	"gorm.io/gorm"
)

// MySql Dsn
func (c *Client) TDengineDsn() string {
	link := c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/"
	return link
}

// mysql open
func (c *Client) TDengineDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return tdengine.Dialect{DSN: dsn}, nil
}
