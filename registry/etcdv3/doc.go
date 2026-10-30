// Package etcdv3 提供基于 etcd v3 的服务注册与发现实现，实现 registry.Registry 接口。
//
// 该包通过租约（lease）与 keepalive 维持服务节点的存活，支持节点注册、注销、
// 列表查询与前缀监听；Config 同时提供普通构造与按节点单例构造两种方式。
package etcdv3
