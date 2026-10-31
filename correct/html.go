package correct

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
)

// ignoreTagNames 为需要跳过内部文本排版的标签名：pre、script、style、textarea（小写）。
var ignoreTagNames = map[string]struct{}{
	"pre":      {},
	"script":   {},
	"style":    {},
	"textarea": {},
}

// isIgnoreTag 判断开始标签原始字节 data 是否属于需要忽略内部内容的标签。
//
// PERF-9-16 [P2]：替代原先每个开始标签执行一次的 ignoreTagsRe 正则；与原正则一致地
// 做大小写不敏感的前缀匹配（保留原正则对 `<prex>` 这类带前缀标签名的宽松匹配行为）。
func isIgnoreTag(data []byte) bool {
	if len(data) == 0 || data[0] != '<' {
		return false
	}
	end := 1
	for end < len(data) {
		c := data[end]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			break
		}
		end++
	}
	if end == 1 {
		return false
	}
	name := make([]byte, end-1)
	for i := 1; i < end; i++ {
		c := data[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		name[i-1] = c
	}
	for tag := range ignoreTagNames {
		if bytes.HasPrefix(name, []byte(tag)) {
			return true
		}
	}
	return false
}

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
	// PERF-9-16 [P2] 按原文长度预分配输出缓冲容量，避免小步动态扩容。
	w := bytes.Buffer{}
	w.Grow(len(body))
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
			// PERF-9-16 [P2] 直接以 string 写入缓冲，去掉 formated → []byte 的第二次显式转换分配。
			if _, err := w.WriteString(formated); err != nil {
				return out, err
			}
		case html.StartTagToken:
			// PERF-9-16 [P2] 用标签名字节前缀判定替代每标签一次正则匹配。
			if isIgnoreTag(data) {
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
