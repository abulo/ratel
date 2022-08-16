package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/abulo/ratel/v3/app"
	"github.com/abulo/ratel/v3/client/grpc/resolver"
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
	fmt.Println(global.Path)
	//获取配置文件目录
	configPath := global.GetEnvironment(global.Path+"/config/env", "configDir")
	if util.Empty(configPath) {
		panic("configPath is empty")
	}

	//配置加载 toml 文件
	dirs := make([]string, 0)
	dirs = append(dirs, global.Path+"/config/"+configPath)

	global.InitConfig(dirs...)
	global.InitMongoDB()
	global.InitRedis()
	global.InitMysql()
	global.InitElasticSearch()
	global.InitSession("session")
	global.InitTrace()
}

// Engine ...
type Engine struct {
	app.Application
}

// NewEngine ...
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

// GrpcClient ...
func (eng *Engine) GrpcClient() error {
	etcd := etcdv3.New()
	etcd.Endpoints = []string{"172.18.1.13:2379"}
	etcd.Secure = false
	etcd.Prefix = "abulo"
	etcd.ConnectTimeout = 2 * time.Second
	resolver.Register("etcd", etcd.MustBuild())
	return nil
}

// HttpServer ...
func (eng *Engine) HttpServer() error {
	client := xgin.New()
	client.Host = "172.18.1.5"
	client.Port = 17777
	client.Deployment = "golang"
	client.DisableMetric = false
	client.DisableTrace = false
	client.DisableSlowQuery = true
	client.ServiceAddress = "172.18.1.5:17777"
	client.SlowQueryThresholdInMilli = 10000
	client.Mode = gin.DebugMode
	server := client.Build()
	server.SetTrustedProxies([]string{"0.0.0.0/0"})
	t, err := loadGlobTemplate(initial.Core.Path + "/view")
	if err != nil {
		panic(err)
	}
	server.LoadHTMLFiles(t...)
	server.Use(gin.Logger())

	Route(server.Engine)
	if gin.IsDebugging() {
		gin.App.Table.Render()
	}
	return eng.Serve(server)
}

// Route ...
func Route(ctx *gin.Engine) {

	ctx.GET("/index", "index", index)
}

func index(ctx *gin.Context) {

	db := initial.Core.Store.LoadSQL("mysql").Read()

	sql := "SELECT video.id as video_id,video.* ,provider.title as provider_title,GROUP_CONCAT(sp.title) as sp_title FROM video LEFT JOIN provider_video ON video.id = provider_video.video_id LEFT JOIN provider ON provider_video.provider_id = provider.id LEFT JOIN sp_provider ON provider.id = sp_provider.provider_id LEFT JOIN sp ON sp_provider.sp_id = sp.id "

	db.NewQuery(ctx.Request.Context()).QueryRows(sql).ToMap()

	items := make([]map[string]interface{}, 0)
	err := initial.Core.Store.LoadMongoDB("mongodb").Collection("op_logger").FindMany(ctx.Request.Context(), &items)

	fmt.Println(err)

	ctx.HTML(http.StatusOK, "manager/login.html", gin.H{
		"redirect": "ddd",
	})

	// config := grpc.New()
	// config.Address = "etcd:///grpc"
	// config.BalancerName = balancer.NameSmoothWeightRoundRobin
	// // config.BalancerName = p2c.Name
	// config.Block = false
	// config.DialTimeout = 1 * time.Second
	// config.OnDialError = "info"

	// client := love.NewLoveClient(config.Build())
	// newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second)
	// defer cancel()
	// name := util.ToString(util.Timestamp())
	// r, err := client.Confession(newCtx, &love.Request{Name: name})
	// // if err != nil {
	// // log.Fatalf("could not greet: %v", err)
	// // }
	// // log.Printf("Greeting: %s", r.GetResult())

	// ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "无效数据", "data": gin.H{
	// 	"data": name,
	// 	"req":  r,
	// 	"err":  err,
	// }})
}
func main() {

	mongodbClient := initial.Core.Store.LoadMongoDB("mongodb")
	loggerHook := mongo.DefaultWithURL(mongodbClient)
	defer loggerHook.Flush()
	logger.Logger.AddHook(loggerHook)
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Logger.SetReportCaller(true)
	if initial.Core.Config.Bool("DisableDebug", true) {
		logger.Logger.SetOutput(io.Discard)
	} else {
		logger.Logger.SetOutput(os.Stdout)
	}

	eng := NewEngine()
	if err := eng.Run(); err != nil {
		logger.Logger.Panic(err.Error())
	}
}

//加载模板文件
// func loadTemplate(funcMap template.FuncMap, r render.Delims) (*template.Template, error) {
// 	t := template.New("").Delims(r.Left, r.Right).Funcs(funcMap)

// 	for _, name := range view.AssetNames() {
// 		if !strings.HasSuffix(name, ".html") {
// 			continue
// 		}
// 		asset, err := view.Asset(name)
// 		if err != nil {
// 			continue
// 		}
// 		name := strings.Replace(name, "view/", "", 1)
// 		t, err = t.New(name).Parse(string(asset))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return t, nil
// }

func loadGlobTemplate(dir string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, filepath.FromSlash(path))
		}
		return nil
	})
	return fileList, err
}
