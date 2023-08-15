package filter

import (
	"sort"
)

const sortedSkipList = "\n\r!\"#$%&'()*+-:;=@[]^_{|}~¤§¨°±·×÷ˉΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψω—―‖‘’“”…‰※€℃№ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ←↑→↓∈∏∑√∝∞∠∥∧∨∩∪∫∮∴∵∶∷∽≈≌≠≡≤≥≮≯⊙⊥⌒①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛─━│┃┄┅┆┇┈┉┊┋┌┍┎┐┑┒┓└┕┖┗┘┙┚┛├┝┞┟┠┡┢┣┤┥┦┧┨┩┪┫┬┭┮┯┰┱┲┳┴┵┶┷┸┹┺┻┼┽┾┿╀╁╂╃╄╅╆╇╈╉╊╋■□▲△◆◇○◎●★☆♀♂、、、。。〃々〈〉《《》》「」『』【【】】〓〔〕〖〗㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩︿！＂＃＆＇（）＋，，－．／：；＜＝＞？？＠［＼］＿｀｛｜｝～￣"

func SortedSkipList() string {
	return sortedSkipList
}

type Skip struct {
	list []rune
}

func (_this *Skip) Set(s string) {
	list := []rune(s)
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	_this.list = list
}

func (_this *Skip) SetSorted(s string) {
	list := []rune(s)
	_this.list = list
}

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

func (_this *Skip) String() string {
	return string(_this.list)
}
