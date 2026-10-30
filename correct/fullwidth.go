package correct

import (
	"regexp"
)

var (
	// fullwidthMaps 记录需要替换的半角标点到对应全角标点的映射。
	fullwidthMaps = map[string]string{
		",": "，",
		".": "。",
		";": "；",
		":": "：",
		"!": "！",
		"?": "？",
		"~": "～",
	}

	spcialPunctuations = `[.:]`
	normalPunctuations = `[,\!\?~]`

	punctuationWithLeftCJKRe        = regexp.MustCompile(normalPunctuations + `[` + cjk + `]+`)
	punctuationWithRightCJKRe       = regexp.MustCompile(`[` + cjk + `]+` + normalPunctuations)
	punctuationWithSpeicalCJKRe     = regexp.MustCompile(`[` + cjk + `]+` + spcialPunctuations + `[` + cjk + `]+`)
	punctuationWithSpeicalLastCJKRe = regexp.MustCompile(`[` + cjk + `]+` + spcialPunctuations + "$")
	punctuationsRe                  = regexp.MustCompile(`(` + spcialPunctuations + `|` + normalPunctuations + `)`)
)

// fullwidth 将紧邻 CJK 字符的半角标点替换为对应全角标点，覆盖标点位于 CJK 左侧、
// 右侧、两侧以及位于文本末尾的情形。
func fullwidth(text string) (out string) {
	out = text

	out = punctuationWithLeftCJKRe.ReplaceAllStringFunc(out, fullwidthReplacePart)
	out = punctuationWithRightCJKRe.ReplaceAllStringFunc(out, fullwidthReplacePart)
	out = punctuationWithSpeicalCJKRe.ReplaceAllStringFunc(out, fullwidthReplacePart)
	out = punctuationWithSpeicalLastCJKRe.ReplaceAllStringFunc(out, fullwidthReplacePart)

	return
}

func fullwidthReplacePart(part string) string {
	part = punctuationsRe.ReplaceAllStringFunc(part, func(str string) string {
		str = fullwidthMaps[str]
		return str
	})

	return part

}
