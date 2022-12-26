package util

import (
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// ZhCharToFirstPinyin 将中文转换成小写字母
func ZhCharToFirstPinyin(p string) string {
	var a = pinyin.NewArgs()
	var s string
	a.Style = pinyin.FirstLetter
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s += string(pinyin.Pinyin(string(r), a)[0][0])
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			s += string(r)
		}
	}
	return StrToLower(s)
}

// ZhCharToPinyin 将中文转换成小写字母
func ZhCharToPinyin(p string) string {
	var a = pinyin.NewArgs()
	var s string
	a.Style = pinyin.Normal
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s += string(pinyin.Pinyin(string(r), a)[0][0])
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			s += string(r)
		}
	}
	return StrToLower(s)
}
