package base

import (
	"errors"
	"os"
	"path"
	"time"

	"github.com/abulo/ratel/v2/config"
	"github.com/abulo/ratel/v2/config/toml"
	"github.com/abulo/ratel/v2/stores/mysql"
	"github.com/abulo/ratel/v2/stores/query"
	"github.com/abulo/ratel/v2/util"
	"github.com/spf13/cast"
)

var Config *config.Config
var Query *query.Query
var Path string

// InitPath 初始化路径
func InitPath() error {
	wd, err := os.Getwd()
	if err != nil {
		// fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	Path = wd
	return nil
}

// InitConfig 初始化
func InitConfig() error {
	configPath := path.Join(Path, "mysql.toml")
	if !util.FileExists(configPath) {
		err := errors.New("数据库配置文件不存在")
		// fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	//加载配置文件
	Config = config.New("dao")
	Config.AddDriver(toml.Driver)
	Config.LoadFiles(configPath)
	return nil
}

// InitQuery 初始化数据查询
func InitQuery() error {
	//创建数据链接
	opt := &mysql.Config{}
	if Username := cast.ToString(Config.String("mysql.Username")); Username != "" {
		opt.Username = Username
	}
	if Password := cast.ToString(Config.String("mysql.Password")); Password != "" {
		opt.Password = Password
	}
	if Host := cast.ToString(Config.String("mysql.Host")); Host != "" {
		opt.Host = Host
	}
	if Port := cast.ToString(Config.String("mysql.Port")); Port != "" {
		opt.Port = Port
	}
	if Charset := cast.ToString(Config.String("mysql.Charset")); Charset != "" {
		opt.Charset = Charset
	}
	if Database := cast.ToString(Config.String("mysql.Database")); Database != "" {
		opt.Database = Database
	}
	if Local := cast.ToString(Config.String("mysql.Local")); Local != "" {
		opt.Local = Local
	}
	// # MaxOpenConns 连接池最多同时打开的连接数
	// MaxOpenConns = 128
	// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxIdleConns = 32
	// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
	// MaxLifetime = 10
	// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
	// MaxIdleTime = 5
	if MaxLifetime := cast.ToInt(Config.Int("mysql.MaxLifetime")); MaxLifetime > 0 {
		opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
	}
	if MaxIdleTime := cast.ToInt(Config.Int("mysql.MaxIdleTime")); MaxIdleTime > 0 {
		opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
	}
	if MaxIdleConns := cast.ToInt(Config.Int("mysql.MaxIdleConns")); MaxIdleConns > 0 {
		opt.MaxIdleConns = cast.ToInt(MaxIdleConns)
	}
	if MaxOpenConns := cast.ToInt(Config.Int("mysql.MaxOpenConns")); MaxOpenConns > 0 {
		opt.MaxOpenConns = cast.ToInt(MaxOpenConns)
	}
	opt.DriverName = "mysql"
	opt.DisableMetric = cast.ToBool(Config.Bool("mysql.DisableMetric"))
	opt.DisableTrace = cast.ToBool(Config.Bool("mysql.DisableTrace"))
	Query = mysql.NewClient(opt)
	return nil
}

// InitBase 初始化数据
func InitBase() error {
	if err := InitPath(); err != nil {
		return err
	}
	if err := InitConfig(); err != nil {
		return err
	}
	if err := InitQuery(); err != nil {
		return err
	}
	return nil
}
