package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abulo/ratel/v2/store/mysql"
	"github.com/abulo/ratel/v2/store/proxy"
	"github.com/abulo/ratel/v2/store/query"
	"github.com/abulo/ratel/v2/util"
)

//MongoDB 代理
var Store *proxy.ProxyPool

// func init() {

// }

type AdminPermission struct {
	ID         int64              `db:"id" json:"id"`
	ParentID   int64              `db:"parent_id" json:"parent_id"` //父ID
	Title      string             `db:"title" json:"title"`         // 权限名称
	Handle     string             `db:"handle" json:"handle"`       //路由别名
	Weight     int64              `db:"weight" json:"weight"`       //权重
	URI        string             `db:"url,-" json:"url"`
	CreateDate query.NullDateTime `db:"create_date"`
	UpdateDate query.NullDateTime `db:"update_date"`
}

func init() {
	Store = proxy.NewProxyPool()
	optm := &mysql.Config{}
	optm.Username = "root"
	optm.Password = "mysql"
	optm.Host = "127.0.0.1"
	optm.Port = "3306"
	optm.Charset = "utf8mb4"
	optm.Database = "ratel"
	optm.MaxLifetime = time.Duration(2) * time.Minute
	optm.MaxIdleTime = time.Duration(2) * time.Minute
	optm.MaxIdleConns = 100
	optm.MaxOpenConns = 100
	optm.DriverName = "mysql"
	conn := mysql.New(optm)
	proxyPool := proxy.NewProxySQL()
	proxyPool.SetWrite(conn)
	proxyPool.SetRead(conn)
	Store.StoreSQL("mysql", proxyPool)
}

func main() {

	for {
		db := Store.LoadSQL("mysql").Read()
		sql := "select * from admin_user where username='admin'"

		ddd, err := db.NewQuery(context.Background()).QueryRow(sql).ToMap()

		fmt.Println(util.Now().GoString(), ddd, err)
	}

	// optm := &mysql.Config{}
	// optm.Username = "root"
	// optm.Password = "mysql"
	// optm.Host = "127.0.0.1"
	// optm.Port = "3306"
	// optm.Charset = "utf8mb4"
	// optm.Database = "ratel"
	// optm.MaxLifetime = 100
	// optm.MaxIdleTime = 100
	// optm.MaxIdleConns = 100
	// optm.MaxOpenConns = 100
	// MySQL = mysql.NewProxyPool()
	// proxy := mysql.NewProxy()
	// proxy.SetWrite(mysql.New(optm))
	// MySQL.SetNameSpace("common", proxy)

	// a1 := new(AdminPermission)
	// a1.Title = "张三"
	// a1.UpdateDate = query.NewNullDateTime()
	// a1.CreateDate = query.NewDateTime(util.Now())
	// a1.ParentID = 0
	// a1.Handle = "abulo1"
	// a1.Weight = 1
	// // ParentID   int64            `db:"parent_id" json:"parent_id"` //父ID
	// // Title      string           `db:"title" json:"title"`         // 权限名称
	// // Handle     string           `db:"handle" json:"handle"`       //路由别名
	// // Weight     int64            `db:"weight" json:"weight"`       //权重
	// // URI

	// db := MySQL.NameSpace("common").Write()
	// ctx := context.TODO()
	// sql, err := db.NewQuery(ctx).Table("admin_permission").Where("id", 487).Update(a1)
	// fmt.Println(sql, err)

	// var result AdminPermission
	// err := db.NewQuery(ctx).Table("admin_permission").Where("id", 527).Row().ToStruct(&result)
	// fmt.Println(err, result)
}
