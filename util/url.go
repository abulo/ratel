package util

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// ParseURL 解析 URL 并按 component 位掩码返回各组成部分。
//
// 位掩码取值：1 scheme、2 host、4 port、8 user、16 pass、32 path、64 query、
// 128 fragment；component 为 -1 时返回全部八项。URL 非法时返回错误。
func ParseURL(str string, component int) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if component == -1 {
		component = 1 | 2 | 4 | 8 | 16 | 32 | 64 | 128
	}
	var components = make(map[string]string)
	if (component & 1) == 1 {
		components["scheme"] = u.Scheme
	}
	if (component & 2) == 2 {
		components["host"] = u.Hostname()
	}
	if (component & 4) == 4 {
		components["port"] = u.Port()
	}
	if (component & 8) == 8 {
		components["user"] = u.User.Username()
	}
	if (component & 16) == 16 {
		components["pass"], _ = u.User.Password()
	}
	if (component & 32) == 32 {
		components["path"] = u.Path
	}
	if (component & 64) == 64 {
		components["query"] = u.RawQuery
	}
	if (component & 128) == 128 {
		components["fragment"] = u.Fragment
	}
	return components, nil
}

// URLEncode 对字符串做查询参数编码，空格被编码为加号 "+"。
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

// URLDecode 解码查询参数字符串；编码非法时返回错误。
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// RawURLEncode 按 RFC 3986 做原始 URL 编码，空格被编码为 %20 而非加号。
func RawURLEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// RawURLDecode 解码原始 URL 编码字符串；编码非法时返回错误。
func RawURLDecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// HTTPBuildQuery 将 url.Values 编码为以 & 连接的查询字符串。
func HTTPBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// Base64Encode 返回字符串的标准 Base64 编码结果。
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode 解码标准 Base64 字符串；长度不是 4 的倍数时自动补齐等号，
// 输入非法时返回空字符串与错误。
func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
