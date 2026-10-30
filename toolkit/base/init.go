package base

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/stores/sql"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// 以下包级变量为代码生成流程共享的运行时状态，由 InitBase 系列函数初始化，非并发安全。
var (
	// Config 为加载完成的 toolkit.toml 配置实例。
	Config *config.Config
	// Query 为已建立连接的数据库查询句柄。
	Query *gorm.DB
	// Path 为当前项目的工作目录。
	Path *string
	// Url 为通过 SetUrl 设置的模板或仓库地址。
	Url *string
	// Driver 为当前数据库驱动名，默认 mysql。
	Driver *string
	// DefaultDefault 为未在配置中指定驱动时使用的默认驱动名 mysql。
	DefaultDefault = "mysql"
)

// SetUrl 将模板或仓库地址保存到包级变量 Url。
func SetUrl(url string) {
	Url = &url
}

// InitPath 初始化当前工作目录 Path，并将数据库驱动重置为默认驱动；获取当前目录失败时返回错误。
func InitPath() error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	Path = &wd
	Driver = &DefaultDefault
	return nil
}

// InitConfig 读取工作目录下的 toolkit.toml 并加载到包级变量 Config；配置文件缺失或解析失败时返回错误。
func InitConfig() error {
	configPath := path.Join(*Path, "toolkit.toml")
	if !util.FileExists(configPath) {
		err := errors.New("配置文件不存在")
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	Config = config.New()
	if err := Config.LoadFile(configPath); err != nil {
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	return nil
}

// InitQuery 依据 Config 中的 db.* 配置项创建数据库客户端并保存到 Query；创建失败时返回错误。
//
// 指标采集、链路追踪与调试日志在此流程中固定关闭。
func InitQuery() error {
	opts := make([]sql.Option, 0)
	if Host := cast.ToString(Config.String("db.Host")); Host != "" {
		opts = append(opts, sql.WithHost(Host))
	}
	if Port := cast.ToString(Config.String("db.Port")); Port != "" {
		opts = append(opts, sql.WithPort(Port))
	}
	if Username := cast.ToString(Config.String("db.Username")); Username != "" {
		opts = append(opts, sql.WithUsername(Username))
	}
	if Password := cast.ToString(Config.String("db.Password")); Password != "" {
		opts = append(opts, sql.WithPassword(Password))
	}
	if Charset := cast.ToString(Config.String("db.Charset")); Charset != "" {
		opts = append(opts, sql.WithCharset(Charset))
	}
	if Database := cast.ToString(Config.String("db.Database")); Database != "" {
		opts = append(opts, sql.WithDatabase(Database))
	}
	if ParseTime := cast.ToBool(Config.Bool("db.ParseTime")); ParseTime {
		opts = append(opts, sql.WithParseTime(ParseTime))
	}
	if TimeZone := cast.ToString(Config.String("db.TimeZone")); TimeZone != "" {
		opts = append(opts, sql.WithTimeZone(TimeZone))
	}
	if SslMode := cast.ToBool(Config.Bool("db.SslMode")); SslMode {
		opts = append(opts, sql.WithSslMode(SslMode))
	}
	if DialTimeOut := cast.ToString(Config.String("db.DialTimeOut")); DialTimeOut != "" {
		opts = append(opts, sql.WithDialTimeOut(DialTimeOut))
	}
	if ReadTimeOut := cast.ToString(Config.String("db.ReadTimeOut")); ReadTimeOut != "" {
		opts = append(opts, sql.WithReadTimeOut(ReadTimeOut))
	}
	if MaxIdleConns := cast.ToInt(Config.Int("db.MaxIdleConns")); MaxIdleConns > 0 {
		opts = append(opts, sql.WithMaxIdleConns(MaxIdleConns))
	}
	if MaxOpenConns := cast.ToInt(Config.Int("db.MaxOpenConns")); MaxOpenConns > 0 {
		opts = append(opts, sql.WithMaxOpenConns(MaxOpenConns))
	}
	if MaxLifetime := cast.ToInt(Config.Int("db.MaxLifetime")); MaxLifetime > 0 {
		opts = append(opts, sql.WithMaxLifetime(time.Duration(MaxLifetime)*time.Second))
	}
	if MaxIdleTime := cast.ToInt(Config.Int("db.MaxIdleTime")); MaxIdleTime > 0 {
		opts = append(opts, sql.WithMaxIdleTime(time.Duration(MaxIdleTime)*time.Second))
	}
	opts = append(opts, sql.WithEnableMetric(false))
	opts = append(opts, sql.WithEnableTrace(false))
	if DriverName := cast.ToString(Config.String("db.DriverName")); DriverName != "" {
		opts = append(opts, sql.WithDriverName(DriverName))
		Driver = &DriverName
	}
	opts = append(opts, sql.WithEnableDebug(false))
	client, err := sql.NewClient(opts...)
	if err != nil {
		return err
	}
	Query, _ = client.SqlConn()
	return nil
}

// InitBase 按顺序初始化工作目录、配置文件与数据库连接，任一步骤失败时直接返回该步骤的错误。
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
