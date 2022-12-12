package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/abulo/ratel/v3/core/env"
	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/logger/mongo"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/stores/null"
	"github.com/abulo/ratel/v3/util"
	"github.com/sirupsen/logrus"
)

func init() {
	// var cstZone = time.FixedZone("GMT", 8*3600)
	// var cstZone = time.FixedZone("CST", 8*3600)
	var cstZone, _ = time.LoadLocation("PRC")
	time.Local = cstZone
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

type Test struct {
	Id            int64          `db:"id" json:"id" form:"id" uri:"id" xml:"id" proto:"id"`                                                                    //DataType:bigint
	NullString    null.String    `db:"null_string" json:"nullString" form:"nullString" uri:"nullString" xml:"nullString" proto:"nullString"`                   //DataType:varchar
	NullDate      null.Date      `db:"null_date" json:"nullDate" form:"nullDate" uri:"nullDate" xml:"nullDate" proto:"nullDate"`                               //DataType:date
	NullDatetime  null.DateTime  `db:"null_datetime" json:"nullDatetime" form:"nullDatetime" uri:"nullDatetime" xml:"nullDatetime" proto:"nullDatetime"`       //DataType:datetime
	NullYear      null.Int32     `db:"null_year" json:"nullYear" form:"nullYear" uri:"nullYear" xml:"nullYear" proto:"nullYear"`                               //DataType:year
	NullTime      null.CTime     `db:"null_time" json:"nullTime" form:"nullTime" uri:"nullTime" xml:"nullTime" proto:"nullTime"`                               //DataType:time
	NullTimestamp null.TimeStamp `db:"null_timestamp" json:"nullTimestamp" form:"nullTimestamp" uri:"nullTimestamp" xml:"nullTimestamp" proto:"nullTimestamp"` //DataType:timestamp
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
	ctx := context.TODO()

	db := initial.Core.Store.LoadSQL("mysql").Write()
	var res Test
	err := db.NewBuilder(ctx).Table("test").Where("id", 1).Row().ToStruct(&res)
	// fmt.Println(err)
	// fmt.Println("=====")
	fmt.Println(res)
	// fmt.Println("=====")
	str, _ := json.Marshal(res)
	fmt.Println(string(str))

	var update Test
	// update.NullString
	update.NullString = null.NewString("", false)

	id, err := db.NewBuilder(ctx).Table("test").Where("id", 1).Update(update)

	fmt.Println(id, err)
}
