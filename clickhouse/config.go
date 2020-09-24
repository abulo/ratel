package clickhouse

import (
	"database/sql"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/abulo/ratel/logger"
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
}

//URI 构造数据库连接
func (config *Config) URI() string {

	//tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000
	return "tcp://" +
		config.Host + ":" +
		config.Port + "?" +
		"username=" + config.Username +
		"&password=" + config.Password +
		"&database=" + config.Database +
		"&read_timeout=" + util.ToString(config.ReadTimeout) +
		"&write_timeout=" + util.ToString(config.WriteTimeout) +
		"&alt_hosts=" + config.LoadBalance
}

//connect 数据库连接
func connect(config *Config) *sql.DB {
	//数据库连接
	db, err := sql.Open("clickhouse", config.URI())
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	if err = db.Ping(); err != nil {
		logger.Logger.Fatal(err.Error())
	}
	return db
}

//New 新连接
func New(config *Config) *QueryDb {
	db := connect(config)
	return &QueryDb{db: db}
}
