package sql

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// ClickHouseDsn 依据 Client 当前配置拼装 ClickHouse 连接字符串，附带连接与读取超时参数。
func (c *Client) ClickHouseDsn() string {
	link := "clickhouse://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Database + "?dial_timeout=" + c.DialTimeOut + "&read_timeout=" + c.ReadTimeOut
	return link
}

// ClickHouseDialector 构造 ClickHouse 的 GORM 方言实例；DSN 构造失败时返回对应错误。
func (c *Client) ClickHouseDialector() (gorm.Dialector, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	return clickhouse.Open(dsn), nil
}
