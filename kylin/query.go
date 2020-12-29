package kylin

import (
	"encoding/json"
)

//Query 查询条件封装
type Query struct {
	SQL           string `json:"sql"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	AcceptPartial bool   `json:"acceptPartial"`
	Project       string `json:"project"`
}

//GetBytes 对象转json
func (query *Query) GetBytes() (body []byte) {
	var err error
	body, err = json.Marshal(query)
	if err != nil {
		return nil
	}
	return
}

//QueryResult 查询后正常的返回结果
type QueryResult struct {
	ColumnMetas       []*Column     `json:"columnMetas"`
	Result            []interface{} `json:"results"`
	Cube              string        `json:"cube"`
	AffectedRowCount  int           `json:"affectedRowCount"`
	IsException       bool          `json:"isException"`
	ExceptionMessage  string        `json:"exceptionMessage"`
	Duration          int           `json:"duration"`
	TotalScanCount    int           `json:"totalScanCount"`
	HitExceptionCache bool          `json:"hitExceptionCache"`
	StorageCacheUsed  bool          `json:"storageCacheUsed"`
	Partial           bool          `json:"partial"`
}

//QueryOut 返回上层数据
type QueryOut struct {
	ColumnMetas []string      `json:"columnMetas"`
	Result      []interface{} `json:"results"`
}

//Column 字段
type Column struct {
	IsNullable         int    `json:"isNullable"`
	DisplaySize        int    `json:"displaySize"`
	Label              string `json:"label"`
	Name               string `json:"name"`
	SchemaName         string `json:"schemaName"`
	CatelogName        string `json:"catelogName"`
	TableName          string `json:"tableName"`
	Precision          int    `json:"precision"`
	Scale              int    `json:"scale"`
	ColumnType         int    `json:"columnType"`
	ColumnTypeName     string `json:"columnTypeName"`
	ReadOnly           bool   `json:"readOnly"`
	AutoIncrement      bool   `json:"autoIncrement"`
	CaseSensitive      bool   `json:"caseSensitive"`
	Searchable         bool   `json:"searchable"`
	Currency           bool   `json:"currency"`
	Signed             bool   `json:"signed"`
	DefinitelyWritable bool   `json:"definitelyWritable"`
	Writable           bool   `json:"writable"`
}
