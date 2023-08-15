package filter

import (
	"unicode/utf8"
)

// Search 表示一个 tireRoot 树的搜索器
type Search struct {
	trieWriter *TrieWriter
}

// TrieWriter 返回关联的 TrieWriter
func (_this *Search) TrieWriter() *TrieWriter {
	return _this.trieWriter
}

// decodeBytes 使用 utf8 解码字节数组并返回 rune 和字节数
func decodeBytes(s []byte) (r rune, size int) {
	return utf8.DecodeRune(s)
}

// Find 在 tireRoot 树中搜索敏感词并将结果写入 w
func (_this *Search) Find(s []byte) []*Result {
	return _this.findByAC(s, false)
}

// HasSens 检查字节数组 s 是否包含敏感词
func (_this *Search) HasSens(s []byte) (has bool) {
	return len(_this.findByAC(s, true)) > 0
}

// Replace 将字节数组 s 中的所有敏感词替换为 new 并返回替换后的字节数组
func (_this *Search) Replace(s []byte, new byte) []byte {
	data := make([]byte, len(s))
	start := 0
	for _, r := range _this.Find(s) {
		copy(data[start:], s[start:r.Start])
		start = r.Start
		copy(data[start:], repeatByte(new, r.End-r.Start+1))
		start = r.End + 1
	}
	copy(data[start:], s[start:])
	return data
}

// ReplaceRune 将字节数组 s 中的所有敏感词替换为 new 并返回替换后的字节数组
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
		copy(data[start+add:], repeatBytes(rBytes, al))
		add += al - (r.End - r.Start + 1)
		start = r.End + 1
	}
	copy(data[start+add:], s[start:])
	return data
}

// repeatBytes 返回一个由 b 字节切片重复组成长度为 l 的字节数组
func repeatBytes(b []byte, l int) []byte {
	res := make([]byte, l)
	for s := 0; s < l; s += len(b) {
		copy(res[s:], b)
	}
	return res
}

// repeatByte 返回一个由 b 重复 n 次组成的字节数组
func repeatByte(b byte, n int) []byte {
	res := make([]byte, n)
	for i := range res {
		res[i] = b
	}
	return res
}

// findSub 在字节数组 s 的前 end 个字节中查找子串 sub 的开始位置
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

// findByAC 是 Aho-Corasick 算法实现的核心函数，用于在 tireRoot 树中搜索敏感词并返回结果
func (_this *Search) findByAC(s []byte, single bool) (list []*Result) {
	n := len(s)
	trieRoot := _this.trieWriter.trie() // 获取 trieRoot 树根节点
	skipper := _this.trieWriter.Skip()  // 获取跳过字符的规则

	for i := 0; i < n; {
		v, l := decodeBytes(s[i:])   // 解码以 s[i:] 开头的字节数组，并返回第一个 rune 和其所占字节数
		node, ok := trieRoot.next[v] // 在 trieRoot 中查找下一个节点
		if !ok {                     // 如果找不到，则继续向后搜索
			i += l
			continue
		}

		j := i
		var (
			word []byte  // 记录匹配到的字符串
			res  *Result // 记录匹配到的敏感词及相关信息
			skip int     // 记录跳过的无意义字符数量
		)
		for {
			word = append(word, s[j:j+l]...) // 将找到的字符加入 word 中
			if node.end {                    // 如果当前节点是一个单词的结尾，则说明找到了一个敏感词
				sub := string(word[len(word)-int(node.len):])            // 记录匹配的子串
				end := j + l - 1                                         // 敏感词在 s 中的结束位置
				start := findSub(s, end, sub)                            // 找到敏感词在 s 中的起始位置
				res = &Result{sub, string(s[start : end+1]), start, end} // 记录敏感词相关信息
			}

			j += l
			// 跳过一些无意义的字符
			for {
				v, l = decodeBytes(s[j:])
				if j < n && skipper.ShouldSkip(v) { // 如果当前字符是无意义字符，则跳过
					j += l
					skip += l
				} else {
					break
				}
			}

			// 如果找不到尝试fail指针
			if next := node.next[v]; next != nil {
				node = next
			} else {
				if res == nil && node.fail.end {
					sub := string(word[len(word)-int(node.fail.len):])       // 记录匹配的子串
					end := j - 1                                             // 敏感词在 s 中的结束位置
					start := findSub(s, end, sub)                            // 找到敏感词在 s 中的起始位置
					res = &Result{sub, string(s[start : end+1]), start, end} // 记录敏感词相关信息
				}
				if res != nil {
					list = append(list, res)
					word = word[:0]
					if single {
						return
					}
				}
				res = nil
				if node = node.fail.next[v]; node == nil {
					i = j
					break
				} else {
					i = j + l - int(node.len) - skip
				}
			}
		}
	}
	return
}
