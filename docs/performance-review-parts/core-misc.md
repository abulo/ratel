# core/misc 包组性能审查报告

> 审查范围（仅本章节）：`core/app`、`core/call`、`core/component`、`core/constant`、`core/cycle`、`core/ecode`、`core/env`、`core/executor`、`core/goroutine`、`core/hooks`、`core/hostname`、`core/ratelimit`、`core/resource`、`core/singleton`、`core/stack` 及其下全部 `.go` 文件（含测试文件通读；性能条目仅针对生产代码）。
> 审查方式：逐文件人工通读源码，结合注释理解 API 语义；对疑似热点回到原文核对 `文件:行号` 与函数名；跨模块调用点仅用于校准触发频率与影响量级，不展开其他模块审查。
> 优先级定义：**P0** = 热路径上的明显性能缺陷，应优先修改；**P1** = 中等影响、性价比高的改进；**P2** = 低影响或改进空间有限的点。

---

## 一、条目清单（按子包分组；组内 P0 → P1 → P2 排序）

### core/resource

#### [P1] `core/resource/breaker.go:454` `RollingWindow.Add`：每个受保护请求都要获取一把排他写锁
- 问题：`Add` 是熔断器滚动窗口的唯一写入入口，内部 `rw.lock.Lock()`（455 行）排他加锁，临界区内执行 `updateOffset()`（457 行，含跨桶复位循环）与桶累加。熔断器热路径上，成功调用经 `googleBreaker.markSuccess()`（416–418 行）调 `Add(1)`，失败调用经 `markFailure()`（420–422 行）调 `Add(0)`，即**每一次受保护的 DB / RPC 请求都产生一次排他写锁**。同一把 `sync.RWMutex`（138 行）还被 `Reduce` 的读路径复用，写者与读者互斥，高并发下存在写者饥饿与锁排队。
- 影响评估：场景为接入熔断器的下游调用（仓内 MongoDB 等客户端以 `brk.DoWithAcceptable` 包装每次操作，见跨模块提示第 3 条）。量级判断：每请求 1 次 `Lock/Unlock` + 1 次 `updateOffset` 复位循环；在高 QPS 下游上，写锁串行化所有成功/失败计数更新，`pprof` 的 mutex profile 中会出现该 RWMutex 热点。当前每请求还要与 `Reduce` 的读锁（见下一条）在同一把锁上交替。
- 优化建议：
  1. 桶内 `Sum` / `Count` 改为 `atomic.AddInt64` 增量更新，`Add` 不再需要排他锁；仅在跨时间桶、需要复位旧桶时才短暂加锁。
  2. 或把“计数写入”与“窗口轮转复位”拆成两把锁：计数用原子操作无锁化，轮转复位走低频的独立锁。
  3. 复位循环（`updateOffset` 内 497–499 行）仅在跨桶时执行，可与计数路径解耦。

#### [P1] `core/resource/breaker.go:463` `RollingWindow.Reduce`：每次放行判定都在锁内全量遍历 40 个时间桶
- 问题：熔断放行路径 `googleBreaker.allow()`（380 行）→ `accept()`（363 行）→ `history()`（424–431 行）→ `stat.Reduce(...)`（463 行）。`Reduce` 在 `RLock` 下通过 `window.reduce`（522–526 行）遍历 `diff` 个桶并逐个回调累加 `Sum` / `Count`；桶数由常量 `buckets = 40`（29 行）固定。也就是说**每次放行判定都要在共享读锁上走完最多 40 个桶的循环**，且 `history()` 每次重新求和，窗口内没有任何累计缓存。
- 影响评估：场景同熔断器热路径（每次 DB/RPC 调用进入 `accept()` 时触发）。量级判断：每请求约 40 次桶指针解引用 + 两次 int64 累加，外加一次 RLock/RLock；与上一条的写锁在同一 RWMutex 上叠加，单请求窗口侧即“1 次读锁 + 最多 40 次迭代 + 1 次写锁”。桶数固定为 40，复杂度 O(buckets) 而非 O(1)，在熔断判定频繁的服务上属于稳定存在的固定开销。
- 优化建议：在窗口上维护累计值 `totalSum`、`totalCount`：`Add` 时原子累加；`resetBucket`（528–530 行）复位旧桶时同步从累计值扣减该桶的 `Sum` / `Count`；`history()` 直接读累计值，将放行判定降到 O(1)，`Reduce` 仅在需要按桶输出明细时保留遍历。

#### [P2] `core/resource/breaker.go:347` `Proba.TrueOnProba`：用互斥锁串行化 `rand.Rand` 随机数获取
- 问题：`Proba` 持有非并发安全的 `*rand.Rand`（154 行），`TrueOnProba` 每次取随机值都 `p.lock.Lock()`（348 行）包裹 `p.r.Float64()`（349 行）。
- 影响评估：该随机判定只在 `accept()` 算出 `dropRatio > 0`（369–373 行）即系统过载、需要按概率拒绝时才进入；正常流量下不触碰，锁竞争有限。
- 优化建议：改用并发安全的随机源（`math/rand/v2` 全局源，或基于 `atomic` 的 xorshift 无锁随机数），去掉互斥锁；保留现有锁也可接受。

#### [P2] `core/resource/breaker.go:490` `RollingWindow.updateOffset`：同一写锁临界区内重复获取当前时间
- 问题：`updateOffset` 先调 `span()`（491 行），`span()` 内部 `Since(rw.lastTime)`（482 行）会取一次 `time.Now()`；函数末尾 502 行又调用 `Now()` 取第二次当前时间用于时间对齐。两次系统时钟读取都在 `Add` 的排他写锁内。
- 影响评估：仅跨桶时进入 `span > 0` 分支；常态同桶内提前 return（492–494 行），重复取时钟的开销发生在桶轮转瞬间，频率低。
- 优化建议：由调用方传入 `now time.Time`（或在 `Add` 入口取一次时间向下传递），`span()` 与对齐逻辑复用同一时刻，消除锁内重复 `time.Now()`。

#### [P1] `core/resource/sourcemanager.go:50` `ResourceManager.GetResource`：缓存命中路径仍然先走单飞全局锁，再走资源表读锁，每调用双重加锁
- 问题：`GetResource` 无论缓存是否命中，第一步都进入 `manager.singleFlight.Do(key, fn)`（51 行）；`fn` 内部才做 `manager.lock.RLock()` 查 `resources` map（52–54 行）。已缓存资源的命中路径也要：单飞组件全局 `Mutex` 加/解锁 + 一次 map 哈希查找，再叠加资源表 `RWMutex` 读加/解锁 + 第二次 map 哈希查找，**没有缓存命中的快速旁路**。
- 影响评估：场景已核实——`stores/redis` 的 `getRedis(r)` 在每条命令执行前调用 `GetResource`（见跨模块提示第 2 条），即 Redis 命令热路径上每次都支付这双重锁与双重 map 查找；SQL / Mongo / ES 等客户端初始化路径同样经过该函数。量级判断：命中路径固定 2 把锁（单飞全局 Mutex + 资源表 RWMutex）、2 次 `map[string]io.Closer` 哈希；在命令级高频调用下属于可直接移除的固定开销。
- 优化建议：在进入单飞之前先做一次 `manager.lock.RLock()` 快速查表，命中立即返回已缓存资源；仅未命中时才释放读锁、进入 `singleFlight.Do` 执行创建与写回。单飞只承担“并发创建同一 key 合并”的职责，不参与缓存命中路径。

#### [P2] `core/resource/singleflight.go:56` `flightGroup.createCall`：全部 key 共用一把互斥锁，不同 key 的并发查找互相阻塞
- 问题：`flightGroup.calls` 由单一 `sync.Mutex`（25 行）保护，`createCall` 每次都先 `g.lock.Lock()`（57 行）再查 map；不同 key 的并发调用在这把全局锁上串行化。等待方命中后在锁外 `c.wg.Wait()`（60 行），这一点是正确的。
- 影响评估：当前仓内单飞主要用于资源管理器创建合并（见上一条），初始化阶段调用为主；多 key 并发同时缺失时才出现锁竞争。配合上一条修复（命中旁路）后，本锁只在创建瞬间短暂持有，影响进一步降低。
- 优化建议：按 key 哈希分片为 N 把锁（或 `sync.Map` + 每 key 的 `call` 结构），让不同 key 的查找互不阻塞；若保持现状，至少保证命中路径不再经过本组件。

### core/hostname

#### [P1] `core/hostname/hostname.go:6` `Hostname`：每次调用都发起 `os.Hostname()` 系统调用，主机名在进程生命周期内恒定却从不缓存
- 问题：`Hostname()` 函数体内每次执行 `os.Hostname()`（7 行），该调用内部发起系统调用并新分配返回字符串；失败时才回退 `"unknown"`。进程主机名在整个运行期间不会改变，但本包没有任何缓存。
- 影响评估：场景已核实——该函数被 gRPC 服务端两个 trace 拦截器、HTTP（gin/hertz）trace 中间件、SQL 拦截器及 Redis trace 钩子调用，其中多数位于**每请求 / 每命令**路径（具体位置见跨模块提示第 1 条）。量级判断：每请求一次主机名系统调用 + 一次字符串堆分配；系统调用本身单次不慢，但在高 QPS 下属于纯重复劳动，且该值可在启动时一次性确定。`core/env` 包已经在 `init()` 里缓存过一次 `hostName`（`core/env/const.go:115–119`），说明框架内已有缓存先例。
- 优化建议：包级 `sync.Once`（或 `init()`）里调用一次 `os.Hostname()` 存入包级变量，`Hostname()` 直接返回缓存值；可直接复用 `env.HostName()` 的缓存结果，避免第二次系统调用。

### core/call

#### [P1] `core/call/call.go:9` `Caller`：每次都做运行时栈回溯并重新解析函数名，结果按调用点恒定却无缓存
- 问题：`Caller(skip)` 每次调用 `runtime.Caller(skip)`（10 行）完成栈帧回溯，再 `runtime.FuncForPC(pc).Name()`（11 行）解析并分配新的函数名字符串；`file`、`lineNo` 同样每次重新获取。相同调用点的 PC 在运行期间恒定，函数名解析结果可缓存，但当前实现没有 `sync.Map` 之类的缓存层。
- 影响评估：场景已核实——Redis `ProcessHook` / `ProcessPipelineHook` 在**每条命令**路径上无条件调用 `call.Caller(11)`（`stores/redis/hook.go:60`、`:92`，未受 trace 开关门控）；MongoDB 集合操作在 `EnableTrace` 打开时调用 `call.Caller(6)`；SQL、ES 拦截器也有调用点（见跨模块提示第 1 条）。量级判断：`runtime.Caller` 需要回溯指定栈帧、`FuncForPC().Name()` 需要查符号表并分配字符串，属于运行时偏重的操作；在无条件调用的 Redis 命令路径上，这部分开销叠加在每条命令上。
- 优化建议：
  1. 包内增加 `sync.Map`，以 `pc` 为键缓存 `(fn, file, lineNo)` 三元组，同 PC 第二次起直接命中缓存，避免重复符号解析与字符串分配。
  2. 调用侧（其他模块）应在 trace 关闭时跳过本调用；本包内无法判断门控，属跨模块配合项。
  3. 若仅需要函数名而不需要文件/行号，可在 API 上拆分轻量版本。

### core/goroutine

#### [P2] `core/goroutine/goroutine.go:37` `RestrictParallel`：先为每个任务创建 goroutine，再在信号量通道上排队，受限的是执行并发而不是 goroutine 数量
- 问题：返回的闭包对每个 `fn` 立即 `go func(fn func())`（43 行），协程启动后才在 45 行 `channel <- struct{}{}` 获取容量为 `restrict` 的信号量；信号量未就绪时协程阻塞在通道发送上，直到被调度执行。任务数 N 会立即产生 N 个 goroutine，其中 `N - restrict` 个处于排队阻塞状态。
- 影响评估：场景为一次性并发任务编排，任务规模通常有限；问题是排队 goroutine 提前占用初始栈（Go 协程 2KB 起步）与调度队列，任务列表很大时比 worker 池方案多一份协程创建/排队开销。本包 `parallel.go:50–79` 的 `RestrictParallelWithErrorChan` 已经采用 worker 池模式，可作为对照。
- 优化建议：改为固定 `restrict` 个 worker 协程从 `jobs` 通道取任务，或在 `go` 之前先获取信号量；与 `RestrictParallelWithErrorChan` 的实现对齐。

#### [P2] `core/goroutine/goroutine.go:57` `GoDirect`：每次调用新建依赖注入容器并通过反射完成参数映射与调用
- 问题：`GoDirect` 每次执行 `inject.New()`（58 行），循环 `inj.Map(arg)`（59–61 行），随后在协程内 `inj.Invoke(fn)`（72 行）完成反射式参数解析与调用；容器不复用，反射每次重新走类型匹配。
- 影响评估：全仓检索确认当前**没有生产代码调用 `GoDirect`**（仅本包定义），暂无实际热点；若未来在任务回调等高频场景使用，容器构造与反射 `Invoke` 的开销会叠加在每次调用上。
- 优化建议：保持其“仅启动期 / 低频使用”的定位；热路径一律使用 `Go` / `SafeGo` / `Parallel` 等无反射变体；如确认不再使用可考虑移除该 API。

### core/stack

#### [P2] `core/stack/stack.go:33` `DeferStack.Clean`：执行用户闭包的整个过程持有写锁，闭包再次入栈存在重入死锁风险
- 问题：`Clean` 在 34 行 `ds.mu.Lock()`，`defer Unlock` 覆盖整个循环；36–38 行在持锁状态下逆序调用每个用户闭包 `v()`。若闭包内部回调 `Push`（24 行，同样要获取写锁），将在同一 goroutine 上自死锁；同时多个并发 `Push` 会被整个清理过程阻塞。
- 影响评估：`DeferStack` 用于注册 / 退出清理，`Clean` 属于生命周期收尾路径，不是请求热路径；问题主要是锁粒度过粗与潜在重入死锁，性能影响在并发清理场景下才显现。另注：`Clean` 执行后不清空 `fns`（注释 32 行已说明），重复调用会重复执行，属既有设计。
- 优化建议：临界区内把 `ds.fns` 拷贝到本地切片并置空（或仅取快照），释放写锁后再逆序执行闭包；用户代码在锁外运行，避免阻塞与重入死锁。

### core/cycle

#### [P2] `core/cycle/cycle.go:48` `Cycle.Done`：等待全部任务结束的协程在持有互斥锁的状态下阻塞于 `wg.Wait()`
- 问题：`Done()` 的一次性等待协程在 51 行 `c.mu.Lock()`，`defer Unlock` 包裹 53 行 `c.wg.Wait()`——锁在整个等待期间一直被持有。等待期间任何并发 `Run()`（35 行定义，36 行 `c.mu.Lock()`）与 `Close()`（67 行定义，68 行 `c.mu.Lock()`）都会阻塞在同一把锁上，直到全部任务退出。
- 影响评估：`Cycle` 用于后台任务生命周期管理，`Done` 通常在停止阶段调用；按框架使用时序（先 `Run` 登记、后 `Done` 等待）常态下不会在等待期间再 `Run`，所以实际竞争有限；属于锁使用方式不规范，极端时序下会造成不必要的串行阻塞。
- 优化建议：锁内仅完成“登记等待意图 + 启动协程”，协程内不要持有 `c.mu` 执行 `wg.Wait()`；可借助 `sync.WaitGroup` 直接派生完成通知通道，或将锁保护范围缩小到状态变量。

### core/ecode

#### [P2] `core/ecode/code.go:64` `ExtractCodes`：非 nil 错误每次新建状态包装对象、`spb.Status` 与空 Details 切片
- 问题：错误非空时，69–75 行每次都分配一个 `spbStatus` 包装对象、一个 `spb.Status` 结构体，并用 `make([]*any.Any, 0)` 新建空 Details 切片；`status.FromError(e)` 对非 gRPC 错误会回退为 `codes.Unknown`。`e == nil` 时 65–67 行直接返回共享 `OK`，成功路径无分配，这一点是正确的。
- 影响评估：该函数被 gRPC 服务端 / 客户端指标与日志拦截器调用，成功请求走 nil 快路径；仅**错误返回**路径产生上述 3 处小对象分配。错误路径本身不在主吞吐面上，分配量级小，但在错误率偏高的服务上会累积临时对象与 GC 压力。
- 优化建议：Details 改为惰性初始化（访问时才创建切片）；可考虑为常见错误码复用预分配状态实例，减少错误路径临时对象。

#### [P2] `core/ecode/status.go:72` `spbStatus.String`：使用 `json.Marshal` 反射序列化整个状态结构
- 问题：`String()` 直接 `json.Marshal(s)`（73 行），依赖反射遍历结构体字段生成 JSON，再把字节切片转成字符串；错误忽略。
- 影响评估：全仓未检索到对 `spbStatus` 调用 `.String()` 的生产点（需实测确认是否在日志格式化隐式触发，如 `%v` / `%s` 打印状态对象）。若未来在错误日志热路径上打印状态，反射 JSON 序列化的开销明显高于手写拼接。
- 优化建议：改为手写拼接 `"code=..., message=..."` 的字符串格式，去掉反射；保留 JSON 输出作为独立的 `JSON()` 方法。

#### [P2] `core/ecode/status.go:40` `spbStatus.GetMessage`：每个 ext 键值对都执行两次 `fmt.Sprintf("%v")` 反射格式化
- 问题：循环内 52 行格式化键、54 行格式化值，`fmt.Sprintf("%v", v)` 对每个 `any` 走类型判断与反射格式化；无 ext 时仅走 `buf.WriteString`，开销很低。
- 影响评估：本方法在指标标签路径使用时通常不传 ext（`server/xgrpc/interceptor.go:64` 调用 `code.GetMessage()` 无参）；带 ext 的调用多在错误日志场景，频率低。
- 优化建议：对 `string` / 数值等常见类型直接写入 `bytes.Buffer`，仅对其他类型回退 `fmt.Sprintf`，减少反射判断。

### core/singleton

#### [P2] `core/singleton/singleton.go:12` `genkey`：通过 `cast.ToString` 把模块编号转字符串，每次单例读写都做一次类型反射式转换
- 问题：`genkey` 用 `cast.ToString(int(module))`（13 行）把枚举编号转字符串后与 `key` 拼接；`spf13/cast` 的 `ToString` 对整数走类型分支与字符串构造，`Load`（17–19 行）、`Store`（22–24 行）每次都重新拼一次完整键。
- 影响评估：单例主要在客户端 / 注册中心初始化阶段按模块 + 键查询与写入，不属于每请求热路径；开销为一次小字符串构造，量级低。
- 优化建议：改用 `strconv.Itoa(int(module))` 拼接键，去掉 `cast` 依赖；固定模块编号的键可在注册时预生成。

### core/ratelimit

#### [P2] `core/ratelimit/ratelimit.go:137` `Bucket.Take` 系列：每次取用都持有 `sync.Mutex`，桶内负令牌记账使无锁改造受限（当前仓内无生产调用方，需实测）
- 问题：`Take`（137–142 行）、`TakeMaxDuration`（146–150 行）、`TakeAvailable`（154–158 行）、`Available`（179–188 行）每次都对 `tb.mu` 加锁，临界区为 `currentTick` 计算、令牌调整与扣减的算术操作；`clock.Now()` 在锁外获取（140 行），这一点正确。
- 影响评估：全仓检索确认该令牌桶当前**没有非测试生产代码引用**，现有锁竞争程度无从在线上验证；一旦在限流中间件等每请求路径接入，单桶高并发下该 Mutex 会成为热点。由于桶需要维护可负的 `availableTokens`（207、220 行）与 `latestTick` 一致性，简单换成 atomic 并不能直接保证记账正确。
- 优化建议：接入前在目标并发模型下做 benchmark / mutex profile；若确为热点，可参考 `golang.org/x/time/rate` 的令牌桶状态机评估 atomic 化，或为高频只读的 `Available` 分离轻量读路径；当前保持实现亦可接受。

---

## 二、条目统计

| 子包 | P0 | P1 | P2 | 小计 |
|---|---|---|---|---|
| core/resource | 0 | 3（`RollingWindow.Add`、`RollingWindow.Reduce`、`ResourceManager.GetResource`） | 3（`Proba.TrueOnProba`、`RollingWindow.updateOffset`、`flightGroup.createCall`） | 6 |
| core/hostname | 0 | 1（`Hostname`） | 0 | 1 |
| core/call | 0 | 1（`Caller`） | 0 | 1 |
| core/goroutine | 0 | 0 | 2（`RestrictParallel`、`GoDirect`） | 2 |
| core/stack | 0 | 0 | 1（`DeferStack.Clean`） | 1 |
| core/cycle | 0 | 0 | 1（`Cycle.Done`） | 1 |
| core/ecode | 0 | 0 | 3（`ExtractCodes`、`spbStatus.String`、`spbStatus.GetMessage`） | 3 |
| core/singleton | 0 | 0 | 1（`genkey`） | 1 |
| core/ratelimit | 0 | 0 | 1（`Bucket.Take` 系列） | 1 |
| core/app | 0 | 0 | 0 | 0 |
| core/component | 0 | 0 | 0 | 0 |
| core/constant | 0 | 0 | 0 | 0 |
| core/env | 0 | 0 | 0 | 0 |
| core/executor | 0 | 0 | 0 | 0 |
| core/hooks | 0 | 0 | 0 | 0 |
| **合计** | **0** | **5** | **12** | **17** |

---

## 三、高价值候选列表（Top 10，供全局 Top 10 评选）

| # | 位置 | 一句话问题 | 预期收益 |
|---|---|---|---|
| 1 | `core/hostname/hostname.go:6` `Hostname` | 每请求 / 每命令都重复发起主机名系统调用并分配字符串，应启动期缓存 | 消除所有 trace 中间件与 DB 钩子热路径上的一次系统调用 + 堆字符串分配；改动极小、风险极低 |
| 2 | `core/resource/sourcemanager.go:50` `GetResource` | 缓存命中仍走单飞全局锁 + 资源表读锁双重加锁，缺少命中快速旁路 | Redis 等命令级调用每次省掉 2 把锁与 2 次 map 哈希查找；锁竞争与固定指令开销同步下降 |
| 3 | `core/resource/breaker.go:463` `RollingWindow.Reduce` | 每次放行判定在锁内全量遍历 40 个桶重新求和 | 熔断判定从 O(40) 桶循环降为 O(1) 累计值读取，降低每个受保护请求的固定开销 |
| 4 | `core/resource/breaker.go:454` `RollingWindow.Add` | 每个成功 / 失败计数都走排他写锁，与读路径在同一 RWMutex 上竞争 | 计数写入 atomic 化后消除每请求一次排他锁排队，缓解高 QPS 下游上的写锁竞争 |
| 5 | `core/call/call.go:9` `Caller` | 每条 Redis 命令无条件做栈回溯 + 符号解析，PC 恒定结果却无缓存 | 同调用点二次起命中 PC 缓存，省去运行时栈回溯与函数名字符串分配；配合调用侧门控收益更大 |
| 6 | `core/resource/singleflight.go:56` `flightGroup.createCall` | 全部 key 共用一把互斥锁，不同 key 并发查找互相阻塞 | 分片后多 key 并发创建不再串行化；与第 2 条配合，把单飞锁完全移出缓存命中路径 |
| 7 | `core/cycle/cycle.go:48` `Cycle.Done` | 等待协程持锁阻塞于 `wg.Wait()`，阻塞期间其他 `Run` / `Close` 全部排队 | 锁粒度修正，消除停止阶段潜在的锁等待；生命周期收尾时序更稳健 |
| 8 | `core/ecode/code.go:64` `ExtractCodes` | 错误路径每次新建状态对象、`Status` 与空 Details 切片 | 错误率偏高的服务上减少临时小对象分配与 GC 压力；改动局限于错误路径 |
| 9 | `core/goroutine/goroutine.go:37` `RestrictParallel` | 任务 goroutine 先全部创建再抢信号量，并发受限而协程数不受限 | 大任务批量场景下减少排队协程的栈与调度开销；与包内已有 worker 池实现对齐 |
| 10 | `core/ecode/status.go:72` `spbStatus.String` | `json.Marshal` 反射序列化状态（是否处于热路径需实测） | 若确有日志调用点，手写拼接可去掉一次反射序列化；无调用点则维持现状 |

---

## 四、拿不准、需实测验证的条目清单

1. **`core/ratelimit/ratelimit.go:137` `Bucket.Take` 系列锁竞争**：仓内当前无生产调用方，Mutex 在单桶高并发下的竞争程度只能等接入限流中间件后，用 mutex profile / benchmark 实测；atomic 化改造是否值得需要数据支撑。
2. **`core/ecode/status.go:72` `spbStatus.String` 是否处于热路径**：未检索到显式 `.String()` 调用，需确认是否有日志以 `%v` / `%s` 隐式触发；若仅在错误展示路径调用，优先级维持 P2。
3. **`core/call/call.go:9` `Caller` 各调用点的门控状态**：Redis 钩子两处为无条件调用；MongoDB 各操作是否走 `Caller` 取决于 `EnableTrace` 运行时取值，需在目标服务配置下用 trace 开关对照 benchmark，量化开关关闭后的节省。
4. **`core/resource/breaker.go:454` / `:463` 熔断窗口锁与遍历的实际占比**：需在接入熔断器的下游客户端上跑 `pprof -mutexprofile` 与 CPU profile，确认 RWMutex 排队与 40 桶遍历在整体请求耗时中的占比，再决定 atomic 化与累计值缓存的投入力度。
5. **`core/resource/sourcemanager.go:50` 双重锁在 Redis 命令路径上的实际占比**：建议对 `getRedis` 命中路径做 before/after benchmark，量化两次锁 + 两次 map 查找的固定成本；若命令实现内还有其他锁，需要合并观察避免重复归因。
6. **`core/resource/breaker.go:347` `Proba.TrueOnProba` 锁**：仅在过载（`dropRatio > 0`）时进入，正常流量不竞争；是否替换随机源取决于线上过载触发频率。
7. **`core/ecode/code.go:64` `ExtractCodes` 错误路径分配量级**：分配是否形成可观 GC 压力，取决于服务错误率与错误路径 QPS，需结合运行时堆 profile 判断。

---

## 五、跨模块提示（超出本章节范围，仅登记线索，不展开审查）

1. **每请求 / 每命令重复取主机名与调用栈**（对应本报告第 1、5 条的调用侧）：
   - `stores/redis/hook.go:57`、`:89` 每条命令路径调用 `hostname.Hostname()`；`:60`、`:92` 每条命令无条件调用 `call.Caller(11)`，未见 trace 开关门控。
   - `server/xgrpc/interceptor.go:93`、`:137`，`server/xgin/middleware.go:52`，`server/xhertz/middleware.go:75`，`stores/sql/interceptor.go:94` 在每请求 trace 路径调用 `hostname.Hostname()`。
   - 建议：调用侧在 trace / 指标关闭时跳过这些采集；主机名统一取缓存值。
2. **`stores/redis/redis.go:199` `getRedis(r)` 在每条命令执行前调用 `clientManager.GetResource`**（201、217、233、253 行四个变体），是 `ResourceManager.GetResource` 双重锁问题（本报告第 2 条）在命令热路径上的直接触发点。
3. **`stores/mongodb/collection.go` 各写 / 读操作经 `c.brk.DoWithAcceptable(...)` 包装**（如 `:119`），是熔断器 `Add` / `Reduce` 每请求锁与桶遍历（本报告第 3、4 条）的触发路径；`stores/mongodb/clientmanager.go:177`、`:201`、`:211` 的 `GetBreaker` 为模型创建期一次调用，不在热路径。
4. **`core/goroutine/init.go:29` 的 `util.FunctionName(fn)`** 基于 `reflect.ValueOf(...).Pointer()` + `runtime.FuncForPC`（实现位于 `util/util.go:121–123`），仅在 panic 恢复分支（冷路径）执行，未列为性能条目；若未来把 `try` 推广到高频无 panic 路径，需重新评估。
5. **`core/ecode/code.go:34` `_codes sync.Map`** 在 `add`（58 行）中只写不读，状态对象常驻进程内存；属于启动期登记缓存，非热路径问题，仅登记备查。
6. **`core/app/app.go:134` `Schedule`、`:140` `Job`** 直接写共享 `workers` 切片与 `jobs` map，未像 `Serve`（126–131 行）那样使用 `app.smu`；均为启动期注册调用，不构成性能问题，属并发一致性线索，转正确性审查组参考。

---

## 六、已通读但未发现性能问题的包 / 导出 API（覆盖说明）

- **`core/app`**：`Application` 为启动、信号监听、停止生命周期编排；`Run` / `Stop` / `GracefulStop` / `startServers` / `startWorkers` / `startJobs` 均为一次性生命周期流程，不在请求热路径。`initialize` 使用 `sync.Once`，停止使用 `stopOnce`，使用方式正确。`exitHandler` 中 `debug.FreeOSMemory()`（348 行）在信号等待前执行一次，属一次性操作，未列为条目。
- **`core/component`**：仅定义 `Component` / `Manager` / `Builder` 接口与 `ComponentFunc` 适配器，无可执行实现代码。
- **`core/constant`**：纯常量、环境变量名与模块枚举；`service.go` 的 `ServiceKind.String()`（28–33 行）为 4 项 map 查找，开销可忽略，不凑数列条目。
- **`core/env`**：`InitEnv` 启动期读取环境变量；各 `AppXxx()` getter 为全局字符串直接返回，无分配、无锁；`AppID()` 空串分支判断代价可忽略。`Set*` 写操作均在配置期完成。
- **`core/executor`**：`Register` / `Run` / `Stop` / `GracefulStop` 为启动 / 停止一次性流程，`sync.Map` 注册表读写频率低；无请求热路径。
- **`core/hooks`**：`Register` 为启动期追加，`Do` 仅在生命周期固定阶段（BeforeRun / BeforeStop 等）各触发一次；`fmt.Sprintf` 与彩色输出不在请求路径。
- **`core/resource/atomicduration.go`**：`AtomicDuration` 全部经 `atomic.LoadInt64` / `StoreInt64` / `CompareAndSwapInt64` 访问，实现正确，无锁竞争问题。
- **`core/resource/batcherror.go`**：`BatchError` 仅在批量关闭 / 错误聚合场景使用；`bytes.Buffer` 拼接方式合适，无热路径分配问题。
- **`core/resource/nopbreaker.go`**：空操作熔断器，所有方法直接透传，无额外开销。
- **`core/resource/breakers.go`**：`GetBreaker` 使用 RLock 快查 + Lock 双检锁（41–58 行），写法正确；共享熔断器 map 仅在按名获取时短临界区访问，未发现热点。
- **`core/stack`**：`defer.go` 仅为全局 `DeferStack` 的 `Register` / `Clean` 转发；除 `Clean` 锁粒度条目外，`Push` 为启动 / 注册期调用，无热路径问题。
- **`core/ecode/unified.go`**：纯日志类别与模块名字符串常量。
- **测试文件**：`core/cycle/cycle_test.go`、`core/ratelimit/ratelimit_test.go`、`core/ratelimit/reader_test.go`、`core/stack/stack_test.go`、`core/stack/defer_test.go` 已通读，不含生产实现。
