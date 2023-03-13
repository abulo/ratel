package util

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

//////////// Mathematical Functions ////////////

// Abs abs()
func Abs(number float64) float64 {
	return math.Abs(number)
}

// Rand rand()
// Range: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// Round round()
func Round(value float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Trunc((value+0.5/p)*p) / p
}

// Floor floor()
func Floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil ceil()
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Pi pi()
func Pi() float64 {
	return math.Pi
}

// Max max()
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

// Min min()
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

// DecBin Decbin decbin()
func DecBin(number int64) string {
	return strconv.FormatInt(number, 2)
}

// BinDec Bindec bindec()
func BinDec(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 10), nil
}

// Hex2bin hex2bin()
func Hex2bin(data string) (string, error) {
	i, err := strconv.ParseInt(data, 16, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 2), nil
}

// Bin2hex bin2hex()
func Bin2hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		// If input is not binary number
		if err.(*strconv.NumError).Err == strconv.ErrSyntax {
			byteArray := []byte(str)
			var out string
			for i := 0; i < len(byteArray); i++ {
				out += strconv.FormatInt(int64(byteArray[i]), 16)
			}
			return out, nil
		} else {
			return "", err
		}
	}
	return strconv.FormatInt(i, 16), nil
}

// Dechex dechex()
func Dechex(number int64) string {
	return strconv.FormatInt(number, 16)
}

// Hexdec hexdec()
func Hexdec(str string) (int64, error) {
	return strconv.ParseInt(str, 16, 0)
}

// Decoct decoct()
func Decoct(number int64) string {
	return strconv.FormatInt(number, 8)
}

// Octdec Octdec()
func Octdec(str string) (int64, error) {
	return strconv.ParseInt(str, 8, 0)
}

// BaseConvert base_convert()
func BaseConvert(number string, frombase, tobase int) (string, error) {
	i, err := strconv.ParseInt(number, frombase, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, tobase), nil
}

// IsNan is_nan()
func IsNan(val float64) bool {
	return math.IsNaN(val)
}

// Divide 查询是否被整除
func Divide(m, n any) bool {
	r := math.Mod(cast.ToFloat64(m), cast.ToFloat64(n))
	return r == 0
}

// Add ...
func Add(m, n any) int {
	return cast.ToInt(m) + cast.ToInt(n)
}
