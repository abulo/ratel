package proxy

import (
	"github.com/abulo/ratel/v3/util"
	"gorm.io/gorm"
)

// SQL 代理持有写库与读库两组 *gorm.DB，运行期随机选择目标节点，用于读写分离场景。
//
// 两组切片应在初始化阶段一次性装配完成；切片未做并发保护，运行期不应并发调用
// SetWrite、SetRead。取库时写库列表必须非空，否则 Write 会在随机选节点时 panic；
// 读库列表为空时 Read 回退到写库。
type SQL struct {
	write []*gorm.DB
	read  []*gorm.DB
}

// NewSQL 构造一个空的 SQL 代理，随后通过 SetWrite、SetRead 装配写库与读库实例。
func NewSQL() *SQL {
	return &SQL{}
}

// SetWrite 追加一个写库实例；须在对外提供服务前完成装配。
func (proxy *SQL) SetWrite(db *gorm.DB) {
	proxy.write = append(proxy.write, db)
}

// SetRead 追加一个读库实例；须在对外提供服务前完成装配。
func (proxy *SQL) SetRead(db *gorm.DB) {
	proxy.read = append(proxy.read, db)
}

// Write 从写库列表中随机返回一个写库实例；写库列表为空时会触发 panic，调用方须保证已装配。
func (proxy *SQL) Write() *gorm.DB {
	length := len(proxy.write)
	// PERF-1-8 [P1] 单写库直接返回下标 0，跳过随机数选择；多写库仍随机负载均衡。
	if length == 1 {
		return proxy.write[0]
	}
	write := util.Rand(0, length-1)
	return proxy.write[write]
}

// Read 从读库列表中随机返回一个读库实例；读库列表为空时回退到 Write 选取写库。
func (proxy *SQL) Read() *gorm.DB {
	length := len(proxy.read)
	if length < 1 {
		return proxy.Write()
	}
	// PERF-1-8 [P1] 单读库直接返回下标 0，跳过随机数选择；多读库仍随机负载均衡。
	if length == 1 {
		return proxy.read[0]
	}
	read := util.Rand(0, length-1)
	return proxy.read[read]
}

// StoreSQL 将 SQL 代理按分组名注册到代理池，同名分组会被覆盖。
func (proxyPool *Proxy) StoreSQL(group string, proxy *SQL) {
	proxyPool.m.Store(group, proxy)
}

// LoadSQL 按分组名取出已注册的 SQL 代理；分组不存在时返回 nil。
func (proxyPool *Proxy) LoadSQL(group string) *SQL {
	if f, ok := proxyPool.m.Load(group); ok {
		return f.(*SQL)
	}
	return nil
}
