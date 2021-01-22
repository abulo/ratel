package nlpword

import (
	"fmt"
	"strconv"
)

const (
	TAG_VAL_ZERO  = 0 //不显示词其他信息
	TAG_VAL_ONE   = 1 //显示词典中的词性
	TAG_VAL_TWO   = 2 //显示词的词频，distance
	TAG_VAL_THREE = 3 //显示词的情感词性，积极词(po)，否定词(ga)，消极词(ne)，程度副词(de)，主张词(cl)，无词性(x)。
)

type OutPut struct {
	parter *Parter
	parts  []*Part
	tag    int
	result string
}

func (op *OutPut) GetOutPut() *OutPut {
	return op
}

//字符串显示分词
func (op *OutPut) ToStrings() (output string) {
	if len(op.parts) <= 0 {
		return ``
	}

	xchars := make([]byte, 0)
	for i := range op.parts {
		for _, w := range op.parts[i].Word() {
			if w.pos == "x" && w.chars[0] != SPACE {
				xchars = append(xchars, w.chars...)
			} else {
				if len(xchars) > 0 {
					if op.tag == TAG_VAL_ONE {
						output += fmt.Sprintf("%s/%s| ", string(xchars), "x")
					} else if op.tag == TAG_VAL_TWO {
						output += fmt.Sprintf("%s/%f| ", string(xchars), 0.)
					} else if op.tag == TAG_VAL_THREE {
						output += op.EmotionLabel(string(xchars))
					} else {
						output += fmt.Sprintf("%s|", string(xchars))
					}
					xchars = make([]byte, 0)
				}

				if op.tag == TAG_VAL_ONE {
					output += fmt.Sprintf("%s/%s| ", string(w.chars), w.pos)
				} else if op.tag == TAG_VAL_TWO {
					output += fmt.Sprintf("%s/%f|", string(w.chars), w.distance)
				} else if op.tag == TAG_VAL_THREE {
					output += op.EmotionLabel(string(w.chars))
				} else {
					output += fmt.Sprintf("%s|", string(w.chars))
				}
			}
		}
	}
	if len(xchars) > 0 {
		if op.tag == TAG_VAL_ONE {
			output += fmt.Sprintf("%s/%s| ", string(xchars), "x")
		} else if op.tag == TAG_VAL_TWO {
			output += fmt.Sprintf("%s/%f| ", string(xchars), 0.)
		} else if op.tag == TAG_VAL_THREE {
			output += op.EmotionLabel(string(xchars))
		} else {
			output += fmt.Sprintf("%s|", string(xchars))
		}
		xchars = make([]byte, 0)
	}

	op.result = output

	return output
}

//字符数组显示分词
func (op *OutPut) ToTexts() (output []string) {
	if len(op.parts) <= 0 {
		return nil
	}

	output = make([]string, 0)
	xchars := make([]byte, 0)
	for i := range op.parts {
		for _, w := range op.parts[i].Word() {
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

//得到分词整句的情感值
type OpRet string

func (r OpRet) GetEmoScore() (output string) {
	if len(r) < 0 {
		return ``
	}

	score := NewEmotion(string(r)).DoScore().GetScore()
	scoreval := strconv.FormatFloat(score, 'f', 2, 64)

	return string(r) + "=" + scoreval
}

//给分词加情感词性标注
func (op *OutPut) EmotionLabel(xtext string) (output string) {
	keyId := GetKeysId(xtext)
	item := op.parter.emodict.emowords.Get(&EmoWord{keyId: keyId})

	if item != nil {
		record := item.(*EmoWord)
		output = fmt.Sprintf("%s/%s| ", record.word, record.pos)
	} else {
		output = fmt.Sprintf("%s/%s| ", xtext, "x")
	}

	return output
}
