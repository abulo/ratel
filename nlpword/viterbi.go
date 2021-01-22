package nlpword

import (
	"sort"
)

const MIN_FLOAT = -3.14e100

var (
	prevStatus = make(map[byte][]byte)
	startProb  = make(map[byte]float64)
)

type probState struct {
	prob  float64
	state byte
}

type probStates []*probState

func init() {
	prevStatus['B'] = []byte{'E', 'S'}
	prevStatus['M'] = []byte{'M', 'B'}
	prevStatus['S'] = []byte{'S', 'E'}
	prevStatus['E'] = []byte{'B', 'M'}

	startProb['B'] = -0.26268660809250016
	startProb['E'] = -3.14e+100
	startProb['M'] = -3.14e+100
	startProb['S'] = -1.4652633398537678
}

func (ps probStates) Len() int {
	return len(ps)
}

func (ps probStates) Less(i, j int) bool {
	if ps[i].prob == ps[j].prob {
		return ps[i].state < ps[j].state
	}

	return ps[i].prob < ps[j].prob
}

func (ps probStates) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func Viterbi(obs []rune, states []byte) (float64, []byte) {
	path := make(map[byte][]byte)
	vtb := make([]map[byte]float64, len(obs))
	vtb[0] = make(map[byte]float64)

	for _, v := range states {
		if val, ok := probEmit[v][obs[0]]; ok {
			vtb[0][v] = val + startProb[v]
		} else {
			vtb[0][v] = val + MIN_FLOAT
		}

		path[v] = []byte{v}
	}

	n := 1
	for ; n < len(obs); n++ {
		newPath := make(map[byte][]byte)
		vtb[n] = make(map[byte]float64)

		for _, vv := range states {
			var emitP float64
			pss := make(probStates, 0)

			if val, ok := probEmit[vv][obs[n]]; ok {
				emitP = val
			} else {
				emitP = MIN_FLOAT
			}

			for _, ps := range prevStatus[vv] {
				var transP float64
				if tp, ok := probTrans[ps][vv]; ok {
					transP = tp
				} else {
					transP = MIN_FLOAT
				}

				prob := vtb[n-1][ps] + transP + emitP
				pss = append(pss, &probState{prob: prob, state: ps})
			}

			sort.Sort(sort.Reverse(pss))
			vtb[n][vv] = pss[0].prob

			pp := make([]byte, len(path[pss[0].state]))
			copy(pp, path[pss[0].state])
			newPath[vv] = append(pp, vv)
		}

		path = newPath
	}

	pss0 := make(probStates, 0)
	for _, s := range []byte{'E', 'S'} {
		pss0 = append(pss0, &probState{vtb[len(obs)-1][s], s})
	}
	sort.Sort(sort.Reverse(pss0))

	v := pss0[0]
	return v.prob, path[v.state]
}
