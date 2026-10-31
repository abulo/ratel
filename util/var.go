package util

import (
	"reflect"
	"strings"
)

// Empty 判断给定值是否为空。
//
// nil、长度为 0 的字符串/数组、长度为 0 或为 nil 的切片与 map、false、零值数字、
// nil 指针与接口均视为空；其余类型与零值做深度比较后判定。
func Empty(val any) bool {
	if val == nil {
		return true
	}
	// PERF-7-29 [P2] 标量与常见容器先走类型断言直比快路径，避免反射 ValueOf；结构体、指针、接口及其他具体类型仍落到下方反射兜底
	switch v := val.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	case bool:
		return !v
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// IsNumeric 判断值是否为数字类型或合法数字字符串。
//
// 整数值、浮点值与复数直接视为数字；字符串先去除首尾空白，允许前导正负号、小数
// 点与科学计数法（e/E），也支持 0x/0X 开头的十六进制写法；空字符串与其他类型
// 返回 false。
func IsNumeric(val any) bool {
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64, complex64, complex128:
		return true
	case string:
		if v == "" {
			return false
		}
		v = strings.TrimSpace(v)
		if v[0] == '-' || v[0] == '+' {
			if len(v) == 1 {
				return false
			}
			v = v[1:]
		}
		if len(v) > 2 && v[0] == '0' && (v[1] == 'x' || v[1] == 'X') {
			for _, h := range v[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		p, s, l := 0, 0, len(v)
		for i, ch := range v {
			if ch == '.' {
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if ch == 'e' || ch == 'E' {
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if ch < '0' || ch > '9' {
				return false
			}
		}
		return true
	}

	return false
}
