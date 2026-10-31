# 性能审查章节：core/logger、core/task、core/metric、core/trace

> 审查范围（仅限本章节，未越界）：
> - `core/logger/`：`entry/`、`queue/`、`es/`、`mongo/`、`sql/` 子目录及包内全部 `.go`（测试桩文件已通读，非热路径，未单独列条目）。
> - `core/task/`：`cron/`、`consistenthash/`、`driver/` 子目录及包内全部 `.go`。
> - `core/metric/` 全部 `.go`；`core/trace/` 全部 `.go`。
>
> 审查方式：逐文件人工通读源码，先借注释理解 API 语义，再定位热点路径与可疑实现；每条结论均回读原文核对 `文件:行号` 与函数名。本机未发现可用 `go` 工具链，未运行 `go vet` / `staticcheck` / benchmark，涉及分配量级与耗时的判断均在条目内标注「需实测验证」。
>
> 优先级定义：**P0**＝同步热路径上的明显性能缺陷，应优先修改；**P1**＝中等影响、性价比高的改进；**P2**＝低影响或改进空间有限的点。本章节未发现可定级 P0 的同步热路径硬缺陷（P0 数量为 0）。

---

## 一、条目清单（按子包分组；组内 P0 → P1 → P2 排序）

### core/logger/queue

#### [P1] `core/logger/queue/list_queue.go:57` `ListQueue.dispatcher`：空队列 10ms 固定 sleep 轮询；默认无长度上限；Push 两次加锁
- 问题：`ListQueue` 的消费循环存在三个相互关联的效率问题：
  1. **固定周期轮询**：`core/logger/queue/list_queue.go:57` `func (q *ListQueue) dispatcher()`，队列空时 `core/logger/queue/list_queue.go:68` `time.Sleep(time.Millisecond * 10)`，新任务入队后最坏要等满 10ms 才被取出，属于纯轮询延迟；空转期间反复 `RLock` 检查 `Len()`（`list_queue.go:59-65`）。
  2. **默认无界队列**：`core/logger/queue/list_queue.go:11` `NewListQueue` 内部调用 `core/logger/queue/list_queue.go:16` `NewListQueueWithMaxLen(maxThread, 0)`，`maxLen=0` 表示不限制长度；后端写入变慢时链表会持续增长，没有背压上限与丢弃策略（`Push` 仅在 `maxLen > 0` 时才检查长度，`list_queue.go:87-94`）。
  3. **Push 双临界区**：`core/logger/queue/list_queue.go:82` `func (q *ListQueue) Push(job Jober)`，先 `RLock/RUnlock` 查长度（`list_queue.go:88-93`），再 `Lock/Unlock` 入队（`list_queue.go:97-99`），一次入队两次锁操作；检查与入队之间也存在长度竞态（正确性问题见附记）。
- 影响评估：轮询延迟主要影响低负载 / 队列间歇为空时的日志落盘时效（最坏 +10ms/条）；无界增长在后端故障或慢节点下表现为内存随积压日志单调上升，是内存可用性风险。该队列是否为默认启用实现需结合启动配置确认（**需实测验证**：线上走 `queue.Run` 还是 `queue.RunListQueue`）；若是默认路径，10ms 轮询在高频日志场景不可忽视。
- 优化建议：
  1. 用 `sync.Cond`（或带缓冲信号 channel）替代 sleep 轮询：`Push` 后 `Signal()`，dispatcher 在 `Cond.Wait()` 上阻塞，任务到达即唤醒；
  2. 默认 `maxLen` 给一个合理上限（复用 `defaultOptions.maxQueues=512` 口径），超限按策略丢弃最旧 / 最新日志并计数；
  3. `Push` 改为单个 `Lock` 临界区内完成长度检查 + `PushBack`，去掉先 RLock 后 Lock 的两段式。

#### [P2] `core/logger/queue/chan_queue.go:76` `Queue.Push`：有界通道满时在业务同步栈内阻塞发送，缺少非阻塞降级策略
- 问题：`core/logger/queue/chan_queue.go:70-77` `func (q *Queue) Push(job Jober)`：`q.wg.Add(1)` 后 `q.jobQueue <- job`；缓冲（`maxQueues`，默认 512，见 `core/logger/es/hook.go:17` 等三个钩子的 `defaultOptions`）写满时发送方阻塞。`Fire` 由 logrus 在业务请求 goroutine 内同步调用，后端消费变慢时业务线程会被这条通道发送卡住，形成反向背压；没有 `select { case q.jobQueue <- job: default: /* 丢弃/计数 */ }` 的降级路径。
- 影响评估：正常负载下无影响；故障 / 突发日志洪峰场景下表现为请求尾延迟被日志后端拖累。触发概率取决于后端写入速率与 512 缓冲、2 worker 的匹配度，**需实测验证**线上积压水位。
- 优化建议：评估场景取舍——保持阻塞（强不丢日志）时可考虑把入队改为带超时的 `select`，超时后丢弃并累加丢弃计数指标；或至少在文档 / 默认值上明确“阻塞式背压”语义。注意 `ListQueue` 默认静默丢弃，两种队列策略不一致，建议统一为可配置项。

### core/logger/es

#### [P1] `core/logger/es/exec.go:42` `defaultExec.Exec`：每条日志一次独立 ES Index HTTP 请求，且每条都生成 UUID，未做批量聚合
- 问题：`Exec` 对每条日志单独 `Index().Id(uuid.New().String()).BodyJson(entry).Do(ctx)`，worker 循环内也是一条一次请求，没有 bulk 批量缓冲；同时显式调用 `uuid.New().String()` 生成文档 ID，而 ES 在不传 ID 时会自动生成 `_id`，该 UUID 对日志场景没有业务价值。
  - 位置：`core/logger/es/exec.go:40` `func (e *defaultExec) Exec(entry *entry.Entry) error`；`core/logger/es/exec.go:42` `_, err := e.client.Index().Index(e.index).Id(uuid.New().String()).BodyJson(entry).Do(ctx)`。
- 影响评估：后台 worker 路径。日志量级 = ES 请求量级，HTTP 头部、`uuid.New()`（版本 1，含时间戳读取与随机源）与 `BodyJson` 反射序列化开销按条重复；后端积压时 512 容量队列（`core/logger/es/hook.go:17`）会先被打满。内网 ES 单条 Index 通常 1~数 ms，批量后单条边际成本可降一个量级；具体吞吐提升幅度**需实测验证**。
- 优化建议：
  1. worker 内引入 `BulkIndexer`（或定时 + 条数双阈值刷盘：累计 N 条或每 M 毫秒 flush 一次），用 ES Bulk API 一次提交多条；
  2. 去掉 `.Id(uuid.New().String())`，由 ES 自动生成 `_id`；
  3. 复用 ES 客户端连接（已复用），批量时优先使用客户端自带的 bulk flush 与并发控制。

### core/logger/mongo

#### [P1] `core/logger/mongo/exec.go:51` `defaultExec.Exec`：每条日志先执行一次 MongoDB `Ping` 再插入，额外多一次网络往返
- 问题：`Exec` 在每次写日志前都调用 `e.client.NewCollection(e.collection)`；而 `stores/mongodb` 的 `NewCollection` 内部第一步就是 `m.Ping()`（见 `stores/mongodb/clientmanager.go:196-198`，本章节仅引用、不展开审查），随后才构造 `DecoratedCollection` 包装器。也就是说每写一条日志，除 `InsertOne` 外还多一次 Ping 网络往返、一次熔断器 `resource.GetBreaker` 查表和一次包装器构造。
  - 位置：`core/logger/mongo/exec.go:40` `func (e *defaultExec) Exec(entry *entry.Entry) error`；`core/logger/mongo/exec.go:51` `handler, err := e.client.NewCollection(e.collection)`；`core/logger/mongo/exec.go:55` `handler.InsertOne(ctx, item)`。
- 影响评估：位于后台 worker 消费路径（默认 `maxWorkers=2`，见 `core/logger/mongo/hook.go:18`），不阻塞业务请求；但日志量较高时，每条 Mongo 日志 = Ping RTT + 插入 RTT 两次往返，写入吞吐被 RTT 直接腰斩，worker 更容易成为队列瓶颈，进而把压力回传至 `core/logger/queue` 的有界通道。量级取决于 Mongo RTT（内网通常 0.1~1ms/次）；**Ping 次数 / 日志条数 1:1 是确定事实，收益幅度需实测**。
- 优化建议：在 `NewExec` / `NewExecWithURL`（`core/logger/mongo/exec.go:23-38`）构造期一次性解析并缓存 `*DecoratedCollection`（或 `*mongo.Collection`）到 `defaultExec` 结构体字段；`Exec` 内直接复用缓存句柄执行 `InsertOne`，删除每条日志的 `NewCollection` 调用。若担心句柄失效，可在 `InsertOne` 报错时惰性重建一次，而不是每条都 Ping。

### core/logger/sql

#### [P2] `core/logger/sql/exec.go:50` `defaultExec.Exec`：每条日志一次 `json.Marshal` + GORM 单条 `Create`，未批量
- 问题：`Exec` 内 `core/logger/sql/exec.go:50` `data, _ := json.Marshal(entry.Data)`，随后 `core/logger/sql/exec.go:52` `e.client.Create(daoItem).Error`，每条日志独立一次 INSERT；`json.Marshal` 的返回错误被 `_` 忽略（正确性问题见附记）。`null.StringFrom` / `null.TimeStampFrom` / `null.JSONFrom`（`exec.go:44-51`）也按条构造包装值。
- 影响评估：worker 后台路径；单条 INSERT 在事务 / 连接复用下开销可接受，但日志高频时数据库往返次数 = 日志条数，且 `json.Marshal` 反射开销按条重复。仅当 SQL 钩子实际挂载时才有影响，**需实测验证**该 sink 是否启用。
- 优化建议：worker 内按条数 / 时间双阈值攒批，使用 GORM `CreateInBatches` 一次插入多行；`entry.Data` 的 JSON 序列化结果在批量构造时统一处理。

### core/logger（es / mongo / sql 三钩子共享模式）

#### [P1] `core/logger/es/hook.go:158` / `core/logger/mongo/hook.go:149` / `core/logger/sql/hook.go:149` `Hook.Fire`：每条日志都调用一次 `os.Hostname()`，主机名应在构造期缓存
- 问题：三个后端钩子的 `Fire` 在同步调用栈内（logrus 触发 Hook 时）每条日志都执行 `hostName, err := os.Hostname()`。`os.Hostname` 标准库不在 `os` 层缓存结果，每次落到系统调用（`gethostname` / `uname`）并返回字符串拷贝；而主机名在进程生命周期内恒定不变。
  - 位置列表：
    - `core/logger/es/hook.go:151` `func (h *Hook) Fire(entryLogrus *logrus.Entry) error`，`core/logger/es/hook.go:158` `hostName, err := os.Hostname()`；
    - `core/logger/mongo/hook.go:142` `func (h *Hook) Fire(entryLogrus *logrus.Entry) error`，`core/logger/mongo/hook.go:149` `hostName, err := os.Hostname()`；
    - `core/logger/sql/hook.go:142` `func (h *Hook) Fire(entryLogrus *logrus.Entry) error`，`core/logger/sql/hook.go:149` `hostName, err := os.Hostname()`。
- 影响评估：`Fire` 位于业务日志同步路径，每个被该 Hook 覆盖级别的日志条目都会走到。单次 `os.Hostname` 系统调用本身很快（微秒级），但在高频日志（每请求多条 Info/Debug）下属于 1:1 重复的固定开销，且三个钩子各自重复调用一次；消除后为纯收益，风险为零。单次耗时占比**需实测验证**（`pprof` 看 syscall 与 Fire 占比）。
- 优化建议：在 `Hook` 结构体新增 `hostName string` 字段，于 `New(...)`（`core/logger/es/hook.go:116`、`core/logger/mongo/hook.go:107`、`core/logger/sql/hook.go:107`）构造期调用一次 `os.Hostname()`（失败回填 `"unknown"`），`Fire` 直接读字段；三个钩子复用同一实现思路即可。

#### [P2] `core/logger/es/hook.go:156` / `core/logger/mongo/hook.go:147` / `core/logger/sql/hook.go:147` `Hook.Fire`：条目转换过程中每条日志多次小分配（Sprintf、Level 大写转换、Now、Entry/Job/闭包）
- 问题：三个 `Fire` 在同步路径内重复执行以下可优化的小分配 / 计算（同一模式三处，位置并列）：
  1. `fmt.Sprintf("%s:%d", caller.File, caller.Line)` 拼接调用位置（`es/hook.go:156`、`mongo/hook.go:147`、`sql/hook.go:147`）；
  2. `level := entryLogrus.Level.String()` 后 `strings.ToUpper(level)`（`es/hook.go:162,169`；`mongo/hook.go:153,160`；`sql/hook.go:153,160`）——logrus 级别只有固定 7 个取值，结果集恒定；
  3. `Timestamp: util.Now()`（`es/hook.go:165`、`mongo/hook.go:156`、`sql/hook.go:156`）：`util.Now()` 内部 `time.Now().In(cstZone)`（`util/time.go:127-130`），`.In` 会复制 `time.Time` 结构体并绑定 Location；
  4. `queue.NewJob(newEntry, func(v any){ h.exec(...) })`（`es/hook.go:172`、`mongo/hook.go:163`、`sql/hook.go:163`）：每条日志同时在堆上分配 `entry.Entry`、`job` 结构体、捕获 `h` 的闭包，并经 `any` 装箱。
- 影响评估：单条开销为纳秒~几十纳秒级的小分配集合，不构成明显延迟；但在每请求多条日志的服务里，分配次数随日志条数线性增长，会增加 GC 压力与小对象堆占用。各项收益幅度**需实测验证**（建议 `BenchmarkHookFire` + `benchstat` 看 allocs/op）。
- 优化建议：
  1. 构造静态映射 `map[logrus.Level]string`（预存 `INFO`/`ERROR`… 大写串），去掉 `String()` + `ToUpper` 两次字符串处理；
  2. 调用位置用 `strconv.AppendInt` 写入复用 `[]byte` / `strings.Builder`，避免 `fmt.Sprintf` 的反射与格式解析；
  3. `entry.Entry` 与 `queue.job` 可纳入 `sync.Pool`（worker 执行完回收），闭包改为预定义的固定回调 + 结构体字段，消除每条日志的闭包分配；
  4. 时间戳可直接 `time.Now()`（后端格式化再做时区转换），避免 `util.Now().In(...)` 的冗余转换。

### core/task（根包）

#### [P1] `core/task/task.go:208` `Task.allowThisNodeRun`：集群稳态下每次任务触发都无条件弹出最近任务集合并新建 goroutine，空集也不例外
- 问题：开启 `WithClusterStable`（`core/task/option.go:67-71`）后，任务每次被调度都会进入 `allowThisNodeRun`；稳态分支内 `core/task/task.go:207-209` 不判断 `PopAllJobs()` 是否为空，直接 `go d.reRunRecentJobs(d.recentJobs.PopAllJobs())`。即：正常稳态期间每次任务触发都发生——加锁弹出堆（`core/task/recentjobpacker.go:80-90`）、分配空 `[]string`、创建一个 goroutine，而 `reRunRecentJobs` 入口还会打一条 `Infof("reRunRecentJobs: length=%d", ...)`（`core/task/task.go:285`），空切片时 `length=0` 同样入日志队列（该日志又会走 logger Hook 热路径，形成自激开销）。
  - 位置：`core/task/task.go:196` `func (d *Task) allowThisNodeRun(jobName string) (ok bool)`；`core/task/task.go:208` `go d.reRunRecentJobs(d.recentJobs.PopAllJobs())`；`core/task/task.go:284` `func (d *Task) reRunRecentJobs(jobNames []string)`。
- 影响评估：位于调度主路径（每个注册任务每次触发一次）。稳态集群中 `recentJobs` 通常为空，开销为纯浪费：每次触发 1 次互斥锁 + 1 次 goroutine 调度 + 1 条 Info 日志及其队列入队 / 后端写链路。任务数 × 触发频率越高，累积越明显；在分钟级高频任务上 goroutine 创建与空日志的开销不可忽略。是否显著取决于任务触发密度，**需实测验证**。
- 优化建议：先取出结果、判空再起 goroutine：
  ```go
  if d.recentJobs != nil {
      if names := d.recentJobs.PopAllJobs(); len(names) > 0 {
          go d.reRunRecentJobs(names)
      }
  }
  ```
  同时 `reRunRecentJobs` 在空切片时直接 return，避免入口 `Infof`；`core/task/task.go:211` 与 `207` 重复的 `d.recentJobs != nil` 判断可合并。

### core/task/driver

#### [P1] `core/task/driver/etcddriver.go:83` `EtcdDriver.watcher`：Watch 使用 `context.Background()`，`Stop` 不会取消监听协程，存在 goroutine 与 watch 流泄漏风险
- 问题：`watcher` 在 `core/task/driver/etcddriver.go:81-93`，`core/task/driver/etcddriver.go:83` `e.cli.Watch(context.Background(), prefix, ...)` 绑定后台上下文；而 `Stop`（`core/task/driver/etcddriver.go:192-196`）只调用 `e.cancel()` 取消 `e.ctx`，该上下文仅被 `keepHeartBeat`（`core/task/driver/etcddriver.go:141-158`）监听。`watcher` 的 `range rch` 永远收不到取消信号，`Stop` 后 watcher goroutine 与底层 etcd watch 流继续挂着；若驱动被 `Stop` 后再次 `Start`（`core/task/driver/etcddriver.go:180-189` 会再次 `go e.watcher(...)`），旧 watcher 不会退出，goroutine 数量随启停次数累积。
- 影响评估：非高频请求路径，属于生命周期缺陷；在进程长期运行、任务管理器反复重建（热重启、配置重载、多 Task 实例轮换）的场景下表现为 goroutine 缓慢泄漏与 etcd 侧残留 watch 连接。单次泄漏体积小，但无上限增长。是否会在当前部署形态触发反复 `Start/Stop` **需实测验证**。
- 优化建议：将 `e.ctx` 传入 `Watch`：`e.cli.Watch(e.ctx, prefix, clientv3.WithPrefix())`；`Stop` 的 `e.cancel()` 即可同时让 `watcher` 的 `range` 退出。并在 `Start` 中避免重复启动 watcher（`started` 标志目前在 etcd 驱动上缺失，对比 `core/task/driver/redisdriver.go:62` 的写法）。

#### [P2] `core/task/driver/rediszsetdriver.go:55` `RedisZSetDriver.GetNodes`：在互斥锁临界区内执行 Redis 网络查询
- 问题：`core/task/driver/rediszsetdriver.go:54` `func (rd *RedisZSetDriver) GetNodes(ctx context.Context) (...)`，`core/task/driver/rediszsetdriver.go:55-56` `rd.Lock(); defer rd.Unlock()` 覆盖整段 `ZRangeByScore` 网络调用。锁实际保护的是 `started` 等状态字段，远程查询不需要持有该锁。
- 影响评估：`GetNodes` 由节点池刷新循环每 3 秒调用一次（`core/task/nodepool.go:138`），频率低；但锁覆盖 I/O 会与 `Start` / `Stop`（`rediszsetdriver.go:72`、`91`）相互阻塞，网络抖动时临界区被拉长。影响面小，P2。对照 `core/task/driver/redisdriver.go:86-89` 的 `scan` 实现并未持锁。
- 优化建议：锁内只读取需要保护的状态（或干脆去掉 `GetNodes` 的锁），Redis 查询放到锁外执行。

### core/task/cron

#### [P2] `core/task/cron/cron.go:229` `Cron.run`：每次事件唤醒都对全量条目做 `sort.Sort`；`removeEntry` 线性重建切片
- 问题：调度主循环 `core/task/cron/cron.go:218` `func (c *Cron) run()` 每轮外层迭代都执行 `core/task/cron/cron.go:229` `sort.Sort(byTime(c.entries))`；内层 `select` 处理完一条触发 / 新增 / 删除 / 快照事件后 `break` 回到外层重新全量排序。`core/task/cron/cron.go:318-326` `removeEntry` 通过遍历整个切片过滤目标 ID 再 `append` 重建，为 O(n)。
- 影响评估：条目数量通常为个位数~数十个，`n log n` 排序与 O(n) 重建在该量级下成本可忽略；仅在注册任务极多（成百上千）且频繁动态增删时才会显现。属实现层面的常数优化空间，**当前规模下收益有限**。
- 优化建议：触发 / 增删后只更新受影响条目位置，或维护一个最小堆（`heap.Interface`，按 `Next` 排序）替代全量排序；`removeEntry` 改为切片元素移动 + 截断，避免整体重建分配。

#### [P2] `core/task/cron/cron.go:283` `Cron.startJob`：每次触发新建 goroutine，无并发上限与并发互斥包装
- 问题：`core/task/cron/cron.go:282` `func (c *Cron) startJob(j Job)` 直接 `c.jobWaiter.Go(func(){ j.Run() })`；若上一个周期任务尚未结束，下一周期触发仍会再起一个 goroutine 并发执行同一任务，框架层没有最大并发数或 `Recover` 包装（默认 `Chain` 为空，`core/task/cron/cron.go:100`）。
- 影响评估：这是 robfig/cron 风格的默认设计取舍，不是必然缺陷；仅在任务执行时长超过调度周期、或任务本身不可重入时，表现为 goroutine 与任务并发数随超期周期累积。列为 P2 供配置链选择包装器时参考。
- 优化建议：通过 `cron.WithChain`（`core/task/cron/option.go:32-36`）叠加 `Recover` 与互斥 / 单飞包装器；如需框架默认兜底，可在 `Chain` 中默认放入 recover 包装，避免单个任务 panic 拖垮调度器（可靠性问题，附记）。

### core/task/consistenthash

#### [P2] `core/task/consistenthash/consistenthash.go:79` `Map.Get`：每次任务归属判定都做一次 `string → []byte` 拷贝分配
- 问题：`core/task/consistenthash/consistenthash.go:75` `func (m *Map) Get(key string) string`，`core/task/consistenthash/consistenthash.go:79` `hash := int(m.hash([]byte(key)))`；`crc32.ChecksumIEEE` 只接收 `[]byte`，每次调用把 `key` 拷贝成字节切片。`Get` 在每次任务触发的 `CheckJobAvailable`（`core/task/nodepool.go:105`）中执行一次。
- 影响评估：纯内存操作，一次小对象分配（长度 = 任务名）；触发频率 = 任务调度频率，量级小，列为 P2。`Add`（`consistenthash.go:63`）在节点刷新（默认 3s 一次，`core/task/nodepool.go:133-148`）时才执行，`strconv.Itoa(i)+"_"+key` 拼接 + `[]byte` 转换同样低频，不单列。
- 优化建议：保存一个 `hash.Hash32`（`crc32.NewIEEE()`）实例，用 `hash.Reset(); hash.WriteString(key); sum32()` 避免 string→[]byte 拷贝（`Hash` 函数类型可扩展为接收 `string`）；`Map` 本身仍由 `NodePool` 的 `rwMut` 保护，无需额外加锁。

### core/metric

#### [P2] `core/metric/counter.go:53` / `core/metric/gauge.go:53` / `core/metric/histogram.go:43` / `core/metric/summary.go:40`：上报包装器在 Prometheus 向量之上又包一层变参方法，请求路径每次都走 `WithLabelValues` 向量查找
- 问题：四个向量包装器的上报方法都把 `labels ...string` 原样转发给 Prometheus 的 `WithLabelValues`：
  - `core/metric/counter.go:52-54` `Inc`、`core/metric/counter.go:57-59` `Add`；
  - `core/metric/gauge.go:52-54` `Inc`、`core/metric/gauge.go:57-59` `Add`、`core/metric/gauge.go:62-64` `Set`；
  - `core/metric/histogram.go:42-44` `Observe`；
  - `core/metric/summary.go:39-41` `Observe`。
  - 需要说明：包装器这层 `labels...` 转发本身不会新建第二次切片（调用方为变参构造的切片被原样传入 Prometheus 方法）；可优化点在于每次请求都经过包装方法 → `WithLabelValues` 向量内的标签元组查找（`prometheus.CounterVec` 内部 `children` map + 标签校验），而热路径上标签组合通常是固定的，查找结果本可在初始化期一次性解析为子指标句柄。
- 影响评估：指标上报位于每个请求的拦截器热路径；单次为一次方法间接调用 + 向量 map 查找，Prometheus 内部 `Child` 虽带缓存，但标签名 / 数量校验与 map 读仍在每次上报时执行。量级低，是否构成可见 allocs/op 或延迟**需实测验证**（`WithLabelValues` 自身的 map 查找开销主体在 Prometheus 库内，本包只新增了一层薄包装，不单独立项评级）。
- 优化建议：为标签数量固定的高频场景提供定参方法（如 `Inc2(l1, l2 string)` / `Observe2(v, l1, l2 string)`），去掉一层变参展开；更彻底的做法是业务侧在初始化时缓存解析后的 `prometheus.Counter` / `prometheus.Histogram` 子指标句柄，请求路径直接调用子指标方法，绕过每次的向量查找。

### core/trace

#### [P2] `core/trace/trace.go:50` `Tracer.Start`：每次调用重复读取全局 propagator，并追加构造 SpanKind 选项切片
- 问题：`core/trace/trace.go:46` `func (t *Tracer) Start(...)`：
  - `core/trace/trace.go:48` 与 `core/trace/trace.go:55` 两次调用 `otel.GetTextMapPropagator()`；
  - `core/trace/trace.go:50` `opts = append(opts, trace.WithSpanKind(t.kind))`，每次 `Start` 都向调用方传入的变参切片追加一个选项（`opts` 为 nil 时触发新切片分配）。
- 影响评估：每个进出请求的 Span 启动执行一次；`GetTextMapPropagator` 为全局接口读取（内部原子加载），开销很小；`append` 分配只在 `opts` 容量不足时发生。整体 P2，优化空间有限。
- 优化建议：`NewTracer`（`core/trace/trace.go:38-40`）时把 propagator 与 `trace.WithSpanKind(t.kind)` 预存为 `Tracer` 字段；`Start` 复用字段，避免每次读取全局值与重复 `append`。

#### [P2] `core/trace/propagation.go:50` `Jaeger.Inject`：headers 切片无容量预分配，flags 用 `fmt.Sprintf("%x", ...)` 构造
- 问题：`core/trace/propagation.go:48` `func (jaeger Jaeger) Inject(...)`：`core/trace/propagation.go:50` `headers := []string{}` 后按 3~4 段 `append`（`propagation.go:54`、`56`、`58`、`60`），切片会经历一次扩容；flags 三个分支分别 `fmt.Sprintf("%x", ...)`（`propagation.go:56/58/60`）格式化固定取值。
- 影响评估：仅在 Span 结束 / 下游注入头时执行，且受采样比例约束；单条分配极小，P2。`debugFromContext`（`propagation.go:160-168`）每次 `ctx.Value` 查找，属 context 标准做法，不单列。
- 优化建议：`headers := make([]string, 0, 4)`；flags 固定为 `"02"` / `"01"` / `"00"` 三个常量字符串，去掉 `fmt.Sprintf`。

---

## 二、条目统计（按子包分组）

> 计数口径：跨 es/mongo/sql 三个钩子的同一模式合并为一条条目，计数记在“共享”行，子包行只计各自独有条目，避免重复计算。

| 子包 | P0 | P1 | P2 | 备注 |
|---|---:|---:|---:|---|
| `core/logger/entry` | 0 | 0 | 0 | 纯数据结构，无热点逻辑 |
| `core/logger/queue` | 0 | 1 | 1 | P1：`list_queue.go` dispatcher/Push；P2：`chan_queue.go` Push 背压 |
| `core/logger/es` | 0 | 1 | 0 | P1：`exec.go` 单条请求/UUID |
| `core/logger/mongo` | 0 | 1 | 0 | P1：`exec.go` 每条 Ping |
| `core/logger/sql` | 0 | 0 | 1 | P2：`exec.go` 单条 Marshal/Create |
| 共享：es/mongo/sql 三钩子 `Fire` | 0 | 1 | 1 | P1：Hostname 缓存；P2：条目转换分配合并条目 |
| **core/logger 小计** | **0** | **4** | **3** | |
| `core/task`（根包） | 0 | 1 | 0 | `task.go:208` 空回放 goroutine |
| `core/task/cron` | 0 | 0 | 2 | 全量排序/线性移除；startJob 无并发上限 |
| `core/task/consistenthash` | 0 | 0 | 1 | `Get` 的 string→[]byte 拷贝 |
| `core/task/driver` | 0 | 1 | 1 | P1：etcd watcher 泄漏；P2：zset GetNodes 锁覆盖 I/O |
| **core/task 小计** | **0** | **2** | **4** | |
| `core/metric` | 0 | 0 | 1 | 四包装器变参转发合并为一条 |
| `core/trace` | 0 | 0 | 2 | `trace.go` propagator/opts；`propagation.go` headers 切片 |
| **总计** | **0** | **6** | **10** | |

---

## 三、高价值候选列表（供全局 Top 10 评选）

> 排序按“预计收益 / 改动成本”性价比；行号已逐条回读核对。

| # | 位置 | 问题一句话 | 预期收益 |
|---:|---|---|---|
| 1 | `core/logger/mongo/exec.go:51` `Exec` | 每条日志多一次 Mongo Ping RTT，构造期缓存 collection 句柄即可消除 | Mongo 日志写入 RTT 减半，worker 吞吐接近翻倍（需实测基线） |
| 2 | `core/logger/es/exec.go:42` `Exec` | 每条日志一次 ES Index 请求 + 冗余 UUID，未做批量 | 改 Bulk 后 ES 写吞吐随批量大小提升，单条边际序列化/HTTP 开销大幅下降（需实测 bulk 阈值） |
| 3 | `core/logger/es/hook.go:158` / `mongo/hook.go:149` / `sql/hook.go:149` `Fire` | 每条日志重复 `os.Hostname()`，构造期缓存一次 | 去掉同步日志路径上每条日志的主机名系统调用，改动极小、零行为风险 |
| 4 | `core/task/task.go:208` `allowThisNodeRun` | 稳态下每次任务触发都空跑一次 goroutine + 空日志回放 | 开启 `WithClusterStable` 时消除调度路径上每触发一次的锁、goroutine 与一条空 Info 日志链路开销 |
| 5 | `core/logger/queue/list_queue.go:68` `dispatcher` | 10ms sleep 轮询 + 默认无界队列 | 消除低负载下最坏 10ms 入队延迟；补上内存上限，避免后端故障时无界积压 |
| 6 | `core/task/driver/etcddriver.go:83` `watcher` | Watch 绑定 `context.Background()`，Stop 不取消，反复启停泄漏 goroutine/watch 流 | 根治驱动反复重建场景的 goroutine 与连接泄漏 |
| 7 | `core/logger/sql/exec.go:50-52` `Exec` | 单条 `json.Marshal` + 单条 `Create`，未攒批 | SQL 日志 sink 的数据库往返次数下降为 1/N（前提：该 sink 实际启用，需确认） |
| 8 | `core/logger/es/hook.go:156` 等三钩子 `Fire` | 每条日志 Sprintf/ToUpper/时间转换/Entry+Job+闭包多分配 | 降低每请求多条日志场景的 GC 分配次数（allocs/op 可量化，需 benchmark） |
| 9 | `core/metric/counter.go:53` 等四处上报包装器 | 请求路径每次走 `WithLabelValues` 向量查找，可初始化期缓存子指标句柄 | 高 QPS 拦截器路径减少每次的向量 map 查找，属低风险微调 |
| 10 | `core/task/consistenthash/consistenthash.go:79` `Get` | string→[]byte 每次拷贝 | 消除任务归属判定上的一次小对象分配，收益较小但改动直接 |

---

## 四、拿不准、需实测验证的条目清单

1. **`os.Hostname()` 单次耗时与占比**（P1，三钩子 `Fire`）：系统调用在目标 macOS/Linux 上的实际耗时、是否已被运行时 / OS 缓存，需 `pprof` + Hook 级 benchmark 确认 Fire 中该调用的时间占比。
2. **`Fire` 每条目分配次数与 GC 影响**（P2，三钩子转换分配）：需 `BenchmarkHookFire` 带 `-benchmem` 测量 allocs/op、B/op，并在真实日志速率下看 GC pause 占比；`sync.Pool` 回收后的净收益需对照实测。
3. **`ListQueue` 是否为线上默认实现**（P1，`list_queue.go`）：需确认启动走 `queue.Run`（有界 channel）还是 `queue.RunListQueue`（无界链表）；若从不启用，该 P1 降级为 P2。
4. **`chan_queue.go` 满队列阻塞概率**（P2）：需在目标日志 QPS、512 缓冲、2 worker、后端延迟分布下观察通道水位与阻塞时长，再决定是否需要非阻塞降级。
5. **Mongo/ES/SQL 三个 sink 的启用情况**：`mongo/exec.go:51`、`es/exec.go:42`、`sql/exec.go:50` 的收益以对应钩子实际挂载为前提，需先核实现网配置。
6. **etcd watcher 泄漏的复现条件**（P1，`etcddriver.go:83`）：需在测试环境执行 `Start → Stop → Start` 循环，用 goroutine profile 验证旧 watcher 是否残留；当前部署形态若从不反复重建驱动，紧迫性下调。
7. **`metric` 包装器向量查找量级**（P2）：需高 QPS 压测对比“包装层 + `WithLabelValues`”与“缓存子指标句柄直调”的 allocs/op 与延迟；Prometheus 内部 `WithLabelValues` 的 map 成本另算、不在本包。
8. **`cron.go:229` 全量排序的实际规模**（P2）：需统计线上注册任务总数；任务数 < 50 时该条目可不做改造。
9. **`allowThisNodeRun` 空回放频率**（P1，`task.go:208`）：仅在 `WithClusterStable` 开启且任务调度较密时才显著；需按实际任务频率估算 goroutine 创建次数。

---

## 五、跨模块提示（超出本章节范围，仅登记线索，不展开审查）

1. **`stores/mongodb/clientmanager.go:196-203` `MongoDB.NewCollection`**：内部每条先 `Ping()` 再构造包装器，是 `core/logger/mongo/exec.go:51` 每条日志多一次 RTT 的直接来源；建议全局审视该方法是否应去掉 Ping（或提供不 Ping 的缓存句柄 API），供其他使用方一并受益。
2. **`util/time.go:127-130` `util.Now`**：`time.Now().In(time.Local)` 在日志 / 时间热路径被反复调用；是否应提供不做时区转换的轻量时间戳，归 `util` 包审查范围。
3. **logrus 本体格式化器（JSON/Text）的缓冲复用、级别判断顺序**：本包仅持有全局 `logrus.Logger`（`core/logger/logger.go:10-11`），格式化器与缓冲策略在 `sirupsen/logrus` 依赖侧，不在本章节范围。
4. **`stores/elasticsearch` 客户端 bulk 能力**：`core/logger/es/exec.go` 批量改造需要确认底层客户端是否已暴露 BulkIndexer，归 `stores/elasticsearch` 审查范围。

---

## 六、已通读但未发现性能问题的包 / 导出 API（覆盖说明）

- **`core/logger/entry/entry.go`**：`Entry` 为纯数据载体结构体，无方法、无锁、无分配逻辑；热点只发生在钩子构造它的路径（已列入共享 P2 条目）。
- **`core/logger/logger.go`**：仅在包加载时创建全局 `logrus.Logger` 与 `logrus.Entry`，属一次性初始化。
- **`core/logger/queue/job.go`**：`job` / `syncJob` 为回调包装结构体；`NewSyncJob` 的 `result` 通道容量为 1，与同步等待语义匹配；`NewJob` 的分配计入共享 P2 条目，未单独列项。
- **`core/logger/queue/worker.go`**：worker 池生命周期正确——`Start` 注册空闲通道，`dispatcher` 执行完任务后重新入池，`Stop` 经 `quit` 通道退出并注销；`WaitGroup` 计数与 `Push` / worker 配对，未发现 goroutine 泄漏或计数错误。
- **`core/logger/es/doc.go`**：`IndexString` 返回静态 mapping 字符串，`CreateIndex` 仅在钩子构造期执行一次索引存在性检查，不在日志热路径。
- **`core/task/cron/parser.go`**：表达式 `Parse` 只在 `AddJob` 注册期调用一次，运行期不重复解析；`getField` / `getRange` / `getBits` 均为注册期一次性位集合计算。
- **`core/task/cron/spec.go`**：`SpecSchedule.Next` 为标准逐级进位匹配，循环上界由年/月/日/时/分/秒字段位数约束，不存在无界循环；`dayMatches` 为两次位与判断。
- **`core/task/cron/chain.go`**：`Then` 仅在注册期按包装链顺序包裹一次，运行期直接调用最外层包装。
- **`core/task/cron/constantdelay.go`**：`Every` 构造期固定间隔；`Next` 为一次 `t.Add(...)`，无分配热点。
- **`core/task/cron/option.go`**：纯选项闭包，注册期应用一次。
- **`core/task/nodepool.go`**：`CheckJobAvailable`（93-109 行）为 RLock 内纯内存一致性哈希查找；`waitingForHashRing`（133-148 行）按 `updateDuration`（默认 3s）周期刷新成员列表，频率低；`updateHashRing` 仅在成员变化时重建哈希环。
- **`core/task/driver/redisdriver.go`**：心跳按 `timeout/2` ticker 续期，`Stop` 取消 `runtimeCtx` 后心跳协程注销节点并退出；`scan` 不在锁内执行；未发现泄漏。
- **`core/task/driver/option.go` / `util.go`**：选项结构体与键前缀 / 节点 ID 生成，`GetNodeId` 仅在 `Init` 期调用一次。
- **`core/metric/metric.go`**：全部向量在包加载时一次性 `MustRegister`；`init()` 只写一次 `BuildInfoGauge`；运行期热点为各上报包装器（已列 P2）。
- **`core/trace/config.go`**：`Build` 为启动期构造 TracerProvider、批量导出器与 resource，不在请求路径。
- **`core/trace/carrier.go`**：`MetadataReaderWriter` 的 `Get` / `Set` / `Keys` 为 gRPC metadata 的薄适配，`Keys` 仅在传播器需要枚举键时使用，单次开销可忽略。
- **`core/trace/util.go`**：`ParseFullMethod` / `PeerAttr` 每次调用构造 1~2 个 `attribute.KeyValue`，属 Span 属性常规小分配；量级低于本章节 P2 门槛，未单列条目。
- **`core/trace/propagation.go` `Extract` 路径**：仅在收到 Jaeger 头时执行一次 split / hex 解析，无循环与锁；热点集中在 `Inject`（已列 P2）。
- **测试文件**：`core/metric/*_test.go`、`core/logger/queue/*_test.go` 均为桩测试 / 基准用例，不含生产实现。

---

## 七、附记：本章节阅读中发现的非性能问题（正确性 / 健壮性，不在性能条目内展开）

- `core/task/task.go:287` `reRunRecentJobs` 直接读 `d.jobs[jobName]`，未持有 `d.jobsRWMut`，与 `addJob` / `Remove` 的写锁存在 map 并发读写竞态（可能触发 fatal map read/write），建议补 RLock。
- `core/logger/queue/list_queue.go:87-99`：长度检查与入队分属两个临界区，`maxLen` 上限在并发 Push 下可能被瞬时越过。
- `core/logger/sql/exec.go:50`：`json.Marshal(entry.Data)` 的错误被 `_` 丢弃；`core/logger/sql/exec.go:17` 包级全局 `tableName` 在多个 SQL 日志实例并存时会被 `NewExec`（26 行）/ `NewExecWithURL`（35 行）互相覆盖。
- `core/task/driver/redisdriver.go:104`、`core/task/driver/rediszsetdriver.go:122`：`Del` 将 `rd.nodeID` 重复传入两次。
- `core/task/nodepool.go:168-180` `equalRing` 会就地排序入参 `a`（驱动返回的节点切片），可能影响调用方后续使用。
- `core/task/cron/cron.go:283` `startJob` 默认链路无 `Recover` 包装，单个任务 panic 可能导致进程崩溃（可靠性问题）。
