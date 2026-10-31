# stores 包性能审查报告（stores/null、stores/sql、stores/mongodb、stores/elasticsearch）

> 审查范围：`stores/null`（含 `stores/null/convert`，27 个非测试 `.go`）、`stores/sql`（8 个非测试 `.go`）、`stores/mongodb`（3 个非测试 `.go`）、`stores/elasticsearch`（2 个非测试 `.go`，`elastic_test.go` 不在审查范围）。
> 审查方式：逐文件人工通读，热点路径重点检查内存分配、反射、锁/并发、I/O 与连接池、日志/指标/追踪开销；每条条目均已回读原文核对 `文件:行号` 与函数名。
> 说明：`stores/elasticsearch/elastic.go:170` 指标开关条件反向为已知 FIXME（逻辑正确性问题），本报告不重复上报；其性能角度另在相关条目内独立描述。

---

## 一、stores/sql

### P1

#### [P1] `stores/sql/interceptor.go:92` `TraceInterceptor`：闭包级 `attrs` 切片被每次请求并发追加并无限增长（数据竞争 + 内存泄漏）

- 问题：`TraceInterceptor` 在工厂闭包内创建 `attrs := []attribute.KeyValue{...}`（第 92 行），该切片被后续所有查询共享；每次被包装的处理器执行时，第 103–107 行又在该共享切片上 `attrs = append(attrs, semconv.CodeFunction(fn), semconv.CodeFilepath(file), semconv.CodeLineNumber(line))`。GORM 回调在高并发下由多个 goroutine 同时进入同一个包装函数，共享切片的 `append` 构成数据竞争；同时切片只增不减，随查询次数无限膨胀，且历次栈帧属性会被累积进后续每一个 Span 的 `attrs...`。
- 影响评估：仅 `EnableTrace=true` 时生效。场景为开启链路追踪的高 QPS 服务：① 每次 append 都可能触发切片扩容拷贝（容量不足时 O(n) 复制）；② 切片终身不回收，进程长时间运行后内存单调上涨；③ 并发 append 在 Go 内存模型下属于未定义行为，存在切片头损坏/运行时 panic 风险（竞争维度另列，本条目从性能与资源泄漏角度上报）。量级随累计查询次数线性增长，与当前并发无关地持续占用堆内存。
- 优化建议：在每个请求的处理器内部局部构造属性切片，不要复用闭包级变量：
  ```go
  return func(scope *gorm.DB) {
      if ctx := scope.Statement.Context; ctx != nil {
          fn, file, line := call.Caller(4)
          attrs := []attribute.KeyValue{
              semconv.RPCSystemKey.String("gorm"),
              semconv.HostName(hostname.Hostname()),
              semconv.CodeFunction(fn),
              semconv.CodeFilepath(file),
              semconv.CodeLineNumber(line),
          }
          // ... 仅用本次局部 attrs 启动 Span
      }
  }
  ```
  静态属性（RPC 系统名、主机名）可在工厂层只准备一份只读模板，请求时 `append` 到新的局部切片中。

#### [P1] `stores/sql/logger.go:87` `WrapLogger.Trace`：未判定任何日志级别就先构造 SQL 字符串、执行栈回溯并分配 Fields 映射（`WrapEntry.Trace` 同构，第 162–181 行）

- 问题：`WrapLogger.Trace`（第 85–104 行）执行顺序为：
  1. 第 87 行 `sql, _ := fc()`：无条件调用 GORM 传入的 `fc` 闭包，即 `Dialector.Explain(stmt.SQL.String(), stmt.Vars...)`，把 SQL 模板与全部参数替换为完整 SQL 字符串；
  2. 第 88 行 `fields := logrus.Fields{}`：无条件分配映射；
  3. 第 89–91 行：只要 `SourceField != ""` 就调用 `utils.FileWithLineNum()`（底层 `runtime.Caller(3)` 栈回溯 + `runtime.FuncForPC` 取函数名）。
  真正的日志分支判断（错误 92–96、慢查询 97–100、Debug 101–103）在这些开销之后才发生。而 `Client.Open`（`client.go:298`）无条件 `SetSourceField("gorm")`，意味着无论 `EnableDebug` 是否开启、是否慢查询、是否出错，每一条 SQL 都会支付一次栈回溯；`fc()` 也在 Debug 关闭且无错误、非慢查询时白白执行一次 SQL 拼装。`WrapEntry.Trace`（第 162–181 行，对应第 164、165、166–168 行）是完全相同的复制实现。
- 影响评估：这是所有 SQL 查询的必经回调，覆盖读、写、Raw 全部操作。栈回溯（`runtime.Caller` + `FuncForPC.Name`）在标准库层面属于微秒级但不可忽略的固定开销，在高 QPS 下会累计为可观的 CPU 占比；`fc()` 生成完整 SQL 字符串伴随一次字符串/参数格式化分配。开销与最终是否输出日志完全无关——Debug 关闭的生产环境也照付。
- 优化建议：先判定级别再产生开销：
  ```go
  func (l *WrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
      elapsed := time.Since(begin)
      // 1) 先判断本次是否需要记录日志；不需要则直接返回，不调用 fc()、不做栈回溯
      acceptable := Acceptable(err) == nil
      needSource := l.SourceField != ""
      if acceptable && l.SlowThreshold == 0 && !l.Debug {
          return
      }
      // 2) 仅在需要记录时才构造 SQL 与源码位置
      sql, _ := fc()
      fields := logrus.Fields{}
      if needSource {
          fields[l.SourceField] = utils.FileWithLineNum()
      }
      // ... 原错误 / 慢查询 / Debug 分支
  }
  ```
  注意慢查询分支需要 `elapsed > l.SlowThreshold` 才记录，可在 `fc()` 之前先比较 `elapsed`；`WrapEntry.Trace` 同步改造。

#### [P1] `stores/sql/client.go:327` `Open`：连接池参数取配置零值直接下发，缺少兜底默认值，空闲连接数为 0 会造成连接反复新建

- 问题：第 327–330 行直接把 `c.MaxIdleConns / c.MaxOpenConns / c.MaxLifetime / c.MaxIdleTime` 传给 `database/sql`：
  ```go
  sqlDB.SetMaxIdleConns(c.MaxIdleConns)
  sqlDB.SetMaxOpenConns(c.MaxOpenConns)
  sqlDB.SetConnMaxLifetime(c.MaxLifetime)
  sqlDB.SetConnMaxIdleTime(c.MaxIdleTime)
  ```
  调用方未通过 `WithMaxIdleConns` 等 Option 配置时，这些字段保持零值：`SetMaxIdleConns(0)` 表示连接池不保留任何空闲连接（每次用完即关闭）；`SetMaxOpenConns(0)` 表示不限制最大打开连接；`SetConnMaxLifetime(0)` / `SetConnMaxIdleTime(0)` 表示连接永不回收。`NewClient`（`client.go:107`）本身不设置任何默认值。
- 影响评估：在业务代码只配置主机/账号、未调连接池 Option 的常见使用方式下，每个请求获取连接后归还即关闭，下一个请求需重新完成 TCP 三次握手、TLS 协商与数据库认证，连接无法复用；并发升高时又因 `MaxOpenConns=0`（不设上限）可能无限制新建连接，同时威胁数据库端连接上限。该问题在服务启动后持续影响每一次数据库访问，量级随 QPS 与连接重建频次放大，属于资源层结构性缺陷。
- 优化建议：在 `Open` 内对零值参数填充合理默认值（例如 `MaxIdleConns` 默认 10、`MaxOpenConns` 默认 100、`MaxLifetime` 默认 1 小时、`MaxIdleTime` 默认 30 分钟），仅当调用方显式传入正值时才覆盖；或在 `NewClient` 的 Option 装配阶段即写入默认配置。

### P2

#### [P2] `stores/sql/interceptor.go:101` `TraceInterceptor`：每次查询都执行 SQL 全量展开、运行时栈回溯与 metadata 拷贝（追踪开启路径）

- 问题：第 101 行 `scope.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)` 把参数回填进 SQL；第 102 行 `call.Caller(4)` 做一次 4 层栈回溯；第 108–113 行每次从 context 取出 `metadata.MD` 后执行 `md.Copy()`（无 metadata 时新建空 `metadata.MD{}`）。这些开销位于被包装的 SQL 回调内部，追踪开关打开后每次操作都执行。
- 影响评估：仅 `EnableTrace=true` 时触发。`Explain` 涉及字符串拼接与参数格式化分配；`Caller(4)` 为栈遍历；`md.Copy()` 拷贝底层 map。单次成本不高，但叠加在每条 SQL 的追踪路径上，高 QPS 下属于可压缩的固定开销。
- 优化建议：静态 metadata 可复用同一只读实例；栈回溯信息若仅用于定位调用方，可改为在 GORM logger 已有的源码位置基础上复用（见上一条 logger 改造）；`Explain` 仅在 Span 需要记录完整语句时执行。

#### [P2] `stores/sql/client.go:301` `Open`：`NowFunc` 每次取时间都冗余调用 `.Local()`

- 问题：第 301–303 行：
  ```go
  NowFunc: func() time.Time {
      return time.Now().Local()
  },
  ```
  `time.Now()` 返回的已是本地时区时间，`.Local()` 再做一次位置归一化与位置查找，每次 GORM 写入 `CreatedAt/UpdatedAt` 等时间戳字段都会调用。
- 影响评估：仅发生在 insert/update 写入路径，单次成本极小；`.Local()` 内部会做 `time.Location` 比较与可能的时区映射判断，属可消除的冗余调用。
- 优化建议：直接 `return time.Now()`；如需固定本地时区可在初始化期一次性 `time.Local` 处理，而不是每次回调转换。

#### [P2] `stores/sql/ext.go:19` `Pagination.GetOffset` / `GetLimit`、`Scope.GetStartString` 等：目标类型已确定仍经由 cast 反射转换指针

- 问题：`GetOffset`（第 19–24 行）、`GetLimit`（第 27–32 行）、`GetStartString`/`GetEndString`（第 44–57 行）、`GetStartInt64`/`GetEndInt64`（第 60–73 行）、`GetStartFloat64`/`GetEndFloat64`（第 76–89 行）、`Sort.GetName`（第 101–106 行）、`GetDesc`（第 109–114 行）接收的分别是 `*int64`、`*string`、`*bool` 等已知具体类型，却统一交给 `cast.ToInt64` / `cast.ToString` / `cast.ToBool`，cast 内部要做 `reflect.TypeOf` 与类型分支判断。
- 影响评估：这些方法属于请求参数绑定阶段，每请求最多调用数次，不在数据库查询热路径上；反射分发开销有限，但完全可以消除。
- 优化建议：已知类型直接解引用返回，例如：
  ```go
  func (obj *Pagination) GetOffset() (int64, error) {
      if obj.Offset == nil {
          return 0, errors.New("offset is nil")
      }
      return *obj.Offset, nil
  }
  ```
  其余 `Scope`/`Sort` 取值方法同构改造。

---

## 二、stores/mongodb

### P1

#### [P1] `stores/mongodb/clientmanager.go:173` `NewModel`：每创建一个 Model 都先发起一次 `Ping` 网络往返（`NewCollection` 第 197 行同构）

- 问题：`NewModel`（第 172–184 行）在第 173 行先 `m.Ping()`，成功后才 `newCollection(...)`；`NewCollection`（第 196–203 行）在第 197 行同样先 `Ping()`；`MustNewCollection`（第 206–213 行）在第 207 行自行调用 `m.Ping()`；`MustNewModel`（第 187–193 行）通过第 188 行调用 `NewModel` 间接 Ping。`Ping` 会向 MongoDB 服务端发起一次真实命令并等待响应。
- 影响评估：若业务层在每次请求处理时调用 `NewModel`/`NewCollection` 取集合句柄（本封装的公开 API 形态容易被这样使用），每个请求都要额外支付一次 Ping 的网络 RTT（同机一般亚毫秒到数毫秒，跨机房/网络抖动时显著放大），且该 RTT 与实际业务读写串行叠加。若业务在启动期一次性缓存 Model，则无影响——**调用方式需实测确认**（见文末待实测清单）。
- 优化建议：启动建立客户端时（`NewMongoDBClient`，第 127–137 行已有 `cli.Ping`）完成连通性校验；`NewModel`/`NewCollection` 去掉每次 Ping，直接构造 `DecoratedCollection`。健康检查由独立的后台探活或连接池自身机制承担。

#### [P1] `stores/mongodb/collection.go:653` `FindMany`：手写反射逐文档 `reflect.New` + `reflect.Append`，结果切片从零容量反复扩容，未使用驱动优化的 `cursor.All`

- 问题：`FindMany`（第 613–685 行）在第 653–671 行自行实现游标遍历：
  ```go
  val := reflect.ValueOf(documents)          // 653
  ...
  slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)  // 659：初始容量 0
  itemTyp := val.Elem().Type().Elem()
  for result.Next(ctx) {
      item := reflect.New(itemTyp)            // 662：每文档一次堆分配 + 反射
      err = result.Decode(item.Interface())
      ...
      slice = reflect.Append(slice, reflect.Indirect(item))  // 669：反射追加，切片反复扩容拷贝
  }
  val.Elem().Set(slice)                        // 671
  ```
  同文件 `Aggregate`（第 374 行）直接使用 `cursor.All(ctx, result)`，而 `FindMany` 放弃了驱动内置的批量解码路径：每条文档都经历 `reflect.New` 分配指针、`Decode` 到指针、`reflect.Indirect` 取址、`reflect.Append` 反射拷贝；切片初始容量为 0，按 Go 切片增长策略（1→2→4→8…）多次重新分配并整体拷贝。
- 影响评估：列表查询是典型高频读路径。返回 N 条文档时产生 N 次反射指针分配、N 次反射字段写入与 log₂N 次切片整体扩容拷贝；结果集越大，反射与扩容分摊成本越高，且 `cursor.All` 在驱动内部可基于元素类型做缓存式解码，本实现完全错失。
- 优化建议：改为 `result, err := c.Collection.Find(...); err = result.All(ctx, documents)`（与 `Aggregate` 第 374 行保持一致），由 mongo-driver 直接解码到调用方传入的切片；类型校验在调用前完成。若必须保留手写遍历，应先 `cursor.RemainingBatchLength` 或目标切片容量预估后预分配 `reflect.MakeSlice(typ, 0, capHint)`，避免反复扩容。

#### [P1] `stores/mongodb/collection.go:829` `BeforeCreate` / `collection.go:871` `BeforeUpdate`：写入前用反射把结构体二次转换为 `bson.M`，驱动随后再次编码，形成重复反射与重复分配

- 问题：`BeforeCreate`（第 829–867 行）对 `reflect.Struct` 分支：第 846–848 行逐字段 `data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()`（每字段一次 tag 读取 + `Interface()` 装箱为 `any` + map 写入），第 849 行 `dataVal := reflect.ValueOf(data)` 再包一层反射值，第 850 行与第 854 行对 `"Id"` 字段做了两次独立的 `FieldByName` 反射查找；`BeforeUpdate`（第 871–901 行）结构体分支第 888–894 行同构，第 889 行每字段额外 `Tag.Lookup("over")`，第 895 行同样多余地 `reflect.ValueOf(data)`。转换得到 `bson.M` 后，mongo-driver 在 `InsertOne`/`InsertMany`/`UpdateOne`/`UpdateMany` 内部还要对这份 map 再做一次 BSON 编码；而驱动本身原生支持直接写入结构体（带 bson tag），本层反射属于叠加在驱动之上的第二遍类型遍历。
- 影响评估：覆盖 `InsertOne`（第 262 行）、`InsertMany`（第 306 行，切片内逐元素递归调用，文档数 × 字段数）、`UpdateOne`（第 441 行）、`UpdateMany`（第 527 行）全部写路径。单文档字段越多，逐字段反射、`Interface()` 装箱与 map 分配越多；批量插入时成本按文档数放大。`InsertMany` 的切片分支（第 837–842 行）每个元素还要做一次 `.(bson.M)` 类型断言。
- 优化建议：优先方案是去掉结构体→`bson.M` 的预转换，直接把结构体（或指针）传给驱动，由驱动使用其带缓存的 bson 编解码完成；`_id` 生成可在结构体构造层处理。若必须保留 `bson.M` 形态：① 结构体类型到字段列表的映射可用 `sync.Map` 按 `reflect.Type` 缓存，避免每次遍历字段与重复读 tag；② 合并第 850/854 行两次 `FieldByName` 为一次；③ 去掉第 849、895 行无意义的 `reflect.ValueOf(data)`，直接操作 `data` map。

### P2

#### [P2] `stores/mongodb/collection.go:121` 等 15 个操作方法：方法内部重复 `NewTracer`、`call.Caller(6)`、`metadata.New(nil)`，未在集合句柄创建期一次性复用

- 问题：以下每个公开方法的追踪分支内都重复执行相同三段构造（`NewTracer` → `Caller(6)` → `metadata.New(nil)`）：
  - `CreateIndex`：第 121、125、131 行
  - `ListIndexes`：第 164、168、174 行
  - `DropIndex`：第 217、221、227 行
  - `InsertOne`：第 264、268、274 行
  - `InsertMany`：第 308、312、318 行
  - `Aggregate`：第 350、354、360 行
  - `UpdateOrInsert`：第 398、402、408 行
  - `UpdateOne`：第 443、447、453 行
  - `UpdateOneRaw`：第 486、490、496 行
  - `UpdateMany`：第 529、533、539 行
  - `FindOne`：第 571、575、581 行
  - `FindMany`：第 618、622、628 行
  - `Delete`：第 694、698、704 行
  - `Drop`：第 747、751、757 行
  - `Count`：第 789、793、799 行

  其中 `NewTracer`（`core/trace/trace.go:38`）每次都新建 `Tracer` 包装结构并向 provider 查询名为 `"ratel"` 的 tracer；`metadata.New(nil)` 每次分配一个空 map；`call.Caller(6)` 每次做固定 6 层栈回溯。
- 影响评估：仅 `EnableTrace=true` 时触发。每个操作额外支付一次 tracer 包装分配、一次空 metadata 分配与一次栈遍历；15 处复制粘贴导致后续优化需要逐方法修改，也容易出现本类不一致。
- 优化建议：在 `newCollection`（第 56–65 行）构造期创建一次 `tracer` 与基础属性，存入 `DecoratedCollection`；抽出公共的 `startSpan(ctx, op)` 方法统一处理 caller 属性与 metadata，各方法只传入操作名。`metadata.New(nil)` 在客户端 Span 场景不需要注入上游传播时可省略 carrier。

#### [P2] `stores/mongodb/collection.go:146` 等：指标标签每次操作重复调用 `c.Database().Name()` 与 `c.Name()`，未在集合句柄上缓存

- 问题：除 `UpdateOne`（第 468–472 行标签口径不同，另见跨模块提示）外，每个方法的指标分支都重复取库名与集合名拼标签，位置包括：第 146、148、150 行；199、201、203；245、247、249；289、291、293；332、334、336；379、381、383；424、426、428；510、512、514；553、555、557；600、602、604；676、678、680；730、732、734；772、774、776；815、817、819。
- 影响评估：`c.Database()` 每次返回一个新的 `*Database` 值，再取 `.Name()`；单次成本很低，但每次指标上报都重复这条调用链，标签字符串也会随 `WithLabelValues` 再做一次 hash/查找。属于可在构造期一次性固化的派生数据。
- 优化建议：在 `newCollection` 时缓存 `dbName := c.Database().Name()`、`collName := c.Name()` 到 `DecoratedCollection` 字段，指标上报直接引用缓存值。

#### [P2] `stores/mongodb/collection.go:590` `FindOne` / `collection.go:637` `FindMany`：每次查询都分配 `options.FindOneOptions` / `options.FindOptions`，且 Skip、Limit 恒传非空指针

- 问题：`FindOne` 第 590–594 行构造 `&options.FindOneOptions{Skip: &c.skip, Sort: c.sort, Projection: c.fields}`；`FindMany` 第 637–642 行构造 `&options.FindOptions{Skip: &c.skip, Limit: &c.limit, Sort: c.sort, Projection: c.fields}`。即使 `skip`、`limit` 为零值，也把字段地址传进去；每次查询新建 options 结构体与其内联指针。
- 影响评估：每次读查询一次小结构体分配；零值指针传入驱动后仍会参与选项解析。量级小，随读 QPS 累计。
- 优化建议：零值时不传指针（`var skip *int64; if c.skip > 0 { skip = &c.skip }`）；同一集合句柄上可缓存不含条件的基础 options，每次仅覆盖差异字段。

---

## 三、stores/elasticsearch

### P1

#### [P1] `stores/elasticsearch/elastic.go:159` `RoundTrip`：每个写请求都把请求体全量读入内存，但读出的字节在追踪与指标中从未被使用

- 问题：第 158–165 行：
  ```go
  contentLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
  if r.Body != nil && contentLength < MaxContentLength {
      buf, err := io.ReadAll(r.Body)
      if err != nil {
          return nil, err
      }
      r.Body = io.NopCloser(bytes.NewBuffer(buf))
  }
  ```
  回读请求体的目的（按注释“请求体未超过 MaxContentLength 时回读以便后续复用”）本是供追踪记录请求体，但实际 Span 只在第 152–153 行记录 `elastic.method` 与 `elastic.url`，body 字节自始至终没有被读取使用。代价是：对每个带 body 的 ES 写请求（Index/Bulk/Update/Delete/UpdateByQuery 等）先 `io.ReadAll` 全量读入内存分配 `buf`，再包成 `bytes.Buffer` 重新挂回 `r.Body`，形成一次额外的全量拷贝与内存占用；超过 `MaxContentLength`（64KB）的大 body 则跳过回读、直接透传。
- 影响评估：所有开启 `EnableMetric` 或 `EnableTrace` 的 ES 客户端都会安装该 transport（第 91 行），写接口每条请求都多一次 body 全量读 + 全量分配；Bulk 批量请求体可达数 MB，虽然 >64KB 会被跳过，但 64KB 以内的批量分片仍受影响。该读写完全是无效功，不产生任何输出。
- 优化建议：删除第 158–165 行的 body 回读逻辑，`r.Body` 直接透传给 `t.Transport.RoundTrip(r)`；若后续确有追踪记录请求体的需求，再按开关单独实现并限制采样比例。

#### [P1] `stores/elasticsearch/elastic.go:110` `ESTraceServerInterceptor`：使用零值 `http.Transport`，空闲连接上限与各类超时均未配置，高并发下连接复用受限

- 问题：第 109–117 行：
  ```go
  newESTracedTransport.Transport = &http.Transport{}
  ```
  零值 `http.Transport` 的 `MaxIdleConnsPerHost` 默认为 2（每个 ES 节点仅保留 2 个空闲连接），且 `DialContext`、`TLSHandshakeTimeout`、`ResponseHeaderTimeout`、`IdleConnTimeout` 均保持零值默认；未显式设置连接池大小与拨号超时。
- 影响评估：当 ES 客户端被多 goroutine 并发使用时，每节点超过 2 个空闲连接后，新归还的连接会被关闭，下一次请求需要重新建立 TCP（及可能的 TLS）连接，连接复用率低；缺少拨号/响应头超时会在 ES 节点网络抖动时长时间挂起调用 goroutine。该 transport 在 `EnableMetric || EnableTrace` 时安装（第 90–92 行），属于线上常见配置。
- 优化建议：显式构造带参数的 `http.Transport`，例如设置 `MaxIdleConnsPerHost`（按并发量级取 32~100）、`MaxIdleConns`、`IdleConnTimeout: 90 * time.Second`、`DialContext` 带 `Timeout`/`KeepAlive`、`TLSHandshakeTimeout`、`ResponseHeaderTimeout`；可在 `http.DefaultTransport` 基础上复制后改字段以继承成熟默认值。

### P2

#### [P2] `stores/elasticsearch/elastic.go:138` `RoundTrip`：每次 HTTP 调用内部重复创建 Tracer、做栈回溯、分配 metadata；Span 名硬编码为 "Count"

- 问题：第 138 行 `globalTrace.NewTracer(...)` 每次请求新建 tracer 包装；第 142 行 `call.Caller(6)` 每次请求栈回溯；第 143–147 行局部属性追加；第 148 行 `metadata.New(nil)` 每次分配空 map；第 149 行 Span 名固定写死 `"Count"`（与实际操作无关，追踪语义问题，顺带记录）。
- 影响评估：仅 `EnableTrace=true` 时触发，单次成本小，与 mongo 条目同源；tracer 可在 transport 构造期创建一次。
- 优化建议：在 `ESTraceServerInterceptor`（第 108–117 行）内创建一次 `tracer` 与静态属性，`RoundTrip` 只取请求方法/路径；`metadata.New(nil)` 无上游传播提取需求时移除。

#### [P2] `stores/elasticsearch/elastic.go:158` `RoundTrip`：每次请求都 `strconv.Atoi` 解析 `Content-Length` 头

- 问题：第 158 行对每个出站请求读取 `Content-Length` 字符串并做一次整数解析；解析失败时 `contentLength` 保持 0，随后进入第 159 行判断。在删除上一条无效 body 回读后，该行随之失去意义；保留期间每次请求多一次字符串到整数转换。
- 影响评估：单次成本极低，属跟随 P1 条目一起清理的附属点。
- 优化建议：随 body 回读逻辑一并删除；若保留 body 采样，可直接使用 `r.ContentLength`（`http.Request` 已解析的 `int64`），无需再 `Atoi`。

---

## 四、stores/null（含 stores/null/convert）

### P2

#### [P2] `stores/null/bool.go:103` 等 24 处：全部可空类型 `MarshalJSON` 空值分支调用 `json.Marshal(nil)` 产生反射序列化与额外分配，应返回包级共享的 `"null"` 字节切片

- 问题：每个类型在 `!Valid` 分支都写 `return json.Marshal(nil)`，该调用走 `encoding/json` 反射入口最终产出 `null`，伴随函数分发与字节分配。位置列表（`file:line`）：
  - `bool.go:103`、`int.go:98`、`int8.go:103`、`int16.go:101`、`int32.go:103`、`int64.go:92`
  - `uint.go:98`、`uint8.go:104`、`uint16.go:104`、`uint32.go:104`、`uint64.go:96`
  - `float32.go:97`、`float64.go:92`
  - `string.go:78`、`bytes.go:99`、`bytes.go:102`、`json.go:134`、`json.go:137`、`byte.go:104`
  - `time.go:60`、`datetime.go:64`、`date.go:64`、`ctime.go:63`、`timestamp.go:64`
- 影响评估：`MarshalJSON` 在结构体被 JSON 响应序列化时逐字段调用；一个包含多个可空字段的 DTO，每个空字段都独立触发一次 `json.Marshal(nil)`。反射入口在高频序列化场景下比直接返回静态字节明显更贵；单次几十到上百纳秒量级，字段多、QPS 高时累计可观。
- 优化建议：在 `bytes.go` 已有的包级变量 `NullBytes = []byte("null")`（`bytes.go:13`）基础上，空值分支直接 `return NullBytes, nil`（注意 `NullBytes` 为共享切片，调用方不应修改；`encoding/json` 不会修改返回的 Marshaler 输出）。

#### [P2] `stores/null/json.go:65` `JSON.Unmarshal` 与 `json.go:122` `JSON.Marshal`：多出一轮“序列化→反序列化”往返，可直接对原始字节操作

- 问题：`Unmarshal`（第 65–76 行）先 `json.Marshal(j)`（触发 `MarshalJSON` 取回内部原始 JSON 字节），再 `json.Unmarshal(res, dest)`；`Marshal`（第 122–129 行）先 `json.Marshal(obj)` 得到 `res`，再 `json.Unmarshal(res, j)` 由 `UnmarshalJSON` 复制进内部切片。两次辅助方法各自多走一遍编码/解码：`Unmarshal` 本可直接 `json.Unmarshal(j.JSON, dest)`；`Marshal` 本可直接 `j.JSON = res`。
- 影响评估：属于手工辅助方法，调用频次取决于业务是否在每次文档/对象转换时使用；一旦用于高频对象转换，额外一轮 json 编解码 + 中间字节分配即被放大。**实际调用频次需实测确认**。
- 优化建议：
  ```go
  func (j JSON) Unmarshal(dest any) error {
      if dest == nil { return errors.New("destination is nil, not a valid pointer to an object") }
      return json.Unmarshal(j.JSON, dest)
  }
  func (j *JSON) Marshal(obj any) error {
      res, err := json.Marshal(obj)
      if err != nil { return err }
      j.JSON = res
      j.Valid = true
      j.Set = true
      return nil
  }
  ```

#### [P2] `stores/null/string.go:80` `String.MarshalJSON`：纯字符串仍走 `json.Marshal` 全量编码器；各 `Result()` 方法对已知具体类型再走 cast 反射

- 问题：① `String.MarshalJSON`（第 76–81 行）有效值分支 `json.Marshal(s.String)`，对纯字符串也要经过 JSON 编码器的引号包裹与转义扫描，可用 `strconv.Quote` 直接产出带引号字节。② 各类型 `Result()` 对已经是具体类型的值调用 `cast.ToString` / `cast.ToString(...)`，cast 内部以 `reflect.TypeOf` 做类型分发：`string.go:155`、`bytes.go:166`、`json.go:201`、`byte.go:175`、`uint.go:162`、`uint8.go:168`、`uint16.go:168`、`uint32.go:168`、`uint64.go:172`。
- 影响评估：`MarshalJSON` 位于响应序列化路径；`Result()` 一般用于调试、日志或 SQL 文本拼接，频次较低，故整体定为 P2。
- 优化建议：字符串序列化改为 `return []byte(strconv.Quote(s.String)), nil`；`Result()` 直接使用 `strconv.FormatInt/FormatUint/FormatFloat` 或 `string(...)`，去掉 cast 依赖。

> 补充说明：`stores/null/convert/convert.go` 移植自标准库 `database/sql/convert.go`，常见类型（string、[]byte、time.Time、nil）在第 27–103 行先用类型开关直判直写，反射仅作为回退分支；`Scan` 路径整体结构与标准库一致，未发现需要上报的性能缺陷。

---

## 五、条目统计

### 按子包 × 优先级

| 子包 | P0 | P1 | P2 | 小计 |
|---|---:|---:|---:|---:|
| `stores/sql` | 0 | 3 | 3 | 6 |
| `stores/mongodb` | 0 | 3 | 3 | 6 |
| `stores/elasticsearch` | 0 | 2 | 2 | 4 |
| `stores/null`（含 `convert`） | 0 | 0 | 3 | 3 |
| **合计** | **0** | **8** | **11** | **19** |

### 条目清单索引

**stores/sql（P1）**
1. `interceptor.go:92` `TraceInterceptor`：共享 attrs 切片并发追加、无限增长
2. `logger.go:87` `WrapLogger.Trace` / `logger.go:164` `WrapEntry.Trace`：未判日志级别即构造 SQL、栈回溯、Fields 映射
3. `client.go:327` `Open`：连接池零值参数无兜底默认

**stores/sql（P2）**
4. `interceptor.go:101` `TraceInterceptor`：每次查询 Explain + 栈回溯 + metadata 拷贝
5. `client.go:301` `Open`：`NowFunc` 冗余 `.Local()`
6. `ext.go:19` 等 `Pagination/Scope/Sort` 取值方法：已知类型仍走 cast 反射

**stores/mongodb（P1）**
7. `clientmanager.go:173` `NewModel` / `clientmanager.go:197` `NewCollection`：每次创建先 Ping 一次网络往返
8. `collection.go:653` `FindMany`：手写反射逐文档解码 + 切片零容量反复扩容
9. `collection.go:829` `BeforeCreate` / `collection.go:871` `BeforeUpdate`：反射结构体→bson.M 二次转换

**stores/mongodb（P2）**
10. `collection.go:121` 等 15 个操作方法：重复 `NewTracer` / `Caller(6)` / `metadata.New(nil)`
11. `collection.go:146` 等：指标标签重复取 `Database().Name()` / `Name()`，未缓存
12. `collection.go:590` `FindOne` / `collection.go:637` `FindMany`：每次分配 options，Skip/Limit 恒传指针

**stores/elasticsearch（P1）**
13. `elastic.go:159` `RoundTrip`：无效的请求体全量回读与内存分配
14. `elastic.go:110` `ESTraceServerInterceptor`：零值 `http.Transport`，空闲连接与超时未配置

**stores/elasticsearch（P2）**
15. `elastic.go:138` `RoundTrip`：每次重复创建 Tracer、栈回溯、metadata 分配
16. `elastic.go:158` `RoundTrip`：重复 `strconv.Atoi` 解析 Content-Length

**stores/null（P2）**
17. `bool.go:103` 等 24 处：空值 `MarshalJSON` 调用 `json.Marshal(nil)`
18. `json.go:65` `JSON.Unmarshal` / `json.go:122` `JSON.Marshal`：多余编解码往返
19. `string.go:80` `String.MarshalJSON` + 各 `Result()`：字符串编码与 cast 反射

---

## 六、本范围内高价值候选（供全局 Top 10 评选）

| # | 位置 | 一句话问题 | 预期收益 |
|---|---|---|---|
| 1 | `stores/mongodb/collection.go:653-671` `FindMany` | 放弃 `cursor.All`，手写反射逐文档 `reflect.New`/`reflect.Append`，切片从零扩容 | 列表接口降低每条文档的反射分配与多次切片拷贝；结果集越大收益越明显，高 QPS 读接口 CPU 与 GC 压力显著下降 |
| 2 | `stores/mongodb/collection.go:829-867` / `871-901` `BeforeCreate`/`BeforeUpdate` | 写前反射结构体转 `bson.M`，驱动再编码一遍，字段逐次装箱、`FieldByName` 重复查找 | 写路径去掉一遍结构体反射遍历与 map/装箱分配；批量插入收益按文档数放大，可消除一整层重复反射 |
| 3 | `stores/sql/logger.go:87,89-91,164,166-168` `WrapLogger.Trace`/`WrapEntry.Trace` | 每条 SQL 在判定日志级别前先构造完整 SQL、做 `runtime.Caller` 栈回溯、分配 Fields 映射 | 生产环境 Debug 关闭时直接省去每查询一次栈回溯与一次 SQL 字符串分配，覆盖全部 SQL 读/写/Raw 操作 |
| 4 | `stores/elasticsearch/elastic.go:159-165` `RoundTrip` | 每个写请求全量回读并拷贝未被使用的请求体 | 删除纯无效功：每条 ES 写请求省一次 body 全量读入与一次缓冲分配，降低内存占用与拷贝 CPU |
| 5 | `stores/mongodb/clientmanager.go:173,197` `NewModel`/`NewCollection` | 每次取集合句柄前先 Ping 一次服务端往返 | 若按请求取句柄，每个请求去掉一次串行 RTT；改为启动期校验后热点路径延迟直接下降一个网络往返量级 |
| 6 | `stores/sql/client.go:327-330` `Open` | 未配置时连接池零值直接下发，空闲连接为 0 导致用完即关、最大连接无上限 | 恢复连接复用，消除反复 TCP/认证重建；防止高并发下连接失控，属于稳定性 + 延迟双重收益 |
| 7 | `stores/elasticsearch/elastic.go:110` `ESTraceServerInterceptor` | 零值 `http.Transport` 每节点仅 2 个空闲连接、无拨号/响应头超时 | 提升 ES 连接复用率，消除空闲连接频繁关闭重建；超时配置避免网络抖动时 goroutine 长时间悬挂 |
| 8 | `stores/sql/interceptor.go:92,103-107` `TraceInterceptor` | 闭包共享 attrs 切片被并发 append 并终身增长 | 消除追踪路径上的切片竞争、无限增长内存泄漏与重复扩容；修复后每个 Span 使用独立属性切片 |
| 9 | `stores/mongodb/collection.go:121` 等 15 处 | 每个操作方法内部重复 `NewTracer`、`Caller(6)`、`metadata.New(nil)` | 追踪开启时每操作省去 tracer 包装/空 metadata 分配与一次栈遍历；集中到构造期复用，改造后维护成本同步下降 |
| 10 | `stores/null` 24 处 `MarshalJSON` 空值分支 | 空字段经 `json.Marshal(nil)` 反射产出 `"null"` | 序列化路径每空字段省去一次 json 反射入口与分配；DTO 空字段多时 GC 对象数下降，改动小、风险低 |

---

## 七、拿不准、需实测验证的条目

1. **`stores/sql/client.go:344` `SqlConn` 的调用频次**：若业务在启动期缓存 `*gorm.DB`，缓存管理器查找路径开销可忽略；若每个请求都调 `SqlConn()`，每次获取都要经过 `core/resource` 的 singleflight 与 RWMutex（见跨模块提示 1）。需通过业务代码调用链或 pprof `mutex`/`alloc` 采样确认该路径是否落在热路径上。
2. **`stores/mongodb/clientmanager.go:173,197` Ping 是否落在请求热路径**：取决于 `NewModel`/`NewCollection` 是启动期一次性创建还是按请求创建；建议在两种使用方式下分别压测 `NewModel` 开销，再决定是否降级条目优先级。
3. **`stores/sql/logger.go:87` `fc()` 的实际分配规模**：`Explain` 生成 SQL 字符串的成本随参数数量与 SQL 文本长度变化；建议用 `benchstat` 对比 Debug 关闭时 Trace 改造前后的 `allocs/op` 与 `ns/op`，量化栈回溯（`FileWithLineNum`）与 SQL 拼装各自占比。
4. **`stores/mongodb/collection.go:829` `BeforeCreate`/`871` `BeforeUpdate` 的输入类型分布**：若业务大多直接传 `bson.M`/`map`，会走 `default` 分支（第 859–865、898–900 行），反射开销很小；若以结构体为主，第 2、9 条候选收益成立。需结合实际入参统计。
5. **`stores/null/json.go:65,122` `JSON.Unmarshal`/`Marshal` 的调用频次**：该辅助方法若仅在低频配置/导入场景使用，P2 评级保持；若用于高频对象转换，可上调为 P1。
6. **`stores/elasticsearch/elastic.go:110` 连接池瓶颈**：零值 transport 每节点 2 空闲连接是否真正造成排队，取决于 ES 并发度与服务端响应时长；建议压测观察 `net/http` 连接新建次数（`Dial` 计数）与等待空闲连接耗时后再定优先级。
7. **`stores/mongodb/collection.go:659` 切片扩容倍率与 `cursor.All` 差距**：需在相同结果集规模（如 100/1000/5000 条）下 benchmark 手写反射实现与 `cursor.All` 的分配数与耗时，作为改造前后对比基线。

---

## 八、跨模块提示（不在本审查范围内，仅登记线索）

1. **`core/resource/sourcemanager.go:50` `GetResource` 与 `core/resource/singleflight.go:36` `Do`**：即使资源已缓存命中，每次调用仍要走 singleflight 互斥锁 `createCall`（锁 + map 查找 + 新建 `call`/`WaitGroup`）→ 内部 RLock 查 map → `makeCall` 再次加锁删除 key → `wg.Done()`；稳态命中每次产生两次互斥锁临界区、两次 map 操作与一次 `call` 堆分配。本范围内 `stores/sql/client.go:359`、`stores/mongodb/clientmanager.go:108`、`stores/elasticsearch/elastic.go:64` 均经此路径获取缓存连接。建议资源管理器为已存在资源提供无 singleflight 的快路径（如 `sync.Map` 或只读 map + RLock），由该模块的审查方评估。
2. **`core/metric/counter.go:52` `Inc(labels ...string)` 与 `core/metric/histogram.go:42` `Observe(labels ...string)`**：可变长标签每次调用都构造一个 `[]string` 切片再传入 `WithLabelValues`；mongo/sql/es 每个操作上报均经历一次。建议模块内预建标签切片或提供定长标签 API，由该模块审查方评估。
3. **`core/trace/trace.go:38` `NewTracer`**：每次调用新建包装结构并 `otel.Tracer("ratel")`；stores 三个客户端都在操作方法内重复调用，建议在客户端/transport 构造期创建一次（见本报告第 10、15 条）。
4. **`core/call/call.go:9` `Caller`**：`runtime.Caller` + `runtime.FuncForPC(pc).Name()` 栈回溯被 sql/mongo/es 的追踪路径每次操作调用；是否可复用 GORM 日志已采集的源码位置，由 trace/logger 模块统一评估。
5. **`stores/mongodb/collection.go:468-472` `UpdateOne`**：指标标签只传了 `("mongodb", "ERR"/"OK")`，与其余方法的四标签口径不一致，属于指标标签基数/口径问题，非性能问题，仅登记提示，不在本报告条目内展开。
