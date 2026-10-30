package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlDsn 依据 Client 当前配置拼装 MySQL 连接字符串，按 ParseTime 追加 parseTime 参数。
func (c *Client) MySqlDsn() string {
	link := c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?charset=" + c.Charset + "&loc=" + c.TimeZone
	if c.ParseTime {
		link = link + "&parseTime=true"
	} else {
		link = link + "&parseTime=false"
	}
	return link
}

// MySqlDialector 构造 MySQL 的 GORM 方言实例；DSN 构造失败时返回对应错误。
func (c *Client) MySqlDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return mysql.Open(dsn), nil
}
