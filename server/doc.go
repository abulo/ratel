// Package server 定义各类服务端通用的抽象契约与服务注册信息模型。
//
// 该包提供 Server 生命周期接口（启动、立即停止、优雅停止、信息查询、健康检查），
// 以及 ServiceInfo、Service、Route 等注册描述结构和函数式配置选项，供 HTTP、gRPC、
// Hertz、异步任务等具体服务端实现复用。
package server
