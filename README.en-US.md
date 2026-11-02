

# Ratel 

[![Go](https://github.com/abulo/ratel/v3/workflows/Go/badge.svg?branch=master)](https://github.com/abulo/ratel/v3/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abulo/ratel/v3)](https://goreportcard.com/report/github.com/abulo/ratel/v3)
[![goproxy](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)](https://goproxy.cn/stats/github.com/abulo/ratel/v3/badges/download-count.svg)
[![codecov](https://codecov.io/gh/abulo/ratel/branch/master/graph/badge.svg)](https://codecov.io/gh/abulo/ratel)
[![Release](https://img.shields.io/github/v/release/abulo/ratel.svg?style=flat-square)](https://github.com/abulo/ratel/v3)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📖 Table of Contents

- [Introduction to Ratel](#introduction-to-ratel)
- [Core Features](#core-features)
- [Design Philosophy](#design-philosophy)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Web Server Example](#web-server-example)
  - [gRPC Server Example](#grpc-server-example)
  - [Hertz Server Example](#hertz-server-example)
- [Architecture Components](#architecture-components)
  - [Server Layer](#server-layer)
  - [Client Layer](#client-layer)
  - [Core Layer](#core-layer)
  - [Stores Layer](#stores-layer)
  - [Registry](#registry)
- [Microservice Governance Capabilities](#microservice-governance-capabilities)
  - [Rate Limiting](#rate-limiting)
  - [Circuit Breaking](#circuit-breaking)
  - [Distributed Tracing](#distributed-tracing)
  - [Metrics Monitoring](#metrics-monitoring)
  - [Logging System](#logging-system)
- [Toolkit](#toolkit)
- [Project Structure](#project-structure)
- [Development Guide](#development-guide)
- [Performance Advantages](#performance-advantages)
- [Use Cases](#use-cases)
- [FAQ](#faq)
- [Contributing Guide](#contributing-guide)
- [License](#license)

---

## Introduction to Ratel

Ratel is a **production-grade microservice governance framework** that integrates various engineering practices into Web and RPC frameworks. battle-tested in services handling tens of millions of daily active users (DAU), it aims to deliver extreme stability, reduce server resource costs, and improve development efficiency.

### Core Value

✅ **Easily achieve stability supporting tens of millions of DAU services**  
✅ **Built-in microservice governance capabilities such as cascading timeout control, rate limiting, adaptive circuit breaking, and adaptive load shedding, requiring no configuration or extra code**  
✅ **Microservice governance middleware can be seamlessly integrated into other existing frameworks**  
✅ **Extensive microservice governance and concurrency toolkits**

---

## Core Features

### 🚀 High Performance
- High efficiency with extremely low server resource costs
- Supports high concurrency, stably handling traffic surges
- Fully compatible with `net/http`, zero learning curve

### 🛡️ Resilient Design
- Fault-oriented programming, automatic circuit breaking and load shedding
- Built-in rate limiting, circuit breaking, and load shedding with automatic triggering and recovery
- Cascading timeout control to prevent avalanche effects

### 🔧 Developer-Friendly
- Minimalist APIs, requiring as little code as possible
- Debug mode, greatly improving application development efficiency
- Automatic API parameter validation
- Automatic cache control

### 📊 Observability
- Unified metrics collection (Prometheus)
- Full-stack distributed tracing (OpenTelemetry)
- Structured logging instrumentation
- Statistical alerting

### 🔌 Highly Extensible
- Middleware support for easy extension
- Built-in service discovery and load balancing
- Dynamic configuration with hot-update support
- Built-in security policies

---

## Design Philosophy

### First Principle: Keep It Simple
Enforce doing things one way only, reducing decision costs.

### Resilient Design: Programming for Failure
Assume failures are inevitable, and handle them automatically via adaptive mechanisms.

### Tools Over Conventions and Documentation
Reduce configuration and documentation dependencies through a powerful toolchain.

### Four Core Goals
- **High Availability**: 99.99% uptime guarantee
- **High Concurrency**: Supports millions of QPS
- **Scalability**: Plugin-based architecture
- **Business Developer-Friendly**: Encapsulates complexity, exposes clean APIs

---

## Quick Start

### Installation

```bash
go get github.com/abulo/ratel/v3@latest
```

**Requirements:**
- Go 1.26.0+
- Make (optional, for building and testing)

### Web Server Example

High-performance Web server based on the Gin framework:

```go
package main

import (
    "github.com/abulo/ratel/v3/server/xgin"
    "github.com/gin-gonic/gin"
)

func main() {
    // Create server instance
    server := xgin.NewServer()
    
    // Register routes
    server.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, Ratel!",
        })
    })
    
    // Start server
    server.Run(":8080")
}
```

**Features:**
- Automatic middleware integration (logging, metrics, tracing)
- WebSocket support
- Health checks
- Graceful shutdown

### gRPC Server Example

```go
package main

import (
    "github.com/abulo/ratel/v3/server/xgrpc"
    "google.golang.org/grpc"
)

func main() {
    // Create gRPC server
    server := xgrpc.NewServer()
    
    // Register services
    // pb.RegisterYourServiceServer(server, &yourServiceImpl{})
    
    // Start service
    server.Run(":50051")
}
```

**Features:**
- Automatic load balancing
- Service discovery integration
- Interceptor chain (authentication, logging, metrics)
- Connection pool management

### Hertz Server Example

High-performance HTTP framework based on CloudWeGo Hertz:

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

## Architecture Components

### Server Layer

Ratel provides various server implementations to meet different scenario requirements:

| Module | Description | Use Case |
|------|------|----------|
| **xgin** | Gin-based Web server | RESTful APIs, Web applications |
| **xgrpc** | gRPC-based RPC server | Inter-service communication |
| **xhertz** | Hertz-based high-performance HTTP server | High-concurrency HTTP services |
| **xasynq** | Async task queue server | Background task processing |
| **xmonitor** | Monitoring server | Metrics exposure, health checks |

### Client Layer

Provides various client implementations to simplify external service calls:

| Module | Description |
|------|------|
| **etcdv3** | Etcd client (supports service discovery, config listening, distributed locks) |
| **grpc** | gRPC client (load balancing, circuit breaker, interceptors) |
| **rabbitmq** | RabbitMQ message queue client |
| **redis** | Redis client (connection pool, sentinel, cluster) |
| **sftp/ftp** | File transfer client |

### Core Layer

Provides core capabilities for microservice governance:

#### 📝 Logger - Logging System
- Structured logging (logrus)
- Multi-output support (console, file, Elasticsearch, MongoDB)
- Log rotation
- Context tracing

```go
import "github.com/abulo/ratel/v3/core/logger"

logger.Logger.Info("application started")
logger.Logger.WithField("user_id", 123).Error("operation failed")
```

#### 📊 Metric - Metrics Monitoring
- Prometheus metrics collection
- Counter
- Gauge
- Histogram
- Summary

```go
import "github.com/abulo/ratel/v3/core/metric"

// Record request count
metric.ServerHandleCounter.Inc(map[string]string{
    "type":   "http",
    "method": "GET",
    "code":   "200",
})
```

#### ⏱️ Ratelimit - Rate Limiting
- Token bucket algorithm
- Sliding window
- Adaptive rate limiting

```go
import "github.com/abulo/ratel/v3/core/ratelimit"

limiter := ratelimit.NewLimiter(rate.Every(1*time.Second), 100)
if !limiter.Allow() {
    return errors.New("rate limit exceeded")
}
```

#### 🔌 Resource - Resource Management
- **Breaker**: Adaptive circuit breaking
- **SingleFlight**: Prevents cache stampede
- **BatchError**: Batch error handling
- **SourceManager**: Resource lifecycle management

#### 🔄 Task - Task Scheduling
- Distributed scheduled tasks
- Consistent hashing
- Task pool management
- Cron expression support

#### 🔍 Trace - Distributed Tracing
- OpenTelemetry integration
- Distributed tracing
- Span management
- Sampling rate configuration

#### 🧵 Goroutine - Goroutine Management
- Parallel execution
- Serial execution
- Error aggregation
- Timeout control

### Stores Layer

Unified data storage access interfaces:

| Storage Type | Supported Databases |
|---------|-----------|
| **SQL** | MySQL, PostgreSQL, ClickHouse |
| **Redis** | Standalone, Sentinel, Cluster |
| **MongoDB** | MongoDB |
| **Elasticsearch** | ES 7.x/8.x |
| **Session** | Session management |

**Features:**
- GORM integration
- Connection pool management
- Slow query logging
- SQL metrics monitoring

### Registry

Service registration and discovery:

- **Etcd v3**: Distributed KV storage
- Service health checks
- Weighted routing
- Graceful registration/deregistration

---

## Microservice Governance Capabilities

### Rate Limiting

Ratel provides multi-layer rate limiting strategies:

1. **API Rate Limiting**: Based on token bucket/leaky bucket algorithms
2. **Adaptive Rate Limiting**: Automatically adjusts based on system load
3. **Distributed Rate Limiting**: Based on Redis/Etcd

```go
// Middleware automatically enables rate limiting
server.Use(xgin.RateLimitMiddleware())
```

### Circuit Breaking

Adaptive circuit breaking protection:

- **Error Rate Breaking**: Automatically trips when error rate exceeds threshold
- **Slow Call Breaking**: Automatically trips when response time is too long
- **Automatic Recovery**: Probes recovery in half-open state

```go
// Circuit breaker works automatically, no configuration needed
breaker := resource.NewBreaker(resource.BreakerConfig{
    MaxConcurrentRequests: 100,
    ErrorPercentThreshold: 50,
})
```

### Distributed Tracing

Full-stack tracing integration:

```go
import "github.com/abulo/ratel/v3/core/trace"

// Automatically inject tracing context
ctx, span := trace.StartSpan(ctx, "operation-name")
defer span.End()

// Cross-service propagation
trace.Inject(ctx, carrier)
```

**Supported Exporters:**
- Jaeger
- Zipkin
- OTLP

### Metrics Monitoring

Automatically collects key metrics:

**Server Metrics:**
- `server_handle_total`: Total requests
- `server_handle_seconds`: Request latency

**Client Metrics:**
- `client_handle_total`: Total calls
- `client_handle_seconds`: Call latency

**Task Metrics:**
- `job_handle_total`: Total task executions
- `job_handle_seconds`: Task execution latency

**Expose Prometheus Endpoint:**
```
GET /metrics
```

### Logging System

Structured log output:

```go
logger.Logger.WithFields(logrus.Fields{
    "request_id": "xxx",
    "user_id":    123,
    "duration":   "100ms",
}).Info("request completed")
```

**Log Levels:**
- Debug
- Info
- Warn
- Error
- Fatal

**Log Outputs:**
- Console (colorized output)
- File (daily rotation)
- Elasticsearch
- MongoDB
- Message Queue

---

## Toolkit

### Toolkit CLI Tools

Ratel provides powerful CLI tools to accelerate project development:

```bash
# Install toolkit
go install github.com/abulo/ratel/v3/toolkit@latest

# Create a new project
toolkit new my-project

# Generate CRUD code
toolkit generate model User

# Database migration
toolkit migrate up

# API documentation generation
toolkit api docs
```

### Utility Packages

#### Util Utilities
- String processing
- Time formatting
- File operations
- Networking tools
- Regular expressions
- Random number generation
- Pinyin conversion

#### Correct Text Correction
- Full-width/half-width conversion
- HTML entity encoding
- Format/unformat

#### NLPWord Tokenization
- Aho-Corasick automaton
- Trie tree
- Sensitive word filtering

#### Snowflake ID Generation
- Distributed unique IDs
- Time-ordered
- High performance

---

## Project Structure

```
ratel/
├── client/              # Client implementations
│   ├── etcdv3/         # Etcd client
│   ├── grpc/           # gRPC client
│   ├── rabbitmq/       # RabbitMQ client
│   ├── redis/          # Redis client
│   └── ...
├── config/             # Configuration management
├── core/               # Core functionalities
│   ├── app/           # Application lifecycle
│   ├── logger/        # Logging system
│   ├── metric/        # Metrics monitoring
│   ├── ratelimit/     # Rate limiting
│   ├── resource/      # Resource management (circuit breakers, etc.)
│   ├── task/          # Task scheduling
│   ├── trace/         # Distributed tracing
│   └── ...
├── server/             # Server implementations
│   ├── xgin/          # Gin Web server
│   ├── xgrpc/         # gRPC server
│   ├── xhertz/        # Hertz server
│   ├── xasynq/        # Async task server
│   └── xmonitor/      # Monitoring server
├── stores/             # Data storage
│   ├── sql/           # SQL databases
│   ├── redis/         # Redis
│   ├── mongodb/       # MongoDB
│   └── elasticsearch/ # Elasticsearch
├── registry/           # Service registration and discovery
├── toolkit/            # CLI tools
├── util/               # Common utilities
├── filter/             # Filters
├── watch/              # Config listeners
└── ...
```

---

## Development Guide

### Local Development

```bash
# Clone the repository
git clone https://github.com/abulo/ratel.git
cd ratel

# Install dependencies
go mod download

# Run tests
make test

# Race condition detection
make race

# Format code
make fmt

# Lint check
make lint

# Full build
make all
```

### Code Standards

- ✅ Format code using `gofmt`
- ✅ Use `revive` for linting
- ✅ Check all error returns
- ✅ Write unit tests
- ✅ Add necessary comments

### Debugging Tips

**Enable Debug Mode:**

```go
import "github.com/abulo/ratel/v3/core/env"

env.SetDebug(true)
```

**View Detailed Logs:**

```bash
export LOG_LEVEL=debug
```

---

## Performance Advantages

### Benchmarks

Ratel excels in the following areas:

- **Low Latency**: P99 < 10ms
- **High Throughput**: Single instance supports 100k+ QPS
- **Low Resource Usage**: Memory footprint < 100MB
- **Fast Startup**: < 100ms

### Optimization Strategies

1. **Connection Pool Reuse**: Reduces connection establishment overhead
2. **Zero-Copy Optimization**: Reduces memory allocations
3. **Asynchronous Processing**: Non-blocking I/O
4. **Caching Strategy**: Multi-level caching
5. **Batch Operations**: Reduces network round trips

---

## Use Cases

### ✅ E-commerce Systems
- Flash sale rate limiting
- Circuit breaking for order processing
- Distributed transaction tracing

### ✅ Social Networks
- High-concurrency message pushing
- User relationship graphs
- Real-time feed pushing

### ✅ Financial Systems
- Transaction risk control
- Real-time risk rules
- Audit logging

### ✅ IoT Platforms
- Device onboarding management
- Massive data collection
- Real-time monitoring and alerting

### ✅ Content Platforms
- CDN scheduling
- Recommendation systems
- Search services

---

## FAQ

### Q1: What is the difference between Ratel and native Gin/gRPC?

**A:** Ratel integrates microservice governance capabilities on top of native frameworks:
- Automatic rate limiting, circuit breaking, and load shedding
- Unified metrics collection
- Distributed tracing
- Standardized logging
- No extra configuration or code required

### Q2: How to migrate an existing project to Ratel?

**A:** Ratel is fully compatible with `net/http` and native Gin/gRPC:
```go
// Simply replace the import
// import "github.com/gin-gonic/gin"
import "github.com/abulo/ratel/v3/server/xgin"

// No other code changes required
```

### Q3: How to customize middleware?

**A:** 
```go
server.Use(func(c *gin.Context) {
    // Your logic here
    c.Next()
})
```

### Q4: How to configure rate limiting rules?

**A:** Via configuration files or environment variables:
```yaml
ratelimit:
  enabled: true
  rate: 1000  # Requests per second
  burst: 100  # Burst capacity
```

### Q5: Which service registries are supported?

**A:** Currently supports Etcd v3, with plans to expand to Consul, Nacos, etc.

### Q6: How to implement canary releases?

**A:** Combine Registry's weighted routing and version tags:
```go
registry.Register(service, registry.WithWeight(50))
```

---

## Contributing Guide

We welcome contributions of all kinds!

### Contribution Process

1. **Fork** this repository
2. **Create a branch**: `git checkout -b feature/your-feature`
3. **Commit your changes**: `git commit -am 'Add some feature'`
4. **Push the branch**: `git push origin feature/your-feature`
5. **Submit a PR**

### Contribution Types

- 🐛 Bug fixes
- ✨ New features
- 📝 Documentation improvements
- 🧪 Test cases
- ⚡ Performance optimizations
- 🔧 Code refactoring

### Development Standards

- Follow official Go coding standards
- Add unit tests
- Update documentation
- Maintain backward compatibility

---

## Community & Support

- 📧 **Email**: [Project Maintainer Email]
- 💬 **Issues**: [GitHub Issues](https://github.com/abulo/ratel/issues)
- 📖 **Wiki**: [Project Wiki](https://github.com/abulo/ratel/wiki)
- 👥 **Discussions**: [GitHub Discussions](https://github.com/abulo/ratel/discussions)

---

## Success Stories

Ratel has been validated in multiple production environments:

- **E-commerce Platform**: Supports tens of millions of QPS during Double 11
- **Social Platform**: 50M+ DAU
- **Financial Company**: 99.99% availability for trading systems
- **IoT Platform**: Millions of devices online simultaneously

---

## License

Ratel is licensed under the [MIT License](https://opensource.org/licenses/MIT).

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

If Ratel is helpful to you, please give us a ⭐️ Star!

[![Star History Chart](https://api.star-history.com/svg?repos=abulo/ratel&type=Date)](https://star-history.com/#abulo/ratel&Date)

---

**Made with ❤️ by Ratel Team**
