package util

import (
	"crypto/rand"
	"math/big"
)

// RandomBytes random_bytes()
func RandomBytes(length int) ([]byte, error) {
	bs := make([]byte, length)
	_, err := rand.Read(bs)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

// RandomInt random_int()
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

// StrPad 填充字符串
func StrPad(str1, str2 string, i int) string {
	n := i - len(str1) - len(str2)
	if n > 0 {
		for i := 0; i < n; i++ {
			str1 = str1 + "0"
		}
	}
	return str1 + str2
}
