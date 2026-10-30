// Package task 提供分布式定时任务管理能力，基于 cron 调度器与一致性哈希节点池在集群内分配任务。
//
// 通过 NewTask 或 NewTaskWithOption 构造 Task，注册 Job 或函数后调用 Start、Run 启动；
// 节点池在集群成员变更时自动进入升级状态以避免任务重复执行，并支持启动任务恢复与近期任务回放。
package task
