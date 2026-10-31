# filter / correct / nlpword / lang 包组性能审查报告

> 审查范围（仅本章节）：`filter`（`doc.go`、`filter.go`、`result.go`、`search.go`、`skip.go`、`trie.go`）、`correct`（`doc.go`、`format.go`、`fullwidth.go`、`halfwidth.go`、`html.go`、`strategery.go`、`unformat.go`）、`nlpword`（`doc.go`、`ac.go`、`filter.go`、`link.go`、`trie.go`）、`lang`（`lang.go`）及其子目录 `lang/i18n`（`doc.go`、`i18n.go`）下全部 `.go` 文件；测试文件仅通读核对，性能条目只针对生产代码。
> 审查方式：逐文件人工通读源码，先借注释理解 API 语义，再定位热点路径与可疑实现；每个可疑点回到原文核对 `文件:行号` 与函数名；`go vet ./filter/... ./correct/... ./nlpword/... ./lang/...` 已通过，无告警。
> 优先级定义：**P0** = 热路径上的明显性能缺陷，应优先修改；**P1** = 中等影响、性价比高的改进；**P2** = 低影响或改进空间有限的点。

---

## 一、条目清单（按子包分组；组内 P0 → P1 → P2 排序）

### nlpword

#### [P1] `nlpword/ac.go:9` `(*ac).fail`：失配字符在查询热路径上沿失败指针逐层回溯，未预构建压缩转移表
- 问题：标准 Aho-Corasick 自动机在建链阶段把“状态 × 字符 → 后继状态”的转移表固化下来，查询每个字符 O(1)；本实现 `fail`（9–24 行）在失配时从 `node.Failure` 起逐层做 `Children[c]` 查找（12 行），找不到就继续 `node = node.Failure`（17 行），直到根节点。转移关系只保留了直接子节点，失败链上的转移目标没有被缓存。
- 影响评估：触发面为四个导出入口——`Trie.Replace`（`nlpword/trie.go:103`）、`Trie.Validate`（`nlpword/trie.go:159`）、`Trie.FindIn`（`nlpword/trie.go:184`）、`Trie.FindAll`（`nlpword/trie.go:190`），每个文本字符失配都可能进入 `fail`。随机文本通常回溯 1~2 跳；对抗性输入（长词前缀后紧跟失配字符）每次失配最坏回溯 L 跳，L = 词库最长词深度，整体查询最坏 O(n·L)，而压缩转移表的标准 AC 为 O(n)。
- 优化建议：在 `BuildFailureLinks`（`nlpword/trie.go:28`）阶段同步构建完整 goto 转移表（状态到后继的映射；可惰性填充并缓存），`next` / `fail` 合并为一次 O(1) 转移查询，查询期不再逐层上溯。

#### [P1] `nlpword/ac.go:36` `(*ac).output`（合并 `firstOutput:49`、`replace:62`）：每个字符位置都递归遍历整条失败链收集后缀命中，无 output-link 压缩，命中即做切片→字符串分配
- 问题：`output`（36–46 行）、`firstOutput`（49–59 行）、`replace`（62–73 行）在每个扫描位置都从当前节点沿 `Failure` 递归到底；即使当前节点不是结束节点，也会把整条失败链走一遍，没有预计算“最近的结束后缀节点”（output link / dictionary suffix link）。命中时 `output:42`、`firstOutput:55` 立即执行 `string(runes[...])` 堆分配；`Validate`（`nlpword/trie.go:175`）只需要判断“是否命中”，仍先构造完整命中字符串。此外 `FindAll` 的 `ac.results`（`nlpword/trie.go:197` 起步为 nil）在 `output` 内反复 append，无容量预估，命中密集时切片多次翻倍扩容拷贝。
- 影响评估：`FindAll` 对每个位置都付出递归与整条链遍历；词库深、文本命中密集时递归深度与临时字符串数量同步增长；`Validate` / `FindIn` 虽提前返回第一条命中，但命中字符串的分配在返回前已经发生。递归深度等于失败链长度，深词库下还存在栈深风险（正确性线索，不在本报告展开）。
- 优化建议：
  1. 建链时为每个节点预计算 output link（指向最近的结束后缀节点），收集命中时沿 output link 跳跃，跳过非结束节点。
  2. `firstOutput` 先返回命中节点与字符区间，调用方确有需要再切片转 `string`；`Validate` 路径只判断结束标记，不构造字符串。
  3. 递归改为迭代；`FindAll` 结果切片按文本长度 / 最短词长预留初始容量。

#### [P1] `nlpword/trie.go:125` `Trie.Filter`：完全未使用失败指针，失配后回退到 `left` 位置从头重扫，共享前缀被反复下钻
- 问题：`Filter` 循环内（136–151 行）只查 `parent.Children[runes[position]]`（137 行）；失配即把 `parent` 重置为 `tree.Root`（141 行）、`position = left`（142 行）回退重扫。`BuildFailureLinks` 构建出的失败指针在本函数完全不参与转移，本质是朴素 trie 回溯匹配，而不是 AC 自动机扫描。
- 影响评估：当文本与词库共享长前缀后失配（如关键词“abcdef…”、文本“abcdef…x”），已下钻的前缀在回退后要从根节点重复走一遍；最坏情况下每个字符段重复遍历深度 L，整体 O(n·L) 量级。该函数是 `Filter.Filter`（`nlpword/filter.go:106`）的底层实现，敏感词删除类调用都会经过。
- 优化建议：与 `Validate` / `FindAll` 统一走失败指针（配合本包第一条的压缩转移表）：失配时沿 fail 链转移到可匹配状态继续扫描，而不是回到 `left` 重新下钻；在保留当前删除语义的前提下消除前缀重复遍历。（注：本函数现有回退 / 命中区间维护方式是否完全等价于“删除全部命中词”属正确性范畴，需另组核对。）

#### [P2] `nlpword/trie.go:28` `BuildFailureLinks`（含 `bfs:47`、`nlpword/link.go:13` `Push`、`:28` `Pop`）：构建期 goroutine + channel 逐节点传递，队列用 `any` 装箱并为每个元素分配链表节点
- 问题：`bfs`（47–66 行）为单次遍历启动一个 goroutine，每个节点通过 `ch <- n`（57 行）做一次 channel 收发；队列 `LinkList`（`nlpword/link.go:6–10`）的 `Push`（13–25 行）为每个元素新建 `listNode`（14 行）并把 `*Node` 装箱成 `any`（15 行），`Pop`（28–37 行）返回 `any`，`nlpword/trie.go:56` 再做一次 `.(*Node)` 类型断言。`BuildFailureLinks` 建链时（30–40 行）从 `node.Parent` 逐层上溯查找失败转移，没有直接利用父节点已经建好的 `Failure`。
- 影响评估：仅发生在词库加载 / 变更后首次匹配触发的重建阶段（`Filter.updateFailureLink`，`nlpword/filter.go:92–97`），不在请求热路径；万级词库重建时，节点数 N 对应 N 次 channel 收发、N 个 `listNode` 堆对象与 N 次接口装箱，构建时间与临时对象数量随 N 线性放大。
- 优化建议：BFS 改用 `[]*Node` 切片队列（两个切片互换复用），去掉 goroutine 与 channel；队列元素直接用 `*Node` 具体类型，去掉 `any` 装箱与类型断言；建链按标准 AC：`fail[child] = transition(fail[parent], c)`，利用父节点失败指针一次定位，避免逐层上溯。

#### [P2] `nlpword/trie.go:213` `NewNode`（含 `NewRootNode:221`）：每个节点无条件分配空 `map[rune]*Node`，叶子节点同样支付一次 map 堆分配
- 问题：`NewNode` 216 行、`NewRootNode` 225 行创建即 `make(map[rune]*Node, 0)`；叶节点之后永远不会写入子节点，空 map 只占内存。
- 影响评估：词库节点数为各词字符并集，通常数万；每个空 map 包含 map header 与底层 bucket 的堆分配，且会被 GC 扫描；词库越大，空 map 的总内存与扫描开销越明显。构建期一次性成本，不影响查询吞吐。
- 优化建议：`Children` 初始保持 `nil`，首次插入子节点时再 `make`；Go 中对 nil map 做读取（`m[k]`）安全返回零值，`add`（`nlpword/trie.go:87`）现有读取逻辑无需改动。

#### [P2] `nlpword/filter.go:120` `Filter.FindIn`（含 `Validate:135`、`RemoveNoise:142`）：检测前对整段文本先跑一遍噪声正则并生成清洗后新串
- 问题：`FindIn`（122 行）与 `Validate`（137 行）在交给 trie 之前都先 `RemoveNoise`；`RemoveNoise`（142–144 行）执行 `noise.ReplaceAllString(text, "")`，正则 `[\|\s&%$@*]+`（29 行）对整段文本做一次全量匹配扫描。文本不含噪声字符时，Go 正则虽返回原串、不分配结果串，但扫描过程本身仍然走完一遍；含噪声时再额外分配一个清洗后字符串供后续匹配。
- 影响评估：每次敏感词检测多一次 O(n) 正则扫描（必要时再加一次整串分配），叠加在 trie 匹配之前；检测类接口若被高频调用（评论 / 发帖校验），这部分是词表匹配之外的固定前置开销。`FindAll`（127 行）、`Filter`（106 行）、`Replace`（112 行）当前不做噪声剥离，行为不一致。
- 优化建议：先用极轻量的线性预检（针对少量噪声字符逐字节 / rune 判断是否存在），不含噪声时直接跳过正则；或把噪声字符跳过逻辑并入 trie 状态机在匹配过程中跳过，避免整串复制。

### filter

#### [P1] `filter/skip.go:37` `Skip.ShouldSkip`：扫描热路径上每个字符都走有序切片二分查找，应为 O(1) 成员判定
- 问题：`ShouldSkip`（37–56 行）先做首尾边界比较（42 行），再在 45–54 行做二分查找；默认跳过集合 `sortedSkipList`（`filter/skip.go:7`）约 500 个字符。`Search.findByAC` 在扫描主循环内对每个字符调用一次（`filter/search.go:141`）；词库构建期 `Insert`（`filter/trie.go:62`）、`InsertBytes`（`filter/trie.go:163`）同样逐字符调用。
- 影响评估：扫描热路径上每字符约 `log2(500) ≈ 9` 次切片下标比较 + 一次函数调用；中文长文本高频检测时，这是与 trie 转移并列的每字符固定成本。二分查找适合需要有序枚举的场景，成员判定用 map 更合适。
- 优化建议：`Set`（22–26 行）/`SetSorted`（31–34 行）在维护有序 rune 切片的同时构建 `map[rune]struct{}`，`ShouldSkip` 改为 map 命中判定 O(1)；有序切片仅在 `String`（59–61 行）等需要枚举的接口保留。

#### [P1] `filter/search.go:97` `findSub`（调用点 `:134`、`:155`）：每条命中都从结尾反向字节扫描定位起始位置，正向扫描本可记录偏移
- 问题：命中分支 132–135 行与失败后缀命中分支 152–156 行都调用 `findSub`；`findSub`（97–106 行）从 `end` 起逐字节反向与 `sub` 比对直到匹配。`node.len` / `node.fail.len` 已经记录词库词条的字节长度（`filter/trie.go:70`、`:173`），但因为被跳过字符夹在原文中间，当前实现用反扫来补齐真实起始位置。
- 影响评估：每条命中付出 O(命中长度) 的反向字节比较；命中密集的长文本上，反扫与正向 trie 扫描叠加；两个分支在反扫前后还各产生 `string(word[...])`（132、153 行）与 `string(s[start:end+1])`（135、156 行）两次字符串分配。
- 优化建议：正向扫描进入当前片段 / 首次下钻节点时记录起始字节偏移 `j`，命中时直接使用记录的 `start`，配合已累计的 `skip` 偏移量校正跨距，去掉反向扫描；`Result` 所需的两条字符串保留（返回值必需），但定位过程不再重复比较。

#### [P1] `filter/search.go:29` `Search.HasSens`：只需要布尔存在性，single 模式下命中前仍构造 `Result` 与两条字符串
- 问题：`HasSens`（29–31 行）调用 `findByAC(s, true)` 后只取 `len(list) > 0`；但 `findByAC` 即使 `single=true`，在 135 / 156 行返回前已经构造 `*Result` 并完成 `string(...)` 两次转换，159 行 append 进结果切片，直到 161–163 行才返回。
- 影响评估：存在性判定（风控 / 审核前置校验）是最轻量的调用形态，却为第一次命中支付一次结构体分配 + 两条字符串堆分配；无命中路径不受影响。调用越频繁，这部分无用分配的 GC 压力越明显。
- 优化建议：`findByAC` 内按 `single` 分支：命中即置标志并立即返回，不构造 `Result`、不做字符串转换、不 append；`HasSens` 直接返回布尔。

#### [P2] `filter/search.go:34` `Search.Replace`（含 `ReplaceRune:49`、`repeatBytes:79`、`repeatByte:88`）：每条命中先分配重复填充的临时切片再整块拷贝进输出，临时切片随即丢弃
- 问题：`Replace` 40 行调用 `repeatByte`（88–93 行）先 `make([]byte, n)` 填充，再 `copy` 进 `data`；`ReplaceRune` 70 行调用 `repeatBytes`（79–85 行）同样先 `make([]byte, al)` 填充再拷贝；`ReplaceRune` 64 行还做了一次 `[]byte(string(new))` 的字符串→字节转换。
- 影响评估：每条命中多一次临时切片分配 + 一次整块拷贝，临时对象随命中密度进入 GC；单次替换长度小，属中等偏低开销，命中密集的批量替换场景累积可见。
- 优化建议：直接在 `data` 的目标区间填充：单字节替换用循环赋值；rune 替换先把 `new` 编码到固定长度小数组，按区间重复写入；去掉中间临时切片。

#### [P2] `filter/trie.go:25` `NewTrieWriter`（含 `Insert:58`、`InsertBytes:154`）：每个 trie 节点无条件分配空 map，叶子节点同样分配
- 问题：根节点 26 行、新建节点 66 行与 169 行均为 `map[rune]*trie{}` 创建即分配；叶节点后续不再写入子节点。
- 影响评估：与 `nlpword` 同类问题，构建期内存占用随节点数线性增长；空 map 堆对象与 GC 扫描量在大词库上累积。
- 优化建议：子节点 map 惰性创建（初始 nil，首次插入时 `make`），nil map 读取安全。

#### [P2] `filter/trie.go:192` `BuildFail`（含 `String:224`、`Array:262`、`appendRunes:254`）：BFS 每层都 make 临时切片并整体拷贝；前缀 rune 在每个节点上整段复制
- 问题：`BuildFail` 201–203 行每层 `temp := make([]*trie, len(queue)); copy(temp, queue); queue = queue[:0]`，可在两个切片间互换复用；`String` 234–236 行、`Array` 270–272 行每层同样 make + copy；`appendRunes`（254–259 行）为每个节点 `make([]rune, len+1)` 并 `copy` 整段前缀，前缀随深度重复复制。
- 影响评估：`BuildFail` 为一次性构建；`String` / `Array` 是调试输出与词库导出工具，不在请求热路径；大词库遍历导出时，节点数 × 前缀深度的拷贝量明显。
- 优化建议：BFS 用两个队列切片互换，消除每层分配；前缀路径改为 DFS + 栈式 rune 切片（push/pop 复用同一底层数组），避免每节点复制整段前缀。

#### [P2] `filter/filter.go:24` `skipStr`（`NewSearch:95` 默认分支 97 行）：默认跳过字符常量每次构造扫描器都重新解码为 rune 切片
- 问题：`skipStr` 30 行 `[]rune(SortedSkipList())`；`NewSearch` 未传自定义 skip 时 97 行 `[]rune(sortedSkipList)`，每次都把约 500 字符的长常量重新 UTF-8 解码并分配 rune 切片。
- 影响评估：扫描器通常启动期创建一次，影响有限；若运行期按多词库反复 `NewSearch` / `Strings`（36 行）/ `File`（45 行），重复解码与分配会累积。
- 优化建议：包级 `sync.Once` 预解码默认 rune 切片，构造时直接引用；自定义 skip 分支保持现有排序逻辑。

### correct

#### [P1] `correct/format.go:55` `Format`：一条文本串行执行约 22 趟全文正则扫描，每趟都分配一个中间字符串
- 问题：`Format` 流水线为：`halfwidth`（58 行，内含 `halfTimeRe` 1 趟，`correct/halfwidth.go:32`）→ `fullwidth`（59 行，`correct/fullwidth.go:34–37` 共 4 趟 `ReplaceAllStringFunc`）→ 遍历 9 条 `strategery`（61–63 行；`correct/strategery.go:47`、`:49`、`:57`、`:59`，合计 14 趟 `ReplaceAllString`）→ `spaceDashWithHans`（65 行，`correct/format.go:41–43` 共 3 趟）。合计约 22 趟全文正则扫描，每趟 `ReplaceAllString` / `ReplaceAllStringFunc` 都产出一个新字符串，上一趟结果作为下一趟输入，中间串随即丢弃。
- 影响评估：若 `Format` 进入内容发布、评论、消息体等逐条文本的排版流水线，单条文本固定支付约 22 次正则 DFA 全文扫描 + 约 22 次字符串堆分配；文本越长，扫描与拷贝成本线性叠加。仓内当前无生产调用方（见第四节待实测清单），量级需在真实文本样本上 benchmark 确认。
- 优化建议：
  1. 合并互补的补空格 / 去空格模式为少量组合正则，减少全文扫描趟数；
  2. 或实现单遍 rune 状态机，一次扫描完成 CJK 与英文 / 数字 / 符号之间的补 / 去空格，避免多趟中间串；
  3. 正则处理尽量在 `[]byte` 缓冲上复用，减少 string↔[]byte 与中间字符串分配。

#### [P2] `correct/halfwidth.go:19` `halfwidth`：rune 切片零值起步反复扩容，同一字符重复查两次 map
- 问题：20 行 `runes := []rune{}` 无初始容量，append 按翻倍策略多次扩容拷贝；22 行查出 `newChar`，24 行又执行一次 `charWidthMap[char]` 重复查找。
- 影响评估：被 `Format` 每条文本调用，切片扩容 2~3 次与重复 map 查找属小开销，叠加在多趟排版流程内；独立调用时同样存在。
- 优化建议：`make([]rune, 0, utf8.RuneCountInString(text))` 预留容量；直接复用 22 行查到的 `newChar`。

#### [P2] `correct/html.go:36` `processHTML`：文本节点经历 []byte→string→[]byte 两次复制；每个开始标签执行一次正则匹配；输出缓冲未预分配
- 问题：64 行 `fn(string(data))` 把词法分析得到的 `data` 拷贝成 string 交给排版函数；65 行 `w.Write([]byte(formated))` 再把结果 string 转回 []byte 写入缓冲，形成两次复制。69 行每个 `StartTagToken` 都执行 `ignoreTagsRe.Match(data)` 正则匹配，而需要识别的只有 pre / script / style / textarea 四类标签。37 行 `bytes.Buffer{}` 零起步，随输出动态扩容。
- 影响评估：HTML 文本节点越多、文档越长，双拷贝与标签正则次数越多；标签名均为 ASCII，正则 DFA 启动开销相对字节比较偏高。`FormatHTML`（20 行）/ `UnformatHTML`（28 行）共用本函数。
- 优化建议：排版回调签名改为接受 `[]byte`（或内部直接 `w.WriteString(formated)` 去掉第二次转换）；四类标签用字节前缀比较或小型 map 判定，去掉每标签正则；缓冲按 `len(body)` 预分配容量。

### lang/i18n

#### [P1] `lang/i18n/i18n.go:145` `I18n.Tr`：每次翻译读取都经第三方 `ini.GetValue` 获取排他写锁，并发翻译在该 Mutex 上串行化
- 问题：`Tr` 在 145 行调用 `l.data[lang].GetValue(key)`，回退路径 177 行 `l.data[fl].String(key)` 最终也进入同一个 `GetValue`。第三方 `gookit/ini/v2 v2.3.2` 的 `GetValue`（`manage.go:27–31`）在 `opts.Readonly == false`（默认值，`options.go:37–46` 未设置 Readonly）时执行 `c.lock.Lock()` 排他写锁、`defer Unlock()`；本包所有实例都通过 `ini.New()` 构造（`lang/i18n/i18n.go:269`、`:301–302`、`:324`），未启用 `Readonly`。也就是说只读查询走的是写锁，并发读取之间也不能共享。
- 影响评估：`T` / `Tr` / `Dt` / `DefTr`（59、64、69、74 行及实例方法 114、119、124、132）是每条本地化消息的入口；高并发下不同 goroutine 的翻译查询会在各语言 `Ini` 实例的同一把 `sync.RWMutex` 写锁上排队，读请求互相串行。锁临界区内还包含 `formatKey`（TrimSpace / Trim）、`splitSectionAndKey`（Contains / SplitN）等字符串处理。实际竞争程度取决于翻译 QPS，需 benchmark 验证（见第四节）。
- 优化建议：
  1. 词库加载完成前用 `ini.NewWithOptions(ini.Readonly)` 构造实例（注意第三方 `WithOptions` 在数据加载后调用会 panic，`ini.go:140–143`，必须在加载前设置），让 `GetValue` 跳过加锁；
  2. 或在本包内把已加载的翻译表扁平化为 `map[string]string`（按 `lang + section + key` 组合键），查询直接读 map，绕开第三方锁；
  3. 若保留 `Ini`，可在 `Tr` 层加一层只读缓存。

#### [P1] `lang/i18n/i18n.go:182` `renderMessage`：ReplaceMode 每次翻译都现场构建 `strings.NewReplacer` 替换器并编译其内部 trie
- 问题：`renderMessage` 在 `TransMode == ReplaceMode` 时，`[]string` 分支 188 行、map / 交替参数分支 210 行，每次调用都 `strings.NewReplacer(ss...)`；`NewReplacer` 内部会对替换键排序并构建匹配 trie，属一次性编译成本，却在每条消息渲染时重复支付。map 分支 194–197 行还先遍历 map 拼接 `{key}` 与值，交替参数分支 200–207 行逐参数 `toString`。
- 影响评估：所有带参数、使用 ReplaceMode 的本地化消息都经过这里；提示语 / 错误文案在请求路径高频渲染时，重复编译替换器的开销叠加在每条消息上。SprintfMode 分支 184 行 `fmt.Sprintf` 为另一模式，不在本条目范围。
- 优化建议：替换键数量通常很少，改用直接的 `strings.ReplaceAll` 循环，或自行实现单遍 `{key}` 扫描替换，避免每次构建 trie；若同一模板 `msg` 被反复渲染，可按模板缓存解析结果（替换值仍动态传入）。

#### [P2] `lang/i18n/i18n.go:132` `I18n.Tr`：先查 `languages` map 再查 `data` map，每次翻译两次哈希查找；回退路径 `transFromFallback:171` 同样双重查找
- 问题：133 行 `HasLang(lang)`（383–385 行查 `languages`），145 行再 `l.data[lang]` 查第二个 map；147 行进入 `transFromFallback` 后，173 行再次 `HasLang(fl)`、177 行再 `l.data[fl]`。`HasKey`（161–168 行）同样是两次查找。
- 影响评估：每次翻译两次 map 哈希，属轻量固定开销；在高频翻译调用上累积，量级低于上一条锁竞争。
- 优化建议：直接取 `ld := l.data[lang]`，以非 nil 判定语言是否已注册，合并两次查找；回退路径与 `HasKey` 同步简化。

#### [P2] `lang/i18n/i18n.go:319` `NewLang`（含 `LoadFile:329`、`LoadString:344`、`DelLang:389`）：`data` / `languages` map 在翻译读路径与动态加载写路径之间没有任何同步措施
- 问题：`Tr`（132）、`HasLang`（383）等读路径直接访问两个 map；`NewLang`（320–325 行）、`LoadFile`（330–331 行）、`LoadString`（345–346 行）、`DelLang`（392–393 行）直接写 map；整个 `I18n` 类型没有 `sync.Mutex` / `RWMutex` / `sync.Map`。
- 影响评估：严格按“启动期 `Init` 加载完成后只读”的用法，稳态并发读取是安全的，不构成性能热点；若运行期动态调用 `NewLang` / `LoadFile` / `LoadString` / `DelLang` 与 `Tr` 并发，存在 map 并发读写导致 panic 的风险（并发安全 / 正确性线索，本条目仅登记，不展开正确性审查）。
- 优化建议：文档明确“仅启动期加载、运行期只读”；如需运行期热更新，为 `data` / `languages` 增加 `RWMutex`（读路径 RLock、写路径 Lock），或改用 `sync.Map`。

### lang

`lang/lang.go` 仅类型别名与空结构体占位，无可执行逻辑，无性能条目（见第六节覆盖说明）。

---

## 二、条目统计

| 子包 | P0 | P1 | P2 | 小计 |
|---|---|---|---|---|
| nlpword | 0 | 3（`ac.fail`、`ac.output` 系列、`Trie.Filter`） | 3（`BuildFailureLinks`/`bfs`/`LinkList`、`NewNode` 空 map、`FindIn/Validate` 噪声正则） | 6 |
| filter | 0 | 3（`ShouldSkip`、`findSub`、`HasSens`） | 4（`Replace`/`ReplaceRune` 临时切片、节点空 map、`BuildFail`/`String`/`Array` 队列与前缀拷贝、`skipStr` 默认解码） | 7 |
| correct | 0 | 1（`Format` 多趟正则） | 2（`halfwidth` 扩容与重复查找、`processHTML` 双拷贝与标签正则） | 3 |
| lang/i18n | 0 | 2（`Tr` 第三方写锁、`renderMessage` 重复编译 Replacer） | 2（`Tr` 双重 map 查找、`NewLang` 系列并发同步缺失） | 4 |
| lang | 0 | 0 | 0 | 0 |
| **合计** | **0** | **9** | **11** | **20** |

---

## 三、高价值候选列表（Top 10，供全局 Top 10 评选）

| # | 位置 | 一句话问题 | 预期收益 |
|---|---|---|---|
| 1 | `lang/i18n/i18n.go:145` `I18n.Tr` | 只读翻译查询经第三方 `ini.GetValue` 获取排他写锁，并发翻译互相串行 | 启用 Readonly 或扁平化缓存后，消除每请求一把写锁的排队竞争；高并发消息翻译吞吐直接受益，改动集中在实例构造与查询层 |
| 2 | `filter/skip.go:37` `Skip.ShouldSkip` | 扫描热路径每个字符对约 500 元素有序切片做二分查找 | 改为 map 成员判定后，每字符从约 9 次比较降到 O(1) 哈希；敏感词扫描主循环固定开销明显下降 |
| 3 | `correct/format.go:55` `Format` | 单条文本串行约 22 趟全文正则扫描并分配 22 个中间字符串 | 合并正则或单遍状态机可把全文扫描趟数降到个位数，同时砍掉大部分中间串分配；逐条文本排版流水线的 CPU 与 GC 压力同步下降 |
| 4 | `filter/search.go:29` `Search.HasSens` | 仅需布尔结果，single 模式仍先构造 `Result` 和两条字符串 | 命中即短路返回，存在性判定路径消除结构体与字符串分配；轻量风控校验场景 GC 压力下降，改动小 |
| 5 | `filter/search.go:97` `findSub` | 每条命中反向扫描定位起始字节，正向扫描已有位置信息却未记录 | 记录片段起始偏移后，命中处理去掉 O(命中长度) 的反向比较；命中密集文本上减少重复字节比较 |
| 6 | `lang/i18n/i18n.go:182` `renderMessage` | ReplaceMode 每条消息都现场构建 `strings.NewReplacer` 并编译内部 trie | 换成少量键的直接替换 / 单遍扫描后，去掉每次调用的 trie 编译；带参数本地化消息的 CPU 与临时对象减少 |
| 7 | `nlpword/ac.go:9` `(*ac).fail` | 失配时在查询路径逐层回溯失败链，缺少压缩转移表 | 建链阶段固化转移表后，查询每字符 O(1)；对抗性长前缀文本从 O(n·L) 回落到标准 AC 的 O(n) |
| 8 | `nlpword/ac.go:36` `(*ac).output` 系列 | 每位置递归遍历整条失败链收集后缀命中，无 output-link | 预计算 output link 后跳过非结束节点；命中收集递归深度与递归调用次数下降，`Validate` 路径可把构造信息最小化 |
| 9 | `nlpword/trie.go:125` `Trie.Filter` | 未使用失败指针，失配回退 `left` 从头重扫 trie | 接入失败指针转移后消除共享前缀重复下钻，长前缀对抗文本最坏复杂度从 O(n·L) 降到 O(n)（需先核对删除语义正确性） |
| 10 | `nlpword/trie.go:213` `NewNode`（含 `filter` 包同类空 map 分配） | 每个节点无条件分配空 map，叶子节点空占一份堆对象 | map 惰性创建后，大词库节点数对应的空 map 分配与 GC 扫描量大幅减少；改动机械、风险低 |

---

## 四、拿不准、需实测验证的条目清单

1. **`lang/i18n/i18n.go:145` `Tr` 第三方写锁的实际竞争程度**：仓内当前无生产调用方，锁排队程度只能在接入 i18n 的服务上用 mutex profile / 并发 benchmark 实测；若翻译 QPS 低、各语言实例独立，竞争可能并不显著。
2. **`correct/format.go:55` `Format` 22 趟正则的实际占比**：仓内当前无生产调用方；需用真实中英文混排样本（不同长度、不同 CJK / 英文比例）做 benchmark，量化正则扫描与中间串分配在总耗时 / 分配中的占比，再决定合并正则或重写状态机的投入力度。
3. **`nlpword/ac.go:9` `fail` 逐层回溯的跳数分布**：随机文本通常 1~2 跳，回溯开销有限；只有在长前缀 / 对抗性文本上才接近最坏 O(n·L)，需构造代表性词库与文本样本测量实际跳数分布。
4. **`nlpword/trie.go:125` `Trie.Filter` 回退重扫的真实复杂度**：该函数同时存在命中区间维护方式是否等价于“删除全部命中词”的正确性疑问；性能收益结论需要在正确性核对通过后，再用 benchmark 量化。
5. **`filter/search.go:97` `findSub` 反扫的占比**：跳过字符夹在原文中间是反扫存在的原因，改造为正向记录偏移前需用样本确认跳过字符的出现频率与命中长度分布。
6. **`nlpword/filter.go:120` 噪声正则前置扫描的必要性**：`FindIn` / `Validate` 剥离噪声而 `FindAll` / `Filter` / `Replace` 不剥离，行为差异是否为有意设计需产品语义确认；预检跳过正则能省多少取决于文本含噪声字符的比例。
7. **`lang/i18n/i18n.go:132` 双重 map 查找收益**：两次哈希查找的固定开销很小，是否值得合并取决于翻译调用频次；建议随第 1 条 benchmark 一并观察。
8. **`correct/html.go:36` `processHTML` 双拷贝与标签正则**：实际开销取决于文档中文本节点数量与标签密度；需在典型 HTML 样本上对比改造前后的 allocs/op 与 ns/op。
9. **第三方 `ini.GetValue` 内部 `formatKey` / `splitSectionAndKey` 字符串处理占比**：已确认默认 `ParseVar=false`（不进入变量引用解析），但 TrimSpace / Contains / SplitN 的单次成本未量化；纳入第 1 条 benchmark 一并观察。

---

## 五、跨模块提示（超出本章节范围，仅登记线索，不展开审查）

1. **本包组在仓内暂无生产调用方**：全仓检索未发现除 `filter` / `correct` / `nlpword` / `lang` 自身之外的非测试代码 import `github.com/abulo/ratel/v3/{filter,correct,nlpword,lang,...}`；上述热路径影响判断基于 API 定位（敏感词检测 / 文本排版 / 消息翻译），接入方在服务中的调用位置与频次需由对应业务模块核实。
2. **`lang/i18n` 依赖的 `github.com/gookit/ini/v2 v2.3.2`**：`GetValue`（`manage.go:27–31`）在非 Readonly 模式下对只读查询使用排他 `Lock()`；`WithOptions`（`ini.go:140–143`）在数据加载后设置选项会 panic。本报告第一条 P1 的修复必须在加载前构造实例，该第三方行为供全局选型 / 替换决策参考。
3. **`correct/html.go` 依赖 `github.com/tdewolff/parse/v2` HTML 词法分析器**：词法流本身的性能不在本章节审查范围；若 `FormatHTML` 进入高频 HTML 处理路径，建议另对词法遍历与 `bytes.Buffer` 增长做专项 profile。
4. **`nlpword/filter.go:56` `LoadNetWordDict`**：每次调用新建 `http.Client`（57–59 行，5 秒超时），未复用连接池；属词库离线加载路径，不在翻译 / 检测热路径，仅登记备查，供加载器模块统一评估 HTTP 客户端复用。
5. **`nlpword/filter.go:92` `updateFailureLink` 版本号字段**（`buildVer` / `updatedVer`，21–22 行）无并发保护，与 `lang/i18n` 同属“文档声明非并发安全”的并发一致性线索，转正确性审查组参考，本报告不作为性能条目。

---

## 六、已通读但未发现性能问题的包 / 导出 API（覆盖说明）

- **`lang/lang.go`**：仅 `PlaceholderType`（空结构体）、`AnyType = any` 类型别名与包级空实例 `Placeholder`，无可执行逻辑与分配，无性能条目。
- **`filter/result.go`**：`Result` 为纯数据结构；`String()`（20–22 行）用 `fmt.Sprintf` 拼接调试描述，仅调试 / 日志场景调用，不在匹配热路径，不单独列条目。`ResultWriter.Write`（13–16 行）仅做计数器自增，开销可忽略。
- **`filter/filter.go`**：`Strings`（36–40）、`File`（45–51）、各 `Set*` Option（62–92）为启动期扫描器构造；`NewSearch`（95–105）默认分支的重复解码已列 P2，其余 Option 闭包调用为一次性配置。
- **`filter/trie.go` `InsertScanner`（99–103）、`InsertFile`（106–115）、`InsertReader`（118–120）、`insertReader`（122–151）、`InsertBytes`（154–188）**：均为词库加载期逻辑；64KB 读取缓冲复用合理，`InsertBytes` 逐 rune 解码为构建期成本；除空 map 分配（P2）外未发现热路径问题。`WriteString`（88–91）为 `io.Writer` 适配，转发 `Insert`。
- **`filter/search.go` `Find`（24–26）、`decodeBytes`（19–21）**：`Find` 仅转发 `findByAC`；`decodeBytes` 为 `utf8.DecodeRune` 直包装，编译器可内联，不单独列条目。
- **`nlpword/filter.go` `LoadWordDict`（39–52）、`Load`（75–89）、`AddWord`（100–103）、`updateFailureLink`（92–97）、`FindAll`（127–130）、`Filter`（106–109）、`Replace`（112–115）**：加载期逻辑；`updateFailureLink` 版本比对避免重复建链，写法正确；`AddWord` 中 `time.Now().UnixNano()`（102 行）为词库变更期调用，不在匹配热路径。`Load` 使用 `bufio.ReadLine` 逐行读取，属加载期开销。
- **`nlpword/trie.go` `Add`（76–80）、`add`（82–100）、`NewTrie`（69–73）、`FindIn`（184–187）、`Node` 各判定方法（235–247）**：`add` 内 `position == len(runes)-1`（96 行）为循环内不变量判断，开销微小，未单列；`FindIn` 仅转发 `Validate`；判定方法为内联布尔返回。
- **`nlpword/ac.go` `next`（27–33）**：单次 map 查找，无额外开销。
- **`correct/unformat.go`**：`Unformat`（16–24 行）仅一趟 `removeSpaceRe.ReplaceAllString`（17 行）加可选回调，无多趟叠加，未发现问题。
- **`correct/strategery.go`**：`newStrategery`（17–32 行）在 `init` 期编译 4 个正则，属一次性启动成本；`format` / `addSpace` / `removeSpace`（34–63 行）的多趟扫描已并入 `Format` P1 条目，不重复列项。
- **`correct/fullwidth.go`**：四趟扫描与段内嵌套正则已并入 `Format` P1 条目；`fullwidthMaps`（9–17 行）为静态 map，查询无额外构建成本。
- **测试文件**：`filter/*_test.go`、`correct/*_test.go`、`nlpword/*_test.go`、`lang/i18n/i18n_test.go` 已通读，均为空桩测试或使用样例，不含生产实现；本报告条目全部针对生产代码。
