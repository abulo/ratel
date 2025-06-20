package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySql Dsn
func (c *Client) MySqlDsn() string {
	link := c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?charset=" + c.Charset + "&loc=" + c.TimeZone
	if c.ParseTime {
		link = link + "&parseTime=true"
	} else {
		link = link + "&parseTime=false"
	}
	return link
}

// mysql open
func (c *Client) MySqlDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return mysql.Open(dsn), nil
}
