// Package proxy 提供按分组注册与获取各类数据存储后端的代理池。
//
// Proxy 以分组名为键维护 Redis、MongoDB、Elasticsearch 与 SQL 代理实例；各代理包装
// 对应存储客户端，SQL 代理额外支持写库与读库列表及随机选库。实例一般在初始化阶段装配，
// 运行期通过 Load 系列方法按分组取出使用。
package proxy
