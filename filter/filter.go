package filter

import (
	"sort"
)

type ResultWriter struct {
	count int
}

func (_this *ResultWriter) Write(_ *Result) (stop bool) {
	_this.count++
	return true
}

func (_this *ResultWriter) Len() int {
	return _this.count
}

func skipStr(skip ...string) []rune {
	if len(skip) > 0 {
		runes := []rune(skip[0])
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		return runes
	}
	return []rune(SortedSkipList())
}

// Strings 将输入的敏感词列表转换成 tireRoot 树
func Strings(words []string, skip ...string) *Search {
	search := NewSearch(SetSortedRunesSkip(skipStr(skip...)))
	search.TrieWriter().InsertWords(words).BuildFail()
	return search
}

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

type Option func(options *options)

func SetWriter(w *TrieWriter) Option {
	return func(options *options) {
		options.writer = w
	}
}

func SetSortedRunesSkip(s []rune) Option {
	return func(options *options) {
		skip := &Skip{list: s}
		options.skip = skip
	}
}

func SetSortedSkip(s string) Option {
	return func(options *options) {
		skip := &Skip{}
		skip.SetSorted(s)
		options.skip = skip
	}
}

func SetSkip(s string) Option {
	return func(options *options) {
		skip := &Skip{}
		skip.Set(s)
		options.skip = skip
	}
}

func NewSearch(opts ...Option) *Search {
	opt := &options{
		skip:   &Skip{list: []rune(sortedSkipList)},
		writer: NewTrieWriter(),
	}
	for _, o := range opts {
		o(opt)
	}
	opt.writer.setSkip(opt.skip)
	return &Search{opt.writer}
}
