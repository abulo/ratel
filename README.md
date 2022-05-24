# Ratel 


[![Go](https://github.com/abulo/ratel/v3/workflows/Go/badge.svg?branch=master)](https://github.com/abulo/ratel/v3/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abulo/ratel/v3)](https://goreportcard.com/report/github.com/abulo/ratel/v3)
[![goproxy](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)
[![codecov](https://codecov.io/gh/abulo/ratel/branch/master/graph/badge.svg)](https://codecov.io/gh/abulo/ratel)
[![Release](https://img.shields.io/github/v/release/abulo/ratel.svg?style=flat-square)](https://github.com/abulo/ratel/v3)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


## Ratel 介绍

Ratel 一套微服务治理框架,集成了各种工程实践的 web 和 rpc 框架。

优点:
* 轻松获得支撑千万日活服务的稳定性
* 内建级联超时控制、限流、自适应熔断、自适应降载等微服务治理能力，无需配置和额外代码
* 微服务治理中间件可无缝集成到其它现有框架使用
* 大量微服务治理和并发工具包

## Ratel 目的

* 高效的性能
* 简洁的语法
* 广泛验证的工程效率
* 极致的部署体验
* 极低的服务端资源成本

## Ratel 思考

* 保持简单，第一原则
* 弹性设计，面向故障编程
* 工具大于约定和文档
* 高可用
* 高并发
* 易扩展
* 对业务开发友好，封装复杂度
* 约束做一件事只有一种方式


## Ratel 特性
* 统一的指标采集
* 链路追踪
* 日志埋点
* 统一错误处理
* 动态配置
* 安全策略
* Debug 模式 等，可以极大的提高应用开发效率
* 强大的工具支持，尽可能少的代码编写
* 极简的接口
* 完全兼容 net/http
* 支持中间件，方便扩展
* 高性能
* 面向故障编程，弹性设计
* 内建服务发现、负载均衡
* 内建限流、熔断、降载，且自动触发，自动恢复
* API 参数自动校验
* 超时级联控制
* 自动缓存控制
* 链路跟踪、统计报警等
* 高并发支撑，稳定保障了疫情期间每天的流量洪峰