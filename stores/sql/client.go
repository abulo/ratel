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

// 数据库驱动名称常量，供 Client.DriverName 选择 DSN 与方言构造分支。
const (
	// MySQLDriverName 为 MySQL 驱动标识。
	MySQLDriverName = "mysql"
	// ClickHouseDriverName 为 ClickHouse 驱动标识。
	ClickHouseDriverName = "clickhouse"
	// PostgreSQLDriverName 为 PostgreSQL 驱动标识。
	PostgreSQLDriverName = "postgres"
	// TDengineDriverName 为 TDengine 驱动标识；当前 Dsn 与 Dialector 尚未实现该分支。
	TDengineDriverName = "tdengine"
)

// PERF-2-3 [P1] 连接池兜底默认值：调用方未显式配置（零值）时使用，
// 避免 SetMaxIdleConns(0) 导致空闲连接用完即关、SetMaxOpenConns(0) 导致连接数无上限。
const (
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 100
	defaultConnMaxLifetime = time.Hour
	defaultConnMaxIdleTime = 30 * time.Minute
)

// Client 保存关系型数据库的全部连接配置，由各 With* Option 在 NewClient 阶段装配。
//
// 零值不可直接使用，须先通过 Option 填写 Host、Port、Username、Password、Database、
// DriverName 等必填项后再调用 Open；Client 为配置载体，应在初始化阶段一次性写入，
// 不保证并发修改安全。
type Client struct {
	// Host 为数据库主机地址（IP 或域名）。
	Host string
	// Port 为数据库服务端口。
	Port string
	// Username 为数据库登录用户名。
	Username string
	// Password 为数据库登录密码。
	Password string
	// Charset 为数据库字符集，仅 MySQL DSN 使用。
	Charset string
	// Database 为要连接的数据库名称。
	Database string
	// ParseTime 指定是否在 DSN 中开启时间解析，仅 MySQL DSN 使用。
	ParseTime bool
	// TimeZone 为数据库时区，MySQL 与 PostgreSQL DSN 使用。
	TimeZone string
	// SslMode 指定是否启用 SSL 连接，仅 PostgreSQL DSN 使用。
	SslMode bool
	// DialTimeOut 为建立连接的超时时间，仅 ClickHouse DSN 使用。
	DialTimeOut string
	// ReadTimeOut 为读取超时时间，仅 ClickHouse DSN 使用。
	ReadTimeOut string
	// MaxIdleConns 为连接池中空闲连接的最大数量，应小于 MaxOpenConns。
	MaxIdleConns int
	// MaxOpenConns 为连接池允许同时打开的最大连接数。
	MaxOpenConns int
	// MaxLifetime 为连接池中连接的最大存活时长，超过后会被回收重建。
	MaxLifetime time.Duration
	// MaxIdleTime 为连接池中连接的最大空闲时长，超过后会被关闭。
	MaxIdleTime time.Duration
	// EnableMetric 是否采集 SQL 操作指标。
	EnableMetric bool
	// EnableTrace 是否挂载 SQL 链路追踪拦截器。
	EnableTrace bool
	// EnableDebug 是否开启 GORM 调试日志输出。
	EnableDebug bool
	// DriverName 指定使用的数据库驱动，取值见本包驱动名称常量。
	DriverName string
	// Name 为连接实例名称，参与连接缓存的去重标识。
	Name string
}

// ClientManager 包装已建立连接的 *gorm.DB，对外提供关闭连接与获取底层实例的能力。
//
// 实例由 Open 创建并可按 DSN 缓存复用；底层 *gorm.DB 的并发安全由 GORM 自身保证；
// Close 关闭底层连接池后不应再复用该实例。
type ClientManager struct {
	*gorm.DB
}

// Close 关闭底层数据库连接池，返回底层驱动 Close 的错误；未建立连接时返回对应错误。
func (c *ClientManager) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Gorm 返回底层 *gorm.DB 实例，供直接执行 GORM 链式查询使用。
func (c *ClientManager) Gorm() *gorm.DB {
	return c.DB
}

// Option 是 Client 的配置选项函数，以函数式选项模式在 NewClient 中装配配置。
type Option func(r *Client)

var connManager = resource.NewResourceManager()

// NewClient 依据传入的 Option 构造 Client 配置对象。
//
// 该方法只完成配置装配，不立即建立数据库连接；实际连接在 Open 时才发起，
// 返回的 error 恒为 nil。
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// WithHost 设置数据库主机地址。
func WithHost(host string) Option {
	return func(r *Client) {
		r.Host = host
	}
}

// WithPort 设置数据库服务端口。
func WithPort(port string) Option {
	return func(r *Client) {
		r.Port = port
	}
}

// WithUsername 设置数据库登录用户名。
func WithUsername(username string) Option {
	return func(r *Client) {
		r.Username = username
	}
}

// WithPassword 设置数据库登录密码。
func WithPassword(password string) Option {
	return func(r *Client) {
		r.Password = password
	}
}

// WithCharset 设置数据库字符集，仅 MySQL DSN 使用。
func WithCharset(charset string) Option {
	return func(r *Client) {
		r.Charset = charset
	}
}

// WithDatabase 设置要连接的数据库名称。
func WithDatabase(database string) Option {
	return func(r *Client) {
		r.Database = database
	}
}

// WithParseTime 设置是否开启 DSN 中的时间解析开关，仅 MySQL DSN 使用。
func WithParseTime(parseTime bool) Option {
	return func(r *Client) {
		r.ParseTime = parseTime
	}
}

// WithTimeZone 设置数据库时区，MySQL 与 PostgreSQL DSN 使用。
func WithTimeZone(timeZone string) Option {
	return func(r *Client) {
		r.TimeZone = timeZone
	}
}

// WithSslMode 设置是否启用 SSL 连接，仅 PostgreSQL DSN 使用。
func WithSslMode(sslMode bool) Option {
	return func(r *Client) {
		r.SslMode = sslMode
	}
}

// WithDialTimeOut 设置建立连接的超时时间，仅 ClickHouse DSN 使用。
func WithDialTimeOut(dialTimeOut string) Option {
	return func(r *Client) {
		r.DialTimeOut = dialTimeOut
	}
}

// WithReadTimeOut 设置读取超时时间，仅 ClickHouse DSN 使用。
func WithReadTimeOut(readTimeOut string) Option {
	return func(r *Client) {
		r.ReadTimeOut = readTimeOut
	}
}

// WithMaxIdleConns 设置连接池中空闲连接的最大数量，应小于最大打开连接数。
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(r *Client) {
		r.MaxIdleConns = maxIdleConns
	}
}

// WithMaxOpenConns 设置连接池允许同时打开的最大连接数。
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(r *Client) {
		r.MaxOpenConns = maxOpenConns
	}
}

// WithMaxLifetime 设置连接池中连接的最大存活时长。
func WithMaxLifetime(maxLifetime time.Duration) Option {
	return func(r *Client) {
		r.MaxLifetime = maxLifetime
	}
}

// WithMaxIdleTime 设置连接池中连接的最大空闲时长。
func WithMaxIdleTime(maxIdleTime time.Duration) Option {
	return func(r *Client) {
		r.MaxIdleTime = maxIdleTime
	}
}

// WithEnableMetric 设置是否采集 SQL 操作指标，传入 true 开启。
func WithEnableMetric(disableMetric bool) Option {
	return func(r *Client) {
		r.EnableMetric = disableMetric
	}
}

// WithEnableTrace 设置是否挂载 SQL 链路追踪拦截器，传入 true 开启。
func WithEnableTrace(disableTrace bool) Option {
	return func(r *Client) {
		r.EnableTrace = disableTrace
	}
}

// WithDriverName 设置数据库驱动名称，取值见本包驱动名称常量。
func WithDriverName(driverName string) Option {
	return func(r *Client) {
		r.DriverName = driverName
	}
}

// WithName 设置连接实例名称，参与连接缓存的去重标识。
func WithName(Name string) Option {
	return func(r *Client) {
		r.Name = Name
	}
}

// WithEnableDebug 设置是否开启 GORM 调试日志输出，传入 true 开启。
func WithEnableDebug(disableDebug bool) Option {
	return func(r *Client) {
		r.EnableDebug = disableDebug
	}
}

// Dsn 按 DriverName 生成对应驱动的连接字符串。
//
// DriverName 为 mysql、clickhouse、postgres 时返回对应 DSN；其余驱动名称返回
// 包含驱动名的错误，当前不支持 TDengine。
func (c *Client) Dsn() (string, error) {
	switch c.DriverName {
	case MySQLDriverName:
		return c.MySqlDsn(), nil
	case ClickHouseDriverName:
		return c.ClickHouseDsn(), nil
	case PostgreSQLDriverName:
		return c.PostgreSQLDsn(), nil
	default:
		return "", fmt.Errorf("driverName not support : %s", c.DriverName)
	}
}

// Dialector 按 DriverName 生成 GORM 方言实例；驱动不受支持时返回 Dsn 产生的错误。
func (c *Client) Dialector() (gorm.Dialector, error) {
	switch c.DriverName {
	case MySQLDriverName:
		return c.MySqlDialector()
	case ClickHouseDriverName:
		return c.ClickHouseDialector()
	case PostgreSQLDriverName:
		return c.PostgreSQLDialector()
	default:
		return nil, fmt.Errorf("driverName not support : %s", c.DriverName)
	}
}

// Open 依据当前配置建立数据库连接，完成方言初始化、连接池参数配置与连通性 Ping 校验。
//
// 成功时返回包装 *gorm.DB 的 ClientManager；方言构造、GORM 初始化、连接池设置或
// Ping 任一环节失败都会返回错误，Ping 失败时会主动关闭已建立的连接。开启指标采集或
// 链路追踪时，会向 GORM 回调注册对应拦截器。
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
		// PERF-2-5 [P2] time.Now() 返回值已是本地时区，去掉冗余的 .Local() 归一化调用。
		NowFunc: func() time.Time {
			return time.Now()
		},
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, err
	}

	// 按配置决定是否注册指标采集与链路追踪拦截器。
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
	// PERF-2-3 [P1] 仅对未显式配置（零值）的连接池参数填充兜底默认值，显式传入的配置保持不变。
	maxIdleConns := c.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = defaultMaxIdleConns
	}
	maxOpenConns := c.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = defaultMaxOpenConns
	}
	maxLifetime := c.MaxLifetime
	if maxLifetime <= 0 {
		maxLifetime = defaultConnMaxLifetime
	}
	maxIdleTime := c.MaxIdleTime
	if maxIdleTime <= 0 {
		maxIdleTime = defaultConnMaxIdleTime
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	client := &ClientManager{DB: db}
	return client, nil
}

// SqlConn 获取按 DSN 缓存复用的数据库连接，返回底层 *gorm.DB。
//
// 首次获取时调用 Open 建立连接并缓存；DSN 构造失败或连接建立失败时返回错误。
func (c *Client) SqlConn() (*gorm.DB, error) {
	conn, err := c.getCachedSqlConn()
	if err != nil {
		return nil, err
	}

	return conn.Gorm(), nil
}

// getCachedSqlConn 按 DSN 从连接管理器获取缓存的 ClientManager，缓存未命中时调用 Open 创建。
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

// SQL 驱动不返回结果元数据时使用的错误变量，由 Acceptable 判定为可接受错误。
var (
	// ErrLastInsertIdIsNotSupported 表示当前驱动不支持获取 LastInsertId。
	ErrLastInsertIdIsNotSupported = errors.New("LastInsertId is not supported by this driver")
	// ErrRowsAffectedIsNotSupported 表示当前驱动不支持获取 RowsAffected。
	ErrRowsAffectedIsNotSupported = errors.New("RowsAffected is not supported by this driver")
	// ErrNoLastInsertIdAvailable 表示执行 DDL 语句后无法取得 LastInsertId。
	ErrNoLastInsertIdAvailable = errors.New("no LastInsertId available after DDL statement")
	// ErrNoRowsAffectedAvailable 表示执行 DDL 语句后无法取得 RowsAffected。
	ErrNoRowsAffectedAvailable = errors.New("no RowsAffected available after DDL statement")
)

// Acceptable 判定错误是否属于可忽略的常规情况：nil、上下文取消、记录不存在、驱动不返回
// 元数据等错误一律视为可接受并返回 nil；其余错误原样返回，交由调用方或日志层上报。
func Acceptable(err error) error {
	if err == nil || ErrorIn(err, sql.ErrNoRows, sql.ErrTxDone, context.Canceled, ErrLastInsertIdIsNotSupported, ErrRowsAffectedIsNotSupported, ErrNoLastInsertIdAvailable, ErrNoRowsAffectedAvailable, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// ErrorIn 判断 err 是否为 errs 中的任意一个错误，基于 errors.Is 进行包装错误匹配。
func ErrorIn(err error, errs ...error) bool {
	for _, each := range errs {
		if errors.Is(err, each) {
			return true
		}
	}
	return false
}
