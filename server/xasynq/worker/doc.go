// Package worker 提供基于 Asynq 的后台任务工作者封装。
//
// 该包通过 Config 配置 Redis 客户端、asynq 运行选项与任务处理器路由表，由 Build 创建
// 实现 server.Server 接口的 Server，负责从 Redis 队列拉取任务并按 Mux 注册的处理器消费执行。
package worker
