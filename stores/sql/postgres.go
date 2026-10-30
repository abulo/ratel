package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgreSQLDsn 依据 Client 当前配置拼装 PostgreSQL 连接字符串，按 SslMode 决定 sslmode 取值，
// 并附加 TimeZone 参数。
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

// PostgreSQLDialector 构造 PostgreSQL 的 GORM 方言实例；DSN 构造失败时返回对应错误。
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
