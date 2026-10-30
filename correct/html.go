package correct

import (
	"bytes"
	"io"
	"regexp"

	"github.com/pkg/errors"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
)

var (
	// ignoreTagsRe 匹配需要跳过内部文本排版的标签：pre、script、style、textarea。
	ignoreTagsRe = regexp.MustCompile("(?mi)<(pre|script|style|textarea)")
)

// FormatHTML 解析 HTML 文本，仅对普通文本节点执行 Format 排版；标签、属性以及
// pre、script、style、textarea 内部内容保持原样。HTML 词法解析失败时返回错误。
func FormatHTML(body string, options ...Option) (out string, err error) {
	return processHTML(body, func(text string) string {
		return Format(text, options...)
	})
}

// UnformatHTML 解析 HTML 文本，仅对普通文本节点执行 Unformat 去除空格；
// 其余内容保持原样。HTML 词法解析失败时返回错误。
func UnformatHTML(body string, options ...UnformatOption) (out string, err error) {
	return processHTML(body, func(text string) string {
		return Unformat(text, options...)
	})
}

// processHTML 遍历 HTML 词法流，对文本节点调用 fn 处理后写回；标签与被忽略标签的
// 内部内容原样写回；遇到非 EOF 的词法错误时返回错误。
func processHTML(body string, fn func(plainText string) string) (out string, err error) {
	w := &bytes.Buffer{}
	lex := html.NewLexer(parse.NewInputString(body))
	out = body

	ignoreTag := false

	for {
		t, data := lex.Next()

		switch t {
		case html.ErrorToken:
			if lex.Err() == io.EOF {
				return w.String(), nil
			}

			err = errors.Errorf("Error on line %v", lex.Err())
			return
		case html.TextToken:
			if ignoreTag {
				if _, err := w.Write(data); err != nil {
					return out, err
				}

				ignoreTag = false
				continue
			}

			formated := fn(string(data))
			if _, err := w.Write([]byte(formated)); err != nil {
				return out, err
			}
		case html.StartTagToken:
			if ignoreTagsRe.Match(data) {
				ignoreTag = true
			}

			if _, err := w.Write(data); err != nil {
				return out, err
			}

		default:
			if _, err := w.Write(data); err != nil {
				return out, err
			}
		}
	}
}
