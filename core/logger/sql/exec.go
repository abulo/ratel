package sql

import (
	"context"
	"encoding/json"

	"github.com/abulo/ratel/v3/core/logger/entry"
	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/stores/sql"
)

// ExecCloser 将logrus条目写入数据库并关闭数据库
type ExecCloser interface {
	Exec(entry *entry.Entry) error
}

type defaultExec struct {
	client    sql.SqlConn
	tableName string
	canClose  bool
}

// NewExec create an exec instance
func NewExec(client sql.SqlConn, tableName string) ExecCloser {
	return &defaultExec{
		client:    client,
		tableName: tableName,
		canClose:  true,
	}
}

// NewExecWithURL create an exec instance
func NewExecWithURL(client sql.SqlConn, tableName string) ExecCloser {
	return &defaultExec{
		client:    client,
		tableName: tableName,
		canClose:  true,
	}
}

func (e *defaultExec) Exec(entry *entry.Entry) error {
	daoItem := &Dao{}
	daoItem.Host = null.StringFrom(entry.Host)
	daoItem.File = null.StringFrom(entry.File)
	daoItem.Func = null.StringFrom(entry.Func)
	daoItem.Message = null.StringFrom(entry.Message)
	daoItem.Level = null.StringFrom(entry.Level)
	data, _ := json.Marshal(entry.Data)
	daoItem.Data = null.JSONFrom(data)
	builder := sql.NewBuilder()
	query, args, err := builder.Table(e.tableName).Insert(data)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = e.client.Insert(ctx, query, args...)
	return err
}

type Dao struct {
	Id        *int64        `db:"id,-" json:"id"`
	Host      null.String   `db:"host" json:"host"`
	Timestamp null.DateTime `db:"timestamp" json:"timestamp"`
	File      null.String   `db:"file" json:"file"`
	Func      null.String   `db:"func" json:"func"`
	Message   null.String   `db:"message" json:"message"`
	Level     null.String   `db:"level" json:"level"`
	Data      null.JSON     `db:"data" json:"data"`
}
