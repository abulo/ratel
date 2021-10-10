package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abulo/ratel"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/logger/hook"
	"github.com/abulo/ratel/mongodb"
	"github.com/abulo/ratel/mysql"
	"github.com/abulo/ratel/redis"
	"github.com/abulo/ratel/util"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	ratel.Ratel
}

//MongoDB 代理
var MongoDB *mongodb.Proxy = mongodb.NewProxy()
var Redis *redis.Proxy = redis.NewProxy()
var MySQL *mysql.ProxyPool = mysql.NewProxyPool()

// func init() {

// }

type AdminPermission struct {
	ID         int64              `db:"id" json:"id"`
	ParentID   int64              `db:"parent_id" json:"parent_id"` //父ID
	Title      string             `db:"title" json:"title"`         // 权限名称
	Handle     string             `db:"handle" json:"handle"`       //路由别名
	Weight     int64              `db:"weight" json:"weight"`       //权重
	URI        string             `db:"url,-" json:"url"`
	CreateDate mysql.NullDateTime `db:"create_date"`
	UpdateDate mysql.NullDateTime `db:"update_date"`
}

func main() {
	mongodb.SetTrace(true)
	opt := &mongodb.Config{}
	opt.URI = "mongodb://root:654321@127.0.0.1:27017/admin_request_log?authSource=admin"
	opt.MaxConnIdleTime = 5
	opt.MaxPoolSize = 64
	opt.MinPoolSize = 10
	MongoDB.SetNameSpace("common", mongodb.New(opt))

	redis.SetTrace(true)

	optr := &redis.Config{}
	optr.KeyPrefix = "abulo"
	optr.Password = ""
	optr.PoolSize = 10
	optr.Database = 0
	optr.Hosts = []string{"127.0.0.1:6379"}
	optr.Type = false

	Redis.SetNameSpace("common", redis.New(optr))

	loggerHook := hook.DefaultWithURL(MongoDB.NameSpace("common"))
	defer loggerHook.Flush()
	logger.Logger.AddHook(loggerHook)
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Logger.SetReportCaller(true)
	logger.Logger.SetOutput(os.Stdout)
	// eng := NewEngine()

	// if err := eng.Run(); err != nil {
	// 	logger.Logger.Panic(err)
	// }

	optm := &mysql.Config{}
	optm.Username = "root"
	optm.Password = "mysql"
	optm.Host = "127.0.0.1"
	optm.Port = "3306"
	optm.Charset = "utf8mb4"
	optm.Database = "ratel"
	optm.ConnMaxLifetime = 100
	optm.ConnMaxIdleTime = 100
	optm.MaxIdleConns = 100
	optm.MaxOpenConns = 100
	MySQL = mysql.NewProxyPool()
	proxy := mysql.NewProxy()
	proxy.SetWrite(mysql.New(optm))
	MySQL.SetNameSpace("common", proxy)

	a1 := new(AdminPermission)
	a1.Title = "张三"
	a1.UpdateDate = mysql.NewNullDateTime()
	a1.CreateDate = mysql.NewDateTime(util.Now())
	a1.ParentID = 0
	a1.Handle = "abulo1"
	a1.Weight = 1
	// ParentID   int64            `db:"parent_id" json:"parent_id"` //父ID
	// Title      string           `db:"title" json:"title"`         // 权限名称
	// Handle     string           `db:"handle" json:"handle"`       //路由别名
	// Weight     int64            `db:"weight" json:"weight"`       //权重
	// URI

	db := MySQL.NameSpace("common").Write()
	ctx := context.TODO()
	sql, err := db.NewQuery(ctx).Table("admin_permission").Where("id", 487).Update(a1)
	fmt.Println(sql, err)

	// var result AdminPermission
	// err := db.NewQuery(ctx).Table("admin_permission").Where("id", 527).Row().ToStruct(&result)
	// fmt.Println(err, result)
}

// func NewEngine() *Engine {
// 	eng := &Engine{}
// 	if err := eng.Startup(
// 		eng.serveHTTP,
// 		eng.serveHTTPTwo,
// 	); err != nil {
// 		logger.Logger.Panic("startup", err)
// 	}
// 	eng.Tracer("ratel", "127.0.0.1:6831")
// 	return eng
// }
// func (eng *Engine) serveHTTP() error {
// 	config := &http.Config{
// 		Host: "127.0.0.1",
// 		Port: 7777,
// 		Mode: gin.DebugMode,
// 		Name: "admin",
// 	}
// 	server := config.Build()
// 	server.Use(trace.HTTPMetricServerInterceptor())
// 	server.Use(trace.HTTPTraceServerInterceptor())
// 	server.GET("/ping", "ping", func(ctx *gin.Context) {
// 		// e := Redis.NameSpace("common").Set(ctx.Request.Context(), "aaaaa", "daadasd", time.Minute*5).Err()
// 		ctx.JSON(200, gin.H{
// 			"status": "7777",
// 		})
// 	})

// 	return eng.Serve(server)
// }

// func (eng *Engine) serveHTTPTwo() error {
// 	config := &monitor.Config{
// 		Host:    "127.0.0.1",
// 		Port:    17777,
// 		Network: "tcp4",
// 		Name:    "monitor",
// 	}
// 	// monitor.HandleFunc("/metrics", func(w ohttp.ResponseWriter, r *ohttp.Request) {
// 	// 	promhttp.Handler().ServeHTTP(w, r)
// 	// })
// 	server := config.Build()

// 	server.HandleFunc("/metrics", func(w ohttp.ResponseWriter, r *ohttp.Request) {
// 		promhttp.Handler().ServeHTTP(w, r)
// 	})

// 	// server.Use(trace.HTTPTraceServerInterceptor())

// 	// 	server.InitFuncMap()
// 	// 	pprof.Register(server.Engine)

// 	return eng.Serve(server)
// }
