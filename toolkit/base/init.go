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

var (
	Config         *config.Config // 配置文件
	Query          *gorm.DB       // 数据库连接
	Path           *string        // 路径
	Url            *string        // url
	Driver         *string        // 驱动
	DefaultDefault = "mysql"      // 默认驱动
)

func SetUrl(url string) {
	Url = &url
}

// InitPath 初始化路径
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

// InitConfig 初始化
func InitConfig() error {
	configPath := path.Join(*Path, "toolkit.toml")
	if !util.FileExists(configPath) {
		err := errors.New("配置文件不存在")
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	//加载配置文件
	Config = config.New()
	if err := Config.LoadFile(configPath); err != nil {
		fmt.Println("初始化目录错误:", color.RedString(err.Error()))
		return err
	}
	return nil
}

// InitQuery 初始化数据查询
func InitQuery() error {
	opts := make([]sql.Option, 0)
	// DisableMetric = false # 关闭指标采集
	// DisableTrace = false # 关闭链路追踪
	// DriverName = "mysql" # 数据库驱动名称
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
	if DisableMetric := cast.ToBool(Config.Bool("db.DisableMetric")); DisableMetric {
		opts = append(opts, sql.WithDisableMetric(DisableMetric))
	}
	if DisableTrace := cast.ToBool(Config.Bool("db.DisableTrace")); DisableTrace {
		opts = append(opts, sql.WithDisableTrace(DisableTrace))
	}
	if DriverName := cast.ToString(Config.String("db.DriverName")); DriverName != "" {
		opts = append(opts, sql.WithDriverName(DriverName))
		Driver = &DriverName
	}
	client, err := sql.NewClient(opts...)
	if err != nil {
		return err
	}
	Query, _ = client.SqlConn()
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
