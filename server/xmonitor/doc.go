// Package xmonitor 提供独立的监控与治理 HTTP 服务端。
//
// 该包通过 Config 配置监听地址，Build 创建实现 server.Server 接口的 Server；服务端在
// 独立路由表上挂载 pprof 性能分析、已注册路由清单与模块构建信息等内置治理接口。
package xmonitor
