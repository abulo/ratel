package mysql

import (
	"time"

	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/store/query"
	_ "github.com/go-sql-driver/mysql"
)

//Config 数据库配置
type Config struct {
	Username        string        //账号 root
	Password        string        //密码
	Host            string        //host localhost
	Port            string        //端口 3306
	Charset         string        //字符编码 utf8mb4
	Database        string        //默认连接数据库
	ConnMaxLifetime time.Duration //设置一个连接的最长生命周期，因为数据库本身对连接有一个超时时间的设置，如果超时时间到了数据库会单方面断掉连接，此时再用连接池内的连接进行访问就会出错, 因此这个值往往要小于数据库本身的连接超时时间
	ConnMaxIdleTime time.Duration //设置连接的生命周期的最大
	MaxIdleConns    int           //设置闲置的连接数,连接池里面允许Idel的最大连接数, 这些Idel的连接 就是并发时可以同时获取的连接,也是用完后放回池里面的互用的连接, 从而提升性能
	MaxOpenConns    int           //设置最大打开的连接数，默认值为0表示不限制。控制应用于数据库建立连接的数量，避免过多连接压垮数据库。
	DriverName      string
}

//New 新连接
func New(config *Config) *query.QueryDb {
	db, err := query.NewSqlConn(config.DriverName, config.URI())
	if err != nil {
		logger.Logger.Panic(err)
	}
	return &query.QueryDb{DB: db, DriverName: config.DriverName}
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

// //connect 数据库连接
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
