package sql

import "github.com/abulo/ratel/v3/util"

func (c *Client) mysqlDns() string {
	return c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?charset=" + c.Charset + "&loc=" + c.TimeZone + "&parseTime=true"
}

func (c *Client) clickhouseDns() string {
	link := "clickhouse://"
	if !util.Empty(c.Username) {
		link = link + c.Username
	}
	link = link + ":"
	if !util.Empty(c.Password) {
		link = link + c.Password
	}
	link = link + "@"
	link = link + util.Implode(",", c.Addr)
	link = link + "/"
	if !util.Empty(c.Database) {
		link = link + c.Database
	}
	param := make([]string, 0)
	if !util.Empty(c.DialTimeout) {
		param = append(param, "dial_timeout="+c.DialTimeout)
	}
	if !util.Empty(c.MaxExecutionTime) {
		param = append(param, "max_execution_time="+c.MaxExecutionTime)
	}
	if !util.Empty(c.OpenStrategy) {
		param = append(param, "connection_open_strategy="+c.OpenStrategy)
	}
	if c.Compress {
		param = append(param, "compress=true")
	} else {
		param = append(param, "compress=false")
	}
	if !c.DisableDebug {
		param = append(param, "debug=true")
	} else {
		param = append(param, "debug=false")
	}
	return link + "?" + util.Implode("&", param)
}

func (c *Client) postgresDns() string {
	return c.Username + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?charset=" + c.Charset + "&loc=" + c.TimeZone + "&parseTime=true"
}
