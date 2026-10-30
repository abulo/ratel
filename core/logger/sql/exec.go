package sql

import (
	"encoding/json"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/null"
	"gorm.io/gorm"
)

// ExecCloser 抽象日志条目的后端写入器，由各实现负责把条目提交到对应存储。
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

// tableName 保存当前日志表名，由 NewExec 创建写入器时覆盖，供 Dao.TableName 读取。
var tableName string

type defaultExec struct {
	client   *gorm.DB
	canClose bool
}

// NewExec 基于给定的 GORM 客户端与日志表名创建写入器。
func NewExec(client *gorm.DB, loggerTableName string) ExecCloser {
	tableName = loggerTableName
	return &defaultExec{
		client:   client,
		canClose: true,
	}
}

// NewExecWithURL 基于给定的 GORM 客户端与日志表名创建写入器，当前实现与 NewExec 行为一致。
func NewExecWithURL(client *gorm.DB, loggerTableName string) ExecCloser {
	tableName = loggerTableName
	return &defaultExec{
		client:   client,
		canClose: true,
	}
}

func (e *defaultExec) Exec(entry *entry.Entry) error {
	daoItem := &Dao{}
	daoItem.Host = null.StringFrom(entry.Host)
	daoItem.File = null.StringFrom(entry.File)
	daoItem.Func = null.StringFrom(entry.Func)
	daoItem.Message = null.StringFrom(entry.Message)
	daoItem.Level = null.StringFrom(entry.Level)
	daoItem.Timestamp = null.TimeStampFrom(entry.Timestamp)
	data, _ := json.Marshal(entry.Data)
	daoItem.Data = null.JSONFrom(data)
	return e.client.Create(daoItem).Error
}

// Dao 是日志条目的 GORM 持久化模型，表名由 TableName 返回的全局 tableName 决定。
type Dao struct {
	// Id 为自增主键。
	Id        *int64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	// Host 为产生日志的主机名。
	Host      null.String    `gorm:"column:host" json:"host"`
	// Timestamp 为日志产生时间。
	Timestamp null.TimeStamp `gorm:"column:timestamp" json:"timestamp"`
	// File 为调用位置，格式为 文件路径:行号。
	File      null.String    `gorm:"column:file" json:"file"`
	// Func 为产生日志的函数全限定名。
	Func      null.String    `gorm:"column:func" json:"func"`
	// Message 为日志正文。
	Message   null.String    `gorm:"column:message" json:"message"`
	// Level 为大写形式的日志级别名称。
	Level     null.String    `gorm:"column:level" json:"level"`
	// Data 为日志附加字段的 JSON 序列化结果。
	Data      null.JSON      `gorm:"column:data" json:"data"`
}

// TableName 返回日志模型对应的数据库表名。
func (Dao) TableName() string {
	return tableName
}
