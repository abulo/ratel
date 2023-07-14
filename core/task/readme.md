# Task

A lightweight distributed cron job library for distributed system

用于分布式系统的轻量分布式定时任务库

---

## 介绍

从 [libi/dcron](https://github.com/libi/dcron) 衍生而来，使用一段时间后发现这个框架存在一些缺陷。

- 不支持本地定时任务
- 节点退出后不能取消注册，在特定场景下会发生任务分配到了已退出的节点上，导致任务无法执行

特性：

在 [libi/dcron](https://github.com/libi/dcron) 基础上，支持：

- 支持本地定时任务
- 支持节点退出取消注册
- 支持懒选举
- 替换[robfig/cron](https://github.com/robfig/cron)为[tovenja/cron](https://github.com/tovenja/cron)库，完全兼容且性能更高

Todo：

- [ ] 任务执行二次确认
- [ ] 完善驱动
- [ ] 完善单元测试
- [ ] 支持自定义logger

## 使用

### 导包

```go
import "cloud/Task"
```

### 示例

```go
package main

import (
	"fmt"
	"cloud/Task"
	"taskdriver/redis"
	"time"
)

func main() {
	stop := cronJob()

	defer stop()
}

func cronJob() func() {

	driver := redis.NewDriver(clientRedis())

	cron := Task.NewTask("test-service", driver, Task.WithLazyPick(true))

	// 分布式任务
	_ = cron.AddFunc("job1", Task.JobDistributed, "*/1 * * * *", func() {
		fmt.Println("执行job1: ", time.Now().Format("15:04:05"))
	})

	// 本地任务
	_ = cron.AddFunc("job2", Task.JobLocaled, "*/1 * * * *", func() {
		fmt.Println("执行job2: ", time.Now().Format("15:04:05"))
	})

	cron.Start()

	return func() { cron.Stop() }
}
```

### Option说明

- 兼容 [tovenja/cron](https://github.com/tovenja/cron) Option
    - WithLocation(loc *time.Location) Option
    - WithSeconds() Option
    - WithParser(p cron.ScheduleParser) Option
    - WithChain(wrappers ...cron.JobWrapper) Option
    - WithLogger(logger cron.Logger) Option

- 自定义节点刷新间隔
    - WithNodeUpdateInterval(dur time.Duration) Option

- 懒选举(任务执行时，拉取节点列表并选举)
    - WithLazyPick(lazy bool) Option