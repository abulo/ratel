package filter

import (
	"slices"
)

// ResultWriter 是仅统计命中数量的结果收集器，收集过程始终返回 true 表示不中断扫描。
type ResultWriter struct {
	count int
}

// Write 记录一次命中并返回 true 以继续扫描后续内容。
func (_this *ResultWriter) Write(_ *Result) (stop bool) {
	_this.count++
	return true
}

// Len 返回已记录的命中数量。
func (_this *ResultWriter) Len() int {
	return _this.count
}

// skipStr 返回跳过集合的 rune 切片；传入自定义字符串时先排序，未传入时使用内置默认集合。
//
// PERF-9-13 [P2]：默认集合通过 sync.Once 预解码缓存，构造时直接引用共享切片，避免重复解码长常量。
func skipStr(skip ...string) []rune {
	if len(skip) > 0 {
		runes := []rune(skip[0])
		slices.Sort(runes)
		return runes
	}
	runes, _ := defaultSkipSet()
	return runes
}

// Strings 用词表 words 构建扫描器并完成失败指针构建。
//
// 可选 skip 参数传入自定义跳过分隔符字符串，不传时使用内置默认跳过集合。
func Strings(words []string, skip ...string) *Search {
	search := NewSearch(SetSortedRunesSkip(skipStr(skip...)))
	search.TrieWriter().InsertWords(words).BuildFail()
	return search
}

// File 从字典文件加载词库构建扫描器并完成失败指针构建。
//
// 可选 skip 参数传入自定义跳过分隔符字符串；文件打开失败时 InsertFile 会直接 panic。
func File(filename string, skip ...string) (search *Search, err error) {
	search = NewSearch(SetSortedRunesSkip(skipStr(skip...)))
	writer := search.TrieWriter()
	writer.InsertFile(filename)
	writer.BuildFail()
	return &Search{writer}, nil
}

type options struct {
	writer *TrieWriter
	skip   *Skip
}

// Option 用于在 NewSearch 中以函数式方式配置词库写入器或跳过集合。
type Option func(options *options)

// SetWriter 指定自定义 TrieWriter 作为词库载体。
func SetWriter(w *TrieWriter) Option {
	return func(options *options) {
		options.writer = w
	}
}

// SetSortedRunesSkip 直接使用已排序的 rune 切片作为跳过集合；调用方须保证切片有序。
func SetSortedRunesSkip(s []rune) Option {
	return func(options *options) {
		skip := &Skip{list: s}
		// PERF-9-7 [P1]：同步构建成员判定 map，与二分查找成员语义一致。
		skip.set = buildSkipSet(s)
		options.skip = skip
	}
}

// SetSortedSkip 使用已排序字符串作为跳过集合；调用方须保证字符有序。
func SetSortedSkip(s string) Option {
	return func(options *options) {
		skip := &Skip{}
		skip.SetSorted(s)
		options.skip = skip
	}
}

// SetSkip 使用指定字符串重建跳过集合，内部自动排序。
func SetSkip(s string) Option {
	return func(options *options) {
		skip := &Skip{}
		skip.Set(s)
		options.skip = skip
	}
}

// NewSearch 按传入 Option 配置创建扫描器；未指定时使用内置默认跳过集合与空词库。
func NewSearch(opts ...Option) *Search {
	opt := &options{
		skip:   defaultSkip(),
		writer: NewTrieWriter(),
	}
	for _, o := range opts {
		o(opt)
	}
	opt.writer.setSkip(opt.skip)
	return &Search{opt.writer}
}

// defaultSkip 返回包级缓存的内置默认跳过集合（PERF-9-13 [P2]）。
func defaultSkip() *Skip {
	list, set := defaultSkipSet()
	return &Skip{list: list, set: set}
}
