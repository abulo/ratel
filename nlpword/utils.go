package nlpword

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

//将英文词转化为小写
func ToLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func FilterPartNil(part []*Part) []*Part {
	newPart := make([]*Part, 0, len(part))
	for _, item := range part {
		if item != nil {
			newPart = append(newPart, item)
		}
	}

	return newPart
}

//字符串显示分词 (老方法保留)
//tag 0 不显示词其他信息
//tag 1 显示词的词性
//tag 2 显示词的词频，distance
func PartToStrings(part []*Part, tag int) (output string) {
	if len(part) <= 0 {
		return ``
	}

	xchars := make([]byte, 0)
	for i := range part {
		for _, w := range part[i].Word() {
			if w.pos == "x" && w.chars[0] != SPACE {
				xchars = append(xchars, w.chars...)
			} else {
				if len(xchars) > 0 {
					if tag&1 == 1 {
						output += fmt.Sprintf("%s/%s| ", string(xchars), "x")
					} else if tag&2 == 2 {
						output += fmt.Sprintf("%s/%f| ", string(xchars), 0.)
					} else {
						output += fmt.Sprintf("%s|", string(xchars))
					}
					xchars = make([]byte, 0)
				}

				if tag&1 == 1 {
					output += fmt.Sprintf("%s/%s| ", string(w.chars), w.pos)
				} else if tag&2 == 2 {
					output += fmt.Sprintf("%s/%f|", string(w.chars), w.distance)
				} else {
					output += fmt.Sprintf("%s|", string(w.chars))
				}
			}
		}
	}
	if len(xchars) > 0 {
		if tag&1 == 1 {
			output += fmt.Sprintf("%s/%s| ", string(xchars), "x")
		} else if tag&2 == 2 {
			output += fmt.Sprintf("%s/%f| ", string(xchars), 0.)
		} else {
			output += fmt.Sprintf("%s|", string(xchars))
		}
		xchars = make([]byte, 0)
	}

	return output
}

//字符数组显示分词 (老方法保留)
func PartToTexts(part []*Part) (output []string) {
	if len(part) <= 0 {
		return nil
	}

	output = make([]string, 0)
	xchars := make([]byte, 0)
	for i := range part {
		for _, w := range part[i].Word() {
			if w.pos == "x" && w.chars[0] != SPACE {
				xchars = append(xchars, w.chars...)
			} else {
				if len(xchars) > 0 {
					output = append(output, string(xchars))
					xchars = make([]byte, 0)
				}

				output = append(output, string(w.chars))
			}
		}
	}
	if len(xchars) > 0 {
		output = append(output, string(xchars))
		xchars = make([]byte, 0)
	}

	return output
}

func GetKeysId(keywords ...string) (id uint64) {
	if len(keywords) == 0 {
		return 0
	}

	if len(keywords) == 1 {
		r := []byte(keywords[0])
		Md5Inst := md5.New()
		Md5Inst.Write(r)
		uret := Md5Inst.Sum([]byte(""))
		id = binary.BigEndian.Uint64(uret)
	} else if len(keywords) > 1 {
		for _, keyword := range keywords {
			r := []byte(keyword)
			Md5Inst := md5.New()
			Md5Inst.Write(r)
			uret := Md5Inst.Sum([]byte(""))
			id = binary.BigEndian.Uint64(uret)
		}
	}

	return
}
