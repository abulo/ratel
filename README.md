# Ratel

[![Go](https://github.com/abulo/ratel/v3/workflows/Go/badge.svg?branch=v3)](https://github.com/abulo/ratel/v3/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abulo/ratel/v3)](https://goreportcard.com/report/github.com/abulo/ratel/v3)
[![goproxy](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)
[![codecov](https://codecov.io/gh/abulo/ratel/branch/v3/graph/badge.svg)](https://codecov.io/gh/abulo/ratel)
[![Release](https://img.shields.io/github/v/release/abulo/ratel.svg?style=flat-square)](https://github.com/abulo/ratel/v3)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Ratel 是一套生产级 Go 微服务治理框架，整合 Web/RPC 服务端、数据存储客户端与限流、熔断、链路追踪、指标、日志等治理能力，开箱即用。

- **服务端**：复用 Gin、gRPC、Hertz、Asynq 等成熟框架，统一封装为 `server.Server` 生命周期接口。
- **客户端与存储**：统一封装 etcd、RabbitMQ、FTP/SFTP、Redis、SQL（MySQL/PostgreSQL/ClickHouse）、MongoDB、Elasticsearch 等中间件。
- **治理能力**：令牌桶限流、Google SRE 自适应熔断、SingleFlight 防击穿、OpenTelemetry 链路追踪、Prometheus 指标、多后端异步日志。
- **配套工具**：分布式任务调度、雪花 ID、敏感词过滤、文本排版、i18n、图像处理、代码生成脚手架、热编译。

项目组织参考 Kratos 的分层方式，使用 Go 1.27 泛型与标准库，注释遵循 Go 官方文档规范（每个库包均有 `doc.go` 包注释，导出 API 带 Godoc）。

---

## 目录

- [特性](#特性)
- [架构设计](#架构设计)
- [安装](#安装)
- [快速开始](#快速开始)
  - [Gin HTTP 服务](#1-gin-http-服务)
  - [gRPC 服务](#2-grpc-服务)
  - [Hertz 服务](#3-hertz-服务)
  - [组合多个服务端](#4-组合多个服务端)
- [应用生命周期](#应用生命周期)
- [服务端详解](#服务端详解)
  - [xgin](#xgin)
  - [xgrpc](#xgrpc)
  - [xhertz](#xhertz)
  - [xasynq](#xasynq)
  - [xmonitor](#xmonitor)
- [客户端详解](#客户端详解)
  - [etcdv3](#etcdv3)
  - [grpc 客户端](#grpc-客户端)
  - [rabbitmq](#rabbitmq)
  - [ftp / sftp](#ftp--sftp)
- [存储层详解](#存储层详解)
  - [redis](#redis)
  - [sql](#sql)
  - [mongodb](#mongodb)
  - [elasticsearch](#elasticsearch)
  - [null 可空类型](#null-可空类型)
  - [proxy 代理池](#proxy-代理池)
  - [session 会话存储](#session-会话存储)
- [治理能力详解](#治理能力详解)
  - [限流 ratelimit](#限流-ratelimit)
  - [熔断 resource](#熔断-resource)
  - [goroutine 编排](#goroutine-编排)
  - [生命周期 cycle 与 hooks](#生命周期-cycle-与-hooks)
- [可观测性](#可观测性)
  - [日志 logger](#日志-logger)
  - [指标 metric](#指标-metric)
  - [链路追踪 trace](#链路追踪-trace)
- [任务调度](#任务调度)
  - [crontab 调度器](#crontab-调度器)
  - [分布式任务 core/task](#分布式任务-coretask)
  - [分布式执行器 executor](#分布式执行器-executor)
  - [worker 与 job 接口](#worker-与-job-接口)
- [工具包](#工具包)
  - [util](#util)
  - [snowflake](#snowflake)
  - [correct](#correct)
  - [nlpword / filter](#nlpword--filter)
  - [i18n](#i18n)
  - [img](#img)
  - [watch 热编译](#watch-热编译)
  - [config 配置读取](#config-配置读取)
- [命令行工具 toolkit](#命令行工具-toolkit)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [常见问题](#常见问题)
- [贡献](#贡献)
- [许可证](#许可证)

---

## 特性

### 多协议服务端

| 服务端 | 底层框架 | 场景 |
|---|---|---|
| `server/xgin` | Gin | RESTful API、Web 应用、WebSocket |
| `server/xgrpc` | gRPC | 微服务间 RPC 通信 |
| `server/xhertz` | CloudWeGo Hertz | 高并发 HTTP 服务 |
| `server/xasynq` | Asynq | 定时任务调度与后台任务消费 |
| `server/xmonitor` | net/http | pprof、路由清单、构建信息等治理接口 |

所有服务端统一实现 `server.Server` 接口（`Serve`/`Stop`/`GracefulStop`/`Info`/`Health`），由 `core/app` 统一管理生命周期，支持信号驱动的优雅启停。

### 统一客户端与存储

- **etcd**：KV 读写、前缀/多键批量查询、分布式互斥锁、租约会话、前缀监听、BasicAuth/TLS、链路追踪。
- **gRPC 客户端**：etcd 服务发现、SWR 平滑加权轮询与 P2C 最少负载均衡、调试/超时/追踪/日志/指标拦截器链。
- **RabbitMQ**：断线自动重连、通道复用、队列声明绑定 Builder、可配置重试策略、断线后消费自动恢复。
- **Redis**：单节点、集群、哨兵、分片四种模式；命令按数据类型分组；统一经熔断器执行；可挂载追踪与指标钩子。
- **SQL**：MySQL、PostgreSQL、ClickHouse；基于 GORM；DSN 级连接缓存复用；分页/排序/范围参数辅助；指标与追踪拦截器。
- **MongoDB**：链式条件构造、熔断/指标/追踪增强集合、写入前结构体到 BSON 自动转换。
- **Elasticsearch**：节点缓存复用、HTTP 传输层追踪与指标。
- **可空类型**：替代 `sql.NullXxx` 的完整可空类型集，统一实现 Scanner/Valuer/JSON/Text 接口。
- **代理池**：按分组统一管理 Redis/Mongo/ES/SQL 实例，SQL 支持读写库分离与随机选库。
- **会话存储**：基于 Redis 的 JSON 会话，点分路径读写嵌套字段。

### 服务治理

- **限流**：令牌桶，支持固定间隔、固定速率、每次补充量子三种填充方式；阻塞/限时/非阻塞三种取令牌模式；限流 Reader/Writer。
- **熔断**：基于 Google SRE 过载保护模型的自适应熔断器，支持自定义错误判定与降级回调。
- **防击穿**：SingleFlight 并发单飞。
- **批量错误**：BatchError 批量错误聚合。
- **goroutine 编排**：串行/并行/限流并发/延迟启动/panic 恢复/错误收集。
- **生命周期**：并发生命周期控制器（cycle）、分阶段钩子（hooks，defer 语义逆序执行）。

### 可观测性

- **指标**：Prometheus 命名空间 `ratel`，预置服务端/客户端/任务/缓存等场景的 Counter、Gauge、Histogram、Summary 向量。
- **追踪**：OpenTelemetry，支持 HTTP/gRPC 两种 OTLP 导出协议，兼容 Jaeger 与 W3C TraceContext 传播。
- **日志**：基于 logrus 的结构化日志；控制台彩色输出；Elasticsearch、MongoDB、SQL 三套异步钩子，后台队列写入。

### 任务与并发

- **cron 调度**：五/六字段表达式与 `@yearly` 等预定义描述符，超时控制、动态增删、优雅关闭。
- **分布式任务**：cron 调度 + 一致性哈希节点池，集群内分配任务，节点变更自动进入升级状态避免重复执行；支持任务恢复与近期任务回放。
- **分布式执行器**：xxl-job 风格的任务执行器接口与全局管理。
- **Asynq**：定时任务管理器与后台任务工作者。

### 配套工具

- 雪花 ID 生成器（原子操作、并发安全、可自定义位宽）。
- AC 自动机敏感词过滤（`nlpword` 与 `filter` 两套实现）。
- 中英文混排文本自动排版（`correct`）。
- i18n 多语言消息管理（`lang/i18n`）。
- 图像读取/变换/合成/水印/文字绘制（`img`，基于 libvips）。
- Go 项目热编译监视（`watch`）。
- 代码生成脚手架（`toolkit`）：项目创建、DAO/接口/模型/Vue 代码生成、更新日志、脚手架升级。

---

## 架构设计

### 分层结构

```mermaid
flowchart TB
    %% ============ 层级：自上而下 ============

    subgraph APP["应用层 core/app —— 组件组装与生命周期"]
        App["Application<br/>New / Run / Serve / Schedule / Job / Executor<br/>RegisterHooks / SetRegistry / Stop / GracefulStop"]
    end

    subgraph SRV["服务端层 server/ —— 统一 server.Server 生命周期接口"]
        Xgin["xgin<br/>Gin HTTP 服务"]
        Xgrpc["xgrpc<br/>gRPC 服务"]
        Xhertz["xhertz<br/>Hertz HTTP 服务"]
        Xasynq["xasynq<br/>Asynq 异步任务服务"]
        Xmonitor["xmonitor<br/>监控 / pprof 服务"]
    end

    subgraph CLI["客户端层 client/ —— 中间件客户端"]
        CEtcd["etcdv3<br/>配置中心 / 注册中心客户端"]
        CGrpc["grpc<br/>服务发现 + 拦截器链"]
        CRabbit["rabbitmq<br/>发布订阅客户端"]
        CFtp["ftp / sftp<br/>文件传输客户端"]
    end

    subgraph STO["存储层 stores/ —— 数据访问"]
        SRedis["redis<br/>单机 / 集群 / 哨兵 / Ring 四模式"]
        SSql["sql<br/>GORM 封装 · MySQL / PostgreSQL / ClickHouse"]
        SMongo["mongodb<br/>mongo-driver 封装"]
        SEs["elasticsearch"]
        SNull["null<br/>数据库可空类型"]
        SProxy["proxy<br/>多客户端代理池"]
        SSession["session<br/>Redis 会话存储"]
    end

    subgraph REG["注册中心 registry/"]
        RInterface["registry.Registry 接口<br/>Register / Deregister / GetService"]
        REtcd["etcdv3 实现<br/>租约续期 · 服务发现"]
    end

    subgraph CORE["治理核心 core/ —— 横切能力"]
        CRatelimit["ratelimit 令牌桶限流"]
        CResource["resource 熔断 / SingleFlight / 资源管理器"]
        CTrace["trace OpenTelemetry 链路追踪"]
        CMetric["metric Prometheus 指标"]
        CLogger["logger 多后端异步日志"]
        CGoroutine["goroutine 并发编排"]
        CTask["task 分布式定时任务"]
        CCycle["cycle 生命周期控制"]
        CHooks["hooks 阶段钩子"]
        CExecutor["executor 分布式执行器"]
    end

    subgraph TOOL["工具层 util/ 等 —— 通用能力"]
        TUtil["util 通用工具"]
        TCorrect["correct 文本排版"]
        TNlp["nlpword 敏感词检测"]
        TFilter["filter 文本过滤"]
        TLang["lang / i18n 国际化"]
        TImg["img 图像处理"]
        TSnow["snowflake 雪花 ID"]
        TCron["crontab 调度器"]
        TWatch["watch 文件监控"]
        TConfig["config 配置加载"]
    end

    subgraph INFRA["外部基础设施"]
        IEtcd["etcd 集群"]
        IRabbit["RabbitMQ"]
        IRedis["Redis"]
        ISql["MySQL / PostgreSQL / ClickHouse"]
        IMongo["MongoDB"]
        IEs["Elasticsearch"]
        IFtp["FTP / SFTP"]
        IOtel["OpenTelemetry Collector"]
        IProm["Prometheus"]
    end

    %% ============ 生命周期托管 ============
    APP -->|"Run 统一启动"| SRV
    APP -->|"Schedule 托管"| Xasynq
    APP -->|"Job 托管"| CTask
    APP -->|"Executor 托管"| CExecutor
    APP -->|"阶段钩子"| CHooks
    APP -->|"SetRegistry 注入"| REG

    %% ============ 服务端调用链 ============
    Xgin -->|"业务调用"| CGrpc
    Xgin -->|"业务调用"| CRabbit
    Xgrpc -->|"业务调用"| CGrpc
    Xhertz -->|"业务调用"| CGrpc
    Xasynq -->|"消费任务"| SRedis

    %% ============ 数据访问 ============
    SRV -->|"数据访问"| STO
    SSession -->|"基于"| SRedis
    SProxy -->|"统一访问"| SRedis
    SProxy -->|"统一访问"| SMongo
    SProxy -->|"统一访问"| SEs
    SProxy -->|"统一访问"| SSql

    %% ============ 服务发现 / 配置 ============
    CGrpc -->|"服务发现"| REG
    SRV -->|"启动登记"| REG
    REtcd -->|"读写节点"| IEtcd
    CEtcd -->|"读写"| IEtcd

    %% ============ 中间件连接 ============
    CRabbit -->|"连接"| IRabbit
    SRedis -->|"连接"| IRedis
    SSql -->|"连接"| ISql
    SMongo -->|"连接"| IMongo
    SEs -->|"连接"| IEs
    CFtp -->|"连接"| IFtp

    %% ============ 治理横切（虚线） ============
    STO -.->|"命令经熔断器 / 指标 / 追踪"| CORE
    CLI -.->|"拦截器链（追踪 / 指标 / 访问日志 / 超时）"| CORE
    CORE -->|"trace 上报"| IOtel
    CORE -->|"metric 暴露"| IProm

    %% ============ 任务调度内部依赖 ============
    CTask -->|"生成 ID"| TSnow
    CTask -->|"节点注册"| CEtcd
    TCron -->|"底层调度"| CTask
    TConfig -->|"配置来源"| CEtcd

    %% ============ 样式 ============
    classDef app fill:#1f77b4,stroke:#0d4f8b,color:#fff,stroke-width:2px;
    classDef server fill:#2ca02c,stroke:#1a7a1a,color:#fff;
    classDef client fill:#ff7f0e,stroke:#b35a00,color:#fff;
    classDef store fill:#9467bd,stroke:#6b3fa0,color:#fff;
    classDef registry fill:#17becf,stroke:#0e7f8f,color:#fff;
    classDef core fill:#d62728,stroke:#9c1c1e,color:#fff;
    classDef tool fill:#7f7f7f,stroke:#5a5a5a,color:#fff;
    classDef infra fill:#8c564b,stroke:#5f3a33,color:#fff;

    class App app;
    class Xgin,Xgrpc,Xhertz,Xasynq,Xmonitor server;
    class CEtcd,CGrpc,CRabbit,CFtp client;
    class SRedis,SSql,SMongo,SEs,SNull,SProxy,SSession store;
    class RInterface,REtcd registry;
    class CRatelimit,CResource,CTrace,CMetric,CLogger,CGoroutine,CTask,CCycle,CHooks,CExecutor core;
    class TUtil,TCorrect,TNlp,TFilter,TLang,TImg,TSnow,TCron,TWatch,TConfig tool;
    class IEtcd,IRabbit,IRedis,ISql,IMongo,IEs,IFtp,IOtel,IProm infra;
```

**分层说明：**

| 层 | 目录 | 职责 |
|---|---|---|
| 应用层 | `core/app` | 组装 Server/Worker/Job/Executor，信号驱动启停，阶段钩子，注册中心注入 |
| 服务端层 | `server/` | 统一封装 Gin / gRPC / Hertz / Asynq / 监控，实现 `server.Server` 接口 |
| 客户端层 | `client/` | etcd、gRPC、RabbitMQ、FTP/SFTP 客户端，含拦截器链与连接管理 |
| 存储层 | `stores/` | Redis（四模式）、SQL（GORM）、MongoDB、Elasticsearch、可空类型、代理池、会话 |
| 注册中心 | `registry/` | `Registry` 接口与 etcdv3 实现（租约、服务发现） |
| 治理核心 | `core/` | 限流、熔断、追踪、指标、日志、并发编排、分布式任务、生命周期、钩子、执行器 |
| 工具层 | `util/` 等 | 通用工具、文本处理、图像、ID 生成、配置加载、文件监控、i18n |
| 外部基础设施 | — | etcd / RabbitMQ / Redis / SQL / MongoDB / ES / FTP·SFTP / OTel / Prometheus |

**关键链路：**

- **业务请求**：应用层 `Run` 启动服务端 → 服务端经客户端层调用下游、经存储层读写数据 → 全部命令穿过治理核心（熔断 → 指标 → 追踪）。
- **服务发现**：服务端启动向注册中心登记；gRPC 客户端从注册中心解析目标节点；etcdv3 实现负责租约续期与节点读写。
- **可观测性**：trace 上报 OpenTelemetry Collector，metric 暴露给 Prometheus，日志经 logger 多后端异步写出。
- **任务调度**：`task` 依赖 `snowflake` 生成任务 ID、依赖 etcd 客户端做节点注册，`crontab` 提供底层 cron 调度。

### 应用生命周期

`core/app` 负责组装全部组件：

1. 应用启动时依次启动注册的 Server、Worker、Job、Executor。
2. Server 启动后自动向注册中心登记（`registry.DefaultRegisterer`）。
3. 监听系统信号（SIGINT/SIGTERM 等），收到信号后执行清理：停止 Server 与 Worker、关闭后台任务、触发阶段钩子。
4. 支持 `Stop`（立即停止）与 `GracefulStop`（优雅停止）两种方式。

### 设计约定

- **配置即代码**：各组件均以 `New()` + 链式 `WithXxx` 选项构造，配置项集中、可读、可静态检查。
- **统一生命周期**：Server/Worker/Job 均实现框架约定接口，由应用统一管理，业务代码无需关心启停细节。
- **治理内建**：存储与客户端命令默认经过熔断器，可开关式挂载指标与追踪，无需业务侧侵入。
- **包注释即文档**：每个库包都有 `doc.go` 包注释，导出 API 带 Godoc，可直接用 `go doc` / pkg.go.dev 查阅。

---

## 安装

### 环境要求

| 依赖 | 版本要求 | 说明 |
|---|---|---|
| Go | 1.27+ | 见 `go.mod` 的 `go` 指令 |
| libvips | 可选 | 仅 `img` 图像处理包需要；macOS 可用 `brew install vips` |
| pkg-config | 可选 | 编译 `img` 包时 libvips 需要 |

### 获取模块

```bash
go get github.com/abulo/ratel/v3@latest
```

---

## 快速开始

### 1. Gin HTTP 服务

```go
package main

import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/server/xgin"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 构建 Gin 服务端：监听 8080，启用指标与慢查询中间件
	httpServer := xgin.New().
		WithPort(8080).
		WithEnableMetric(true).
		WithEnableSlowQuery(true).
		Build()

	// 2. 注册路由：Server 内嵌 gin.Engine，Gin 原生 API 全部可用
	httpServer.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Ratel!"})
	})

	// 3. 组装应用并启动：阻塞等待退出信号，支持优雅停机
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	if err := application.Run(httpServer); err != nil {
		panic(err)
	}
}
```

### 2. gRPC 服务

```go
package main

import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/server/xgrpc"
)

func main() {
	// 构建 gRPC 服务端（Server 内嵌 grpc.Server，原生 gRPC API 可用）
	cfg := xgrpc.New()
	cfg.Host = "0.0.0.0"
	cfg.Port = 50051
	cfg.EnableTrace = true
	cfg.EnableMetric = true
	grpcServer, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	// 注册你的服务实现，例如：
	// pb.RegisterYourServiceServer(grpcServer, &yourServiceImpl{})

	application, err := app.New()
	if err != nil {
		panic(err)
	}
	if err := application.Run(grpcServer); err != nil {
		panic(err)
	}
}
```

### 3. Hertz 服务

```go
package main

import (
	"context"

	ratelapp "github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/server/xhertz"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	hertzServer := xhertz.New().
		WithPort(8081).
		WithEnableMetric(true).
		WithEnableTrace(true).
		Build()

	// Server 内嵌 hertz 引擎，Hertz 原生 API 可用
	hertzServer.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, map[string]string{"message": "Hello from Hertz!"})
	})

	application, err := ratelapp.New()
	if err != nil {
		panic(err)
	}
	if err := application.Run(hertzServer); err != nil {
		panic(err)
	}
}
```

> 注意：`core/app` 与 `hertz/pkg/app` 包名相同，同时导入时需给其中一个加别名（上例使用 `ratelapp`）。

### 4. 组合多个服务端

一个应用可以同时托管多个服务端：

```go
func main() {
	httpServer := xgin.New().WithPort(8080).Build()
	grpcServer := mustBuildGRPC()
	monitorServer := xmonitor.New().WithPort(9090).Build() // 独立监控端，暴露 pprof 等

	application, err := app.New()
	if err != nil {
		panic(err)
	}
	// Run 可接收任意多个 server.Server
	if err := application.Run(httpServer, grpcServer, monitorServer); err != nil {
		panic(err)
	}
}
```

---

## 应用生命周期

`core/app` 提供 `Application` 应用生命周期管理器。

### Application 方法

| 方法 | 说明 |
|---|---|
| `New(fns ...func() error) (*Application, error)` | 创建应用，`fns` 为启动前执行的初始化函数 |
| `DefaultApp() *Application` | 创建默认应用（跳过参数解析与配置加载） |
| `Serve(s ...server.Server) error` | 注册服务端 |
| `Schedule(w worker.Worker) error` | 注册后台工作单元 |
| `Job(runner job.Runner) error` | 注册启动后并发执行的任务 |
| `Executor(e executor.Executor)` | 注册分布式任务执行器 |
| `SetRegistry(reg registry.Registry)` | 设置注册中心 |
| `Run(servers ...server.Server) error` | 启动全部组件并阻塞等待退出信号 |
| `Stop() error` | 立即停止应用 |
| `GracefulStop(ctx context.Context) error` | 优雅停止应用 |
| `RegisterHooks(stage hooks.Stage, fns ...func())` | 注册生命周期钩子 |
| `WithOptions(options ...Option)` | 应用选项 |

### 生命周期钩子

```go
import "github.com/abulo/ratel/v3/core/hooks"

func main() {
	// 注册启动钩子（在配置加载前执行）
	hooks.Register(hooks.Stage_BeforeLoadConfig, func() {
		// 自定义初始化逻辑
	})

	application, err := app.New()
	// ...
}
```

钩子按阶段分组，同一阶段按注册顺序**逆序**执行（模拟 defer 语义）。

### 信号处理

`Run` 会监听系统信号：

- **SIGINT / SIGTERM**：触发优雅停止流程。
- 停止流程：执行 `BeforeStop` 钩子 → 停止全部 Server/Worker → 等待后台任务结束 → 执行 `AfterStop` 钩子 → 关闭生命周期。

---

## 服务端详解

所有服务端统一实现 `server.Server` 接口，并支持：

| 方法 | 说明 |
|---|---|
| `Serve() error` | 启动服务并阻塞监听 |
| `Stop() error` | 立即停止 |
| `GracefulStop(ctx context.Context) error` | 优雅停止（等待在途请求完成） |
| `Info() *server.ServiceInfo` | 返回注册信息（名称、地址、协议、标签） |
| `Health() bool` | 健康检查 |

### xgin

基于 Gin 的 HTTP 服务端。

```go
server := xgin.New().
	WithHost("0.0.0.0").          // 监听主机
	WithPort(8080).               // 监听端口
	WithMode(gin.ReleaseMode).    // 运行模式（默认 ReleaseMode）
	WithDeployment("prod").       // 部署组（服务注册流量隔离）
	WithServiceAddress("").       // 注册中心登记地址，空则用 Host:Port
	WithEnableMetric(true).       // 启用指标采集中间件
	WithEnableTrace(true).        // 启用链路追踪中间件
	WithEnableSlowQuery(true).    // 启用慢查询与 panic 恢复中间件
	WithSlowQueryThresholdInMilli(500). // 慢查询阈值（毫秒）
	Build()
```

**配置项：**

| 字段 | 默认值 | 说明 |
|---|---|---|
| `Host` | 空 | 监听主机地址 |
| `Port` | 0 | 监听端口，构建时由实际绑定端口回填 |
| `Mode` | `ReleaseMode` | Gin 运行模式（`ReleaseMode`/`DebugMode`/`TestMode`） |
| `EnableMetric` | false | 指标采集中间件 |
| `EnableTrace` | false | 服务端链路追踪中间件 |
| `EnableSlowQuery` | false | 慢请求记录与 panic 恢复中间件 |
| `SlowQueryThresholdInMill` | 500 | 慢查询阈值（毫秒） |
| `Deployment` | 空 | 部署组 |
| `ServiceAddress` | 空 | 注册中心登记地址 |

**WebSocket：**

```go
type MyWebSocket struct {
	// 实现自己的协议升级与处理逻辑
}

// 注册 WebSocket 处理器为 GET 路由
server.Upgrade(&MyWebSocket{Pattern: "/ws"})
```

### xgrpc

基于 gRPC 的服务端。

```go
cfg := xgrpc.New()
cfg.Host = "0.0.0.0"
cfg.Port = 50051
cfg.Network = "tcp4"
cfg.EnableTrace = true
cfg.EnableMetric = true
cfg.EnableTLS = false // 开启后需配置 CaFile/CertFile/PrivateFile

// 追加拦截器
cfg.WithUnaryInterceptor(myUnaryInterceptor)
cfg.WithStreamInterceptor(myStreamInterceptor)

grpcServer, err := cfg.Build()     // 返回 (*Server, error)
// 或 grpcServer := cfg.MustBuild() // 构建失败直接 panic
```

**配置项（与 xgrpc 风格一致，字段赋值 + `Build`）：**

| 字段 | 默认值 | 说明 |
|---|---|---|
| `Host` / `Port` | 空 / 0 | 监听地址，端口构建时回填 |
| `Network` | `tcp4` | 监听网络类型 |
| `Deployment` | `constant.DefaultDeployment` | 部署组 |
| `EnableTrace` | false | 服务端链路追踪拦截器 |
| `EnableMetric` | false | 指标采集拦截器 |
| `EnableTLS` | false | 双向 TLS 认证 |
| `CaFile` / `CertFile` / `PrivateFile` | 空 | TLS 证书文件路径 |
| `SlowQueryThresholdInMill` | 500 | 预留配置项 |

**panic 恢复：** 子包 `server/xgrpc/recovery` 提供一元与流式恢复拦截器，未配置回调时把 panic 包装为 `PanicError` 返回，避免单请求 panic 导致进程退出。

### xhertz

基于 CloudWeGo Hertz 的高性能 HTTP 服务端，配置项与 xgin 对齐（`WithHost`/`WithPort`/`WithMode`/`WithEnableMetric`/`WithEnableTrace`/`WithEnableSlowQuery`/`WithDeployment`/`WithServiceAddress`/`WithSlowQueryThresholdInMilli`），`Build()` 返回内嵌 hertz 引擎的 `*Server`。

```go
server := xhertz.New().WithPort(8081).Build()
server.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, map[string]string{"message": "hello"})
})
```

### xasynq

基于 Asynq 的定时任务与后台任务组件，分 `manager` 与 `worker` 两个子包。

**manager（定时任务管理器）：** 按已注册的定时计划周期性将任务写入 Redis 队列。配置项与 xgin/xhertz 对齐（`WithHost`/`WithPort`/`WithDeployment`/`WithEnableMetric`/`WithEnableTrace`/`WithMode` 等），Asynq 运行选项通过 `Config.Option` 字段传入。

```go
import (
	"github.com/abulo/ratel/v3/server/xasynq/manager"
	"github.com/hibiken/asynq"
)

mgrCfg := manager.New().
	WithHost("0.0.0.0").
	WithPort(0)

// Build 前配置 Asynq 运行选项（含 Redis 连接与定时任务注册表）
mgrCfg.Option = asynq.PeriodicTaskManagerOpts{
	RedisConnOpt: asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
	// 在此注册定时任务：cronspec -> 任务处理器
}

mgr := mgrCfg.Build() // 实现 server.Server 接口
// application.Run(mgr)
```

**worker（任务工作者）：** 从 Redis 队列拉取任务并按处理器路由表消费。`Config.Redis`（`redis.UniversalClient`）、`Config.Mux`（`*asynq.ServeMux`）与 `Config.Option`（`asynq.Config`）在 Build 前赋值。

```go
import (
	"github.com/abulo/ratel/v3/server/xasynq/worker"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

wkCfg := worker.New().
	WithHost("0.0.0.0").
	WithPort(0)

// Build 前配置 Redis、处理器路由表与运行选项
wkCfg.Redis = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
wkCfg.Mux = asynq.NewServeMux()
wkCfg.Mux.HandleFunc("task:email", func(ctx context.Context, t *asynq.Task) error {
	// 处理邮件任务
	return nil
})
wkCfg.Option = asynq.Config{Concurrency: 10}

wk := wkCfg.Build() // 实现 server.Server 接口
// application.Run(wk)
```

### xmonitor

独立的监控与治理 HTTP 服务端，在独立路由表上挂载 pprof 性能分析、已注册路由清单与模块构建信息等内置治理接口。

```go
monitorServer := xmonitor.New().
	WithPort(9090).
	Build()

// application.Run(httpServer, monitorServer)
```

---

## 客户端详解

### etcdv3

基于 `go.etcd.io/etcd/client/v3` 的 etcd 客户端封装。

```go
import (
	"github.com/abulo/ratel/v3/client/etcdv3"
)

// Config 使用字段赋值或 Set 前缀方法填充（非 With 链式）
cfg := etcdv3.New() // 默认 BasicAuth=false、ConnectTimeout=5s、Secure=false
cfg.SetEndpoints([]string{"127.0.0.1:2379"})
cfg.SetBasicAuth(true)
cfg.SetUserName("root")
cfg.SetPassword("secret")
cfg.SetConnectTimeout(5 * time.Second)

client, err := etcdv3.NewClient(cfg)
if err != nil {
	panic(err)
}
```

**Client 能力（内嵌原生 `clientv3.Client`，原生 API 直接可用）：**

| 方法 | 说明 |
|---|---|
| `GetKeyValue(ctx, key)` | 读取单个键 |
| `GetValues(ctx, keys...)` | 批量读取多个键 |
| `GetPrefix(ctx, prefix)` | 前缀读取，返回 `map[string]string` |
| `DelPrefix(ctx, prefix)` | 前缀删除 |
| `GetLeaseSession(ctx, opts...)` | 创建租约会话（keepalive） |
| `NewMutex(key, opts...)` | 创建分布式互斥锁 |
| `WatchPrefix(ctx, prefix)` | 前缀变更监听 |

**分布式锁：**

```go
mutex, err := client.NewMutex("/lock/order")
if err != nil {
	panic(err)
}
if err := mutex.Lock(ctx); err != nil {
	panic(err)
}
defer mutex.Unlock(ctx)
// 临界区
```

### grpc 客户端

基于 etcd 服务发现的 gRPC 客户端。

```go
import (
	grpcc "github.com/abulo/ratel/v3/client/grpc"
	regEtcd "github.com/abulo/ratel/v3/registry/etcdv3"
)

// 1. etcd 服务发现配置（注意：使用 registry/etcdv3 的 Config；
// Set 方法链式调用会返回内嵌的 client 配置类型，应分开调用保持类型）
etcdCfg := regEtcd.New()
etcdCfg.SetEndpoints([]string{"127.0.0.1:2379"})

// 2. gRPC 客户端配置
conn := grpcc.New().
	SetName("order-service").               // 客户端实例名（单例键与标签）
	SetAddress("order.svc").                // 目标服务名（由 etcd 解析器解析）
	SetBalancerName("p2c").                 // 负载均衡：roundrobin / swr / p2c
	SetEtcd(etcdCfg).                       // 必填：etcd 服务发现
	SetReadTimeout(1 * time.Second).        // 单次调用超时
	SetSlowThreshold(600 * time.Millisecond).
	SetEnableTraceInterceptor(true).
	SetEnableMetricInterceptor(true).
	SetEnableAccessInterceptor(true).
	SetEnableTimeoutInterceptor(true).
	MustBuild() // 或 Build() 返回 (conn, error)
```

**配置项：**

| 字段 | 默认值 | 说明 |
|---|---|---|
| `BalancerName` | `roundrobin` | 负载均衡策略（`roundrobin`/`swr`/`p2c`） |
| `DialTimeout` | 3s | 拨号超时 |
| `ReadTimeout` | 1s | 单次调用超时（超时拦截器用） |
| `SlowThreshold` | 600ms | 慢调用判定阈值 |
| `EnableDebug` | false | 调试拦截器（打印地址/方法/响应/错误） |
| `EnableTraceInterceptor` | false | 链路追踪拦截器 |
| `EnableAidInterceptor` | false | 应用 ID 元数据透传 |
| `EnableTimeoutInterceptor` | false | 调用超时与慢调用统计 |
| `EnableMetricInterceptor` | false | 指标采集 |
| `EnableAccessInterceptor` | false | 访问日志拦截器 |
| `Etcd` | 必填 | etcd 服务发现配置，缺失会空指针 |

**负载均衡：** `client/grpc/balancer` 提供 SWR（平滑加权轮询）与 P2C（随机二抽样最少负载）均衡器，`p2c` 子包按在途请求数选择负载更低的节点。

### rabbitmq

基于 `streadway/amqp` 的 RabbitMQ 客户端封装，支持断线自动重连。

```go
import (
	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

// 1. 建立连接（可断线重连）
conn, err := rabbitmq.Dial("amqp://guest:guest@127.0.0.1:5672/", rabbitmq.DefaultTimesRetry())
if err != nil {
	panic(err)
}

// 2. 声明并绑定队列（队列名、路由键、交换器）
queue := rabbitmq.NewQueueBuilder(conn).
	SetQueueDeclareOpts(func(builder *rabbitmq.QueueDeclareOptsBuilder) *rabbitmq.QueueDeclareOpts {
		return builder.SetDurable(true).SetAutoDelete(false).Build()
	}).
	Build()
if err := queue.DeclareAndBind("task.queue", "task.key", "exchange"); err != nil {
	panic(err)
}

// 3. 获取通道并发送消息（通道由连接管理，复用常驻通道）
ch, err := conn.Channel()
if err != nil {
	panic(err)
}
if err := ch.SendOpts("exchange", "task.key", []byte("hello"), rabbitmq.DefaultSendOpts()); err != nil {
	// 发送失败处理
}

// 4. 消费消息（ConsumerFunc 返回 true 终止消费）
// ch.ReceiveOpts("task.queue", func(delivery *amqp.Delivery) bool {
// 	处理消息，返回 true 表示确认并退出
// 	return true
// }, rabbitmq.DefaultReceiveOpts())
```

**关键类型：**

| 类型 | 说明 |
|---|---|
| `Connection` | 可断线重连的连接，`Dial`/`NewConnection` 构造 |
| `Channel` | 复用连接的通道，`Connection.Channel()` 获取；提供 `Send`/`SendOpts`/`Receive`/`ReceiveOpts` |
| `QueueBuilder` | 队列声明与绑定 Builder（`SetQueueDeclareOpts`/`SetQueueBindOpts`/`SetRetryable`/`Build`） |
| `Queue` | 队列载体，`DeclareAndBind`/`RetryDeclareAndBind` 声明绑定 |
| `Producer` / `Consumer` | 便捷生产者（`Send`）与消费者（`Receive`/`Get`）封装 |
| `SendOptsBuilder` / `ReceiveOptsBuilder` / `QueueDeclareOptsBuilder` / `QueueBindOptsBuilder` | 选项 Builder |
| `TimesRetry` / `CtxRetry` | 重试策略 |
| `MessagePlainTransient` 等 | 常用消息工厂 |

### ftp / sftp

文件传输客户端：

```go
// FTP：配置字段 + Config.New 拨号 + Login 认证
ftpConfig := &ftp.Config{
	Host:     "127.0.0.1",
	Port:     "21",
	User:     "user",
	Password: "pass",
	Timeout:  10 * time.Second,
}
ftpClient, err := ftpConfig.New()
if err != nil {
	panic(err)
}
defer ftpClient.Quit()
if err := ftpClient.Login(); err != nil {
	panic(err)
}

// SFTP：配置字段 + Config.New（底层为 SSH 连接 + SFTP 通道）
sftpConfig := &sftp.Config{
	Host:     "127.0.0.1",
	Port:     "22",
	User:     "user",
	Password: "pass",
	Timeout:  10 * time.Second,
}
sftpClient, err := sftpConfig.New()
if err != nil {
	panic(err)
}
defer sftpClient.SFTPClient.Close()
```

两个客户端均提供目录切换、文件/目录增删查、上传下载（SFTP 支持递归目录与传输进度统计）等操作。

---

## 存储层详解

### redis

基于 go-redis 的 Redis 客户端封装，支持四种连接模式。

```go
import "github.com/abulo/ratel/v3/stores/redis"

// 单节点
client, err := redis.NewRedisClient(
	redis.WithClientType("normal"),
	redis.WithAddr("127.0.0.1:6379"),
	redis.WithPassword(""),
	redis.WithDatabase(0),
	redis.WithPoolSize(20),
)

// 集群
client, err = redis.NewRedisClient(
	redis.WithClientType("cluster"),
	redis.WithHosts([]string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}),
)

// 哨兵
client, err = redis.NewRedisClient(
	redis.WithClientType("failover"),
	redis.WithMasterName("mymaster"),
	redis.WithHosts([]string{"127.0.0.1:26379", "127.0.0.1:26380"}),
)

// 分片
client, err = redis.NewRedisClient(
	redis.WithClientType("ring"),
	redis.WithAddrs(map[string]string{
		"shard1": "127.0.0.1:6379",
		"shard2": "127.0.0.1:6380",
	}),
)
```

**配置项（Option）：**

| 选项 | 说明 |
|---|---|
| `WithClientType` | `normal`（单节点）/`cluster`/`failover`/`ring` |
| `WithAddr` | 单节点地址 `host:port` |
| `WithHosts` | 集群/哨兵节点列表 |
| `WithAddrs` | 分片名称到地址映射 |
| `WithMasterName` | 哨兵主节点名称 |
| `WithDatabase` | 数据库编号 |
| `WithPassword` | 连接密码 |
| `WithDialTimeout` / `WithReadTimeout` / `WithWriteTimeout` | 超时配置 |
| `WithPoolSize` / `WithPoolTimeout` / `WithIdleTimeout` / `WithIdleCheckFrequency` / `WithConnMaxLifetime` | 连接池配置 |
| `WithEnableMetric` / `WithEnableTrace` | 指标与追踪钩子 |
| `WithName` | 实例名称（单例键） |

**命令调用：** 命令按数据类型组织为分组方法（字符串、哈希、列表、集合、有序集合、位图、流、地理、JSON、搜索、脚本、Pub/Sub 等），统一经熔断器执行；键不存在等可接受错误不会作为失败上报。底层 go-redis 客户端可通过 `redis.OriginalClient(client)` 获取。

### sql

基于 GORM 的关系型数据库客户端封装，支持 MySQL、PostgreSQL、ClickHouse。

```go
import "github.com/abulo/ratel/v3/stores/sql"

// 1. 配置并打开连接：NewClient 只装配配置，Open 才建立连接（返回内嵌 *gorm.DB 的 ClientManager）
client, err := sql.NewClient(
	sql.WithDriverName(sql.MySQLDriverName), // 或 sql.PostgreSQLDriverName / sql.ClickHouseDriverName
	sql.WithHost("127.0.0.1"),
	sql.WithPort("3306"),
	sql.WithUsername("root"),
	sql.WithPassword("secret"),
	sql.WithDatabase("mydb"),
	sql.WithCharset("utf8mb4"),
	sql.WithParseTime(true),
	sql.WithMaxOpenConns(100),
	sql.WithMaxIdleConns(20),
	sql.WithMaxLifetime(1*time.Hour),
	sql.WithEnableMetric(true),
	sql.WithEnableTrace(true),
)
if err != nil {
	panic(err)
}
manager, err := client.Open() // (*ClientManager, error)
if err != nil {
	panic(err)
}

// 2. 执行 GORM 查询（manager 内嵌 *gorm.DB，GORM 原生 API 可用）
// manager.Model(&User{}).Find(&users)
```

**配置项（Option）：**

| 选项 | 说明 |
|---|---|
| `WithDriverName` | `mysql` / `postgres` / `clickhouse` |
| `WithHost` / `WithPort` / `WithUsername` / `WithPassword` / `WithDatabase` | 连接信息 |
| `WithCharset` / `WithParseTime` / `WithTimeZone` / `WithSslMode` | 驱动特定 DSN 参数 |
| `WithDialTimeOut` / `WithReadTimeOut` | ClickHouse 超时 |
| `WithMaxOpenConns` / `WithMaxIdleConns` / `WithMaxIdleTime` / `WithMaxLifetime` | 连接池 |
| `WithEnableDebug` / `WithEnableMetric` / `WithEnableTrace` | 开关 |

**查询辅助：**

```go
// 分页参数载体（字段为指针，未设置时取值返回错误）
offset := int64(0)
limit := int64(20)
page := &sql.Pagination{Offset: &offset, Limit: &limit}

// 范围与排序参数载体
scope := &sql.Scope{Start: ptr("2024-01-01"), End: ptr("2024-12-31")}
sort := &sql.Sort{Name: ptr("created_at"), Desc: ptr(true)}

// 取值后拼入 GORM 链式查询（manager 内嵌 *gorm.DB）
offsetVal, _ := page.GetOffset()
limitVal, _ := page.GetLimit()
sortName, _ := sort.GetName()
manager.WithContext(ctx).Model(&User{}).
	Offset(int(offsetVal)).
	Limit(int(limitVal)).
	Order(sortName).
	Find(&users)
```

> `ptr` 为取地址辅助（`func ptr[T any](v T) *T { return &v }`），可按需自行定义。

**拦截器：** `sql.MetricInterceptor()` 与 `sql.TraceInterceptor()` 可通过 `sql.RegisterInterceptor(manager, client, interceptors...)` 注册，自动采集 SQL 指标与链路追踪（`Client.Open` 已在 EnableMetric/EnableTrace 开启时自动装配）。

### mongodb

基于 mongo-driver 的 MongoDB 客户端封装。

```go
import "github.com/abulo/ratel/v3/stores/mongodb"

mongoClient, err := mongodb.NewMongoDBClient("mongodb://127.0.0.1:27017",
	mongodb.WithTimeout(5*time.Second),
	mongodb.WithMaxPoolSize(100),
	mongodb.WithMinPoolSize(10),
	mongodb.WithMaxConnIdleTime(1*time.Minute),
)
if err != nil {
	panic(err)
}

// 获取集合增强句柄（collection 为集合名，库名取自客户端 Name）
model, err := mongoClient.NewModel("users")
if err != nil {
	panic(err)
}
// model.Collection 为 *mongodb.DecoratedCollection，包装了熔断、指标与追踪
```

**DecoratedCollection 能力：**

- 链式条件构造（封装原生 Filter 构建）。
- 故障保护：操作经熔断器执行。
- 指标采集与链路追踪自动挂载。
- `BeforeCreate`/`BeforeUpdate` 在写入前完成结构体到 BSON 文档的转换与 `_id` 生成。

### elasticsearch

基于 olivere/elastic v7 的 Elasticsearch 客户端封装。

```go
import "github.com/abulo/ratel/v3/stores/elasticsearch"

esClient := elasticsearch.NewClient(&elasticsearch.Config{
	URL:          []string{"http://127.0.0.1:9200"},
	Username:     "elastic",
	Password:     "secret",
	EnableMetric: true,
	EnableTrace:  true,
})

// esClient 内嵌原生 elastic.Client，全部原生 API 可用
// esClient.Search("index").Query(...).Do(ctx)
```

客户端按节点地址缓存复用；`ESTracedTransport` 在 HTTP 传输层挂载链路追踪与指标采集，`MaxContentLength`（64KB）限制参与追踪记录的请求体大小。

### null 可空类型

提供完整可空基础类型集，用于替代 `sql.NullXxx` 处理数据库与 JSON 中的空值。

**已实现类型：** `Bool`、`Byte`、`Bytes`、`String`、`Int`、`Int8/16/32/64`、`Uint`、`Uint8/16/32/64`、`Float32/64`、`Time`、`Date`、`DateTime`、`Timestamp`、`CTime`、`JSON`。

每个类型由三个状态位组成：

- **取值**：实际值。
- **`Valid`**：是否持有非空有效值。
- **`Set`**：是否被显式赋值（包含赋值为 null）。

统一实现 `sql.Scanner`、`driver.Valuer`、`json.Marshaler`、`json.Unmarshaler`、`encoding.TextMarshaler`、`encoding.TextUnmarshaler` 接口；空值在 JSON 中序列化为 `null`，读写数据库映射为 SQL `NULL`。

```go
import "github.com/abulo/ratel/v3/stores/null"

type User struct {
	ID    int64     `json:"id"`
	Name  null.String `json:"name"`
	Email null.String `json:"email"`
}

// JSON 空值处理
b, _ := json.Marshal(User{ID: 1, Name: null.StringFrom("ratel"), Email: null.String{}})
// {"id":1,"name":"ratel","email":null}
```

### proxy 代理池

按分组注册与获取各类数据存储后端的代理池。

```go
import "github.com/abulo/ratel/v3/stores/proxy"

p := proxy.NewProxy()

// 初始化阶段：构造代理载体并注入底层客户端，再注册到代理池
redisProxy := proxy.NewRedis()
redisProxy.Store(redisClient)
p.StoreRedis("default", redisProxy)

mongoProxy := proxy.NewMongoDB()
mongoProxy.Store(mongoClient)
p.StoreMongoDB("default", mongoProxy)

esProxy := proxy.NewElasticSearch()
esProxy.Store(esClient)
p.StoreElasticSearch("default", esProxy)

// SQL 代理支持读写分离：SetWrite/SetRead 装配 *gorm.DB，Write/Read 随机选库
sqlProxy := proxy.NewSQL()
sqlProxy.SetWrite(sqlClient)       // 写库（必填）
sqlProxy.SetRead(sqlReadReplica)   // 读库（可选，缺省回退写库）
p.StoreSQL("write", sqlProxy)

// 运行期按分组取出使用
r := p.LoadRedis("default")   // *redis.Client
m := p.LoadMongoDB("default") // *mongodb.MongoDB
e := p.LoadElasticSearch("default") // *elasticsearch.Client
s := p.LoadSQL("write")       // *proxy.SQL
db := s.Write()               // 随机写库；读多库场景用 s.Read()
```

SQL 代理额外支持写库与读库列表及随机选库（读写分离）。

### session 会话存储

基于 Redis 的会话存储，整份会话以 JSON 形式保存在单个 Redis 键下。

```go
import "github.com/abulo/ratel/v3/stores/session"

sess := &session.Session{
	Driver: redisClient, // 底层 Redis 客户端，不可为 nil
	Name:   "session:user:1001",
	TTL:    3600, // 秒，0 表示不过期
}

// 点分路径读写嵌套字段
_ = sess.Put(ctx, "profile.name", "ratel")
_ = sess.Put(ctx, "profile.level", 3)
name := sess.Get(ctx, "profile.name")   // "ratel"

// 删除字段 / 销毁整份会话
sess.Remove(ctx, "profile.level")
sess.Destroy(ctx)
```

---

## 治理能力详解

### 限流 ratelimit

令牌桶限流实现。

```go
import "github.com/abulo/ratel/v3/core/ratelimit"

// 固定间隔填充：每 100ms 补充 1 个令牌，容量 100
bucket := ratelimit.NewBucket(100*time.Millisecond, 100)

// 固定速率填充：每秒补充 rate 个令牌
bucket = ratelimit.NewBucketWithRate(1000, 100)

// 每次补充量子：每 100ms 补充 5 个令牌
bucket = ratelimit.NewBucketWithQuantum(100*time.Millisecond, 100, 5)
```

**Bucket 方法：**

| 方法 | 行为 |
|---|---|
| `Take(count)` | 取令牌，返回等待时长；桶空时按需等待 |
| `TakeMaxDuration(count, maxWait)` | 限时等待取令牌，超时返回 `(0, false)` |
| `TakeAvailable(count)` | 非阻塞取令牌，返回实际取到的数量 |
| `Wait(count)` | 阻塞直到取到令牌（无返回值） |
| `WaitMaxDuration(count, maxWait)` | 限时等待，返回是否在时限内取到（bool） |
| `Available()` / `Capacity()` / `Rate()` | 查询当前可用令牌、容量、速率 |

**限流 Reader/Writer：**

```go
limitedReader := ratelimit.Reader(r, bucket)   // 限制读取速率
limitedWriter := ratelimit.Writer(w, bucket)   // 限制写入速率
```

桶内部加锁，可被多个 goroutine 并发使用。

### 熔断 resource

基于 Google SRE 过载保护模型的自适应熔断器。

```go
import "github.com/abulo/ratel/v3/core/resource"

// 1. 获取或创建命名熔断器
breaker := resource.GetBreaker("my-service")

// 2. 包装业务调用
err := breaker.Do(func() error {
	return callDownstream()
})

// 3. 自定义错误判定：只有指定错误才算失败
err = breaker.DoWithAcceptable(func() error {
	return callDownstream()
}, func(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
})

// 4. 降级回调：熔断时执行 fallback
err = breaker.DoWithFallback(func() error {
	return callDownstream()
}, func(err error) error {
	return fallbackResult()
})

// 5. 手动管理请求生命周期
promise, err := breaker.Allow()
if err != nil {
	return resource.ErrServiceUnavailable
}
defer promise.Reject("failed") // 或成功后 promise.Accept()；Accept/Reject 须且仅调用其一
// 业务逻辑
```

**Breaker 接口：**

| 方法 | 说明 |
|---|---|
| `Allow() (Promise, error)` | 放行时返回 Promise，成功调用 `Accept`，失败调用 `Reject` |
| `Do(req)` | 放行时执行 `req`，拒绝时返回错误；panic 记为失败后重抛 |
| `DoWithAcceptable(req, acceptable)` | 自定义错误是否视为成功 |
| `DoWithFallback(req, fallback)` | 拒绝时执行降级回调 |
| `DoWithFallbackAcceptable(...)` | 组合降级与自定义判定 |

**其它资源组件：**

```go
// SingleFlight 防击穿
sf := resource.NewSingleFlight()
v, err := sf.Do("key", func() (any, error) { return loadFromDB() })

// 批量错误聚合
batch := &resource.BatchError{}
batch.Add(err1)
batch.Add(err2)
if err := batch.Err(); err != nil { ... }

// 资源管理器（按名称管理可关闭资源）
manager := resource.NewResourceManager()
```

### goroutine 编排

`core/goroutine` 提供 goroutine 编排与安全执行工具。

| 函数 | 行为 |
|---|---|
| `Go(fn)` | 安全启动 goroutine（自动捕获 panic） |
| `DelayGo(delay, fn)` | 延迟启动 |
| `SafeGo(fn, rec)` | 启动并回调恢复函数处理 panic |
| `GoDirect(fn, args...)` | 反射调用任意函数 |
| `Serial(fns...)` / `SerialWithError(fns...)` | 串行执行 |
| `Parallel(fns...)` | 并行执行，返回等待函数 |
| `ParallelWithError(fns...)` / `ParallelWithErrorChan(fns...)` | 并行并收集错误 |
| `RestrictParallel(restrict, fns...)` | 限流并发 |
| `RestrictParallelWithErrorChan(concurrency, fns...)` | 限流并发收集错误 |

```go
import "github.com/abulo/ratel/v3/core/goroutine"

goroutine.Go(func() {
	// 安全执行，panic 不会导致进程退出
})

// 并行执行并等待
wait := goroutine.Parallel(
	func() { processA() },
	func() { processB() },
	func() { processC() },
)
wait()
```

### 生命周期 cycle 与 hooks

```go
// cycle：并发生命周期控制器
import "github.com/abulo/ratel/v3/core/cycle"

cyc := cycle.NewCycle()
cyc.Run(func() error { /* 后台任务 */ })
cyc.Run(func() error { /* 另一个后台任务 */ })
// 等待全部任务结束并获取错误
if err := <-cyc.Wait(); err != nil { ... }
```

```go
// hooks：分阶段钩子（逆序执行，defer 语义）
import "github.com/abulo/ratel/v3/core/hooks"

hooks.Register(hooks.Stage_BeforeRun, func() { /* ... */ })
hooks.Register(hooks.Stage_AfterStop, func() { /* ... */ })
hooks.Do(hooks.Stage_BeforeRun)
```

---

## 可观测性

### 日志 logger

基于 logrus 的进程内全局日志器。

```go
import "github.com/abulo/ratel/v3/core/logger"

logger.Logger.Info("application started")
logger.Logger.WithFields(logrus.Fields{
	"user_id": 123,
	"action":  "login",
}).Error("operation failed")
```

**多后端异步输出（钩子）：**

```go
import (
	"github.com/abulo/ratel/v3/core/logger"
	eslogger "github.com/abulo/ratel/v3/core/logger/es"
	mongologger "github.com/abulo/ratel/v3/core/logger/mongo"
	sqllogger "github.com/abulo/ratel/v3/core/logger/sql"
)

// Elasticsearch 钩子（构造失败会 panic）
esHook := eslogger.New(eslogger.WithClient(esClient))
logger.Logger.AddHook(esHook)

// MongoDB 钩子
mongoHook := mongologger.New(mongologger.WithClient(mongoClient))
logger.Logger.AddHook(mongoHook)

// SQL 钩子
sqlHook := sqllogger.New(sqllogger.WithDB(gormDB))
logger.Logger.AddHook(sqlHook)
```

日志条目由 `core/logger/queue` 的异步队列在后台写入，不阻塞业务调用；`core/logger/entry` 定义标准化的日志条目结构（主机名、时间戳、调用位置、级别、附加字段）。

### 指标 metric

基于 Prometheus 的指标采集封装，命名空间缺省为 `ratel`，构造时自动注册到默认注册器。

```go
import "github.com/abulo/ratel/v3/core/metric"

// 预置向量（标签按定义顺序：type、method、client、code）
metric.ServerHandleCounter.Inc("http", "GET", "my-service", "200")

// 自定义向量
counter := metric.NewCounterVec("my_metric_total", []string{"label"})
counter.Inc("value")
```

**预置指标：**

| 指标 | 说明 |
|---|---|
| `ServerHandleCounter` / `ServerHandleSeconds` | 服务端请求总数 / 耗时 |
| `ClientHandleCounter` / `ClientHandleSeconds` | 客户端调用总数 / 耗时 |
| `JobHandleCounter` / `JobHandleSeconds` | 任务执行总数 / 耗时 |

### 链路追踪 trace

OpenTelemetry 链路追踪封装。

```go
import (
	"github.com/abulo/ratel/v3/core/trace"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// 初始化全局 TracerProvider
cfg := trace.NewConfig(
	trace.WithName("my-service"),
	trace.WithEndpoint("127.0.0.1:4317"), // OTLP 端点
	trace.WithProtocol("grpc"),           // grpc / http
	trace.WithSampler(1.0),               // 采样率 0-1
)
trace.SetGlobalTracer(cfg.Build())

// 在业务代码中创建 span（carrier 用于注入/提取传播上下文）
tracer := trace.NewTracer(oteltrace.SpanKindServer)
ctx, span := tracer.Start(ctx, "operation", propagation.MapCarrier{})
defer span.End()
```

支持 HTTP/gRPC 两种 OTLP 导出协议，兼容 Jaeger 与 W3C TraceContext 传播格式；提供面向服务端、客户端、生产者、消费者角色的 Tracer 包装。中间件自动注入：xgin/xhertz 服务端中间件、xgrpc 拦截器、各存储与客户端钩子会自动携带追踪上下文。

---

## 任务调度

### crontab 调度器

基于 cron 表达式的定时任务调度器。

```go
import "github.com/abulo/ratel/v3/crontab"

// 1. 创建任务管理器
manager := crontab.NewTaskManager()

// 2. 注册任务（五/六字段 cron 或 @yearly/@monthly 等描述符）
manager.AddTask("hourly-job", crontab.NewTask("hourly-job", "0 * * * *", func(ctx context.Context) error {
	// 任务逻辑
	return nil
}, crontab.TimeoutOption(30*time.Second))) // 可选超时控制

// 3. 启动调度
manager.StartTask()

// 4. 优雅关闭
<-manager.GracefulShutdown()
```

**TaskManager 方法：**

| 方法 | 说明 |
|---|---|
| `AddTask(name, task)` | 注册任务 |
| `DeleteTask(name)` / `ClearTask()` | 删除任务 |
| `StartTask()` / `StopTask()` | 启动/停止调度 |
| `GracefulShutdown()` | 等待全部任务结束 |
| `WorkerStart()` / `WorkerStop()` | 实现 worker.Worker 接口，可被应用托管 |

### 分布式任务 core/task

基于 cron 调度器与一致性哈希节点池的分布式定时任务，在集群内分配任务，避免重复执行。

```go
import (
	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/task"
	"github.com/abulo/ratel/v3/core/task/driver"
)

// 1. 创建驱动（Redis 普通键 / Redis 有序集合 / etcd 租约）
etcdClient, _ := etcdv3.NewClient(etcdv3.New().SetEndpoints([]string{"127.0.0.1:2379"}))
drv := driver.NewEtcdDriver(etcdClient)
// 或 drv := driver.NewRedisDriver(redisClient)
// 或 drv := driver.NewRedisZSetDriver(redisClient)

// 2. 构造分布式任务管理器（CronOptionSeconds 等选项经 NewTaskWithOption 传入）
t := task.NewTaskWithOption("order-service", drv, task.CronOptionSeconds())

// 3. 注册任务（AddFunc/AddJob；cron 表达式按 6 段含秒解析）
t.AddFunc("sync-order", "0 */5 * * * *", func() {
	// 任务逻辑，集群内仅一个节点执行
})

// 4. 启动
t.Start()
```

**核心概念：**

- **节点池（NodePool）**：基于一致性哈希（可配副本数 `WithHashReplicas`）在集群成员间分配任务；成员变更时自动进入升级状态（`NodePoolStateUpgrading`）避免任务重复执行。
- **任务恢复**：启动时支持恢复（`WithRecoverFunc`）与近期任务回放（`NewRecentJobPacker`）。
- **驱动（driver）**：`NewRedisDriver`/`NewRedisZSetDriver`/`NewEtcdDriver` 三种节点注册与发现实现。

### 分布式执行器 executor

xxl-job 风格的任务执行器接口与全局管理。

```go
import "github.com/abulo/ratel/v3/core/executor"

// 1. 实现执行器接口（GetAddress/RegXJob/Run/Stop/GracefulStop）
type MyExecutor struct {
	address string
}

func (e *MyExecutor) GetAddress() string { return e.address }
func (e *MyExecutor) RegXJob(jobs ...executor.XJob) {}
func (e *MyExecutor) Run() error                      { return nil }
func (e *MyExecutor) Stop()                           {}
func (e *MyExecutor) GracefulStop()                   {}

// 2. 注册执行器并启动
executor.Register(":9999", &MyExecutor{address: ":9999"})
if err := executor.Run(); err != nil { ... }

// 3. 停止
// executor.Stop() 或 executor.GracefulStop()
```

### worker 与 job 接口

```go
// worker.Worker：后台工作单元（应用托管）
type MyWorker struct{}

func (w *MyWorker) WorkerStart() error { /* 启动 */ }
func (w *MyWorker) WorkerStop() error  { /* 停止 */ }

application, _ := app.New()
_ = application.Schedule(&MyWorker{}) // 注册后台工作单元

// job.Runner：启动后并发执行的任务（Run 无返回值，通常为常驻循环）
type MyJob struct{}

func (j *MyJob) Run() { /* 任务逻辑 */ }

_ = application.Job(&MyJob{}) // 注册任务
```

---

## 工具包

### util

通用工具函数集合（纯函数，可并发调用）。

```go
import "github.com/abulo/ratel/v3/util"

// 字符串与切片
util.InArray("x", []string{"a", "b", "x"})      // true
util.ArrayRemoveRepeatedElement([]string{"a", "a", "b"}) // ["a","b"]
util.Implode(",", []string{"a", "b"})           // "a,b"
util.Explode(",", "a,b")                        // ["a","b"]

// JSON
util.JSONString(obj)              // 序列化（去装箱优化）
util.JSONDecode([]byte(s), &dst)  // 反序列化

// 时间
util.Timestamp()                  // 当前时间戳
util.FormatDateTime(time.Now())   // 格式化时间
util.DateTimeParse(s)             // 解析时间

// 数学
util.Max(1, 2) / util.Min(1, 2)   // 最大值 / 最小值（float64 可变参数）
util.RandomInt(0, 100)            // 随机整数（(int, error)）
bytes, _ := util.RandomBytes(32)  // 随机字节（([]byte, error)）

// 文件与网络
util.FileExists(path)             // 文件是否存在
util.GetHostByName("example.com") // 域名解析

// 其它
util.ZhCharToPinyin("中国")        // 拼音转换
util.Uniqid("prefix_")            // 唯一 ID
util.Md5(s) / util.Sha1(s)        // 摘要
util.MapToStruct(map[string]any)  // 结构转换（返回 *pb.Struct）
util.StructToMap(pbStruct)        // 结构转换（返回 map[string]any）
```

涵盖字符串、切片与映射、时间日期、数学计算、文件系统、网络、URL、正则信息提取、随机数、拼音以及 protobuf 结构转换等能力。

### snowflake

雪花算法 int64 唯一 ID 生成器。

```go
import "github.com/abulo/ratel/v3/snowflake"

// 一般应用配置：时间戳 40 位 + 序列 7 位 + 机器 16 位
config := snowflake.CommonConfig

// 或自定义位宽与机器编号
config = snowflake.ConstructConfigWithMachineID(40, 7, 16, 1)

// 生成 ID（原子操作，并发安全）
id := config.GenInt64ID()
```

### correct

中英文混排文本的自动排版与清理。

```go
import "github.com/abulo/ratel/v3/correct"

// 排版：CJK 与英文/数字/符号之间自动补空格，规范化全角字符
out := correct.Format("这是一段English文本，包含123数字。")
// 这是一段 English 文本，包含 123 数字。

// 排版 HTML（只处理标签之间的文本节点，保留标签与代码块）
out, err := correct.FormatHTML(htmlBody)

// 反排版（去除自动补充的空格）
orig := correct.Unformat(out)
```

### nlpword / filter

基于 AC 自动机的敏感词处理工具（两套实现）。

```go
// nlpword：检测、过滤、替换
import "github.com/abulo/ratel/v3/nlpword"

filter := nlpword.New()
_ = filter.LoadWordDict("words.txt")      // 从文件加载
// filter.LoadNetWordDict(url)             // 从网络加载
// filter.AddWord("敏感词")                // 程序内添加

replaced := filter.Replace("这句话包含敏感词", '*') // 替换命中词
ok, hit := filter.Validate("测试文本")             // 是否含敏感词（返回命中词）
// 或 filter.FindIn(text) (bool, string) / filter.FindAll(text) []string

// filter：扫描器（词表/文件构建，自定义跳过分隔符）
import "github.com/abulo/ratel/v3/filter"

search := filter.Strings([]string{"敏感词1", "敏感词2"}, "@#$%^&*!?，。；：、")
has := search.HasSens([]byte("文本内容"))           // 是否含敏感词
hits := search.Find([]byte("文本内容"))             // 全部命中收集
replaced = string(search.ReplaceRune([]byte("文本内容"), '*'))
```

### i18n

多语言消息管理。

```go
import "github.com/abulo/ratel/v3/lang/i18n"

// 初始化默认实例：语言目录 + 默认语言 + 语言映射（键=语言名，值=显示名）
i18n.Init("./locales", "zh-CN", map[string]string{
	"zh-CN": "简体中文",
	"en-US": "English",
})
// 默认 FileMode 下按 langDir/<语言名>.ini 加载语言文件（DirMode 按语言目录加载）

// 按指定语言翻译（缺失键自动回退到备用语言，支持参数渲染）
msg := i18n.Tr("zh-CN", "hello", "Ratel") // "你好, Ratel"
// 或使用默认语言：i18n.Dt("hello", "Ratel")
// 也可自行构建实例：instance := i18n.NewWithInit(langDir, defLang, languages)
```

### img

图像读取、变换、合成与保存（基于 libvips，需安装依赖）。

```go
import (
	"github.com/abulo/ratel/v3/img"
	"github.com/disintegration/imaging"
)

image, err := img.OpenLocalFile("input.jpg")
if err != nil {
	panic(err)
}

// 链式变换（就地修改源图，方法返回自身）：等比缩放 → 居中裁剪 → 高斯模糊
image.Resize(800, 0, imaging.Lanczos)
image.Crop(800, 600, imaging.Center)
image.Blur(5)

// 叠加水印
wm, _ := img.OpenLocalFile("logo.png")
image.AddWaterMark(wm, imaging.Center, 20, 20, 0.5)

// 保存（按文件扩展名推断编码格式；JPEG 质量可用 imaging.JPEGQuality 控制）
if err := img.SaveToFile(image, "output.jpg", imaging.JPEGQuality(80)); err != nil {
	panic(err)
}
```

支持缩放、裁剪、模糊、灰度、对比度/亮度/饱和度调整、旋转、粘贴、叠加、水印、JPEG 压缩，以及基于 `img/fontx` 的文字绘制；方法返回自身支持链式调用。

### watch 热编译

Go 项目热编译监视（文件变更自动构建并重启）。

```go
import "github.com/abulo/ratel/v3/watch"

err := watch.Watch(ctx, &watch.Options{
	// 配置监听目录、入口文件、输出名称、文件后缀与编译参数
})
```

### config 配置读取

基于 Viper 的配置读取封装。

```go
import "github.com/abulo/ratel/v3/config"

conf := config.New()
_ = conf.LoadDir("./configs")     // 目录批量加载
// _ = conf.LoadFile("config.yaml") // 或加载单个文件

// 类型化读取（缺失时返回默认值）
port := conf.Int("server.port", 8080)
name := conf.String("server.name", "ratel")
timeout := conf.Duration("server.timeout", 5*time.Second)
```

支持从单个文件或目录批量加载 JSON、TOML、YAML 配置。

---

## 命令行工具 toolkit

`toolkit` 是 ratel 的命令行脚手架程序（基于 Cobra），用于项目创建与代码生成。

```bash
# 安装
go install github.com/abulo/ratel/v3/toolkit@latest

# 查看版本与帮助
toolkit version
toolkit --help
```

**子命令：**

| 命令 | 功能 |
|---|---|
| `toolkit project` | 项目脚手架：从远程 layout 仓库克隆模板，生成后端/前端项目，或在现有项目中追加服务目录 |
| `toolkit module` | 数据模型层代码生成：交互式采集模型路径、表名、软删除策略，生成 module/proto/service/convert 四类代码 |
| `toolkit dao` | DAO 代码生成：读取表元数据，生成 GORM 模型结构体（支持单表/多表） |
| `toolkit api` | HTTP 接口层代码生成：按 Web 框架（gin/hertz）生成接口处理函数与路由注册片段 |
| `toolkit vue` | Vue3 + TypeScript 前端代码生成：接口定义、API 请求方法模块、列表页面 |
| `toolkit build` | 热编译：监听目录文件变更，自动重新编译并重启 |
| `toolkit changelog` | 更新日志生成：按约定式提交类型分类整理，写入 CHANGELOG.md |
| `toolkit upgrade` / `toolkit init` | 脚手架升级与初始化：安装最新 toolkit 本体及 protoc 工具 |

```bash
# 示例：交互式创建项目
toolkit project

# 示例：生成 DAO（指定表）
toolkit dao --table users,orders
```

---

## 项目结构

```
ratel/
├── client/              # 客户端封装
│   ├── etcdv3/         # etcd 客户端（KV/锁/租约/监听）
│   ├── grpc/           # gRPC 客户端（服务发现/负载均衡/拦截器）
│   ├── grpc/balancer/  # SWR、P2C 负载均衡器
│   ├── rabbitmq/       # RabbitMQ 客户端（断线重连）
│   ├── ftp/            # FTP 客户端
│   └── sftp/           # SFTP 客户端
├── config/             # Viper 配置读取
├── core/               # 核心治理能力
│   ├── app/            # 应用生命周期管理
│   ├── logger/         # 日志（es/mongo/sql 后端 + queue 异步队列）
│   ├── metric/         # Prometheus 指标
│   ├── ratelimit/      # 令牌桶限流
│   ├── resource/       # 熔断器、SingleFlight、批量错误、资源管理
│   ├── task/           # 分布式任务与 cron 调度
│   ├── executor/       # xxl-job 风格分布式执行器
│   ├── trace/          # OpenTelemetry 链路追踪
│   ├── goroutine/      # goroutine 编排
│   ├── cycle/          # 并发生命周期控制
│   ├── hooks/          # 生命周期阶段钩子
│   ├── ecode/          # 统一错误码模型
│   ├── singleton/      # 单例存储
│   ├── constant/       # 全局常量
│   ├── env/            # 运行环境元信息
│   ├── hostname/       # 主机名获取
│   ├── call/           # 调用栈信息
│   ├── component/      # 组件抽象
│   └── stack/          # 调用栈工具
├── server/             # 服务端实现
│   ├── xgin/           # Gin HTTP 服务端
│   ├── xgrpc/          # gRPC 服务端
│   ├── xhertz/         # Hertz HTTP 服务端
│   ├── xasynq/         # Asynq 定时任务（manager/worker）
│   └── xmonitor/       # 监控服务端（pprof 等）
├── stores/             # 数据存储
│   ├── redis/          # Redis（单机/集群/哨兵/分片）
│   ├── sql/            # SQL（MySQL/PostgreSQL/ClickHouse）
│   ├── mongodb/        # MongoDB
│   ├── elasticsearch/  # Elasticsearch
│   ├── null/           # 可空类型集
│   ├── proxy/          # 存储代理池
│   └── session/        # Redis 会话存储
├── registry/           # 服务注册与发现（etcdv3 实现）
├── toolkit/            # 命令行脚手架
├── util/               # 通用工具
├── correct/            # 中英文混排校正
├── nlpword/            # 敏感词处理（AC 自动机）
├── filter/             # 敏感词扫描
├── img/                # 图像处理（libvips）
├── lang/               # i18n 与占位类型
├── snowflake/          # 雪花 ID
├── crontab/            # cron 调度器
├── watch/              # 热编译监视
└── worker/             # 后台工作单元与任务
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

# 运行单元测试
make test

# 静态检查
make vet       # go vet
make errcheck  # 错误检查
make lint      # revive lint

# 格式化与格式检查
make fmt       # gofmt -w
make fmtcheck  # gofmt 检查

# 完整构建（fmt + errcheck + lint + build）
make all
```

### Makefile 目标

| 目标 | 说明 |
|---|---|
| `make all` | 完整构建（fmt → errcheck → lint → build） |
| `make test` | 运行单元测试（`go test -short`） |
| `make vet` | 静态检查（`go vet ./...`） |
| `make lint` | revive 检查（自动安装 revive） |
| `make fmt` / `make fmtcheck` | 格式化 / 格式检查 |
| `make errcheck` | 错误处理检查（自动安装 errcheck） |
| `make cmt` | 导出函数自动注释（gocmt） |

### 代码规范

- 使用 `gofmt` 格式化代码，`revive` 做 Lint 检查。
- 导出标识符编写 Godoc 文档注释（包注释统一收敛到 `doc.go`，以 `Package xxx` 开头）。
- 检查所有错误返回，错误处理不吞异常。
- 编写单元测试，覆盖关键路径。
- 提交信息遵循约定式提交：`feat`/`fix`/`docs`/`perf`/`chore`/`refactor`/`test` 等。

### 持续集成

仓库配置了 GitHub Actions（`.github/workflows/go.yml`），在 push 与 PR 时自动构建并上报 Codecov；`release.yml` 负责发布流程。

---

## 常见问题

### Q1: Ratel 与原生 Gin/gRPC 有什么区别？

Ratel 在原生框架之上提供统一封装：Server 生命周期管理、自动中间件（指标/追踪/慢查询）、注册中心集成、治理能力（限流/熔断）与配套工具，业务代码无需关心这些基础设施。

### Q2: 如何迁移现有项目到 Ratel？

Ratel 服务端内嵌原生引擎，Gin/gRPC/Hertz 原生 API 全部可用：

```go
// 替换导入即可，其余代码无需修改
// import "github.com/gin-gonic/gin"
import "github.com/abulo/ratel/v3/server/xgin"
```

### Q3: 如何自定义中间件？

```go
// xgin（内嵌 gin.Engine）
server.Use(func(c *gin.Context) {
	// 自定义逻辑
	c.Next()
})

// xgrpc
cfg.WithUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// 自定义逻辑
	return handler(ctx, req)
})
```

### Q4: 如何启用限流？

```go
// 在业务代码中直接使用令牌桶
bucket := ratelimit.NewBucketWithRate(1000, 100)
if bucket.TakeAvailable(1) == 0 {
	return errors.New("rate limit exceeded")
}
```

### Q5: 如何配置服务注册？

```go
import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/registry"
	regEtcd "github.com/abulo/ratel/v3/registry/etcdv3"
)

// 方式一：使用包级默认注册器（初始为本地实现）
application, _ := app.New()
application.SetRegistry(registry.DefaultRegisterer)

// 方式二：生产环境接入 etcd 注册中心
// （Set 方法会返回内嵌的 client 配置类型，应分开调用保持类型）
regConfig := regEtcd.New()
regConfig.SetEndpoints([]string{"127.0.0.1:2379"})
reg, err := regConfig.Build() // 或 regConfig.MustBuild()
if err != nil {
	panic(err)
}
application.SetRegistry(reg)
// Server 启动后自动向注册中心登记
```

### Q6: 如何实现优雅停机？

`application.Run()` 已内置信号监听与优雅停机：收到 SIGINT/SIGTERM 后自动执行 `GracefulStop` 流程（停止 Server/Worker、等待后台任务结束、触发钩子）。

### Q7: 哪些模块需要额外系统依赖？

仅 `img` 图像处理包需要 libvips 与 pkg-config；其余包均为纯 Go 依赖。

---

## 贡献

欢迎以 Issue、Pull Request 或 Discussion 的形式参与贡献。提交 PR 前请确保：

- 通过 `make all` 全部检查。
- 新增或修改的代码带有相应单元测试。
- 导出 API 补充 Godoc 注释，包注释写入 `doc.go`。
- 保持向后兼容；破坏性变更需在 PR 描述中说明。

---

## 许可证

Ratel 采用 [Apache License 2.0](LICENSE)。
