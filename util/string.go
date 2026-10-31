package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/spf13/cast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// tagPatterns 保存需要在 Tag 中统一替换为 "/" 的标签分隔符及其替换值，按两两配对排列。
var tagPatterns = []string{
	",", "/",
	"，", "/",
	"|", "/",
	"丨", "/",
	"、", "/",
	";", "/",
}

// StrTrim 去除字符串中的全部空白字符（包括中间位置的空格、制表符等），而非仅去除两端。
func StrTrim(str string) string {
	return strings.Join(strings.Fields(str), "")
}

// StrPos 返回 needle 在 haystack 中从 offset 起首次出现的字节位置；
// 找不到、haystack 为空或 offset 越界时返回 -1。
func StrPos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// StrIPos 忽略大小写返回 needle 在 haystack 中从 offset 起首次出现的字节位置；
// 找不到或 offset 越界时返回 -1。
func StrIPos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	haystack = haystack[offset:]
	if offset < 0 {
		offset += length
	}
	// PERF-7-4 [P2] 窗口 EqualFold 忽略大小写匹配，避免先把整段 haystack 转小写产生一次整串拷贝
	pos := indexFold(haystack, needle)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// StrRPos 返回 needle 在 haystack 中最后一次出现的字节位置。
//
// offset 为正时限定在 offset 之后的范围内搜索，为负时限定在其之前的范围内；
// 找不到或 offset 越界时返回 -1。
func StrRPos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// StrRIPos 忽略大小写返回 needle 在 haystack 中最后一次出现的字节位置；
// 范围规则与 StrRPos 相同，找不到或越界时返回 -1。
func StrRIPos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	// PERF-7-4 [P2] 窗口 EqualFold 反向匹配，避免整段 haystack 转小写拷贝
	pos = lastIndexFold(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// indexFold 返回 needle 在 haystack 中首次忽略大小写出现的字节位置；找不到返回 -1；needle 为空串返回 0。
func indexFold(haystack, needle string) int {
	nl := len(needle)
	if nl == 0 {
		return 0
	}
	for i := 0; i+nl <= len(haystack); i++ {
		if strings.EqualFold(haystack[i:i+nl], needle) {
			return i
		}
	}
	return -1
}

// lastIndexFold 返回 needle 在 haystack 中最后一次忽略大小写出现的字节位置；找不到返回 -1；needle 为空串返回 len(haystack)。
func lastIndexFold(haystack, needle string) int {
	nl := len(needle)
	if nl == 0 {
		return len(haystack)
	}
	for i := len(haystack) - nl; i >= 0; i-- {
		if strings.EqualFold(haystack[i:i+nl], needle) {
			return i
		}
	}
	return -1
}

// StrReplace 将 subject 中前 count 处 search 替换为 replace；count 为负数时替换全部。
func StrReplace(search, replace, subject string, count int) string {
	return strings.Replace(subject, search, replace, count)
}

// StrToUpper 将字符串全部字符转换为大写。
func StrToUpper(str string) string {
	return strings.ToUpper(str)
}

// StrToLower 将字符串全部字符转换为小写。
func StrToLower(str string) string {
	return strings.ToLower(str)
}

// UCfirst 将字符串首字符转换为大写，其余字符保持不变；空串返回空字符串。
func UCfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// LCFirst 将字符串首字符转换为小写，其余字符保持不变；空串返回空字符串。
func LCFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// PERF-7-3 [P2] cases.Title 大写器可并发复用，包级构造一次，避免每次调用重复初始化语言规则
var titleCaser = cases.Title(language.English)

// UCWords 将字符串按单词做标题式大小写处理，每个单词首字母大写。
func UCWords(str string) string {
	return titleCaser.String(str)
}

// SubStr 按字节截取子串。
//
// start 为负时从末尾倒数；length 为 -1 时取到末尾，为 0 时返回空字符串，
// 小于 -1 时原样返回整串。
func SubStr(str string, start int, length int) string {
	if start < 0 {
		start = max(len(str)+start, 0)
	}
	if length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := min(int(start)+length, len(str))
	return str[start:end]
}

// StrRev 按 Unicode 字符（rune）反转字符串。
func StrRev(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ParseStr 将查询字符串解析进 result 指向的嵌套 map。
//
// 支持 f1=m 形式的平铺键、f[a]=m 形式的嵌套 map，以及 f[]=m、f[a][]=m 形式的
// 切片；同名键类型冲突或 URL 解码失败时返回错误。暂不支持多层嵌套切片，
// 与 PHP parse_str 的行为存在差异。
func ParseStr(encodedString string, result map[string]any) error {
	var build func(map[string]any, []string, any) error

	build = func(result map[string]any, keys []string, value any) error {
		length := len(keys)
		key := strings.Trim(keys[0], "'\"")
		if length == 1 {
			result[key] = value
			return nil
		}

		if keys[1] == "" && length == 2 {
			// TODO(owner): 暂不支持更深层的嵌套切片。
			if key == "" {
				return nil
			}
			val, ok := result[key]
			if !ok {
				result[key] = []any{value}
				return nil
			}
			children, ok := val.([]any)
			if !ok {
				return fmt.Errorf("expected type '[]any' for key '%s', but got '%T'", key, val)
			}
			result[key] = append(children, value)
			return nil
		}

		if keys[1] == "" && length > 2 && keys[2] != "" {
			val, ok := result[key]
			if !ok {
				result[key] = []any{}
				val = result[key]
			}
			children, ok := val.([]any)
			if !ok {
				return fmt.Errorf("expected type '[]any' for key '%s', but got '%T'", key, val)
			}
			if l := len(children); l > 0 {
				if child, ok := children[l-1].(map[string]any); ok {
					if _, ok := child[keys[2]]; !ok {
						_ = build(child, keys[2:], value)
						return nil
					}
				}
			}
			child := map[string]any{}
			_ = build(child, keys[2:], value)
			result[key] = append(children, child)

			return nil
		}

		val, ok := result[key]
		if !ok {
			result[key] = map[string]any{}
			val = result[key]
		}
		children, ok := val.(map[string]any)
		if !ok {
			return fmt.Errorf("expected type 'map[string]any' for key '%s', but got '%T'", key, val)
		}

		return build(children, keys[1:], value)
	}

	parts := strings.SplitSeq(encodedString, "&")
	for part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}
		for key[0] == ' ' {
			key = key[1:]
		}
		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		var first strings.Builder
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first.WriteString("_")
			} else {
				first.WriteString(string(chr))
			}
			if chr == '[' {
				first.WriteString(keys[0][i+1:])
				break
			}
		}
		keys[0] = first.String()

		if err := build(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

// NumberFormat 按小数位数 decimals、小数点 decPoint 与千分位分隔符 thousandsSep
// 格式化数字；负数保留负号，小数部分按四舍五入处理。
func NumberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}

// ChunkSplit 将字符串按 chunklen 长度切分，每段末尾追加 end；end 为空时默认 \r\n。
func ChunkSplit(body string, chunklen uint, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, erunes := []rune(body), []rune(end)
	l := uint(len(runes))
	if l <= 1 || l < chunklen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(erunes))
	var i uint
	for i = 0; i < l; i += chunklen {
		if i+chunklen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunklen]...)
		}
		ns = append(ns, erunes...)
	}
	return string(ns)
}

// StrWordCount 按空白字符分词，返回单词组成的切片。
func StrWordCount(str string) []string {
	return strings.Fields(str)
}

// WordWrap 按宽度 width 折行，折行处插入 br。
//
// cut 为 true 时允许在单词中间强制截断；br 为空字符串时 panic，width 为 0 且
// cut 为 true 时 panic。
func WordWrap(str string, width uint, br string, cut bool) string {
	strlen := len(str)
	brlen := len(br)
	linelen := int(width)

	if strlen == 0 {
		return ""
	}
	if brlen == 0 {
		panic("break string cannot be empty")
	}
	if linelen == 0 && cut {
		panic("can't force cut when width is zero")
	}

	current, laststart, lastspace := 0, 0, 0
	var ns []byte
	for current = 0; current < strlen; current++ {
		if str[current] == br[0] && current+brlen < strlen && str[current:current+brlen] == br {
			ns = append(ns, str[laststart:current+brlen]...)
			current += brlen - 1
			lastspace = current + 1
			laststart = lastspace
		} else if str[current] == ' ' {
			if current-laststart >= linelen {
				ns = append(ns, str[laststart:current]...)
				ns = append(ns, br[:]...)
				laststart = current + 1
			}
			lastspace = current
		} else if current-laststart >= linelen && cut && laststart >= lastspace {
			ns = append(ns, str[laststart:current]...)
			ns = append(ns, br[:]...)
			laststart = current
			lastspace = current
		} else if current-laststart >= linelen && laststart < lastspace {
			ns = append(ns, str[laststart:lastspace]...)
			ns = append(ns, br[:]...)
			lastspace++
			laststart = lastspace
		}
	}

	if laststart != current {
		ns = append(ns, str[laststart:current]...)
	}
	return string(ns)
}

// StrLen 返回字符串的字节长度。
func StrLen(str string) int {
	return len(str)
}

// MbStrLen 返回字符串的 Unicode 字符数（rune 数），多字节字符按一个字符计。
func MbStrLen(str string) int {
	return utf8.RuneCountInString(str)
}

// StrRepeat 将 input 重复 multiplier 次后返回。
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// StrStr 返回从 needle 末字符位置起（仍包含 needle 的末字符）的子串；
// needle 为空或找不到时返回空字符串。
func StrStr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// StrTr 实现字符或子串替换。
//
// 传入一个 map[string]string 时按键值对替换对应子串；传入两个字符串时按位置
// 逐字符映射替换；参数数量不符时原样返回 haystack；实参类型与上述形式不符时触发
// 类型断言 panic。
func StrTr(haystack string, params ...any) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}
		// PERF-7-8 [P2] 容量按 length*2 预分配后直接 append，避免原 make 零长度再 append 扩容到 4 倍并留出空串槽位
		oldnew := make([]string, 0, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}
		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)
		if trlen > lt {
			trlen = lt
		}
		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8
		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}
		return string(str)
	}

	return haystack
}

// StrShuffle 随机打乱字符串的字符顺序并返回新字符串。
func StrShuffle(str string) string {
	runes := []rune(str)
	// PERF-7-14 [P2] 复用进程级共享随机源，不再每次新建 rand.Rand 与随机源（预期收益需实测验证）
	s := make([]rune, len(runes))
	for i, v := range sharedPerm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Trim 去除字符串两端空白；提供 characterMask 时改为去除两端指定字符集中的字符。
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// LTrim 去除字符串左端空白；提供 characterMask 时改为去除左端指定字符集中的字符。
func LTrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, characterMask[0])
}

// RTrim 去除字符串右端空白；提供 characterMask 时改为去除右端指定字符集中的字符。
func RTrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, characterMask[0])
}

// Explode 按 delimiter 切分字符串 str，返回各段组成的切片。
func Explode(delimiter, str string) []string {
	return strings.Split(str, delimiter)
}

// Chr 将 ASCII 或 Unicode 码点转换为对应字符组成的字符串。
func Chr(ascii int) string {
	return string(rune(ascii))
}

// Ord 返回字符串首字符的 Unicode 码点。
func Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

// Nl2br 将字符串中的换行符替换为 HTML 换行标签。
//
// isXhtml 为 true 时使用 <br />，否则使用 <br>；连续的 CR/LF 组合只插入一个标签。
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

// JSONDecode 将 JSON 字节流反序列化到 val 指向的对象；失败返回错误。
func JSONDecode(data []byte, val any) error {
	return json.Unmarshal(data, val)
}

// JSONEncode 将 val 序列化为 JSON 字节；失败返回错误。
func JSONEncode(val any) ([]byte, error) {
	return json.Marshal(val)
}

// JSONString 将对象序列化为 JSON 字符串；序列化错误被忽略，失败时返回空字符串。
func JSONString(obj any) string {
	jsByte, _ := json.Marshal(obj)
	// PERF-7-1 [P1] 去掉 cast.ToString(any) 的接口装箱与类型断言，直接 []byte→string，序列化结果与原实现逐字节一致
	return string(jsByte)
}

// AddSlashes 为字符串中的单引号、双引号与反斜杠添加反斜杠转义。
func AddSlashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// StripSlashes 去除字符串中的反斜杠转义；连续两个反斜杠视为一个字面反斜杠并跳过其后者。
func StripSlashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// QuoteMeta 为字符串中的正则元字符添加反斜杠转义。
func QuoteMeta(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// HtmlEntities 转义字符串中的 HTML 特殊字符。
func HtmlEntities(str string) string {
	return html.EscapeString(str)
}

// HTMLEntityDecode 反转义字符串中的 HTML 实体。
func HTMLEntityDecode(str string) string {
	return html.UnescapeString(str)
}

// Md5 返回字符串的 MD5 十六进制摘要。
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Md5File 计算文件内容的 MD5 十六进制摘要。
//
// 小于 1MB 的文件整体读取，较大文件分块读取；打开或读取失败时返回空字符串与错误。
func Md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	var size int64 = 1048576
	hash := md5.New()

	if fi.Size() < size {
		// PERF-7-5 [P2] 复用已打开的文件句柄分块读取，避免小文件分支再次 os.ReadFile 二次打开磁盘文件
		if _, err := io.Copy(hash, f); err != nil {
			return "", err
		}
	} else {
		b := make([]byte, size)
		for {
			n, err := f.Read(b)
			if err != nil {
				break
			}

			hash.Write(b[:n])
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Sha1 返回字符串的 SHA1 十六进制摘要。
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha1File 计算整个文件内容的 SHA1 十六进制摘要；读取失败时返回空字符串与错误。
func Sha1File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha1.New()
	// PERF-7-6 [P2] 与 Md5File 对齐，1MB 缓冲分块读取，避免大文件一次性整包读入内存
	if _, err := io.CopyBuffer(hash, f, make([]byte, 1048576)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Crc32 返回字符串的 IEEE CRC32 校验值。
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// Levenshtein 计算两字符串之间的编辑距离，costIns、costRep、costDel 分别指定
// 插入、替换、删除代价；任一字符串长度超过 255 时返回 -1。
func Levenshtein(str1, str2 string, costIns, costRep, costDel int) int {
	var maxLen = 255
	l1 := len(str1)
	l2 := len(str2)
	if l1 == 0 {
		return l2 * costIns
	}
	if l2 == 0 {
		return l1 * costDel
	}
	if l1 > maxLen || l2 > maxLen {
		return -1
	}

	p1 := make([]int, l2+1)
	p2 := make([]int, l2+1)
	var c0, c1, c2 int
	var i1, i2 int
	for i2 := 0; i2 <= l2; i2++ {
		p1[i2] = i2 * costIns
	}
	for i1 = range l1 {
		p2[0] = p1[0] + costDel
		for i2 = range l2 {
			if str1[i1] == str2[i2] {
				c0 = p1[i2]
			} else {
				c0 = p1[i2] + costRep
			}
			c1 = p1[i2+1] + costDel
			if c1 < c0 {
				c0 = c1
			}
			c2 = p2[i2] + costIns
			if c2 < c0 {
				c0 = c2
			}
			p2[i2+1] = c0
		}
		tmp := p1
		p1 = p2
		p2 = tmp
	}
	c0 = p1[l2]

	return c0
}

// SimilarText 返回两字符串最长公共片段的累计匹配长度。
//
// percent 非 nil 时通过该指针写入相似度百分比（0–100）；两个字符串均为空时返回 0。
func SimilarText(first, second string, percent *float64) int {
	// PERF-7-7 [P2] 入口处把两个字符串各转一次 []byte 向下传递，消除递归每层重复 []byte(str)/string(...) 往返拷贝（算法复杂度未变更，预期收益需实测验证）
	var similarText func(b1, b2 []byte, len1, len2 int) int
	similarText = func(b1, b2 []byte, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		for i := range len1 {
			for j := range len2 {
				for l := 0; (i+l < len1) && (j+l < len2) && (b1[i+l] == b2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(b1, b2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				sum += similarText(b1, b2, len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText([]byte(first), []byte(second), l1, l2)
	if percent != nil {
		*percent = float64(sim*200) / float64(l1+l2)
	}
	return sim
}

// Soundex 计算字符串的 Soundex 语音编码（固定 4 位）；空字符串会 panic。
func Soundex(str string) string {
	if str == "" {
		panic("str: cannot be an empty string")
	}
	// 编码表按字母 A–Z 顺序排列，索引即字母在字母表中的位置，0 表示不产生编码。
	table := [26]rune{
		'0', '1', '2', '3',
		'0', '1', '2',
		'0',
		'0', '2', '2', '4', '5', '5',
		'0', '1', '2', '6', '2', '3',
		'0', '1',
		'0', '2',
		'0', '2',
	}
	last, code, small := -1, 0, 0
	sd := make([]rune, 4)
	for i := 0; i < len(str) && small < 4; i++ {
		char := str[i]
		if char < '\u007F' && 'a' <= char && char <= 'z' {
			code = int(char - 'a' + 'A')
		} else {
			code = int(char)
		}
		if code >= 'A' && code <= 'Z' {
			if small == 0 {
				sd[small] = rune(code)
				small++
				last = int(table[code-'A'])
			} else {
				code = int(table[code-'A'])
				if code != last {
					if code != 0 {
						sd[small] = rune(code)
						small++
					}
					last = code
				}
			}
		}
	}
	for ; small < 4; small++ {
		sd[small] = '0'
	}
	return string(sd)
}

// PERF-7-2 [P1] tagPatterns 为静态分隔符表，Replacer 包级构建一次并发复用，避免每次调用现场 NewReplacer
var tagReplacer = strings.NewReplacer(tagPatterns...)

// Tag 将标签字符串中的常见分隔符统一替换为 "/"。
func Tag(i any) string {
	return tagReplacer.Replace(cast.ToString(i))
}

// StringToBytes 零拷贝将 string 转换为 []byte。
//
// 返回的切片与原字符串共享底层内存，不得修改其内容，否则会破坏字符串的不可变性。
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
