# 性能审查报告：stores/redis、stores/proxy、stores/session

> 审查范围（仅以下三个包目录及其 `.go` 文件）：
>
> - `stores/redis/`（25 个 `.go`：核心装配 `redis.go` / `options.go` / `command.go` / `hook.go` / `utils.go` / `client.go`，以及 18 个 `*_commands.go` 命令封装文件，共 457 个 `(r *Client)` 命令方法）
> - `stores/proxy/`（6 个 `.go`：`proxy.go`、`redis.go`、`mongodb.go`、`elasticsearch.go`、`sql.go`、`doc.go`）
> - `stores/session/`（3 个 `.go`：`driver.go`、`driver_test.go`、`doc.go`）
>
> 审查方式：静态通读 + 行号回读核对；熔断器、资源管理器、随机数工具等被调用方实现位于本范围之外（`core/resource/`、`util/`），正文条目锚定在范围内的调用点，根因实现统一在文末「跨模块提示」中单列，不在本章节展开审查。
>
> 总体结论：`stores/redis` 的 18 个命令封装文件经逐行核对，**全部是同一模式的薄包装**（`r.brk.DoWithAcceptable(func(){ conn, err := getRedis(r); ... conn.Xxx(getCtx(ctx), ...).Result() })`），封装层自身没有自定义循环、JSON/反射编解码、正则或额外 buffer；问题集中在**每条命令必经的外围链路**：熔断器统计、资源管理器/SingleFlight 查找、追踪钩子分配，以及 `stores/session` 的整文档读-改-写模型、`stores/proxy` 的随机选库实现。

---

## 一、stores/redis

#### [P0] `stores/redis/redis.go:199` `getClient`（同 `getCluster:215`、`getFailover:231`、`getRing:247`，经 `getRedis:180` 被每条命令调用）：每条命令都经过全局 ResourceManager + SingleFlight 锁查找，底层实例未在 `Client` 上缓存

- 问题：所有命令方法（共 457 个，覆盖 `stores/redis/*_commands.go`）执行前先调用 `getRedis(r)`（`redis.go:180-195`），再进入 `getClient/getCluster/getFailover/getRing`，后者统一走 `clientManager.GetResource(driverName, create)`（`redis.go:201`、`217`、`233`、`253`）。`Client` 结构体（`options.go:18-66`）只保存配置与 `brk`，**没有任何字段缓存解析后的 `RedisNode` 实例**；因此即使底层 `*redis.Client` / `*ClusterClient` / `*Ring` 早已创建并驻留在全局 `clientManager` 中，每条命令仍然要重复走一遍：
  1. `singleFlight.Do(key, fn)`（`core/resource/singleflight.go:36-44`）：进入 `createCall` 先取全局 `g.lock` **互斥锁**（`singleflight.go:57`）查 `calls` map，结束后 `makeCall` 内再次取同一把锁删除 key（`singleflight.go:74-76`）；
  2. `fn` 内部 `manager.lock.RLock()` 查 `resources` map（`core/resource/sourcemanager.go:52-54`），未命中才创建；命中路径也要完成一次 RLock/RUnlock。
  即：**每条命令、每个客户端实例，共享全局 SingleFlight 互斥锁，叠加资源表读写锁**，目的仅仅是取出一个本可以在 `Client` 上常驻的连接实例。
- 影响评估：高频短命令（`Get`/`Set`/`Incr`/`HGet` 等，见 `string_commands.go:51`、`hash_commands.go:37`、`set_commands.go:10`）是框架最热的数据访问路径。每条命令在网络往返之外多付出：1 次 SingleFlight 互斥锁临界区（创建路径锁 + 清理路径锁）+ 1 次资源表 RLock 查找 + 2 次 map 哈希 + 类型断言。锁是**跨全部 Redis 客户端实例共享**的同一把（`redis.go:197` 的包级 `clientManager`），高并发下所有实例的命令在这把锁上排队；量级判断：单次临界区为微秒级，但在万级 QPS、多个 `*Client` 并存时，互斥锁串行化带来的等待与上下文切换会叠加到每条命令的尾延迟上；具体占用比例需 pprof 实测（见文末「需实测验证清单」第 1、2 项）。
- 优化建议：
  1. 在 `Client` 上增加缓存字段（如 `node RedisNode`，配 `sync.Once` 或首检 + `atomic.Pointer[RedisNode]` 保护），`getRedis` 成功解析后赋值；后续命令直接读取缓存字段，跳过 `clientManager.GetResource`。
  2. `GetResource` 命中路径可做无锁 fast path：先 `atomic.LoadPointer` / 无锁读 map，未命中再进入 SingleFlight 创建流程（该改造在 `core/resource`，见跨模块提示）。
  3. `driverName` 与解析结果绑定，缓存实例时一并缓存，避免随每条命令重建（见 P2 条目「driverName 每条命令重建」）。

#### [P1] `stores/redis/string_commands.go:52` `Get`（同一包装模式遍布全部 `*_commands.go`，代表性热点读路径）：每条命令强制经过熔断器，成功/失败都要做一次全桶历史扫描（读锁 + 遍历 40 个桶）与一次写锁计数

- 问题：每个命令方法都以 `r.brk.DoWithAcceptable(func() error { ... }, acceptable)` 包裹（如 `string_commands.go:52-59`、`command.go:11-21`、`hash_commands.go:38-46`）。`DoWithAcceptable` → `doReq`（`core/resource/breaker.go:390-414`）在**执行业务命令前后各做一次窗口统计**：
  - 进入时 `b.accept()` → `b.history()`（`breaker.go:363-378`、`424-431`）→ `stat.Reduce(...)`（`breaker.go:463-479`）：取 `stat.lock.RLock()`，遍历**全部 40 个时间桶**（`breaker.go:29` `buckets = 40`）累加 `Sum`/`Count`；
  - 命令返回后，`acceptable(err)` 为真走 `b.markSuccess()` → `b.stat.Add(1)`（`breaker.go:416-418`、`454-459`）：取 `stat.lock.Lock()` 写锁，`updateOffset()`（`breaker.go:490-505`）计算 `span` 并循环重置已过期桶，过程中多次调用 `Now()`/`Since()`（`breaker.go:481-488`、`502`）。
  也就是说，**每一条成功的 Redis 命令**都要完成：一次 RLock + 40 桶遍历求和，一次 Lock + 过期桶重置循环 + 多次时钟读取；熔断器实现位于 `core/resource`（跨模块），但「每条命令都必须经过熔断器」这一决策在 `stores/redis` 范围内。
- 影响评估：`Get`/`Set` 这类命令本身服务端处理在亚毫秒到毫秒级；外围每命令 2 把锁（RWMutex 读锁 + 写锁）+ 固定 40 桶遍历 + 多次 `time.Since`/`time.Now` 调用，在高并发下 `stat.lock` 成为同一 `Client` 上所有命令的共享竞争点；`Reduce` 的 40 桶遍历是固定 O(40) 工作，与命令频率成正比。是否构成尾延迟主因需实测对比「直接调用底层 client、不套熔断器」的 benchmark（见需实测清单第 1 项）；但该开销对**每一条**命令都成立，属于热点路径上的结构性叠加。
- 优化建议：
  1. 范围可控方案：为高频只读命令提供不经过熔断器的 fast-path 包装（熔断器主要应对下游故障熔断，健康的单节点短命令可不纳入失败统计）；或熔断器开关按 `Client` 配置关闭。
  2. 跨模块根因改造（供 `core/resource` 审查组参考）：`history()` 改为增量维护 accepts/total 计数（`Add` 时同步原子累加，`Reduce` 仅在需要判断拒绝比例时读取）；`Add` 内 `updateOffset` 的过期桶重置改为惰性/后台异步；`stat` 锁降为原子计数 + 分片窗口，避免每命令读写双锁。

#### [P1] `stores/redis/hook.go:53` `OpenTraceHook.ProcessHook`（同 `ProcessPipelineHook:85`）：闭包共享 `attrs` 切片被每条命令反复 `append`，属性跨命令累积、底层数组只增不减；且每条命令做一次 11 层栈回溯 + 全量命令字符串拼接

- 问题：`ProcessHook` 在钩子注册时创建外层 `attrs := []attribute.KeyValue{...}`（`hook.go:55-58`），返回的闭包**捕获该切片**；每次执行命令时：
  - `hook.go:61-67` 连续 `attrs = append(attrs, semconv.CodeFunction(fn), semconv.CodeFilepath(file), semconv.CodeLineNumber(line), semconv.DBStatement(cmdString))`——append 写入的是闭包**同一个底层数组**。第一条命令后切片扩容，第二条命令在已含 2~7 个属性的切片上继续追加，导致：(a) 后续每个 span 携带前面命令残留的源码位置与旧命令语句；(b) 底层数组随命令种类/参数持续增长、永不回收，运行时间越久内存占用越高。这同时是正确性（错误属性）与内存增长问题。
  - `hook.go:60` `fn, file, line := call.Caller(11)` → `core/call/call.go:9-12` 每次命令执行 `runtime.Caller(11)` 回溯 11 层栈帧，再 `runtime.FuncForPC(pc).Name()` 取函数全名；`hook.go:66` `rediscmd.CmdString(cmd)` 把整条命令（key 与全部参数）格式化成字符串作为 `DBStatement` 属性写入，参数越长分配越大。
  - `ProcessPipelineHook` 同构：`hook.go:92` `call.Caller(11)`、`:93-100` 对共享 `attrs` 反复 append、`:99` `rediscmd.CmdsString(cmds)` 汇总整条管道命令串。
  - 另：`hook.go:68-73` / `101-106` 每条命令 `metadata.FromIncomingContext` 后 `md.Copy()` 复制 metadata。
- 影响评估：仅在 `EnableTrace` 开启（`redis.go:204-206` 挂载 `OpenTraceHook`）时触发；一旦开启，每条 Redis 命令叠加一次栈回溯（`runtime.Caller` 在深栈下是明确的 CPU 热点）、一次命令串格式化分配、metadata 拷贝，以及 `attrs` 累积带来的单调内存增长。`attrs` 累积问题与命令数量正相关，长跑服务会持续泄漏式占用；命令体较大（大 key、多参数 MSet/HSet）时 `CmdString` 分配与命令体量同阶。
- 优化建议：
  1. 修复累积：在返回闭包内部、每条命令局部 `attrs := make([]attribute.KeyValue, 0, 4)` 构建本次 span 属性，禁止复用外层切片；基础属性（`RPCSystem`、`HostName`）可在闭包外固定一份，命令级属性局部追加。
  2. 栈回溯：源码位置（function/file/line）属性对线上排障价值有限，可配置关闭；确需保留时，仅对错误命令采样记录，避免成功命令逐条 `runtime.Caller`。
  3. `DBStatement` 只记录命令名（`cmd.FullName()`），参数串做长度截断（如 256/512 字节）或按开关关闭；`CmdsString` 同理。
  4. `md.Copy()` 仅在确实需要透传 metadata 时执行；trace 链路优先使用 `trace.ContextWithRemoteSpanContext` 机制，避免每命令复制整份 metadata。

#### [P2] `stores/redis/redis.go:200` `getClient`（同 `getCluster:216`、`getFailover:232`、`getRing:249-252`）：缓存键 `driverName` 在每条命令上重新拼接/Join/map 遍历生成

- 问题：配置在 `NewRedisClient` 阶段即固定，但 `driverName` 在每次 `getClient/getCluster/getFailover/getRing` 内现算：
  - `redis.go:200`：`r.Addr + "@" + r.Name`（两次字符串拼接 + 一次合并分配）；
  - `redis.go:216`：`util.Implode(";", r.Hosts) + "@" + r.Name`——`util/array.go:209-219` 的 `Implode` 内部每次新建 `bytes.Buffer`、循环 `WriteString`、最后 `buf.String()` 整体拷贝；
  - `redis.go:232`：`r.MasterName + "://" + util.Implode(";", r.Hosts) + "@" + r.Name`；
  - `redis.go:249-252`：`for k, v := range r.Addrs { driverName += k + ":" + v + ";" }`，每次命令遍历整张分片 map，map 迭代顺序不确定，且 `+=` 逐次触发字符串重新分配。
- 影响评估：属于每命令的小字符串分配叠加项：单节点路径分配最小（2 段拼接）；集群/哨兵按地址数量线性放大；ring 路径是 O(分片数) 遍历 + O(累计长度) 的多次字符串重分配。相比 P0 条目，单次量级小，但与命令 QPS 成正比，属于可低成本消除的分配热点。
- 优化建议：在 `NewRedisClient` 或首次实例创建成功后，把 `driverName` 作为 `Client` 的缓存字段生成一次；ring 模式初始化阶段用预分配长度的 `strings.Builder` 生成；配合 P0 条目把实例与 key 一起缓存，命令路径上不再出现任何拼接。

#### [P2] `stores/redis/scripting_commands.go:10` `Eval`（同 `EvalRO:36`、`EvalSha:23`、`EvalShaRO:49`）：每次执行都把完整 Lua 脚本源码随命令发送，未使用 go-redis 的 `Script`（EVALSHA + NOSCRIPT 自动回退）机制

- 问题：`Eval`（`scripting_commands.go:10-20`）与 `EvalRO`（`:36-46`）直接把 `script string` 透传给 `conn.Eval`/`conn.EvalRO`，每次调用都在命令体中携带整段 Lua 源码；`EvalSha` 虽然存在（`:23-33`），但上层需要调用方自行管理 SHA1、自行处理 NOSCRIPT 重载，封装层没有提供 `go-redis` `*redis.Script` 那种「首次 EVAL、之后 EVALSHA、收到 NOSCRIPT 自动重新 SCRIPT LOAD/EVAL」的标准封装。
- 影响评估：对**频繁复用同一段脚本**的业务（限流、扣减、分布式锁、复合校验等），每次网络往返都重复传输脚本体，服务端也需反复解析；脚本越大、调用频次越高，网络带宽与服务端解析开销越明显。短小、一次性脚本影响有限；是否构成瓶颈取决于实际脚本大小与复用频次，需结合业务用法实测（需实测清单第 6 项）。
- 优化建议：
  1. 增加基于 `redis.NewScript(src)` 的封装（如 `NewScriptRunner(script)` 返回持有 `*redis.Script` 的句柄，`Run(ctx, keys, args...)` 内部走 EVALSHA/`EvalShaRO` 自动回退）。
  2. 对启动期固定的脚本，在客户端初始化时 `ScriptLoad` 预热，运行期一律走 `EvalSha`。
  3. 现有 `Eval`/`EvalRO` 保留给一次性脚本使用。

#### [P2] `stores/redis/options.go:78` `GetRingClientConfig`（同 `GetFailoverClientConfig:98-99`、`GetClusterClientConfig:118-119`、`GetClientConfig:141-142`）：`IdleCheckFrequency` 被错位映射到 `ConnMaxLifetime`；未显式配置的超时以零值直接透传；`PoolSize` 默认 64 且集群按节点计数

- 问题：
  1. 语义错位（`options.go:79`、`:99`、`:119`、`:142`）：`ConnMaxLifetime: o.IdleCheckFrequency`。`IdleCheckFrequency` 在注释中是「空闲连接检查频率」（`options.go:54-55`），而 `ConnMaxLifetime` 是连接**最大存活时间**；两者被当成同一个字段赋值，未配置 `IdleCheckFrequency`（零值 0）时 `ConnMaxLifetime = 0`，go-redis 侧对应「连接永不因寿命到期回收」。真正的空闲检查周期 go-redis 使用自身内部默认值，框架提供的 `WithIdleCheckFrequency`（`redis.go:167-171`）实际没有作用到检查频率上。
  2. 零值透传：`DialTimeout`/`ReadTimeout`/`WriteTimeout`/`PoolTimeout`/`IdleTimeout` 在四个 Config 方法中直接原样写入（`options.go:73-78` 等），框架未设置非零兜底默认值；`redis.go:29-31` 只对 `PoolSize` 做了默认 64，超时依赖 go-redis 内部默认行为，配置语义与注释（`options.go:42-52` 所述 5 秒/3 秒默认值）并不由框架保证。
  3. `PoolSize` 默认 64（`redis.go:29-31`）；go-redis 集群模式下该值是**每节点**连接数（`options.go:48` 注释亦注明），节点较多时总连接数 = 64 × 节点数，容易超出 Redis `maxclients` 或本机 fd 上限；低并发单节点场景 64 又偏大。
- 影响评估：属于资源配置健壮性问题：连接生命周期不回收会导致长连接长期复用、对端半开连接问题暴露不充分；零值超时在极端网络/挂死场景下的行为取决于库内兜底；连接数放大在集群部署中才会显现。对单命令 CPU/分配没有直接影响，故定 P2；真实影响程度需结合部署节点数、并发量与压测实测（需实测清单第 7 项）。
- 优化建议：
  1. 拆分语义：新增 `ConnMaxLifetime` 独立选项（建议默认 30 分钟量级），`IdleCheckFrequency` 映射到 go-redis 真正的空闲检查周期（若库版本未暴露该选项则在文档中明确删除该旋钮）。
  2. 在四个 Config 方法内对零值超时补框架级默认（Dial 5s、Read/Write 3s、PoolTimeout = ReadTimeout + 1s、IdleTimeout 30min 量级），与注释语义对齐。
  3. `PoolSize` 默认值按模式推导：集群/哨兵按「每节点」给出建议值并在文档/配置校验中提示总量；提供最小/最大值校验。

#### [P2] `stores/redis/string_commands.go:52` `Get`（同一闭包模式遍布全部 457 个 `*_commands.go` 方法）：每个包装方法在命令路径上分配一个逃逸闭包，并附带 `DoWithAcceptable` 接口回调层

- 问题：每个命令方法都形如 `err = r.brk.DoWithAcceptable(func() error { conn, err := getRedis(r); ...; return err }, acceptable)`（`string_commands.go:52-59`，其余文件逐一同构，已用静态扫描确认无例外）。该匿名函数在调用时被 `doReq(req func() error, ...)` 以接口值持有，通常逃逸到堆上，每次命令产生一次闭包分配；闭包内捕获 `r`、`ctx`、命名返回值 `val/err` 与入参。`acceptable` 包级变量（`utils.go:19-21`）本身无需每命令分配，但 `req` 接口包装伴随一次函数值/接口构造。
- 影响评估：这是包装层叠加在 go-redis 原生命令分配之上的额外开销；go-redis 内部每个命令本就有 `*XxxCmd` 分配，闭包占比相对有限，故定 P2。量级需要 benchmark 用 `benchstat -b` 对比「原生 client 直连」与「本封装」的 allocs/op 差值（需实测清单第 3 项）；管理类低频命令（ACL、Cluster、Config、SlowLog 等）无需优化，热点读/写命令值得评估。
- 优化建议：
  1. 优先落地 P0 条目（缓存 `node`）后，热点方法可提供无闭包的直连变体：先取缓存 `node`，再直接 `node.Get(ctx, key).Result()`，仅在需要熔断时走回调路径。
  2. 若保留统一包装，可将 `acceptable` 判定内联进 `doReq`，减少接口层间接；避免在闭包内重复 `getCtx(ctx)` 之外的局部对象分配。

---

## 二、stores/proxy

#### [P1] `stores/proxy/sql.go:34` `SQL.Write`（同 `SQL.Read:41`）：每次选库都调用 `util.Rand`，而 `util.Rand` 内部每次新建 `math/rand` 实例与随机种子

- 问题：`Write()`（`sql.go:34-38`）执行 `write := util.Rand(0, len-1)`；`Read()`（`sql.go:41-48`）读库列表非空时同样 `read := util.Rand(0, len-1)`。根因在范围外的 `util/math.go:22-34`：`Rand` 每次调用都执行 `r := rand.New(rand.NewSource(time.Now().UnixNano()))`（`math.go:32`），即每次选库：
  1. 新建一个 `*rand.Rand`（含内部 `lockedSource` 与源状态，堆分配）；
  2. `time.Now().UnixNano()` 取系统时钟做种子；
  3. `r.Intn(...)` 后该实例立即丢弃。
  另外以同一纳秒时间戳为种子：并发下多个 goroutine 在同一纳秒内进入会得到相同序列，导致选库倾斜（分布问题，顺带记录）。
- 影响评估：`Write`/`Read` 处于每个 SQL 请求选择数据源的必经路径（读写分离代理），属于请求级热点。每次调用 1 次随机对象 + 源状态堆分配 + 一次系统调用级时钟读取；在高并发 DB 请求路径上属于可避免的纯 CPU/分配开销。随机种子重复还会削弱读库负载均衡的均匀性。具体 allocs/op 与时钟开销需 benchmark 实测（需实测清单第 5 项）。
- 优化建议：
  1. 在 `SQL` 结构体内持有一个复用的随机源：可使用带 `sync.Mutex`/`sync.RWMutex` 保护的 `*rand.Rand`（初始化种子一次），或 `golang.org/x/exp/rand`（原生并发安全）；`Write`/`Read` 改为 `rng.Intn(n)`。
  2. 空列表与 `len == 1` 的快速返回：单节点时直接返回下标 0，跳过随机数。
  3. 根因 `util/math.go:22` 的 `Rand` 为全局共用工具，其他调用方同样受影响，列入跨模块提示；范围内先以 `SQL` 自持随机源规避。

#### [P2] `stores/proxy/redis.go:30` `Proxy.LoadRedis`（同 `LoadMongoDB:28`、`LoadElasticSearch:28`、`LoadSQL:56`）：每次 `Load` 做一次 `sync.Map.Load` + 类型断言

- 问题：各 `Load*` 方法（`redis.go:30-35`、`mongodb.go:28-33`、`elasticsearch.go:28-33`、`sql.go:56-61`）都是 `if f, ok := proxyPool.m.Load(group); ok { return f.(*X).X }` 的标准 `sync.Map` 读路径。`sync.Map` 对读多写少场景已做优化（命中走 read-only map 无锁读），类型断言为单步操作；未发现额外分配或循环。
- 影响评估：开销很小，属于注册期装配 + 请求期一次性查找的常规模式；仅当同一请求内反复 `Load` 同一分组时才会重复该查找。本条仅作覆盖性记录，不建议投入优化成本。
- 优化建议：业务侧可在请求上下文/装配阶段把解析后的后端实例注入，避免每次 `Load`；`Proxy` 本身无需改造。

---

## 三、stores/session

#### [P1] `stores/session/driver.go:30` `Session.Put`：先 `Exists` 再 `Get` 再 `Set` 三趟 Redis 往返；单字段写入把整份会话 JSON 读出、反序列化、改一个字段、再整体序列化写回

- 问题：`Put`（`driver.go:30-63`）流程：
  1. `:40` `h.Exists(ctx, session.Name)` 判断键是否存在（第 1 趟 RTT）；
  2. `:41` `content, _ = h.Get(ctx, session.Name)` 再把整份会话 JSON 取回（第 2 趟 RTT）；
  3. `:46` `json.Unmarshal([]byte(content), &m)` 整份反序列化为 `map[string]any`；
  4. `:51` `strings.Split(key, ".")` 拆路径，`:57` `setSliceMap` 写入一个字段；
  5. `:60` `json.Marshal(m)` 整份重新序列化；`:61` `h.Set(...)` 整体写回（第 3 趟 RTT，返回错误被忽略）。
  其中 `Exists` 与 `Get` 冗余：`Get` 在键不存在时返回 `redis.Nil`，一次 `Get` 即可区分存在与否，`Exists` 这一趟 RTT 完全可以去掉。更根本的问题是数据模型：**整份会话存为单个 JSON 文档**，任何一个顶层/嵌套字段的写入都要承担整份文档的读出、JSON 解码、整份 JSON 编码与写回，复杂度与会话大小 N 成正比；且读-改-写之间没有 `WATCH`/乐观锁或 Lua 保护，并发写相互覆盖（正确性背景，不计入性能评分）。
- 影响评估：会话存储通常挂在登录态用户请求链路上，`Put` 频次随用户操作增长。会话字段越多、单会话体积越大，每字段写的 CPU（两次 JSON 编解码 + `string`→`[]byte` 拷贝）与网络（3 趟 RTT）线性上升；多实例并发写同一会话时还会因覆盖导致数据丢失。在中型会话（几 KB ~ 几十 KB）、高频写场景下，JSON 编解码可能占到该方法 CPU 的主要部分；实际占比需结合会话体积与 QPS 实测（需实测清单第 4 项）。
- 优化建议：
  1. 立即消除冗余 RTT：去掉 `:40` `Exists`，单次 `Get` 并忽略 `redis.Nil`（键不存在时使用空 `{}` map），三趟降为两趟。
  2. 模型改造（收益最大）：顶层字段改用 Redis Hash 存储——`HSet(name, field, value)` 单字段直写，省去整文档读出与编解码；嵌套结构按需用二级 hash 或独立子键；仅当确需整文档原子语义时保留 JSON 模型。
  3. 保留 JSON 模型的过渡方案：把 `Get` 结果的 `[]byte` 用 `sync.Pool` 复用；`Marshal` 复用 buffer；写回前 `WATCH(name)` + `MULTI/EXEC` 保证读-改-写一致性（同时修复覆盖问题）。

#### [P2] `stores/session/driver.go:67` `Session.Get`：为读取一个字段，把整份会话 JSON 反序列化为 `map[string]any`，再沿点分路径下钻

- 问题：`Get`（`driver.go:67-85`）：`:76` `h.Get` 取回整份 JSON；`:77` `json.Unmarshal([]byte(content), &m)` 整份解码；`:79` `strings.Split(key, ".")`；`:84` `getSliceMap`（`:135-154`）逐层类型断言下钻。一次字段读取承担整份会话的 JSON 解码与 `map` 分配，`content`（string）到 `[]byte` 再产生一次完整拷贝；解码错误被 `_` 忽略，键不存在时 `m` 为 nil（读 nil map 在 Go 中安全返回零值）。
- 影响评估：`Get` 只有 1 趟 RTT（优于 `Put`），CPU 成本集中在整文档 JSON 解码；会话小（数百字节）时影响不明显，会话变大后每次字段读取都按 O(N) 解码，读路径高频时 CPU 与 GC 压力显著。等级 P2，优先级低于 `Put`。
- 优化建议：
  1. 配合 `Put` 的 Hash 模型：顶层字段 `HGet(name, field)` 直取，完全避免 JSON 解码。
  2. 过渡期避免全量 `map[string]any` 分配：改用只读路径提取的库（如 `gjson`）按路径取标量；`[]byte` 拷贝使用 pooled buffer。
  3. `getSliceMap` 内 `:147` `itMap = v.(map[string]any)` 未做类型保护，类型不符会 panic（正确性背景），改造模型时一并处理。

#### [P2] `stores/session/driver.go:89` `Session.Remove`：删除一个字段仍需整会话 Get + 全量反序列化 + 全量序列化 + Set 写回，两个分支重复同构代码

- 问题：`Remove`（`driver.go:89-113`）：`:99` `h.Get` 整份取回；`:100` `json.Unmarshal` 整份解码；`:102` `strings.Split`；`depth < 2` 分支 `:105` `delete(m, key)`、`:106` `Marshal`、`:107` `Set`；嵌套分支 `:110` `delSliceMap`（`:157-176`）、`:111` `Marshal`、`:112` `Set`。单字段删除同样承担整文档往返与两次编解码，且顶层/嵌套两个分支各自重复一遍 `Marshal`+`Set`；`Get`/`Set` 错误均被忽略。
- 影响评估：删除操作频次通常低于 `Put`/`Get`，故定 P2；模型不变时其 CPU 特征与 `Put` 相同（整文档 RMW），会话体积大时放大。
- 优化建议：Hash 模型下 `HDel(name, field)` 单字段删除，一步完成、无编解码；过渡期与 `Put` 共用同一份读-改-写辅助函数（含 `WATCH` 保护与 buffer 复用），消除两个分支的重复 `Marshal`/`Set`。

---

## 四、条目统计

### stores/redis（共 7 条：P0 = 1，P1 = 2，P2 = 4）

| 级别 | 条目 |
|---|---|
| P0 | 1. `redis.go:199` 等 — 每命令经全局 ResourceManager/SingleFlight 锁查找，实例未缓存 |
| P1 | 2. `string_commands.go:52` 等 — 每命令经熔断器全桶扫描 + 读写双锁 |
| P1 | 3. `hook.go:53` — attrs 跨命令累积 + 每命令栈回溯与命令串拼接 |
| P2 | 4. `redis.go:200` 等 — driverName 每命令重建/拼接/map 遍历 |
| P2 | 5. `scripting_commands.go:10` — Eval 每次发送全脚本，无 EVALSHA 缓存封装 |
| P2 | 6. `options.go:78` 等 — IdleCheckFrequency 误映射 ConnMaxLifetime、零值超时透传、PoolSize 语义 |
| P2 | 7. `string_commands.go:52` 等 — 每方法闭包逃逸分配（模式遍布全部命令文件） |

### stores/proxy（共 2 条：P0 = 0，P1 = 1，P2 = 1）

| 级别 | 条目 |
|---|---|
| P1 | 1. `sql.go:34` — Write/Read 每次新建 math/rand 实例选库 |
| P2 | 2. `redis.go:30` 等 — Load* 每次 sync.Map.Load + 断言（标准用法，开销小） |

### stores/session（共 3 条：P0 = 0，P1 = 1，P2 = 2）

| 级别 | 条目 |
|---|---|
| P1 | 1. `driver.go:30` — Put 三趟 RTT + 整会话 JSON 读-改-写 |
| P2 | 2. `driver.go:67` — Get 单字段读取整会话 JSON 解码 |
| P2 | 3. `driver.go:89` — Remove 单字段删除整会话读-改-写 |

**合计：P0 = 1，P1 = 4，P2 = 7（总计 12 条）。**

---

## 五、高价值候选列表（供全局 Top 10 评选，按预期收益排序）

| # | 位置 | 一句话问题 | 预期收益 |
|---|---|---|---|
| 1 | `stores/redis/redis.go:199`（`getClient`/`getCluster`/`getFailover`/`getRing`，`getRedis:180`） | 每条命令经全局 SingleFlight 互斥锁 + 资源表双锁查找，底层实例未在 `Client` 缓存 | 消除每命令 3 次锁临界区与两次 map 查找，去掉跨实例共享锁竞争；高频短命令尾延迟与锁等待直接下降，改造范围小、收益最确定 |
| 2 | `stores/redis/string_commands.go:52` 等全部命令方法（根因 `core/resource/breaker.go:390`） | 每条命令熔断器 history 全 40 桶扫描（RLock）+ Add（写锁 + 过期桶循环重置） | 去掉每命令两把共享锁与固定 O(40) 遍历；高并发下同一 `Client` 上的命令不再在熔断器锁上排队，CPU 与尾延迟双降 |
| 3 | `stores/redis/hook.go:53` `ProcessHook`（`:85` `ProcessPipelineHook`） | 闭包共享 attrs 反复 append 造成属性累积与内存只增不减；每命令 11 层 `runtime.Caller` + 全量 `CmdString` 分配 | trace 开启路径消除单调内存增长；去掉成功命令的栈回溯与大命令串分配，显著降低 trace 模式 CPU 与堆分配 |
| 4 | `stores/session/driver.go:30` `Put` | Exists+Get+Set 三趟 RTT，单字段写整会话 JSON 读-改-写 | 先省掉 Exists 一趟 RTT；改 Hash 模型后单字段写一次 HSet 完成，消除整文档两次 JSON 编解码，会话越大收益越大 |
| 5 | `stores/proxy/sql.go:34` `Write` / `:41` `Read` | 每次选库新建 `math/rand` 实例并以当前纳秒为种子 | 每请求去掉一次随机对象堆分配与时钟读取；修复同纳秒并发种子重复导致的读库倾斜；改造仅在 `SQL` 结构体加一个复用随机源 |
| 6 | `stores/redis/redis.go:249-252`（同 `:200`/`:216`/`:232`） | driverName 每命令重建：字符串拼接、Implode、ring map 遍历 | 配合候选 #1 缓存实例后一并消除；ring/集群模式去掉每命令的 map 遍历与多次字符串重分配 |
| 7 | `stores/redis/scripting_commands.go:10` `Eval`（`:36` `EvalRO`） | 每次执行传输完整 Lua 脚本，无 Script/EvalSha 自动回退封装 | 复用型脚本改走 EVALSHA 后命令体不再携带脚本源码，降低网络传输与服务端解析；脚本越大、频次越高收益越大 |
| 8 | `stores/redis/options.go:78`（`:99`/`:119`/`:142`） | IdleCheckFrequency 被错映射为 ConnMaxLifetime，连接生命周期不回收；超时零值透传 | 修正连接回收语义，避免长连接永不退役；默认超时兜底避免挂死场景；属稳定性收益，性能收益需压测确认 |
| 9 | `stores/session/driver.go:67` `Get`（`:89` `Remove`） | 单字段读/删承担整会话 JSON 解码（与 Hash 模型改造一并落地） | 读路径消除整文档解码与 map 分配；大会话场景下 CPU 与 GC 压力明显下降 |
| 10 | `stores/redis/string_commands.go:52` 等（全部命令方法） | 每方法逃逸闭包 + 接口回调层 | 热点命令直连化后减少每命令一次闭包堆分配；收益受 go-redis 原生命令分配基数影响，需 benchmark 确认后再排优先级 |

---

## 六、需实测验证的条目清单（静态阅读无法定量，不得作为已确认结论上报）

1. **候选 #1（SingleFlight/资源表锁）**：需在多 `*Client`、高并发下用 pprof `mutex`/`block` 与 `go tool trace` 测量 `singleflight.go:57`、`sourcemanager.go:52/64` 的锁等待时长；对比「`Client` 缓存实例」前后的 allocs/op 与 P99 延迟，确认锁竞争是否真正构成瓶颈。
2. **候选 #2（熔断器每命令统计）**：需 benchmark 对比 `r.brk.DoWithAcceptable` 包装前后的单命令耗时，量化 `breaker.go:425 Reduce`（40 桶遍历）与 `breaker.go:455 Lock` 的 CPU 占比；确认 `stat.lock` 在并发下的竞争排队长度。
3. **候选 #10（包装闭包分配）**：需 `benchstat -b` 对比原生 `*redis.Client` 与本封装的 `allocs/op`、`B/op`，确认闭包逃逸分配的绝对量级；再决定是否值得为热点方法做直连变体。
4. **候选 #4/#9（会话 JSON 模型）**：需在不同会话体积（0.5KB / 5KB / 20KB）下分别测量 `Put`/`Get`/`Remove` 的 RTT 次数、JSON `Unmarshal`/`Marshal` CPU 耗时与 GC 占比；验证 `Exists` 冗余 RTT 在网络耗时中的占比。
5. **候选 #5（util.Rand 新建实例）**：需 benchmark 测量 `util.Rand` 单次调用的分配次数与 `time.Now().UnixNano()` 开销；并验证高并发同纳秒种子重复是否真实造成选库倾斜（分布直方图）。
6. **候选 #7（Eval 全脚本传输）**：需结合业务脚本长度与调用频次统计；小脚本、低频调用下收益有限，不应仅凭代码形态定级为高优。
7. **候选 #8（连接池配置）**：需压测验证零值超时在慢网络/对端挂死时的实际行为；集群模式下统计 `PoolSize=64 × 节点数` 的实际连接总量与 `ErrPoolTimeout` 触发频率；`ConnMaxLifetime=0` 时长连接的存活率需长稳测试确认。
8. **`hook.go` attrs 累积（候选 #3）**：需长跑压测观察 `attrs` 底层数组容量与对应 span 属性数量的增长曲线，量化内存泄漏速率；`Caller(11)` 与 `CmdString` 的 CPU/分配需在 trace 开关开启状态下单独采样。

## 七、跨模块提示（仅记录线索，不在本章节展开审查）

- **`core/resource/breaker.go`（熔断器实现，本范围外）**：`doReq:390` 每命令 `accept()` + `markSuccess/markFailure`；`history:424` 经 `Reduce:463` 全桶扫描；`Add:454` 写锁 + `updateOffset:490` 过期桶循环；`buckets=40`（`:29`）。建议由对应审查组评估窗口统计的锁粒度、遍历复杂度与时钟调用次数。
- **`core/resource/sourcemanager.go:50` / `core/resource/singleflight.go:36`（资源管理器与单飞，本范围外）**：命中已缓存资源时仍走 SingleFlight 全局互斥锁与资源表 RLock；建议增加无锁 fast path。`singleflight.go:57` 为 `sync.Mutex`（非 RWMutex），读多场景可优化。
- **`util/math.go:22` `Rand`（本范围外）**：`:32` 每次 `rand.New(rand.NewSource(time.Now().UnixNano()))`，新建实例 + 纳秒种子；除 `stores/proxy/sql.go` 外，其他 `util.Rand` 调用方同样受影响，建议工具层改为进程级复用并发安全随机源。
- **`util/array.go:209` `Implode`（本范围外）**：每次新建 `bytes.Buffer`，可针对已知长度场景预分配；`stores/redis/redis.go:216/232` 的集群/哨兵 driverName 每次命令调用它。
- **同模式对照线索**：`stores/sql`、`stores/mongodb`、`stores/elasticsearch` 等存储封装若同样采用 `brk.DoWithAcceptable` + `clientManager.GetResource` + 每命令闭包模式，应纳入全局统一优化（本章节仅确认 `stores/redis` 存在该模式）。

## 八、非性能类观察（正确性/接线问题，供对应审查组参考，不计入本报告 P0/P1/P2 性能条目）

1. **`stores/redis/redis.go:228` `getCluster` 类型断言不匹配（正确性，集群模式必现 panic）**：`redis.NewClusterClient`（`redis.go:219`）在 go-redis v9.23.0 中返回 `*redis.ClusterClient`（`osscluster.go:1399`），而 `:228` 执行 `val.(*redis.Client)`，运行时取值会触发类型断言 panic。`getFailover:244` 断言 `*redis.Client` 与 `NewFailoverClient` 返回类型（`*redis.Client`，`sentinel.go:546`）一致，无此问题；`getRing:264` 断言 `*redis.Ring` 一致。建议由正确性审查组修复 `getCluster` 的返回类型与 `RedisNode` 适配。
2. **`stores/redis` 的 `EnableMetric` 开关与 `MetricsHook` 未接线**：`redis.go:204-206`/`220-222`/`236-238`/`256-258` 仅在 `EnableTrace` 时挂载 `OpenTraceHook`；`options.go:59` 的 `EnableMetric` 与 `hook.go:118` 的 `MetricsHook` 在任何装配路径上都未被 `AddHook`，指标开关实际不生效（接线遗漏，非性能缺陷）。
3. **`stores/session/driver.go` 错误处理与并发安全（正确性背景）**：`Put:41/61`、`Get:76/77`、`Remove:99/100/107/112` 忽略 Redis 与 JSON 错误；整文档读-改-写无 `WATCH`/事务保护，并发写互相覆盖；`setSliceMap:127`、`getSliceMap:147`、`delSliceMap:169` 的 `v.(map[string]any)` 无类型保护，路径中段类型不符会 panic。
4. **`stores/redis/options.go:79/99/119/142` 字段语义错位**（已在性能条目候选 #8 中记录其资源回收影响；此处标注其配置语义错误本质：`IdleCheckFrequency` 与 `ConnMaxLifetime` 是两个独立概念，应分别映射独立选项）。
5. **`stores/redis/command.go:607-663`**：整段 `MGetByPipeline` 为注释掉的历史实现，其中包含 `time.Now()` 日志、`sync.WaitGroup` + 锁 + 无界 `errors` channel 的旧设计；当前不参与编译，仅提示后续若恢复 Pipeline 批量 MGet 时，应直接使用 go-redis `Pipelined` 与按槽位分组的成熟实现，不要复用注释中的旧模型。
