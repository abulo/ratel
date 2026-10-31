package filter

import (
	"unicode/utf8"
)

// Search 是绑定到指定 TrieWriter 词库上的敏感词扫描器，通过 NewSearch、Strings
// 或 File 构造；不保证并发安全。
type Search struct {
	trieWriter *TrieWriter
}

// TrieWriter 返回扫描器关联的词库写入器。
func (_this *Search) TrieWriter() *TrieWriter {
	return _this.trieWriter
}

// decodeBytes 使用 UTF-8 解码字节切片首个 rune，返回该 rune 与其占用的字节数。
func decodeBytes(s []byte) (r rune, size int) {
	return utf8.DecodeRune(s)
}

// Find 在词库中扫描字节切片 s，返回全部命中的敏感词结果列表；无命中时返回空列表。
func (_this *Search) Find(s []byte) []*Result {
	list, _ := _this.findByAC(s, false)
	return list
}

// HasSens 判断字节切片 s 是否包含敏感词，命中即返回 true，不收集完整命中结果。
func (_this *Search) HasSens(s []byte) bool {
	// PERF-9-9 [P1]：直接消费 findByAC 返回的命中标志，single 模式不构造 Result 与字符串。
	_, hit := _this.findByAC(s, true)
	return hit
}

// Replace 将 s 中全部敏感词替换为单个 new 字节并返回新的字节切片；替换不改变字节宽度。
func (_this *Search) Replace(s []byte, new byte) []byte {
	data := make([]byte, len(s))
	start := 0
	for _, r := range _this.Find(s) {
		copy(data[start:], s[start:r.Start])
		start = r.Start
		// PERF-9-10 [P2] 直接在目标区间填充替换字节，去掉 repeatByte 临时切片分配与二次整块拷贝。
		for k := start; k <= r.End; k++ {
			data[k] = new
		}
		start = r.End + 1
	}
	copy(data[start:], s[start:])
	return data
}

// ReplaceRune 将 s 中全部敏感词替换为 new rune 并返回新的字节切片；
// 按 rune 的 UTF-8 长度调整输出切片长度，可能与输入长度不同。
func (_this *Search) ReplaceRune(s []byte, new rune) []byte {
	results := _this.Find(s)
	newLen := utf8.RuneLen(new)
	re := make([]int, len(results))
	addLen := 0
	for i, r := range results {
		rl := 0
		for range r.Matched {
			rl += newLen
		}
		re[i] = rl
		addLen += rl - (r.End - r.Start + 1)
	}

	data := make([]byte, len(s)+addLen)
	rBytes := []byte(string(new))
	start, add := 0, 0
	for i, r := range results {
		copy(data[start+add:], s[start:r.Start])
		start = r.Start
		al := re[i]
		// PERF-9-10 [P2] 直接在目标区间重复写入编码字节，去掉 repeatBytes 临时切片分配与二次拷贝。
		dst := data[start+add : start+add+al]
		for k := 0; k < al; k += len(rBytes) {
			copy(dst[k:], rBytes)
		}
		add += al - (r.End - r.Start + 1)
		start = r.End + 1
	}
	copy(data[start+add:], s[start:])
	return data
}

// findSub 在 s 的前 end 个字节中反向定位子串 sub 的起始字节位置。
//
// 备注（PERF-9-8 [P1]，未实施）：报告建议改为正向记录偏移以去掉反扫；在被跳过字符与
// 词条字符相同的重复片段上，反向贪心对齐与正向接受字符偏移可能定位到不同的起始位置，
// 无法证明两种定位在全部样本下结果一致，故保留原反向扫描实现，待补充等价样本与单独评审后再改。
func findSub(s []byte, end int, sub string) (idx int) {
	idx, j := end, len(sub)-1
	for ; j >= 0; idx-- {
		if s[idx] == sub[j] {
			j--
		}
	}
	idx++
	return
}

// findByAC 是 Aho-Corasick 扫描的核心实现，在词库自动机上扫描 s 并收集命中结果；
// single 为 true 时命中第一条结果即停止扫描并返回（命中标志 hit 置为 true）。
func (_this *Search) findByAC(s []byte, single bool) (list []*Result, hit bool) {
	n := len(s)
	trieRoot := _this.trieWriter.trie()
	skipper := _this.trieWriter.Skip()

	for i := 0; i < n; {
		v, l := decodeBytes(s[i:])
		node, ok := trieRoot.next[v]
		if !ok {
			i += l
			continue
		}

		j := i
		var (
			word []byte
			res  *Result
			skip int
		)
		for {
			word = append(word, s[j:j+l]...)
			if node.end {
				if single {
					// PERF-9-9 [P1] 存在性判定：命中即短路返回，不构造 Result、不转换字符串、不 append。
					return nil, true
				}
				sub := string(word[len(word)-int(node.len):])
				end := j + l - 1
				start := findSub(s, end, sub)
				res = &Result{sub, string(s[start : end+1]), start, end}
			}

			j += l
			for {
				v, l = decodeBytes(s[j:])
				if j < n && skipper.ShouldSkip(v) {
					j += l
					skip += l
				} else {
					break
				}
			}

			if next := node.next[v]; next != nil {
				node = next
			} else {
				if res == nil && node.fail.end {
					if single {
						// PERF-9-9 [P1] 失败后缀命中同样命中即短路。
						return nil, true
					}
					sub := string(word[len(word)-int(node.fail.len):])
					end := j - 1
					start := findSub(s, end, sub)
					res = &Result{sub, string(s[start : end+1]), start, end}
				}
				if res != nil {
					list = append(list, res)
					word = word[:0]
					if single {
						return list, true
					}
				}
				res = nil
				if node = node.fail.next[v]; node == nil {
					i = j
					break
				} else {
					i = j + l - int(node.len) - skip
				}
				break
			}
		}
	}
	return
}
