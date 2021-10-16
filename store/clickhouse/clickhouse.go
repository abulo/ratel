package clickhouse

import (
	"database/sql"

	_ "github.com/abulo/clickhouse-go"
	"github.com/abulo/ratel/logger"
	"github.com/abulo/ratel/store/base"
	"github.com/abulo/ratel/util"
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
	LoadBalance  string //负载均衡
	DriverName   string
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
func connect(config *Config) *sql.DB {
	//数据库连接
	db, err := sql.Open(config.DriverName, config.URI())
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	if err = db.Ping(); err != nil {
		logger.Logger.Fatal(err.Error())
	}
	return db
}

//New 新连接
func New(config *Config) *base.QueryDb {
	db := connect(config)
	return &base.QueryDb{DB: db, DriverName: config.DriverName}
}
