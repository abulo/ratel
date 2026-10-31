# 性能审查报告：toolkit（代码生成 / 脚手架 CLI）

> 审查范围（仅以下目录及其下所有 `.go` 文件，测试文件 `*_test.go` 不在性能审查范围内）：
>
> - `toolkit/` 根目录：`main.go`
> - `toolkit/api/`：`api.go`、`gin.go`、`hertz.go`、`doc.go`
> - `toolkit/base/`：`init.go`、`install.go`、`install_compatible.go`、`mod.go`、`path.go`、`query.go`、`repo.go`、`util.go`、`vcs_url.go`、`doc.go`
> - `toolkit/build/`：`build.go`、`doc.go`
> - `toolkit/change/`：`change.go`、`doc.go`
> - `toolkit/dao/`：`dao.go`、`doc.go`
> - `toolkit/module/`：`module.go`、`tpl_convert.go`、`tpl_module.go`、`tpl_proto.go`、`tpl_service.go`、`doc.go`
> - `toolkit/project/`：`add.go`、`backend.go`、`new.go`、`doc.go`
> - `toolkit/upgrade/`：`upgrade.go`、`doc.go`
> - `toolkit/vue/`：`interface.go`、`method.go`、`page.go`、`vue.go`、`doc.go`
>
> 审查方式：逐文件静态通读（非测试源码共 38 个文件、6451 行），每条目回读原文核对 `文件:行号` 与函数名；被调用方实现位于范围之外（`util/`、`stores/sql/`、`watch/`）的，正文只锚定范围内的调用点，根因统一在文末「跨模块提示」中单列，不在本章节展开审查。
>
> 总体结论：`toolkit` 是一次性 CLI 代码生成工具，**不存在服务端热路径**，因此本章节 **P0 = 0**。问题集中在四类：① 批量生成时重复执行模板解析 / 类型映射表构建；② 渲染期对同一字段名、同一表名反复做驼峰转换；③ 成员判断使用反射型 `util.InArray` 做线性扫描；④ 数据库元数据查询与 git 子进程串行发起。多数条目为 P2，单表 / 少量标签场景下耗时以模板外部进程（`go fmt`、`goimports`、`protoc`、`git clone/pull`、数据库 RTT）为主，CLI 自身 CPU 与内部分配占比需要 benchmark / memprofile 实测确认（见文末清单）。

---

## 一、toolkit 根目录（main.go）

通读 `toolkit/main.go:22-56`：仅做 cobra 命令注册与 `rootCmd.Execute()`，无循环、无 I/O、无缓冲 / 分配热点，**本目录无性能条目**。

---

## 二、toolkit/api

### P2

#### [P2] `toolkit/api/api.go:92` `Run`（同 `api.go:101` `TableItem`、`api.go:106` `TableIndex`、`api.go:111` `TablePrimary`）：4 次 information_schema 元数据查询严格串行，网络 RTT 依次累加

- 问题：`Run` 在一次命令内顺序执行 4 次独立的数据库元数据查询：`TableColumn`（`api.go:92`）→ `TableItem`（`api.go:101`）→ `TableIndex`（`api.go:106`）→ `TablePrimary`（`api.go:111`）。4 条 SQL 之间无数据依赖（后续拼装只使用各自返回结果），但代码中是上一条 `err != nil` 返回后才发起下一条，无法利用数据库连接池的并发能力。
- 影响评估：单次命令固定叠加 4 次 RTT。本地 / 同机房 MySQL 单次 information_schema 查询通常为毫秒级，影响不明显；当数据库在远端（跨网段、容器远端、只读副本）或 information_schema 缓存未命中时，串行 RTT 之和会成为交互式命令可感知的等待。量级判断：耗时 = 4 × 单次 RTT + 4 × 查询执行时间；具体占比需实测（见需实测清单第 1 项）。同类串行模式同时存在于 `toolkit/module/module.go:88,97,102,107` 与 `toolkit/vue/vue.go:71,81,86,91`，此处合并记录一次。
- 优化建议：
  1. 用 `errgroup.Group` 并发发起 4 个只读查询（GORM 句柄在连接池之上支持并发使用），任一失败则取消；失败处理与现有 `err != nil` 直接返回保持一致。
  2. 更彻底的方向：`TableItem` 与 `TablePrimary` 的元数据可在 `TableColumn` / `TableList` 结果上合并为一条 SQL，减少往返次数（该 SQL 调整落在 `toolkit/base/query.go`）。

#### [P2] `toolkit/api/api.go:157` `Run`（共享模式，遍布 `api.go`、`module.go`、`vue.go` 三处命令装配流程）：同一表名 `base.CamelStr(tableItem.TableName)` 在单次命令中重复计算 12 次以上

- 问题：`Run` 内每追加一个 `base.Method` 与一条 `multiSelect` 候选项，就重新调用一次 `base.CamelStr(tableItem.TableName)`。`api.go` 内的调用点为：`157`、`170`、`183`、`196`、`208`、`209`、`210`、`211`、`219`、`231`、`239`、`251`、`255`、`303`、`331`；`toolkit/module/module.go` 同构调用点：`151`、`164`、`177`、`190`、`202`、`203`、`204`、`205`、`214`、`226`、`233`、`245`、`249`、`297`、`322`；`toolkit/vue/vue.go` 同构调用点：`132`、`143`、`154`、`165`、`175`、`176`、`177`、`178`、`187`、`197`、`204`、`214`、`218`、`264`、`291`。输入 `tableItem.TableName` 在整个流程中不变，但每次 `CamelStr`（`toolkit/base/util.go:12-16`）都执行两次 `strings.Replace` 全串扫描 + 一次 `util.UCWords` + 再一次 `strings.Replace`，每次都产生新字符串分配。
- 影响评估：单次命令约 12~16 次重复转换，单次转换为微秒级字符串处理，累计分配数十个临时字符串；不属于耗时主因，但属于确定性可消除的重复计算。表名越长、方法越多，重复次数线性增长。
- 优化建议：在 `Run` 进入方法装配前计算一次 `tableCamel := base.CamelStr(tableItem.TableName)`，后续 `Name: tableCamel + "Create"` 等处直接复用局部变量；`module.go`、`vue.go` 两处做相同改造。

#### [P2] `toolkit/api/api.go:121` `Run`（共享模式，遍布 `api.go`、`module.go`、`vue.go`）：索引 / 方法 / 模板去重与筛选全程使用反射型 `util.InArray` 做线性查找

- 问题：装配流程中所有成员判断都走 `util.InArray`（实现位于 `util/array.go:223-243`，见跨模块提示第 2 条），每次调用都先 `reflect.ValueOf`，再逐元素 `reflect.DeepEqual`。范围内调用点：
  - 索引字段去重（`selectIndex`）：`api.go:121`、`module.go:116`、`vue.go:100`；
  - 索引字段过滤（`useIndex`）：`api.go:282`、`module.go:274`、`vue.go:243`；
  - 条件字段过滤（`field`）：`api.go:289`、`module.go:281`、`vue.go:250`；
  - 索引类型白名单（`idx = ["Item","List"]`）：`api.go:300`、`module.go:294`、`vue.go:261`（固定 2 元素切片，仍走反射）；
  - 方法名按勾选结果筛选：`api.go:372`、`api.go:375`、`module.go:364`、`vue.go:343`、`vue.go:346`；
  - 模板类型勾选判断：`module.go:382,386,389,392`、`vue.go:365,374,383,387,391`；
  - 多表参数过滤：`toolkit/dao/dao.go:66`。
- 影响评估：切片长度普遍为个位数到十几个元素，单次判断开销为「反射包装 + 线性比较」，量级为微秒级；在方法 / 索引多的宽表场景，外层循环 × 内层线性扫描形成 O(n·m)，且每次反射包装都产生分配。整体不是命令耗时主因，但改造方式统一、成本低。
- 优化建议：在 `Run` 内把需要反复查找的勾选集合（`multiSelected`、`multiPageSelected`、`useIndex`、`tplSelected`、`tables`）构建为 `map[string]struct{}`，循环内做 O(1) 原生 map 查找；固定白名单 `idx = []string{"Item","List"}` 直接改为 `map[string]struct{}` 或 `switch`。

> `toolkit/api/gin.go`、`toolkit/api/hertz.go`：两文件为模板字符串（`GinTemplate:15-444`、`HertzTemplate:15-436`），仅在 `Generate`（`api.go:401-419`）中解析一次；模板内 `{{CamelStr .Table.TableName}}` 等函数调用的重复开销归入 `toolkit/base` 章节字段级转换条目。两文件自身无循环、无独立分配热点。

---

## 三、toolkit/base

### P1

#### [P1] `toolkit/base/util.go:554` `Props`（同 `util.go:573` `PropsSearch`）：条件字段成员判断在外层字段循环内逐次调用反射型 `util.InArray`，构成 O(n·m) 反射线性扫描

- 问题：两个函数都先把 `condition` 展开成字符串切片 `conditionList`（`util.go:555-558`、`util.go:574-577`），随后遍历整张表字段 `table`，对每个字段调用 `util.InArray(item.ColumnName, conditionList)`（`util.go:561`、`util.go:580`）判断是否为查询条件。`util.InArray` 内部是 `reflect.ValueOf` + 逐元素 `reflect.DeepEqual`（`util/array.go:223-243`），没有把 `conditionList` 预建成集合。
- 影响评估：复杂度为 O(表字段数 n × 条件字段数 m)，每次比较附带反射包装分配。普通表字段 10~30 个、条件字段 2~6 个，单次渲染开销在百微秒量级；宽表（50+ 字段、多组合索引条件）下扫描次数与反射分配同比例放大。`Props` 是活路径（`toolkit/vue/page.go:287` 在表格列片段模板中调用）；`PropsSearch` 当前仅在 `page.go:41` 的 FuncMap 中注册、模板正文尚未调用，属同一实现模式的备用路径，改造时一并处理即可。
- 优化建议：进入循环前将条件字段构建为 `map[string]struct{}`：
  ```go
  condSet := make(map[string]struct{}, len(condition))
  for _, c := range condition { condSet[c.ColumnName] = struct{}{} }
  // 循环内：if _, ok := condSet[item.ColumnName]; ok { ... }
  ```
  复杂度降为 O(n+m)，同时去掉每次判断的反射开销。

#### [P1] `toolkit/base/util.go:88` `ModuleDaoConvertProto`（同 `util.go:373` `ApiToProto`、`util.go:777` `ModuleProtoConvertDao`，以及 `util.go:72` `Convert`、`util.go:286` `ProtoRequest`、`util.go:316` `TypeScriptCondition`、`util.go:529` `ModuleProtoConvertMap`、`util.go:542` `Rule`、`util.go:589` `Json`、`util.go:682` `TypeScript`）：同一字段在单个 helper 内反复调用 `CamelStr` / `Helper`，转换结果未在元数据加载阶段缓存

- 问题：每个代码片段 helper 遍历 `Condition` / `TableColumn` 一次，而在每个字段的各 `switch` 分支里，`CamelStr(item.ColumnName)` 与 `Helper(item.ColumnName)` 被重复调用 2~4 次。典型位置：
  - `ModuleDaoConvertProto`（`util.go:88-283`）：每个可空 / 非空分支的 `if` 判断行与赋值行各调一次 `CamelStr`（如 `94`、`96`、`101`、`103` …… 全函数实测 78 处 `CamelStr` 调用）；
  - `ApiToProto`（`util.go:373-526`）：每分支 `Helper(item.ColumnName)` 一次、`CamelStr(item.ColumnName)` 1~2 次（如 `380`、`382`、`385`，`424`、`426`、`429` ……）；
  - `ModuleProtoConvertDao`（`util.go:777-972`）：每分支 `CamelStr` 出现 3 次（判断、目标字段、源字段，如 `784`、`786`、`791`、`793` ……）；
  - 其余 helper 每个字段至少 1~2 次 `Helper` 调用（`ProtoRequest:290`、`TypeScriptCondition:322-363`、`ModuleProtoConvertMap:532-534`、`Rule:546`、`Json:595-673`、`TypeScript:688-766`）。
  `Helper`（`util.go:19-22`）内部先调 `CamelStr` 再取首字母小写，即一次 `Helper` 已经包含一次完整 `CamelStr`；同一字段在同一分支内既调 `Helper` 又调 `CamelStr`，等于同一输入转换两遍。转换的输入在元数据加载后就固定不变，但没有任何缓存字段。
- 影响评估：一次完整 `mp` / `api` / `vue` 命令会串行触发 4~10 个 helper，每个 helper 遍历全部字段；宽表 30~50 字段时，渲染期 `CamelStr`/`Helper` 调用量级达到数百到上千次，每次包含两次 `strings.Replace` 全串扫描、一次 `util.UCWords`（内部每次新建 `cases.Title` 分词器，见跨模块提示第 1 条）与若干字符串分配。单次 CLI 运行下这部分累计 CPU / 分配在百毫秒以内量级（**需 benchmark 实测确认**），但属于纯重复劳动，消除后所有渲染 helper 同时受益。
- 优化建议：
  1. 在 `toolkit/base/query.go` 加载字段元数据时（`TableColumn:212-242`、`TablePrimary:264-289`），把 `CamelStr(ColumnName)` 与 `Helper(ColumnName)` 计算一次，挂到 `Column` 新增的附加字段（如 `CamelName`、`LowerName`，`Column` 定义见 `query.go:24-41`）；模板与 helper 直接读取预计算值。
  2. 短期不改结构体时，至少在每个 helper 的 `for _, item := range` 循环顶部计算一次 `fieldCamel := base.CamelStr(item.ColumnName)`、`fieldHelper := base.Helper(item.ColumnName)`，分支内复用局部变量，避免单字段 2~4 次重复转换。
  3. `util.UCWords` 每次重建分词器的根因在 `util` 包，列入跨模块提示。

### P2

#### [P2] `toolkit/base/query.go:234` `TableColumn`（同 `query.go:283` `TablePrimary`）：类型映射表 `NewDataType()` 在每次元数据加载时整体重建

- 问题：`TableColumn` 在扫描结果回填循环内调用 `dataType := NewDataType()`（`query.go:234`），`TablePrimary` 同样调用 `dataType := NewDataType()`（`query.go:283`）。`NewDataType`（`query.go:128-169`）每次都 `make(map[string]DataType)` 并逐条赋值 33 个类型映射条目，返回全新 map；表字段类型在程序运行期间永不变化。
- 影响评估：`toolkit/dao` 批量生成场景（`dao.go:64-75` 逐表循环）下，每张表触发一次 `TableColumn`（表多时还会伴随 `TablePrimary`），映射表被整体重建 N 次；单次重建为几十次 map 插入与结构体字面量分配，累计分配随表数线性增长。单表命令影响可忽略。
- 优化建议：将映射表提为包级只读变量（`var dataTypeMap = NewDataType()`，或直接用包级 `map` 字面量初始化），`TableColumn` / `TablePrimary` 只读引用；注意现有注释 `NewDataType` 已声明「调用方不应修改返回结果」，提为全局只读表与现有语义一致。

#### [P2] `toolkit/base/query.go:229` `TableColumn`：对同一结果切片执行两遍独立循环，类型归一化与类型映射 / 位置序号回填可合并

- 问题：`query.go:229-233` 第一遍循环只做 `SqlResultConvert` 归一化；`query.go:235-239` 第二遍循环再挂 `DataTypeMap`、写 `PosiTion`，其中 `PosiTion = cast.ToInt64(newKey)`（`query.go:238`）每个字段走一次 `cast` 转换。两遍循环读取同一份 `res`，循环本身开销很小，但属于可合并的不变量重复遍历。
- 影响评估：字段数为 n 时多一次 O(n) 遍历 + n 次 `cast.ToInt64`；量级可忽略，仅作为代码简化项记录。
- 优化建议：合并为单循环；位置序号直接使用 `int64(k+1)`，去掉 `cast.ToInt64` 包装。

#### [P2] `toolkit/base/util.go:73` `Convert`（同 `ModuleDaoConvertProto:89`、`ApiToProto:374`、`ModuleProtoConvertDao:778`、`Props:559`、`Rule:543`、`Json:590`、`TypeScript:683` 等全部 strings.Builder helper）：逐字段 `builder.WriteString(fmt.Sprintf(...))`，静态文本与动态插值混在 Sprintf 中，且 builder 未预估容量

- 问题：代码片段 helper 普遍写法为每段代码两行：先 `builder.WriteString(fmt.Sprintf("	...%s...", 动态值))`，再 `builder.WriteString("\n")`。`fmt.Sprintf` 即使格式串只插值 1~2 个参数，也要走格式解析与参数反射式取值并分配临时字符串；静态骨架文本本可直接 `WriteString`。同时 `strings.Builder` 以零容量起步，宽表输出数百行代码时按 2 倍扩容多次搬迁底层缓冲。
- 影响评估：每条 Sprintf 产生 1 个临时字符串 + 1 次 builder 写入；宽表 × 多 helper 时累计数百次 Sprintf 分配。builder 扩容为 O(已写入字节) 的均摊成本，总拷贝量不高。属于确定性的分配冗余，非耗时主因。
- 优化建议：
  1. 纯静态骨架直接 `builder.WriteString("...")`；仅动态片段使用 `Sprintf`，或改用 `strconv` / 直接 `builder.WriteString(prefix); builder.WriteString(value); builder.WriteString(suffix)`。
  2. helper 入口按字段数预估容量：`builder.Grow(estimate)`（如 `len(condition) * 128`），减少扩容搬迁。

#### [P2] `toolkit/base/path.go:40` `copyFile`（由 `copyDir:61-94` 递归逐文件调用）：替换字面量的 `[]byte` 转换在每个文件内重复进行，且每轮替换都对整个文件缓冲做一次全量扫描；模板文件整体读入内存

- 问题：
  1. `copyFile` 接收 `replaces []string`（`old, new` 成对），循环内 `bytes.ReplaceAll(buf, []byte(old), []byte(next))`（`path.go:56`）——`[]byte(old)` / `[]byte(next)` 在**每个被复制文件**上重新分配并完成一次字符串→字节拷贝；同一组替换字面量在整个目录树内重复转换数百次。
  2. 每一对 `old→new` 都对整份 `buf` 做一次独立全串扫描，替换对数为 k 时总扫描量为 k × 文件大小。
  3. 函数先用 `os.ReadFile(src)` 把整文件读入内存（`path.go:46`），再整体 `os.WriteFile`（`path.go:58`）；`CopyTo` / `CopyFrontTo` / `CopyToV2`（`toolkit/base/repo.go:108-139`）复制整套 layout 模板仓库，前端模板可能包含较大的二进制资源。
- 影响评估：模板源码普遍为 KB 级文本，替换扫描与内存占用不构成瓶颈；当模板仓库混入较大二进制文件时（二进制不应做文本替换，当前代码仍会整文件读入并尝试 `ReplaceAll`），内存峰值与无效扫描随资源大小上升。
- 优化建议：
  1. 在 `copyDir` 入口把 `replaces` 预转换为 `[][]byte`（`oldB, newB []byte` 成对），`copyFile` 直接使用，避免逐文件重复分配；
  2. 对非文本后缀（图片、字体、压缩包等）跳过 `ReplaceAll`，直接流式拷贝（`io.Copy` + `bufio`）；文本文件在替换前可先判断 `bytes.Contains(buf, oldB)` 再决定是否扫描，减少无命中文件的全量扫描；
  3. 较大文件改用 `io.Copy` 流式读写，避免整文件驻留内存。

#### [P2] `toolkit/base/util.go:60` `InMethod`：对方法切片线性扫描，内部再次委托反射型 `util.InArray`；被 Vue 页面模板在单次渲染中反复调用

- 问题：`InMethod(arr []Method, target ...string)`（`util.go:60-69`）外层 `for _, item := range arr` 逐个方法调用 `util.InArray(item.Type, target)`（`util.go:63`）。`util.InArray` 又对 `target` 切片做反射线性扫描，即判定一次「某类型是否存在于方法列表」为 O(方法数 × 目标类型数) 且附带反射包装。该函数在 `toolkit/vue/page.go` 模板中被反复调用（调用位置见 `toolkit/vue` 章节条目）。
- 影响评估：方法列表通常 ≤10 个、目标类型每次 1 个，单次判断为微秒级；问题在于模板同一轮渲染重复判定同一批方法十几次，每次都重新线性扫描 + 反射包装，结果在渲染期间恒定不变。
- 优化建议：在组装 `moduleParam.Method` 后（或 `GeneratePage` 入口）构建一次 `map[string]struct{}`（方法类型集合），模板函数改为 `func InMethodSet(methods map[string]struct{}, target string) bool` 做 O(1) 查找；或在 `ModuleParam` 上预计算类型集合字段，模板直接读取。

#### [P2] `toolkit/base/mod.go:24` `ModuleVersion`：每次调用都派生 `go mod graph` 子进程并把整份依赖图读入内存后逐行扫描

- 问题：`ModuleVersion`（`mod.go:24-44`）执行 `go mod graph`，整份输出写入 `bytes.Buffer`（`mod.go:25-31`），再用 `bufio.Reader` 逐行 `ReadLine` 扫描匹配目标模块（`mod.go:32-43`）。`RatelMod`（`mod.go:47-59`）在一次命令流程中调用一次；`go mod graph` 输出规模随依赖数增长，且子进程冷启动本身耗时明显。
- 影响评估：单次命令一次调用，属于进程外部耗时（`go` 子进程 + 依赖图规模），Go 侧缓冲与逐行扫描占比很小；记录为低优先级，避免与外部子进程耗时混淆。
- 优化建议：可改用 `go list -m -f '{{.Path}} {{.Version}}' github.com/abulo/ratel/v3` 定向查询，避免拉取整张依赖图；结果在 `RatelMod` 生命周期内缓存一次，避免重复调用。

#### [P2] `toolkit/base/install.go:15` `GoInstall`（`install_compatible.go:14` 回退版同构）：多个工具串行 `go install`，任一步等待网络 / 编译完成后才开始下一个

- 问题：`GoInstall`（`install.go:15-29`）`for _, p := range path` 循环内顺序 `exec.Command("go", "install", p).Run()`，每个工具独立联网拉取与编译；`toolkit/upgrade/upgrade.go:30-37` 一次性传入 6 个工具（toolkit、protoc-gen-go、protoc-gen-go-grpc、protoc-gen-openapi、go-mod-upgrade、protoc-go-inject-tag），总等待为 6 次串行安装之和。
- 影响评估：CPU 与内存分配不是问题，耗时主要是网络下载与编译；串行导致总时长 ≈ Σ 单工具时长，并行可压缩为近似最大值。
- 优化建议：用 `errgroup.Group` 限定并发度（如 2~3，避免同时打满网络与本地编译资源）并行执行安装命令，汇总错误信息；保持现有「任一失败返回错误」的语义。

---

## 四、toolkit/build

通读 `toolkit/build/build.go:47-117`：入口文件列表经校验后，`errgroup.Group`（`build.go:109-116`）已并行启动各 `watch.Watch`，监听配置按入口逐项构造，无串行热点、无重复分配；真正的文件监听与编译逻辑在范围外 `watch` 包（见跨模块提示第 3 条）。**本包范围内无独立性能条目。**

---

## 五、toolkit/change

### P1

#### [P1] `toolkit/change/change.go:117` `classifyCommits`：正则表达式在函数体内 `MustCompile`，每个发布标签分类时重复编译一次

- 问题：`re := regexp.MustCompile(...)`（`change.go:117`）位于 `classifyCommits` 函数内部；`Run` 在标签循环中每个标签调用一次 `classifyCommits`（`change.go:233`）。即仓库有 N 个标签，该正则（提交类型 + 可选修饰 + 消息 + 7~40 位哈希的分组模式）就被完整编译 N 次，编译产物在函数返回后即丢弃。
- 影响评估：Go 正则编译包含语法解析、NFA 构建与优化，单次编译开销显著高于一次匹配；标签数量随项目发布历史增长（数十到上百个 tag 很常见），N 次重复编译是确定性可消除的固定成本。匹配本身在提交条数级别，不是瓶颈；重复编译才是。
- 优化建议：将正则提为包级 `var`，`init` 或声明处 `regexp.MustCompile` 一次，`classifyCommits` 内直接复用（与 `toolkit/base/vcs_url.go:13` 已采用的包级 `scpSyntaxRe` 模式保持一致）。

### P2

#### [P2] `toolkit/change/change.go:98` `classifyCommits`：16 个分类桶切片在每次调用时整体重建

- 问题：函数开头用字面量初始化 `classified map[string][]string`（`change.go:98-115`），一次性分配 16 个空切片头并写入 map；每次标签调用都重建同一份结构。分类键集合在程序内固定。
- 影响评估：单次调用为 1 个 map + 16 个空切片的分配，开销很小；标签多时重复 N 次，属于确定性冗余分配。
- 优化建议：分类键列表与桶顺序提为包级常量；每次调用只需 `classified := make(map[string][]string, 16)`（或复用一个预定义桶顺序切片），分类追加逻辑不变。

#### [P2] `toolkit/change/change.go:211` `Run`：分类输出顺序切片 `categories` 在标签循环内部重复声明分配

- 问题：`categories := []string{...16 个分类键...}`（`change.go:211-228`）位于 `for i := range tags` 循环体内，每个标签都重新构造同一字面量切片，仅用于后续按固定顺序遍历 `classifiedCommits`。
- 影响评估：每标签一次 16 元素切片分配，量级很小；与上一条目同源，合并改造即可。
- 优化建议：`categories` 提为包级只读切片变量，循环内直接遍历。

#### [P2] `toolkit/change/change.go:193` `Run`（同 `change.go:207` `getCommits`）：每个标签串行派生两个 git 子进程（取日期 + 取区间提交），标签之间同样串行

- 问题：标签循环内先 `getTagDate(tag1)`（`change.go:193`，子进程见 `getTagDate:42-52`），返回后再 `getCommits(tag1, tag2)`（`change.go:207`，子进程见 `getCommits:56-71`）；标签之间也是严格串行。两次 git 子进程都要经历进程冷启动、仓库对象库读取；标签数为 N 时共 2N 次子进程派生，串行等待累加。
- 影响评估：耗时主体是 git 子进程本身（仓库规模、对象库冷热），Go 侧缓冲开销可忽略；大仓库 + 多标签时串行等待明显。**是否值得并行化需实测**：git 命令在同一仓库内并发读取通常安全，但要保证最终 Markdown 输出顺序仍按标签倒序。
- 优化建议：先批量并发取各标签日期与提交（`errgroup` 限定并发度），结果按原标签顺序写回；或同一标签的 `getTagDate` 与 `getCommits` 并发发起；汇总阶段再按 `categories` 固定顺序输出，不改变文档结构。

---

## 六、toolkit/dao

### P1

#### [P1] `toolkit/dao/dao.go:100` `GenerateDao`：模板字符串获取、FuncMap 注册与 `template.Parse` 在逐表循环中每张表重复执行一次

- 问题：`Run` 逐表循环（`dao.go:64-75`）内对每张表调用 `GenerateDao`（`dao.go:74`）；`GenerateDao` 每次执行：
  1. `DaoTemplate()`（`dao.go:105`，模板函数体 `dao.go:123-157`）——每次重新检查 `base.Config.Exists`、`FileExists`，自定义模板场景下还会重新 `FileGetContents`；
  2. `template.FuncMap{...4 个函数...}` 重新构造（`dao.go:100-104`）；
  3. `template.New(...).Parse(...)`（`dao.go:100-105`）重新完成词法分析、语法树构建与模板编译。
  内置模板字符串对所有表完全相同，解析结果与 FuncMap 在批量生成期间恒定不变，却被重复编译 N 次。
- 影响评估：`text/template` 解析编译的 CPU 开销显著高于一次渲染；批量生成几十张表时，重复解析的累计时间与表数成正比，是 dao 命令内 Go 侧最明显的可消除重复劳动。具体占比（相对于每张表的 DB 查询、文件写入、末尾一次 `go fmt` / `goimports`）需计时实测（见需实测清单第 2 项）；内置模板场景收益最确定。
- 优化建议：
  1. 在 `Run` 循环外完成一次 `DaoTemplate()` 读取与 `template.Must(...Parse(...))`，`GenerateDao` 接收预解析好的 `*template.Template`，循环内仅 `Execute`；
  2. FuncMap 提为包级变量复用；自定义模板文件在命令开始时读取一次并缓存，避免逐表重复 `FileExists` / 读文件。

### P2

#### [P2] `toolkit/dao/dao.go:70` `Run`：批量生成时逐表串行发起 `TableColumn` 元数据查询，表多时 RTT 串行累加

- 问题：循环 `dao.go:64-75` 中每张表先 `base.TableColumn(...)`（`dao.go:70`），成功后才 `GenerateDao`；表之间无数据依赖。与 `toolkit/api` 章节串行查询条目同源，但这里是「跨表串行」，表数量级更大。
- 影响评估：本地库影响有限；远端库 + 多表时总耗时 ≈ 表数 × 单次 RTT。模板渲染与文件写入可在查询之后串行，瓶颈首先在 DB 往返。**并发收益需结合连接池上限实测**。
- 优化建议：用带缓冲的工作池 / `errgroup` 限制并发度（如 4~8，匹配 `db.MaxOpenConns` 配置，见 `toolkit/base/init.go:108-110`）并发取字段元数据，渲染与写文件按表串行或同样受限并发；保持现有「单表失败跳过、不中断其余表」语义（`dao.go:71-73`）。

---

## 七、toolkit/module

### P2

#### [P2] `toolkit/module/module.go:88` `Run`（同 `module.go:97`、`module.go:102`、`module.go:107`）：4 次元数据查询串行发起

- 问题：与 `toolkit/api/api.go:92,101,106,111` 同一模式：`TableColumn` → `TableItem` → `TableIndex` → `TablePrimary` 顺序执行，互不依赖。改造方案见 `toolkit/api` 章节第一条目，此处不重复展开。
- 影响评估：同 api 条目；`mp` 命令后续还要渲染 module / proto / service / convert 四类文件并执行 `protoc` 与 `protoc-go-inject-tag`（`tpl_proto.go:53-58`），DB 串行 RTT 在总时长中占比相对更小。
- 优化建议：同 api 条目，`errgroup` 并发只读查询。

#### [P2] `toolkit/module/tpl_module.go:19` `GenerateModule`（共享模式：module 包 4 个 + vue 包 3 个生成函数）：7 个生成函数各自重复注册同一套模板函数映射

- 问题：7 个 `Generate*` 函数各自 `template.FuncMap{...}` 重复列出核心 10 个公共函数（`Convert`、`SymbolChar`、`Char`、`Helper`、`CamelStr`、`Add`、`ModuleProtoConvertDao`、`ModuleDaoConvertProto`、`ModuleProtoConvertMap`、`ApiToProto`），位置如下：
  - module 包：`tpl_module.go:20-29`（10 个）、`tpl_proto.go:23-33`（11 个，额外 `ProtoRequest`）、`tpl_service.go:20-29`（10 个）、`tpl_convert.go:20-29`（10 个）；
  - vue 包：`interface.go:21-32`（12 个，额外 `TypeScriptCondition`、`TypeScript`）、`method.go:42-52`（11 个，额外本地 `ApiList`）、`page.go:25-41`（17 个，额外 `TypeScriptCondition`、`TypeScript`、`Json`、`InMethod`、`Rule`、`Props`、`PropsSearch`）。
  每次 `template.New(...).Funcs(...)` 都重建该 map；公共部分在 7 处逐字重复维护。
- 影响评估：单次命令每个模板解析一次，map 重建开销为十几次键值注册，量级很小；主要价值是消除重复维护（新增公共模板函数需 7 处同步），性能上仅为冗余小分配。
- 优化建议：提取包级公共 `var commonFuncs = template.FuncMap{...}`（10 个核心函数），各 `Generate*` 以 `template.New(...).Funcs(commonFuncs)` 为基础，需要额外函数时再局部追加；与 `toolkit/dao/dao.go:100-104` 的单表重复解析条目配套，模板解析实例同样可在命令启动时预构建缓存。

> 本包其余共享模式条目（`Run` 内表名 `CamelStr` 重复计算、`InArray` 反射线性判断）见 `toolkit/api` 章节合并条目，调用点行号已在该条目中逐项列出；FuncMap 重复注册条目已合并至本包 `tpl_module.go:19` 一条，vue 三个生成函数的位置一并列在该条目中。

---

## 八、toolkit/project

通读 `toolkit/project/new.go:28-59`、`new.go:63-94`、`add.go:24-57`、`backend.go:41-82`：

- 目录存在性判断、覆盖确认、`Repo.CopyTo / CopyFrontTo / CopyToV2` 委托 `toolkit/base/repo.go`，文件拷贝性能问题归入 `toolkit/base/path.go` `copyFile` 条目；
- `backend.go:66-69` 使用缓冲大小为 1 的 `done chan error` 启动 goroutine 执行 `NewFront`，配合 `select` 等待 `ctx.Done()` / `done`，goroutine 写完即退出，无泄漏、无无界队列；
- 模板仓库克隆 / 拉取耗时由 `git clone` / `git pull` 子进程决定（`toolkit/base/repo.go:69-105`，`Pull:69-84`、`Clone:89-105`），属外部进程，不在本包 Go 代码优化范围。

**本包范围内无独立性能条目。**

---

## 九、toolkit/upgrade

通读 `toolkit/upgrade/upgrade.go:29-41`：仅把 6 个工具路径传给 `base.GoInstall`，串行安装问题已归入 `toolkit/base/install.go:15` 条目；本包无独立性能条目。

---

## 十、toolkit/vue

### P2

#### [P2] `toolkit/vue/method.go:18` `ApiList`：每次模板调用都新建 7 项 map 并为每个 key 重新走 `PathChar` 字符串拼装

- 问题：模板在每个方法分支调用一次 `ApiList`（`method.go:99,104,109,114,119,124,130,134`，共 7 类操作 × 最多 2 个 List 分支）；函数内部每次 `list := make(map[string]string)`（`method.go:20`）并逐键调用 `base.PathChar`（`method.go:21-27`）拼装路径，命中后再从 map 取回。`apiUrl` 在一次命令内固定，7 条路径结果恒定不变。
- 影响评估：单次渲染调用 7~8 次，每次 1 个 map + 7 个切片 / 字符串分配；量级很小，属于确定性冗余。
- 优化建议：`apiUrl` 由 `GenerateMethod` 通过 `base.SetUrl` 注入（`method.go:40`）后，在 `GenerateMethod` 入口一次性算出 7 条路径存为局部变量 / 闭包；`ApiList` 改为直接 `switch key` 返回预置片段，去掉每次调用建 map；`base.PathChar` 的切片 `append` 拼装（`toolkit/base/util.go:46-52`）也可直接字符串拼接。

#### [P2] `toolkit/vue/page.go:24` `GeneratePage`（模板内 `InMethod` 反复调用，函数实现见 `toolkit/base/util.go:60`）：同一批方法类型在单页模板渲染中重复判定 30 余次，每次重新线性扫描 + 反射包装

- 问题：`PageTemplate` 内对每种操作类型重复调用 `InMethod .Method "..."`，涉及位置：`page.go:88,100,108,111,114,117,128,133,142,151,196,199,202,205,208,211,219,222,225,228,231,234,237,294,297,300,303,306,337,341,351,359,367,375,385,397`。每次判定都重新遍历 `moduleParam.Method` 并经 `util.InArray` 反射比较；方法集合在渲染期间不变，30 余次判定结果完全可复用。
- 影响评估：方法数 ≤10，单次判定微秒级；累计为 30+ 次重复扫描与反射分配，属于冗余而非瓶颈。优化方向与 `toolkit/base/util.go:60` 条目一致：预构建类型集合后 O(1) 查找。
- 优化建议：在 `GeneratePage` 入口把 `moduleParam.Method` 的类型集合预建成 `map[string]struct{}`，注册为新的模板函数（或在 `ModuleParam` 上预计算），模板内 `{{if HasMethod "Create"}}` 直接查 map。

#### [P2] `toolkit/vue/vue.go:71` `Run`（同 `vue.go:81`、`vue.go:86`、`vue.go:91`）：4 次元数据查询串行发起

- 问题：与 `toolkit/api/api.go:92,101,106,111` 同一串行模式（`TableColumn` → `TableItem` → `TableIndex` → `TablePrimary`），改造方案见该条目，此处不重复展开。
- 影响评估：vue 命令后续渲染 interface / method / page 三类模板，DB 串行 RTT 占比取决于网络环境，需实测。
- 优化建议：`errgroup` 并发只读查询，错误处理语义保持现状。

> 本包 `vue.go` 中表名 `CamelStr` 重复计算（`132,143,154,165,175-178,187,197,204,214,218,264,291`）与 `InArray` 反射判断（`100,243,250,261,343,346,365,374,383,387,391`）已分别归入 `toolkit/api` 章节两条合并条目；`interface.go:20`、`method.go:41`、`page.go:24` 三处 FuncMap 重复注册已合并至 `toolkit/module/tpl_module.go:19` 条目。`interface.go` 模板仅在命令内解析一次，自身无独立热点。

---

## 十一、条目统计

| 子包 | P0 | P1 | P2 | 备注 |
|---|---|---|---|---|
| `toolkit/`（main.go） | 0 | 0 | 0 | 仅命令注册 |
| `toolkit/api` | 0 | 0 | 3 | 串行查询、表名重复转换（合并）、InArray 反射判断（合并） |
| `toolkit/base` | 0 | 2 | 7 | Props/PropsSearch 二次扫描；字段级转换缓存；NewDataType 重建；双循环；Sprintf/builder；copyFile；InMethod；go mod graph；GoInstall 串行 |
| `toolkit/build` | 0 | 0 | 0 | 已用 errgroup 并行 |
| `toolkit/change` | 0 | 1 | 3 | 正则每标签重编译；分类桶重建；categories 切片重建；git 子进程串行 |
| `toolkit/dao` | 0 | 1 | 1 | 模板逐表重复解析；跨表串行查询 |
| `toolkit/module` | 0 | 0 | 2 | 串行查询；FuncMap 重复注册 |
| `toolkit/project` | 0 | 0 | 0 | 委托 base 拷贝，无独立热点 |
| `toolkit/upgrade` | 0 | 0 | 0 | 委托 GoInstall |
| `toolkit/vue` | 0 | 0 | 3 | ApiList 建 map；InMethod 重复判定；串行查询 |
| **合计** | **0** | **4** | **19** | 跨包合并条目已在归属包内计数一次 |

---

## 十二、范围内高价值候选（供全局 Top 10 评选）

按「改造确定性 × 可消除重复劳动量 × 覆盖代码面」排序：

1. `toolkit/dao/dao.go:100-105`（`GenerateDao`，循环调用点 `dao.go:74`）——批量生成时每张表重复获取模板、注册 FuncMap、`template.Parse`；命令启动时预解析一次复用。**预期收益**：表数 N 时消除 N−1 次完整模板编译，批量 dao 生成命令的 Go 侧 CPU 与分配随表数线性下降。
2. `toolkit/change/change.go:117`（`classifyCommits`）——正则提为包级 `MustCompile`，避免每个 tag 重复编译。**预期收益**：标签数 N 时消除 N−1 次正则编译；改造成本极低（一处上移），无行为变化。
3. `toolkit/base/util.go:554-586`（`Props` / `PropsSearch`）——条件字段集合预建 `map[string]struct{}`，消除 O(n·m) 反射 `InArray` 扫描。**预期收益**：宽表渲染一次降为 O(n+m)，并去掉每次成员判断的反射包装分配。
4. `toolkit/base/util.go:88,373,777`（`ModuleDaoConvertProto` / `ApiToProto` / `ModuleProtoConvertDao` 等全部渲染 helper）——字段名驼峰 / 小驼峰结果在元数据加载阶段预计算挂到 `Column`。**预期收益**：一次命令内数百到上千次重复字符串转换与 `cases.Title` 分词器重建全部消除，覆盖 api / module / vue 全部渲染链路。
5. `toolkit/base/query.go:234,283`（`TableColumn` / `TablePrimary`）——`NewDataType()` 提为包级只读映射表。**预期收益**：批量 dao 场景消除每表一次的 33 项映射表重建；改动小、语义不变。
6. `toolkit/base/path.go:40-59`（`copyFile`，被 `copyDir` 递归调用）——替换字面量预转 `[]byte` 一次，非文本流式拷贝并先 `Contains` 命中再替换。**预期收益**：整套模板仓库复制时避免逐文件重复字节分配与多轮全量扫描；大二进制资源降低内存峰值。
7. `toolkit/api/api.go:157`（`Run`，共享 `module.go` / `vue.go` 装配流程）——表名 `CamelStr` 提到命令级局部变量复用；`InArray` 集合判断改原生 map。**预期收益**：单次命令消除十余次重复转换与一批反射成员判断，三处命令统一受益。
8. `toolkit/dao/dao.go:70`（`Run` 批量循环）——跨表 `TableColumn` 查询限定并发度并行。**预期收益**：远端库多表场景总耗时由 N × RTT 降至接近 max + 排队（需结合连接池上限实测）。
9. `toolkit/change/change.go:191-207`（`Run` 标签循环）——标签级 git 子进程受限并发，输出按原顺序回填。**预期收益**：多标签仓库 changelog 生成总时长由串行 Σ 压缩为近似并行最大值。
10. `toolkit/vue/method.go:18-32`（`ApiList`）——7 条路径预计算 + `switch` 返回，去掉每次调用建 map。**预期收益**：去掉页面渲染时每方法一次的 map 与切片分配；顺带简化 `base.PathChar` 拼装。

---

## 十三、拿不准、需实测验证的条目清单

1. **4 次串行元数据查询的 RTT 占比**（`api.go:92,101,106,111`、`module.go:88,97,102,107`、`vue.go:71,81,86,91`）：本地 MySQL 下 information_schema 查询为毫秒级，串行化收益不明显；需在目标数据库网络环境下分别计时「串行」与「errgroup 并发」两种实现，确认并发是否受 GORM 连接池 / 驱动限制。
2. **dao 模板重复解析的实际耗时**（`dao.go:100-105`）：需对批量 N 张表做命令级计时（`time toolkit dao [tables...]`），并用 benchmark 对比「逐表 Parse」与「预解析复用」的 `allocs/op` 与 `B/op`；注意末尾 `go fmt` / `goimports` 子进程耗时可能掩盖 Go 侧差异，需单独拆分计时。
3. **字段级 `CamelStr` / `Helper` 重复转换的分配占比**（`base/util.go` 各 helper）：需 `benchmem` 跑宽表（30~50 字段）渲染，对比预缓存方案；`util.UCWords` 内 `cases.Title` 重建的成本是否占大头，需结合跨模块提示第 1 条一起 profile。
4. **`copyFile` 整文件读入内存的实际影响**（`base/path.go:46,58`）：取决于 layout 模板仓库是否包含较大二进制资源；需统计模板仓库文件大小分布后判断流式拷贝是否必要。
5. **`ModuleVersion` 派生 `go mod graph` 的耗时归属**（`base/mod.go:24-44`）：绝大部分耗时应在 `go` 子进程冷启动与依赖图构建，Go 侧缓冲扫描占比很低；是否值得换 `go list` 定向查询需实测子进程总耗时。
6. **change 标签 git 子进程并行化收益与安全性**（`change.go:191-207`）：需确认同一仓库内并发 `git log` 是否触发锁竞争 / 对象库读放大；并行后输出顺序由代码回填保证，不影响文档结构，但并发度需要实测调优。
7. **`InMethod` / `InArray` 反射开销在 CLI 场景的可感知度**（`base/util.go:60`、`vue/page.go` 模板）：方法与条件切片均为个位数长度，反射包装的单次开销可能被交互提示、DB 查询、外部子进程完全掩盖；改造优先级按代码整洁度与宽表上限考虑，而非当前实测延迟。

---

## 十四、跨模块提示（仅记录线索，不在本章节审查）

1. **`util/string.go:155-158` `UCWords`**：每次调用 `cases.Title(language.English)` 都新建一个 title caser；caser 可安全复用，应提为包级变量。该函数被 `toolkit/base/util.go:12` `CamelStr` 高频调用，是条目「字段级转换缓存」之外的第二层根因；改造后所有调用方（不止 toolkit）受益。
2. **`util/array.go:223-243` `InArray`**：基于 `reflect.ValueOf` + 逐元素 `reflect.DeepEqual`，无类型特化；toolkit 内大量成员判断（`api.go`、`module.go`、`vue.go`、`dao.go`、`util.go:561,580,63`）都经此函数。建议 `util` 侧补充强类型 `StringInSlice([]string, string) bool` / `map` 集合构造辅助，toolkit 侧逐步替换；`util/var.go:12-35` `Empty` 同样为反射实现，被 `build.go:93,96,99,102` 等调用，一并记录备查。
3. **`watch` 包**：`toolkit/build/build.go:112` 的 `watch.Watch(ctx, s)` 为文件监听与自动重编译主体，目录递归、监听粒度、事件合并、重复编译抑制等性能行为不在本章节范围内。
4. **`stores/sql` 连接池**：`toolkit/base/init.go:70-130` `InitQuery` 透传 `MaxIdleConns` / `MaxOpenConns` / `MaxLifetime` / `MaxIdleTime` 配置，并发元数据查询时的池上限行为由 `stores/sql` 与底层驱动决定，本章节未展开。
5. **模板外部进程**：`go fmt`（`api.go:435`、`dao.go:77`、`tpl_module.go:46`、`tpl_service.go:47`、`tpl_convert.go:46`）、`goimports`（`api.go:440`、`dao.go:82`、`tpl_module.go:51`、`tpl_service.go:52`、`tpl_convert.go:51`）、`protoc`（`tpl_proto.go:53`）、`protoc-go-inject-tag`（`tpl_proto.go:57`）、`git clone/pull`（`base/repo.go:69-105`）均为外部子进程；`vue` 三个生成函数只写 `.ts` 文件，不派生格式化子进程。外部进程耗时与优化空间不属于 toolkit Go 代码审查范围。
