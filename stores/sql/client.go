package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/core/resource"
	"gorm.io/gorm"
)

const (
	MySQLDriverName      = "mysql"
	ClickHouseDriverName = "clickhouse"
	PostgreSQLDriverName = "postgres"
	TDengineDriverName   = "tdengine"
)

type Client struct {
	Host           string        // 数据库 IP
	Port           string        // 数据库端口
	Username       string        // 数据库用户名
	Password       string        // 数据库密码
	Charset        string        // 数据库字符集
	Database       string        // 数据库名称
	ParseTime      bool          // 是否解析时间
	TimeZone       string        // 数据库时区  mysql & postgresql 专用
	SslMode        bool          // 数据库SSL模式  postgresql 专用
	DialTimeOut    string        // 连接超时时间 clickhouse 专用
	ReadTimeOut    string        // 读取超时时间 clickhouse 专用
	MaxIdleConns   int           // 连接池里最大空闲连接数。必须要比maxOpenConns小
	MaxOpenConns   int           // 连接池最多同时打开的连接数
	MaxLifetime    time.Duration // 连接池里面的连接最大存活时长
	MaxIdleTime    time.Duration // 连接池里面的连接最大空闲时长
	EnableMetric   bool          // 开启指标采集
	EnableTrace    bool          // 开启链路追踪
	EnableDebug    bool          // 关闭调试模式
	ConnectionMode string        // tdengine 连接模式(native/rest/websocket)
	DriverName     string        // 数据库驱动名称
}

type ClientManager struct {
	*gorm.DB
}

func (c *ClientManager) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *ClientManager) Gorm() *gorm.DB {
	return c.DB
}

// Option 选项
type Option func(r *Client)

var connManager = resource.NewResourceManager()

// NewClient 创建一个新的客户端
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// WithHost 设置数据库IP
func WithHost(host string) Option {
	return func(r *Client) {
		r.Host = host
	}
}

// WithPort 设置数据库端口
func WithPort(port string) Option {
	return func(r *Client) {
		r.Port = port
	}
}

// WithUsername 设置数据库用户名
func WithUsername(username string) Option {
	return func(r *Client) {
		r.Username = username
	}
}

// WithPassword 设置数据库密码
func WithPassword(password string) Option {
	return func(r *Client) {
		r.Password = password
	}
}

// WithCharset 设置数据库字符集
func WithCharset(charset string) Option {
	return func(r *Client) {
		r.Charset = charset
	}
}

// WithDatabase 设置数据库名称
func WithDatabase(database string) Option {
	return func(r *Client) {
		r.Database = database
	}
}

// WithParseTime 是否解析时间
func WithParseTime(parseTime bool) Option {
	return func(r *Client) {
		r.ParseTime = parseTime
	}
}

// WithTimeZone 设置数据库时区
func WithTimeZone(timeZone string) Option {
	return func(r *Client) {
		r.TimeZone = timeZone
	}
}

// WithSslMode 设置数据库SSL模式
func WithSslMode(sslMode bool) Option {
	return func(r *Client) {
		r.SslMode = sslMode
	}
}

// WithDialTimeOut 设置连接超时时间
func WithDialTimeOut(dialTimeOut string) Option {
	return func(r *Client) {
		r.DialTimeOut = dialTimeOut
	}
}

// WithReadTimeOut 设置读取超时时间
func WithReadTimeOut(readTimeOut string) Option {
	return func(r *Client) {
		r.ReadTimeOut = readTimeOut
	}
}

// WithMaxIdleConns 设置连接池里最大空闲连接数
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(r *Client) {
		r.MaxIdleConns = maxIdleConns
	}
}

// WithMaxOpenConns 设置连接池最多同时打开的连接数
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(r *Client) {
		r.MaxOpenConns = maxOpenConns
	}
}

// WithMaxLifetime 设置连接池里面的连接最大存活时长
func WithMaxLifetime(maxLifetime time.Duration) Option {
	return func(r *Client) {
		r.MaxLifetime = maxLifetime
	}
}

// WithMaxIdleTime 设置连接池里面的连接最大空闲时长
func WithMaxIdleTime(maxIdleTime time.Duration) Option {
	return func(r *Client) {
		r.MaxIdleTime = maxIdleTime
	}
}

// WithEnableMetric 开启指标采集
func WithEnableMetric(disableMetric bool) Option {
	return func(r *Client) {
		r.EnableMetric = disableMetric
	}
}

// WithEnableTrace 开启链路追踪
func WithEnableTrace(disableTrace bool) Option {
	return func(r *Client) {
		r.EnableTrace = disableTrace
	}
}

// WithDriverName 设置数据库驱动名称
func WithDriverName(driverName string) Option {
	return func(r *Client) {
		r.DriverName = driverName
	}
}

// WithEnableDebug 关闭调试模式
func WithEnableDebug(disableDebug bool) Option {
	return func(r *Client) {
		r.EnableDebug = disableDebug
	}
}

func (c *Client) Dsn() (string, error) {
	switch c.DriverName {
	case MySQLDriverName:
		return c.MySqlDsn(), nil
	case ClickHouseDriverName:
		return c.ClickHouseDsn(), nil
	case PostgreSQLDriverName:
		return c.PostgreSQLDsn(), nil
	case TDengineDriverName:
		return c.TDengineDsn(), nil
	default:
		return "", fmt.Errorf("driverName not support : %s", c.DriverName)
	}
}

// Dsn 数据库连接
func (c *Client) Dialector() (gorm.Dialector, error) {
	switch c.DriverName {
	case MySQLDriverName:
		return c.MySqlDialector()
	case ClickHouseDriverName:
		return c.ClickHouseDialector()
	case PostgreSQLDriverName:
		return c.PostgreSQLDialector()
	case TDengineDriverName:
		return c.TDengineDialector()
	default:
		return nil, fmt.Errorf("driverName not support : %s", c.DriverName)
	}
}

// Open 打开数据库连接
func (c *Client) Open() (*ClientManager, error) {
	dialector, err := c.Dialector()
	if err != nil {
		return nil, err
	}
	newLogger := NewLogger(logger.Logger)
	newLogger.SetDebug(c.EnableDebug)
	newLogger.SetSourceField("gorm")
	newLogger.SetSkipErrRecordNotFound(true)
	db, err := gorm.Open(dialector, &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, err
	}

	// 这里需要去判断要不要开启链路追踪和指标采集
	interceptors := []Interceptor{}
	if c.EnableMetric {
		interceptors = append(interceptors, MetricInterceptor())
	}
	if c.EnableTrace {
		interceptors = append(interceptors, TraceInterceptor())
	}
	if len(interceptors) > 0 {
		RegisterInterceptor(db, c, interceptors...)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(c.MaxLifetime)
	sqlDB.SetConnMaxIdleTime(c.MaxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	client := &ClientManager{DB: db}
	return client, nil
}

// GetConn 获取数据库连接
func (c *Client) SqlConn() (*gorm.DB, error) {
	conn, err := c.getCachedSqlConn()
	if err != nil {
		return nil, err
	}

	return conn.Gorm(), nil
}

// getCachedSqlConn 获取缓存的数据库连接
func (c *Client) getCachedSqlConn() (*ClientManager, error) {
	dsn, err := c.Dsn()
	if err != nil {
		return nil, err
	}
	val, err := connManager.GetResource(dsn, func() (io.Closer, error) {
		conn, err := c.Open()
		if err != nil {
			return nil, err
		}

		return conn, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*ClientManager), nil
}

var (
	ErrLastInsertIdIsNotSupported = errors.New("LastInsertId is not supported by this driver")
	ErrRowsAffectedIsNotSupported = errors.New("RowsAffected is not supported by this driver")
	ErrNoLastInsertIdAvailable    = errors.New("no LastInsertId available after DDL statement")
	ErrNoRowsAffectedAvailable    = errors.New("no RowsAffected available after DDL statement")
)

func Acceptable(err error) error {
	if err == nil || ErrorIn(err, sql.ErrNoRows, sql.ErrTxDone, context.Canceled, ErrLastInsertIdIsNotSupported, ErrRowsAffectedIsNotSupported, ErrNoLastInsertIdAvailable, ErrNoRowsAffectedAvailable, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// In checks if the given err is one of errs.
func ErrorIn(err error, errs ...error) bool {
	for _, each := range errs {
		if errors.Is(err, each) {
			return true
		}
	}
	return false
}
