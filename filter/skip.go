package filter

import (
	"slices"
	"sync"
)

const sortedSkipList = "\n\r!\"#$%&'()*+-:;=@[]^_{|}~¤§¨°±·×÷ˉΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψω—―‖‘’“”…‰※€℃№ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ←↑→↓∈∏∑√∝∞∠∥∧∨∩∪∫∮∴∵∶∷∽≈≌≠≡≤≥≮≯⊙⊥⌒①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛─━│┃┄┅┆┇┈┉┊┋┌┍┎┐┑┒┓└┕┖┗┘┙┚┛├┝┞┟┠┡┢┣┤┥┦┧┨┩┪┫┬┭┮┯┰┱┲┳┴┵┶┷┸┹┺┻┼┽┾┿╀╁╂╃╄╅╆╇╈╉╊╋■□▲△◆◇○◎●★☆♀♂、、、。。〃々〈〉《《》》「」『』【【】】〓〔〕〖〗㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩︿！＂＃＆＇（）＋，，－．／：；＜＝＞？？＠［＼］＿｀｛｜｝～￣"

// defaultSkipOnce / defaultSkipRunes：PERF-9-13 [P2] 包级 sync.Once 预解码默认跳过常量，
// 避免每次构造扫描器都重新解码约 500 字符的长常量并分配 rune 切片。
var (
	defaultSkipOnce   sync.Once
	defaultSkipRunes  []rune
	defaultSkipValues map[rune]struct{}
)

// defaultSkipSet 惰性初始化并返回内置默认跳过集合的 rune 切片与成员判定 map。
func defaultSkipSet() ([]rune, map[rune]struct{}) {
	defaultSkipOnce.Do(func() {
		defaultSkipRunes = []rune(sortedSkipList)
		defaultSkipValues = make(map[rune]struct{}, len(defaultSkipRunes))
		for _, r := range defaultSkipRunes {
			defaultSkipValues[r] = struct{}{}
		}
	})
	return defaultSkipRunes, defaultSkipValues
}

// SortedSkipList 返回内置默认的跳过分隔符字符集合字符串。
func SortedSkipList() string {
	return sortedSkipList
}

// Skip 维护一组在匹配时需要跳过的字符，同时保留有序 rune 切片供枚举，
// 并构建成员判定 map 供热路径 O(1) 查询（PERF-9-7 [P1]）。
//
// 零值可用：未设置字符集合时 ShouldSkip 始终返回 false；该类型不保证并发安全。
type Skip struct {
	list []rune
	set  map[rune]struct{}
}

// Set 用指定字符串中的字符重建跳过集合，内部排序后供枚举，并构建成员判定 map。
func (_this *Skip) Set(s string) {
	list := []rune(s)
	slices.Sort(list)
	_this.list = list
	_this.set = buildSkipSet(list)
}

// SetSorted 直接使用指定字符串中的字符作为跳过集合。
//
// 调用方须保证字符已按 rune 顺序排好；成员判定与顺序无关，map 直接由切片构建（PERF-9-7 [P1]）。
func (_this *Skip) SetSorted(s string) {
	list := []rune(s)
	_this.list = list
	_this.set = buildSkipSet(list)
}

// buildSkipSet 由 rune 切片构建成员判定 map；切片为空时返回 nil。
func buildSkipSet(list []rune) map[rune]struct{} {
	if len(list) == 0 {
		return nil
	}
	set := make(map[rune]struct{}, len(list))
	for _, r := range list {
		set[r] = struct{}{}
	}
	return set
}

// ShouldSkip 判断指定字符是否属于跳过集合；集合为空时返回 false。
//
// PERF-9-7 [P1]：由有序切片二分查找改为 map O(1) 成员判定，成员语义与原二分查找完全一致。
func (_this *Skip) ShouldSkip(r rune) bool {
	_, ok := _this.set[r]
	return ok
}

// String 返回跳过集合对应的字符串。
func (_this *Skip) String() string {
	return string(_this.list)
}
