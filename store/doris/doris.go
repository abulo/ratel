package mysql

import (
	"time"

	"github.com/abulo/ratel/v2/logger"
	"github.com/abulo/ratel/v2/store/query"
	_ "github.com/go-sql-driver/mysql"
)

//Config 数据库配置
type Config struct {
	Username     string        //账号 root
	Password     string        //密码
	Host         string        //host localhost
	Port         string        //端口 3306
	Charset      string        //字符编码 utf8mb4
	Database     string        //默认连接数据库
	MaxOpenConns int           //连接池最多同时打开的连接数
	MaxIdleConns int           //连接池里最大空闲连接数。必须要比maxOpenConns小
	MaxLifetime  time.Duration //连接池里面的连接最大存活时长
	MaxIdleTime  time.Duration //连接池里面的连接最大空闲时长
	DriverName   string
	Trace        bool
}

//New 新连接
func New(config *Config) *query.QueryDb {
	opt := &query.Opt{
		MaxOpenConns: config.MaxOpenConns,
		MaxIdleConns: config.MaxIdleConns,
		MaxLifetime:  config.MaxLifetime,
		MaxIdleTime:  config.MaxIdleTime,
	}

	db, err := query.NewSqlConn(config.DriverName, config.URI(), opt)
	if err != nil {
		logger.Logger.Panic(err)
	}
	return &query.QueryDb{DB: db, DriverName: config.DriverName, Trace: config.Trace}
}

//URI 构造数据库连接
func (config *Config) URI() string {
	return config.Username + ":" +
		config.Password + "@tcp(" +
		config.Host + ":" +
		config.Port + ")/" +
		config.Database + "?charset=" +
		config.Charset + "&loc=" + time.Local.String() + "&parseTime=true"
}

//connect 数据库连接
// func connect(config *Config) *sql.DB {
// 	//数据库连接
// 	db, err := sql.Open(config.DriverName, config.URI())
// 	if err != nil {
// 		logger.Logger.Fatal(err.Error())
// 	}
// 	if err = db.Ping(); err != nil {
// 		logger.Logger.Fatal(err.Error())
// 	}
// 	db.SetMaxIdleConns(config.MaxIdleConns)
// 	db.SetMaxOpenConns(config.MaxOpenConns)
// 	db.SetConnMaxLifetime(config.ConnMaxLifetime)
// 	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
// 	return db
// }
