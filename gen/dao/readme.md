```golang
import (
	"time"

	"github.com/abulo/ratel/v3/gen/mysql2struct"
	"github.com/abulo/ratel/v3/store/mysql"
	"github.com/abulo/ratel/v3/util"
)

var MySQL *mysql.ProxyPool = mysql.NewProxyPool()

func main() {
	opt := &mysql.Config{}
	opt.Username = "root"
	opt.Password = "mysql"
	opt.Host = "127.0.0.1"
	opt.Port = "3306"
	opt.Charset = "utf8mb4"
	opt.Database = "xmt"
	opt.DriverName = "mysql"
	opt.MaxLifetime = time.Duration(1) * time.Minute
	opt.MaxIdleTime = time.Duration(1) * time.Minute
	opt.MaxIdleConns = util.ToInt(64)
	opt.MaxOpenConns = util.ToInt(64)

	conn := mysql.New(opt)

	mysql2struct.MysqlToStruct(conn, "xmt", "sss", "dddd")
}

```