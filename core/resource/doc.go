// Package resource 提供服务治理与资源管理相关的基础组件。
//
// 该包包含并发单飞（SingleFlight）、基于 Google SRE 过载保护模型的熔断器（Breaker）、
// 批量错误聚合、原子时长以及按名称管理可关闭资源的资源管理器，供各客户端与存储层复用。
package resource
