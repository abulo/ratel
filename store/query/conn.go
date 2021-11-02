package query

import (
	"context"
	"database/sql"
	"io"
	"sync"
)

var connManager = NewResourceManager()

type pingedDB struct {
	*sql.DB
	once sync.Once
}

func getCachedSqlConn(driverName, server string, opt *Opt) (*pingedDB, error) {
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		conn, err := newDBConnection(driverName, server, opt)
		if err != nil {
			return nil, err
		}

		return &pingedDB{
			DB: conn,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*pingedDB), nil
}

func NewSqlConn(driverName, server string, opt *Opt) (*sql.DB, error) {
	pdb, err := getCachedSqlConn(driverName, server, opt)
	if err != nil {
		return nil, err
	}

	pdb.once.Do(func() {
		err = pdb.PingContext(context.TODO())
	})
	if err != nil {
		return nil, err
	}

	return pdb.DB, nil
}

func newDBConnection(driverName, datasource string, opt *Opt) (*sql.DB, error) {
	conn, err := sql.Open(driverName, datasource)
	if err != nil {
		return nil, err
	}

	// MaxOpenConns int           //连接池最多同时打开的连接数
	// MaxIdleConns int           //连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxLifetime  time.Duration //连接池里面的连接最大存活时长
	// MaxIdleTime  time.Duration //连接池里面的连接最大空闲时长

	conn.SetMaxOpenConns(opt.MaxOpenConns)
	conn.SetMaxIdleConns(opt.MaxIdleConns)
	conn.SetConnMaxLifetime(opt.MaxLifetime)
	conn.SetConnMaxIdleTime(opt.MaxIdleTime)

	// SetMaxOpenConns(maxOpenConns)
	// 连接池最多同时打开的连接数。

	// 这个maxOpenConns理应要设置得比mysql服务器的max_connections值要小。

	// 一般设置为： 服务器cpu核心数 * 2 + 服务器有效磁盘数。参考这里

	// 可用show variables like 'max_connections'; 查看服务器当前设置的最大连接数。

	// SetMaxIdleConns(maxIdleConns)
	// 连接池里最大空闲连接数。必须要比maxOpenConns小；

	// SetConnMaxIdleTime(maxIdleTime)
	// 连接池里面的连接最大空闲时长。

	// 当连接持续空闲时长达到maxIdleTime后，该连接就会被关闭并从连接池移除，哪怕当前空闲连接数已经小于SetMaxIdleConns(maxIdleConns)设置的值。

	// 连接每次被使用后，持续空闲时长会被重置，从0开始从新计算；

	// 用show processlist; 可用查看mysql服务器上的连接信息，Command表示连接的当前状态，Command为Sleep时表示休眠、空闲状态，Time表示此状态的已持续时长

	// SetConnMaxLifetime(maxLifeTime)
	// 连接池里面的连接最大存活时长。

	// maxLifeTime必须要比mysql服务器设置的wait_timeout小，否则会导致golang侧连接池依然保留已被mysql服务器关闭了的连接。

	// mysql服务器的wait_timeout默认是8 hour，可通过show variables like 'wait_timeout'查看。

	return conn, nil
}
