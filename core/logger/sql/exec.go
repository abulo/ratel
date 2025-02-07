package sql

import (
	"encoding/json"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/null"
	"gorm.io/gorm"
)

// ExecCloser 将logrus条目写入数据库并关闭数据库
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

var tableName string

type defaultExec struct {
	client   *gorm.DB
	canClose bool
}

// NewExec create an exec instance
func NewExec(client *gorm.DB, loggerTableName string) ExecCloser {
	tableName = loggerTableName
	return &defaultExec{
		client:   client,
		canClose: true,
	}
}

// NewExecWithURL create an exec instance
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

type Dao struct {
	Id        *int64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Host      null.String    `gorm:"column:host" json:"host"`
	Timestamp null.TimeStamp `gorm:"column:timestamp" json:"timestamp"`
	File      null.String    `gorm:"column:file" json:"file"`
	Func      null.String    `gorm:"column:func" json:"func"`
	Message   null.String    `gorm:"column:message" json:"message"`
	Level     null.String    `gorm:"column:level" json:"level"`
	Data      null.JSON      `gorm:"column:data" json:"data"`
}

// 表名重写为 tableName
func (Dao) TableName() string {
	return tableName
}
