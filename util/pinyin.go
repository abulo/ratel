package util

import (
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// PERF-7-24 [P2] 两种样式的 pinyin.Args 包级复用，不再每次调用 NewArgs 构造
var (
	firstLetterArgs = func() pinyin.Args {
		a := pinyin.NewArgs()
		a.Style = pinyin.FirstLetter
		return a
	}()
	normalArgs = func() pinyin.Args {
		a := pinyin.NewArgs()
		a.Style = pinyin.Normal
		return a
	}()
)

// ZhCharToFirstPinyin 将中文字符串转换为各汉字拼音首字母组成的小写字符串。
//
// 非中文的数字与字母原样保留，其余字符被忽略；结果统一转为小写。
func ZhCharToFirstPinyin(p string) string {
	// PERF-7-24 [P2] 整串一次 Pinyin 取得全部拼音二维切片，再按原 rune 顺序对齐，替代逐字符 Pinyin 调用带来的重复字典查找与切片分配
	all := pinyin.Pinyin(p, firstLetterArgs)
	var s strings.Builder
	hi := 0
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s.WriteString(all[hi][0])
			hi++
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			s.WriteString(string(r))
		}
	}
	return StrToLower(s.String())
}

// ZhCharToPinyin 将中文字符串转换为全拼小写字符串。
//
// 非中文的数字与字母原样保留，其余字符被忽略；结果统一转为小写。
func ZhCharToPinyin(p string) string {
	// PERF-7-24 [P2] 整串一次 Pinyin 取得全部拼音二维切片，再按原 rune 顺序对齐，替代逐字符 Pinyin 调用
	all := pinyin.Pinyin(p, normalArgs)
	var s strings.Builder
	hi := 0
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s.WriteString(all[hi][0])
			hi++
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			s.WriteString(string(r))
		}
	}
	return StrToLower(s.String())
}
