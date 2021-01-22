package nlpword

import (
	"math"
)

type MMseg struct {
	clen   int
	chars  []rune
	dict   *Dictionary
	chunks [][]*Word
}

func NewMMsge(chars []rune, dict *Dictionary) *MMseg {
	return &MMseg{
		clen:   len(chars),
		chars:  chars,
		dict:   dict,
		chunks: make([][]*Word, 0),
	}
}

func (mm *MMseg) GetChunks() (chunks [][]*Word) {
	chars := mm.chars
	charsLen := mm.clen
	if len(mm.chars) == 0 {
		return nil
	}

	//words1
	words1, wnum := mm.dict.FindWord(chars[0:])
	if wnum == 0 {
		return nil
	} else if wnum == 1 {
		chunks = append(chunks, words1)
		mm.chunks = append(mm.chunks, chunks...)
		return mm.chunks
	}

	for i := wnum - 1; i >= 0; i-- {
		//words1
		char1 := []rune(string(words1[i].chars))
		charl1 := len(char1)
		left1 := charsLen - charl1
		if left1 == 0 {
			chunks = append(chunks, words1[i:i+1])
			mm.chunks = append(mm.chunks, chunks...)
			return mm.chunks
		}

		//words2
		words2, wnum2 := mm.dict.FindWord(chars[charl1:])
		if wnum2 == 0 {
			return nil
		}
		for j := wnum2 - 1; j >= 0; j-- {
			char2 := []rune(string(words2[j].chars))
			charl2 := charl1 + len(char2)
			left2 := charsLen - charl2
			if left2 == 0 {
				c := []*Word{words1[i], words2[j]}
				chunks = append(chunks, c)
				continue
			}

			//words3
			words3, wnum3 := mm.dict.FindWord(chars[charl2:])
			if wnum3 == 0 {
				return
			}
			for k := wnum3 - 1; k >= 0; k-- {
				c := []*Word{words1[i], words2[j], words3[k]}
				chunks = append(chunks, c)
			}
		}
	}

	mm.chunks = append(mm.chunks, chunks...)

	return mm.chunks
}

func (mm *MMseg) FilterByMMsegRules() []*Word {
	var choice1, choice2, choice3, choice4 [][]*Word

	chunks := mm.chunks
	length := len(chunks)
	maxLength := 0

	//rule1: 取最大匹配的chunk
	for i := 0; i < length; i++ {
		var l int
		words := chunks[i]
		for j := 0; j < len(words); j++ {
			crune := []rune(string(words[j].chars))
			l += len(crune)
			if l > maxLength {
				maxLength = l
				choice1 = [][]*Word{chunks[i]}
			} else if l == maxLength {
				choice1 = append(choice1, chunks[i])
			}
		}
	}
	if len(choice1) == 1 {
		return choice1[0]
	}

	//rule2: 取平均词长最大的chunk
	avgLen := 0.
	for i := 0; i < len(choice1); i++ {
		avg := average(choice1[i])
		if avg > avgLen {
			avgLen = avg
			choice2 = [][]*Word{choice1[i]}
		} else if avg == avgLen {
			choice2 = append(choice2, choice1[i])
		}
	}
	if len(choice2) == 1 {
		return choice2[0]
	}

	//rule3: 取词长标准差最小的chunk
	smallestV := 65536. //large enough number
	for i := 0; i < len(choice2); i++ {
		v := variance(choice2[i])
		if v < smallestV {
			smallestV = v
			choice3 = [][]*Word{choice2[i]}
		} else if v == smallestV {
			choice3 = append(choice3, choice2[i])
		}
	}
	if len(choice3) == 1 {
		return choice3[0]
	}

	//rule4: 取单字词自由语素度之和最大的chunk
	smf := 0.
	for i := 0; i < len(choice3); i++ {
		dm := morphemicFreedom(choice3[i])
		if dm > smf {
			smf = dm
			choice4 = [][]*Word{choice3[i]}
		} else if dm == smf {
			choice4 = append(choice4, choice3[i])
		}
	}

	return choice4[0]
}

//求平均词长
func average(in []*Word) float64 {
	wnum := 0
	denominator := 0
	for j := 0; j < len(in); j++ {
		rword := []rune(string(in[j].chars))
		wnum += len(rword)
		denominator++
	}

	return float64(wnum) / float64(denominator)
}

//求词长标准差
func variance(in []*Word) float64 {
	avg := average(in)
	cumulative := 0.
	denominator := 0.

	for j := 0; j < len(in); j++ {
		rword := []rune(string(in[j].chars))
		v := float64(len(rword)) - avg
		cumulative += v * v
		denominator++
	}

	return math.Sqrt(cumulative / denominator)
}

//求单字词自由语素度之和
func morphemicFreedom(in []*Word) (out float64) {
	for i := 0; i < len(in); i++ {
		rword := []rune(string(in[i].chars))
		if 1 == len(rword) {
			//add offset 3 to prevent negative log value
			out += math.Log(float64(3 + in[i].distance))
		}
	}

	return out
}
