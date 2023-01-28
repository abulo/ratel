package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

// MarshalHTML ...
func MarshalHTML(v interface{}) template.HTML {
	a, _ := json.Marshal(v)
	return template.HTML(a)
}

// MarshalJS ...
func MarshalJS(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

// JS ...
func JS(v string) template.JS {
	return template.JS(v)
}

// Static ...
func Static(p, v string) string {
	return "/" + p + "/" + v
}

// DomainStatic ...
func DomainStatic(p, v string) string {
	return p + "/" + v
}

// UnescapeString ...
func UnescapeString(v interface{}) string {
	a, _ := json.Marshal(v)
	return html.EscapeString(string(a))
}

// GetAppRootPath 获取应用程序根目录
func GetAppRootPath() string {
	return GetParentDirectory(GetCurrentDirectory())
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// GetParentDirectory 获取上级目录
func GetParentDirectory(dir string) string {
	return substr(dir, 0, strings.LastIndex(dir, "/"))
}

// GetCurrentDirectory 获取当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Random 随机字符串
func Random() string {
	return uuid.New().String() + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// VerifyPhone 验证电话号码是否正确
func VerifyPhone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

// VerifyEmail 验证邮箱
func VerifyEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyIPv4 判断是否是否 ipv4 地址
func VerifyIPv4(address string) bool {
	res, _ := regexp.MatchString(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`, address)
	return res
}

// VerifyIPv6 判断是否是 ipv6 地址
func VerifyIPv6(address string) bool {
	res, _ := regexp.MatchString(`^\s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?\s*$`, address)
	return res
}

// FunctionName ...
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type intSlice []int

// Len ...
func (s intSlice) Len() int { return len(s) }

// Swap ...
func (s intSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less ...
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

// GetMissingElement ...
func GetMissingElement(arr []int) int {
	if arr == nil || len(arr) <= 0 {
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

// NewReplacer ...
func NewReplacer(endpoint string, values ...interface{}) string {
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
				key = fmt.Sprint(v)
			} else {
				params[key] = fmt.Sprint(v)
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

// ByteReplace ...
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
