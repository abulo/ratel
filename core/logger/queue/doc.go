// Package queue 提供可复用的异步任务队列，供日志钩子在后台异步提交日志条目。
//
// 包内提供基于带缓冲通道的 Queue 与基于双向链表的 ListQueue 两种实现，
// 以及异步任务 Jober、可等待结果的同步任务 SyncJober 接口；队列由固定
// 数量的 worker 协程消费任务，队列未运行时 Push 会静默丢弃任务。
package queue
