// Package xgin 提供基于 Gin 的 HTTP 服务端封装。
//
// 该包通过 Config 以函数式选项配置监听地址、运行模式以及指标、链路追踪、慢查询等中间件，
// 由 Build 创建实现 server.Server 接口的 Server；同时提供 HTTP 错误类型、状态码常量、
// WebSocket 升级支持以及 gRPC 响应转发所需的 Protobuf 包装消息。
package xgin
