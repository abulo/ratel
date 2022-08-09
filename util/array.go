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

//////////// Array(Slice/Map) Functions ////////////

// ArrayFill array_fill()
func ArrayFill(startIndex int, num uint, value interface{}) map[int]interface{} {
	m := make(map[int]interface{})
	var i uint
	for i = 0; i < num; i++ {
		m[startIndex] = value
		startIndex++
	}
	return m
}

// ArrayFlip array_flip()
func ArrayFlip(m map[interface{}]interface{}) map[interface{}]interface{} {
	n := make(map[interface{}]interface{})
	for i, v := range m {
		n[v] = i
	}
	return n
}

// ArrayKeys array_keys()
func ArrayKeys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key := range elements {
		keys[i] = key
		i++
	}
	return keys
}

// ArrayValues array_values()
func ArrayValues(elements map[interface{}]interface{}) []interface{} {
	i, vals := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}

// ArrayMerge array_merge()
func ArrayMerge(ss ...[]interface{}) []interface{} {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]interface{}, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}

// ArrayChunk array_chunk()
func ArrayChunk(s []interface{}, size int) [][]interface{} {
	if size < 1 {
		panic("size: cannot be less than 1")
	}
	length := len(s)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, s[i*size:end])
		i++
	}
	return n
}

// ArrayPad array_pad()
func ArrayPad(s []interface{}, size int, val interface{}) []interface{} {
	if size == 0 || (size > 0 && size < len(s)) || (size < 0 && size > -len(s)) {
		return s
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(s)
	tmp := make([]interface{}, n)
	for i := 0; i < n; i++ {
		tmp[i] = val
	}
	if size > 0 {
		return append(s, tmp...)
	}
	return append(tmp, s...)
}

// ArraySlice array_slice()
func ArraySlice(s []interface{}, offset, length uint) []interface{} {
	if offset > uint(len(s)) {
		panic("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(s)) {
		return s[offset:end]
	}
	return s[offset:]
}

// ArrayRand array_rand()
func ArrayRand(elements []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]interface{}, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// ArrayColumn array_column()
func ArrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

// ArrayPush array_push()
// Push one or more elements onto the end of slice
func ArrayPush(s *[]interface{}, elements ...interface{}) int {
	*s = append(*s, elements...)
	return len(*s)
}

// ArrayPop array_pop()
// Pop the element off the end of slice
func ArrayPop(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	ep := len(*s) - 1
	e := (*s)[ep]
	*s = (*s)[:ep]
	return e
}

// ArrayUnshift array_unshift()
// Prepend one or more elements to the beginning of a slice
func ArrayUnshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// ArrayShift array_shift()
// Shift an element off the beginning of slice
func ArrayShift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayKeyExists array_key_exists()
func ArrayKeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}

// ArrayCombine array_combine()
func ArrayCombine(s1, s2 []interface{}) map[interface{}]interface{} {
	if len(s1) != len(s2) {
		panic("the number of elements for each slice isn't equal")
	}
	m := make(map[interface{}]interface{}, len(s1))
	for i, v := range s1 {
		m[v] = s2[i]
	}
	return m
}

// ArrayReverse array_reverse()
func ArrayReverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Implode implode()
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

// InArray in_array()
// haystack supported types: slice, array or map
func InArray(needle interface{}, haystack interface{}) bool {
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

// ArrayRandMap 随即
func ArrayRandMap(elements []map[string]string) []map[string]string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]map[string]string, len(elements))
	for i, v := range r.Perm(len(elements)) {
		n[i] = elements[v]
	}
	return n
}

// InMultiArray in_array()
// haystack supported types: slice, array or map
func InMultiArray(haystack interface{}, needle ...interface{}) bool {
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

// 判断二维数组里面是不是只有一条数据
func MultiArray(haystack interface{}) bool {
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

func InterfaceToString(data []interface{}) (s []string) {
	for _, v := range data {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return s
}

// 分割字符串 p 字符串, split 分隔符 , space 是否需要保留文字中的空格
func SplitString(p, split string, space bool) string {
	var res string
	for _, c := range p {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += string(split)
		} else {
			//如果不保留空格空格
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
	newArys := make(map[interface{}]interface{}, 0)
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

func ArrayPluck(data []map[string]string, value string) []string {
	res := make([]string, 0)
	for _, v := range data {
		res = append(res, v[value])
	}
	return res
}

// ArrayRemoveRepeatedElement 数组去重
func ArrayRemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
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

func ArrayKeyPluck(data []map[string]string, value, key string) map[string]string {
	res := make(map[string]string, 0)
	for _, v := range data {
		res[v[key]] = v[value]
	}
	return res
}

func ArrayMultiPluck(data []map[string]string, key string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	for _, v := range data {
		res[v[key]] = v
	}
	return res
}

// jsonStringToObject attempts to unmarshall a string as JSON into
// the object passed as pointer.
func jsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}

func AryMapStringToAryMapInterface(d []map[string]string) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	for _, v := range d {
		data = append(data, MapStringToMapInterface(v))
	}
	return data
}

func MapStringToMapInterface(d map[string]string) map[string]interface{} {
	data := make(map[string]interface{}, 0)

	for k, v := range d {
		data[k] = interface{}(v)
	}
	return data
}

func AryMapInterfaceToAryMapString(d []map[string]interface{}) []map[string]string {
	data := make([]map[string]string, 0)
	for _, v := range d {
		data = append(data, MapInterfaceToMapString(v))
	}
	return data
}

func MapInterfaceToMapString(d map[string]interface{}) map[string]string {
	data := make(map[string]string, 0)

	for k, v := range d {
		data[k] = cast.ToString(v)
	}
	return data
}

func ArgStringToAryInterface(d []string) []interface{} {
	data := make([]interface{}, 0)
	for _, v := range d {
		data = append(data, interface{}(v))
	}
	return data
}

func AryInterfaceToArgString(d []interface{}) []string {
	data := make([]string, 0)
	for _, v := range d {
		data = append(data, cast.ToString(v))
	}
	return data
}

func InterfaceToAryMapStringInterface(in interface{}) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	newData := in.([]interface{})
	for _, v := range newData {
		data = append(data, cast.ToStringMap(v))
	}
	return data
}

func InterfaceToAryMapStringString(in interface{}) []map[string]string {
	data := make([]map[string]string, 0)
	newData := in.([]interface{})
	for _, v := range newData {
		data = append(data, cast.ToStringMapString(v))
	}
	return data
}

// ArrayStringUniq 数组去重
func ArrayStringUniq(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
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
