package util

import (
	"crypto/rand"
	"math/big"
)

// RandomBytes 使用加密安全随机源生成 length 字节的随机切片。
//
// 随机源读取失败时返回 nil 与错误。
func RandomBytes(length int) ([]byte, error) {
	bs := make([]byte, length)
	_, err := rand.Read(bs)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

// RandomInt 使用加密安全随机源返回区间 [min, max] 内的随机整数。
//
// min 必须小于等于 max，否则 panic；min 等于 max 时直接返回 min。
func RandomInt(min, max int) (int, error) {
	if min > max {
		panic("argument #1 must be less than or equal to argument #2")
	}

	if min == max {
		return min, nil
	}
	nb, err := rand.Int(rand.Reader, big.NewInt(int64(max+1-min)))
	if err != nil {
		return 0, err
	}
	return int(nb.Int64()) + min, nil
}

// StrPad 在 str1 右侧补 0，使拼接 str2 后的总长度达到 i，最终返回补齐后的字符串。
func StrPad(str1, str2 string, i int) string {
	n := i - len(str1) - len(str2)
	if n > 0 {
		for range n {
			str1 = str1 + "0"
		}
	}
	return str1 + str2
}
