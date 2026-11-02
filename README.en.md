# Ratel

[![Go](https://github.com/abulo/ratel/v3/workflows/Go/badge.svg?branch=v3)](https://github.com/abulo/ratel/v3/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abulo/ratel/v3)](https://goreportcard.com/report/github.com/abulo/ratel/v3)
[![goproxy](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)
[![codecov](https://codecov.io/gh/abulo/ratel/branch/v3/graph/badge.svg)](https://codecov.io/gh/abulo/ratel)
[![Release](https://img.shields.io/github/v/release/abulo/ratel.svg?style=flat-square)](https://github.com/abulo/ratel/v3)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Ratel is a production-grade Go microservice governance framework that integrates Web/RPC servers, data storage clients, and governance capabilities such as rate limiting, circuit breaking, distributed tracing, metrics, and logging — out of the box.

- **Servers**: reuses mature frameworks such as Gin, gRPC, Hertz, and Asynq, uniformly wrapped behind the `server.Server` lifecycle interface.
- **Clients & Storage**: unified wrappers for etcd, RabbitMQ, FTP/SFTP, Redis, SQL (MySQL/PostgreSQL/ClickHouse), MongoDB, Elasticsearch, and more.
- **Governance**: token-bucket rate limiting, Google SRE adaptive circuit breaking, SingleFlight request coalescing, OpenTelemetry tracing, Prometheus metrics, and multi-backend async logging.
- **Companion tools**: distributed task scheduling, Snowflake IDs, sensitive-word filtering, text typography, i18n, image processing, code generation scaffolding, and hot compilation.

The project is organized in the layered style of Kratos, uses Go 1.27 generics and the standard library, and follows the official Go documentation conventions (every library package has a `doc.go` package comment; all exported APIs carry Godoc comments).

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Gin HTTP Server](#1-gin-http-server)
  - [gRPC Server](#2-grpc-server)
  - [Hertz Server](#3-hertz-server)
  - [Composing Multiple Servers](#4-composing-multiple-servers)
- [Application Lifecycle](#application-lifecycle)
- [Server Packages](#server-packages)
  - [xgin](#xgin)
  - [xgrpc](#xgrpc)
  - [xhertz](#xhertz)
  - [xasynq](#xasynq)
  - [xmonitor](#xmonitor)
- [Client Packages](#client-packages)
  - [etcdv3](#etcdv3)
  - [gRPC Client](#grpc-client)
  - [rabbitmq](#rabbitmq)
  - [ftp / sftp](#ftp--sftp)
- [Storage Packages](#storage-packages)
  - [redis](#redis)
  - [sql](#sql)
  - [mongodb](#mongodb)
  - [elasticsearch](#elasticsearch)
  - [null Nullable Types](#null-nullable-types)
  - [proxy Proxy Pool](#proxy-proxy-pool)
  - [session Storage](#session-storage)
- [Governance Capabilities](#governance-capabilities)
  - [Rate Limiting ratelimit](#rate-limiting-ratelimit)
  - [Circuit Breaker resource](#circuit-breaker-resource)
  - [goroutine Orchestration](#goroutine-orchestration)
  - [Lifecycle cycle & hooks](#lifecycle-cycle--hooks)
- [Observability](#observability)
  - [Logging logger](#logging-logger)
  - [Metrics metric](#metrics-metric)
  - [Tracing trace](#tracing-trace)
- [Task Scheduling](#task-scheduling)
  - [crontab Scheduler](#crontab-scheduler)
  - [Distributed Tasks core/task](#distributed-tasks-coretask)
  - [Distributed Executor executor](#distributed-executor-executor)
  - [worker & job Interfaces](#worker--job-interfaces)
- [Utility Packages](#utility-packages)
  - [util](#util)
  - [snowflake](#snowflake)
  - [correct](#correct)
  - [nlpword / filter](#nlpword--filter)
  - [i18n](#i18n)
  - [img](#img)
  - [watch Hot Compilation](#watch-hot-compilation)
  - [config Configuration](#config-configuration)
- [CLI Toolkit](#cli-toolkit)
- [Project Structure](#project-structure)
- [Development Guide](#development-guide)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## Features

### Multi-Protocol Servers

| Server | Underlying framework | Use case |
|---|---|---|
| `server/xgin` | Gin | RESTful APIs, web applications, WebSocket |
| `server/xgrpc` | gRPC | RPC communication between microservices |
| `server/xhertz` | CloudWeGo Hertz | High-concurrency HTTP services |
| `server/xasynq` | Asynq | Scheduled task dispatching and background task consumption |
| `server/xmonitor` | net/http | Governance endpoints: pprof, route listing, build info |

All servers uniformly implement the `server.Server` interface (`Serve`/`Stop`/`GracefulStop`/`Info`/`Health`) and are managed by `core/app`, with signal-driven graceful start/stop.

### Unified Clients & Storage

- **etcd**: KV read/write, prefix/multi-key batch queries, distributed mutex locks, lease sessions, prefix watching, BasicAuth/TLS, distributed tracing.
- **gRPC client**: etcd service discovery, SWR smooth weighted round-robin and P2C least-loaded load balancing, debug/timeout/tracing/logging/metrics interceptor chain.
- **RabbitMQ**: automatic reconnection on disconnect, channel reuse, queue declare/bind builder, configurable retry policies, automatic consumption recovery after disconnect.
- **Redis**: standalone, cluster, sentinel, and ring modes; commands grouped by data type; all executed through the circuit breaker; optional tracing and metrics hooks.
- **SQL**: MySQL, PostgreSQL, ClickHouse; based on GORM; DSN-level connection cache reuse; pagination/sorting/scope parameter helpers; metrics and tracing interceptors.
- **MongoDB**: chainable condition builders, circuit breaker/metrics/tracing enhanced collections, automatic struct-to-BSON conversion on write.
- **Elasticsearch**: node cache reuse, HTTP transport-level tracing and metrics.
- **Nullable types**: a complete nullable type set replacing `sql.NullXxx`, uniformly implementing Scanner/Valuer/JSON/Text interfaces.
- **Proxy pool**: group-based unified management of Redis/Mongo/ES/SQL instances; SQL supports read/write separation with random selection.
- **Session storage**: Redis-based JSON sessions with dotted-path nested field read/write.

### Service Governance

- **Rate limiting**: token bucket with three fill modes (fixed interval, fixed rate, per-tick quantum); blocking/timed/non-blocking token acquisition; rate-limited Reader/Writer.
- **Circuit breaking**: adaptive circuit breaker based on the Google SRE overload protection model, with custom error judgment and fallback callbacks.
- **Request coalescing**: SingleFlight concurrency deduplication.
- **Batch errors**: BatchError aggregation.
- **goroutine orchestration**: serial/parallel/rate-limited concurrency, delayed start, panic recovery, error collection.
- **Lifecycle**: concurrent lifecycle controller (cycle), staged hooks (hooks, executed in reverse order with defer semantics).

### Observability

- **Metrics**: Prometheus namespace `ratel`, with pre-built Counter, Gauge, Histogram, and Summary vectors for server/client/task/cache scenarios.
- **Tracing**: OpenTelemetry with HTTP/gRPC OTLP export protocols, compatible with Jaeger and W3C TraceContext propagation.
- **Logging**: logrus-based structured logging; console color output; Elasticsearch, MongoDB, and SQL async hooks writing via background queues.

### Tasks & Concurrency

- **cron scheduling**: five/six-field expressions and predefined descriptors such as `@yearly`, with timeout control, dynamic add/remove, and graceful shutdown.
- **Distributed tasks**: cron scheduling + consistent-hash node pool, tasks distributed within the cluster; node changes automatically enter an upgrading state to avoid duplicate execution; supports task recovery and recent-job replay.
- **Distributed executor**: xxl-job style task executor interface and global management.
- **Asynq**: periodic task manager and background task worker.

### Companion Utilities

- Snowflake ID generator (atomic, concurrency-safe, configurable bit layout).
- Aho-Corasick sensitive-word filtering (two implementations: `nlpword` and `filter`).
- Automatic CJK/Latin mixed text typography (`correct`).
- i18n multilingual message management (`lang/i18n`).
- Image read/transform/composite/watermark/text drawing (`img`, based on libvips).
- Go project hot compilation watcher (`watch`).
- Code generation scaffolding (`toolkit`): project creation, DAO/interface/model/Vue code generation, changelog, scaffold upgrade.

---

## Architecture

### Layered Architecture

```mermaid
flowchart TB
    %% ============ Layers: top to bottom ============

    subgraph APP["Application Layer core/app — Component Assembly & Lifecycle"]
        App["Application<br/>New / Run / Serve / Schedule / Job / Executor<br/>RegisterHooks / SetRegistry / Stop / GracefulStop"]
    end

    subgraph SRV["Server Layer server/ — Unified server.Server Lifecycle Interface"]
        Xgin["xgin<br/>Gin HTTP server"]
        Xgrpc["xgrpc<br/>gRPC server"]
        Xhertz["xhertz<br/>Hertz HTTP server"]
        Xasynq["xasynq<br/>Asynq async task server"]
        Xmonitor["xmonitor<br/>Monitor / pprof server"]
    end

    subgraph CLI["Client Layer client/ — Middleware Clients"]
        CEtcd["etcdv3<br/>Config center / registry client"]
        CGrpc["grpc<br/>Service discovery + interceptor chain"]
        CRabbit["rabbitmq<br/>Pub/sub client"]
        CFtp["ftp / sftp<br/>File transfer client"]
    end

    subgraph STO["Storage Layer stores/ — Data Access"]
        SRedis["redis<br/>Standalone / cluster / sentinel / ring"]
        SSql["sql<br/>GORM wrapper · MySQL / PostgreSQL / ClickHouse"]
        SMongo["mongodb<br/>mongo-driver wrapper"]
        SEs["elasticsearch"]
        SNull["null<br/>Database nullable types"]
        SProxy["proxy<br/>Multi-client proxy pool"]
        SSession["session<br/>Redis session store"]
    end

    subgraph REG["Registry registry/"]
        RInterface["registry.Registry interface<br/>Register / Deregister / GetService"]
        REtcd["etcdv3 implementation<br/>Lease renewal · service discovery"]
    end

    subgraph CORE["Governance Core core/ — Cross-cutting Capabilities"]
        CRatelimit["ratelimit Token-bucket rate limiting"]
        CResource["resource Circuit breaker / SingleFlight / resource manager"]
        CTrace["trace OpenTelemetry tracing"]
        CMetric["metric Prometheus metrics"]
        CLogger["logger Multi-backend async logging"]
        CGoroutine["goroutine Concurrency orchestration"]
        CTask["task Distributed scheduled tasks"]
        CCycle["cycle Lifecycle control"]
        CHooks["hooks Stage hooks"]
        CExecutor["executor Distributed executor"]
    end

    subgraph TOOL["Tool Layer util/ etc — Common Capabilities"]
        TUtil["util Utility functions"]
        TCorrect["correct Text typography"]
        TNlp["nlpword Sensitive word detection"]
        TFilter["filter Text filtering"]
        TLang["lang / i18n"]
        TImg["img Image processing"]
        TSnow["snowflake Snowflake ID"]
        TCron["crontab Scheduler"]
        TWatch["watch File watching"]
        TConfig["config Config loading"]
    end

    subgraph INFRA["External Infrastructure"]
        IEtcd["etcd cluster"]
        IRabbit["RabbitMQ"]
        IRedis["Redis"]
        ISql["MySQL / PostgreSQL / ClickHouse"]
        IMongo["MongoDB"]
        IEs["Elasticsearch"]
        IFtp["FTP / SFTP"]
        IOtel["OpenTelemetry Collector"]
        IProm["Prometheus"]
    end

    %% ============ Lifecycle management ============
    APP -->|"Start via Run"| SRV
    APP -->|"Managed via Schedule"| Xasynq
    APP -->|"Managed via Job"| CTask
    APP -->|"Managed via Executor"| CExecutor
    APP -->|"Stage hooks"| CHooks
    APP -->|"Inject via SetRegistry"| REG

    %% ============ Server call chain ============
    Xgin -->|"Business calls"| CGrpc
    Xgin -->|"Business calls"| CRabbit
    Xgrpc -->|"Business calls"| CGrpc
    Xhertz -->|"Business calls"| CGrpc
    Xasynq -->|"Consume tasks"| SRedis

    %% ============ Data access ============
    SRV -->|"Data access"| STO
    SSession -->|"Built on"| SRedis
    SProxy -->|"Unified access"| SRedis
    SProxy -->|"Unified access"| SMongo
    SProxy -->|"Unified access"| SEs
    SProxy -->|"Unified access"| SSql

    %% ============ Service discovery / config ============
    CGrpc -->|"Service discovery"| REG
    SRV -->|"Register on start"| REG
    REtcd -->|"Read/write nodes"| IEtcd
    CEtcd -->|"Read/write"| IEtcd

    %% ============ Middleware connections ============
    CRabbit -->|"Connect"| IRabbit
    SRedis -->|"Connect"| IRedis
    SSql -->|"Connect"| ISql
    SMongo -->|"Connect"| IMongo
    SEs -->|"Connect"| IEs
    CFtp -->|"Connect"| IFtp

    %% ============ Governance cross-cutting (dashed) ============
    STO -.->|"Commands via breaker / metrics / tracing"| CORE
    CLI -.->|"Interceptor chain (tracing / metrics / access log / timeout)"| CORE
    CORE -->|"Trace export"| IOtel
    CORE -->|"Metrics expose"| IProm

    %% ============ Task scheduling internal dependencies ============
    CTask -->|"Generate IDs"| TSnow
    CTask -->|"Node registration"| CEtcd
    TCron -->|"Underlying scheduling"| CTask
    TConfig -->|"Config source"| CEtcd

    %% ============ Styles ============
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

**Layer overview:**

| Layer | Directory | Responsibility |
|---|---|---|
| Application layer | `core/app` | Assembles Server/Worker/Job/Executor, signal-driven start/stop, stage hooks, registry injection |
| Server layer | `server/` | Unified wrappers for Gin / gRPC / Hertz / Asynq / monitor, implementing the `server.Server` interface |
| Client layer | `client/` | etcd, gRPC, RabbitMQ, FTP/SFTP clients with interceptor chains and connection management |
| Storage layer | `stores/` | Redis (four modes), SQL (GORM), MongoDB, Elasticsearch, nullable types, proxy pool, session |
| Registry | `registry/` | `Registry` interface and the etcdv3 implementation (lease, service discovery) |
| Governance core | `core/` | Rate limiting, circuit breaking, tracing, metrics, logging, concurrency orchestration, distributed tasks, lifecycle, hooks, executor |
| Tool layer | `util/` etc. | Utility functions, text processing, image, ID generation, config loading, file watching, i18n |
| External infrastructure | — | etcd / RabbitMQ / Redis / SQL / MongoDB / ES / FTP·SFTP / OTel / Prometheus |

**Key paths:**

- **Business request**: the application layer starts servers via `Run` → servers call downstream services through the client layer and read/write data through the storage layer → every command passes through the governance core (circuit breaker → metrics → tracing).
- **Service discovery**: servers register with the registry on start; gRPC clients resolve target nodes from the registry; the etcdv3 implementation handles lease renewal and node read/write.
- **Observability**: traces are exported to the OpenTelemetry Collector, metrics are exposed to Prometheus, and logs are written asynchronously through the logger's multi-backend.
- **Task scheduling**: `task` depends on `snowflake` for task IDs and on the etcd client for node registration; `crontab` provides the underlying cron scheduling.

### Application Lifecycle

`core/app` is responsible for assembling all components:

1. On application start, registered Servers, Workers, Jobs, and Executors are started in sequence.
2. After a Server starts, it automatically registers with the registry (`registry.DefaultRegisterer`).
3. System signals (SIGINT/SIGTERM, etc.) are listened for; on receipt, cleanup is performed: stop Servers and Workers, shut down background tasks, and trigger stage hooks.
4. Both `Stop` (immediate) and `GracefulStop` (graceful) are supported.

### Design Conventions

- **Configuration as code**: every component is constructed with `New()` plus chainable `WithXxx` options — centralized, readable, and statically checkable.
- **Unified lifecycle**: Server/Worker/Job all implement framework-defined interfaces and are managed uniformly by the application; business code does not need to care about start/stop details.
- **Built-in governance**: storage and client commands pass through the circuit breaker by default, with metrics and tracing attachable via switches, without business-side intrusion.
- **Package comments as documentation**: every library package has a `doc.go` package comment; exported APIs carry Godoc comments, readable via `go doc` / pkg.go.dev.

---

## Installation

### Requirements

| Dependency | Version | Notes |
|---|---|---|
| Go | 1.27+ | see the `go` directive in `go.mod` |
| libvips | optional | only needed by the `img` image processing package; on macOS use `brew install vips` |
| pkg-config | optional | required by libvips when compiling the `img` package |

### Get the module

```bash
go get github.com/abulo/ratel/v3@latest
```

---

## Quick Start

### 1. Gin HTTP Server

```go
package main

import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/server/xgin"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Build the Gin server: listen on 8080, enable metrics and slow-query middleware
	httpServer := xgin.New().
		WithPort(8080).
		WithEnableMetric(true).
		WithEnableSlowQuery(true).
		Build()

	// 2. Register routes: the Server embeds gin.Engine, so all native Gin APIs are available
	httpServer.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Ratel!"})
	})

	// 3. Assemble the application and start: blocks until an exit signal, with graceful shutdown
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	if err := application.Run(httpServer); err != nil {
		panic(err)
	}
}
```

### 2. gRPC Server

```go
package main

import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/server/xgrpc"
)

func main() {
	// Build the gRPC server (the Server embeds grpc.Server; native gRPC APIs are available)
	cfg := xgrpc.New()
	cfg.Host = "0.0.0.0"
	cfg.Port = 50051
	cfg.EnableTrace = true
	cfg.EnableMetric = true
	grpcServer, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	// Register your service implementation, for example:
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

### 3. Hertz Server

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

	// The Server embeds the Hertz engine; native Hertz APIs are available
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

> Note: `core/app` and `hertz/pkg/app` share the same package name; when importing both, alias one of them (the example above uses `ratelapp`).

### 4. Composing Multiple Servers

A single application can host several servers at once:

```go
func main() {
	httpServer := xgin.New().WithPort(8080).Build()
	grpcServer := mustBuildGRPC()
	monitorServer := xmonitor.New().WithPort(9090).Build() // standalone monitor endpoint exposing pprof, etc.

	application, err := app.New()
	if err != nil {
		panic(err)
	}
	// Run accepts any number of server.Server instances
	if err := application.Run(httpServer, grpcServer, monitorServer); err != nil {
		panic(err)
	}
}
```

---

## Application Lifecycle

`core/app` provides the `Application` lifecycle manager.

### Application methods

| Method | Description |
|---|---|
| `New(fns ...func() error) (*Application, error)` | Creates an application; `fns` are initialization functions run before startup |
| `DefaultApp() *Application` | Creates a default application (skips flag parsing and config loading) |
| `Serve(s ...server.Server) error` | Registers servers |
| `Schedule(w worker.Worker) error` | Registers a background worker |
| `Job(runner job.Runner) error` | Registers a task executed concurrently after startup |
| `Executor(e executor.Executor)` | Registers a distributed task executor |
| `SetRegistry(reg registry.Registry)` | Sets the registry |
| `Run(servers ...server.Server) error` | Starts all components and blocks until an exit signal |
| `Stop() error` | Stops the application immediately |
| `GracefulStop(ctx context.Context) error` | Stops the application gracefully |
| `RegisterHooks(stage hooks.Stage, fns ...func())` | Registers lifecycle hooks |
| `WithOptions(options ...Option)` | Application options |

### Lifecycle hooks

```go
import "github.com/abulo/ratel/v3/core/hooks"

func main() {
	// Register a startup hook (executed before config loading)
	hooks.Register(hooks.Stage_BeforeLoadConfig, func() {
		// custom initialization logic
	})

	application, err := app.New()
	// ...
}
```

Hooks are grouped by stage; hooks in the same stage execute in **reverse** registration order (defer semantics).

### Signal handling

`Run` listens for system signals:

- **SIGINT / SIGTERM**: triggers the graceful shutdown flow.
- Shutdown flow: run `BeforeStop` hooks → stop all Servers/Workers → wait for background tasks to finish → run `AfterStop` hooks → close the lifecycle.

---

## Server Packages

All servers uniformly implement the `server.Server` interface and support:

| Method | Description |
|---|---|
| `Serve() error` | Starts the server and blocks listening |
| `Stop() error` | Stops immediately |
| `GracefulStop(ctx context.Context) error` | Stops gracefully (waits for in-flight requests) |
| `Info() *server.ServiceInfo` | Returns registration info (name, address, protocol, tags) |
| `Health() bool` | Health check |

### xgin

Gin-based HTTP server.

```go
server := xgin.New().
	WithHost("0.0.0.0").          // listen host
	WithPort(8080).               // listen port
	WithMode(gin.ReleaseMode).    // run mode (ReleaseMode by default)
	WithDeployment("prod").       // deployment group (traffic isolation for registration)
	WithServiceAddress("").       // address registered with the registry; empty means Host:Port
	WithEnableMetric(true).       // enable metrics collection middleware
	WithEnableTrace(true).        // enable tracing middleware
	WithEnableSlowQuery(true).    // enable slow-query and panic recovery middleware
	WithSlowQueryThresholdInMilli(500). // slow query threshold (ms)
	Build()
```

**Configuration:**

| Field | Default | Description |
|---|---|---|
| `Host` | empty | Listen host |
| `Port` | 0 | Listen port; backfilled from the actual bound port at build time |
| `Mode` | `ReleaseMode` | Gin run mode (`ReleaseMode`/`DebugMode`/`TestMode`) |
| `EnableMetric` | false | Metrics collection middleware |
| `EnableTrace` | false | Server-side tracing middleware |
| `EnableSlowQuery` | false | Slow-request logging and panic recovery middleware |
| `SlowQueryThresholdInMill` | 500 | Slow query threshold (ms) |
| `Deployment` | empty | Deployment group |
| `ServiceAddress` | empty | Address registered with the registry |

**WebSocket:**

```go
type MyWebSocket struct {
	// implement your own protocol upgrade and handling logic
}

// Register a WebSocket handler as a GET route
server.Upgrade(&MyWebSocket{Pattern: "/ws"})
```

### xgrpc

gRPC-based server.

```go
cfg := xgrpc.New()
cfg.Host = "0.0.0.0"
cfg.Port = 50051
cfg.Network = "tcp4"
cfg.EnableTrace = true
cfg.EnableMetric = true
cfg.EnableTLS = false // when enabled, configure CaFile/CertFile/PrivateFile

// Append interceptors
cfg.WithUnaryInterceptor(myUnaryInterceptor)
cfg.WithStreamInterceptor(myStreamInterceptor)

grpcServer, err := cfg.Build()     // returns (*Server, error)
// or grpcServer := cfg.MustBuild() // panics on build failure
```

**Configuration (consistent with the xgrpc style: field assignment + `Build`):**

| Field | Default | Description |
|---|---|---|
| `Host` / `Port` | empty / 0 | Listen address; port backfilled at build time |
| `Network` | `tcp4` | Listen network type |
| `Deployment` | `constant.DefaultDeployment` | Deployment group |
| `EnableTrace` | false | Server-side tracing interceptor |
| `EnableMetric` | false | Metrics collection interceptor |
| `EnableTLS` | false | Mutual TLS authentication |
| `CaFile` / `CertFile` / `PrivateFile` | empty | TLS certificate file paths |
| `SlowQueryThresholdInMill` | 500 | Reserved configuration |

**Panic recovery:** the sub-package `server/xgrpc/recovery` provides unary and stream recovery interceptors; when no callback is configured, panics are wrapped into a `PanicError` and returned, preventing a single request panic from crashing the process.

### xhertz

A high-performance HTTP server based on CloudWeGo Hertz, with configuration aligned with xgin (`WithHost`/`WithPort`/`WithMode`/`WithEnableMetric`/`WithEnableTrace`/`WithEnableSlowQuery`/`WithDeployment`/`WithServiceAddress`/`WithSlowQueryThresholdInMilli`); `Build()` returns a `*Server` that embeds the Hertz engine.

```go
server := xhertz.New().WithPort(8081).Build()
server.GET("/hello", func(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, map[string]string{"message": "hello"})
})
```

### xasynq

An Asynq-based scheduled-task and background-task component with two sub-packages: `manager` and `worker`.

**manager (periodic task manager):** periodically writes tasks to the Redis queue according to registered schedules. Configuration is aligned with xgin/xhertz (`WithHost`/`WithPort`/`WithDeployment`/`WithEnableMetric`/`WithEnableTrace`/`WithMode`, etc.); Asynq runtime options are passed through the `Config.Option` field.

```go
import (
	"github.com/abulo/ratel/v3/server/xasynq/manager"
	"github.com/hibiken/asynq"
)

mgrCfg := manager.New().
	WithHost("0.0.0.0").
	WithPort(0)

// Configure Asynq runtime options before Build (Redis connection and periodic task registry)
mgrCfg.Option = asynq.PeriodicTaskManagerOpts{
	RedisConnOpt: asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
	// register periodic tasks here: cronspec -> task handler
}

mgr := mgrCfg.Build() // implements server.Server
// application.Run(mgr)
```

**worker (task worker):** pulls tasks from the Redis queue and consumes them via the handler routing table. Assign `Config.Redis` (`redis.UniversalClient`), `Config.Mux` (`*asynq.ServeMux`), and `Config.Option` (`asynq.Config`) before Build.

```go
import (
	"github.com/abulo/ratel/v3/server/xasynq/worker"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

wkCfg := worker.New().
	WithHost("0.0.0.0").
	WithPort(0)

// Configure Redis, handler routing table, and runtime options before Build
wkCfg.Redis = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
wkCfg.Mux = asynq.NewServeMux()
wkCfg.Mux.HandleFunc("task:email", func(ctx context.Context, t *asynq.Task) error {
	// handle email task
	return nil
})
wkCfg.Option = asynq.Config{Concurrency: 10}

wk := wkCfg.Build() // implements server.Server
// application.Run(wk)
```

### xmonitor

A standalone monitoring and governance HTTP server that mounts built-in governance endpoints — pprof profiling, registered route listing, and module build info — on its own route table.

```go
monitorServer := xmonitor.New().
	WithPort(9090).
	Build()

// application.Run(httpServer, monitorServer)
```

---

## Client Packages

### etcdv3

An etcd client wrapper based on `go.etcd.io/etcd/client/v3`.

```go
import (
	"github.com/abulo/ratel/v3/client/etcdv3"
)

// Config is filled via field assignment or Set-prefixed methods (not With chains)
cfg := etcdv3.New() // defaults: BasicAuth=false, ConnectTimeout=5s, Secure=false
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

**Client capabilities (embeds the native `clientv3.Client`; native APIs work directly):**

| Method | Description |
|---|---|
| `GetKeyValue(ctx, key)` | Reads a single key |
| `GetValues(ctx, keys...)` | Reads multiple keys in batch |
| `GetPrefix(ctx, prefix)` | Prefix read, returns `map[string]string` |
| `DelPrefix(ctx, prefix)` | Prefix delete |
| `GetLeaseSession(ctx, opts...)` | Creates a lease session (keepalive) |
| `NewMutex(key, opts...)` | Creates a distributed mutex |
| `WatchPrefix(ctx, prefix)` | Watches prefix changes |

**Distributed lock:**

```go
mutex, err := client.NewMutex("/lock/order")
if err != nil {
	panic(err)
}
if err := mutex.Lock(ctx); err != nil {
	panic(err)
}
defer mutex.Unlock(ctx)
// critical section
```

### gRPC Client

A gRPC client with etcd-based service discovery.

```go
import (
	grpcc "github.com/abulo/ratel/v3/client/grpc"
	regEtcd "github.com/abulo/ratel/v3/registry/etcdv3"
)

// 1. etcd service discovery config (note: use the registry/etcdv3 Config;
// chaining Set methods returns the embedded client config type, so call them
// separately to preserve the type)
etcdCfg := regEtcd.New()
etcdCfg.SetEndpoints([]string{"127.0.0.1:2379"})

// 2. gRPC client config
conn := grpcc.New().
	SetName("order-service").               // client instance name (singleton key and tag)
	SetAddress("order.svc").                // target service name (resolved by the etcd resolver)
	SetBalancerName("p2c").                 // load balancing: roundrobin / swr / p2c
	SetEtcd(etcdCfg).                       // required: etcd service discovery
	SetReadTimeout(1 * time.Second).        // per-call timeout
	SetSlowThreshold(600 * time.Millisecond).
	SetEnableTraceInterceptor(true).
	SetEnableMetricInterceptor(true).
	SetEnableAccessInterceptor(true).
	SetEnableTimeoutInterceptor(true).
	MustBuild() // or Build() which returns (conn, error)
```

**Configuration:**

| Field | Default | Description |
|---|---|---|
| `BalancerName` | `roundrobin` | Load balancing strategy (`roundrobin`/`swr`/`p2c`) |
| `DialTimeout` | 3s | Dial timeout |
| `ReadTimeout` | 1s | Per-call timeout (used by the timeout interceptor) |
| `SlowThreshold` | 600ms | Slow-call threshold |
| `EnableDebug` | false | Debug interceptor (prints address/method/response/error) |
| `EnableTraceInterceptor` | false | Tracing interceptor |
| `EnableAidInterceptor` | false | Application ID metadata propagation |
| `EnableTimeoutInterceptor` | false | Call timeout and slow-call statistics |
| `EnableMetricInterceptor` | false | Metrics collection |
| `EnableAccessInterceptor` | false | Access log interceptor |
| `Etcd` | required | etcd service discovery config; missing it causes a nil pointer |

**Load balancing:** `client/grpc/balancer` provides SWR (smooth weighted round-robin) and P2C (power of two choices) balancers; the `p2c` sub-package picks the node with the lower in-flight request count.

### rabbitmq

A RabbitMQ client wrapper based on `streadway/amqp`, with automatic reconnection on disconnect.

```go
import (
	"github.com/abulo/ratel/v3/client/rabbitmq"
	"github.com/streadway/amqp"
)

// 1. Establish a connection (auto-reconnects on disconnect)
conn, err := rabbitmq.Dial("amqp://guest:guest@127.0.0.1:5672/", rabbitmq.DefaultTimesRetry())
if err != nil {
	panic(err)
}

// 2. Declare and bind a queue (queue name, routing key, exchange)
queue := rabbitmq.NewQueueBuilder(conn).
	SetQueueDeclareOpts(func(builder *rabbitmq.QueueDeclareOptsBuilder) *rabbitmq.QueueDeclareOpts {
		return builder.SetDurable(true).SetAutoDelete(false).Build()
	}).
	Build()
if err := queue.DeclareAndBind("task.queue", "task.key", "exchange"); err != nil {
	panic(err)
}

// 3. Get a channel and send a message (channels are managed by the connection and reused)
ch, err := conn.Channel()
if err != nil {
	panic(err)
}
if err := ch.SendOpts("exchange", "task.key", []byte("hello"), rabbitmq.DefaultSendOpts()); err != nil {
	// handle send failure
}

// 4. Consume messages (ConsumerFunc returns true to stop consuming)
// ch.ReceiveOpts("task.queue", func(delivery *amqp.Delivery) bool {
// 	// handle the message; return true to acknowledge and exit
// 	return true
// }, rabbitmq.DefaultReceiveOpts())
```

**Key types:**

| Type | Description |
|---|---|
| `Connection` | Auto-reconnecting connection, built via `Dial`/`NewConnection` |
| `Channel` | Channel reusing the connection, obtained via `Connection.Channel()`; provides `Send`/`SendOpts`/`Receive`/`ReceiveOpts` |
| `QueueBuilder` | Queue declare/bind builder (`SetQueueDeclareOpts`/`SetQueueBindOpts`/`SetRetryable`/`Build`) |
| `Queue` | Queue carrier; `DeclareAndBind`/`RetryDeclareAndBind` declare and bind |
| `Producer` / `Consumer` | Convenient producer (`Send`) and consumer (`Receive`/`Get`) wrappers |
| `SendOptsBuilder` / `ReceiveOptsBuilder` / `QueueDeclareOptsBuilder` / `QueueBindOptsBuilder` | Options builders |
| `TimesRetry` / `CtxRetry` | Retry strategies |
| `MessagePlainTransient` etc. | Common message factories |

### ftp / sftp

File transfer clients:

```go
// FTP: config fields + Config.New dial + Login authentication
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

// SFTP: config fields + Config.New (underlying SSH connection + SFTP channel)
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

Both clients provide directory switching, file/directory CRUD, upload/download (SFTP supports recursive directories and transfer progress statistics), and more.

---

## Storage Packages

### redis

A Redis client wrapper based on go-redis, supporting four connection modes.

```go
import "github.com/abulo/ratel/v3/stores/redis"

// Standalone
client, err := redis.NewRedisClient(
	redis.WithClientType("normal"),
	redis.WithAddr("127.0.0.1:6379"),
	redis.WithPassword(""),
	redis.WithDatabase(0),
	redis.WithPoolSize(20),
)

// Cluster
client, err = redis.NewRedisClient(
	redis.WithClientType("cluster"),
	redis.WithHosts([]string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}),
)

// Sentinel
client, err = redis.NewRedisClient(
	redis.WithClientType("failover"),
	redis.WithMasterName("mymaster"),
	redis.WithHosts([]string{"127.0.0.1:26379", "127.0.0.1:26380"}),
)

// Ring (sharding)
client, err = redis.NewRedisClient(
	redis.WithClientType("ring"),
	redis.WithAddrs(map[string]string{
		"shard1": "127.0.0.1:6379",
		"shard2": "127.0.0.1:6380",
	}),
)
```

**Configuration (Option):**

| Option | Description |
|---|---|
| `WithClientType` | `normal` (standalone) / `cluster` / `failover` / `ring` |
| `WithAddr` | Standalone address `host:port` |
| `WithHosts` | Cluster/sentinel node list |
| `WithAddrs` | Shard name to address mapping |
| `WithMasterName` | Sentinel master name |
| `WithDatabase` | Database number |
| `WithPassword` | Connection password |
| `WithDialTimeout` / `WithReadTimeout` / `WithWriteTimeout` | Timeout configuration |
| `WithPoolSize` / `WithPoolTimeout` / `WithIdleTimeout` / `WithIdleCheckFrequency` / `WithConnMaxLifetime` | Connection pool configuration |
| `WithEnableMetric` / `WithEnableTrace` | Metrics and tracing hooks |
| `WithName` | Instance name (singleton key) |

**Command calls:** commands are organized into grouped methods by data type (strings, hashes, lists, sets, sorted sets, bitmaps, streams, geo, JSON, search, scripts, Pub/Sub, etc.), all executed through the circuit breaker; acceptable errors such as missing keys are not reported as failures. The underlying go-redis client can be obtained via `redis.OriginalClient(client)`.

### sql

A GORM-based relational database client wrapper supporting MySQL, PostgreSQL, and ClickHouse.

```go
import "github.com/abulo/ratel/v3/stores/sql"

// 1. Configure and open the connection: NewClient only assembles the config;
// Open establishes the connection (returns a ClientManager that embeds *gorm.DB)
client, err := sql.NewClient(
	sql.WithDriverName(sql.MySQLDriverName), // or sql.PostgreSQLDriverName / sql.ClickHouseDriverName
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

// 2. Run GORM queries (manager embeds *gorm.DB; native GORM APIs are available)
// manager.Model(&User{}).Find(&users)
```

**Configuration (Option):**

| Option | Description |
|---|---|
| `WithDriverName` | `mysql` / `postgres` / `clickhouse` |
| `WithHost` / `WithPort` / `WithUsername` / `WithPassword` / `WithDatabase` | Connection info |
| `WithCharset` / `WithParseTime` / `WithTimeZone` / `WithSslMode` | Driver-specific DSN parameters |
| `WithDialTimeOut` / `WithReadTimeOut` | ClickHouse timeouts |
| `WithMaxOpenConns` / `WithMaxIdleConns` / `WithMaxIdleTime` / `WithMaxLifetime` | Connection pool |
| `WithEnableDebug` / `WithEnableMetric` / `WithEnableTrace` | Switches |

**Query helpers:**

```go
// Pagination parameter carrier (fields are pointers; getting a value from an unset field returns an error)
offset := int64(0)
limit := int64(20)
page := &sql.Pagination{Offset: &offset, Limit: &limit}

// Range and sort parameter carriers
scope := &sql.Scope{Start: ptr("2024-01-01"), End: ptr("2024-12-31")}
sort := &sql.Sort{Name: ptr("created_at"), Desc: ptr(true)}

// Get values and plug them into a GORM chain (manager embeds *gorm.DB)
offsetVal, _ := page.GetOffset()
limitVal, _ := page.GetLimit()
sortName, _ := sort.GetName()
manager.WithContext(ctx).Model(&User{}).
	Offset(int(offsetVal)).
	Limit(int(limitVal)).
	Order(sortName).
	Find(&users)
```

> `ptr` is an address helper (`func ptr[T any](v T) *T { return &v }`); define it yourself as needed.

**Interceptors:** `sql.MetricInterceptor()` and `sql.TraceInterceptor()` can be registered via `sql.RegisterInterceptor(manager, client, interceptors...)` to automatically collect SQL metrics and tracing (`Client.Open` already wires them automatically when EnableMetric/EnableTrace are on).

### mongodb

A MongoDB client wrapper based on mongo-driver.

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

// Get an enhanced collection handle (collection is the collection name; the database comes from the client Name)
model, err := mongoClient.NewModel("users")
if err != nil {
	panic(err)
}
// model.Collection is a *mongodb.DecoratedCollection wrapping circuit breaker, metrics, and tracing
```

**DecoratedCollection capabilities:**

- Chainable condition construction (wraps native filter building).
- Fault protection: operations execute through the circuit breaker.
- Metrics collection and tracing are attached automatically.
- `BeforeCreate`/`BeforeUpdate` convert structs to BSON documents and generate `_id` before writing.

### elasticsearch

An Elasticsearch client wrapper based on olivere/elastic v7.

```go
import "github.com/abulo/ratel/v3/stores/elasticsearch"

esClient := elasticsearch.NewClient(&elasticsearch.Config{
	URL:          []string{"http://127.0.0.1:9200"},
	Username:     "elastic",
	Password:     "secret",
	EnableMetric: true,
	EnableTrace:  true,
})

// esClient embeds the native elastic.Client; all native APIs are available
// esClient.Search("index").Query(...).Do(ctx)
```

Clients are cached and reused by node address; `ESTracedTransport` attaches tracing and metrics collection at the HTTP transport layer, and `MaxContentLength` (64KB) limits the request body size recorded in traces.

### null Nullable Types

A complete set of nullable base types, replacing `sql.NullXxx` for handling null values in databases and JSON.

**Implemented types:** `Bool`, `Byte`, `Bytes`, `String`, `Int`, `Int8/16/32/64`, `Uint`, `Uint8/16/32/64`, `Float32/64`, `Time`, `Date`, `DateTime`, `Timestamp`, `CTime`, `JSON`.

Each type consists of three state fields:

- **Value**: the actual value.
- **`Valid`**: whether a non-null valid value is held.
- **`Set`**: whether it was explicitly assigned (including assigning null).

All types uniformly implement `sql.Scanner`, `driver.Valuer`, `json.Marshaler`, `json.Unmarshaler`, `encoding.TextMarshaler`, and `encoding.TextUnmarshaler`; null values serialize to `null` in JSON and map to SQL `NULL` when reading/writing the database.

```go
import "github.com/abulo/ratel/v3/stores/null"

type User struct {
	ID    int64     `json:"id"`
	Name  null.String `json:"name"`
	Email null.String `json:"email"`
}

// JSON null handling
b, _ := json.Marshal(User{ID: 1, Name: null.StringFrom("ratel"), Email: null.String{}})
// {"id":1,"name":"ratel","email":null}
```

### proxy Proxy Pool

A proxy pool that registers and retrieves various data storage backends by group.

```go
import "github.com/abulo/ratel/v3/stores/proxy"

p := proxy.NewProxy()

// Initialization: construct proxy carriers, inject underlying clients, then register into the pool
redisProxy := proxy.NewRedis()
redisProxy.Store(redisClient)
p.StoreRedis("default", redisProxy)

mongoProxy := proxy.NewMongoDB()
mongoProxy.Store(mongoClient)
p.StoreMongoDB("default", mongoProxy)

esProxy := proxy.NewElasticSearch()
esProxy.Store(esClient)
p.StoreElasticSearch("default", esProxy)

// The SQL proxy supports read/write separation: SetWrite/SetRead take *gorm.DB, Write/Read pick a DB at random
sqlProxy := proxy.NewSQL()
sqlProxy.SetWrite(sqlClient)       // write DB (required)
sqlProxy.SetRead(sqlReadReplica)   // read DB (optional; falls back to the write DB)
p.StoreSQL("write", sqlProxy)

// At runtime, fetch by group and use
r := p.LoadRedis("default")   // *redis.Client
m := p.LoadMongoDB("default") // *mongodb.MongoDB
e := p.LoadElasticSearch("default") // *elasticsearch.Client
s := p.LoadSQL("write")       // *proxy.SQL
db := s.Write()               // random write DB; use s.Read() for multi-read-DB scenarios
```

The SQL proxy additionally supports write/read DB lists and random selection (read/write separation).

### session Storage

Redis-based session storage; a whole session is stored as JSON under a single Redis key.

```go
import "github.com/abulo/ratel/v3/stores/session"

sess := &session.Session{
	Driver: redisClient, // underlying Redis client; must not be nil
	Name:   "session:user:1001",
	TTL:    3600, // seconds; 0 means no expiry
}

// Read/write nested fields via dotted paths
_ = sess.Put(ctx, "profile.name", "ratel")
_ = sess.Put(ctx, "profile.level", 3)
name := sess.Get(ctx, "profile.name")   // "ratel"

// Remove a field / destroy the whole session
sess.Remove(ctx, "profile.level")
sess.Destroy(ctx)
```

---

## Governance Capabilities

### Rate Limiting ratelimit

A token-bucket rate limiter implementation.

```go
import "github.com/abulo/ratel/v3/core/ratelimit"

// Fixed interval refill: 1 token every 100ms, capacity 100
bucket := ratelimit.NewBucket(100*time.Millisecond, 100)

// Fixed rate refill: rate tokens per second
bucket = ratelimit.NewBucketWithRate(1000, 100)

// Per-tick quantum: 5 tokens every 100ms
bucket = ratelimit.NewBucketWithQuantum(100*time.Millisecond, 100, 5)
```

**Bucket methods:**

| Method | Behavior |
|---|---|
| `Take(count)` | Takes tokens, returns the wait duration; waits as needed when the bucket is empty |
| `TakeMaxDuration(count, maxWait)` | Timed wait for tokens; returns `(0, false)` on timeout |
| `TakeAvailable(count)` | Non-blocking token take; returns the number actually taken |
| `Wait(count)` | Blocks until tokens are available (no return value) |
| `WaitMaxDuration(count, maxWait)` | Timed wait; returns whether tokens were acquired in time (bool) |
| `Available()` / `Capacity()` / `Rate()` | Query current available tokens, capacity, and rate |

**Rate-limited Reader/Writer:**

```go
limitedReader := ratelimit.Reader(r, bucket)   // limit read rate
limitedWriter := ratelimit.Writer(w, bucket)   // limit write rate
```

The bucket is internally locked and safe for concurrent use by multiple goroutines.

### Circuit Breaker resource

An adaptive circuit breaker based on the Google SRE overload protection model.

```go
import "github.com/abulo/ratel/v3/core/resource"

// 1. Get or create a named breaker
breaker := resource.GetBreaker("my-service")

// 2. Wrap a business call
err := breaker.Do(func() error {
	return callDownstream()
})

// 3. Custom error judgment: only the specified errors count as failures
err = breaker.DoWithAcceptable(func() error {
	return callDownstream()
}, func(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
})

// 4. Fallback callback: executed when the breaker is open
err = breaker.DoWithFallback(func() error {
	return callDownstream()
}, func(err error) error {
	return fallbackResult()
})

// 5. Manually manage the request lifecycle
promise, err := breaker.Allow()
if err != nil {
	return resource.ErrServiceUnavailable
}
defer promise.Reject("failed") // or promise.Accept() on success; call exactly one of Accept/Reject
// business logic
```

**Breaker interface:**

| Method | Description |
|---|---|
| `Allow() (Promise, error)` | Returns a Promise when allowed; call `Accept` on success, `Reject` on failure |
| `Do(req)` | Runs `req` when allowed, returns an error when rejected; panics count as failures and are re-thrown |
| `DoWithAcceptable(req, acceptable)` | Customize whether an error counts as success |
| `DoWithFallback(req, fallback)` | Runs a fallback callback when rejected |
| `DoWithFallbackAcceptable(...)` | Combines fallback and custom judgment |

**Other resource components:**

```go
// SingleFlight request coalescing
sf := resource.NewSingleFlight()
v, err := sf.Do("key", func() (any, error) { return loadFromDB() })

// Batch error aggregation
batch := &resource.BatchError{}
batch.Add(err1)
batch.Add(err2)
if err := batch.Err(); err != nil { ... }

// Resource manager (manages closeable resources by name)
manager := resource.NewResourceManager()
```

### goroutine Orchestration

`core/goroutine` provides goroutine orchestration and safe execution utilities.

| Function | Behavior |
|---|---|
| `Go(fn)` | Safely starts a goroutine (automatically recovers panics) |
| `DelayGo(delay, fn)` | Delayed start |
| `SafeGo(fn, rec)` | Starts and calls a recovery function to handle panics |
| `GoDirect(fn, args...)` | Reflectively invokes any function |
| `Serial(fns...)` / `SerialWithError(fns...)` | Serial execution |
| `Parallel(fns...)` | Parallel execution; returns a wait function |
| `ParallelWithError(fns...)` / `ParallelWithErrorChan(fns...)` | Parallel execution with error collection |
| `RestrictParallel(restrict, fns...)` | Rate-limited concurrency |
| `RestrictParallelWithErrorChan(concurrency, fns...)` | Rate-limited concurrency with error collection |

```go
import "github.com/abulo/ratel/v3/core/goroutine"

goroutine.Go(func() {
	// safe execution; a panic will not crash the process
})

// Parallel execution and wait
wait := goroutine.Parallel(
	func() { processA() },
	func() { processB() },
	func() { processC() },
)
wait()
```

### Lifecycle cycle & hooks

```go
// cycle: concurrent lifecycle controller
import "github.com/abulo/ratel/v3/core/cycle"

cyc := cycle.NewCycle()
cyc.Run(func() error { /* background task */ })
cyc.Run(func() error { /* another background task */ })
// Wait for all tasks to finish and collect errors
if err := <-cyc.Wait(); err != nil { ... }
```

```go
// hooks: staged hooks (reverse order, defer semantics)
import "github.com/abulo/ratel/v3/core/hooks"

hooks.Register(hooks.Stage_BeforeRun, func() { /* ... */ })
hooks.Register(hooks.Stage_AfterStop, func() { /* ... */ })
hooks.Do(hooks.Stage_BeforeRun)
```

---

## Observability

### Logging logger

A process-wide global logger based on logrus.

```go
import "github.com/abulo/ratel/v3/core/logger"

logger.Logger.Info("application started")
logger.Logger.WithFields(logrus.Fields{
	"user_id": 123,
	"action":  "login",
}).Error("operation failed")
```

**Multi-backend async output (hooks):**

```go
import (
	"github.com/abulo/ratel/v3/core/logger"
	eslogger "github.com/abulo/ratel/v3/core/logger/es"
	mongologger "github.com/abulo/ratel/v3/core/logger/mongo"
	sqllogger "github.com/abulo/ratel/v3/core/logger/sql"
)

// Elasticsearch hook (panics on construction failure)
esHook := eslogger.New(eslogger.WithClient(esClient))
logger.Logger.AddHook(esHook)

// MongoDB hook
mongoHook := mongologger.New(mongologger.WithClient(mongoClient))
logger.Logger.AddHook(mongoHook)

// SQL hook
sqlHook := sqllogger.New(sqllogger.WithDB(gormDB))
logger.Logger.AddHook(sqlHook)
```

Log entries are written in the background by the async queue in `core/logger/queue`, without blocking business calls; `core/logger/entry` defines the standardized log entry structure (hostname, timestamp, caller location, level, extra fields).

### Metrics metric

A Prometheus-based metrics collection wrapper; the namespace defaults to `ratel` and is automatically registered with the default registry at construction.

```go
import "github.com/abulo/ratel/v3/core/metric"

// Pre-built vectors (labels in definition order: type, method, client, code)
metric.ServerHandleCounter.Inc("http", "GET", "my-service", "200")

// Custom vector
counter := metric.NewCounterVec("my_metric_total", []string{"label"})
counter.Inc("value")
```

**Pre-built metrics:**

| Metric | Description |
|---|---|
| `ServerHandleCounter` / `ServerHandleSeconds` | Server request count / duration |
| `ClientHandleCounter` / `ClientHandleSeconds` | Client call count / duration |
| `JobHandleCounter` / `JobHandleSeconds` | Job execution count / duration |

### Tracing trace

OpenTelemetry tracing wrapper.

```go
import (
	"github.com/abulo/ratel/v3/core/trace"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Initialize the global TracerProvider
cfg := trace.NewConfig(
	trace.WithName("my-service"),
	trace.WithEndpoint("127.0.0.1:4317"), // OTLP endpoint
	trace.WithProtocol("grpc"),           // grpc / http
	trace.WithSampler(1.0),               // sampling rate 0-1
)
trace.SetGlobalTracer(cfg.Build())

// Create a span in business code (the carrier is used to inject/extract propagation context)
tracer := trace.NewTracer(oteltrace.SpanKindServer)
ctx, span := tracer.Start(ctx, "operation", propagation.MapCarrier{})
defer span.End()
```

Supports HTTP/gRPC OTLP export protocols, compatible with Jaeger and the W3C TraceContext propagation format; provides Tracer wrappers for server, client, producer, and consumer roles. Middleware injection is automatic: xgin/xhertz server middleware, xgrpc interceptors, and the storage/client hooks all carry tracing context automatically.

---

## Task Scheduling

### crontab Scheduler

A scheduled task scheduler based on cron expressions.

```go
import "github.com/abulo/ratel/v3/crontab"

// 1. Create a task manager
manager := crontab.NewTaskManager()

// 2. Register a task (five/six-field cron or descriptors such as @yearly/@monthly)
manager.AddTask("hourly-job", crontab.NewTask("hourly-job", "0 * * * *", func(ctx context.Context) error {
	// task logic
	return nil
}, crontab.TimeoutOption(30*time.Second))) // optional timeout control

// 3. Start scheduling
manager.StartTask()

// 4. Graceful shutdown
<-manager.GracefulShutdown()
```

**TaskManager methods:**

| Method | Description |
|---|---|
| `AddTask(name, task)` | Registers a task |
| `DeleteTask(name)` / `ClearTask()` | Deletes tasks |
| `StartTask()` / `StopTask()` | Starts/stops scheduling |
| `GracefulShutdown()` | Waits for all tasks to finish |
| `WorkerStart()` / `WorkerStop()` | Implements the worker.Worker interface; can be hosted by the application |

### Distributed Tasks core/task

Distributed scheduled tasks based on a cron scheduler and a consistent-hash node pool; tasks are distributed within the cluster to avoid duplicate execution.

```go
import (
	"github.com/abulo/ratel/v3/client/etcdv3"
	"github.com/abulo/ratel/v3/core/task"
	"github.com/abulo/ratel/v3/core/task/driver"
)

// 1. Create a driver (Redis plain key / Redis sorted set / etcd lease)
etcdClient, _ := etcdv3.NewClient(etcdv3.New().SetEndpoints([]string{"127.0.0.1:2379"}))
drv := driver.NewEtcdDriver(etcdClient)
// or drv := driver.NewRedisDriver(redisClient)
// or drv := driver.NewRedisZSetDriver(redisClient)

// 2. Build the distributed task manager (options such as CronOptionSeconds are passed via NewTaskWithOption)
t := task.NewTaskWithOption("order-service", drv, task.CronOptionSeconds())

// 3. Register tasks (AddFunc/AddJob; cron expressions are parsed with 6 fields including seconds)
t.AddFunc("sync-order", "0 */5 * * * *", func() {
	// task logic; executed by only one node in the cluster
})

// 4. Start
t.Start()
```

**Core concepts:**

- **NodePool**: distributes tasks among cluster members via consistent hashing (configurable replica count `WithHashReplicas`); on member changes it automatically enters an upgrading state (`NodePoolStateUpgrading`) to avoid duplicate execution.
- **Task recovery**: supports recovery (`WithRecoverFunc`) and recent-job replay (`NewRecentJobPacker`) on start.
- **driver**: three node registration/discovery implementations — `NewRedisDriver`/`NewRedisZSetDriver`/`NewEtcdDriver`.

### Distributed Executor executor

An xxl-job style task executor interface with global management.

```go
import "github.com/abulo/ratel/v3/core/executor"

// 1. Implement the executor interface (GetAddress/RegXJob/Run/Stop/GracefulStop)
type MyExecutor struct {
	address string
}

func (e *MyExecutor) GetAddress() string { return e.address }
func (e *MyExecutor) RegXJob(jobs ...executor.XJob) {}
func (e *MyExecutor) Run() error                      { return nil }
func (e *MyExecutor) Stop()                           {}
func (e *MyExecutor) GracefulStop()                   {}

// 2. Register the executor and start
executor.Register(":9999", &MyExecutor{address: ":9999"})
if err := executor.Run(); err != nil { ... }

// 3. Stop
// executor.Stop() or executor.GracefulStop()
```

### worker & job Interfaces

```go
// worker.Worker: background worker unit (hosted by the application)
type MyWorker struct{}

func (w *MyWorker) WorkerStart() error { /* start */ }
func (w *MyWorker) WorkerStop() error  { /* stop */ }

application, _ := app.New()
_ = application.Schedule(&MyWorker{}) // register the background worker

// job.Runner: a task executed concurrently after startup (Run has no return value; usually a resident loop)
type MyJob struct{}

func (j *MyJob) Run() { /* task logic */ }

_ = application.Job(&MyJob{}) // register the task
```

---

## Utility Packages

### util

A collection of general utility functions (pure functions, safe for concurrent calls).

```go
import "github.com/abulo/ratel/v3/util"

// Strings and slices
util.InArray("x", []string{"a", "b", "x"})      // true
util.ArrayRemoveRepeatedElement([]string{"a", "a", "b"}) // ["a","b"]
util.Implode(",", []string{"a", "b"})           // "a,b"
util.Explode(",", "a,b")                        // ["a","b"]

// JSON
util.JSONString(obj)              // serialize (with unboxing optimization)
util.JSONDecode([]byte(s), &dst)  // deserialize

// Time
util.Timestamp()                  // current timestamp
util.FormatDateTime(time.Now())   // format time
util.DateTimeParse(s)             // parse time

// Math
util.Max(1, 2) / util.Min(1, 2)   // max / min (variadic float64)
util.RandomInt(0, 100)            // random int ((int, error))
bytes, _ := util.RandomBytes(32)  // random bytes (([]byte, error))

// Files and network
util.FileExists(path)             // whether the file exists
util.GetHostByName("example.com") // DNS resolution

// Others
util.ZhCharToPinyin("中国")        // pinyin conversion
util.Uniqid("prefix_")            // unique ID
util.Md5(s) / util.Sha1(s)        // digests
util.MapToStruct(map[string]any)  // struct conversion (returns *pb.Struct)
util.StructToMap(pbStruct)        // struct conversion (returns map[string]any)
```

Covers strings, slices and maps, time and date, math, filesystem, network, URL, regex info extraction, random numbers, pinyin, and protobuf struct conversion.

### snowflake

A Snowflake-algorithm int64 unique ID generator.

```go
import "github.com/abulo/ratel/v3/snowflake"

// Typical application config: 40-bit timestamp + 7-bit sequence + 16-bit machine
config := snowflake.CommonConfig

// Or custom bit layout and machine ID
config = snowflake.ConstructConfigWithMachineID(40, 7, 16, 1)

// Generate an ID (atomic, concurrency-safe)
id := config.GenInt64ID()
```

### correct

Automatic typography and cleanup for mixed CJK/Latin text.

```go
import "github.com/abulo/ratel/v3/correct"

// Typography: auto-inserts spaces between CJK and English/digits/symbols, normalizes full-width characters
out := correct.Format("这是一段English文本，包含123数字。")
// 这是一段 English 文本，包含 123 数字。

// Typography for HTML (only text nodes between tags are processed; tags and code blocks are preserved)
out, err := correct.FormatHTML(htmlBody)

// Reverse typography (removes the auto-inserted spaces)
orig := correct.Unformat(out)
```

### nlpword / filter

Aho-Corasick based sensitive-word processing utilities (two implementations).

```go
// nlpword: detection, filtering, replacement
import "github.com/abulo/ratel/v3/nlpword"

filter := nlpword.New()
_ = filter.LoadWordDict("words.txt")      // load from file
// filter.LoadNetWordDict(url)             // load from network
// filter.AddWord("敏感词")                // add in-program

replaced := filter.Replace("这句话包含敏感词", '*') // replace matched words
ok, hit := filter.Validate("测试文本")             // whether it contains sensitive words (returns the hit)
// or filter.FindIn(text) (bool, string) / filter.FindAll(text) []string

// filter: scanner (built from word lists/files, with custom skip separators)
import "github.com/abulo/ratel/v3/filter"

search := filter.Strings([]string{"敏感词1", "敏感词2"}, "@#$%^&*!?，。；：、")
has := search.HasSens([]byte("文本内容"))           // whether it contains sensitive words
hits := search.Find([]byte("文本内容"))             // collect all hits
replaced = string(search.ReplaceRune([]byte("文本内容"), '*'))
```

### i18n

Multilingual message management.

```go
import "github.com/abulo/ratel/v3/lang/i18n"

// Initialize the default instance: language dir + default language + language map (key=language name, value=display name)
i18n.Init("./locales", "zh-CN", map[string]string{
	"zh-CN": "简体中文",
	"en-US": "English",
})
// In default FileMode, language files are loaded from langDir/<language>.ini (DirMode loads per-language directories)

// Translate for a specified language (missing keys fall back to the backup language; supports argument rendering)
msg := i18n.Tr("zh-CN", "hello", "Ratel") // "你好, Ratel"
// or use the default language: i18n.Dt("hello", "Ratel")
// you can also build your own instance: instance := i18n.NewWithInit(langDir, defLang, languages)
```

### img

Image reading, transformation, composition, and saving (based on libvips; requires the dependency installed).

```go
import (
	"github.com/abulo/ratel/v3/img"
	"github.com/disintegration/imaging"
)

image, err := img.OpenLocalFile("input.jpg")
if err != nil {
	panic(err)
}

// Chainable transformations (mutate the source image in place; methods return the receiver):
// proportional resize → center crop → gaussian blur
image.Resize(800, 0, imaging.Lanczos)
image.Crop(800, 600, imaging.Center)
image.Blur(5)

// Overlay a watermark
wm, _ := img.OpenLocalFile("logo.png")
image.AddWaterMark(wm, imaging.Center, 20, 20, 0.5)

// Save (format inferred from the file extension; JPEG quality controlled via imaging.JPEGQuality)
if err := img.SaveToFile(image, "output.jpg", imaging.JPEGQuality(80)); err != nil {
	panic(err)
}
```

Supports resize, crop, blur, grayscale, contrast/brightness/saturation adjustment, rotation, paste, overlay, watermark, JPEG compression, and text drawing via `img/fontx`; methods return the receiver for chainable calls.

### watch Hot Compilation

Go project hot compilation watcher (auto-build and restart on file changes).

```go
import "github.com/abulo/ratel/v3/watch"

err := watch.Watch(ctx, &watch.Options{
	// configure watch dirs, entry file, output name, file extensions, and build arguments
})
```

### config Configuration

A Viper-based configuration reading wrapper.

```go
import "github.com/abulo/ratel/v3/config"

conf := config.New()
_ = conf.LoadDir("./configs")     // batch load a directory
// _ = conf.LoadFile("config.yaml") // or load a single file

// Typed reads (return the default when missing)
port := conf.Int("server.port", 8080)
name := conf.String("server.name", "ratel")
timeout := conf.Duration("server.timeout", 5*time.Second)
```

Supports loading JSON, TOML, and YAML configs from a single file or a directory in batch.

---

## CLI Toolkit

`toolkit` is ratel's command-line scaffolding program (based on Cobra) for project creation and code generation.

```bash
# Install
go install github.com/abulo/ratel/v3/toolkit@latest

# Check version and help
toolkit version
toolkit --help
```

**Subcommands:**

| Command | Function |
|---|---|
| `toolkit project` | Project scaffolding: clones a template from the remote layout repo, generates backend/frontend projects, or appends service directories to an existing project |
| `toolkit module` | Data model layer code generation: interactively collects model path, table names, and soft-delete strategy; generates module/proto/service/convert code |
| `toolkit dao` | DAO code generation: reads table metadata and generates GORM model structs (single/multi-table) |
| `toolkit api` | HTTP API layer code generation: generates handler functions and route registration snippets per web framework (gin/hertz) |
| `toolkit vue` | Vue3 + TypeScript frontend code generation: interface definitions, API request modules, list pages |
| `toolkit build` | Hot compilation: watches directory file changes and automatically rebuilds and restarts |
| `toolkit changelog` | Changelog generation: categorizes conventional-commit types and writes to CHANGELOG.md |
| `toolkit upgrade` / `toolkit init` | Scaffold upgrade and initialization: installs the latest toolkit binary and protoc tools |

```bash
# Example: interactively create a project
toolkit project

# Example: generate DAOs (specified tables)
toolkit dao --table users,orders
```

---

## Project Structure

```
ratel/
├── client/              # client wrappers
│   ├── etcdv3/         # etcd client (KV/lock/lease/watch)
│   ├── grpc/           # gRPC client (service discovery/load balancing/interceptors)
│   ├── grpc/balancer/  # SWR, P2C balancers
│   ├── rabbitmq/       # RabbitMQ client (auto-reconnect)
│   ├── ftp/            # FTP client
│   └── sftp/           # SFTP client
├── config/             # Viper config reading
├── core/               # core governance capabilities
│   ├── app/            # application lifecycle management
│   ├── logger/         # logging (es/mongo/sql backends + queue async queue)
│   ├── metric/         # Prometheus metrics
│   ├── ratelimit/      # token-bucket rate limiting
│   ├── resource/       # circuit breaker, SingleFlight, batch errors, resource management
│   ├── task/           # distributed tasks and cron scheduling
│   ├── executor/       # xxl-job style distributed executor
│   ├── trace/          # OpenTelemetry tracing
│   ├── goroutine/      # goroutine orchestration
│   ├── cycle/          # concurrent lifecycle control
│   ├── hooks/          # lifecycle stage hooks
│   ├── ecode/          # unified error code model
│   ├── singleton/      # singleton storage
│   ├── constant/       # global constants
│   ├── env/            # runtime environment metadata
│   ├── hostname/       # hostname retrieval
│   ├── call/           # call stack info
│   ├── component/      # component abstraction
│   └── stack/          # call stack utilities
├── server/             # server implementations
│   ├── xgin/           # Gin HTTP server
│   ├── xgrpc/          # gRPC server
│   ├── xhertz/         # Hertz HTTP server
│   ├── xasynq/         # Asynq scheduled tasks (manager/worker)
│   └── xmonitor/       # monitor server (pprof, etc.)
├── stores/             # data storage
│   ├── redis/          # Redis (standalone/cluster/sentinel/ring)
│   ├── sql/            # SQL (MySQL/PostgreSQL/ClickHouse)
│   ├── mongodb/        # MongoDB
│   ├── elasticsearch/  # Elasticsearch
│   ├── null/           # nullable type set
│   ├── proxy/          # storage proxy pool
│   └── session/        # Redis session storage
├── registry/           # service registration and discovery (etcdv3 implementation)
├── toolkit/            # command-line scaffolding
├── util/               # general utilities
├── correct/            # CJK/Latin text typography
├── nlpword/            # sensitive-word processing (AC automaton)
├── filter/             # sensitive-word scanning
├── img/                # image processing (libvips)
├── lang/               # i18n and placeholder types
├── snowflake/          # Snowflake IDs
├── crontab/            # cron scheduler
├── watch/              # hot compilation watcher
└── worker/             # background worker units and tasks
```

---

## Development Guide

### Local Development

```bash
# Clone the project
git clone https://github.com/abulo/ratel.git
cd ratel

# Install dependencies
go mod download

# Run unit tests
make test

# Static checks
make vet       # go vet
make errcheck  # error check
make lint      # revive lint

# Formatting and format checks
make fmt       # gofmt -w
make fmtcheck  # gofmt check

# Full build (fmt + errcheck + lint + build)
make all
```

### Makefile targets

| Target | Description |
|---|---|
| `make all` | Full build (fmt → errcheck → lint → build) |
| `make test` | Run unit tests (`go test -short`) |
| `make vet` | Static checks (`go vet ./...`) |
| `make lint` | revive checks (auto-installs revive) |
| `make fmt` / `make fmtcheck` | Format / format check |
| `make errcheck` | Error-handling check (auto-installs errcheck) |
| `make cmt` | Auto-comment exported functions (gocmt) |

### Code conventions

- Use `gofmt` to format code and `revive` for lint checks.
- Write Godoc comments for exported identifiers (package comments are consolidated into `doc.go`, starting with `Package xxx`).
- Check all error returns; never swallow exceptions.
- Write unit tests covering critical paths.
- Follow conventional commits for commit messages: `feat`/`fix`/`docs`/`perf`/`chore`/`refactor`/`test`, etc.

### Continuous Integration

The repository is configured with GitHub Actions (`.github/workflows/go.yml`) that builds and reports to Codecov on push and PR; `release.yml` handles the release flow.

---

## FAQ

### Q1: What is the difference between Ratel and native Gin/gRPC?

Ratel provides unified wrappers on top of the native frameworks: Server lifecycle management, automatic middleware (metrics/tracing/slow query), registry integration, governance capabilities (rate limiting/circuit breaking), and companion tools — business code does not need to care about this infrastructure.

### Q2: How do I migrate an existing project to Ratel?

Ratel servers embed the native engines, so all native Gin/gRPC/Hertz APIs remain available:

```go
// Just replace the import; the rest of the code needs no changes
// import "github.com/gin-gonic/gin"
import "github.com/abulo/ratel/v3/server/xgin"
```

### Q3: How do I customize middleware?

```go
// xgin (embeds gin.Engine)
server.Use(func(c *gin.Context) {
	// custom logic
	c.Next()
})

// xgrpc
cfg.WithUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// custom logic
	return handler(ctx, req)
})
```

### Q4: How do I enable rate limiting?

```go
// Use the token bucket directly in business code
bucket := ratelimit.NewBucketWithRate(1000, 100)
if bucket.TakeAvailable(1) == 0 {
	return errors.New("rate limit exceeded")
}
```

### Q5: How do I configure service registration?

```go
import (
	"github.com/abulo/ratel/v3/core/app"
	"github.com/abulo/ratel/v3/registry"
	regEtcd "github.com/abulo/ratel/v3/registry/etcdv3"
)

// Option 1: use the package-level default registerer (initially a local implementation)
application, _ := app.New()
application.SetRegistry(registry.DefaultRegisterer)

// Option 2: production — connect an etcd registry
// (Set methods return the embedded client config type, so call them separately to preserve the type)
regConfig := regEtcd.New()
regConfig.SetEndpoints([]string{"127.0.0.1:2379"})
reg, err := regConfig.Build() // or regConfig.MustBuild()
if err != nil {
	panic(err)
}
application.SetRegistry(reg)
// Servers automatically register with the registry after starting
```

### Q6: How do I implement graceful shutdown?

`application.Run()` already has built-in signal listening and graceful shutdown: on SIGINT/SIGTERM it automatically runs the `GracefulStop` flow (stops Servers/Workers, waits for background tasks, triggers hooks).

### Q7: Which modules need extra system dependencies?

Only the `img` image processing package needs libvips and pkg-config; all other packages are pure Go dependencies.

---

## Contributing

Contributions via Issues, Pull Requests, or Discussions are welcome. Before submitting a PR, please ensure:

- `make all` passes all checks.
- New or modified code includes corresponding unit tests.
- Exported APIs have Godoc comments, and package comments are written to `doc.go`.
- Backward compatibility is preserved; breaking changes must be described in the PR description.

---

## License

Ratel is licensed under the [Apache License 2.0](LICENSE).
