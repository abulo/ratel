// Package mongodb 提供基于 mongo-driver 的 MongoDB 客户端封装与集合操作增强。
//
// 该包通过 NewMongoDBClient 按 URI 缓存复用客户端，并以 Option 配置超时与连接池参数；
// Model 与 DecoratedCollection 在原生集合之上提供链式条件构造、故障保护、指标采集
// 与链路追踪；BeforeCreate、BeforeUpdate 负责在写入前完成结构体到 BSON 文档的转换与 _id 生成。
package mongodb
