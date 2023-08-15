package filter

import (
	"fmt"
)

type Result struct {
	Word    string `json:"word"`    // 匹配到的敏感词
	Matched string `json:"matched"` // 匹配到的字符串
	Start   int    `json:"start"`   // 原始字符串中匹配到的起始位置
	End     int    `json:"end"`     // 原始字符串中匹配到的结束位置
}

func (_this *Result) String() string {
	return fmt.Sprintf("word:%s mathced:%s start:%d end:%d;", _this.Word, _this.Matched, _this.Start, _this.End)
}
