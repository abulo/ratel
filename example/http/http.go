package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/logger/mongo"
	"github.com/abulo/ratel/v3/stores/query"
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

// 权限表
type Permission struct {
	Id       int64              `db:"id" json:"id"`              //主键,自增(PRI)
	ParentId int64              `db:"parent_id" json:"parentId"` //父ID
	Title    string             `db:"title" json:"title"`        //权限名称
	Handle   string             `db:"handle" json:"handle"`      //路由别名(UNI)
	Type     string             `db:"type" json:"type"`          //分类(作用前后端)(MUL)
	Weight   int64              `db:"weight" json:"weight"`      //权重
	CreateAt query.NullDateTime `db:"create_at" json:"createAt"` //
	UpdateAt query.NullDateTime `db:"update_at" json:"updateAt"` //
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
	ctx := context.Background()
	data, err := ValidationPermissionByHandleAndType(ctx, "index", "api")
	fmt.Println(data, err)
	s, _ := json.Marshal(data)
	fmt.Println(string(s))
}

func ValidationPermissionByHandleAndType(ctx context.Context, handle, ItemType string) ([]Permission, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	var res []Permission
	err := db.NewBuilder(ctx).Table("permission").Where("type", ItemType).Where("handle", handle).Rows().ToStruct(&res)
	return res, err
}
