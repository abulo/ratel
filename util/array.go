package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"github.com/spf13/cast"
)

// ArrayFill 生成以 startIndex 为起点、连续 num 个键均为 value 的 map[int]any。
func ArrayFill(startIndex int, num uint, value any) map[int]any {
	m := make(map[int]any)
	for range num {
		m[startIndex] = value
		startIndex++
	}
	return m
}

// ArrayFlip 交换 map 的键与值；键值冲突时后遍历到的值覆盖前者。
func ArrayFlip(m map[any]any) map[any]any {
	n := make(map[any]any)
	for i, v := range m {
		n[v] = i
	}
	return n
}

// ArrayKeys 返回 map 的全部键，返回顺序不保证。
func ArrayKeys(elements map[any]any) []any {
	i, keys := 0, make([]any, len(elements))
	for key := range elements {
		keys[i] = key
		i++
	}
	return keys
}

// ArrayValues 返回 map 的全部值，返回顺序不保证。
func ArrayValues(elements map[any]any) []any {
	i, vals := 0, make([]any, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}

// ArrayMerge 按参数顺序将多个切片拼接为一个新切片。
func ArrayMerge(ss ...[]any) []any {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]any, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// ArrayChunk 将切片按 size 大小切分为多个子切片；size 小于 1 时 panic。
func ArrayChunk(s []any, size int) [][]any {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]any
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = min((i+1)*size, length)
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

// ArrayPad 将切片填充到指定长度。
//
// size 为正时在尾部追加 val，为负时在头部前置 val；切片长度已满足要求时原样返回。
func ArrayPad(s []any, size int, val any) []any {
	if size == 0 || (size > 0 && size < len(s)) || (size < 0 && size > -len(s)) {
		return s
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(s)
	tmp := make([]any, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		return append(s, tmp...)
	}
	return append(tmp, s...)
}

// ArraySlice 从 offset 起截取 length 个元素；offset 超出切片长度时 panic，
// length 超出尾部时截到末尾。
func ArraySlice(s []any, offset, length uint) []any {
	if offset > uint(len(s)) {
		panic("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(s)) {
		return s[offset:end]
	}
	return s[offset:]
}

// ArrayRand 返回 elements 随机打乱顺序后的新切片，不修改原切片。
func ArrayRand(elements []any) []any {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]any, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// ArrayColumn 从 map 集合中提取每个元素 columnKey 字段的值组成列表；字段不存在的元素被跳过。
func ArrayColumn(input map[string]map[string]any, columnKey string) []any {
	columns := make([]any, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

// ArrayPush 将 elements 追加到切片尾部，返回追加后的长度。
//
// s 为切片指针，底层切片会被原地修改。
func ArrayPush(s *[]any, elements ...any) int {
	*s = append(*s, elements...)
	return len(*s)
}

// ArrayPop 移除并返回切片尾部元素；空切片返回 nil。
//
// s 为切片指针，底层切片会被原地修改。
func ArrayPop(s *[]any) any {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// ArrayUnshift 将 elements 插入切片头部，返回插入后的长度。
//
// s 为切片指针，底层切片会被原地修改。
func ArrayUnshift(s *[]any, elements ...any) int {
	*s = append(elements, *s...)
	return len(*s)
}

// ArrayShift 移除并返回切片头部元素；空切片返回 nil。
//
// s 为切片指针，底层切片会被原地修改。
func ArrayShift(s *[]any) any {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayKeyExists 判断 map m 中是否存在指定键 key。
func ArrayKeyExists(key any, m map[any]any) bool {
	_, ok := m[key]
	return ok
}

// ArrayCombine 以 s1 为键、s2 为值构造 map；两个切片长度不等时 panic。
func ArrayCombine(s1, s2 []any) map[any]any {
	if len(s1) != len(s2) {
		panic("the number of elements for each slice isn't equal")
	}
	m := make(map[any]any, len(s1))
	for i, v := range s1 {
		m[v] = s2[i]
	}
	return m
}

// ArrayReverse 原地反转切片并返回；原切片内容会被修改。
func ArrayReverse(s []any) []any {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Implode 使用 glue 连接字符串切片 pieces。
func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// InArray 判断 needle 是否存在于切片、数组或 map 中，元素比较使用 DeepEqual；
// haystack 类型不支持时 panic。
func InArray(needle any, haystack any) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

// ArrayRandMap 返回 map 切片随机打乱顺序后的新切片，不修改原切片。
func ArrayRandMap(elements []map[string]string) []map[string]string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]map[string]string, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// InMultiArray 判断 needle 中任一值是否出现在 haystack 的切片、数组或 map 中；
// haystack 类型不支持时 panic。
func InMultiArray(haystack any, needle ...any) bool {
	val := reflect.ValueOf(haystack)
	vals := reflect.ValueOf(needle)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			for j := 0; j < vals.Len(); j++ {
				if reflect.DeepEqual(vals.Index(j).Interface(), val.Index(i).Interface()) {
					return true
				}
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			for _, v := range vals.MapKeys() {
				if reflect.DeepEqual(vals.MapIndex(v).Interface(), val.MapIndex(k).Interface()) {
					return true
				}
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}
	return false
}

// MultiArray 判断切片或数组是否包含多于一条元素；其他类型返回 false。
func MultiArray(haystack any) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		if val.Len() > 1 {
			return true
		}
	default:
		return false
	}
	return false
}

// InterfaceToString 将 any 切片逐个格式化为字符串并返回字符串切片。
func InterfaceToString(data []any) (s []string) {
	for _, v := range data {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return s
}

// SplitString 将 p 中的标点与符号替换为分隔符 split。
//
// space 为 false 时同时去除文字中的空格；处理后对结果去重并以 split 重新拼接。
func SplitString(p, split string, space bool) string {
	var res string
	for _, c := range p {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += string(split)
		} else {
			if !space {
				if unicode.IsSpace(c) {
					continue
				} else {
					res += string(c)
				}
			} else {
				res += string(c)
			}
		}
	}
	reg := regexp.MustCompile(`\/{2,}`)
	res = reg.ReplaceAllString(res, split)
	args := Explode(split, res)
	newAry := make(map[string]string, 0)
	newArys := make(map[any]any, 0)
	for _, v := range args {
		if len(StrTrim(v)) > 0 {
			newAry[StrTrim(v)] = StrTrim(v)
		}
	}
	for _, v := range newAry {
		newArys[v] = v
	}
	data := InterfaceToString(ArrayValues(newArys))
	return Implode(split, data)
}

// ArrayPluck 从 map 切片中提取每个元素 value 字段的值组成字符串切片。
func ArrayPluck(data []map[string]string, value string) []string {
	res := make([]string, 0)
	for _, v := range data {
		res = append(res, v[value])
	}
	return res
}

// ArrayRemoveRepeatedElement 去除字符串切片中的重复元素，保留首次出现的顺序。
func ArrayRemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := range arr {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// ArrayKeyPluck 以每个元素 key 字段为键、value 字段为值构造 map。
func ArrayKeyPluck(data []map[string]string, value, key string) map[string]string {
	res := make(map[string]string, 0)
	for _, v := range data {
		res[v[key]] = v[value]
	}
	return res
}

// ArrayMultiPluck 以每个元素 key 字段为键、整个 map 为值构造嵌套 map。
func ArrayMultiPluck(data []map[string]string, key string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	for _, v := range data {
		res[v[key]] = v
	}
	return res
}

// jsonStringToObject 将 JSON 字符串反序列化到 v 指向的对象；解析失败返回错误。
func jsonStringToObject(s string, v any) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

// AryMapStringToAryMapInterface 将 []map[string]string 逐个转换为 []map[string]any。
func AryMapStringToAryMapInterface(d []map[string]string) []map[string]any {
	data := make([]map[string]any, 0)
	for _, v := range d {
		data = append(data, MapStringToMapInterface(v))
	}
	return data
}

// MapStringToMapInterface 将 map[string]string 的值转换为 any 类型。
func MapStringToMapInterface(d map[string]string) map[string]any {
	data := make(map[string]any, 0)

	for k, v := range d {
		data[k] = any(v)
	}
	return data
}

// AryMapInterfaceToAryMapString 将 []map[string]any 逐个转换为 []map[string]string，
// 值统一转为字符串。
func AryMapInterfaceToAryMapString(d []map[string]any) []map[string]string {
	data := make([]map[string]string, 0)
	for _, v := range d {
		data = append(data, MapInterfaceToMapString(v))
	}
	return data
}

// MapInterfaceToMapString 将 map[string]any 的值转换为字符串并返回新 map。
func MapInterfaceToMapString(d map[string]any) map[string]string {
	data := make(map[string]string, 0)

	for k, v := range d {
		data[k] = cast.ToString(v)
	}
	return data
}

// ArgStringToAryInterface 将字符串切片逐个包装为 any 类型组成新切片。
func ArgStringToAryInterface(d []string) []any {
	data := make([]any, 0)
	for _, v := range d {
		data = append(data, any(v))
	}
	return data
}

// AryInterfaceToArgString 将 any 切片元素逐个转为字符串组成新切片。
func AryInterfaceToArgString(d []any) []string {
	data := make([]string, 0)
	for _, v := range d {
		data = append(data, cast.ToString(v))
	}
	return data
}

// InterfaceToAryMapStringInterface 将 any 类型的切片元素转换为 []map[string]any；
// 传入值不是 []any 时会触发类型断言 panic。
func InterfaceToAryMapStringInterface(in any) []map[string]any {
	data := make([]map[string]any, 0)
	newData := in.([]any)
	for _, v := range newData {
		data = append(data, cast.ToStringMap(v))
	}
	return data
}

// InterfaceToAryMapStringString 将 any 类型的切片元素转换为 []map[string]string；
// 传入值不是 []any 时会触发类型断言 panic。
func InterfaceToAryMapStringString(in any) []map[string]string {
	data := make([]map[string]string, 0)
	newData := in.([]any)
	for _, v := range newData {
		data = append(data, cast.ToStringMapString(v))
	}
	return data
}

// ArrayStringUniq 去除字符串切片中的重复元素，保留首次出现的顺序。
func ArrayStringUniq(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := range arr {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
