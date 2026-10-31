# util / registry / snowflake / config / crontab / watch / worker / img 包组性能审查报告

> 审查范围（仅本章节）：`util` 全部生产 `.go` 文件；`registry` 根目录（`config.go`、`endpoint.go`、`nop.go`、`registry.go`、`service.go`、`update.go`、`doc.go`）；`registry/etcdv3`（`config.go`、`registry.go`、`doc.go`）；`snowflake`（`snowflake.go`、`doc.go`）；`config`（`config.go`、`doc.go`）；`crontab`（`task.go`、`doc.go`）；`watch`（`builder.go`、`logger.go`、`options.go`、`watch.go`、`doc.go`）；`worker`（`worker.go`、`job/job.go`、两个 `doc.go`）；`img`（`img.go`、`fontx/font.go`、两个 `doc.go`）。测试文件已通读，性能条目只针对生产代码。
> 审查方式：逐文件人工通读源码，先借 Godoc 注释理解 API 语义，再定位热点路径与可疑实现；每条结论回到原文核对 `文件:行号` 与函数名；跨模块调用点仅用于校准触发频率与影响量级，不展开其他模块审查。`img` 包因本机缺 libvips/pkg-config 无法编译，按任务要求以静态阅读为准。
> 优先级定义：**P0** = 热路径上的明显性能缺陷，应优先修改；**P1** = 中等影响、性价比高的改进；**P2** = 低影响或改进空间有限的点。

---

## 一、条目清单（按子包分组；组内 P0 → P1 → P2 排序）

### util/string.go

#### [P1] `util/string.go:637` `JSONString`：Marshal 出的 `[]byte` 经 `cast.ToString(any)` 多一层接口装箱与类型断言，可直接 `string(jsByte)`
- 问题：`JSONString` 先 `json.Marshal(obj)` 得到 `jsByte []byte`（638 行），随后 `return cast.ToString(jsByte)`（639 行）。核对 `spf13/cast@v1.10.0/basic.go:102` 的 `ToStringE`：`[]byte` 分支就是 `return string(s), nil`，与直接转换结果相同；当前实现多出来的开销是：调用点把 `jsByte` 包进 `any` 接口（切片头需堆上装箱小对象），进入 `ToStringE` 再做一次运行时类型 switch 断言与错误返回值解包。需要澄清的是：**大 payload 的 `[]byte → string` 拷贝在两种写法下都只发生一次**，本条目要消除的是装箱分配 + cast 层分发这层固定间接开销，而不是第二次整包字节拷贝。
- 影响评估：场景已核实——`client/grpc/interceptor.go:144`、`:145`、`:156`、`:157`、`:169`、`:170` 在 Unary 访问日志拦截器中对 `req`、`reply` 各调用一次 `util.JSONString(...)`，`:44` 的调试打印还调用一次；即**每条经过该拦截器的 gRPC 调用**都执行 2~3 次本函数。量级判断：单次额外成本很小（一次切片头装箱分配 + 一次类型断言），但在每请求 2~3 次的日志路径上属于可一行移除的固定分配与函数间接层；`json.Marshal` 本身的序列化开销才是大头，归调用侧日志门控问题（见第五节第 1 条）。
- 优化建议：函数体改为 `jsByte, _ := json.Marshal(obj); return string(jsByte)`，去掉 `github.com/spf13/cast` 依赖；直接转换一次得到字符串，省掉装箱与类型 switch。若拦截器在日志级别未开启时也无条件构造 `req/reply` 的 JSON 字符串，调用侧应再加日志级别门控（属跨模块配合项，见第五节）。

#### [P1] `util/string.go:910` `Tag`：静态分隔符表下每次调用都重新 `strings.NewReplacer`
- 问题：`Tag` 在 911 行 `replacer := strings.NewReplacer(tagPatterns...)`，而 `tagPatterns`（28–35 行）是包级静态切片，运行期从不变化。`NewReplacer` 每次都会把成对键值构建成内部前缀树/排好序的替换表，属于重复构建同一结构。
- 影响评估：`Tag` 用于把标签字符串中的逗号、竖线、顿号、分号等统一替换为 `/`，是面向内容处理（标签归一化、关键字拆分）的导出工具；一旦在批量导入、搜索分词、标签渲染等循环中被逐个元素调用，每次都支付一次替换器构建成本（切片遍历 + 内部表分配），替换本体反而只占小头。全仓检索确认当前**无任何调用点**（含测试），风险定位为导出 API 的潜在批量热点，优先级按构造开销固定、修复成本极低保留 P1。
- 优化建议：包级初始化一次 `var tagReplacer = strings.NewReplacer(tagPatterns...)`，`Tag` 内直接 `return tagReplacer.Replace(cast.ToString(i))`，消除每次调用的 Replacer 构造与分配。

#### [P2] `util/string.go:155` `UCWords`：每次调用新建 `cases.Title` 大写器
- 问题：`UCWords` 在 156 行 `caser := cases.Title(language.English)`，157 行 `caser.String(str)`。`cases.Title` 携带语言规则与内部状态，是可复用对象；此处每次调用都重新构造。
- 影响评估：仓内调用点集中在 `toolkit` 代码生成期（`toolkit/base/util.go:14`、`toolkit/api/api.go:298`、`toolkit/vue/vue.go:259` 等），属离线生成路径，频率低；若未来在 Web 表单、用户名展示等在线循环中使用，构造开销会被放大。
- 优化建议：包级 `var titleCaser = cases.Title(language.English)`，函数内复用；注意 `cases.Title` 返回的 `Caser` 按官方实现可并发复用。

#### [P2] `util/string.go:62` `StrIPos` / `util/string.go:103` `StrRIPos`：忽略大小写查找时先把整段 haystack 转小写
- 问题：`StrIPos` 在 72 行 `strings.Index(strings.ToLower(haystack), strings.ToLower(needle))`；`StrRIPos` 在 114 行同样对 `haystack` 整体 `ToLower`。无论 `needle` 是否存在，haystack 都会先被完整复制一份小写副本；`needle` 只做了小写，长字符串上整串拷贝是主要开销。
- 影响评估：字符串定位类工具，非框架内置热路径；仅当在长文本（日志、正文）循环中反复忽略大小写查找时，每次调用多一次与 haystack 等长的堆分配。
- 优化建议：优先使用 `strings.EqualFold` 思路或 `strings.ToLower(needle)` 后配合已大小写无关的查找；标准库没有直接的 `IndexFold`，若必须保留语义，可在短 needle 上做窗口 `EqualFold` 比较，避免整串 `ToLower` 拷贝；至少把 `strings.ToLower(needle)` 提出索引判断之外（当前实现已如此）。

#### [P2] `util/string.go:706` `Md5File`：小文件分支先 `os.Open` + `Stat`，又 `os.ReadFile(path)` 二次打开磁盘文件
- 问题：`Md5File` 707 行 `os.Open(path)`，713 行 `f.Stat()` 判断大小；小文件分支 722 行却再次 `os.ReadFile(path)`，绕开已打开的 `f` 重新发起一次打开与读全量；727–737 行大文件分支才复用 `f` 分块读取。
- 影响评估：文件摘要计算属低频管理操作；问题是小文件场景多一次 `open/read/close` 系统调用往返与一次路径解析，小文件批量摘要时（如扫描目录算哈希）累积可见。
- 优化建议：小文件分支直接 `io.Copy(hash, f)`（或 `io.ReadAll(f)`），复用 707 行已打开的文件句柄；大小判断沿用 `f.Stat()` 结果即可。

#### [P2] `util/string.go:750` `Sha1File`：整文件一次性读入内存，未像 `Md5File` 那样分块
- 问题：`Sha1File` 751 行 `os.ReadFile(path)` 把整个文件读进 `data []byte`，再 `hash.Write(data)`；大文件会在内存里同时驻留「文件全量字节 + 哈希内部状态」，缺少 `Md5File`（718–737 行）已有的分块读取模式。
- 影响评估：摘要工具类 API；只要被用于日志包、归档包、镜像等大文件场景，内存峰值随文件大小线性上升；仓内当前未见大文件调用点，定 P2。
- 优化建议：与 `Md5File` 对齐，`os.Open` 后用固定大小缓冲区（如 1MB）循环 `io.CopyBuffer` 写入 `hash`，保持恒定内存占用。

#### [P2] `util/string.go:818` `SimilarText`：三层嵌套循环 + 递归，递归内反复 `[]byte`/字符串切片拷贝
- 问题：`SimilarText` 的闭包 `similarText` 是 `len1 × len2 × 公共前缀长度` 的三层嵌套循环（824–833 行），命中最长片段后在 838、843 行递归两个子区间；841–843 行每次递归都 `[]byte(str1)`、`[]byte(str2)` 转换并对子串再做字符串切片拷贝。复杂度在字符串较长时明显恶化（最坏 O(n²·m) 量级叠加递归分治开销）。
- 影响评估：编辑距离/相似度类工具，框架内置热点的可能性低；只有被接入模糊搜索、查重、推荐匹配等在线循环时才会成为热点。仓内未见生产调用，**实际触发频率需实测确认**（见第四节第 3 条）。
- 优化建议：若保留语义，公共片段提取可用后缀数组/动态规划降复杂度；递归前先把字符串转一次 `[]rune`/`[]byte` 向下传递，避免每层递归重复转换；长字符串场景建议直接替换为标准库成熟的编辑距离实现或 `Levenshtein`（767 行，已是二维滚动数组实现，复杂度 O(n·m)）。

#### [P2] `util/string.go:482` `StrTr`（map 分支）：`make([]string, length*2)` 后再 `append`，切片实际扩容到 4 倍长度；每次调用现场 `NewReplacer`
- 问题：map 分支 490 行 `oldnew := make([]string, length*2)` 先把切片长度置为 `2*length`，491–496 行再 `append(oldnew, o, n)` 追加 `2*length` 个实参，最终长度 `4*length`（前半段全是空字符串槽位），容量不足时还要再扩容一次；497 行 `strings.NewReplacer(oldnew...)` 每次调用现场构建替换器，且把一倍数量的空 `""` 键值对也传入构建流程参与内部表构建。
- 影响评估：仓内检索确认 `StrTr` 当前**无生产调用点**（仅定义本身），定 P2；一旦在字符串替换循环中使用，每次都多一份冗余切片分配与 Replacer 表构建。
- 优化建议：`oldnew := make([]string, 0, length*2)` 起步直接 `append`；替换器若固定模式可缓存，通用多对替换仍在调用时构建；同时注意前半段空串槽位会被 `NewReplacer` 当作空键处理，属附带的正确性隐患（非本条目性能结论）。

### util/array.go

#### [P1] `util/array.go:223` `InArray`：通用反射遍历 + `DeepEqual`，每次比较都装箱与反射比较
- 问题：`InArray` 通过 `reflect.ValueOf(haystack)`（224 行）判断 Kind，切片/数组逐元素 `val.Index(i).Interface()`（228 行）装箱后调用 `reflect.DeepEqual(needle, ...)`；map 分支（232–237 行）同样先 `MapKeys()` 取键切片再逐个 `Interface()`。`DeepEqual` 内部还要走 `reflect.Type` 比较与递归反射；即使两边都是 `string`、`int` 这类标量，也不做类型快路径。
- 影响评估：场景已核实——`client/etcdv3/config.go:98` `SetEndpoints` 逐个端点去重；`img/img.go:51` 格式白名单判断；`config/config.go:64` 配置扩展名白名单；`util/mongodb.go:48` `ConvertBson` 键存在性判断；仓内现有调用共 34 处生产点（`toolkit` 代码生成包内约 30 处，其余 4 处为启动期配置/白名单判断），单次数组短、当前量级有限；但作为导出通用工具，框架使用者很容易把 `[]string` 包含判断放进请求循环，届时每个元素都要支付反射 `ValueOf`、`Interface()` 装箱、`DeepEqual` 反射三层开销。
- 优化建议：
  1. 增加标量快路径：先 `needle` 与 `haystack` 的具体类型（`[]string`、`[]int`、`[]int64`、`map[string]struct{}` 等）走直接 `==` / map 查找，O(1) map 命中，不再反射。
  2. 反射仅作为非标量/混合类型的兜底分支保留。
  3. map 分支用 `val.MapIndex(k)` 直接取值，避免先构建 `MapKeys()` 切片。

#### [P1] `util/array.go:352` `ArrayRemoveRepeatedElement` 与 `util/array.go:473` `ArrayStringUniq`：O(n²) 双重循环去重
- 问题：两个函数实现完全同构——外层 `for i` 遍历，内层 `for j := i+1; j < len(arr); j++`（356 行 / 477 行）与 `arr[i]` 逐一比较判断重复，命中即跳过。每个元素都要向后扫描剩余全部元素，复杂度 O(n²)；切片 `newArr` 也以 `make([]string, 0)` 起步（353 / 474 行），伴随多次扩容拷贝。
- 影响评估：字符串切片去重是高频数据处理动作；一旦对批量标签、URL、表头、错误码等几百上千元素的切片调用，比较次数按平方增长，元素相等时还会在循环内做逐字节字符串比较。仓内当前仅测试覆盖（`util/array_test.go:526`、`:770`），风险为导出 API 在在线批量处理中的潜在热点。
- 优化建议：用 `map[string]struct{}` 记录已出现元素，单趟 O(n) 遍历完成去重并保留首次出现顺序；切片按 `make([]string, 0, len(arr))` 预分配容量，避免反复扩容。两个函数可合并为一个内部实现。

#### [P1] `util/array.go:308` `SplitString`：循环内 `res += string(c)` 逐字符 O(n²) 拼接 + 每次调用 `regexp.MustCompile` 现场编译 + 双层 map 冗余去重
- 问题：三处问题叠加：
  1. 310–323 行在 rune 循环里反复 `res += string(c)` / `res += string(split)`，Go 字符串不可变，每次拼接都从旧串复制一份新串，长度为 n 的文本累计 O(n²) 次字节拷贝；
  2. 325 行 `reg := regexp.MustCompile(`\\/{2,}`)` 在函数体内现场编译正则，每次调用都走一遍正则语法分析与 NFA 构建；
  3. 328–337 行先构造 `newAry map[string]string` 去重，又构造 `newArys map[any]any` 把同一批键值再装一遍，随后 `ArrayValues` + `InterfaceToString` 转换，两层 map 与两次类型转换纯属冗余。
- 影响评估：文本清洗（标点/符号归一、分隔符统一）工具；若用于批量内容处理（导入文本、字段清洗循环），长文本拼接平方增长 + 每次正则编译在该路径上最显眼。仓内当前仅测试覆盖，优先级按导出 API 批量场景定为 P1。
- 优化建议：
  1. 循环改用 `strings.Builder`（或 `bytes.Buffer`），先估算容量，把 O(n²) 拼接降为 O(n)；
  2. 把 `/{2,}` 正则提为包级 `var` 预编译一次；
  3. 单层 `map[string]struct{}` 完成去重，去掉 `newArys`、`ArrayValues`、`InterfaceToString` 这一串中间转换；
  4. `StrTrim`（331、332 行）对同一元素调用两次，可缓存局部变量。

#### [P2] `util/array.go:257` `InMultiArray`：双层循环 O(n·m) 叠加反射与 `DeepEqual`
- 问题：切片/数组分支 262–268 行对 haystack 每个元素、needle 每个元素做嵌套 `DeepEqual`；map 分支 270–276 行同样先取 `MapKeys()` 再双层嵌套比较。needle 与 haystack 都需要包含判断时，正确做法是把一侧建成集合 O(n+m)，当前实现是 O(n·m) 反射比较。
- 影响评估：通用多值包含判断，现有仓内未见生产调用；与 `InArray` 同属反射型通用工具，被放进批量循环时成本高于 `InArray`（多一层嵌套）。
- 优化建议：先把 needle 建成 `map[any]struct{}`（标量类型另走具体类型 map），再单趟扫描 haystack 命中即返回；保留反射分支仅用于非标量元素。

#### [P2] `util/array.go:298` `InterfaceToString`、`util/array.go:343` `ArrayPluck` 等：`make([]T, 0)` 起步，追加过程多次扩容
- 问题：`InterfaceToString`（298–303 行）、`ArrayPluck`（343–348 行）、`MapStringToMapInterface`（403–410 行）、`MapInterfaceToMapString`（423–430 行）等转换函数均以 `make([]T, 0)` / `make(map, 0)` 起步，随后 `append`，长度已知时仍会触发切片多次扩容拷贝；`InterfaceToString` 每个元素还用 `fmt.Sprintf("%v", v)`（300 行）走反射格式化。
- 影响评估：切片/映射批量转换工具，单次调用规模通常有限；主要是扩容次数与 `fmt.Sprintf` 反射格式化在批量转换时的固定浪费。
- 优化建议：入口处 `make([]string, 0, len(data))`、`make(map[string]any, len(d))` 预分配容量；标量值优先用 `cast.ToString`（避免 `fmt.Sprintf` 的反射路径）。

#### [P2] `util/array.go:122` `ArrayRand`、`util/array.go:247` `ArrayRandMap`、`util/string.go:544` `StrShuffle`：每次调用都 `rand.New(rand.NewSource(time.Now().UnixNano()))` 新建随机源与 `rand.Rand`
- 问题：三处实现同构——`ArrayRand` 122 行、`ArrayRandMap` 247 行、`StrShuffle` 544 行均在函数体内 `r := rand.New(rand.NewSource(time.Now().UnixNano()))`，随即 `r.Perm(...)` 取排列；每次调用都读取一次系统时钟、构造 `rngSource` 与 `*rand.Rand` 两个对象，不复用任何随机源。与 `util/math.go:22` `Rand`（P1，已被选库热路径调用）、`util/init.go:13` `Do`（P2，退避抖动）属同一模式。
- 影响评估：仓内检索确认这三个函数当前**仅有测试调用、无生产调用点**，定 P2；潜在风险与 `Rand` 一致：同纳秒批量并发时种子相同可能产生相同排列序列（正确性隐患，见第四节第 5 条），高频洗牌场景下每次新建源的对象分配叠加。
- 优化建议：与 `Rand`、`Do` 统一为包级复用随机源（互斥锁保护的 `*rand.Rand`，或标准库全局 `rand.Perm` / `math/rand/v2` 的无锁源），三个洗牌函数改为调用共享源，删除函数体内的 `New`/`NewSource`。

### util/math.go

#### [P1] `util/math.go:22` `Rand`：每次取随机数都新建随机源与 `rand.Rand` 实例
- 问题：`Rand` 在 32 行 `r := rand.New(rand.NewSource(time.Now().UnixNano()))`，33 行 `r.Intn(...)`。每次调用都：读一次 `time.Now().UnixNano()`、构造 `rngSource` 结构、构造 `*rand.Rand`，再取一个整数；既没有复用随机源，也没有使用包级全局随机源。
- 影响评估：场景已核实——`stores/proxy/sql.go:36` `Write()` 与 `:46` `Read()` 在**每次选库**时调用 `util.Rand(0, len-1)`，即每条走读写分离代理的 SQL 语句都支付一次新随机源构造；写/读库列表通常只有几个元素，`Intn` 本身极轻，新源构造相对占比不低。量级判断：每查询 1 次 `time.Now` 系统时钟读取 + 两个小对象分配 + 源初始化；全局 `math/rand` 的 `rand.Intn` 内部是加锁复用的长期源，成本远低于新建源。
- 优化建议：改用包级复用的随机源：要么直接用标准库全局 `rand.Intn(max+1-min)+min`（全局源有内部互斥锁，但省去每次新建源）；要么在 `util` 包内维护一把锁保护的 `*rand.Rand`（`sync.Mutex` + 复用 `rand.New(rand.NewSource(...))`），或改用并发安全的 `math/rand/v2`。另注意 `time.Now().UnixNano()` 做种子在同一纳秒批量并发调用时可能产生相同序列，是正确性隐患（见第四节第 5 条）。

### util/time.go

#### [P1] `util/time.go:153` `Date` 与 `util/time.go:141` `StrToTime`：每次调用都用静态 `datePatterns` 重建 `strings.Replacer`
- 问题：`Date` 156 行 `replacer := strings.NewReplacer(datePatterns...)`，157 行 `replacer.Replace(format)`；`StrToTime` 142–143 行同构。`datePatterns`（15–46 行）是包级静态常量切片，两次调用间不会变，但 Replacer 每次都重新构建成对键值与内部表。`Date` 是时间格式化的核心导出函数，还被同包 `UnixTimeFormatDate`（207 行）、`FormatDateTime`（212 行）、`FormatDate`（220 行）、`MonthDate`、`NextMonthDate`、`PrevMonthDate`、`Month`、`Day`、`Hour`、`ToWeekDay` 层层复用（详见下文各条）。
- 影响评估：场景已核实——`core/env/const.go:120` 启动期调用一次；其他 `FormatDateTime/FormatDate` 等封装是面向业务代码的时间展示工具，若在接口响应、日志、报表循环中反复格式化时间，每次都叠加一次 Replacer 构建（16 对替换项的遍历与内部分配）。复杂度不变，属于可一次性消除的重复构造开销。
- 优化建议：包级 `var dateReplacer = strings.NewReplacer(datePatterns...)`，`Date` 与 `StrToTime` 直接复用；Replacer 按标准库实现可并发使用。

#### [P2] `util/time.go:393` `Hour`：循环内反复 `time.ParseDuration(cast.ToString(i)+"h")` 解析同一个递增步长
- 问题：`Hour` 在 400–406 行每轮 `time.ParseDuration(cast.ToString(i)+"h")`（401 行），得到 `t` 后 `start.Add(t)`；步长恒定为 1 小时，循环内却每次做整数转字符串、时长字符串解析两次转换。相邻的 `Month`（359–363 行）、`Day`（379–387 行）每个时间点还各调两次 `Date(...)`（如 361–362、375–376、381–382 行），每次 `Date` 又重建 Replacer（见上一条）。
- 影响评估：按区间批量生成时间轴（小时/天/月列表）的报表辅助函数，调用频率取决于业务报表；区间越大，循环内重复解析与重复格式化次数越多。
- 优化建议：`Hour` 直接 `start.Add(time.Duration(i) * time.Hour)`，去掉 `cast.ToString` 与 `ParseDuration`；`Month`/`Day` 循环内同一时刻只格式化一次字符串与数字，第二次结果由第一次派生，避免重复 `Date` 调用。

#### [P2] `util/time.go:411` `ToWeekDay`：为取年月日连续调用三次 `Date`，每次都重建 Replacer
- 问题：`ToWeekDay` 415 行 `cast.ToUint16(Date("Y", now))`、`Date("m", now)`、`Date("d", now)` 三次独立格式化，每次 `Date` 内部都执行一次 `NewReplacer` 构建（156 行）、一次布局替换与一次 `t.In(...)`；而 `now` 本身已是 `time.Time`，年月日本可直接取。
- 影响评估：中文星期名转换工具；单次调用量级小，但三次重复 Replacer 构建在被批量遍历日期（报表、日历渲染）时累积。
- 优化建议：直接 `now.Date()` 取得 `year, month, day`，删除三次 `Date(...)` 调用；星期计算沿用现有蔡勒公式即可。

### util/misc.go

#### [P2] `util/misc.go:56` `MemoryGetUsage` 与 `util/misc.go:65` `MemoryGetPeakUsage`：每次调用 `new(runtime.MemStats)` 堆分配并执行 `ReadMemStats`
- 问题：`MemoryGetUsage` 57 行 `stat := new(runtime.MemStats)`，58 行 `runtime.ReadMemStats(stat)`；`MemoryGetPeakUsage` 66–67 行同构。`ReadMemStats` 需要 STW 阶段汇总堆统计，本身成本不低；每次调用还在堆上新建一份 `MemStats`（约百字节结构），`realUsage` 参数按注释当前未被使用。
- 影响评估：运行时内存读取通常用于监控指标采集；若被指标采集器按固定频率调用，`ReadMemStats` 的 STW 成本是主要关注点，`new(MemStats)` 分配次之。仓内当前仅测试覆盖，实际采集频率需结合接入情况判断（见第四节第 1 条）。
- 优化建议：包级复用一个 `MemStats` 实例（`var memStats runtime.MemStats`，配合 `sync.Mutex` 保护 `ReadMemStats` 写入）；或仅在监控路径按需提供；评估 `ReadMemStats` 调用频率是否可降低。

#### [P2] `util/misc.go:76` `VersionCompare`：递归比较闭包内每次比较都新建 `forms` 预发布级别 map
- 问题：`VersionCompare` 把 `vcompare`、`canonicalize`、`special` 三个闭包定义在函数体内；`special` 在 209–220 行每次被调用都新建 `forms` map（10 个键值），随后两轮 `for name, order := range forms`（222、231 行）做前缀匹配。版本号比较在同一进程内可能被反复调用（依赖排序、版本校验），map 反复构造与哈希分配没有必要。
- 影响评估：版本比较非请求热路径；主要浪费是每次比较中 `special` 可能被调用多次（每个版本段一次），每次重建同一张静态映射表。
- 优化建议：`forms` 提为包级 `var` 只读 map；三个闭包可改为包级普通函数，避免每次调用 `VersionCompare` 都为闭包环境分配捕获结构。

### util/net.go

#### [P2] `util/net.go:141` `isPrivateIP`：每次判断都对 5 个常量网段重复 `net.ParseCIDR`
- 问题：`isPrivateIP` 在 144–148 行对 `"10.0.0.0/8"`、`"172.16.0.0/12"`、`"192.168.0.0/16"`、`"100.64.0.0/10"`、`"fd00::/8"` 五个常量串逐个 `net.ParseCIDR`，构造 `*net.IPNet` 切片后再 `Contains`；网段是固定常量，每次调用都重新解析 5 次 CIDR 字符串。
- 影响评估：`ExtractIP`（80 行）在启动期遍历本机网卡时对每个候选 IP 调用一次（120 行），频率低；问题是常量解析被重复执行，未来若把本函数用于请求级 IP 判定（内网识别、限流分区），重复解析会进入热路径。
- 优化建议：包级 `var privateBlocks = []*net.IPNet{...}` 在 `init()` 或 `sync.Once` 中解析一次；`isPrivateIP` 直接遍历预解析切片 `Contains`。另建议入参直接收 `net.IP`，避免 120 行 `ip.String()` 来回字符串化。

### util/mongodb.go

#### [P2] `util/mongodb.go:16` `Pipeline`：每个聚合阶段都走 JSON 反序列化 → bson.Marshal → bson.Unmarshal 三重序列化往返
- 问题：`Pipeline` 单文档分支 21 行 `json.Unmarshal([]byte(str), &doc)`，23 行 `bson.Marshal(doc)`，24 行 `bson.Unmarshal(b, &v)`；多阶段分支 28–33 行对每个阶段重复同一往返。JSON 对象先解成 `bson.M`，再 marshal 成 bson 字节，再 unmarshal 成 `bson.D`，两次 bson 编解码之间还隔了一次 JSON 解析。
- 影响评估：MongoDB 聚合管道构造工具；若在每次查询前由字符串动态构造管道，每个阶段支付三重编解码；静态管道应在启动期构造一次复用。是否处于查询热路径需结合调用方判断（见第四节第 6 条）。
- 优化建议：直接 `bson.Unmarshal`（或 `bsonx`）解析原始 JSON 字节到 `[]bson.D`，跳过 `bson.M` 中转与第二次 bson 往返；管道构造结果可缓存复用。

#### [P2] `util/mongodb.go:43` `ConvertBson`：先全量收集键成 `[]string`，再走反射型 `InArray` 判重，随后又全量循环定位同键
- 问题：`ConvertBson` 44–47 行先把 `items` 所有键收集到 `key` 切片；48 行 `InArray(item.Key, key)` 走 `util.InArray` 反射 + `DeepEqual`；命中后 52–64 行再从头遍历 `items` 找同键元素做合并。键存在性判断完全可以在同一次遍历中用 `map[string]int` 完成，当前是「收集切片 + 反射查找 + 二次遍历」三遍。
- 影响评估：bson.D 元素合并辅助函数，单次操作字段数有限；反射查找与两遍遍历在文档字段较多时才显现浪费。
- 优化建议：入口构建 `map[string]int` 记录键下标，单趟遍历完成判重与定位；去掉对 `InArray` 的依赖。

### util/pinyin.go

#### [P2] `util/pinyin.go:13` `ZhCharToFirstPinyin` 与 `util/pinyin.go:30` `ZhCharToPinyin`：逐字符调用 `pinyin.Pinyin`，每字符一次字典查找与切片分配
- 问题：两个函数在循环内对单个 rune 执行 `pinyin.Pinyin(string(r), a)`（19 行 / 36 行）：每字符先 `string(r)` 转换、传入整词接口只取 `[0][0]`；`pinyin.Pinyin` 本身支持整串输入并返回二维切片，逐字符调用等于把一次批量字典查询拆成 n 次，每次都产生内层切片。
- 影响评估：拼音首字母/全拼转换通常用于搜索索引、姓名排序等离线或批量场景；在线高频路径使用的可能性低，定 P2。
- 优化建议：整串传入 `pinyin.Pinyin(p, a)` 一次取得结果二维切片，再遍历拼接；`pinyin.NewArgs()`（14 / 31 行）可按样式（FirstLetter/Normal）包级复用。

### util/scprng.go

#### [P2] `util/scprng.go:40` `StrPad`：循环内 `str1 = str1 + "0"` 补零，O(n²) 字符串拼接
- 问题：`StrPad` 43–45 行 `for range n { str1 = str1 + "0" }`，每次加法都复制旧串；需要补齐 n 个零时累计 O(n²) 字节拷贝，最后 47 行再 `str1 + str2` 拼接一次。
- 影响评估：数字补齐工具，单次补齐位数有限（常见 2~6 位）；只有在循环内批量格式化、补齐位数较大时平方成本才显现。
- 优化建议：`str1 + strings.Repeat("0", n) + str2`，一次分配完成；或 `strings.Builder` 预分配容量。

### util/util.go

#### [P1] `util/util.go:93` `VerifyPhone`、`util/util.go:101` `VerifyEmail`、`util/util.go:109` `VerifyIPv4`、`util/util.go:115` `VerifyIPv6`：每次校验都现场编译正则表达式
- 问题：
  - `VerifyPhone` 96 行 `reg := regexp.MustCompile(regular)`，每次调用编译手机号模式；
  - `VerifyEmail` 104 行 `reg := regexp.MustCompile(pattern)`，每次编译邮箱模式；
  - `VerifyIPv4` 110 行 `regexp.MatchString(...)`，标准库该函数内部每次现场编译 IPv4 模式；
  - `VerifyIPv6` 116 行 `regexp.MatchString(...)`，每次现场编译一个长达数百字符的 IPv6 模式（116 行内联），IPv6 正则编译成本在四个函数中最高。
- 影响评估：表单/请求参数校验类导出 API，典型在线场景是每个注册、登录、资料提交请求对手机号/邮箱/IP 各校验一次；正则编译（语法分析、解析树、NFA）远贵于一次匹配本身，IPv6 模式尤其明显。仓内当前仅测试调用，问题定位为导出 API 接入请求校验后的必然热点。
- 优化建议：四个模式提为包级 `var phoneRegex = regexp.MustCompile(...)` 等预编译一次，函数体内只调用 `reg.MatchString(...)`；`util/regexp.go` 已为其他模式提供了同款预编译范例（68–127 行），保持一致风格即可。

#### [P2] `util/util.go:173` `NewReplacer`：循环内对每个键全量重建 endpoint 字符串，O(k·m) 重复拷贝
- 问题：`NewReplacer` 先把 `values` 收集进 `params` map（180–190 行），196–198 行 `for pk, pv := range params { endpoint = string(ByteReplace([]byte(endpoint), []byte(pk), []byte(pv), 1)) }`：每替换一个键，都把当前 endpoint 转 `[]byte`、替换后再转回 `string`，整串复制一遍；map 遍历顺序还不稳定。
- 影响评估：endpoint 占位符替换工具，调用频率低；键值多时逐键全串复制，复杂度随键数 × 串长增长。
- 优化建议：把成对键值交给 `strings.NewReplacer(...)` 一次构建（与 `Tag`、`Date` 的修复方向一致），单趟扫描完成全部替换；`ByteReplace` 的字节级原语保留给需要就地修改的场景。

#### [P2] `util/util.go:148` `GetMissingElement`：构造完整 1..max 切片 + 排序 + 全等比较，用 O(n log n) 与额外分配解决 O(n) 问题
- 问题：`GetMissingElement` 152 行先 `sort.Sort(intSlice(arr))`；153–157 行再构造 `arrTemp` 从 1 到切片最大值的全部整数；158 行 `equal(arrTemp, arr)` 全量比较；不相等时 161–165 行再遍历一次定位缺口。找缺失最小整数只需一次哈希集合（或排序后一次扫描），当前实现额外产生一个等长切片、一次排序与两次全量遍历。
- 影响评估：工具函数，仓内未见生产调用；仅在大批量连续整数切片上调用时，排序与等长切片分配才明显。
- 优化建议：用 `map[int]struct{}` 标记已有数字，从 1 起探测首个缺失值，O(n) 时间、O(n) 空间且无需排序；若输入可保证有序，直接单趟扫描即可。

### util/var.go

#### [P2] `util/var.go:12` `Empty`：反射入口 + 非基本类型兜底 `DeepEqual`/`Zero` 分配
- 问题：`Empty` 16 行先 `reflect.ValueOf(val)`；string/map/slice/bool/数值等基本类型在 switch 内处理；34 行对未命中分支统一 `reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())`，`reflect.Zero(...).Interface()` 会分配对应零值包装对象。函数不接收任何具体类型快路径，所有入参都先付一次 `reflect.ValueOf`。
- 影响评估：空值判断工具，仓内调用点多为 `toolkit` 代码生成与 `stores/mongodb/clientmanager.go:119` 启动期配置判断；若未来在请求参数循环中对 `string`/`map`/`slice` 频繁判空，反射入口与兜底分配会叠加。
- 优化建议：在反射之前增加类型断言快路径（`string` 判长度、`[]T`/`map[K]V` 判长度与 nil、`bool`、数值零值直接比较），反射仅作为结构体、接口等兜底；`DeepEqual` 兜底分支保留。

### util/init.go

#### [P2] `util/init.go:12` `Do`：重试入口每次都新建随机源生成抖动
- 问题：`Do` 13 行 `r := rand.New(rand.NewSource(time.Now().UnixNano()))`，17 行 `r.Int63n(int64(sleep))` 计算抖动；与 `math.go:22` `Rand` 同款「每次新建源」问题。递归调用 `Do(attempts, 2*sleep, f)`（20 行）时每一层再次新建源。
- 影响评估：带退避抖动的重试包装，唯一生产调用点为 `registry/etcdv3/registry.go:399` `registerAllKvs` 的注册重试；失败路径才产生等待与抖动，正常成功调用不触碰随机源。主要问题是每次重试新建源的分配与同纳秒种子重复风险。
- 优化建议：与 `util.Rand` 统一为包级复用随机源；抖动用共享源 `Int63n` 取得。

### util/exec.go

#### [P2] `util/exec.go:57` `System`：stdout/stderr 拷贝 goroutine 未纳入 `WaitGroup`，`cmd.Wait()` 返回时缓冲可能丢失；`stdBuf` 无界增长
- 问题：`System` 86–87 行用 `io.MultiWriter(os.Stdout, &stdBuf)` 包装管道；95–100 行启动两个 goroutine 分别 `io.Copy` stdout/stderr；102 行 `cmd.Wait()` 直接等待进程退出，没有等待两个拷贝 goroutine 结束，`err1`/`err2` 在 103 行被读取时 goroutine 可能仍在写。`stdBuf`（59 行）无容量上限，命令输出很大时内存随输出增长。
- 影响评估：系统命令执行工具，非框架请求热路径；本条目主要是并发正确性隐患（goroutine 与主流程不同步），性能角度为无界缓冲与悬挂拷贝 goroutine；`Exec`（17 行）用 `CombinedOutput` 没有该问题。
- 优化建议：用 `sync.WaitGroup` 等待两个 `io.Copy` goroutine 结束后再返回最后一行；对 `stdBuf` 设置大小上限或按行裁剪；输出大时直接透传不缓存。

### registry/etcdv3

#### [P1] `registry/etcdv3/registry.go:81` `ListServices`：同一批 KV 被两轮循环各反序列化一次
- 问题：`ListServices` 拿到 `getResp.Kvs` 后，91–100 行第一轮把每个 `kv.Value` `json.Unmarshal` 到 `server.ServiceInfo` 并追加；102–115 行第二轮又把**同一份** `kv.Value` `json.Unmarshal` 到 `registry.Update`，只取 `Addr` 再包装成 `ServiceInfo{Address: service.Addr}` 追加。每条节点数据被 JSON 反序列化两次、分配两个目标结构；两轮还分别打日志。
- 影响评估：服务发现启动/重连时的全量拉取路径，节点数 N 决定成本；N 较大（集群多实例、多服务前缀）时双倍 JSON CPU 与双倍对象分配在一次列表刷新内集中发生，延迟与 GC 压力翻倍。
- 优化建议：单次反序列化到能同时容纳两类字段的结构（`registry.Update` 内嵌 `Metadata *server.ServiceInfo`，见 `registry/update.go:17–24`），一轮解析后按需产出 `ServiceInfo` 列表；或合并两轮循环在同一遍历中完成；避免对同一 `kv.Value` 连续两次 `Unmarshal`。

#### [P2] `registry/etcdv3/registry.go:124` `WatchServices`：每个 PUT/DELETE 事件都对全量四个 map 做 `DeepCopy`
- 问题：监听 goroutine 151 行起，每个事件先 `updateAddrList` / `deleteAddrList` 改动累计快照 `al`，159 行 `out := al.DeepCopy()`，随后经 161 行推送给消费方；`DeepCopy`（`registry/endpoint.go:40`）内部 `newEndpoints()` 创建四个空 map 并 `maps.Copy` 全部条目（51–56 行）。即每来一个节点变更事件，都把当前全量节点/路由/消费者/提供者映射整体拷贝一遍。
- 影响评估：服务发现推送路径，节点规模小、变更不频繁时成本可忽略；事件密集（批量扩缩容、会话重连刷事件）或节点数上千时，单事件全量拷贝叠加在推送通道上，且通道容量仅 10（135 行），丢弃事件时只记 `invalid event` 告警。
- 优化建议：消费方若只需增量，可推送 `Update` 事件（新增/删除 + 单节点信息）而非全量快照；保留全量推送时，评估 `maps.Copy` 的浅拷贝是否满足下游需求；事件合并（同一周期内多个 PUT 合并为一次快照推送）可降低拷贝频次。

#### [P2] `registry/etcdv3/registry.go:187` `Close`：每个缓存键启动一个注销 goroutine，并发数量无上限
- 问题：`Close` 192 行 `reg.kvs.Range` 遍历缓存，194 行每个 key `go func(k any){...}` 注销，191 行 `sync.WaitGroup` 只负责等待；没有并发度限制。注册节点多时一次性拉起等量 goroutine 同时发起 etcd 删除请求。
- 影响评估：进程关闭/注销收尾路径，不在请求热路径；问题是收尾阶段的突发并发与 etcd 短时间内批量 Delete 压力，属于资源使用形态问题。
- 优化建议：用固定 worker 池（如 semaphore 通道限制并发数，N=10~20）消费待注销键；或分批串行注销；每个 goroutine 内 196 行 `context.WithTimeout` + 210 行 `cancel()` 已配对，保持现有释放方式。

### registry（根目录）

#### [P2] `registry/config.go:37` `RegisterBuilder` 与 `registry/config.go:61` `InitDefaultRegister`：`registryBuilder` 裸 map 无互斥保护
- 问题：`registryBuilder`（9 行）是包级 `map[string]Builder`；`RegisterBuilder` 38 行直接读、41 行直接写；`InitDefaultRegister` 67 行直接遍历读取。全程没有 `sync.RWMutex` / `sync.Map` 保护。Go map 在并发读写下会触发运行时 fatal。
- 影响评估：当前构造流程是启动期各模块 `init()` 阶段顺序注册，正常生命周期内不构成并发竞争；问题是缺少保护后的脆弱性——若后续支持运行期动态注册注册器，map 读写会立刻进入数据竞争。性能层面无锁本身是快的，本条按并发稳健性登记 P2。
- 优化建议：`registryBuilder` 改为 `sync.Map`，或加一把 `sync.RWMutex`；注册用写锁、查找用读锁，与 `etcdv3Registry` 已有的锁粒度保持一致。

### snowflake

#### [P1] `snowflake/snowflake.go:77` `Config.GenInt64ID`：CAS 冲突后固定 `time.Sleep(20)` 退避；序列号耗尽分支再次取系统时钟
- 问题：`GenInt64ID` 是 ID 生成热路径：79 行 `atomic.LoadInt64`，83–90 行判定时间片与序列号，93 行 `CompareAndSwapInt64`；CAS 失败进入 96 行 `time.Sleep(time.Duration(20))` 后自旋重试。冲突协程每次都固定睡眠约 20ns 再重试，没有先本地短自旋；85–87 行序列号耗尽分支又调用一次 `time.Now().UnixNano()`（86 行）计算等待窗口，而 `getCurrentTimestamp`（68–71 行）在 82 行已经取过一次 `time.Now().UnixNano()`。
- 影响评估：雪花 ID 通常被用于订单号、消息 ID、日志追踪等**每请求/每事件**都调用一次的路径；单机并发分配同一毫秒段时，CAS 冲突率随并发度上升，固定小睡会把冲突协程的调度延迟放大；86 行在锁外无并发问题，但属于重复系统时钟读取。实际冲突率与 20ns 睡眠的影响程度取决于线上并发分配速率，**需 benchmark 实测**（见第四节第 2 条）。
- 优化建议：
  1. CAS 冲突后先本地自旋若干次（如 3~5 次不睡眠重试），仍失败再进入 `time.Sleep` 退避，减少高并发下的协程睡眠/调度开销；
  2. 82 行取得的 `now` 时间向下传递，序列号耗尽分支复用同一时刻计算对齐等待，避免第二次 `time.Now()`；
  3. 评估 20ns 固定睡眠是否过短导致无效重试、或过长导致尾延迟，结合 pprof 调度 profile 调整；原子操作与位运算本身实现正确，无需改锁。

### config

#### [P2] `config/config.go:51` `LoadDir`：Walk 回调内每个文件新建扩展名切片并调用反射型 `InArray`
- 问题：`LoadDir` 在 52 行 `filepath.Walk`，回调内 56–57 行 `ext := make([]string, 0); ext = append(ext, "json", "toml", "yaml")` 为**每个访问到的文件**重新构造白名单切片；64 行 `util.InArray(strings.TrimPrefix(fileExt, "."), ext)` 走 `util.InArray` 反射 + `DeepEqual` 做包含判断。
- 影响评估：配置目录一次性加载流程，目录文件数量有限；浪费是每文件一次小切片分配 + 一次反射包含判断，量级不大。
- 优化建议：白名单提为包级 `map[string]struct{}{"json":{}, "toml":{}, "yaml":{}}`，回调内直接 map 查找 O(1)；删除每文件的切片构造与 `InArray` 反射调用。

### crontab

#### [P1] `crontab/task.go:459` `TaskManager.run` 与 `crontab/task.go:658` `MapSorter.Less`：每轮调度循环重建排序器、全量排序；比较器每次取最多 4 次 `GetNext`
- 问题：调度主循环 `run` 466–468 行在 `RLock` 内 `NewMapSorter(t.adminTaskList)`（467 行，内部 637–647 行新建两个等长切片并逐个 append），469 行 `sortList.Sort()` 对全部任务做 O(n log n) 排序；只要循环被唤醒（到期、`changed`、停止）就重建并排序一次。排序比较器 `Less` 658–666 行：659 行取一次 `GetNext(...).IsZero()`，662 行再取一次，665 行两侧又各取一次 `GetNext(...).Before(...)`，即**每次比较最多 4 次接口方法调用**，而 `Next` 时间戳在排序期间恒定；`Task.GetNext`（206–208 行）只是字段返回，接口分派 + 方法调用叠加在 O(n log n) 次比较上。
- 影响评估：定时任务调度器，任务数通常几十个以内，单次排序成本有限；但 `run` 循环在 `time.After` 到期后立即重排（479–481 行），高频秒级任务场景下每个调度周期都重复「切片重建 + 全量排序 + 重复 GetNext」；比较器多次取同一字段属于可直接消除的固定开销。
- 优化建议：
  1. 构造 `MapSorter` 时一次性把每个任务的 `Next time.Time` 缓存进并行切片（与 `Vals` 对齐），`Less` 直接比较缓存时间戳，删除 659、662、665 行的重复 `GetNext` 调用；
  2. 仅在任务列表 `changed`（增删改）时重建排序；到期触发只需要找到 `Next <= now` 的任务，可维护最小堆（`container/heap`）替代每轮全量排序，O(n log n) 降为堆顶维护 O(log n)；
  3. `NewMapSorter` 已按 `len(m)` 预分配容量（639–640 行），保留该写法。

#### [P2] `crontab/task.go:550` `StopTask`、`:557` `GracefulShutdown` 与 `:569`/`:592`/`:610` `AddTask`/`DeleteTask`/`ClearTask`：无缓冲 `stop`/`changed` 通道在独立 goroutine 中发送，调度循环退出后发送方永久阻塞
- 问题：`TaskManager` 的 `stop`（41 行）与 `changed`（42 行）均为无缓冲通道。`StopTask` 551 行 `go func(){ t.stop <- true }()`；`GracefulShutdown` 559 行同样 `go func(){ t.stop <- true; t.wait.Wait(); close(done) }()`；`AddTask` 580–582 行、`DeleteTask` 603–605 行、`ClearTask` 621–623 行在调度器已启动时各自 `go func(){ t.changed <- true }()`。`run` 循环在 487 行收到停止信号后 488 行 `markManagerStop` 并 `return`，两个通道都不再有接收方；此后任何一次停止或增删任务触发的发送 goroutine 都会永久阻塞在通道发送上，形成 goroutine 泄漏。
- 影响评估：停止/重启与任务动态增删的生命周期路径，非任务执行热路径；泄漏规模取决于调度器退出后重复调用停止、增删任务的次数，本条按并发资源稳健性登记 P2（主因是阻塞 goroutine 占用调度槽位，而非任务执行期开销）。
- 优化建议：停止信号改用 `context.Context` + `cancel()`（幂等、可重复调用、无阻塞发送）；`changed` 改为缓冲为 1 的通道，发送侧用 `select { case t.changed <- true: default: }` 非阻塞合并信号；`stop` 同理改为缓冲 1 + 非阻塞发送；`GracefulShutdown` 的等待逻辑保持 `wait.Wait()` 即可。

### watch

#### [P2] `watch/builder.go:72` `builder.isIgnore`：每个文件事件都依次跑 excludes 的 `filepath.Match` 与 exts 后缀循环
- 问题：`isIgnore` 81–90 行对 `b.excludes` 逐个 `filepath.Match(p, path)`；95–99 行再对 `b.exts` 逐个 `strings.HasSuffix(path, ext)`。fsnotify 每个创建/写入/重命名事件都完整走一遍两层线性扫描；`filepath.Match` 每次还要解析通配模式。
- 影响评估：`watch` 是开发期热重载工具，事件量与频率都不在生产热路径上，定 P2；仅在大型 monorepo 密集保存文件时，匹配循环与重复 Match 解析会叠加在编译触发前的判断里。
- 优化建议：启动时把 excludes 模式预编译成 `regexp` 或 `filepath.Glob` 已解析结构；exts 转成 `map[string]struct{}` 后缀集合，匹配 O(1)；`filterPaths`（199–225 行）与事件回调复用同一套判定。

### img/img.go

#### [P2] `img/img.go:37` `OpenLocalFile`：同一磁盘文件被读取两到三次
- 问题：45 行 `os.ReadFile(filename)` 已把整个文件读入 `fileBuffer`；49 行用该 buffer 判断类型；webp 分支 58 行又 `bimg.Read(filename)` 第二次从磁盘读入，62 行转 JPEG、66 行写出临时文件；73 行 `imaging.Open(filename)` 第三次从磁盘解码（webp 时读的是转换后的临时 jpeg，非 webp 分支则再次读原始文件）。
- 影响评估：图像处理工具，批量缩略图/水印流水线中每次打开都重复读盘；非 webp 分支 45 行的 `fileBuffer` 解码后并未复用，直接浪费了一次全量读入。
- 优化建议：优先用 `imaging.Decode(bytes.NewReader(fileBuffer), opts...)` 从已加载字节解码；webp 转换也基于 `fileBuffer` 构造 `bimg.NewImage`，避免第二次 `bimg.Read`；确认无需磁盘临时中转后再决定是否保留 `imaging.Open` 的文件路径。

#### [P2] `img/img.go:245` `Image.Compress`：JPEG 编码进缓冲后立刻再解码回 RGBA，全图编码 + 解码往返
- 问题：`Compress` 246 行 `bytes.Buffer`，247 行 `jpeg.Encode(&buffer, img.src, ...)`；251 行 `imaging.Decode(&buffer)` 把刚压缩的 JPEG 字节再解回内存位图，252 行替换 `img.src`。函数名义上是压缩，实际产物仍是未压缩的 RGBA 位图，压缩结果只经过内存缓冲便被丢弃；全图经历一次有损编码 + 一次完整解码。
- 影响评估：图片处理流水线中的环节；被链式调用（`Resize` → `Compress` → `Save`）时，每张图多一次全图 JPEG 编解码 CPU 与一次位图内存分配。是否在在线图片服务热路径需结合接入方式确认（见第四节第 4 条）。
- 优化建议：压缩语义应为「输出 JPEG 字节流」：`Compress` 返回编码后的 `[]byte` 或直接写入目标 `io.Writer`；若链式处理必须保留位图，把压缩放到最后保存阶段，不要在中间往返编解码。

### img/fontx

#### [P2] `img/fontx/font.go:45` `TrueTypeFont.GetMetrics`：每次调用都用当前字号新建 `truetype.Face`，不做缓存
- 问题：`GetMetrics` 46–50 行构造 `truetype.Options{Size, DPI, Hinting}`，51 行 `truetype.NewFace(fo.Font, opt)`，52 行取 `face.Metrics()` 后未保留 face；同一字体、同一字号反复取度量时重复构建排版面。
- 影响评估：字体度量获取通常在文字绘制初始化阶段；与 `img/img.go:261` `DrawText` 配合，批量绘制不同字号文字时会重复建 face。本包静态阅读，未编译运行。
- 优化建议：按 `(font, size, dpi, hinting)` 为键在 `TrueTypeFont` 上缓存 `*truetype.Face`（`sync.Map` 或互斥锁保护）；face 复用后度量读取零构建成本。

---

## 二、条目统计

| 子包 | P0 | P1 | P2 | 小计 |
|---|---|---|---|---|
| util/string.go | 0 | 2（`JSONString`、`Tag`） | 6（`UCWords`、`StrIPos`/`StrRIPos`、`Md5File`、`Sha1File`、`SimilarText`、`StrTr`） | 8 |
| util/array.go | 0 | 3（`InArray`、`ArrayRemoveRepeatedElement`/`ArrayStringUniq`、`SplitString`） | 3（`InMultiArray`、批量切片/map 扩容与 `fmt.Sprintf`、`ArrayRand`/`ArrayRandMap` 随机源重建；`StrShuffle` 同模式见本条位置列表） | 6 |
| util/math.go | 0 | 1（`Rand`） | 0 | 1 |
| util/time.go | 0 | 1（`Date`/`StrToTime`） | 2（`Hour`、`ToWeekDay`） | 3 |
| util/misc.go | 0 | 0 | 2（`MemoryGetUsage`/`MemoryGetPeakUsage`、`VersionCompare`） | 2 |
| util/net.go | 0 | 0 | 1（`isPrivateIP`） | 1 |
| util/mongodb.go | 0 | 0 | 2（`Pipeline`、`ConvertBson`） | 2 |
| util/pinyin.go | 0 | 0 | 1（`ZhCharToFirstPinyin`/`ZhCharToPinyin`） | 1 |
| util/scprng.go | 0 | 0 | 1（`StrPad`） | 1 |
| util/util.go | 0 | 1（`VerifyPhone`/`VerifyEmail`/`VerifyIPv4`/`VerifyIPv6`） | 2（`NewReplacer`、`GetMissingElement`） | 3 |
| util/var.go | 0 | 0 | 1（`Empty`） | 1 |
| util/init.go | 0 | 0 | 1（`Do`） | 1 |
| util/exec.go | 0 | 0 | 1（`System`） | 1 |
| registry/etcdv3 | 0 | 1（`ListServices`） | 2（`WatchServices`、`Close`） | 3 |
| registry（根目录） | 0 | 0 | 1（`RegisterBuilder`/`InitDefaultRegister` 裸 map） | 1 |
| snowflake | 0 | 1（`GenInt64ID`） | 0 | 1 |
| config | 0 | 0 | 1（`LoadDir`） | 1 |
| crontab | 0 | 1（`run`/`MapSorter.Less`） | 1（`StopTask`/`GracefulShutdown`/`AddTask`/`DeleteTask`/`ClearTask` 通道异步发送） | 2 |
| watch | 0 | 0 | 1（`isIgnore`） | 1 |
| img/img.go | 0 | 0 | 2（`OpenLocalFile`、`Compress`） | 2 |
| img/fontx | 0 | 0 | 1（`GetMetrics`） | 1 |
| worker（`worker.go`、`job/job.go`） | 0 | 0 | 0 | 0 |
| **合计** | **0** | **11** | **32** | **43** |

> 说明：合并条目在统计中按「一条」计；如 `ArrayRemoveRepeatedElement`/`ArrayStringUniq`、`Date`/`StrToTime`、四个 IP/手机号/邮箱校验函数等多位置同模式问题已合并为单条，位置列表见正文。

---

## 三、高价值候选列表（Top 10，供全局 Top 10 评选）

| # | 位置 | 一句话问题 | 预期收益 |
|---|---|---|---|
| 1 | `util/math.go:22` `Rand` | 每次选库都新建随机源 + 种子；已被 `stores/proxy/sql.go:36`、`:46` 每次读写库选择调用 | 每条 SQL 查询省去一次 `time.Now` + 随机源/Rand 对象构造与初始化；改动极小（换共享源），消除每查询固定分配 |
| 2 | `util/string.go:637` `JSONString` | Marshal 后经 `cast.ToString(any)` 多一层接口装箱 + 类型 switch；gRPC 访问日志每条调用 2~3 次 | 每个被记录的 gRPC 请求省去 2~3 次切片头装箱分配与 cast 层间接调用；一行改为 `string(jsByte)`，payload 拷贝次数不变但固定开销归零 |
| 3 | `util/util.go:93`/`:101`/`:109`/`:115` 四个校验函数 | 手机号/邮箱/IPv4/IPv6 正则每次请求现场编译，IPv6 模式尤其重 | 接入参数校验后，每请求省去一次正则编译（语法分析 + 解析树 + NFA 构建）；包级预编译一次即可，标准库已有同包 `util/regexp.go` 范例 |
| 4 | `util/array.go:352`/`:473` 两个去重函数 | O(n²) 双重循环逐元素比较去重 | 批量标签/列表去重从平方比较降为单趟 O(n) map 标记；消除长切片上的平方级字符串比较与反复切片扩容 |
| 5 | `registry/etcdv3/registry.go:81` `ListServices` | 同一份 KV 值被两轮循环各 JSON 反序列化一次 | 服务发现全量刷新时 JSON CPU 与对象分配直接减半；合并为一轮解析，节点数越多收益越明显 |
| 6 | `util/time.go:153` `Date` / `:141` `StrToTime` | 静态 `datePatterns` 每次重建 Replacer，且被全部日期格式化封装复用 | 每次时间格式化省去 16 对替换项的 Replacer 构建与内部分配；一次包级缓存修复同时惠及 10 余个 `Format*` 封装 |
| 7 | `snowflake/snowflake.go:77` `GenInt64ID` | CAS 冲突固定睡 20ns，序列号耗尽分支重复取时钟；ID 生成属每请求热路径 | 冲突路径先短自旋再退避，降低高并发分配下的协程调度/睡眠延迟；复用时间戳消除重复系统调用；收益需并发 benchmark 量化 |
| 8 | `util/array.go:223` `InArray` | 反射 `ValueOf` + `Interface()` 装箱 + `DeepEqual` 三层通用包含判断，无标量快路径 | `[]string`/`[]int` 等常见场景改直比/map 查找后，每次包含判断从反射三层降为 O(1) 哈希或 O(n) 直比；波及 30+ 现有调用点 |
| 9 | `crontab/task.go:459` `run` / `:658` `Less` | 每轮调度重建切片并全量排序，比较器重复取最多 4 次 `GetNext` | 排序期缓存 Next 时间戳消除重复接口调用；堆化调度后高频任务周期不必每轮全量重排；任务数多时减少排序与分配 |
| 10 | `util/array.go:308` `SplitString` | 逐字符 `res +=` O(n²) 拼接 + 每次 MustCompile 正则 + 双层冗余 map | 长文本清洗拼接从平方拷贝降为 Builder 线性追加；正则包级预编译；单层 map 去重删除多余类型转换链 |

---

## 四、拿不准、需实测验证的条目清单

1. **`util/misc.go:56`/`:65` 内存统计函数的调用频率**：仓内仅有测试覆盖，是否被线上指标采集器按秒级/分钟级调用未知；`runtime.ReadMemStats` 的 STW 成本只有在高频调用下才显著，需确认接入方调用频率后再决定是否值得加复用缓冲。
2. **`snowflake/snowflake.go:77` `GenInt64ID` 的 CAS 冲突率与 20ns 睡眠影响**：需在目标并发度下对 `GenInt64ID` 写 benchmark（`b.RunParallel`，不同 QPS 档位），观察 CAS 重试次数、协程调度等待与尾延迟；先自旋后睡眠的改造是否优于现状，以 goroutine/调度 profile 与 CPU profile 为准。序列号耗尽分支 86 行重复读取 `time.Now().UnixNano()` 是否构成可观测开销，同样以 profile 结果为准。
3. **`util/string.go:818` `SimilarText` 是否存在在线调用方**：仓内未检索到生产调用；若后续接入模糊匹配/查重循环，需先在目标文本长度分布下测量 O(n²·m) 嵌套循环与递归字节拷贝的实际耗时，再决定是否替换算法。
4. **`img/img.go:245` `Compress` 是否处于在线图片服务流水线**：本包在本机无法编译；需确认调用方是离线批处理还是请求内链式处理，再评估 JPEG 编解码往返的 CPU 占比；若仅离线脚本使用，可维持 P2 不动。
5. **`util/math.go:22` `Rand` 及 `ArrayRand`/`ArrayRandMap`/`StrShuffle`/`init.go:13` `Do` 的随机源改造选型**：全局 `math/rand` 自带互斥锁，高并发下锁竞争与「每次新建无锁源」孰快孰慢需 benchmark；`math/rand/v2` 无锁全局源是候选方案；同时需验证 `UnixNano` 种子在短时间并发批量调用时是否会产生重复随机序列/相同排列（正确性一并验证）。
6. **`util/mongodb.go:16` `Pipeline` 是否在查询路径动态构造**：若管道在启动期构造一次缓存复用，三重编解码开销不进入查询热路径；需结合 MongoDB 调用方确认管道构造时机，再定优先级。
7. **`registry/etcdv3/registry.go:159` 全量 `DeepCopy` 的实际事件频率**：线上节点规模与变更频率决定全量拷贝成本；建议在 etcd 监听回调处打点统计每事件节点数与拷贝耗时，再决定是否改造为增量推送。
8. **`watch/builder.go:72` `isIgnore` 匹配开销**：开发期工具，仅在大型项目密集保存触发编译时才可能感知；不需要为该包做生产级优化，条目按静态阅读登记。

---

## 五、跨模块提示（超出本章节范围，仅登记线索，不展开审查）

1. **`util.JSONString` 的调用侧门控**：`client/grpc/interceptor.go:144`、`:145`、`:156`、`:157`、`:169`、`:170` 在错误 Warn 分支与 info 级别分支无条件对 `req`、`reply` 做 JSON 序列化；即便日志级别低于 info，错误分支仍会序列化两个消息体。`json.Marshal` 的序列化成本（反射 + 整包字节分配）远大于本章节函数内 cast 间接层；建议调用侧在日志级别/采样开关关闭时跳过 `JSONString` 构造，本章节只去掉函数内的 cast 装箱间接层。
2. **`util.InArray` 的现有调用分布**：`client/etcdv3/config.go:98`（启动期端点去重）、`img/img.go:51`（格式白名单）、`config/config.go:64`（扩展名白名单）、`util/mongodb.go:48`（bson 键判重），外加 `toolkit/` 代码生成包内约 30 处（`toolkit/base/util.go`、`toolkit/api/api.go`、`toolkit/dao/dao.go`、`toolkit/vue/vue.go`、`toolkit/module/module.go`）；其中 `toolkit` 为离线生成，前三者为启动期，当前均非请求热路径；改造 `InArray` 快路径时这些点可作为回归对照，无需改动调用方。
3. **`util.Rand` 的调用侧**：`stores/proxy/sql.go:36` `Write()`、`:46` `Read()` 为每条 SQL 选库调用；该代理切片在启动期装配完成，若后续有每请求多次选库的用法，随机源复用收益会进一步放大，属 `stores/proxy` 模块审查时复核的线索。
4. **`registry/registry.go:113` `GetServiceValue` 与 `:120` `GetService`**：`json.Marshal` 后 `string(val)`、`json.Unmarshal([]byte(s), ...)` 是注册/发现路径的标准序列化，单次调用频率低，未单独立条目；若发现 `ListServices` 之外的循环调用点，可与本章节 `ListServices` 双倍反序列化问题一并优化。
5. **`util/pb.go` 反射式 protobuf 结构互转**：`StructToMap`、`MapToStruct`、`StructToPbStruct`、`labEntry` 全量走 `reflect` 与 `structs.Map`；仓内未检索到生产调用方，若未来在 gRPC 元数据透传热路径使用，需另立条目评估反射转换开销，本章节不展开。
6. **`util/exec.go:17` `Exec`**：使用 `CombinedOutput` 一次性收集输出，未发现 `System` 那样的 goroutine 同步问题；两条系统命令实现并存，建议后续统一为带 `WaitGroup` 的流式收集。

---

## 六、已通读但未发现性能问题的包 / 导出 API（覆盖说明）

- **`registry/endpoint.go`**：`Endpoints` 四个 map 结构定义；`DeepCopy`/`DeepCopyInfo` 使用标准库 `maps.Copy` 浅拷贝，无循环内分配；拷贝策略本身被 `WatchServices` 条目（第一节）引用讨论，结构实现本身无额外问题。
- **`registry/nop.go`**：`Local` 本地空实现，注册/注销仅打日志，无外部资源与锁；`ListServices`/`WatchServices` 按设计 panic，不属于性能范畴。
- **`registry/service.go`**：`Service` 仅为 `server.Server` 接口包装，`Serve` 直接委托，无额外层开销。
- **`registry/update.go`**：纯 `Operation` 常量与 `Update` 结构体定义，无可执行逻辑。
- **`registry/etcdv3/config.go`**：`Config` 配置 DTO 与 `New`/`SetNode`/`SetPrefix`/`Build`/`Singleton`/`MustSingleton`；`Singleton` 经 `core/singleton` 加载/存储单例，缓存命中路径短；无请求热路径代码。
- **`snowflake/snowflake.go` 非热路径部分**：`ConstructConfig`/`ConstructConfigWithMachineID` 为启动期构造；`getMachineID`、`getIP`、`getEnv` 仅初始化阶段执行一次；`GetSeqFromID`/`GetTimeFromID` 为纯位运算内联；`getCurrentTimestamp` 一次纳秒读取 + 掩码，无额外分配。`getStartEpochFromEnv` 的环境变量语义写反属已知正确性 FIXME，按任务要求不在本性能章节重复上报。
- **`config/config.go` 读取封装**：`String`/`Int`/`Uint`/`Int64`/`Duration`/`Float`/`Bool` 先 `Exists`（`IsSet`）再取值，两次 Viper 查询；配置读取以启动期为主，两次 map 查找开销可忽略，未单独立条目。`BindStruct` 委托 `UnmarshalKey`，无额外处理。
- **`util/file.go`**：全部为 `os.Stat`/`os.ReadFile`/`os.OpenFile`/`filepath` 等标准库薄封装；`FileMTime`（215 行）为取修改时间先 `os.Open` 再 `fd.Stat()`，可用 `os.Stat` 直接替代减少一次打开，但属低频文件管理操作，未凑数单列；`Copy`（121 行）使用 `io.Copy` 标准缓冲，方式正确；`FGetCsv` 委托 `csv.Reader.ReadAll`。
- **`util/url.go`**：`ParseURL`/`URLEncode`/`URLDecode`/`RawURLEncode`/`RawURLDecode`/`HTTPBuildQuery`/`Base64Encode`/`Base64Decode` 均为标准库一层封装，无自定义循环、无正则、无反射；`RawURLEncode` 一次 `strings.Replace` 可接受。
- **`util/grpc.go`**：`GrpcTime` 纯 `time.Unix(...).Local()` 转换，无分配热点。
- **`util/regexp.go`**：全部模式在包级 `var` 块（68–127 行）`MustCompile` 预编译，`regexpMatch` 单委托 `FindAllString`；实现正确，是本章节多个「应预编译」条目的正面范例。
- **`util/string.go` 未单列的其余导出函数**：`StrReplace`、`Nl2br`、`AddSlashes`、`StripSlashes`、`QuoteMeta`、`HtmlEntities`、`StrRepeat`、`StrStr`、`Explode`、`Implode`、`Chr`、`Ord`、`Soundex`、`ChunkSplit`、`WordWrap`、`NumberFormat`、`StrWordCount`、`UCfirst`/`LCFirst`、`SubStr`、`StrRev`、`Md5`/`Sha1`/`Crc32` 等均为标准库一层包装或单趟线性扫描：`Nl2br`/`AddSlashes`/`QuoteMeta` 已使用 `bytes.Buffer`/`strings.Builder`；`StrRepeat` 委托 `strings.Repeat`；`NumberFormat`（346 行）已按长度预分配 `tmp`；未发现循环内反复分配或二次算法浪费。`StrStr`（474 行）`len([]byte(needle))` 等价于 `len(needle)`，属冗余类型转换，开销可忽略，未单列。
- **`util/array.go` 未单列的其余导出函数**：`ArrayFill`/`ArrayFlip`/`ArrayKeys`/`ArrayValues`/`ArrayMerge`/`ArrayChunk`/`ArrayPad`/`ArraySlice`/`ArrayReverse`/`ArrayPush`/`ArrayPop`/`ArrayUnshift`/`ArrayShift`/`ArrayCombine`/`ArrayKeyExists`/`ArrayColumn`/`ArrayKeyPluck`/`ArrayMultiPluck`/`Implode`/`InterfaceToString` 之外的类型转换函数，均为单趟切片/map 操作；`ArrayMerge`（58–62 行）先统计总长再预分配，`ArrayKeys`/`ArrayValues`（38、48 行）按 `len(elements)` 预分配，实现方式正确。
- **`util/math.go` 未单列函数**：`Abs`/`Round`/`Floor`/`Ceil`/`Pi`/`Max`/`Min`/`DecBin`/`BinDec`/`Hex2bin`/`Bin2hex`/`Dechex`/`Hexdec`/`Decoct`/`Octdec`/`BaseConvert`/`IsNan`/`Divide`/`Add` 均为纯函数或标准库一层包装，无循环、无分配热点。
- **`util/time.go` 未单列函数**：`DateTimeParse`、`FormatDuring`、`GetHourDiffer`、`TimeZone`、`Time`、`ToTime`、`Duration`、`Now`、`Timestamp`、`CheckDate`、`Sleep`/`USleep`、`DateToJulian`/`JulianToDate`、`ISOWeekday`、`StartDate`、`FromYearWeek`、`ordinalInYear` 等为纯计算或标准库一层转换；`Month`/`Day`/`Hour`/`ToWeekDay` 的重复 `Date` 问题已由本节 `Hour`、`ToWeekDay` 两条目覆盖。
- **`util/misc.go` 未单列函数**：`Echo`、`Uniqid`、`Exit`/`Die`、`Getenv`/`Putenv`、`Pack`/`Unpack`、`Ternary` 为极薄包装；`Pack`/`Unpack` 单次二进制读写，未在循环路径上。
- **`util/net.go` 未单列函数**：`GetHostName`、`GetHostByName`、`GetHostByNameL`、`GetHostByAddr` 为显式 DNS 查询封装，`IP2long`/`Long2ip` 为纯字节序转换；网络解析成本由调用方语义决定，非库内可优化点。
- **`util/mongodb.go` 未单列函数**：`indirect`（73–86 行）仅在指针类型上走反射解引用，普通值直接原样返回，实现已做快路径判断。
- **`img/img.go` 未单列的变换函数**：`OpenReader`、`Clone`、`Resize`、`Crop`、`Blur`、`Gray`、`AdjustContrast`、`Sharpen`、`Invert`、`Convolve3x3`/`Convolve5x5`、`AdjustBrightness`/`AdjustGamma`/`AdjustSaturation`/`AdjustSigmoid`、`Rotate`、`Transverse`、`Fit`、`Fill`、`Paste`、`Overlay`、`AddWaterMark`、`Save`/`SaveToFile` 均为 `imaging` 库一行委托；位图分配是图像处理固有成本，包装层未引入额外重复拷贝。`DrawText`（261 行）每次新建 `freetype.Context`，与 `fontx.GetMetrics` 建 face 问题同源，已并入 `fontx` 条目讨论。
- **`crontab/task.go` 未单列部分**：cron 解析器 `parse`/`parseSpec`/`getField`/`getRange`/`getBits`/`parseIntOrName` 仅在任务注册 `SetCron` 时执行一次；`Schedule.Next`（342 行）按时间分量步进查找，仅在 `SetNext` 与任务触发后调用，不在排序比较路径上；`GetStatus`（169 行）已使用 `strings.Builder`。
- **`util/pb.go`**：见第五节第 5 条，当前无生产调用，未立性能条目。
- **`watch/logger.go`**：`loggerWriter.Write`（52 行）把 `[]byte` 转 `string` 交给 Logger，属开发期日志通道；`consoleLogger.Output` 按类型选 writer 输出，无锁、无缓存问题。
- **`watch/options.go`**：`sanitize`/`sanitizeExts`/`getAppName`/`recursivePaths`/`splitArgs` 均为启动期配置归一化；`recursivePaths` 目录展开只执行一次；无运行期热路径。
- **`watch/watch.go`**：`GoVersion` 启动期执行一次 `go version`；`Watch`/`watch`/`initWatcher` 为开发期监听循环，事件防抖 `watcherFreq`（75 行）已限制短时间重复编译；具体匹配开销见 `watch/builder.go:72` 条目。
- **`worker/worker.go`、`worker/job/job.go`**：纯 `Worker`、`Runner` 接口定义，无实现代码与状态，不存在性能问题。
- **测试文件**：`util/*_test.go`、`registry/*_test.go`、`registry/etcdv3/*_test.go`、`snowflake/snowflake_test.go`、`crontab/task_test.go`、`watch/*_test.go`、`img/img_test.go` 已通读，不含生产实现。
