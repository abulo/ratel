package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/abulo/ratel/v3/app"
	"github.com/abulo/ratel/v3/client/grpc"
	"github.com/abulo/ratel/v3/client/grpc/balancer"
	"github.com/abulo/ratel/v3/client/grpc/resolver"
	"github.com/abulo/ratel/v3/example/grpc/love"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/gin"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/logger/mongo"
	"github.com/abulo/ratel/v3/registry/etcdv3"
	"github.com/abulo/ratel/v3/server/xgin"
	"github.com/abulo/ratel/v3/util"
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
		eng.GrpcClient,
	); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("startup")
	}
	return eng
}

func (eng *Engine) GrpcClient() error {
	etcd := etcdv3.New()
	etcd.Endpoints = []string{"172.18.1.13:2379"}
	etcd.Secure = false
	etcd.Prefix = "golang-test"
	etcd.ConnectTimeout = 2 * time.Second
	resolver.Register("etcd", etcd.MustBuild())
	return nil
}

func (eng *Engine) HttpServer() error {
	client := xgin.New()
	client.Host = "127.0.0.1"
	client.Port = 17777
	client.Deployment = "golang"
	client.DisableMetric = false
	client.DisableTrace = false
	client.DisableSlowQuery = true
	client.ServiceAddress = "127.0.0.1:17777"
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

	r.GET("/index", "index", index)
}

func index(ctx *gin.Context) {

	config := grpc.New()
	config.Address = "etcd:///grpc"
	config.BalancerName = balancer.NameSmoothWeightRoundRobin
	// config.BalancerName = p2c.Name
	config.Block = false
	config.DialTimeout = 1 * time.Second
	config.OnDialError = "info"

	client := love.NewLoveClient(config.Build())
	newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second)
	defer cancel()
	name := util.ToString(util.Timestamp())
	r, err := client.Confession(newCtx, &love.Request{Name: name})
	// if err != nil {
	// log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.GetResult())

	ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "无效数据", "data": gin.H{
		"data": name,
		"req":  r,
		"err":  err,
	}})
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
