package util

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

// PERF-7-15 [P1] 进程级共享随机源：包级单例 + 互斥保护，并发安全；替代各函数每次 rand.New(rand.NewSource(...)) 的重复构造（预期收益需实测验证）
var (
	sharedRandMu sync.Mutex
	sharedRand   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// sharedIntn 返回 [0, n) 的随机整数，供包内各随机工具复用共享随机源。
func sharedIntn(n int) int {
	sharedRandMu.Lock()
	v := sharedRand.Intn(n)
	sharedRandMu.Unlock()
	return v
}

// sharedPerm 返回 [0, n) 的随机排列，供包内各随机工具复用共享随机源。
func sharedPerm(n int) []int {
	sharedRandMu.Lock()
	v := sharedRand.Perm(n)
	sharedRandMu.Unlock()
	return v
}

// sharedInt63n 返回 [0, n) 的非负随机 int64，供包内各随机工具复用共享随机源。
func sharedInt63n(n int64) int64 {
	sharedRandMu.Lock()
	v := sharedRand.Int63n(n)
	sharedRandMu.Unlock()
	return v
}

// Abs 返回 number 的绝对值。
func Abs(number float64) float64 {
	return math.Abs(number)
}

// Rand 返回区间 [min, max] 内的随机整数。
//
// min 必须小于等于 max，否则 panic；max 不得超过 2^31-1，否则 panic；
// min 等于 max 时直接返回 min。
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	// PERF-7-15 [P1] 复用进程级共享随机源；区间闭区间 [min, max] 语义与原实现保持一致
	return sharedIntn(max+1-min) + min
}

// Round 将 value 四舍五入到 precision 位小数。
func Round(value float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Trunc((value+0.5/p)*p) / p
}

// Floor 返回不大于 value 的最大整数。
func Floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil 返回不小于 value 的最小整数。
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pi 返回圆周率常量 math.Pi。
func Pi() float64 {
	return math.Pi
}

// Max 返回 nums 中的最大值；参数少于 2 个时 panic。
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max
}

// Min 返回 nums 中的最小值；参数少于 2 个时 panic。
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// DecBin 将整数转换为二进制字符串。
func DecBin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// BinDec 将二进制字符串解析为十进制数字字符串；解析失败时返回错误。
func BinDec(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 10), nil
}

// Hex2bin 将十六进制字符串转换为二进制字符串；解析失败时返回错误。
func Hex2bin(data string) (string, error) {
	i, err := strconv.ParseInt(data, 16, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 2), nil
}

// Bin2hex 将二进制字符串转换为十六进制字符串。
//
// 输入无法按二进制解析时，退化为按原始字节逐个转换为两位十六进制；其他解析错误
// 原样返回。
func Bin2hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		if err.(*strconv.NumError).Err == strconv.ErrSyntax {
			byteArray := []byte(str)
			var out strings.Builder
			for i := range byteArray {
				out.WriteString(strconv.FormatInt(int64(byteArray[i]), 16))
			}
			return out.String(), nil
		} else {
			return "", err
		}
	}
	return strconv.FormatInt(i, 16), nil
}

// Dechex 将整数转换为十六进制字符串。
func Dechex(number int64) string {
	return strconv.FormatInt(number, 16)
}

// Hexdec 将十六进制字符串解析为十进制整数；解析失败时返回错误。
func Hexdec(str string) (int64, error) {
	return strconv.ParseInt(str, 16, 0)
}

// Decoct 将整数转换为八进制字符串。
func Decoct(number int64) string {
	return strconv.FormatInt(number, 8)
}

// Octdec 将八进制字符串解析为十进制整数；解析失败时返回错误。
func Octdec(str string) (int64, error) {
	return strconv.ParseInt(str, 8, 0)
}

// BaseConvert 将 number 从 frombase 进制转换为 tobase 进制字符串；解析失败时返回错误。
func BaseConvert(number string, frombase, tobase int) (string, error) {
	i, err := strconv.ParseInt(number, frombase, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, tobase), nil
}

// IsNan 判断 value 是否为 NaN（非数字）。
func IsNan(val float64) bool {
	return math.IsNaN(val)
}

// Divide 判断 m 是否能被 n 整除；两者先转换为浮点值再取模，余数为 0 时返回 true。
func Divide(m, n any) bool {
	r := math.Mod(cast.ToFloat64(m), cast.ToFloat64(n))
	return r == 0
}

// Add 将 m 与 n 转换为整数后相加并返回结果。
func Add(m, n any) int {
	return cast.ToInt(m) + cast.ToInt(n)
}
