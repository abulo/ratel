// Package xhertz 提供基于 CloudWeGo Hertz 的 HTTP 服务端封装。
//
// 该包通过 Config 以函数式选项配置监听地址、运行模式以及指标、链路追踪、慢查询等中间件，
// 由 Build 创建实现 server.Server 接口的 Server，并提供 HTTP 头部与 CGI 变量名互转等辅助函数。
package xhertz
