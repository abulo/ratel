package nlpword

import (
	"regexp"
)

type Hmm struct {
	regHan  *regexp.Regexp
	regSkip *regexp.Regexp
	hstates []byte //隐状态
}

func NewHmm() *Hmm {
	LoadEmitProb()

	return &Hmm{
		regHan:  regexp.MustCompile(`\p{Han}+`),
		regSkip: regexp.MustCompile(`(\d+\.\d+|[a-zA-Z0-9]+)`),
		hstates: []byte{'B', 'M', 'E', 'S'},
	}
}

func (hmm *Hmm) HmmPart(text string) []string {
	result := make([]string, 0, 10)

	var (
		hans      string
		nonHans   string
		hanLoc    []int
		nonHanLoc []int
	)

	for {
		//匹配汉字
		hanLoc = hmm.regHan.FindStringIndex(text)
		if hanLoc == nil {
			if len(text) == 0 {
				break
			}
		} else if hanLoc[0] == 0 {
			hans = text[hanLoc[0]:hanLoc[1]]
			text = text[hanLoc[1]:]
			for _, han := range hmm.GetViterbiResult(hans) {
				result = append(result, han)
			}
			continue
		}

		//匹配字母数字
		nonHanLoc = hmm.regSkip.FindStringIndex(text)
		if nonHanLoc == nil {
			if len(text) == 0 {
				break
			}
		} else if nonHanLoc[0] == 0 {
			nonHans = text[nonHanLoc[0]:nonHanLoc[1]]
			text = text[nonHanLoc[1]:]
			if nonHans != "" {
				result = append(result, nonHans)
				continue
			}
		}

		loc := locSwitch(text, hanLoc, nonHanLoc)
		if loc == nil {
			result = append(result, text)
			break
		}

		result = append(result, text[:loc[0]])
		text = text[loc[0]:]
	}

	return result
}

func (hmm *Hmm) GetViterbiResult(text string) []string {
	result := make([]string, 0, 10)

	begin, next := 0, 0
	runes := []rune(text)
	_, posList := Viterbi(runes, hmm.hstates)

	for i, v := range runes {
		pos := posList[i]
		switch pos {
		case 'B':
			begin = i
		case 'E':
			result = append(result, string(runes[begin:i+1]))
			next = i + 1
		case 'S':
			result = append(result, string(v))
			next = i + 1
		}
	}

	if next < len(runes) {
		result = append(result, string(runes[next:]))
	}

	return result
}

func locSwitch(text string, hanLoc, nonHanLoc []int) (loc []int) {
	if hanLoc == nil && nonHanLoc == nil {
		if len(text) > 0 {
			return nil
		}
	} else if hanLoc == nil {
		loc = nonHanLoc
	} else if nonHanLoc == nil || hanLoc[0] < nonHanLoc[0] {
		loc = hanLoc
	} else {
		loc = nonHanLoc
	}

	return loc
}
