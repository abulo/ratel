package util

import (
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// ZhCharToFirstPinyin 将中文转换成小写字母
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

// ZhCharToPinyin 将中文转换成小写字母
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
