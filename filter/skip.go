package filter

import (
	"slices"
)

const sortedSkipList = "\n\r!\"#$%&'()*+-:;=@[]^_{|}~¤§¨°±·×÷ˉΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψω—―‖‘’“”…‰※€℃№ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ←↑→↓∈∏∑√∝∞∠∥∧∨∩∪∫∮∴∵∶∷∽≈≌≠≡≤≥≮≯⊙⊥⌒①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛─━│┃┄┅┆┇┈┉┊┋┌┍┎┐┑┒┓└┕┖┗┘┙┚┛├┝┞┟┠┡┢┣┤┥┦┧┨┩┪┫┬┭┮┯┰┱┲┳┴┵┶┷┸┹┺┻┼┽┾┿╀╁╂╃╄╅╆╇╈╉╊╋■□▲△◆◇○◎●★☆♀♂、、、。。〃々〈〉《《》》「」『』【【】】〓〔〕〖〗㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩︿！＂＃＆＇（）＋，，－．／：；＜＝＞？？＠［＼］＿｀｛｜｝～￣"

// SortedSkipList 返回内置默认的跳过分隔符字符集合字符串。
func SortedSkipList() string {
	return sortedSkipList
}

// Skip 维护一组在匹配时需要跳过的字符，基于有序字符集合进行二分查找。
//
// 零值可用：未设置字符集合时 ShouldSkip 始终返回 false；该类型不保证并发安全。
type Skip struct {
	list []rune
}

// Set 用指定字符串中的字符重建跳过集合，并在内部排序后供二分查找使用。
func (_this *Skip) Set(s string) {
	list := []rune(s)
	slices.Sort(list)
	_this.list = list
}

// SetSorted 直接使用指定字符串中的字符作为跳过集合。
//
// 调用方须保证字符已按 rune 顺序排好，否则 ShouldSkip 的二分查找结果不可靠。
func (_this *Skip) SetSorted(s string) {
	list := []rune(s)
	_this.list = list
}

// ShouldSkip 判断指定字符是否属于跳过集合；集合为空时返回 false。
func (_this *Skip) ShouldSkip(r rune) bool {
	left, right := 0, len(_this.list)
	if right == 0 {
		return false
	}
	if r < _this.list[0] || r > _this.list[right-1] {
		return false
	}
	for left < right {
		mid := (left + right) >> 1
		if _this.list[mid] == r {
			return true
		} else if _this.list[mid] > r {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return false
}

// String 返回跳过集合对应的字符串。
func (_this *Skip) String() string {
	return string(_this.list)
}
