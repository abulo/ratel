package clickhouse

import (
	"time"

	_ "github.com/abulo/clickhouse-go"
	"github.com/abulo/ratel/v2/logger"
	"github.com/abulo/ratel/v2/store/query"
	"github.com/abulo/ratel/v2/util"
)

//Config 数据库配置
type Config struct {
	Username     string //账号 root
	Password     string //密码
	Host         string //host localhost
	Port         string //端口 3306
	Database     string //默认连接数据库
	ReadTimeout  int
	WriteTimeout int
	LoadBalance  string        //负载均衡
	MaxOpenConns int           //连接池最多同时打开的连接数
	MaxIdleConns int           //连接池里最大空闲连接数。必须要比maxOpenConns小
	MaxLifetime  time.Duration //连接池里面的连接最大存活时长
	MaxIdleTime  time.Duration //连接池里面的连接最大空闲时长
	DriverName   string
	Trace        bool
}

//URI 构造数据库连接
func (config *Config) URI() string {

	//tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000

	link := "tcp://" + config.Host + ":" + config.Port
	param := make([]string, 0)
	if !util.Empty(config.Username) {
		param = append(param, "username="+config.Username)
	}
	if !util.Empty(config.Password) {
		param = append(param, "password="+config.Password)
	}
	if !util.Empty(config.Database) {
		param = append(param, "database="+config.Database)
	}
	if config.ReadTimeout > 0 {
		param = append(param, "read_timeout="+util.ToString(config.ReadTimeout))
	}
	if config.WriteTimeout > 0 {
		param = append(param, "write_timeout="+util.ToString(config.WriteTimeout))
	}
	if !util.Empty(config.LoadBalance) {
		param = append(param, "alt_hosts="+config.LoadBalance)
	}
	param = append(param, "&parseTime=true")
	return link + "?" + util.Implode("&", param)
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
// 	return db
// }

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
