package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/abulo/ratel/v3/app"
	"github.com/abulo/ratel/v3/client/grpc/balancer"
	"github.com/abulo/ratel/v3/example/grpc/love"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/logger/mongo"
	"github.com/abulo/ratel/v3/registry/etcdv3"
	"github.com/abulo/ratel/v3/server/xgrpc"
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

	// eng.SetRegistry()

	etcd := etcdv3.New()
	etcd.Endpoints = []string{"172.18.1.13:2379"}
	etcd.Secure = false
	etcd.Prefix = "golang-test"
	etcd.ConnectTimeout = 2 * time.Second

	eng.SetRegistry(
		etcd.MustBuild(),
	)

	if err := eng.Startup(
		eng.GrpcServer,
	); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("startup")
	}
	return eng
}

func (eng *Engine) GrpcServer() error {
	client := xgrpc.New()
	client.Name = balancer.NameSmoothWeightRoundRobin
	client.Host = "127.0.0.1"
	client.Port = 18888
	client.Deployment = "golang"
	server := client.MustBuild()
	love.RegisterLoveServer(server.Server, &Greeter{
		server: server,
	})
	return eng.Serve(server)
}

type Greeter struct {
	server *xgrpc.Server
	love.UnimplementedLoveServer
}

func (g *Greeter) Confession(context context.Context, request *love.Request) (*love.Response, error) {
	return &love.Response{
		Result: request.Name,
	}, nil
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
