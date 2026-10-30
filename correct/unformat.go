package correct

import "regexp"

var (
	// removeSpaceRe 匹配 CJK 字符两侧的空白，用于去除 CJK 与相邻字符之间的多余空格。
	removeSpaceRe = regexp.MustCompile(`(` + spaceRe + `+)?(` + cjk + `)(` + spaceRe + `+)?`)
)

// UnformatOption 表示 Unformat 去空格流程执行完毕后追加运行的自定义文本处理步骤。
type UnformatOption interface {
	Unformat(text string) string
}

// Unformat 去除文本中 CJK 字符两侧的空格；可选 options 在去空格之后依次对文本二次加工。
func Unformat(text string, options ...UnformatOption) string {
	text = removeSpaceRe.ReplaceAllString(text, "$2")

	for _, opt := range options {
		text = opt.Unformat(text)
	}

	return text
}
