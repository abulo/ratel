package main

import (
	"io/ioutil"
	"os"

	"github.com/abulo/ratel/v3/app"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/gin"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/logger/mongo"
	"github.com/abulo/ratel/v3/pprof"
	"github.com/abulo/ratel/v3/server/xgin"
	"github.com/sirupsen/logrus"
)

func init() {
	//全局设置
	global := initial.Default()
	//配置加载 toml 文件
	dirs := make([]string, 0)
	dirs = append(dirs, global.Path+"/config")
	global.InitConfig(dirs...)
	global.InitMongoDB()
	global.InitRedis()
	global.InitMysql()
	global.InitElasticSearch()
	global.InitSession("session")
	global.InitTrace()
}

type Engine struct {
	app.Application
}

func NewEngine() *Engine {
	eng := &Engine{}

	if err := eng.Startup(
		eng.HttpServer,
	); err != nil {
		logger.Logger.Panic("startup", "err", err)
	}
	return eng
}

func (eng *Engine) HttpServer() error {
	client := xgin.New()
	client.Host = "0.0.0.0"
	client.Port = 17777
	client.Deployment = "golang"
	client.DisableMetric = false
	client.DisableTrace = false
	client.DisableSlowQuery = true
	client.ServiceAddress = "0.0.0.0:17777"
	client.SlowQueryThresholdInMilli = 10000
	client.Mode = gin.DebugMode
	server := client.Build()
	server.SetTrustedProxies([]string{"0.0.0.0/0"})
	Route(server.Engine)
	if gin.IsDebugging() {
		gin.App.Table.Render()
	}
	return eng.Serve(server)
}
func Route(r *gin.Engine) {
	pprof.Register(r)
}
func main() {
	mongodbClient := initial.Core.Store.LoadMongoDB("mongodb")
	loggerHook := mongo.DefaultWithURL(mongodbClient)
	defer loggerHook.Flush()
	logger.Logger.AddHook(loggerHook)
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Logger.SetReportCaller(true)
	if initial.Core.Config.Bool("DisableDebug", true) {
		logger.Logger.SetOutput(ioutil.Discard)
	} else {
		logger.Logger.SetOutput(os.Stdout)
	}

	eng := NewEngine()
	if err := eng.Run(); err != nil {
		logger.Logger.Panic(err.Error())
	}
}
