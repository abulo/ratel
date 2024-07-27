package base

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"

	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/toml"
	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
)

var Config *config.Config
var Query sql.SqlConn
var Path string
var Url string

func SetUrl(url string) {
	Url = url
}

// InitPath 初始化路径
func InitPath() error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	Path = wd
	return nil
}

// InitConfig 初始化
func InitConfig() error {
	configPath := path.Join(Path, "toolkit.toml")
	if !util.FileExists(configPath) {
		err := errors.New("配置文件不存在")
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
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

	opts := make([]sql.Option, 0)

	if Username := cast.ToString(Config.String("db.Username")); Username != "" {
		opts = append(opts, sql.WithUsername(Username))
	}
	if Password := cast.ToString(Config.String("db.Password")); Password != "" {
		opts = append(opts, sql.WithPassword(Password))
	}
	if Host := cast.ToString(Config.String("db.Host")); Host != "" {
		opts = append(opts, sql.WithHost(Host))
	}
	if Port := cast.ToString(Config.String("db.Port")); Port != "" {
		opts = append(opts, sql.WithPort(Port))
	}
	if Charset := cast.ToString(Config.String("db.Charset")); Charset != "" {
		opts = append(opts, sql.WithCharset(Charset))
	}
	if Database := cast.ToString(Config.String("db.Database")); Database != "" {
		opts = append(opts, sql.WithDatabase(Database))
	}
	if TimeZone := cast.ToString(Config.String("db.TimeZone")); TimeZone != "" {
		opts = append(opts, sql.WithTimeZone(TimeZone))
	}
	if DriverName := cast.ToString(Config.String("db.DriverName")); DriverName != "" {
		opts = append(opts, sql.WithDriverName(DriverName))
	}
	// # MaxOpenConns 连接池最多同时打开的连接数
	// MaxOpenConns = 128
	// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxIdleConns = 32
	// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
	// MaxLifetime = 10
	// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
	// MaxIdleTime = 5
	if MaxLifetime := cast.ToInt(Config.Int("db.MaxLifetime")); MaxLifetime > 0 {
		opts = append(opts, sql.WithMaxLifetime(time.Duration(MaxLifetime)*time.Minute))
	}
	if MaxIdleTime := cast.ToInt(Config.Int("db.MaxIdleTime")); MaxIdleTime > 0 {
		opts = append(opts, sql.WithMaxIdleTime(time.Duration(MaxIdleTime)*time.Minute))
	}
	if MaxIdleConns := cast.ToInt(Config.Int("db.MaxIdleConns")); MaxIdleConns > 0 {
		opts = append(opts, sql.WithMaxIdleConns(MaxIdleConns))
	}
	if MaxOpenConns := cast.ToInt(Config.Int("db.MaxOpenConns")); MaxOpenConns > 0 {
		opts = append(opts, sql.WithMaxOpenConns(MaxOpenConns))
	}
	opts = append(opts, sql.WithDisableMetric(cast.ToBool(Config.Bool("db.DisableMetric"))))
	opts = append(opts, sql.WithDisableTrace(cast.ToBool(Config.Bool("db.DisableTrace"))))
	opts = append(opts, sql.WithDisablePrepare(cast.ToBool(Config.Bool("db.DisablePrepare"))))
	opts = append(opts, sql.WithParseTime(cast.ToBool(Config.Bool("db.ParseTime"))))
	client, err := sql.NewClient(opts...)
	if err != nil {
		return err
	}
	Query = client.NewSqlClient()
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
