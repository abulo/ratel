package main

import (
	ohttp "net/http"
	"os"

	"github.com/abulo/ratel"
	"github.com/abulo/ratel/gin"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/logger/hook"
	"github.com/abulo/ratel/mongodb"
	"github.com/abulo/ratel/redis"
	"github.com/abulo/ratel/server/http"
	"github.com/abulo/ratel/server/monitor"
	"github.com/abulo/ratel/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	ratel.Ratel
}

//MongoDB 代理
var MongoDB *mongodb.Proxy = mongodb.NewProxy()
var Redis *redis.Proxy = redis.NewProxy()

// func init() {

// }

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
	eng := NewEngine()

	if err := eng.Run(); err != nil {
		logger.Logger.Panic(err)
	}
}
func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
		eng.serveHTTPTwo,
	); err != nil {
		logger.Logger.Panic("startup", err)
	}
	eng.Tracer("ratel", "127.0.0.1:6831")
	return eng
}
func (eng *Engine) serveHTTP() error {
	config := &http.Config{
		Host: "127.0.0.1",
		Port: 7777,
		Mode: gin.DebugMode,
		Name: "admin",
	}
	server := config.Build()
	server.Use(trace.HTTPMetricServerInterceptor())
	server.Use(trace.HTTPTraceServerInterceptor())
	server.GET("/ping", "ping", func(ctx *gin.Context) {
		// e := Redis.NameSpace("common").Set(ctx.Request.Context(), "aaaaa", "daadasd", time.Minute*5).Err()
		ctx.JSON(200, gin.H{
			"status": "7777",
		})
	})

	return eng.Serve(server)
}

func (eng *Engine) serveHTTPTwo() error {
	config := &monitor.Config{
		Host:    "127.0.0.1",
		Port:    17777,
		Network: "tcp4",
		Name:    "monitor",
	}
	// monitor.HandleFunc("/metrics", func(w ohttp.ResponseWriter, r *ohttp.Request) {
	// 	promhttp.Handler().ServeHTTP(w, r)
	// })
	server := config.Build()

	server.HandleFunc("/metrics", func(w ohttp.ResponseWriter, r *ohttp.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// server.Use(trace.HTTPTraceServerInterceptor())

	// 	server.InitFuncMap()
	// 	pprof.Register(server.Engine)

	return eng.Serve(server)
}
