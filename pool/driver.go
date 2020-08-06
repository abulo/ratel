package pool

import (
	"io"
	"sync"
	"time"
)

type driverConn struct {
	db *DB
	mu sync.Mutex

	createdAt time.Time //连接创建时间
	inUse     bool      //是否被使用
	ci        io.Closer //连接接口
}

// 判断活跃时间是否到期
func (dc *driverConn) expired(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return dc.createdAt.Add(timeout).Before(time.Now())
}

// 关闭资源
func (dc *driverConn) close() error {
	var err error

	dc.mu.Lock()
	if dc.ci == nil {
		return nil
	}
	err = dc.ci.Close()
	dc.mu.Unlock()

	dc.db.Lock()
	dc.db.numOpen--
	dc.db.maybeOpenNewConnections()
	dc.db.Unlock()

	return err
}

// 获取真实连接Connect返回的值，连接池通用性
func (dc *driverConn) Conn() io.Closer {
	return dc.ci
}

// 回收资源
func (dc *driverConn) Close() error {
	dc.db.Lock()
	if !dc.db.recovery(dc) {
		dc.db.Unlock()
		return dc.close()
	}
	dc.db.Unlock()
	return nil
}
