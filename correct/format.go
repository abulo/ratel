package correct

import (
	"regexp"
)

const (
	cjk     = `\p{Han}|\p{Hangul}|\p{Hanunoo}|\p{Katakana}|\p{Hiragana}|\p{Bopomofo}`
	spaceRe = `[ ]`
)

var (
	strategies   []*strategery
	dashHansRe   = regexp.MustCompile(`([` + cjk + `）】」》”’])([\-]+)([` + cjk + `（【「《“‘])`)
	leftQuoteRe  = regexp.MustCompile(spaceRe + `([（【「《])`)
	rightQuoteRe = regexp.MustCompile(`([）】」》])` + spaceRe)
)

// registerStrategery 注册一条排版规则，规则在 Format 流程中按注册顺序依次生效。
func registerStrategery(one, other string, space, reverse bool) {
	strategies = append(strategies, newStrategery(one, other, space, reverse))
}

func init() {
	registerStrategery(cjk, `[a-zA-Z]`, true, true)

	registerStrategery(cjk, `[0-9]`, true, true)

	registerStrategery(cjk, `[\|+*]`, true, true)
	registerStrategery(cjk, `[@]`, true, false)
	registerStrategery(cjk, `[\[\(‘“]`, true, false)
	registerStrategery(`[’”\]\)!%]`, cjk, true, false)
	registerStrategery(`[”\]\)!]`, `[a-zA-Z0-9]+`, true, false)

	registerStrategery(`[\w`+cjk+`]`, `[，。！？：；）」》】”’]`, false, true)
	registerStrategery(`[‘“【「《（]`, `[\w`+cjk+`]`, false, true)
}

// spaceDashWithHans 在 CJK 字符与连字符之间补空格，并去除左、右引号括号内侧的多余空格。
func spaceDashWithHans(in string) (out string) {
	out = dashHansRe.ReplaceAllString(in, "$1 $2 $3")
	out = leftQuoteRe.ReplaceAllString(out, "$1")
	out = rightQuoteRe.ReplaceAllString(out, "$1")
	return out
}

// Option 表示 Format 内置规则执行完毕后追加运行的自定义文本处理步骤。
type Option interface {
	Format(text string) string
}

// Format 对中英文混排文本进行自动排版：先做全半角规范化，再按内置规则在
// CJK 与英文、数字、符号之间补空格或去除多余空格，最后处理连字符与引号括号；
// 可选的 options 在内置规则之后依次对文本进行二次加工。
func Format(in string, options ...Option) (out string) {
	out = in

	out = halfwidth(out)
	out = fullwidth(out)

	for _, s := range strategies {
		out = s.format(out)
	}

	out = spaceDashWithHans(out)

	for _, opt := range options {
		out = opt.Format(out)
	}

	return
}
