// Package crontab 提供基于 cron 表达式的定时任务调度器。
//
// 该包解析五或六字段 cron 表达式以及 @yearly、@monthly 等预定义描述符，
// 维护任务的上一次、下一次执行时间；后台 goroutine 按时间触发任务，支持任务
// 超时控制、错误记录、任务动态增删以及等待全部任务结束的优雅关闭。
package crontab
