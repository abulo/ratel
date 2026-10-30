package util

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// MarshalHTML 将 v 序列化为 JSON 后包装为 template.HTML 类型，供 HTML 模板原样输出；
// 序列化错误被忽略，失败时内容为空。
func MarshalHTML(v any) template.HTML {
	a, _ := json.Marshal(v)
	return template.HTML(a)
}

// MarshalJS 将 v 序列化为 JSON 后包装为 template.JS 类型，供模板在脚本上下文中使用；
// 序列化错误被忽略，失败时内容为空。
func MarshalJS(v any) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

// JS 将字符串直接包装为 template.JS 类型，不做额外转义。
func JS(v string) template.JS {
	return template.JS(v)
}

// Static 拼接静态资源路径，返回 "/p/v" 形式的相对路径。
func Static(p, v string) string {
	return "/" + p + "/" + v
}

// DomainStatic 拼接静态资源路径，返回 "p/v" 形式的路径（不带前导斜杠）。
func DomainStatic(p, v string) string {
	return p + "/" + v
}

// UnescapeString 将 v 序列化为 JSON 字符串后做 HTML 转义并返回；序列化错误被忽略。
func UnescapeString(v any) string {
	a, _ := json.Marshal(v)
	return html.EscapeString(string(a))
}

// GetAppRootPath 返回当前可执行文件所在目录的上一级目录。
func GetAppRootPath() string {
	return GetParentDirectory(GetCurrentDirectory())
}

// substr 按 rune 截取 s 中从 pos 起长度为 length 的子串。
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := min(pos+length, len(runes))
	return string(runes[pos:l])
}

// GetParentDirectory 返回 dir 中最后一个 "/" 之前的前缀；dir 中不存在斜杠时，
// 会按长度 -1 的子串处理，返回去掉末尾一个字符后的内容。
func GetParentDirectory(dir string) string {
	return substr(dir, 0, strings.LastIndex(dir, "/"))
}

// GetCurrentDirectory 返回当前可执行文件所在目录的绝对路径，路径中的反斜杠统一为斜杠；
// 取绝对路径失败时调用 log.Fatal 直接终止进程。
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Random 返回 UUID 拼接当前纳秒时间戳的随机字符串。
func Random() string {
	return uuid.New().String() + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// VerifyPhone 判断 phone 是否为中国大陆手机号格式。
func VerifyPhone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

// VerifyEmail 判断 email 是否为合法邮箱地址格式。
func VerifyEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyIPv4 判断 address 是否为合法 IPv4 地址；正则编译错误被忽略并按不匹配处理。
func VerifyIPv4(address string) bool {
	res, _ := regexp.MatchString(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`, address)
	return res
}

// VerifyIPv6 判断 address 是否为合法 IPv6 地址；正则编译错误被忽略并按不匹配处理。
func VerifyIPv6(address string) bool {
	res, _ := regexp.MatchString(`^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$`, address)
	return res
}

// FunctionName 返回函数值 i 对应的运行时函数完整名称；i 不是函数类型时会 panic。
func FunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type intSlice []int

func (s intSlice) Len() int { return len(s) }

func (s intSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s intSlice) Less(i, j int) bool { return s[i] < s[j] }

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// GetMissingElement 找出 1 起始连续整数切片 arr 中缺失的最小整数。
//
// 切片为空时返回 1；切片连续完整时返回末尾元素加 1。
func GetMissingElement(arr []int) int {
	if len(arr) <= 0 {
		return 1
	}
	sort.Sort(intSlice(arr))
	var arrTemp []int

	for i := 1; i <= arr[len(arr)-1]; i++ {
		arrTemp = append(arrTemp, i)
	}
	if equal(arrTemp, arr) {
		return arr[len(arr)-1] + 1
	}
	for i, v := range arrTemp {
		if v != arr[i] {
			return v
		}
	}
	return 1
}

// NewReplacer 将 endpoint 中的占位键按成对传入的 values 依次替换。
//
// values 必须按 键1, 值1, 键2, 值2……成对出现，数量为奇数时原样返回 endpoint；
// endpoint 为空或没有可替换内容时原样返回。
func NewReplacer(endpoint string, values ...any) string {
	if len(endpoint) < 1 {
		return endpoint
	}
	if len(values)%2 != 0 {
		return endpoint
	}
	params := make(map[string]string)
	if len(values) > 0 {
		key := ""
		for k, v := range values {
			if k%2 == 0 {
				key = cast.ToString(v)
			} else {
				params[key] = cast.ToString(v)
			}
		}
	}

	if len(params) < 1 || Empty(params) {
		return endpoint
	}

	for pk, pv := range params {
		endpoint = string(ByteReplace([]byte(endpoint), []byte(pk), []byte(pv), 1))
	}
	return endpoint
}

// ByteReplace 将字节切片 s 中前 n 处 old 替换为 new。
//
// n 为 0 时原样返回；n 为负数时替换全部；替换在原切片缓冲内进行，返回结果切片。
func ByteReplace(s, old, new []byte, n int) []byte {
	if n == 0 {
		return s
	}

	if len(old) < len(new) {
		return bytes.Replace(s, old, new, n)
	}

	if n < 0 {
		n = len(s)
	}

	var wid, i, j, w int
	for i, j = 0, 0; i < len(s) && j < n; j++ {
		wid = bytes.Index(s[i:], old)
		if wid < 0 {
			break
		}

		w += copy(s[w:], s[i:i+wid])
		w += copy(s[w:], new)
		i += wid + len(old)
	}

	w += copy(s[w:], s[i:])
	return s[0:w]
}
