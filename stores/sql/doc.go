// Package sql 提供基于 GORM 的关系型数据库客户端封装，支持 MySQL、ClickHouse 与 PostgreSQL。
//
// 该包通过 Client 配合 Option 完成连接配置，由 Dsn 与 Dialector 按 DriverName 生成对应
// 方言；SqlConn 按 DSN 缓存复用已建立的连接，可按需挂载指标采集与链路追踪拦截器。
// 包内同时提供 GORM 日志适配、分页/范围/排序参数辅助类型，以及可接受错误的判定工具。
package sql
