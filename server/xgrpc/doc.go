// Package xgrpc 提供基于 gRPC 的服务端封装。
//
// 该包通过 Config 以函数式选项配置监听地址、网络、TLS 与一元/流式拦截器，由 Build 或
// MustBuild 创建实现 server.Server 接口的 Server；内置服务端链路追踪与指标采集拦截器，
// 并提供拦截器串联能力。
package xgrpc
