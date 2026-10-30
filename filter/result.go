package filter

import (
	"fmt"
)

// Result 描述一次敏感词命中的词库词条、原文匹配片段及其在原文中的字节位置。
type Result struct {
	// Word 为词库中命中的敏感词。
	Word    string `json:"word"`
	// Matched 为原文中实际匹配到的字符串。
	Matched string `json:"matched"`
	// Start 为匹配片段在原始字符串中的起始字节位置。
	Start   int    `json:"start"`
	// End 为匹配片段在原始字符串中的结束字节位置。
	End     int    `json:"end"`
}

// String 返回包含敏感词、匹配串与起止位置的可读调试描述。
func (_this *Result) String() string {
	return fmt.Sprintf("word:%s mathced:%s start:%d end:%d;", _this.Word, _this.Matched, _this.Start, _this.End)
}
