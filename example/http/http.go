package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/logger/mongo"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/gin"
	"github.com/abulo/ratel/v3/pprof"
	"github.com/abulo/ratel/v3/server/xgin"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func init() {
	// 全局设置
	global := initial.New()
	// 获取配置文件目录
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
	global.InitClickHouse()
	global.InitTrace()

	env.SetAppID("1")
	env.SetAppRegion("sc")
	env.SetBuildVersion("00000000")
	env.SetBuildTime("2021-01-01 23;23:23")
	env.SetAppZone("cs")
}

// 权限表
type Permission struct {
	Id       int64          `db:"id" json:"id"`              //主键,自增(PRI)
	ParentId int64          `db:"parent_id" json:"parentId"` //父ID
	Title    string         `db:"title" json:"title"`        //权限名称
	Handle   string         `db:"handle" json:"handle"`      //路由别名(UNI)
	Type     string         `db:"type" json:"type"`          //分类(作用前后端)(MUL)
	Weight   int64          `db:"weight" json:"weight"`      //权重
	CreateAt query.NullTime `db:"create_at" json:"createAt"` //
	UpdateAt query.NullTime `db:"update_at" json:"updateAt"` //
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

type Engine struct {
	app.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	//加载计划任务
	// eng.Schedule(eng.CrontabWork())
	// 注册函数
	// eng.RegisterHooks(hooks.Stage_AfterLoadConfig, eng.BeforeInit)
	if err := eng.Startup(
		eng.NewAdminHttpServer,
		// eng.ApiServer,
	); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panic("startup")
	}
	return eng
}

func (eng *Engine) NewAdminHttpServer() error {
	configAdmin := initial.Core.Config.Get("server.admin")
	cfg := configAdmin.(map[string]interface{})
	//先获取这个服务是否是需要开启
	if disable := cast.ToBool(cfg["Disable"]); disable {
		logger.Logger.Error("server.admin is disabled")
		return nil
	}
	client := xgin.New()
	client.Host = cast.ToString(cfg["Host"])
	client.Port = cast.ToInt(cfg["Port"])
	client.Deployment = cast.ToString(cfg["Deployment"])
	client.DisableMetric = cast.ToBool(cfg["DisableMetric"])
	client.DisableTrace = cast.ToBool(cfg["DisableTrace"])
	client.DisableSlowQuery = cast.ToBool(cfg["DisableSlowQuery"])
	client.ServiceAddress = cast.ToString(cfg["ServiceAddress"])
	client.SlowQueryThresholdInMilli = cast.ToInt64(cfg["SlowQueryThresholdInMilli"])
	server := client.Build()
	if !initial.Core.Config.Bool("DisableDebug", true) {
		client.Mode = gin.DebugMode
		server.Use(gin.Logger())
		server.Use(gin.Recovery())
	} else {
		client.Mode = gin.ReleaseMode
	}
	server.SetTrustedProxies([]string{"0.0.0.0/0"})
	//添加路由
	pprof.Register(server.Engine)
	Route(server.Engine)
	if gin.IsDebugging() {
		gin.App.Table.Render()
	}
	return eng.Serve(server)
}

func Route(r *gin.Engine) {
	groupRoutes := r.Group("/admin/v1")
	{
		groupRoutes.GET("/permission/:id", "test", Show)
	}
}

func Show(ctx *gin.Context) {
	id := cast.ToInt64(ctx.Param("id"))
	//实例化
	data, err := ShowData(ctx.Request.Context(), id)
	err = nil
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	cache := initial.Core.Store.LoadRedis("redis")
	cache.Set(ctx.Request.Context(), "test", "test", -1).Result()
	//成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": data,
	})
}

func ShowData(ctx context.Context, id int64) (Permission, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res Permission
	err := db.NewBuilder(ctx).Table("permission").Where("id", id).Row().ToStruct(&res)
	return res, err
}

//struct{}{}
//[]interface{}
