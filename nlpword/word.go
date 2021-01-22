package nlpword

type Word struct {
	chars    []byte
	rate     int
	pos      string
	distance float32
	parts    []*Part //该词的再次切分
}

type OptionalVal func(*Word) //可参数值

func NewWord(chars []byte, optvals ...OptionalVal) *Word {
	word := &Word{
		chars: chars,
	}

	for _, optval := range optvals {
		optval(word)
	}

	return word
}

func Rate(rate int) OptionalVal {
	return func(w *Word) {
		w.rate = rate
	}
}

func Pos(pos string) OptionalVal {
	return func(w *Word) {
		w.pos = pos
	}
}
