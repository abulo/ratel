// Package logger 提供基于 logrus 的进程内全局日志器与日志入口。
//
// 该包维护共享的默认 Logger 与 LoggerEntry，供业务代码直接调用；
// 将日志异步写入 Elasticsearch、MongoDB、SQL 数据库的钩子由本包的
// es、mongo、sql 子包提供，异步队列实现见 queue 子包。
package logger
