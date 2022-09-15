package main

import (
	"github.com/abulo/ratel/v3/env"
	"github.com/abulo/ratel/v3/example/initial"
	"github.com/abulo/ratel/v3/gen/module"
	"github.com/abulo/ratel/v3/util"
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

func main() {
	filePath := "/home/www/golang/src/ratel/gen/module/template.tpl"
	tableName := "area"
	dao := module.CamelStr(tableName)
	outputDir := "/home/www/golang/src/ratel/gen/lib/module/" + util.StrToLower(module.CamelStr(tableName))
	outputPackage := util.StrToLower(module.CamelStr(tableName))

	db := initial.Core.Store.LoadSQL("mysql")
	module.Run(db.Write(), tableName, outputDir, outputPackage, dao, filePath)
}
