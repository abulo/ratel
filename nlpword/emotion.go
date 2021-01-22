package nlpword

import (
	"strings"
)

const (
	EMO_VAL_DEFAULT = 0    //词情感值，默认值
	EMO_VAL_ONE     = 2    //词情感值，积极词的前一词为程度副词，
	EMO_VAL_TWO     = 1    //词情感值，积极词的前后词什么都不是，消极词的前一词为否定词
	EMO_VAL_THREE   = -0.5 //词情感值，否定词
	EMO_VAL_FOUR    = -1   //词情感值，积极词的前一词为否定词、消极词，后一词为消极词，消极词的前一次什么都不是
	EMO_VAL_FIVE    = -2   //词情感值，消极词的前一词为程度副词
)

type Emotion struct {
	sentence string
	score    float64
}

func NewEmotion(sentence string) *Emotion {
	return &Emotion{
		sentence: sentence,
		score:    0,
	}
}

func (e *Emotion) DoScore() *Emotion {
	if len(e.sentence) == 0 {
		e.score = 0
		return e
	}

	words := strings.Split(e.sentence, "|")
	if len(words) == 0 {
		e.score = 0
		return e
	}

	score := 0.
	for i := 0; i < len(words); i++ {
		word := words[i]

		//first
		if i == 0 {
			pos := getPos(word)
			ppos := "x"
			next := words[i+1]
			npos := getPos(next)

			score = filterScore(score, ppos, pos, npos)
		}

		if len(words) == 1 {
			break
		}

		if i != 0 && i != len(words)-1 && len(words) >= 3 {
			prev := words[i-1]
			ppos := getPos(prev)
			pos := getPos(word)
			next := words[i+1]
			npos := getPos(next)

			score = filterScore(score, ppos, pos, npos)
		}

		//last
		if i == len(words)-1 {
			pos := getPos(word)
			prev := words[i-1]
			ppos := getPos(prev)
			npos := "x"

			score = filterScore(score, ppos, pos, npos)
		}
	}

	e.score = score
	return e
}

func (e *Emotion) GetScore() float64 {
	return e.score
}

func getPos(word string) string {
	if len(word) == 0 {
		return "x"
	}

	texts := strings.Split(word, "/")
	if len(texts) != 2 {
		return "x"
	}

	return texts[1]
}

/**
 * 通过前一个词的词性，当前词性，后一个词性，过滤词性分值
 *
 */
func filterScore(score float64, ppos, cpos, npos string) float64 {
	//fmt.Println("ppos,cpos,npos=", ppos, cpos, npos)
	//fmt.Println("score1=", score)
	switch cpos {
	case "po":
		{
			if ppos == "de" {
				score += EMO_VAL_ONE
			} else if ppos == "ga" || ppos == "ne" || npos == "ne" {
				score += EMO_VAL_FOUR
			} else {
				score += EMO_VAL_TWO
			}
		}
	case "ga":
		{
			score += EMO_VAL_THREE
		}
	case "ne":
		{
			if ppos == "de" {
				score += EMO_VAL_FIVE
			} else if ppos == "ga" {
				score += EMO_VAL_TWO
			} else if ppos == "x" {
				score += EMO_VAL_FOUR
			}
		}
	case "x":
	default:
		{
			score += EMO_VAL_DEFAULT
		}
	}
	//fmt.Println("score2=", score)
	return score
}
