package correct

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// charWidthMap 记录全角字母、数字与全角空格到对应半角字符的映射。
	charWidthMap = map[rune]rune{
		'ａ': 'a', 'ｂ': 'b', 'ｃ': 'c', 'ｄ': 'd', 'ｅ': 'e', 'ｆ': 'f', 'ｇ': 'g', 'ｈ': 'h', 'ｉ': 'i', 'ｊ': 'j', 'ｋ': 'k', 'ｌ': 'l', 'ｍ': 'm', 'ｎ': 'n', 'ｏ': 'o', 'ｐ': 'p', 'ｑ': 'q', 'ｒ': 'r', 'ｓ': 's', 'ｔ': 't', 'ｕ': 'u', 'ｖ': 'v', 'ｗ': 'w', 'ｘ': 'x', 'ｙ': 'y', 'ｚ': 'z', 'Ａ': 'A', 'Ｂ': 'B', 'Ｃ': 'C', 'Ｄ': 'D', 'Ｅ': 'E', 'Ｆ': 'F', 'Ｇ': 'G', 'Ｈ': 'H', 'Ｉ': 'I', 'Ｊ': 'J', 'Ｋ': 'K', 'Ｌ': 'L', 'Ｍ': 'M', 'Ｎ': 'N', 'Ｏ': 'O', 'Ｐ': 'P', 'Ｑ': 'Q', 'Ｒ': 'R', 'Ｓ': 'S', 'Ｔ': 'T', 'Ｕ': 'U', 'Ｖ': 'V', 'Ｗ': 'W', 'Ｘ': 'X', 'Ｙ': 'Y', 'Ｚ': 'Z', '１': '1', '２': '2', '３': '3', '４': '4', '５': '5', '６': '6', '７': '7', '８': '8', '９': '9', '０': '0', '　': ' ',
	}

	// halfTimeRe 匹配数字之间的全角冒号（如 12：00），用于替换为半角冒号。
	halfTimeRe = regexp.MustCompile(`(\d)(：)(\d)`)
)

// halfwidth 将全角字母、数字与空格转换为半角，并把时间写法中的全角冒号替换为半角冒号。
func halfwidth(text string) string {
	// PERF-9-15 [P2] 按 rune 数量预分配切片容量，避免零起步反复扩容拷贝；同一次 map 查找结果直接复用。
	runes := make([]rune, 0, utf8.RuneCountInString(text))
	for _, char := range text {
		if newChar, ok := charWidthMap[char]; ok {
			runes = append(runes, newChar)
		} else {
			runes = append(runes, char)
		}
	}

	out := string(runes)

	out = halfTimeRe.ReplaceAllStringFunc(out, func(part string) string {
		return strings.Replace(part, "：", ":", 1)
	})

	return out
}
