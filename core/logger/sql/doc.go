// Package sql 提供基于 GORM 的 logrus 日志钩子实现。
//
// 钩子将日志条目映射为 Dao 模型后，通过后台队列异步写入指定数据表；
// 写入失败时错误信息输出到配置的错误输出目标。
package sql
