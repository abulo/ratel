package main

import (
	"context"
	"fmt"
	"os"

	cl "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/abulo/ratel/v2"
	"github.com/abulo/ratel/v2/logger"
	"github.com/abulo/ratel/v2/logger/mongo"
	"github.com/abulo/ratel/v2/store/clickhouse"
	"github.com/abulo/ratel/v2/store/elasticsearch"
	"github.com/abulo/ratel/v2/store/mongodb"
	"github.com/abulo/ratel/v2/store/mysql"
	"github.com/abulo/ratel/v2/store/query"
	"github.com/abulo/ratel/v2/store/redis"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	ratel.Ratel
}

//MongoDB 代理
var MongoDB *mongodb.Proxy = mongodb.NewProxy()
var Redis *redis.Proxy = redis.NewProxy()
var MySQL *mysql.ProxyPool = mysql.NewProxyPool()
var ClickHouse *clickhouse.ProxyPool = clickhouse.NewProxyPool()
var Elastic *elasticsearch.Proxy = elasticsearch.NewProxy()

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

func main() {
	// mongodb.SetTrace(true)
	opt := &mongodb.Config{}
	opt.URI = "mongodb://root:654321@127.0.0.1:27017/admin_request_log?authSource=admin"
	opt.MaxConnIdleTime = 5
	opt.MaxPoolSize = 64
	opt.MinPoolSize = 10
	MongoDB.SetNameSpace("common", mongodb.New(opt))

	// esOpt := []elastic.ClientOptionFunc{}

	// urls := make([]string, 0)
	// urls = append(urls, "http://127.0.0.1:9200")
	// esOpt = append(esOpt, elastic.SetURL(urls...))
	// esOpt = append(esOpt, elastic.SetSniff(false))
	// Elastic.SetNameSpace("common", elasticsearch.NewClient(esOpt...))

	optr := &redis.Config{}
	optr.KeyPrefix = ""
	optr.Password = ""
	optr.PoolSize = 10
	optr.Database = 1
	optr.Hosts = []string{"127.0.0.1:6379"}
	optr.Type = false

	Redis.SetNameSpace("common", redis.New(optr))

	// loggerHook := es.DefaultWithURL(Elastic.NameSpace("common"))
	loggerHook := mongo.DefaultWithURL(MongoDB.NameSpace("common"))
	defer loggerHook.Flush()
	logger.Logger.AddHook(loggerHook)
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Logger.SetReportCaller(true)
	logger.Logger.SetOutput(os.Stdout)
	// eng := NewEngine()

	// logger.Logger.Info("adasdasd")

	// logger.Logger.WithFields(logrus.Fields{
	// 	"animal": "walrus",
	// }).Info("A walrus appears")

	// if err := eng.Run(); err != nil {
	// 	logger.Logger.Panic(err)
	// }

	// MaxOpenConns     int           //连接池最多同时打开的连接数
	// MaxIdleConns     int           //连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxLifetime      time.Duration //连接池里面的连接最大存活时长
	// MaxIdleTime      time.Duration //连接池里面的连接最大空闲时长
	// DriverName       string
	// Trace            bool

	addr := make([]string, 0)
	addr = append(addr, "127.0.0.1:9000")
	optm := &clickhouse.Config{}
	optm.Username = "zeus"
	optm.Password = "zeus"
	optm.Addr = addr
	optm.Database = "xmt"
	optm.MaxIdleConns = 100
	optm.MaxOpenConns = 100
	optm.DialTimeout = "3s"
	optm.OpenStrategy = "random"
	optm.DriverName = "clickhouse"
	optm.MaxExecutionTime = "60"
	optm.Compress = true
	optm.Debug = true
	optm.Prepare = false
	ClickHouse = clickhouse.NewProxyPool()
	proxy := clickhouse.NewProxy()
	proxy.SetWrite(clickhouse.New(optm))
	ClickHouse.SetNameSpace("common", proxy)

	// queueName := "account:queue:login"
	// ctx := context.Background()

	ctx := cl.Context(context.Background())

	// redisHandel := Redis.NameSpace("common")
	// data := make([]interface{}, 0)
	// for i := 0; i < 5; i++ {
	// 	if result, err := redisHandel.LPop(ctx, queueName).Result(); err == nil {
	// 		tmp := make(map[string]interface{})
	// 		fmt.Println(result)
	// 		if err = json.Unmarshal([]byte(result), &tmp); err == nil {
	// 			data = append(data, tmp)
	// 		}
	// 	}
	// }

	clickHouseHandel := ClickHouse.NameSpace("common").Write()

	where := make([]interface{}, 0)
	where = append(where, cl.Named("Col1", "xmt"))
	where = append(where, cl.Named("Col2", "account_register"))

	// where = append(where, "xmt")
	// where = append(where, "account_register")
	data, err := clickHouseHandel.NewQuery(ctx).QueryRow("SELECT * FROM information_schema.tables WHERE table_schema = @Col1 AND table_name = @Col2", where...).ToMap()
	fmt.Println(data, err)

	// txClickhouse, err := clickHouseHandel.Begin()
	// if err != nil {
	// 	fmt.Println("1", err)
	// 	return
	// }

	// fmt.Println(data)

	// sql := txClickhouse.NewQuery(ctx).Table("account_login").MultiInsertSQL(data...)
	// fmt.Println(sql)
	// _, err = txClickhouse.NewQuery(ctx).Table("account_login").MultiInsert(data...)
	// fmt.Println("2", err)
	// if err == nil {
	// 	err = txClickhouse.Commit()
	// 	fmt.Println("3", err)
	// }
	// err = txClickhouse.Rollback()
	// fmt.Println("4", err)

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

// func NewEngine() *Engine {
// 	eng := &Engine{}
// 	if err := eng.Startup(
// 		eng.serveHTTP,
// 		// eng.serveHTTPTwo,
// 	); err != nil {
// 		logger.Logger.Panic("startup", err)
// 	}
// 	// eng.Tracer("ratel", "127.0.0.1:6831")
// 	return eng
// }
// func (eng *Engine) serveHTTP() error {
// 	config := &http.Config{
// 		Host:    "127.0.0.1",
// 		Port:    17777,
// 		Mode:    gin.DebugMode,
// 		Name:    "admin",
// 		Network: "tcp",
// 	}
// 	server := config.Build()
// 	// server.Use(trace.HTTPMetricServerInterceptor())
// 	// server.Use(trace.HTTPTraceServerInterceptor())
// 	server.GET("/ping", "ping", func(ctx *gin.Context) {
// 		// e := Redis.NameSpace("common").Set(ctx.Request.Context(), "aaaaa", "daadasd", time.Minute*5).Err()
// 		ctx.JSON(200, gin.H{
// 			"status": "7777",
// 		})
// 	})

// 	if gin.IsDebugging() {
// 		gin.App.Table.Render()
// 	}

// 	// data := [][]string{
// 	// 	[]string{"A", "The Good", "500"},
// 	// 	[]string{"B", "The Very very Bad Man", "288"},
// 	// 	[]string{"C", "The Ugly", "120"},
// 	// 	[]string{"D", "The Gopher", "800"},
// 	// }

// 	// table := tablewriter.NewWriter(os.Stdout)
// 	// table.SetHeader([]string{"Name", "Sign", "Rating"})

// 	// for _, v := range data {
// 	// 	table.Append(v)
// 	// }
// 	// table.Render()
// 	// }

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

// 	server.Use(trace.HTTPTraceServerInterceptor())

// 	server.InitFuncMap()
// 	pprof.Register(server.Engine)

// 	return eng.Serve(server)
// }
