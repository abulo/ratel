package util

import (
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// ZhCharToFirstPinyin 将中文字符串转换为各汉字拼音首字母组成的小写字符串。
//
// 非中文的数字与字母原样保留，其余字符被忽略；结果统一转为小写。
func ZhCharToFirstPinyin(p string) string {
	var a = pinyin.NewArgs()
	var s strings.Builder
	a.Style = pinyin.FirstLetter
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s.WriteString(string(pinyin.Pinyin(string(r), a)[0][0]))
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
	var a = pinyin.NewArgs()
	var s strings.Builder
	a.Style = pinyin.Normal
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s.WriteString(string(pinyin.Pinyin(string(r), a)[0][0]))
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			s.WriteString(string(r))
		}
	}
	return StrToLower(s.String())
}
