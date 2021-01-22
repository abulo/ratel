package nlpword

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MIN_WORD_FREQUENCY = 2
const SPACE = 32

const (
	PART_MODE_ONE   = 1 //普通分词
	PART_MODE_TWO   = 2 //mmseg分词
	PART_MODE_THREE = 3 //隐马尔可夫模型分词
)

//npartword分词器
type Parter struct {
	dict    *Dictionary    //分词字典
	emodict *EmoDictionary //情感词典
}

type jumper struct {
	word    *Word
	miniDis float32
}

func NewParter() *Parter {
	return &Parter{
		dict:    NewDictionary(),
		emodict: NewEmoDictionary(),
	}
}

func (pr *Parter) Dictionary() *Dictionary {
	return pr.dict
}

/**
 * 加载分词词典数据，dictfiles为字典文件路径，可以多个以,分隔
 *
 */
func (pr *Parter) LoadDictionary(dictfiles string) {
	for _, file := range strings.Split(dictfiles, ",") {
		log.Println("载入字典文件-", file)
		dictFile, err := os.Open(file)
		if err != nil {
			log.Fatalln("无法载入字典文件-", file)
		}

		var (
			text, ratetext, pos string
			rate                int
		)
		//读取文件，直到读完
		reader := bufio.NewReader(dictFile)
		for {
			size, _ := fmt.Fscanln(reader, &text, &ratetext, &pos)

			if size == 0 {
				//文件结束
				break
			} else if size < 2 {
				//无效行
				continue
			} else if size == 2 {
				//没有词性
				pos = ``
			}

			var err error
			rate, err = strconv.Atoi(ratetext)
			if err != nil {
				continue
			}

			//过滤词频过小的
			if rate < MIN_WORD_FREQUENCY {
				continue
			}

			//加载到字典树中
			word := NewWord([]byte(text), Rate(rate), Pos(pos))
			pr.dict.AddWord(word)
		}
	}

	//计算每个词的distance
	TotalDistance := float32(math.Log2(float64(pr.dict.TotalRate())))
	if len(pr.dict.words) > 0 {
		for i := range pr.dict.words {
			pr.dict.words[i].distance = TotalDistance - float32(math.Log2(float64(pr.dict.words[i].rate)))
		}
	}

	log.Println("加载字典文件完成")
}

/**
 * 加载情感词典数据，dictfiles为字典文件路径，可以多个以,分隔
 *
 */
func (pr *Parter) LoadEmoDictionary(dictfiles string) {
	for _, file := range strings.Split(dictfiles, ",") {
		log.Println("载入情感词典文件-", file)
		dictFile, err := os.Open(file)
		if err != nil {
			log.Fatalln("无法载入情感词典文件-", file)
		}

		var text, pos string
		reader := bufio.NewReader(dictFile)
		for {
			size, _ := fmt.Fscanln(reader, &text, &pos)

			if size == 0 {
				//文件结束
				break
			} else if size < 2 {
				//无效行
				continue
			} else if size == 1 {
				//没有词性
				pos = ``
			}

			//加载到情感词典btree中
			keyId := GetKeysId(text)

			//log
			//logFile, _ := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
			//loger := log.New(logFile, "npw", log.Ldate|log.Ltime|log.Lshortfile)
			//loger.Println(text, "====", keyId)

			pr.emodict.emowords.ReplaceOrInsert(&EmoWord{
				keyId: keyId,
				word:  text,
				pos:   pos,
			})
		}
	}

	log.Println("加载情感词典文件完成")
}

//执行分词
func (pr *Parter) DoPartWords(text string, partModel int) []*Part {
	chars := []rune(text)

	var outputPart []*Part

	//根据词disctance分词，精细度更高的分词
	if partModel == PART_MODE_ONE {
		var (
			current = 0
			partNum = 0
			jumpers = make([]jumper, len(chars))
		)
		for ; current < len(chars); current++ {
			var baseDis float32
			if current == 0 {
				baseDis = 0
			} else {
				baseDis = jumpers[current-1].miniDis
			}

			words, wnum := pr.dict.FindWord(chars[current:])
			if wnum > 0 {
				for wi := 0; wi < wnum; wi++ {
					wl := len([]rune(string(words[wi].chars)))
					location := current + wl - 1
					updateDis(&jumpers[location], baseDis, words[wi])
				}
			}
		}
		for index := len(chars) - 1; index >= 0; {
			if jumpers[index].word != nil {
				location := index - len([]rune(string(jumpers[index].word.chars))) + 1
				partNum++
				index = location - 1
			} else {
				index--
			}
		}
		outputPart = make([]*Part, partNum)
		for index := len(chars) - 1; index >= 0; {
			if jumpers[index].word != nil {
				location := index - len([]rune(string(jumpers[index].word.chars))) + 1
				partNum--

				part := NewPart(0, 0, []*Word{jumpers[index].word})
				outputPart[partNum] = part

				index = location - 1
			} else {
				index--
			}
		}
	}

	//mmseg方式分词，准确率更高的分词
	if partModel == PART_MODE_TWO {
		var (
			mmseg       *MMseg
			pos         = 0
			charsLength = len(chars)
		)
		for pos < charsLength {
			words := ``
			mmseg = NewMMsge(chars[pos:charsLength], pr.dict)
			mmseg.GetChunks()
			fwords := mmseg.FilterByMMsegRules()
			for _, w := range fwords {
				words += string(w.chars)
				part := NewPart(0, 0, []*Word{w})
				outputPart = append(outputPart, part)
			}

			pos += len([]rune(words))
		}
	}

	//隐马尔可夫模型分词，不走字典查询速度更快
	if partModel == PART_MODE_THREE {
		hmm := NewHmm()
		pstrs := hmm.HmmPart(text)
		for i := 0; i < len(pstrs); i++ {
			w := NewWord([]byte(pstrs[i]), Rate(0), Pos("hmm"))
			part := NewPart(0, 0, []*Word{w})
			outputPart = append(outputPart, part)
		}
	}

	return outputPart
}

//分词
func (pr *Parter) Part(text string, partModel int, tag int) *OutPut {
	words := pr.DoPartWords(text, partModel)

	return &OutPut{
		parter: pr,
		parts:  words,
		tag:    tag,
	}
}

//分词返回字符串，老方法暂时保留
func (pr *Parter) PartWords(text string, partModel int, tag int) string {
	words := pr.DoPartWords(text, partModel)

	return PartToStrings(words, tag)
}

//分词返回字符串数组，老方法暂时保留
func (pr *Parter) PartWordsTexts(text string, partModel int) []string {
	words := pr.DoPartWords(text, partModel)

	return PartToTexts(words)
}

type Text []byte

func splitToChars(text []byte) []Text {
	current := 0
	output := make([]Text, 0, len(text)/3)
	isAlphaNumeric := true
	alphaNumericStart := 0

	for current < len(text) {
		d, size := utf8.DecodeRune(text[current:])
		if size <= 2 && (unicode.IsLetter(d) || unicode.IsNumber(d)) { //数字、字母
			if !isAlphaNumeric {
				isAlphaNumeric = true
				alphaNumericStart = current
			}
		} else { //中文字符、日韩字符
			if isAlphaNumeric {
				isAlphaNumeric = false
				if current != 0 {
					//output = append(output, toLower(text[alphaNumericStart:current]))
					output = append(output, text[alphaNumericStart:current])
				}
			}
			output = append(output, text[current:current+size])
		}

		current += size
	}

	if isAlphaNumeric {
		if current != 0 {
			//output = append(output, toLower(text[alphaNumericStart:current]))
			output = append(output, text[alphaNumericStart:current])
		}
	}

	return output
}

func updateDis(jumper *jumper, baseDis float32, word *Word) {
	newDis := baseDis + word.distance
	if jumper.miniDis == 0 || jumper.miniDis > newDis {
		jumper.miniDis = newDis
		jumper.word = word
	}
}
