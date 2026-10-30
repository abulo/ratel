package correct

import "regexp"

// strategery 封装一条排版规则：space 为 true 时按正向、反向正则在两类字符之间补空格，
// 否则去除两类字符之间已有的空格；reverse 控制是否同时处理反向顺序。
type strategery struct {
	addSpaceRe           *regexp.Regexp
	addSpaceReverseRe    *regexp.Regexp
	removeSpaceRe        *regexp.Regexp
	removeSpaceReverseRe *regexp.Regexp

	space   bool
	reverse bool
}

func newStrategery(one, other string, space, reverse bool) *strategery {
	addSpaceStr := "(" + one + `)(` + other + ")"
	addSpaceReverseStr := "(" + other + `)(` + one + ")"

	removeSpaceStr := `(` + one + `)` + spaceRe + `(` + other + `)`
	removeSpaceReverseStr := "(" + other + `)` + spaceRe + `(` + one + ")"

	return &strategery{
		addSpaceRe:           regexp.MustCompile(addSpaceStr),
		addSpaceReverseRe:    regexp.MustCompile(addSpaceReverseStr),
		removeSpaceRe:        regexp.MustCompile(removeSpaceStr),
		removeSpaceReverseRe: regexp.MustCompile(removeSpaceReverseStr),
		space:                space,
		reverse:              reverse,
	}
}

func (s *strategery) format(in string) (out string) {
	out = in
	if s.space {
		out = s.addSpace(out)
	} else {
		out = s.removeSpace(out)
	}

	return
}

func (s *strategery) addSpace(in string) (out string) {
	out = in
	out = s.addSpaceRe.ReplaceAllString(out, "$1 $2")
	if s.reverse {
		out = s.addSpaceReverseRe.ReplaceAllString(out, "$1 $2")
	}

	return
}

func (s *strategery) removeSpace(in string) (out string) {
	out = in
	out = s.removeSpaceRe.ReplaceAllString(out, "$1 $2")
	if s.reverse {
		out = s.removeSpaceReverseRe.ReplaceAllString(out, "$1 $2")
	}

	return
}
