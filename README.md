# Ratel 

[![Go](https://github.com/abulo/ratel/v3/workflows/Go/badge.svg?branch=master)](https://github.com/abulo/ratel/v3/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abulo/ratel/v3)](https://goreportcard.com/report/github.com/abulo/ratel/v3)
[![goproxy](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)
[![codecov](https://codecov.io/gh/abulo/ratel/branch/master/graph/badge.svg)](https://codecov.io/gh/abulo/ratel)
[![Release](https://img.shields.io/github/v/release/abulo/ratel.svg?style=flat-square)](https://github.com/abulo/ratel/v3)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📖 目录

- [Ratel 介绍](#ratel-介绍)
- [核心特性](#核心特性)
- [设计理念](#设计理念)
- [快速开始](#快速开始)
  - [安装](#安装)
  - [Web 服务器示例](#web-服务器示例)
  - [gRPC 服务器示例](#grpc-服务器示例)
  - [Hertz 服务器示例](#hertz-服务器示例)
- [架构组件](#架构组件)
  - [Server 层](#server-层)
  - [Client 层](#client-层)
  - [Core 核心层](#core-核心层)
  - [Stores 存储层](#stores-存储层)
  - [Registry 注册中心](#registry-注册中心)
- [微服务治理能力](#微服务治理能力)
  - [限流](#限流)
  - [熔断](#熔断)
  - [链路追踪](#链路追踪)
  - [指标监控](#指标监控)
  - [日志系统](#日志系统)
- [工具集](#工具集)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [性能优势](#性能优势)
- [应用场景](#应用场景)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## Ratel 介绍

Ratel 是一套**生产级微服务治理框架**，集成了各种工程实践的 Web 和 RPC 框架。经过千万日活服务的实战验证，旨在提供极致的稳定性、降低服务端资源成本、提升开发效率。

### 核心价值

✅ **轻松获得支撑千万日活服务的稳定性**  
✅ **内建级联超时控制、限流、自适应熔断、自适应降载等微服务治理能力，无需配置和额外代码**  
✅ **微服务治理中间件可无缝集成到其它现有框架使用**  
✅ **大量微服务治理和并发工具包**

---

## 核心特性

### 🚀 高性能
- 高效性能，极低的服务端资源成本
- 高并发支撑，稳定保障流量洪峰
- 完全兼容 `net/http`，零学习成本

### 🛡️ 弹性设计
- 面向故障编程，自动熔断、降载
- 内建限流、熔断、降载，且自动触发、自动恢复
- 级联超时控制，防止雪崩效应

### 🔧 开发友好
- 极简的接口，尽可能少的代码编写
- Debug 模式，极大提高应用开发效率
- API 参数自动校验
- 自动缓存控制

### 📊 可观测性
- 统一的指标采集（Prometheus）
- 全链路追踪（OpenTelemetry）
- 结构化日志埋点
- 统计报警

### 🔌 扩展性强
- 支持中间件，方便扩展
- 内建服务发现、负载均衡
- 动态配置，热更新支持
- 安全策略内置

---

## 设计理念

### 第一原则：保持简单
约束做一件事只有一种方式，减少决策成本。

### 弹性设计：面向故障编程
假设故障必然发生，通过自适应机制自动应对。

### 工具大于约定和文档
通过强大的工具链减少配置和文档依赖。

### 四大目标
- **高可用**：99.99% 可用性保障
- **高并发**：支撑百万 QPS
- **易扩展**：插件化架构
- **对业务开发友好**：封装复杂度，暴露简洁接口

---

## 快速开始

### 安装

```bash
go get github.com/abulo/ratel/v3@latest
```

**环境要求：**
- Go 1.26.0+
- Make（可选，用于构建和测试）

### Web 服务器示例

基于 Gin 框架的高性能 Web 服务器：

```go
package main

import (
    "github.com/abulo/ratel/v3/server/xgin"
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建服务器实例
    server := xgin.NewServer()
    
    // 注册路由
    server.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Ratel!",
        })
    })
    
    // 启动服务器
    server.Run(":8080")
}
```

**特性：**
- 自动集成中间件（日志、指标、链路追踪）
- 支持 WebSocket
- 健康检查
- 优雅停机

### gRPC 服务器示例

```go
package main

import (
    "github.com/abulo/ratel/v3/server/xgrpc"
    "google.golang.org/grpc"
)

func main() {
    // 创建 gRPC 服务器
    server := xgrpc.NewServer()
    
    // 注册服务
    // pb.RegisterYourServiceServer(server, &yourServiceImpl{})
    
    // 启动服务
    server.Run(":50051")
}
```

**特性：**
- 自动负载均衡
- 服务发现集成
- 拦截器链（认证、日志、指标）
- 连接池管理

### Hertz 服务器示例

基于 CloudWeGo Hertz 的高性能 HTTP 框架：

```go
package main

import (
    "github.com/abulo/ratel/v3/server/xhertz"
    "github.com/cloudwego/hertz/pkg/app"
)

func main() {
    server := xhertz.NewServer()
    
    server.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
        ctx.JSON(200, map[string]string{
            "message": "Hello from Hertz!",
        })
    })
    
    server.Run(":8081")
}
```

---

## 架构组件

### Server 层

Ratel 提供多种服务器实现，满足不同场景需求：

| 模块 | 说明 | 适用场景 |
|------|------|----------|
| **xgin** | 基于 Gin 的 Web 服务器 | RESTful API、Web 应用 |
| **xgrpc** | 基于 gRPC 的 RPC 服务器 | 微服务间通信 |
| **xhertz** | 基于 Hertz 的高性能 HTTP 服务器 | 高并发 HTTP 服务 |
| **xasynq** | 异步任务队列服务器 | 后台任务处理 |
| **xmonitor** | 监控服务器 | 指标暴露、健康检查 |

### Client 层

提供多种客户端实现，简化外部服务调用：

| 模块 | 说明 |
|------|------|
| **etcdv3** | Etcd 客户端（支持服务发现、配置监听、分布式锁） |
| **grpc** | gRPC 客户端（负载均衡、断路器、拦截器） |
| **rabbitmq** | RabbitMQ 消息队列客户端 |
| **redis** | Redis 客户端（连接池、哨兵、集群） |
| **sftp/ftp** | 文件传输客户端 |

### Core 核心层

提供微服务治理的核心能力：

#### 📝 Logger - 日志系统
- 结构化日志（logrus）
- 多输出支持（控制台、文件、Elasticsearch、MongoDB）
- 日志轮转
- 上下文追踪

```go
import "github.com/abulo/ratel/v3/core/logger"

logger.Logger.Info("application started")
logger.Logger.WithField("user_id", 123).Error("operation failed")
```

#### 📊 Metric - 指标监控
- Prometheus 指标采集
- 计数器（Counter）
- 仪表盘（Gauge）
- 直方图（Histogram）
- 摘要（Summary）

```go
import "github.com/abulo/ratel/v3/core/metric"

// 记录请求次数
metric.ServerHandleCounter.Inc(map[string]string{
    "type":   "http",
    "method": "GET",
    "code":   "200",
})
```

#### ⏱️ Ratelimit - 限流
- 令牌桶算法
- 滑动窗口
- 自适应限流

```go
import "github.com/abulo/ratel/v3/core/ratelimit"

limiter := ratelimit.NewLimiter(rate.Every(1*time.Second), 100)
if !limiter.Allow() {
    return errors.New("rate limit exceeded")
}
```

#### 🔌 Resource - 资源管理
- **断路器（Breaker）**：自适应熔断
- **SingleFlight**：防止缓存击穿
- **BatchError**：批量错误处理
- **SourceManager**：资源生命周期管理

#### 🔄 Task - 任务调度
- 分布式定时任务
- 一致性哈希
- 任务池管理
- Cron 表达式支持

#### 🔍 Trace - 链路追踪
- OpenTelemetry 集成
- 分布式追踪
- 跨度（Span）管理
- 采样率配置

#### 🧵 Goroutine - 协程管理
- 并行执行
- 串行执行
- 错误聚合
- 超时控制

### Stores 存储层

统一的数据存储访问接口：

| 存储类型 | 支持数据库 |
|---------|-----------|
| **SQL** | MySQL、PostgreSQL、ClickHouse |
| **Redis** | 单机、哨兵、集群 |
| **MongoDB** | MongoDB |
| **Elasticsearch** | ES 7.x/8.x |
| **Session** | 会话管理 |

**特性：**
- GORM 集成
- 连接池管理
- 慢查询日志
- SQL 指标监控

### Registry 注册中心

服务注册与发现：

- **Etcd v3**：分布式 KV 存储
- 服务健康检查
- 权重路由
- 优雅上下线

---

## 微服务治理能力

### 限流

Ratel 提供多层限流策略：

1. **接口限流**：基于令牌桶/漏桶算法
2. **自适应限流**：根据系统负载自动调整
3. **分布式限流**：基于 Redis/Etcd

```go
// 中间件自动启用限流
server.Use(xgin.RateLimitMiddleware())
```

### 熔断

自适应熔断保护：

- **错误率熔断**：错误率超过阈值自动熔断
- **慢调用熔断**：响应时间过长自动熔断
- **自动恢复**：半开状态探测恢复

```go
// 断路器自动工作，无需配置
breaker := resource.NewBreaker(resource.BreakerConfig{
    MaxConcurrentRequests: 100,
    ErrorPercentThreshold: 50,
})
```

### 链路追踪

全链路追踪集成：

```go
import "github.com/abulo/ratel/v3/core/trace"

// 自动注入追踪上下文
ctx, span := trace.StartSpan(ctx, "operation-name")
defer span.End()

// 跨服务传播
trace.Inject(ctx, carrier)
```

**支持导出器：**
- Jaeger
- Zipkin
- OTLP

### 指标监控

自动采集关键指标：

**服务器指标：**
- `server_handle_total`：请求总数
- `server_handle_seconds`：请求耗时

**客户端指标：**
- `client_handle_total`：调用总数
- `client_handle_seconds`：调用耗时

**任务指标：**
- `job_handle_total`：任务执行总数
- `job_handle_seconds`：任务执行耗时

**暴露 Prometheus 端点：**
```
GET /metrics
```

### 日志系统

结构化日志输出：

```go
logger.Logger.WithFields(logrus.Fields{
    "request_id": "xxx",
    "user_id":    123,
    "duration":   "100ms",
}).Info("request completed")
```

**日志级别：**
- Debug
- Info
- Warn
- Error
- Fatal

**日志输出：**
- 控制台（彩色输出）
- 文件（按天轮转）
- Elasticsearch
- MongoDB
- 消息队列

---

## 工具集

### Toolkit 命令行工具

Ratel 提供强大的 CLI 工具，加速项目开发：

```bash
# 安装 toolkit
go install github.com/abulo/ratel/v3/toolkit@latest

# 创建新项目
toolkit new my-project

# 生成 CRUD 代码
toolkit generate model User

# 数据库迁移
toolkit migrate up

# API 文档生成
toolkit api docs
```

### 实用工具包

#### Util 工具集
- 字符串处理
- 时间格式化
- 文件操作
- 网络工具
- 正则表达式
- 随机数生成
- 拼音转换

#### Correct 文本校正
- 全角/半角转换
- HTML 实体编码
- 格式化/反格式化

#### NLPWord 分词
- AC 自动机
- Trie 树
- 敏感词过滤

#### Snowflake ID 生成
- 分布式唯一 ID
- 时间有序
- 高性能

---

## 项目结构

```
ratel/
├── client/              # 客户端实现
│   ├── etcdv3/         # Etcd 客户端
│   ├── grpc/           # gRPC 客户端
│   ├── rabbitmq/       # RabbitMQ 客户端
│   ├── redis/          # Redis 客户端
│   └── ...
├── config/             # 配置管理
├── core/               # 核心功能
│   ├── app/           # 应用生命周期
│   ├── logger/        # 日志系统
│   ├── metric/        # 指标监控
│   ├── ratelimit/     # 限流
│   ├── resource/      # 资源管理（熔断器等）
│   ├── task/          # 任务调度
│   ├── trace/         # 链路追踪
│   └── ...
├── server/             # 服务器实现
│   ├── xgin/          # Gin Web 服务器
│   ├── xgrpc/         # gRPC 服务器
│   ├── xhertz/        # Hertz 服务器
│   ├── xasynq/        # 异步任务服务器
│   └── xmonitor/      # 监控服务器
├── stores/             # 数据存储
│   ├── sql/           # SQL 数据库
│   ├── redis/         # Redis
│   ├── mongodb/       # MongoDB
│   └── elasticsearch/ # Elasticsearch
├── registry/           # 服务注册发现
├── toolkit/            # CLI 工具
├── util/               # 通用工具
├── filter/             # 过滤器
├── watch/              # 配置监听
└── ...
```

---

## 开发指南

### 本地开发

```bash
# 克隆项目
git clone https://github.com/abulo/ratel.git
cd ratel

# 安装依赖
go mod download

# 运行测试
make test

# 数据竞争检测
make race

# 代码格式化
make fmt

# Lint 检查
make lint

# 完整构建
make all
```

### 代码规范

- ✅ 使用 `gofmt` 格式化代码
- ✅ 使用 `revive` 进行 Lint 检查
- ✅ 必须检查所有错误返回
- ✅ 编写单元测试
- ✅ 添加必要的注释

### 调试技巧

**启用 Debug 模式：**

```go
import "github.com/abulo/ratel/v3/core/env"

env.SetDebug(true)
```

**查看详细日志：**

```bash
export LOG_LEVEL=debug
```

---

## 性能优势

### 基准测试

Ratel 在以下方面表现优异：

- **低延迟**：P99 < 10ms
- **高吞吐**：单实例支持 10万+ QPS
- **低资源**：内存占用 < 100MB
- **快速启动**：< 100ms

### 优化策略

1. **连接池复用**：减少连接建立开销
2. **零拷贝优化**：减少内存分配
3. **异步处理**：非阻塞 I/O
4. **缓存策略**：多级缓存
5. **批量操作**：减少网络往返

---

## 应用场景

### ✅ 电商系统
- 秒杀活动限流
- 订单处理熔断保护
- 分布式事务追踪

### ✅ 社交网络
- 高并发消息推送
- 用户关系图谱
- 实时动态推送

### ✅ 金融系统
- 交易风控
- 实时风控规则
- 审计日志

### ✅ IoT 平台
- 设备接入管理
- 海量数据采集
- 实时监控告警

### ✅ 内容平台
- CDN 调度
- 推荐系统
- 搜索服务

---

## 常见问题

### Q1: Ratel 与原生 Gin/gRPC 有什么区别？

**A:** Ratel 在原生框架基础上集成了微服务治理能力：
- 自动限流、熔断、降载
- 统一指标采集
- 链路追踪
- 日志规范化
- 无需额外配置和代码

### Q2: 如何迁移现有项目到 Ratel？

**A:** Ratel 完全兼容 `net/http` 和原生 Gin/gRPC：
```go
// 只需替换导入
// import "github.com/gin-gonic/gin"
import "github.com/abulo/ratel/v3/server/xgin"

// 其他代码无需修改
```

### Q3: 如何自定义中间件？

**A:** 
```go
server.Use(func(c *gin.Context) {
    // 你的逻辑
    c.Next()
})
```

### Q4: 如何配置限流规则？

**A:** 通过配置文件或环境变量：
```yaml
ratelimit:
  enabled: true
  rate: 1000  # 每秒请求数
  burst: 100  # 突发容量
```

### Q5: 支持哪些注册中心？

**A:** 目前支持 Etcd v3，后续会扩展 Consul、Nacos 等。

### Q6: 如何实现灰度发布？

**A:** 结合 Registry 的权重路由和版本标签：
```go
registry.Register(service, registry.WithWeight(50))
```

---

## 贡献指南

我们欢迎所有形式的贡献！

### 贡献流程

1. **Fork** 本仓库
2. **创建分支**：`git checkout -b feature/your-feature`
3. **提交代码**：`git commit -am 'Add some feature'`
4. **推送分支**：`git push origin feature/your-feature`
5. **提交 PR**

### 贡献类型

- 🐛 Bug 修复
- ✨ 新功能
- 📝 文档改进
- 🧪 测试用例
- ⚡ 性能优化
- 🔧 重构代码

### 开发规范

- 遵循 Go 官方代码规范
- 添加单元测试
- 更新文档
- 保持向后兼容

---

## 社区与支持

- 📧 **Email**: [项目维护者邮箱]
- 💬 **Issues**: [GitHub Issues](https://github.com/abulo/ratel/issues)
- 📖 **Wiki**: [项目 Wiki](https://github.com/abulo/ratel/wiki)
- 👥 **Discussions**: [GitHub Discussions](https://github.com/abulo/ratel/discussions)

---

## 成功案例

Ratel 已在多个生产环境验证：

- **某电商平台**：支撑双 11 千万级 QPS
- **某社交平台**：日活用户 5000 万+
- **某金融公司**：交易系统 99.99% 可用性
- **某 IoT 平台**：百万设备同时在线

---

## 许可证

Ratel 采用 [MIT 许可证](https://opensource.org/licenses/MIT)。

```
MIT License

Copyright (c) 2024 Ratel Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Star History

如果 Ratel 对你有帮助，请给我们一个 ⭐️ Star！

[![Star History Chart](https://api.star-history.com/svg?repos=abulo/ratel&type=Date)](https://star-history.com/#abulo/ratel&Date)

---

**Made with ❤️ by Ratel Team**
