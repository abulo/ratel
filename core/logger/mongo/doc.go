// Package mongo 提供基于 MongoDB 的 logrus 日志钩子实现。
//
// 钩子将日志条目转换为 BSON 文档后，通过后台队列异步写入指定集合；
// 写入失败时错误信息输出到配置的错误输出目标。
package mongo
