// Package etcdv3 提供基于 go.etcd.io/etcd/client/v3 的 etcd 客户端封装。
//
// 该包以 Config 链式配置构建客户端，覆盖单键与前缀键值读写、多键批量事务查询、
// 分布式互斥锁、租约会话与前缀变更监听；支持 BasicAuth、TLS 证书配置、自动成员同步
// 与链路追踪拦截器，并提供按名称复用的单例构建入口。
package etcdv3
