package nlpword

type Part struct {
	start int
	end   int
	word  []*Word
}

func NewPart(start, end int, word []*Word) *Part {
	return &Part{
		start: start,
		end:   end,
		word:  word,
	}
}

func (p *Part) Start() int {
	return p.start
}

func (p *Part) End() int {
	return p.end
}

func (p *Part) Word() []*Word {
	return p.word
}
