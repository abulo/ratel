// Package manager 提供基于 Asynq 的定时任务管理器封装。
//
// 该包通过 Config 配置 PeriodicTaskManager 运行选项，由 Build 创建实现 server.Server
// 接口的 Server，按已注册的定时计划周期性将任务写入 Redis 队列，交由任务工作者消费。
package manager
