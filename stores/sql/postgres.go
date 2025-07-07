package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgres Dsn
func (c *Client) PostgreSQLDsn() string {
	link := "host=" + c.Host + " user=" + c.Username + " password=" + c.Password + " dbname=" + c.Database + " port=" + c.Port
	if c.SslMode {
		link = link + " sslmode=enable"
	} else {
		link = link + " sslmode=disable"
	}
	link = link + " TimeZone=" + c.TimeZone
	return link
}

// postgres open
func (c *Client) PostgreSQLDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return postgres.New(postgres.Config{
		DSN:              dsn,
		WithoutReturning: false,
	}), nil
}
